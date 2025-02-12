package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	abci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

//
// =============================================================================
// BLOCK HANDLERS
// =============================================================================
//

// BeginBlocker calls telemetry and then tracks historical info.
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)
	return k.TrackHistoricalInfo(ctx)
}

// EndBlocker calls telemetry and then processes validator updates.
func (k *Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)
	return k.BlockValidatorUpdates(ctx)
}

//
// =============================================================================
// VOTE EXTENSION HANDLERS
// =============================================================================
//

// ExtendVoteHandler is called by all validators to extend the consensus vote
// with additional data to be voted on.
func (k *Keeper) ExtendVoteHandler(ctx context.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	oracleData, err := k.GetSidecarState(ctx, req.Height)
	if err != nil {
		k.Logger(ctx).Error("error retrieving AVS delegations", "height", req.Height, "error", err)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	voteExt, err := k.constructVoteExtension(ctx, req.Height, oracleData)
	if err != nil {
		k.Logger(ctx).Error("error creating vote extension", "height", req.Height, "error", err)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	if voteExt.IsInvalid(k.Logger(ctx)) {
		k.Logger(ctx).Error("invalid vote extension in ExtendVote", "height", req.Height)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	voteExtBz, err := json.Marshal(voteExt)
	if err != nil {
		k.Logger(ctx).Error("error marshalling vote extension", "height", req.Height, "error", err)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	return &abci.ResponseExtendVote{VoteExtension: voteExtBz}, nil
}

// constructVoteExtension builds the vote extension based on oracle data and on-chain state.
func (k *Keeper) constructVoteExtension(ctx context.Context, height int64, oracleData *OracleData) (VoteExtension, error) {
	avsDelegationsHash, err := deriveHash(oracleData.EigenDelegationsMap)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving AVS contract delegation state hash: %w", err)
	}

	ethBurnEventsHash, err := deriveHash(oracleData.EthBurnEvents)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving ethereum burn events hash: %w", err)
	}
	ethereumRedemptionsHash, err := deriveHash(oracleData.Redemptions)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving redemptions hash: %w", err)
	}

	neutrinoResponse, err := k.retrieveBitcoinHeader(ctx)
	if err != nil {
		return VoteExtension{}, err
	}
	bitcoinHeaderHash, err := deriveHash(neutrinoResponse.BlockHeader)
	if err != nil {
		return VoteExtension{}, err
	}

	nonces := make(map[uint64]uint64)
	for _, key := range k.getZenBTCKeyIDs(ctx) {
		requested, err := k.EthereumNonceRequested.Get(ctx, key)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				return VoteExtension{}, err
			}
			requested = false
		}
		if requested {
			nonce, err := k.lookupEthereumNonce(ctx, key)
			if err != nil {
				return VoteExtension{}, err
			}
			nonces[key] = nonce
		}
	}

	voteExt := VoteExtension{
		ZRChainBlockHeight:         height,
		ROCKUSDPrice:               oracleData.ROCKUSDPrice,
		BTCUSDPrice:                oracleData.BTCUSDPrice,
		ETHUSDPrice:                oracleData.ETHUSDPrice,
		EigenDelegationsHash:       avsDelegationsHash[:],
		EthBurnEventsHash:          ethBurnEventsHash[:],
		RedemptionsHash:            ethereumRedemptionsHash[:],
		BtcBlockHeight:             neutrinoResponse.BlockHeight,
		BtcHeaderHash:              bitcoinHeaderHash[:],
		EthBlockHeight:             oracleData.EthBlockHeight,
		EthGasLimit:                oracleData.EthGasLimit,
		EthBaseFee:                 oracleData.EthBaseFee,
		EthTipCap:                  oracleData.EthTipCap,
		SolanaLamportsPerSignature: oracleData.SolanaLamportsPerSignature,
		RequestedStakerNonce:       nonces[k.zenBTCKeeper.GetStakerKeyID(ctx)],
		RequestedEthMinterNonce:    nonces[k.zenBTCKeeper.GetEthMinterKeyID(ctx)],
		RequestedUnstakerNonce:     nonces[k.zenBTCKeeper.GetUnstakerKeyID(ctx)],
		RequestedCompleterNonce:    nonces[k.zenBTCKeeper.GetCompleterKeyID(ctx)],
	}

	return voteExt, nil
}

// VerifyVoteExtensionHandler is called by all validators to verify vote extension data.
func (k *Keeper) VerifyVoteExtensionHandler(ctx context.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	if len(req.VoteExtension) == 0 {
		return ACCEPT_VOTE, nil
	}

	if len(req.VoteExtension) > VoteExtBytesLimit {
		k.Logger(ctx).Error("vote extension is too large", "height", req.Height, "limit", VoteExtBytesLimit, "size", len(req.VoteExtension))
		return REJECT_VOTE, nil
	}

	var voteExt VoteExtension
	if err := json.Unmarshal(req.VoteExtension, &voteExt); err != nil {
		k.Logger(ctx).Debug("error unmarshalling vote extension", "height", req.Height, "error", err)
		return REJECT_VOTE, nil
	}

	if req.Height != voteExt.ZRChainBlockHeight {
		k.Logger(ctx).Error("mismatched height for vote extension", "expected", req.Height, "got", voteExt.ZRChainBlockHeight)
		return REJECT_VOTE, nil
	}

	if voteExt.IsInvalid(k.Logger(ctx)) {
		k.Logger(ctx).Error("invalid vote extension in VerifyVoteExtension", "height", req.Height)
		return REJECT_VOTE, nil
	}

	return ACCEPT_VOTE, nil
}

//
// =============================================================================
// PROPOSAL HANDLERS
// =============================================================================
//

// PrepareProposal is executed only by the proposer to inject oracle data into the block.
func (k *Keeper) PrepareProposal(ctx sdk.Context, req *abci.RequestPrepareProposal) ([]byte, error) {
	if !VoteExtensionsEnabled(ctx) {
		k.Logger(ctx).Debug("vote extensions disabled; not injecting oracle data", "height", req.Height)
		return nil, nil
	}

	voteExt, err := k.GetSuperMajorityVE(ctx, req.Height, req.LocalLastCommit)
	if err != nil {
		k.Logger(ctx).Error("error retrieving supermajority vote extension", "height", req.Height, "error", err)
		return nil, nil
	}

	if voteExt.ZRChainBlockHeight == 0 { // no supermajority vote extension
		return k.marshalOracleData(req, &OracleData{ConsensusData: req.LocalLastCommit})
	}

	if voteExt.ZRChainBlockHeight != req.Height-1 { // vote extension is from previous block
		k.Logger(ctx).Error("mismatched height for vote extension", "height", req.Height, "voteExt.ZRChainBlockHeight", voteExt.ZRChainBlockHeight)
		return nil, nil
	}

	oracleData, _, err := k.getValidatedOracleData(ctx, voteExt)
	if err != nil {
		k.Logger(ctx).Warn("error in getValidatedOracleData; injecting empty oracle data", "height", req.Height, "error", err)
		oracleData = &OracleData{}
	}
	oracleData.ConsensusData = req.LocalLastCommit

	return k.marshalOracleData(req, oracleData)
}

// ProcessProposal is executed by all validators to check whether the proposer prepared valid data.
func (k *Keeper) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	if !k.zrConfig.IsValidator {
		return ACCEPT_PROPOSAL, nil
	}

	if !VoteExtensionsEnabled(ctx) || len(req.Txs) == 0 {
		return ACCEPT_PROPOSAL, nil
	}

	if !ContainsVoteExtension(req.Txs[0], k.txDecoder) {
		k.Logger(ctx).Warn("block does not contain vote extensions, rejecting proposal")
		return REJECT_PROPOSAL, nil
	}

	var recoveredOracleData OracleData
	if err := json.Unmarshal(req.Txs[0], &recoveredOracleData); err != nil {
		return REJECT_PROPOSAL, fmt.Errorf("error unmarshalling oracle data: %w", err)
	}

	// Remove commit info before comparison.
	recoveredOracleDataNoCommitInfo := recoveredOracleData
	recoveredOracleDataNoCommitInfo.ConsensusData = abci.ExtendedCommitInfo{}
	if reflect.DeepEqual(recoveredOracleDataNoCommitInfo, OracleData{}) {
		k.Logger(ctx).Warn("accepting empty oracle data", "height", req.Height)
		return ACCEPT_PROPOSAL, nil
	}

	if err := ValidateVoteExtensions(ctx, k, req.Height, ctx.ChainID(), recoveredOracleData.ConsensusData); err != nil {
		k.Logger(ctx).Error("error validating vote extensions", "height", req.Height, "error", err)
		return REJECT_PROPOSAL, err
	}

	return ACCEPT_PROPOSAL, nil
}

//
// =============================================================================
// PRE-BLOCKER: ORACLE DATA PROCESSING
// =============================================================================
//

// PreBlocker is called before each block to process oracle data and update state.
func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	if !k.shouldProcessOracleData(ctx, req) {
		return nil
	}

	oracleData, ok := k.unmarshalOracleData(ctx, req.Txs[0])
	if !ok {
		return nil
	}

	voteExt, ok := k.validateCanonicalVE(ctx, req.Height, oracleData)
	if !ok {
		return nil
	}

	// Process various state updates.
	k.updateAssetPrices(ctx, oracleData)
	k.updateValidatorStakes(ctx, oracleData)
	k.updateAVSDelegationStore(ctx, oracleData)
	k.storeBitcoinBlockHeader(ctx, oracleData)
	k.storeNewZenBTCBurnEventsEthereum(ctx, oracleData)
	k.storeNewZenBTCRedemptions(ctx, oracleData)

	// Toggle minting and unstaking every other block due to a 1-block delay in processing VEs.
	if ctx.BlockHeight()%2 == 0 {
		k.updateNonces(ctx, oracleData)
		k.processZenBTCStaking(ctx, oracleData)
		k.processZenBTCMints(ctx, oracleData)
		k.processZenBTCBurnEventsEthereum(ctx, oracleData)
		k.processZenBTCRedemptions(ctx, oracleData)
	}

	k.recordNonVotingValidators(ctx, req)
	k.recordMismatchedVoteExtensions(ctx, req.Height, voteExt, oracleData.ConsensusData)
	return nil
}

// shouldProcessOracleData checks if oracle data should be processed for this block.
func (k *Keeper) shouldProcessOracleData(ctx sdk.Context, req *abci.RequestFinalizeBlock) bool {
	if len(req.Txs) == 0 {
		k.Logger(ctx).Debug("no transactions in block")
		return false
	}

	if req.Height == 1 || !VoteExtensionsEnabled(ctx) {
		k.Logger(ctx).Debug("vote extensions not enabled for this block", "height", req.Height)
		return false
	}

	if !ContainsVoteExtension(req.Txs[0], k.txDecoder) {
		k.Logger(ctx).Debug("first transaction does not contain vote extension", "height", req.Height)
		return false
	}

	return true
}

// validateCanonicalVE validates the proposed oracle data against the supermajority vote extension.
func (k *Keeper) validateCanonicalVE(ctx sdk.Context, height int64, oracleData OracleData) (VoteExtension, bool) {
	voteExt, err := k.GetSuperMajorityVE(ctx, height, oracleData.ConsensusData)
	if err != nil {
		k.Logger(ctx).Error("error retrieving supermajority vote extensions", "height", height, "error", err)
		return VoteExtension{}, false
	}

	if reflect.DeepEqual(voteExt, VoteExtension{}) {
		k.Logger(ctx).Warn("accepting empty vote extension", "height", height)
		return voteExt, true
	}

	if err := k.validateOracleData(voteExt, &oracleData); err != nil {
		k.Logger(ctx).Error("error validating oracle data; won't store VE data", "height", height, "error", err)
		return VoteExtension{}, false
	}

	return voteExt, true
}

// getValidatedOracleData retrieves and validates oracle data based on a vote extension.
func (k *Keeper) getValidatedOracleData(ctx sdk.Context, voteExt VoteExtension) (*OracleData, *VoteExtension, error) {
	oracleData, err := k.GetSidecarStateByEthHeight(ctx, voteExt.EthBlockHeight)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching oracle state: %w", err)
	}

	bitcoinData, err := k.sidecarClient.GetBitcoinBlockHeaderByHeight(
		ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{
			ChainName:   k.bitcoinNetwork(ctx),
			BlockHeight: voteExt.BtcBlockHeight,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching bitcoin header: %w", err)
	}

	oracleData.BtcBlockHeight = bitcoinData.BlockHeight
	oracleData.BtcBlockHeader = *bitcoinData.BlockHeader
	oracleData.RequestedStakerNonce = voteExt.RequestedStakerNonce
	oracleData.RequestedEthMinterNonce = voteExt.RequestedEthMinterNonce
	oracleData.RequestedUnstakerNonce = voteExt.RequestedUnstakerNonce
	oracleData.RequestedCompleterNonce = voteExt.RequestedCompleterNonce

	if err := k.validateOracleData(voteExt, oracleData); err != nil {
		return nil, nil, err
	}

	return oracleData, &voteExt, nil
}

//
// =============================================================================
// VALIDATOR & DELEGATION STATE UPDATES
// =============================================================================
//

// updateValidatorStakes updates validator stake values and delegation mappings.
func (k *Keeper) updateValidatorStakes(ctx sdk.Context, oracleData OracleData) {
	validatorInAVSDelegationSet := make(map[string]bool)

	for _, delegation := range oracleData.ValidatorDelegations {
		if delegation.Validator == "" {
			k.Logger(ctx).Debug("empty validator address in delegation; skipping")
			continue
		}

		valAddr, err := sdk.ValAddressFromBech32(delegation.Validator)
		if err != nil {
			k.Logger(ctx).Error("invalid validator address", "validator", delegation.Validator, "error", err)
			continue
		}

		validator, err := k.GetZenrockValidator(ctx, valAddr)
		if err != nil || validator.Status != types.Bonded {
			k.Logger(ctx).Debug("invalid delegation for", "validator", delegation.Validator, "error", err)
			continue
		}

		validator.TokensAVS = math.Int(delegation.Stake)

		if err = k.SetValidator(ctx, validator); err != nil {
			k.Logger(ctx).Error("error setting validator", "validator", delegation.Validator, "error", err)
			continue
		}

		if err = k.ValidatorDelegations.Set(ctx, valAddr.String(), delegation.Stake); err != nil {
			k.Logger(ctx).Error("error setting validator delegations", "validator", delegation.Validator, "error", err)
			continue
		}

		validatorInAVSDelegationSet[valAddr.String()] = true
	}

	k.removeStaleValidatorDelegations(ctx, validatorInAVSDelegationSet)
}

// removeStaleValidatorDelegations removes delegation entries for validators not present in the current AVS data.
func (k *Keeper) removeStaleValidatorDelegations(ctx sdk.Context, validatorInAVSDelegationSet map[string]bool) {
	var validatorsToRemove []string

	if err := k.ValidatorDelegations.Walk(ctx, nil, func(valAddr string, stake math.Int) (bool, error) {
		if !validatorInAVSDelegationSet[valAddr] {
			validatorsToRemove = append(validatorsToRemove, valAddr)
		}
		return true, nil
	}); err != nil {
		k.Logger(ctx).Error("error walking validator delegations", "error", err)
	}

	for _, valAddr := range validatorsToRemove {
		if err := k.ValidatorDelegations.Remove(ctx, valAddr); err != nil {
			k.Logger(ctx).Error("error removing validator delegation", "validator", valAddr, "error", err)
			continue
		}

		if err := k.updateValidatorTokensAVS(ctx, valAddr); err != nil {
			k.Logger(ctx).Error("error updating validator TokensAVS", "validator", valAddr, "error", err)
		}
	}
}

// updateValidatorTokensAVS resets a validator's AVS tokens to zero.
func (k *Keeper) updateValidatorTokensAVS(ctx sdk.Context, valAddr string) error {
	validator, err := k.GetZenrockValidator(ctx, sdk.ValAddress(valAddr))
	if err != nil {
		return fmt.Errorf("error retrieving validator for removal: %w", err)
	}

	validator.TokensAVS = math.ZeroInt()

	if err = k.SetValidator(ctx, validator); err != nil {
		return fmt.Errorf("error updating validator after removal: %w", err)
	}

	return nil
}

// updateAVSDelegationStore updates the AVS delegation store with new delegation amounts.
func (k *Keeper) updateAVSDelegationStore(ctx sdk.Context, oracleData OracleData) {
	for validatorAddr, delegatorMap := range oracleData.EigenDelegationsMap {
		for delegatorAddr, amount := range delegatorMap {
			if err := k.AVSDelegations.Set(ctx, collections.Join(validatorAddr, delegatorAddr), math.NewIntFromBigInt(amount)); err != nil {
				k.Logger(ctx).Error("error setting AVS delegations", "error", err)
			}
		}
	}
}

//
// =============================================================================
// BITCOIN HEADER PROCESSING
// =============================================================================
//

// storeBitcoinBlockHeader stores the Bitcoin header and handles historical header requests.
func (k *Keeper) storeBitcoinBlockHeader(ctx sdk.Context, oracleData OracleData) {
	if oracleData.BtcBlockHeight == 0 || oracleData.BtcBlockHeader.MerkleRoot == "" {
		k.Logger(ctx).Error("invalid bitcoin header data", "height", oracleData.BtcBlockHeight, "merkle", oracleData.BtcBlockHeader.MerkleRoot)
		return
	}

	requestedHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting requested historical Bitcoin headers", "error", err)
		}
		return
	}

	isHistorical := k.isHistoricalHeader(oracleData.BtcBlockHeight, requestedHeaders.Heights)
	headerPreviouslySeen, err := k.BtcBlockHeaders.Has(ctx, oracleData.BtcBlockHeight)
	if err != nil {
		k.Logger(ctx).Error("error checking if Bitcoin header is already stored", "height", oracleData.BtcBlockHeight, "error", err)
		return
	}

	if err := k.BtcBlockHeaders.Set(ctx, oracleData.BtcBlockHeight, oracleData.BtcBlockHeader); err != nil {
		k.Logger(ctx).Error("error storing Bitcoin header", "height", oracleData.BtcBlockHeight, "error", err)
		return
	}

	if isHistorical {
		requestedHeaders.Heights = slices.DeleteFunc(requestedHeaders.Heights, func(height int64) bool {
			return height == oracleData.BtcBlockHeight
		})
		if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
			k.Logger(ctx).Error("error updating requested historical Bitcoin headers", "error", err)
			return
		}
		k.Logger(ctx).Debug("stored historical Bitcoin header and removed request", "height", oracleData.BtcBlockHeight, "remaining_requests", len(requestedHeaders.Heights))
		return
	}

	if headerPreviouslySeen {
		k.Logger(ctx).Debug("bitcoin header previously seen; skipping reorg check", "height", oracleData.BtcBlockHeight)
		return
	}

	if err := k.checkForBitcoinReorg(ctx, oracleData, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error handling potential Bitcoin reorg", "height", oracleData.BtcBlockHeight, "error", err)
	}
}

// isHistoricalHeader checks if the given Bitcoin block height is in the list of requested historical headers.
func (k *Keeper) isHistoricalHeader(height int64, requestedHeights []int64) bool {
	for _, h := range requestedHeights {
		if h == height {
			return true
		}
	}
	return false
}

// checkForBitcoinReorg detects reorgs by requesting previous headers when a new one is received.
func (k *Keeper) checkForBitcoinReorg(ctx sdk.Context, oracleData OracleData, requestedHeaders zenbtctypes.RequestedBitcoinHeaders) error {
	var numHistoricalHeadersToRequest int64 = 20
	if strings.HasPrefix(ctx.ChainID(), "diamond") {
		numHistoricalHeadersToRequest = 6
	}

	prevHeights := make([]int64, 0, numHistoricalHeadersToRequest)
	for i := int64(1); i <= numHistoricalHeadersToRequest; i++ {
		prevHeight := oracleData.BtcBlockHeight - i
		if prevHeight <= 0 {
			break
		}
		prevHeights = append(prevHeights, prevHeight)
	}

	if len(prevHeights) == 0 {
		k.Logger(ctx).Error("no previous heights to request (this should not happen with a valid VE)", "height", oracleData.BtcBlockHeight)
		return nil
	}

	requestedHeaders.Heights = append(requestedHeaders.Heights, prevHeights...)
	if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error setting requested historical Bitcoin headers", "error", err)
		return err
	}

	return nil
}

//
// =============================================================================
// ZENBTC PROCESSING: STAKING, MINTING, BURN EVENTS & REDEMPTIONS
// =============================================================================
//

// checkForUpdateAndDispatchTx is a generic helper that encapsulates the common logic for tx dispatch and nonce management.
// It receives a slice of pending transactions (or redemptions, or burn events) and passes the
// slice to nonceUpdatedCallback. If a second tx exists, it passes that tx to txDispatchCallback.
func checkForUpdateAndDispatchTx[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedNonce uint64,
	pendingTxs []T,
	nonceUpdatedCallback func(tx T) error,
	txDispatchCallback func(tx T) error,
) {
	if len(pendingTxs) == 0 {
		return
	}

	nonceData, err := k.LastUsedEthereumNonce.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting last used ethereum nonce", "keyID", keyID, "error", err)
			return
		}
		nonceData = zenbtctypes.NonceData{Nonce: 0, PrevNonce: 0, Counter: 0, Skip: true}
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			k.Logger(ctx).Error("error setting last used ethereum nonce", "keyID", keyID, "error", err)
		}
		return
	}
	k.Logger(ctx).Info("Nonce info", "nonce", nonceData.Nonce, "prev", nonceData.PrevNonce, "counter", nonceData.Counter, "skip", nonceData.Skip, "requested", requestedNonce)

	if nonceData.Nonce != 0 && requestedNonce == 0 {
		return
	}

	requested, err := k.EthereumNonceRequested.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting ethereum nonce request state", "keyID", keyID, "error", err)
		}
		return
	}
	if !requested {
		return
	}

	nonceUpdated := false
	if requestedNonce != nonceData.PrevNonce {
		if err := nonceUpdatedCallback(pendingTxs[0]); err != nil {
			k.Logger(ctx).Error("nonce update callback error", "keyID", keyID, "error", err)
		}
		k.Logger(ctx).Warn("nonce updated for key", "keyID", keyID, "requestedNonce", requestedNonce, "prevNonce", nonceData.PrevNonce, "currentNonce", nonceData.Nonce)
		nonceData.PrevNonce = nonceData.Nonce
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			k.Logger(ctx).Error("error setting last used Ethereum nonce", "keyID", keyID, "error", err)
		}
		if len(pendingTxs) == 1 {
			if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
				k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
			}
			return
		}
		nonceUpdated = true
	}
	if nonceData.Skip {
		return
	}

	txIndex := 0
	// If tx[0] confirmed on-chain via nonce increment, dispatch tx[1]. If not then retry dispatching tx[0].
	if nonceUpdated {
		txIndex = 1
	}
	if err := txDispatchCallback(pendingTxs[txIndex]); err != nil {
		k.Logger(ctx).Error("tx dispatch callback error", "keyID", keyID, "error", err)
	}
}

// updateNonces handles updating nonce state for keys used for minting and unstaking.
func (k *Keeper) updateNonces(ctx sdk.Context, oracleData OracleData) {
	for _, key := range k.getZenBTCKeyIDs(ctx) {
		requested, err := k.EthereumNonceRequested.Get(ctx, key)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error("error checking nonce request state", "keyID", key, "error", err)
			}
			continue
		}
		if !requested {
			continue
		}

		var currentNonce uint64
		switch key {
		case k.zenBTCKeeper.GetStakerKeyID(ctx):
			currentNonce = oracleData.RequestedStakerNonce
		case k.zenBTCKeeper.GetEthMinterKeyID(ctx):
			currentNonce = oracleData.RequestedEthMinterNonce
		case k.zenBTCKeeper.GetUnstakerKeyID(ctx):
			currentNonce = oracleData.RequestedUnstakerNonce
		case k.zenBTCKeeper.GetCompleterKeyID(ctx):
			currentNonce = oracleData.RequestedCompleterNonce
		default:
			k.Logger(ctx).Error("invalid key ID", "keyID", key)
			continue
		}

		// Avoid erroneously setting nonce to zero if a non-zero nonce exists i.e. blocks with no consensus on VEs.
		nonceData, err := k.LastUsedEthereumNonce.Get(ctx, key)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error("error getting last used ethereum nonce", "keyID", key, "error", err)
				continue
			}
			nonceData = zenbtctypes.NonceData{Nonce: 0, PrevNonce: 0, Counter: 0, Skip: true}
			if err := k.LastUsedEthereumNonce.Set(ctx, key, nonceData); err != nil {
				k.Logger(ctx).Error("error setting last used ethereum nonce", "keyID", key, "error", err)
			}
			continue
		}
		if nonceData.Nonce != 0 && currentNonce == 0 {
			continue
		}

		if err := k.updateNonceState(ctx, key, currentNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "keyID", key, "error", err)
		}
	}
}

// clearEthereumNonceRequest resets the nonce-request flag for a given key.
func (k *Keeper) clearEthereumNonceRequest(ctx sdk.Context, keyID uint64) error {
	k.Logger(ctx).Warn("set EthereumNonceRequested state to false", "keyID", keyID)
	return k.EthereumNonceRequested.Set(ctx, keyID, false)
}

// getPendingMintTransactionsByStatus retrieves up to 2 pending mint transactions matching the given status.
func (k *Keeper) getPendingMintTransactionsByStatus(ctx sdk.Context, status zenbtctypes.MintTransactionStatus) ([]zenbtctypes.PendingMintTransaction, error) {
	var txs []zenbtctypes.PendingMintTransaction
	if err := k.zenBTCKeeper.WalkPendingMintTransactions(ctx, func(id uint64, pending zenbtctypes.PendingMintTransaction) (bool, error) {
		if pending.Status == status {
			txs = append(txs, pending)
			if len(txs) == 2 {
				return true, nil
			}
		}
		return false, nil
	}); err != nil {
		return nil, fmt.Errorf("error walking pending mint transactions: %w", err)
	}
	return txs, nil
}

// processZenBTCStaking processes pending staking transactions.
func (k *Keeper) processZenBTCStaking(ctx sdk.Context, oracleData OracleData) {
	keyID := k.zenBTCKeeper.GetStakerKeyID(ctx)
	requestedNonce := oracleData.RequestedStakerNonce

	depositedTxs, err := k.getPendingMintTransactionsByStatus(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED)
	if err != nil {
		k.Logger(ctx).Error("error getting pending mint transactions", "error", err)
		return
	}
	if len(depositedTxs) == 0 {
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		return
	}

	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedNonce, depositedTxs,
		func(pendingMint zenbtctypes.PendingMintTransaction) error {
			pendingMint.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED
			if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, pendingMint); err != nil {
				return err
			}
			if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), true); err != nil {
				return err
			}
			return nil
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			k.Logger(ctx).Warn("processing zenBTC stake",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", requestedNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)
			unsignedStakeTxHash, unsignedStakeTx, err := k.constructStakeTx(
				ctx,
				getChainIDForEigen(ctx),
				tx.Amount,
				requestedNonce,
				oracleData.EthGasLimit,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
				ctx,
				&treasurytypes.MsgNewSignTransactionRequest{
					Creator:             tx.Creator,
					KeyId:               keyID,
					WalletType:          treasurytypes.WalletType(tx.ChainType),
					UnsignedTransaction: unsignedStakeTx,
					Metadata:            metadata,
					NoBroadcast:         false,
				},
				[]byte(hex.EncodeToString(unsignedStakeTxHash)),
			)
			return err
		},
	)
}

// processZenBTCMints processes pending mint transactions.
func (k *Keeper) processZenBTCMints(ctx sdk.Context, oracleData OracleData) {
	keyID := k.zenBTCKeeper.GetEthMinterKeyID(ctx)
	requestedNonce := oracleData.RequestedEthMinterNonce

	stakedTxs, err := k.getPendingMintTransactionsByStatus(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED)
	if err != nil {
		k.Logger(ctx).Error("error getting pending mint transactions", "error", err)
		return
	}
	if len(stakedTxs) == 0 {
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		return
	}

	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedNonce, stakedTxs,
		func(mint zenbtctypes.PendingMintTransaction) error {
			supply, err := k.zenBTCKeeper.GetSupply(ctx)
			if err != nil {
				return err
			}
			supply.PendingZenBTC -= mint.Amount
			supply.MintedZenBTC += mint.Amount
			if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
				return err
			}
			k.Logger(ctx).Warn("pending mint supply updated",
				"pending_mint_old", supply.PendingZenBTC+mint.Amount,
				"pending_mint_new", supply.PendingZenBTC)
			k.Logger(ctx).Warn("minted supply updated",
				"minted_old", supply.MintedZenBTC-mint.Amount,
				"minted_new", supply.MintedZenBTC)
			mint.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
			if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, mint); err != nil {
				return err
			}
			return nil
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
			if err != nil {
				return err
			}
			feeZenBTC := k.CalculateZenBTCMintFee(
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
				oracleData.EthGasLimit,
				oracleData.BTCUSDPrice,
				oracleData.ETHUSDPrice,
				exchangeRate,
			)
			if oracleData.BTCUSDPrice.IsZero() {
				return nil
			}
			if tx.Caip2ChainId != "eip155:17000" {
				return fmt.Errorf("invalid chain ID: %s", tx.Caip2ChainId)
			}
			chainID, err := types.ExtractEVMChainID(tx.Caip2ChainId)
			if err != nil {
				return err
			}
			unsignedMintTxHash, unsignedMintTx, err := k.constructMintTx(
				ctx,
				tx.RecipientAddress,
				chainID,
				tx.Amount,
				feeZenBTC,
				requestedNonce,
				oracleData.EthGasLimit,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: chainID})
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
				ctx,
				&treasurytypes.MsgNewSignTransactionRequest{
					Creator:             tx.Creator,
					KeyId:               keyID,
					WalletType:          treasurytypes.WalletType(tx.ChainType),
					UnsignedTransaction: unsignedMintTx,
					Metadata:            metadata,
					NoBroadcast:         false,
				},
				[]byte(hex.EncodeToString(unsignedMintTxHash)),
			)
			return err
		},
	)
}

// storeNewZenBTCBurnEventsEthereum stores new burn events coming from Ethereum.
func (k *Keeper) storeNewZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	foundNewBurn := false
	// Loop over each burn event from oracle to check for new ones.
	for _, burn := range oracleData.EthBurnEvents {
		// Check if this burn event already exists
		exists := false
		if err := k.zenBTCKeeper.WalkBurnEvents(ctx, func(id uint64, existingBurn zenbtctypes.BurnEvent) (bool, error) {
			if existingBurn.TxID == burn.TxID &&
				existingBurn.LogIndex == burn.LogIndex &&
				existingBurn.ChainID == burn.ChainID {
				exists = true
				return true, nil
			}
			return false, nil
		}); err != nil {
			k.Logger(ctx).Error("error walking burn events", "error", err)
			continue
		}

		if !exists {
			newBurn := zenbtctypes.BurnEvent{
				TxID:            burn.TxID,
				LogIndex:        burn.LogIndex,
				ChainID:         burn.ChainID,
				DestinationAddr: burn.DestinationAddr,
				Amount:          burn.Amount,
				Status:          zenbtctypes.BurnStatus_BURN_STATUS_BURNED,
			}
			id, err := k.zenBTCKeeper.CreateBurnEvent(ctx, &newBurn)
			if err != nil {
				k.Logger(ctx).Error("error creating burn event", "error", err)
				continue
			}
			k.Logger(ctx).Info("created new burn event", "id", id)
			foundNewBurn = true
		}
	}

	if foundNewBurn {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "error", err)
		}
	}
}

// processZenBTCBurnEventsEthereum processes pending burn events by constructing unstake transactions.
func (k *Keeper) processZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	keyID := k.zenBTCKeeper.GetUnstakerKeyID(ctx)
	requestedNonce := oracleData.RequestedUnstakerNonce

	var burnedEvents []zenbtctypes.BurnEvent
	if err := k.zenBTCKeeper.WalkBurnEvents(ctx, func(id uint64, burnEvent zenbtctypes.BurnEvent) (bool, error) {
		if burnEvent.Status == zenbtctypes.BurnStatus_BURN_STATUS_BURNED {
			burnedEvents = append(burnedEvents, burnEvent)
			if len(burnedEvents) == 2 {
				return true, nil
			}
		}
		return false, nil
	}); err != nil {
		k.Logger(ctx).Error("error walking burn events", "error", err)
		return
	}
	if len(burnedEvents) == 0 {
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		return
	}

	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedNonce, burnedEvents,
		func(burnEvent zenbtctypes.BurnEvent) error {
			burnEvent.Status = zenbtctypes.BurnStatus_BURN_STATUS_UNSTAKING
			if err := k.zenBTCKeeper.SetBurnEvent(ctx, burnEvent.Id, burnEvent); err != nil {
				return err
			}
			return nil
		},
		func(burnEvent zenbtctypes.BurnEvent) error {
			k.Logger(ctx).Warn("processing zenBTC burn unstake",
				"burn_event", burnEvent,
				"nonce", requestedNonce,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap)
			unsignedTxHash, unsignedTx, err := k.constructUnstakeTx(
				ctx,
				getChainIDForEigen(ctx),
				burnEvent.DestinationAddr,
				burnEvent.Amount,
				requestedNonce,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
			if err != nil {
				return err
			}
			creator, err := k.getAddressByKeyID(ctx, keyID, treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
				ctx,
				&treasurytypes.MsgNewSignTransactionRequest{
					Creator:             creator,
					KeyId:               keyID,
					WalletType:          treasurytypes.WalletType_WALLET_TYPE_EVM,
					UnsignedTransaction: unsignedTx,
					Metadata:            metadata,
					NoBroadcast:         false,
				},
				[]byte(hex.EncodeToString(unsignedTxHash)),
			)
			return err
		},
	)
}

// storeNewZenBTCRedemptions processes new redemption events.
func (k *Keeper) storeNewZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	// Find the first INITIATED redemption.
	var firstInitiatedRedemption zenbtctypes.Redemption
	var found bool

	if err := k.zenBTCKeeper.WalkRedemptions(ctx, func(id uint64, r zenbtctypes.Redemption) (bool, error) {
		if r.Status == zenbtctypes.RedemptionStatus_INITIATED {
			firstInitiatedRedemption = r
			found = true
			return true, nil
		}
		return false, nil
	}); err != nil {
		k.Logger(ctx).Error("error finding first initiated redemption", "error", err)
		return
	}

	// If an INITIATED redemption is found, check if it exists in oracleData.
	if found {
		redemptionExists := false
		for _, redemption := range oracleData.Redemptions {
			if redemption.Id == firstInitiatedRedemption.Data.Id {
				redemptionExists = true
				break
			}
		}
		// If not present, mark it as unstaked.
		if !redemptionExists {
			firstInitiatedRedemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
			if err := k.zenBTCKeeper.SetRedemption(ctx, firstInitiatedRedemption.Data.Id, firstInitiatedRedemption); err != nil {
				k.Logger(ctx).Error("error updating redemption status to unstaked", "error", err)
				return
			}
		}
	}

	if len(oracleData.Redemptions) == 0 {
		return
	}

	exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting zenBTC exchange rate", "error", err)
		return
	}

	foundNewRedemption := false

	for _, redemption := range oracleData.Redemptions {
		redemptionExists, err := k.zenBTCKeeper.HasRedemption(ctx, redemption.Id)
		if err != nil {
			k.Logger(ctx).Error("error checking redemption existence", "error", err)
			continue
		}
		if redemptionExists {
			k.Logger(ctx).Debug("redemption already stored", "id", redemption.Id)
			continue
		}

		foundNewRedemption = true
		// Convert zenBTC amount to BTC amount.
		btcAmount := uint64(float64(redemption.Amount) * exchangeRate)
		if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Id, zenbtctypes.Redemption{
			Data: zenbtctypes.RedemptionData{
				Id:                 redemption.Id,
				DestinationAddress: redemption.DestinationAddress,
				Amount:             btcAmount,
			},
			Status: zenbtctypes.RedemptionStatus_INITIATED,
		}); err != nil {
			k.Logger(ctx).Error("error adding redemption to store", "error", err)
			continue
		}
	}

	if foundNewRedemption {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "error", err)
		}
	}
}

// processZenBTCRedemptions processes pending redemption completions.
func (k *Keeper) processZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	keyID := k.zenBTCKeeper.GetCompleterKeyID(ctx)
	requestedNonce := oracleData.RequestedCompleterNonce

	var initiatedRedemptions []zenbtctypes.Redemption
	if err := k.zenBTCKeeper.WalkRedemptions(ctx, func(id uint64, r zenbtctypes.Redemption) (bool, error) {
		if r.Status == zenbtctypes.RedemptionStatus_INITIATED {
			initiatedRedemptions = append(initiatedRedemptions, r)
			if len(initiatedRedemptions) == 2 {
				return true, nil
			}
		}
		return false, nil
	}); err != nil {
		k.Logger(ctx).Error("error getting initiated redemptions", "error", err)
		return
	}
	if len(initiatedRedemptions) == 0 {
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		return
	}

	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedNonce, initiatedRedemptions,
		func(redemption zenbtctypes.Redemption) error {
			redemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
			if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Data.Id, redemption); err != nil {
				return err
			}
			if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), true); err != nil {
				return err
			}
			return nil
		},
		func(redemption zenbtctypes.Redemption) error {
			k.Logger(ctx).Warn("processing zenBTC complete",
				"id", redemption.Data.Id,
				"nonce", requestedNonce,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap)
			unsignedTxHash, unsignedTx, err := k.constructCompleteTx(
				ctx,
				getChainIDForEigen(ctx),
				redemption.Data.Id,
				requestedNonce,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
			if err != nil {
				return err
			}
			creator, err := k.getAddressByKeyID(ctx, keyID, treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
				ctx,
				&treasurytypes.MsgNewSignTransactionRequest{
					Creator:             creator,
					KeyId:               keyID,
					WalletType:          treasurytypes.WalletType_WALLET_TYPE_EVM,
					UnsignedTransaction: unsignedTx,
					Metadata:            metadata,
					NoBroadcast:         false,
				},
				[]byte(hex.EncodeToString(unsignedTxHash)),
			)
			return err
		},
	)
}

//
// =============================================================================
// ORACLE DATA VALIDATION
// =============================================================================
//

// validateHashField derives a hash from the given data and compares it with the expected value.
func validateHashField(fieldName string, expectedHash []byte, data any) error {
	derivedHash, err := deriveHash(data)
	if err != nil {
		return fmt.Errorf("error deriving %s hash: %w", fieldName, err)
	}
	if !bytes.Equal(expectedHash, derivedHash[:]) {
		return fmt.Errorf("%s hash mismatch, expected %x, got %x", fieldName, expectedHash, derivedHash)
	}
	return nil
}

// validateOracleData verifies that the vote extension and oracle data match.
func (k *Keeper) validateOracleData(voteExt VoteExtension, oracleData *OracleData) error {
	if err := validateHashField("AVS contract delegation state", voteExt.EigenDelegationsHash, oracleData.EigenDelegationsMap); err != nil {
		return err
	}
	if err := validateHashField("Ethereum burn events", voteExt.EthBurnEventsHash, oracleData.EthBurnEvents); err != nil {
		return err
	}
	if err := validateHashField("Ethereum redemptions", voteExt.RedemptionsHash, oracleData.Redemptions); err != nil {
		return err
	}

	if voteExt.EthBlockHeight != oracleData.EthBlockHeight {
		return fmt.Errorf("ethereum block height mismatch, expected %d, got %d", voteExt.EthBlockHeight, oracleData.EthBlockHeight)
	}
	if voteExt.EthGasLimit != oracleData.EthGasLimit {
		return fmt.Errorf("ethereum gas limit mismatch, expected %d, got %d", voteExt.EthGasLimit, oracleData.EthGasLimit)
	}
	if voteExt.EthBaseFee != oracleData.EthBaseFee {
		return fmt.Errorf("ethereum base fee mismatch, expected %d, got %d", voteExt.EthBaseFee, oracleData.EthBaseFee)
	}
	if voteExt.EthTipCap != oracleData.EthTipCap {
		return fmt.Errorf("ethereum tip cap mismatch, expected %d, got %d", voteExt.EthTipCap, oracleData.EthTipCap)
	}

	if voteExt.BtcBlockHeight != oracleData.BtcBlockHeight {
		return fmt.Errorf("bitcoin block height mismatch, expected %d, got %d", voteExt.BtcBlockHeight, oracleData.BtcBlockHeight)
	}
	if err := validateHashField("Bitcoin header", voteExt.BtcHeaderHash, &oracleData.BtcBlockHeader); err != nil {
		return err
	}

	if voteExt.RequestedStakerNonce != oracleData.RequestedStakerNonce {
		return fmt.Errorf("requested staker nonce mismatch, expected %d, got %d", voteExt.RequestedStakerNonce, oracleData.RequestedStakerNonce)
	}
	if voteExt.RequestedEthMinterNonce != oracleData.RequestedEthMinterNonce {
		return fmt.Errorf("requested eth minter nonce mismatch, expected %d, got %d", voteExt.RequestedEthMinterNonce, oracleData.RequestedEthMinterNonce)
	}
	if voteExt.RequestedUnstakerNonce != oracleData.RequestedUnstakerNonce {
		return fmt.Errorf("requested unstaker nonce mismatch, expected %d, got %d", voteExt.RequestedUnstakerNonce, oracleData.RequestedUnstakerNonce)
	}
	if voteExt.RequestedCompleterNonce != oracleData.RequestedCompleterNonce {
		return fmt.Errorf("requested completer nonce mismatch, expected %d, got %d", voteExt.RequestedCompleterNonce, oracleData.RequestedCompleterNonce)
	}

	if !voteExt.ROCKUSDPrice.Equal(oracleData.ROCKUSDPrice) {
		return fmt.Errorf("ROCK/USD price mismatch, expected %s, got %s", voteExt.ROCKUSDPrice, oracleData.ROCKUSDPrice)
	}
	if !voteExt.BTCUSDPrice.Equal(oracleData.BTCUSDPrice) {
		return fmt.Errorf("BTC/USD price mismatch, expected %s, got %s", voteExt.BTCUSDPrice, oracleData.BTCUSDPrice)
	}
	if !voteExt.ETHUSDPrice.Equal(oracleData.ETHUSDPrice) {
		return fmt.Errorf("ETH/USD price mismatch, expected %s, got %s", voteExt.ETHUSDPrice, oracleData.ETHUSDPrice)
	}

	return nil
}
