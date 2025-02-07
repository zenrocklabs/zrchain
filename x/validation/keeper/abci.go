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

func (k *Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)
	return k.TrackHistoricalInfo(ctx)
}

func (k *Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)
	return k.BlockValidatorUpdates(ctx)
}

// ExtendVoteHandler is called by all validators to extend the consensus vote with additional data to be voted on.
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

// PrepareProposal is executed only by the proposer (1 validator on rotation) to inject oracle data into the block.
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
		return k.marshalOracleData(req, &OracleData{ConsensusData: req.LocalLastCommit}) // inject empty oracle data
	}

	if voteExt.ZRChainBlockHeight != req.Height-1 { // vote extension is created in ExtendVote step from the previous block
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

// PreBlocker is called before each block to process oracle data and update state.
// We don't return errors in the PreBlocker as this would halt the chain. Instead, we log errors and continue.
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

	k.updateAssetPrices(ctx, oracleData)

	k.updateValidatorStakes(ctx, oracleData)

	k.updateAVSDelegationStore(ctx, oracleData)

	k.storeBitcoinBlockHeader(ctx, oracleData)

	k.storeNewZenBTCBurnEventsEthereum(ctx, oracleData)

	k.storeNewZenBTCRedemptions(ctx, oracleData)

	// Toggle minting + unstaking every other block as VEs originate from block n-1 so nonce requests have 1 block latency
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

// shouldProcessOracleData checks if oracle data should be processed for this block
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

// validateCanonicalVE validates the proposed oracle data against the supermajority vote extension
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

// updateNonces handles updating nonce state for keys used for minting and unstaking
func (k *Keeper) updateNonces(ctx sdk.Context, oracleData OracleData) {
	for _, keyID := range k.getZenBTCKeyIDs(ctx) {
		requested, err := k.EthereumNonceRequested.Get(ctx, keyID)
		if err != nil && !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error checking nonce request state", "keyID", keyID, "error", err)
			continue
		}

		if !requested {
			continue
		}

		var currentNonce uint64
		switch keyID {
		case k.zenBTCKeeper.GetStakerKeyID(ctx):
			currentNonce = oracleData.RequestedStakerNonce
		case k.zenBTCKeeper.GetEthMinterKeyID(ctx):
			currentNonce = oracleData.RequestedEthMinterNonce
		case k.zenBTCKeeper.GetUnstakerKeyID(ctx):
			currentNonce = oracleData.RequestedUnstakerNonce
		case k.zenBTCKeeper.GetCompleterKeyID(ctx):
			currentNonce = oracleData.RequestedCompleterNonce
		default:
			k.Logger(ctx).Error("invalid key ID", "keyID", keyID)
			continue
		}

		// Don't set nonce to zero value erroneously if we already have a non-zero nonce
		lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, keyID)
		if err == nil && lastUsedNonce.Nonce != 0 && currentNonce == 0 {
			continue
		}

		if err := k.updateNonceState(ctx, keyID, currentNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "keyID", keyID, "error", err)
		}
	}
}

func (k *Keeper) getValidatedOracleData(ctx context.Context, voteExt VoteExtension) (*OracleData, *VoteExtension, error) {
	oracleData, err := k.GetSidecarStateByEthHeight(ctx, voteExt.EthBlockHeight)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching oracle state: %w", err)
	}

	bitcoinData, err := k.sidecarClient.GetBitcoinBlockHeaderByHeight(
		ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{ChainName: k.bitcoinNetwork(ctx), BlockHeight: voteExt.BtcBlockHeight},
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

func (k *Keeper) updateValidatorStakes(ctx sdk.Context, oracleData OracleData) {
	validatorInAVSDelegationSet := make(map[string]bool)

	for _, delegation := range oracleData.ValidatorDelegations {
		if delegation.Validator == "" {
			k.Logger(ctx).Debug("empty validator address in delegation; skipping")
			continue
		}

		valAddr, err := sdk.ValAddressFromBech32(delegation.Validator)
		if err != nil {
			k.Logger(ctx).Error("invalid validator address: "+delegation.Validator, "err", err)
			continue
		}

		validator, err := k.GetZenrockValidator(ctx, valAddr)
		if err != nil || validator.Status != types.Bonded {
			k.Logger(ctx).Debug("invalid delegation for "+delegation.Validator, "err", err, "reason", "invalid address / not bonded")
			continue
		}

		validator.TokensAVS = math.Int(delegation.Stake)

		if err = k.SetValidator(ctx, validator); err != nil {
			k.Logger(ctx).Error("error setting validator "+delegation.Validator, "err", err)
			continue
		}

		if err = k.ValidatorDelegations.Set(ctx, valAddr.String(), delegation.Stake); err != nil {
			k.Logger(ctx).Error("error setting validator delegations", "err", err)
			continue
		}

		validatorInAVSDelegationSet[valAddr.String()] = true
	}

	k.removeStaleValidatorDelegations(ctx, validatorInAVSDelegationSet)
}

func (k *Keeper) removeStaleValidatorDelegations(ctx sdk.Context, validatorInAVSDelegationSet map[string]bool) {
	var validatorsToRemove []string

	// First, collect the validators that need to be removed
	if err := k.ValidatorDelegations.Walk(ctx, nil, func(valAddr string, stake math.Int) (bool, error) {
		if !validatorInAVSDelegationSet[valAddr] {
			validatorsToRemove = append(validatorsToRemove, valAddr)
		}
		return true, nil
	}); err != nil {
		k.Logger(ctx).Error("error walking validator delegations", "err", err)
	}

	// Now, remove the collected validators (we can't do it while walking the store)
	for _, valAddr := range validatorsToRemove {
		if err := k.ValidatorDelegations.Remove(ctx, valAddr); err != nil {
			k.Logger(ctx).Error("error removing validator delegation", "validator", valAddr, "err", err)
			continue
		}

		if err := k.updateValidatorTokensAVS(ctx, valAddr); err != nil {
			k.Logger(ctx).Error("error updating validator TokensAVS", "validator", valAddr, "err", err)
		}
	}
}

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

func (k *Keeper) updateAVSDelegationStore(ctx sdk.Context, oracleData OracleData) {
	for validatorAddr, delegatorMap := range oracleData.EigenDelegationsMap {
		for delegatorAddr, amount := range delegatorMap {
			if err := k.AVSDelegations.Set(ctx, collections.Join(validatorAddr, delegatorAddr), math.NewIntFromBigInt(amount)); err != nil {
				k.Logger(ctx).Error("error setting AVS delegations", "err", err)
			}
		}
	}
}

func (k *Keeper) storeBitcoinBlockHeader(ctx sdk.Context, oracleData OracleData) {
	if oracleData.BtcBlockHeight == 0 || oracleData.BtcBlockHeader.MerkleRoot == "" {
		k.Logger(ctx).Error("invalid bitcoin header data", "height", oracleData.BtcBlockHeight, "merkle", oracleData.BtcBlockHeader.MerkleRoot)
	}

	requestedHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting requested historical Bitcoin headers", "err", err)
		return
	}

	// Check if this is a requested historical header
	isHistoricalHeader := false
	for _, height := range requestedHeaders.Heights {
		if height == oracleData.BtcBlockHeight {
			isHistoricalHeader = true
			break
		}
	}

	headerPreviouslySeen, err := k.BtcBlockHeaders.Has(ctx, oracleData.BtcBlockHeight)
	if err != nil {
		k.Logger(ctx).Error("error checking if Bitcoin header is already stored", "height", oracleData.BtcBlockHeight, "err", err)
		return
	}

	if err := k.BtcBlockHeaders.Set(ctx, oracleData.BtcBlockHeight, oracleData.BtcBlockHeader); err != nil {
		k.Logger(ctx).Error("error storing Bitcoin header", "height", oracleData.BtcBlockHeight, "err", err)
		return
	}

	// If it's a historical header, remove it from the requested list and return early
	if isHistoricalHeader {
		requestedHeaders.Heights = slices.DeleteFunc(requestedHeaders.Heights, func(height int64) bool {
			return height == oracleData.BtcBlockHeight
		})

		if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
			k.Logger(ctx).Error("error updating requested historical Bitcoin headers", "err", err)
			return
		}

		k.Logger(ctx).Debug("successfully stored historical Bitcoin header and removed request",
			"height", oracleData.BtcBlockHeight,
			"remaining_requests", len(requestedHeaders.Heights))
		return
	}

	if headerPreviouslySeen {
		k.Logger(ctx).Debug("bitcoin header previously seen; skipping reorg check", "height", oracleData.BtcBlockHeight)
		return
	}

	if err := k.checkForBitcoinReorg(ctx, oracleData, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error handling potential Bitcoin reorg", "height", oracleData.BtcBlockHeight, "err", err)
	}
}

// checkForBitcoinReorg detects reorgs by requesting previous headers when a new one is received
func (k *Keeper) checkForBitcoinReorg(
	ctx sdk.Context,
	oracleData OracleData,
	requestedHeaders zenbtctypes.RequestedBitcoinHeaders,
) error {
	var numHistoricalHeadersToRequest int64 = 20     // default for non-mainnet environments
	if strings.HasPrefix(ctx.ChainID(), "diamond") { // mainnet
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
		k.Logger(ctx).Error("error setting requested historical Bitcoin headers", "err", err)
		return err
	}

	return nil
}

func (k *Keeper) processZenBTCStaking(ctx sdk.Context, oracleData OracleData) {
	requested, err := k.EthereumNonceRequested.Get(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx))
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting EthereumNonceRequested state", "err", err)
		return
	}
	if !requested {
		return
	}

	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx))
	if err != nil {
		k.Logger(ctx).Error("error getting last used Ethereum nonce", "err", err)
		return
	}

	k.Logger(ctx).Info("lastUsedNonce",
		"nonce", lastUsedNonce.Nonce,
		"counter", lastUsedNonce.Counter,
		"skip", lastUsedNonce.Skip,
		"requested_nonce", oracleData.RequestedStakerNonce,
	)

	if lastUsedNonce.Nonce != 0 && oracleData.RequestedStakerNonce == 0 {
		return
	}

	// Gather up to 2 deposited transactions
	var depositedTxs []zenbtctypes.PendingMintTransaction
	err = k.zenBTCKeeper.WalkPendingMintTransactions(ctx, func(id uint64, pending zenbtctypes.PendingMintTransaction) (bool, error) {
		if pending.Status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED {
			depositedTxs = append(depositedTxs, pending)
			if len(depositedTxs) == 2 {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		k.Logger(ctx).Error("error walking pending mint transactions", "err", err)
		return
	}

	// If no deposited tx exists, clear the nonce flag and return
	if len(depositedTxs) == 0 {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), false); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state for staker", "err", err)
		}
		return
	}

	// When a nonce update has been detected, process the staking update
	if oracleData.RequestedStakerNonce != lastUsedNonce.PrevNonce {
		firstTx := depositedTxs[0]

		firstTx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED
		if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, firstTx); err != nil {
			k.Logger(ctx).Error("error setting pending stake transactions", "err", err)
			return
		}

		k.Logger(ctx).Warn("updated stake transaction", "tx", fmt.Sprintf("%+v", firstTx))

		lastUsedNonce.PrevNonce = lastUsedNonce.Nonce
		if err := k.LastUsedEthereumNonce.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), lastUsedNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "err", err)
		}

		// Set EthereumNonceRequested to true for the minter key after successful stake
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state for minter", "err", err)
		}

		// If no more transactions to process, set staker nonce request to false
		if len(depositedTxs) == 1 {
			if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), false); err != nil {
				k.Logger(ctx).Error("error setting EthereumNonceRequested state for staker", "err", err)
			}
		}
		return
	}

	if lastUsedNonce.Skip {
		return
	}

	// Decide which transaction to use for staking
	var stakeTx zenbtctypes.PendingMintTransaction
	if len(depositedTxs) > 1 {
		stakeTx = depositedTxs[1]
	} else {
		stakeTx = depositedTxs[0]
	}

	k.Logger(ctx).Warn("processing zenBTC stake",
		"recipient", stakeTx.RecipientAddress,
		"amount", stakeTx.Amount,
		"nonce", oracleData.RequestedStakerNonce,
		"gas_limit", oracleData.EthGasLimit,
		"base_fee", oracleData.EthBaseFee,
		"tip_cap", oracleData.EthTipCap,
	)

	unsignedStakeTxHash, unsignedStakeTx, err := k.constructStakeTx(
		ctx,
		getChainIDForEigen(ctx),
		stakeTx.Amount,
		oracleData.RequestedStakerNonce,
		oracleData.EthGasLimit,
		oracleData.EthBaseFee,
		oracleData.EthTipCap,
	)
	if err != nil {
		k.Logger(ctx).Error("error constructing stake transaction", "err", err)
		return
	}

	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
	if err != nil {
		k.Logger(ctx).Error("error creating metadata", "err", err)
		return
	}

	if _, err := k.treasuryKeeper.HandleSignTransactionRequest(
		ctx,
		&treasurytypes.MsgNewSignTransactionRequest{
			Creator:             stakeTx.Creator,
			KeyId:               k.zenBTCKeeper.GetStakerKeyID(ctx),
			WalletType:          treasurytypes.WalletType(stakeTx.ChainType),
			UnsignedTransaction: unsignedStakeTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedStakeTxHash)),
	); err != nil {
		k.Logger(ctx).Error("error creating stake transaction", "err", err)
	}
}

func (k *Keeper) processZenBTCMints(ctx sdk.Context, oracleData OracleData) {
	requested, err := k.EthereumNonceRequested.Get(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx))
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting EthereumNonceRequested state", "err", err)
		return
	}
	if !requested {
		return
	}

	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx))
	if err != nil {
		k.Logger(ctx).Error("error getting last used Ethereum nonce", "err", err)
		return
	}

	k.Logger(ctx).Info("lastUsedNonce", "nonce", lastUsedNonce.Nonce, "counter", lastUsedNonce.Counter, "skip", lastUsedNonce.Skip, "requested_nonce", oracleData.RequestedEthMinterNonce)

	if lastUsedNonce.Nonce != 0 && oracleData.RequestedEthMinterNonce == 0 {
		return
	}

	// Gather up to 2 staked transactions
	var stakedTxs []zenbtctypes.PendingMintTransaction
	err = k.zenBTCKeeper.WalkPendingMintTransactions(ctx, func(id uint64, pending zenbtctypes.PendingMintTransaction) (bool, error) {
		if pending.Status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED {
			stakedTxs = append(stakedTxs, pending)
			if len(stakedTxs) == 2 {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		k.Logger(ctx).Error("error walking pending mint transactions", "err", err)
		return
	}

	// If no staked tx exists, clear the nonce flag and return
	if len(stakedTxs) == 0 {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), false); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state for minter", "err", err)
		}
		return
	}

	// When a nonce update has been detected, process the staking update
	if oracleData.RequestedEthMinterNonce != lastUsedNonce.PrevNonce {
		firstTx := stakedTxs[0]

		supply, err := k.zenBTCKeeper.GetSupply(ctx)
		if err != nil {
			k.Logger(ctx).Error("error getting zenBTC supply", "err", err)
			return
		}

		supply.PendingZenBTC -= firstTx.Amount
		supply.MintedZenBTC += firstTx.Amount

		if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
			k.Logger(ctx).Error("error updating zenBTC supply", "err", err)
			return
		}
		k.Logger(ctx).Warn("pending mint supply updated", "pending_mint_old", supply.PendingZenBTC+firstTx.Amount, "pending_mint_new", supply.PendingZenBTC)
		k.Logger(ctx).Warn("minted supply updated", "minted_old", supply.MintedZenBTC-firstTx.Amount, "minted_new", supply.MintedZenBTC)

		firstTx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
		if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, firstTx); err != nil {
			k.Logger(ctx).Error("error setting pending mint transactions", "err", err)
			return
		}

		k.Logger(ctx).Warn("updated mint transaction", "tx", fmt.Sprintf("%+v", firstTx))

		lastUsedNonce.PrevNonce = lastUsedNonce.Nonce
		if err := k.LastUsedEthereumNonce.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), lastUsedNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "err", err)
		}

		// If no more transactions to process, set minter nonce request to false
		if len(stakedTxs) == 1 {
			if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), false); err != nil {
				k.Logger(ctx).Error("error setting EthereumNonceRequested state for minter", "err", err)
			}
		}
		return
	}

	if lastUsedNonce.Skip {
		return
	}

	// Decide which transaction to use for minting
	var mintTx zenbtctypes.PendingMintTransaction
	if len(stakedTxs) > 1 {
		mintTx = stakedTxs[1]
	} else {
		mintTx = stakedTxs[0]
	}

	exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting zenBTC exchange rate", "err", err)
		return
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
		return
	}

	k.Logger(ctx).Warn("processing zenBTC mint",
		"recipient", mintTx.RecipientAddress,
		"amount", mintTx.Amount,
		"nonce", oracleData.RequestedEthMinterNonce,
		"gas_limit", oracleData.EthGasLimit,
		"base_fee", oracleData.EthBaseFee,
		"tip_cap", oracleData.EthTipCap,
		"chain_id", mintTx.Caip2ChainId,
		"fee_zen_btc", feeZenBTC,
	)

	if mintTx.Caip2ChainId != "eip155:17000" {
		k.Logger(ctx).Error("invalid chain ID", "chain_id", mintTx.Caip2ChainId)
		return
	}

	chainID, err := types.ExtractEVMChainID(mintTx.Caip2ChainId)
	if err != nil {
		k.Logger(ctx).Error("error extracting chainId from CAIP-2", "err", err)
	}

	unsignedMintTxHash, unsignedMintTx, err := k.constructMintTx(
		ctx,
		mintTx.RecipientAddress,
		chainID,
		mintTx.Amount,
		feeZenBTC,
		oracleData.RequestedStakerNonce,
		oracleData.EthGasLimit,
		oracleData.EthBaseFee,
		oracleData.EthTipCap,
	)
	if err != nil {
		k.Logger(ctx).Error("error constructing mint transaction", "err", err)
		return
	}

	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: chainID})
	if err != nil {
		k.Logger(ctx).Error("error creating metadata", "err", err)
		return
	}

	if _, err := k.treasuryKeeper.HandleSignTransactionRequest(
		ctx,
		&treasurytypes.MsgNewSignTransactionRequest{
			Creator:             mintTx.Creator,
			KeyId:               k.zenBTCKeeper.GetEthMinterKeyID(ctx),
			WalletType:          treasurytypes.WalletType(mintTx.ChainType),
			UnsignedTransaction: unsignedMintTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedMintTxHash)),
	); err != nil {
		k.Logger(ctx).Error("error creating mint transaction", "err", err)
	}
}

func (k *Keeper) storeNewZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	// Retrieve the current burn events from the store
	burnEvents, err := k.zenBTCKeeper.GetBurnEvents(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("failed to get current burn events", "err", err)
		}
		return
	}

	// Loop over each burn event from oracle to check for new ones
	for _, burn := range oracleData.EthBurnEvents {
		exists := false
		newBurn := zenbtctypes.BurnEvent(burn)
		for _, existingBurn := range burnEvents.Events {
			if reflect.DeepEqual(newBurn, *existingBurn) {
				exists = true
				break
			}
		}
		if !exists {
			burnEvents.Events = append(burnEvents.Events, &newBurn)
		}
	}

	if err := k.zenBTCKeeper.SetBurnEvents(ctx, burnEvents); err != nil {
		k.Logger(ctx).Error("error setting burn events", "err", err)
	}
}

func (k *Keeper) storeNewZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	// First, find the first INITIATED redemption
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
		k.Logger(ctx).Error("error finding first initiated redemption", "err", err)
		return
	}

	// If we found an INITIATED redemption, check if it exists in oracleData
	if found {
		redemptionExists := false
		for _, redemption := range oracleData.Redemptions {
			if redemption.Id == firstInitiatedRedemption.Data.Id {
				redemptionExists = true
				break
			}
		}

		// If the redemption is not in oracleData, mark it as unstaked
		if !redemptionExists {
			firstInitiatedRedemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
			if err := k.zenBTCKeeper.SetRedemption(ctx, firstInitiatedRedemption.Data.Id, firstInitiatedRedemption); err != nil {
				k.Logger(ctx).Error("error updating redemption status to unstaked", "err", err)
				return
			}
		}
	}

	if len(oracleData.Redemptions) == 0 {
		return
	}

	// Get current exchange rate for conversion
	exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting zenBTC exchange rate", "err", err)
		return
	}

	foundNewRedemption := false

	for _, redemption := range oracleData.Redemptions {
		redemptionExists, err := k.zenBTCKeeper.HasRedemption(ctx, redemption.Id)
		if err != nil {
			k.Logger(ctx).Error("error checking redemption existence", "err", err)
			continue
		}
		if redemptionExists {
			k.Logger(ctx).Debug("redemption already stored", "id", redemption.Id)
			continue
		}

		foundNewRedemption = true

		// Convert zenBTC amount to BTC amount
		// redemption.Amount is zenBTC, multiply by BTC/zenBTC rate to get BTC amount
		btcAmount := uint64(float64(redemption.Amount) * exchangeRate)

		if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Id, zenbtctypes.Redemption{
			Data: zenbtctypes.RedemptionData{
				Id:                 redemption.Id,
				DestinationAddress: redemption.DestinationAddress,
				Amount:             btcAmount,
			},
			Status: zenbtctypes.RedemptionStatus_INITIATED,
		}); err != nil {
			k.Logger(ctx).Error("error adding redemption to store", "err", err)
			continue
		}
	}

	if foundNewRedemption {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "err", err)
		}
	}
}

// processZenBTCBurnEventsEthereum processes pending burn events by constructing unstake transactions.
func (k *Keeper) processZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	// Use the unstaker key ID for processing burn events.
	keyID := k.zenBTCKeeper.GetUnstakerKeyID(ctx)

	// Check if a nonce request for unstaking is currently active.
	requested, err := k.EthereumNonceRequested.Get(ctx, keyID)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting EthereumNonceRequested state for burn events", "err", err)
		return
	}
	if !requested {
		return
	}

	// Retrieve the last used Ethereum nonce for the unstaker key.
	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, keyID)
	if err != nil {
		k.Logger(ctx).Error("error getting last used Ethereum nonce for burn events", "err", err)
		return
	}
	if lastUsedNonce.Skip {
		return
	}

	// Retrieve the current burn events from storage.
	burnEvents, err := k.zenBTCKeeper.GetBurnEvents(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("failed to get current burn events", "err", err)
		}
		return
	}

	// If there are no burn events, exit.
	if len(burnEvents.Events) == 0 {
		return
	}

	// --- Nonce update branch ---
	// If the unstaker nonce has changed (meaning the previous unstake succeeded),
	// update the nonce state and remove the processed event.
	if oracleData.RequestedUnstakerNonce != lastUsedNonce.PrevNonce {
		k.Logger(ctx).Warn("unstaker nonce updated for burn events",
			"nonce", oracleData.RequestedUnstakerNonce,
			"prev_nonce", lastUsedNonce.PrevNonce)

		// Update the nonce state.
		lastUsedNonce.PrevNonce = lastUsedNonce.Nonce
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, lastUsedNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state for burn events", "err", err)
		}

		// Set EthereumNonceRequested to true for the completer key after successful unstake.
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state for completer", "err", err)
		}

		// Remove the first processed event.
		newEvents := burnEvents.Events[1:]
		burnEvents.Events = newEvents
		if err := k.zenBTCKeeper.SetBurnEvents(ctx, burnEvents); err != nil {
			k.Logger(ctx).Error("error setting updated burn events", "err", err)
		}

		// If no more burn events remain, clear the unstaker nonce request.
		if len(newEvents) == 0 {
			if err := k.EthereumNonceRequested.Set(ctx, keyID, false); err != nil {
				k.Logger(ctx).Error("error setting EthereumNonceRequested state for unstaker", "err", err)
			}
		}
		return
	}

	// --- Processing branch ---
	// Instead of always processing the first burn event, choose the second if available.
	var burnEvent *zenbtctypes.BurnEvent
	if len(burnEvents.Events) > 1 {
		burnEvent = burnEvents.Events[1]
	} else {
		burnEvent = burnEvents.Events[0]
	}

	k.Logger(ctx).Warn("processing zenBTC burn unstake",
		"burn_event", burnEvent,
		"nonce", oracleData.RequestedUnstakerNonce,
		"base_fee", oracleData.EthBaseFee,
		"tip_cap", oracleData.EthTipCap)

	unsignedTxHash, unsignedTx, err := k.constructUnstakeTx(
		ctx,
		getChainIDForEigen(ctx),
		burnEvent.DestinationAddr,
		burnEvent.Amount,
		oracleData.RequestedUnstakerNonce,
		oracleData.EthBaseFee,
		oracleData.EthTipCap,
	)
	if err != nil {
		k.Logger(ctx).Error("error constructing unstake transaction for burn event", "err", err)
		return
	}

	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
	if err != nil {
		k.Logger(ctx).Error("error creating metadata for burn event unstake tx", "err", err)
		return
	}

	creator, err := k.getAddressByKeyID(ctx, keyID, treasurytypes.WalletType_WALLET_TYPE_NATIVE)
	if err != nil {
		k.Logger(ctx).Error("error getting creator address for burn event unstake tx", "err", err)
		return
	}

	if _, err := k.treasuryKeeper.HandleSignTransactionRequest(
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
	); err != nil {
		k.Logger(ctx).Error("error creating unstake transaction for burn event", "err", err)
	}
}

func (k *Keeper) processZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	requested, err := k.EthereumNonceRequested.Get(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx))
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting EthereumNonceRequested state", "err", err)
		return
	}
	if !requested {
		return
	}

	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx))
	if err != nil {
		k.Logger(ctx).Error("error getting last used Ethereum nonce", "err", err)
		return
	}

	if lastUsedNonce.Skip {
		return
	}

	// Gather up to 2 initiated redemptions
	var initiatedRedemptions []zenbtctypes.Redemption
	err = k.zenBTCKeeper.WalkRedemptions(ctx, func(id uint64, r zenbtctypes.Redemption) (bool, error) {
		if r.Status == zenbtctypes.RedemptionStatus_INITIATED {
			initiatedRedemptions = append(initiatedRedemptions, r)
			if len(initiatedRedemptions) == 2 {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		k.Logger(ctx).Error("error walking redemptions", "err", err)
		return
	}

	// If no initiated redemption exists, clear the nonce flag and return
	if len(initiatedRedemptions) == 0 {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), false); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state for completer", "err", err)
		}
		return
	}

	// When a nonce update has been detected, process the redemption update
	if oracleData.RequestedCompleterNonce != lastUsedNonce.PrevNonce {
		firstRedemption := initiatedRedemptions[0]

		firstRedemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
		if err := k.zenBTCKeeper.SetRedemption(ctx, firstRedemption.Data.Id, firstRedemption); err != nil {
			k.Logger(ctx).Error("error updating redemption status", "err", err)
			return
		}

		lastUsedNonce.PrevNonce = lastUsedNonce.Nonce
		if err := k.LastUsedEthereumNonce.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), lastUsedNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "err", err)
		}

		// Set EthereumNonceRequested to true for the staker key after successful completion
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state for staker", "err", err)
		}

		// If no more redemptions to process, set completer nonce request to false
		if len(initiatedRedemptions) == 1 {
			if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), false); err != nil {
				k.Logger(ctx).Error("error setting EthereumNonceRequested state for completer", "err", err)
			}
		}
		return
	}

	// Decide which redemption to use for completion
	var completeRedemption zenbtctypes.Redemption
	if len(initiatedRedemptions) > 1 {
		completeRedemption = initiatedRedemptions[1]
	} else {
		completeRedemption = initiatedRedemptions[0]
	}

	k.Logger(ctx).Warn("processing zenBTC complete",
		"id", completeRedemption.Data.Id,
		"nonce", oracleData.RequestedCompleterNonce,
		"base_fee", oracleData.EthBaseFee,
		"tip_cap", oracleData.EthTipCap,
	)

	unsignedTxHash, unsignedTx, err := k.constructCompleteTx(
		ctx,
		getChainIDForEigen(ctx),
		completeRedemption.Data.Id,
		oracleData.RequestedCompleterNonce,
		oracleData.EthBaseFee,
		oracleData.EthTipCap,
	)
	if err != nil {
		k.Logger(ctx).Error("error constructing complete transaction", "err", err)
		return
	}

	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
	if err != nil {
		k.Logger(ctx).Error("error creating metadata", "err", err)
		return
	}

	creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
	if err != nil {
		k.Logger(ctx).Error("error getting creator address", "err", err)
		return
	}

	if _, err := k.treasuryKeeper.HandleSignTransactionRequest(
		ctx,
		&treasurytypes.MsgNewSignTransactionRequest{
			Creator:             creator,
			KeyId:               k.zenBTCKeeper.GetCompleterKeyID(ctx),
			WalletType:          treasurytypes.WalletType_WALLET_TYPE_EVM,
			UnsignedTransaction: unsignedTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedTxHash)),
	); err != nil {
		k.Logger(ctx).Error("error creating complete transaction", "err", err)
	}
}

func (k *Keeper) validateOracleData(voteExt VoteExtension, oracleData *OracleData) error {
	eigenDelegationsHash, err := deriveHash(oracleData.EigenDelegationsMap)
	if err != nil {
		return fmt.Errorf("error deriving AVS contract delegation state hash: %w", err)
	}
	if !bytes.Equal(voteExt.EigenDelegationsHash, eigenDelegationsHash[:]) {
		return fmt.Errorf("AVS contract delegation state hash mismatch, expected %x, got %x", voteExt.EigenDelegationsHash, eigenDelegationsHash)
	}

	ethBurnEventsHash, err := deriveHash(oracleData.EthBurnEvents)
	if err != nil {
		return fmt.Errorf("error deriving ethereum burn events hash: %w", err)
	}
	if !bytes.Equal(voteExt.EthBurnEventsHash, ethBurnEventsHash[:]) {
		return fmt.Errorf("ethereum burn events hash mismatch, expected %x, got %x", voteExt.EthBurnEventsHash, ethBurnEventsHash)
	}

	ethereumRedemptionsHash, err := deriveHash(oracleData.Redemptions)
	if err != nil {
		return fmt.Errorf("error deriving redemptions hash: %w", err)
	}
	if !bytes.Equal(voteExt.RedemptionsHash, ethereumRedemptionsHash[:]) {
		return fmt.Errorf("ethereum redemptions hash mismatch, expected %x, got %x", voteExt.RedemptionsHash, ethereumRedemptionsHash)
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

	bitcoinHeaderHash, err := deriveHash(&oracleData.BtcBlockHeader)
	if err != nil {
		return fmt.Errorf("error deriving bitcoin header hash: %w", err)
	}
	if !bytes.Equal(voteExt.BtcHeaderHash, bitcoinHeaderHash[:]) {
		return fmt.Errorf("bitcoin header hash mismatch, expected %x, got %x", voteExt.BtcHeaderHash, bitcoinHeaderHash)
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
