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
	avsDelegationsHash, err := deriveAVSContractStateHash(oracleData.EigenDelegationsMap)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving AVS contract delegation state hash: %w", err)
	}

	ethereumRedemptionsHash, err := deriveRedemptionsHash(oracleData.EthereumRedemptions)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving ethereum redemptions hash: %w", err)
	}
	solanaRedemptionsHash, err := deriveRedemptionsHash(oracleData.SolanaRedemptions)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving solana redemptions hash: %w", err)
	}

	neutrinoResponse, err := k.retrieveBitcoinHeader(ctx)
	if err != nil {
		return VoteExtension{}, err
	}
	bitcoinHeaderHash, err := deriveBitcoinHeaderHash(neutrinoResponse.BlockHeader)
	if err != nil {
		return VoteExtension{}, err
	}

	nonces := make(map[uint64]uint64)
	keys := []uint64{k.GetZenBTCMinterKeyID(ctx), k.GetZenBTCUnstakerKeyID(ctx)}
	for _, key := range keys {
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
		EthereumRedemptionsHash:    ethereumRedemptionsHash[:],
		SolanaRedemptionsHash:      solanaRedemptionsHash[:],
		BtcBlockHeight:             neutrinoResponse.BlockHeight,
		BtcHeaderHash:              bitcoinHeaderHash[:],
		EthBlockHeight:             oracleData.EthBlockHeight,
		EthGasLimit:                oracleData.EthGasLimit,
		EthBaseFee:                 oracleData.EthBaseFee,
		EthTipCap:                  oracleData.EthTipCap,
		SolanaLamportsPerSignature: oracleData.SolanaLamportsPerSignature,
		RequestedEthMinterNonce:    nonces[k.GetZenBTCMinterKeyID(ctx)],
		RequestedEthUnstakerNonce:  nonces[k.GetZenBTCUnstakerKeyID(ctx)],
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

	k.storeNewZenBTCRedemptionsEthereum(ctx, oracleData)

	// Toggle minting + unstaking every other block as VEs originate from block n-1 so nonce requests have 1 block latency
	if ctx.BlockHeight()%2 == 0 {

		k.processZenBTCMints(ctx, oracleData)

		k.processZenBTCRedemptionsEthereum(ctx, oracleData)

		k.updateNonces(ctx, oracleData)
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
	keys := []uint64{k.GetZenBTCMinterKeyID(ctx), k.GetZenBTCUnstakerKeyID(ctx)}
	for _, keyID := range keys {
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
		case k.GetZenBTCMinterKeyID(ctx):
			currentNonce = oracleData.RequestedEthMinterNonce
		case k.GetZenBTCUnstakerKeyID(ctx):
			currentNonce = oracleData.RequestedEthUnstakerNonce
		default:
			k.Logger(ctx).Error("invalid key ID", "keyID", keyID)
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
	oracleData.RequestedEthMinterNonce = voteExt.RequestedEthMinterNonce
	oracleData.RequestedEthUnstakerNonce = voteExt.RequestedEthUnstakerNonce

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
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting requested historical Bitcoin headers", "err", err)
			return
		}
		requestedHeaders = zenbtctypes.RequestedBitcoinHeaders{}
		if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
			k.Logger(ctx).Error("error setting requested historical Bitcoin headers", "err", err)
			return
		}
	}

	// Check if this is a requested historical header
	isRequestedHistoricalHeader := false
	for _, height := range requestedHeaders.Heights {
		if height == oracleData.BtcBlockHeight {
			isRequestedHistoricalHeader = true
			break
		}
	}

	if err := k.BtcBlockHeaders.Set(ctx, oracleData.BtcBlockHeight, oracleData.BtcBlockHeader); err != nil {
		k.Logger(ctx).Error("error storing Bitcoin header", "height", oracleData.BtcBlockHeight, "err", err)
	}

	// If it's a requested historical header, remove it from the requested list and return early
	if isRequestedHistoricalHeader {
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

	if err := k.checkForBitcoinReorg(ctx, oracleData, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error handling potential Bitcoin fork", "height", oracleData.BtcBlockHeight, "err", err)
	}
}

// checkForBitcoinReorg detects reorgs by requesting previous blocks when a new header is received
func (k *Keeper) checkForBitcoinReorg(
	ctx sdk.Context,
	oracleData OracleData,
	requestedHeaders zenbtctypes.RequestedBitcoinHeaders,
) error {
	prevHeights := make([]int64, 0, 6)
	for i := int64(1); i <= 6; i++ {
		prevHeight := oracleData.BtcBlockHeight - i
		if prevHeight <= 0 {
			break
		}
		prevHeights = append(prevHeights, prevHeight)
	}

	if len(prevHeights) == 0 {
		k.Logger(ctx).Error("no previous heights to request (this should never happen)", "height", oracleData.BtcBlockHeight)
		return nil
	}

	requestedHeaders.Heights = append(requestedHeaders.Heights, prevHeights...)

	if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error setting requested historical Bitcoin headers", "err", err)
		return err
	}

	return nil
}

func (k *Keeper) processZenBTCMints(ctx sdk.Context, oracleData OracleData) {
	requested, err := k.EthereumNonceRequested.Get(ctx, k.GetZenBTCMinterKeyID(ctx))
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting EthereumNonceRequested state", "err", err)
		}
		requested = false
		if err := k.EthereumNonceRequested.Set(ctx, k.GetZenBTCMinterKeyID(ctx), requested); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "err", err)
		}
	}
	if !requested {
		return
	}

	pendingMints, err := k.PendingMintTransactions.Get(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting pending mint transactions", "err", err)
		return
	}
	if len(pendingMints.Txs) == 0 {
		return
	}

	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, k.GetZenBTCMinterKeyID(ctx))
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting last used Ethereum nonce", "err", err)
		}
		lastUsedNonce = zenbtctypes.NonceData{Nonce: oracleData.RequestedEthMinterNonce, Counter: 0, Skip: true}
		if err := k.LastUsedEthereumNonce.Set(ctx, k.GetZenBTCMinterKeyID(ctx), lastUsedNonce); err != nil {
			k.Logger(ctx).Error("error setting last used Ethereum nonce", "err", err)
		}
	}

	if lastUsedNonce.Nonce != 0 && oracleData.RequestedEthMinterNonce == 0 {
		return
	}

	lastMintTx := pendingMints.Txs[0]

	// remove last pending tx + update supply (after nonce updated indicating successful mint)
	if oracleData.RequestedEthMinterNonce != lastUsedNonce.Nonce {
		supply, err := k.ZenBTCSupply.Get(ctx)
		if err != nil {
			k.Logger(ctx).Error("error getting zenBTC supply", "err", err)
		}

		supply.MintedZenBTC += lastMintTx.Amount

		if err := k.ZenBTCSupply.Set(ctx, supply); err != nil {
			k.Logger(ctx).Error("error updating zenBTC supply", "err", err)
		}

		lastUsedNonce.Nonce = oracleData.RequestedEthMinterNonce
		lastUsedNonce.Counter = 0
		if err := k.LastUsedEthereumNonce.Set(ctx, k.GetZenBTCMinterKeyID(ctx), lastUsedNonce); err != nil {
			k.Logger(ctx).Error("error setting last used Ethereum nonce", "err", err)
		}

		pendingMints.Txs = pendingMints.Txs[1:]
		if err := k.PendingMintTransactions.Set(ctx, pendingMints); err != nil {
			k.Logger(ctx).Error("error setting pending mint transactions", "err", err)
		}

		if len(pendingMints.Txs) == 0 {
			if err := k.EthereumNonceRequested.Set(ctx, k.GetZenBTCMinterKeyID(ctx), false); err != nil {
				k.Logger(ctx).Error("error setting EthereumNonceRequested state", "err", err)
			}
			return
		}
	}

	if lastUsedNonce.Skip {
		return
	}

	pendingMintTx := pendingMints.Txs[0]

	// baseFeePlusTip := new(big.Int).Add(new(big.Int).SetUint64(oracleData.EthBaseFee), new(big.Int).SetUint64(oracleData.EthTipCap))
	// feeETH := new(big.Int).Mul(baseFeePlusTip, new(big.Int).SetUint64(oracleData.EthGasLimit))

	if oracleData.BTCUSDPrice.IsZero() {
		return
	}
	// ethToBTC := oracleData.ETHUSDPrice.Quo(oracleData.BTCUSDPrice)
	// feeBTCFloat := new(big.Float).Mul(new(big.Float).SetInt(feeETH), new(big.Float).SetFloat64(ethToBTC.MustFloat64()))
	// feeBTCInt, _ := feeBTCFloat.Int(nil)
	// feeBTC := feeBTCInt.Uint64()

	unsignedMintTxHash, unsignedMintTx, err := k.constructMintTx(
		ctx,
		pendingMintTx.RecipientAddress,
		pendingMintTx.ChainId,
		pendingMintTx.Amount,
		// feeBTC,
		0,
		oracleData.RequestedEthMinterNonce,
		oracleData.EthGasLimit,
		oracleData.EthBaseFee,
		oracleData.EthTipCap,
	)
	if err != nil {
		k.Logger(ctx).Error("error constructing mint transaction", "err", err)
	}

	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: pendingMintTx.ChainId})
	if err != nil {
		k.Logger(ctx).Error("error creating metadata", "err", err)
	}

	if _, err := k.treasuryKeeper.HandleSignTransactionRequest(
		ctx,
		&treasurytypes.MsgNewSignTransactionRequest{
			Creator:             pendingMintTx.Creator,
			KeyId:               pendingMintTx.KeyId,
			WalletType:          pendingMintTx.ChainType,
			UnsignedTransaction: unsignedMintTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedMintTxHash)),
	); err != nil {
		k.Logger(ctx).Error("error creating mint transaction", "err", err)
	}
}

func (k *Keeper) storeNewZenBTCRedemptionsEthereum(ctx sdk.Context, oracleData OracleData) {
	if len(oracleData.EthereumRedemptions) == 0 {
		return
	}

	// Get current exchange rate for conversion
	exchangeRate, err := k.GetZenBTCExchangeRate(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting zenBTC exchange rate", "err", err)
		return
	}

	foundNewRedemption := false

	for _, redemption := range oracleData.EthereumRedemptions {
		redemptionExists, err := k.ZenBTCRedemptions.Has(ctx, redemption.Id)
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

		if err := k.ZenBTCRedemptions.Set(ctx, redemption.Id, zenbtctypes.Redemption{
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
		if err := k.EthereumNonceRequested.Set(ctx, k.GetZenBTCUnstakerKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "err", err)
		}
	}
}

func (k *Keeper) processZenBTCRedemptionsEthereum(ctx sdk.Context, oracleData OracleData) {
	// Check if we should process redemptions
	requested, err := k.EthereumNonceRequested.Get(ctx, k.GetZenBTCUnstakerKeyID(ctx))
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting EthereumNonceRequested state", "err", err)
		return
	}
	if !requested {
		return
	}

	// Get last used nonce state
	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, k.GetZenBTCUnstakerKeyID(ctx))
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting last used Ethereum nonce", "err", err)
		return
	}

	if lastUsedNonce.Skip {
		return
	}

	// Get the first INITIATED redemption
	var redemption zenbtctypes.Redemption
	var redemptionID uint64
	var found bool

	if err := k.ZenBTCRedemptions.Walk(ctx, nil, func(id uint64, r zenbtctypes.Redemption) (bool, error) {
		if r.Status == zenbtctypes.RedemptionStatus_INITIATED {
			redemption = r
			redemptionID = id
			found = true
			return true, nil
		}
		return false, nil
	}); err != nil {
		k.Logger(ctx).Error("error finding redemption", "err", err)
		return
	}
	if !found {
		if err := k.EthereumNonceRequested.Set(ctx, k.GetZenBTCUnstakerKeyID(ctx), false); err != nil {
			k.Logger(ctx).Error("error updating nonce request state", "err", err)
		}
		return
	}

	// If nonce changed, previous unstake succeeded - update status
	if oracleData.RequestedEthUnstakerNonce != lastUsedNonce.Nonce {
		redemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
		if err := k.ZenBTCRedemptions.Set(ctx, redemptionID, redemption); err != nil {
			k.Logger(ctx).Error("error updating redemption status", "err", err)
			return
		}
		lastUsedNonce.Nonce = oracleData.RequestedEthUnstakerNonce
		lastUsedNonce.Counter = 0
		if err := k.LastUsedEthereumNonce.Set(ctx, k.GetZenBTCUnstakerKeyID(ctx), lastUsedNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "err", err)
		}
		return
	}

	// Create and sign new unstake transaction
	unsignedTxHash, unsignedTx, err := k.constructUnstakeTx(
		ctx,
		redemption.Data.Id,
		17000, // TODO: make this dynamic
		oracleData.RequestedEthUnstakerNonce,
		oracleData.EthBaseFee,
		oracleData.EthTipCap,
	)
	if err != nil {
		k.Logger(ctx).Error("error constructing unstake transaction", "err", err)
		return
	}

	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: 17000})
	if err != nil {
		k.Logger(ctx).Error("error creating metadata", "err", err)
		return
	}

	creator, err := k.getAddressByKeyID(ctx, k.GetZenBTCUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
	if err != nil {
		k.Logger(ctx).Error("error getting creator address", "err", err)
		return
	}

	if _, err := k.treasuryKeeper.HandleSignTransactionRequest(
		ctx,
		&treasurytypes.MsgNewSignTransactionRequest{
			Creator:             creator,
			KeyId:               k.GetZenBTCUnstakerKeyID(ctx),
			WalletType:          treasurytypes.WalletType_WALLET_TYPE_EVM,
			UnsignedTransaction: unsignedTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedTxHash)),
	); err != nil {
		k.Logger(ctx).Error("error creating unstake transaction", "err", err)
	}
}

func (k *Keeper) validateOracleData(voteExt VoteExtension, oracleData *OracleData) error {
	eigenDelegationsHash, err := deriveAVSContractStateHash(oracleData.EigenDelegationsMap)
	if err != nil {
		return fmt.Errorf("error deriving AVS contract delegation state hash: %w", err)
	}
	if !bytes.Equal(voteExt.EigenDelegationsHash, eigenDelegationsHash[:]) {
		return fmt.Errorf("AVS contract delegation state hash mismatch, expected %x, got %x", voteExt.EigenDelegationsHash, eigenDelegationsHash)
	}

	ethereumRedemptionsHash, err := deriveRedemptionsHash(oracleData.EthereumRedemptions)
	if err != nil {
		return fmt.Errorf("error deriving ethereum redemptions hash: %w", err)
	}
	if !bytes.Equal(voteExt.EthereumRedemptionsHash, ethereumRedemptionsHash[:]) {
		return fmt.Errorf("ethereum redemptions hash mismatch, expected %x, got %x", voteExt.EthereumRedemptionsHash, ethereumRedemptionsHash)
	}

	solanaRedemptionsHash, err := deriveRedemptionsHash(oracleData.SolanaRedemptions)
	if err != nil {
		return fmt.Errorf("error deriving solana redemptions hash: %w", err)
	}
	if !bytes.Equal(voteExt.SolanaRedemptionsHash, solanaRedemptionsHash[:]) {
		return fmt.Errorf("solana redemptions hash mismatch, expected %x, got %x", voteExt.SolanaRedemptionsHash, solanaRedemptionsHash)
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

	bitcoinHeaderHash, err := deriveBitcoinHeaderHash(&oracleData.BtcBlockHeader)
	if err != nil {
		return fmt.Errorf("error deriving bitcoin header hash: %w", err)
	}
	if !bytes.Equal(voteExt.BtcHeaderHash, bitcoinHeaderHash[:]) {
		return fmt.Errorf("bitcoin header hash mismatch, expected %x, got %x", voteExt.BtcHeaderHash, bitcoinHeaderHash)
	}

	if voteExt.RequestedEthMinterNonce != oracleData.RequestedEthMinterNonce {
		return fmt.Errorf("requested Ethereum nonce mismatch, expected %d, got %d", voteExt.RequestedEthMinterNonce, oracleData.RequestedEthMinterNonce)
	}

	if voteExt.RequestedEthUnstakerNonce != oracleData.RequestedEthUnstakerNonce {
		return fmt.Errorf("requested Ethereum nonce mismatch, expected %d, got %d", voteExt.RequestedEthUnstakerNonce, oracleData.RequestedEthUnstakerNonce)
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
