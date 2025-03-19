package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time" // Added time import for timing measurements

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
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
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("BeginBlocker timing", "duration_ms", elapsed.Milliseconds())
	}()

	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	trackStart := time.Now()
	err := k.TrackHistoricalInfo(ctx)
	trackElapsed := time.Since(trackStart)
	k.Logger(ctx).Warn("TrackHistoricalInfo timing", "duration_ms", trackElapsed.Milliseconds())

	return err
}

// EndBlocker calls telemetry and then processes validator updates.
func (k *Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("EndBlocker timing", "duration_ms", elapsed.Milliseconds())
	}()

	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	blockValUpdatesStart := time.Now()
	updates, err := k.BlockValidatorUpdates(ctx)
	blockValUpdatesElapsed := time.Since(blockValUpdatesStart)
	k.Logger(ctx).Warn("BlockValidatorUpdates timing", "duration_ms", blockValUpdatesElapsed.Milliseconds())

	return updates, err
}

//
// =============================================================================
// VOTE EXTENSION HANDLERS
// =============================================================================
//

// ExtendVoteHandler is called by all validators to extend the consensus vote
// with additional data to be voted on.
func (k *Keeper) ExtendVoteHandler(ctx context.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("ExtendVoteHandler timing", "height", req.Height, "duration_ms", elapsed.Milliseconds())
	}()

	getSidecarStart := time.Now()
	oracleData, err := k.GetSidecarState(ctx, req.Height)
	getSidecarElapsed := time.Since(getSidecarStart)
	k.Logger(ctx).Warn("GetSidecarState timing", "height", req.Height, "duration_ms", getSidecarElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error retrieving AVS delegations", "height", req.Height, "error", err)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	constructVEStart := time.Now()
	voteExt, err := k.constructVoteExtension(ctx, req.Height, oracleData)
	constructVEElapsed := time.Since(constructVEStart)
	k.Logger(ctx).Warn("constructVoteExtension timing", "height", req.Height, "duration_ms", constructVEElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error creating vote extension", "height", req.Height, "error", err)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	validateStart := time.Now()
	isInvalid := voteExt.IsInvalid(k.Logger(ctx))
	validateElapsed := time.Since(validateStart)
	k.Logger(ctx).Warn("VoteExt.IsInvalid timing", "height", req.Height, "duration_ms", validateElapsed.Milliseconds())

	if isInvalid {
		k.Logger(ctx).Error("invalid vote extension in ExtendVote", "height", req.Height)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	marshalStart := time.Now()
	voteExtBz, err := json.Marshal(voteExt)
	marshalElapsed := time.Since(marshalStart)
	k.Logger(ctx).Warn("json.Marshal voteExt timing", "height", req.Height, "duration_ms", marshalElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error marshalling vote extension", "height", req.Height, "error", err)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	return &abci.ResponseExtendVote{VoteExtension: voteExtBz}, nil
}

// constructVoteExtension builds the vote extension based on oracle data and on-chain state.
func (k *Keeper) constructVoteExtension(ctx context.Context, height int64, oracleData *OracleData) (VoteExtension, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("constructVoteExtension timing", "height", height, "duration_ms", elapsed.Milliseconds())
	}()

	avsDelegationsHashStart := time.Now()
	avsDelegationsHash, err := deriveHash(oracleData.EigenDelegationsMap)
	avsDelegationsHashElapsed := time.Since(avsDelegationsHashStart)
	k.Logger(ctx).Warn("deriveHash EigenDelegationsMap timing", "duration_ms", avsDelegationsHashElapsed.Milliseconds())

	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving AVS contract delegation state hash: %w", err)
	}

	ethBurnEventsHashStart := time.Now()
	ethBurnEventsHash, err := deriveHash(oracleData.EthBurnEvents)
	ethBurnEventsHashElapsed := time.Since(ethBurnEventsHashStart)
	k.Logger(ctx).Warn("deriveHash EthBurnEvents timing", "duration_ms", ethBurnEventsHashElapsed.Milliseconds())

	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving ethereum burn events hash: %w", err)
	}

	redemptionsHashStart := time.Now()
	redemptionsHash, err := deriveHash(oracleData.Redemptions)
	redemptionsHashElapsed := time.Since(redemptionsHashStart)
	k.Logger(ctx).Warn("deriveHash Redemptions timing", "duration_ms", redemptionsHashElapsed.Milliseconds())

	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving redemptions hash: %w", err)
	}

	btcHeadersStart := time.Now()
	latestHeader, requestedHeader, err := k.retrieveBitcoinHeaders(ctx)
	btcHeadersElapsed := time.Since(btcHeadersStart)
	k.Logger(ctx).Warn("retrieveBitcoinHeaders timing", "duration_ms", btcHeadersElapsed.Milliseconds())

	if err != nil {
		return VoteExtension{}, err
	}

	latestHeaderHashStart := time.Now()
	latestBitcoinHeaderHash, err := deriveHash(latestHeader.BlockHeader)
	latestHeaderHashElapsed := time.Since(latestHeaderHashStart)
	k.Logger(ctx).Warn("deriveHash latestHeader.BlockHeader timing", "duration_ms", latestHeaderHashElapsed.Milliseconds())

	if err != nil {
		return VoteExtension{}, err
	}

	// Only set requested header fields if there's a requested header
	requestedBtcBlockHeight := int64(0)
	var requestedBtcHeaderHash []byte
	if requestedHeader != nil {
		// Get requested Bitcoin header hash
		requestedHeaderHashStart := time.Now()
		requestedBitcoinHeaderHash, err := deriveHash(requestedHeader.BlockHeader)
		requestedHeaderHashElapsed := time.Since(requestedHeaderHashStart)
		k.Logger(ctx).Warn("deriveHash requestedHeader.BlockHeader timing", "duration_ms", requestedHeaderHashElapsed.Milliseconds())

		if err != nil {
			return VoteExtension{}, err
		}
		requestedBtcBlockHeight = requestedHeader.BlockHeight
		requestedBtcHeaderHash = requestedBitcoinHeaderHash[:]
	}

	// Check if Solana blockhash is requested first
	// solanaBlockhashRequested, err := k.SolanaBlockhashRequested.Get(ctx)
	// if err != nil {
	// 	if !errors.Is(err, collections.ErrNotFound) {
	// 		return VoteExtension{}, err
	// 	}
	// 	// Not found means false
	// 	solanaBlockhashRequested = false
	// }

	// Only get Solana recent blockhash if it's requested
	// solanaRecentBlockhash := ""
	// if solanaBlockhashRequested {
	// 	solanaRecentBlockhash, err = k.GetSolanaRecentBlockhash(ctx)
	// 	if err != nil {
	// 		k.Logger(ctx).Error("error getting Solana recent blockhash", "error", err)
	// 		// Non-fatal error, continue with empty string
	// 		solanaRecentBlockhash = ""
	// 	}
	// }

	getNoncesStart := time.Now()
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
	getNoncesElapsed := time.Since(getNoncesStart)
	k.Logger(ctx).Warn("Get nonces timing", "duration_ms", getNoncesElapsed.Milliseconds())

	voteExt := VoteExtension{
		ZRChainBlockHeight:         height,
		ROCKUSDPrice:               oracleData.ROCKUSDPrice,
		BTCUSDPrice:                oracleData.BTCUSDPrice,
		ETHUSDPrice:                oracleData.ETHUSDPrice,
		EigenDelegationsHash:       avsDelegationsHash[:],
		EthBurnEventsHash:          ethBurnEventsHash[:],
		RedemptionsHash:            redemptionsHash[:],
		RequestedBtcBlockHeight:    requestedBtcBlockHeight,
		RequestedBtcHeaderHash:     requestedBtcHeaderHash,
		LatestBtcBlockHeight:       latestHeader.BlockHeight,
		LatestBtcHeaderHash:        latestBitcoinHeaderHash[:],
		EthBlockHeight:             oracleData.EthBlockHeight,
		EthGasLimit:                oracleData.EthGasLimit,
		EthBaseFee:                 oracleData.EthBaseFee,
		EthTipCap:                  oracleData.EthTipCap,
		SolanaLamportsPerSignature: oracleData.SolanaLamportsPerSignature,
		// SolanaRecentBlockhash:      solanaRecentBlockhash,
		RequestedStakerNonce:    nonces[k.zenBTCKeeper.GetStakerKeyID(ctx)],
		RequestedEthMinterNonce: nonces[k.zenBTCKeeper.GetEthMinterKeyID(ctx)],
		RequestedUnstakerNonce:  nonces[k.zenBTCKeeper.GetUnstakerKeyID(ctx)],
		RequestedCompleterNonce: nonces[k.zenBTCKeeper.GetCompleterKeyID(ctx)],
	}

	return voteExt, nil
}

// VerifyVoteExtensionHandler is called by all validators to verify vote extension data.
func (k *Keeper) VerifyVoteExtensionHandler(ctx context.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("VerifyVoteExtensionHandler timing", "height", req.Height, "duration_ms", elapsed.Milliseconds())
	}()

	if len(req.VoteExtension) == 0 {
		return ACCEPT_VOTE, nil
	}

	if len(req.VoteExtension) > VoteExtBytesLimit {
		k.Logger(ctx).Error("vote extension is too large", "height", req.Height, "limit", VoteExtBytesLimit, "size", len(req.VoteExtension))
		return REJECT_VOTE, nil
	}

	unmarshalStart := time.Now()
	var voteExt VoteExtension
	if err := json.Unmarshal(req.VoteExtension, &voteExt); err != nil {
		k.Logger(ctx).Debug("error unmarshalling vote extension", "height", req.Height, "error", err)
		return REJECT_VOTE, nil
	}
	unmarshalElapsed := time.Since(unmarshalStart)
	k.Logger(ctx).Warn("json.Unmarshal voteExt timing", "height", req.Height, "duration_ms", unmarshalElapsed.Milliseconds())

	if req.Height != voteExt.ZRChainBlockHeight {
		k.Logger(ctx).Error("mismatched height for vote extension", "expected", req.Height, "got", voteExt.ZRChainBlockHeight)
		return REJECT_VOTE, nil
	}

	validateStart := time.Now()
	isInvalid := voteExt.IsInvalid(k.Logger(ctx))
	validateElapsed := time.Since(validateStart)
	k.Logger(ctx).Warn("VoteExt.IsInvalid timing", "height", req.Height, "duration_ms", validateElapsed.Milliseconds())

	if isInvalid {
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
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("PrepareProposal timing", "height", req.Height, "duration_ms", elapsed.Milliseconds())
	}()

	if !VoteExtensionsEnabled(ctx) {
		k.Logger(ctx).Debug("vote extensions disabled; not injecting oracle data", "height", req.Height)
		return nil, nil
	}

	getSuperMajorityStart := time.Now()
	voteExt, fieldVotePowers, err := k.GetSuperMajorityVEData(ctx, req.Height, req.LocalLastCommit)
	getSuperMajorityElapsed := time.Since(getSuperMajorityStart)
	k.Logger(ctx).Warn("GetSuperMajorityVEData timing", "height", req.Height, "duration_ms", getSuperMajorityElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error retrieving supermajority vote extension data", "height", req.Height, "error", err)
		return nil, nil
	}

	if len(fieldVotePowers) == 0 { // no field reached consensus
		k.Logger(ctx).Warn("no fields reached consensus in vote extension", "height", req.Height)
		marshalStart := time.Now()
		result, err := k.marshalOracleData(req, &OracleData{ConsensusData: req.LocalLastCommit, FieldVotePowers: fieldVotePowers})
		marshalElapsed := time.Since(marshalStart)
		k.Logger(ctx).Warn("marshalOracleData (no consensus) timing", "height", req.Height, "duration_ms", marshalElapsed.Milliseconds())
		return result, err
	}

	if voteExt.ZRChainBlockHeight != req.Height-1 { // vote extension is from previous block
		k.Logger(ctx).Error("mismatched height for vote extension", "height", req.Height, "voteExt.ZRChainBlockHeight", voteExt.ZRChainBlockHeight)
		return nil, nil
	}

	getValidatedStart := time.Now()
	oracleData, err := k.getValidatedOracleData(ctx, voteExt, fieldVotePowers)
	getValidatedElapsed := time.Since(getValidatedStart)
	k.Logger(ctx).Warn("getValidatedOracleData timing", "height", req.Height, "duration_ms", getValidatedElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Warn("error in getValidatedOracleData; injecting empty oracle data", "height", req.Height, "error", err)
		oracleData = &OracleData{}
	}

	oracleData.ConsensusData = req.LocalLastCommit

	marshalStart := time.Now()
	result, err := k.marshalOracleData(req, oracleData)
	marshalElapsed := time.Since(marshalStart)
	k.Logger(ctx).Warn("marshalOracleData timing", "height", req.Height, "duration_ms", marshalElapsed.Milliseconds())

	return result, err
}

// ProcessProposal is executed by all validators to check whether the proposer prepared valid data.
func (k *Keeper) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("ProcessProposal timing", "height", req.Height, "duration_ms", elapsed.Milliseconds())
	}()

	// Return early if this node is not a validator so non-validators don't need to be running a sidecar
	if !k.zrConfig.IsValidator {
		return ACCEPT_PROPOSAL, nil
	}

	if !VoteExtensionsEnabled(ctx) || len(req.Txs) == 0 {
		return ACCEPT_PROPOSAL, nil
	}

	containsVEStart := time.Now()
	containsVE := ContainsVoteExtension(req.Txs[0], k.txDecoder)
	containsVEElapsed := time.Since(containsVEStart)
	k.Logger(ctx).Warn("ContainsVoteExtension timing", "height", req.Height, "duration_ms", containsVEElapsed.Milliseconds())

	if !containsVE {
		k.Logger(ctx).Warn("block does not contain vote extensions, rejecting proposal")
		return REJECT_PROPOSAL, nil
	}

	unmarshalStart := time.Now()
	var recoveredOracleData OracleData
	if err := json.Unmarshal(req.Txs[0], &recoveredOracleData); err != nil {
		return REJECT_PROPOSAL, fmt.Errorf("error unmarshalling oracle data: %w", err)
	}
	unmarshalElapsed := time.Since(unmarshalStart)
	k.Logger(ctx).Warn("json.Unmarshal oracleData timing", "height", req.Height, "duration_ms", unmarshalElapsed.Milliseconds())

	// Check for empty oracle data - if it's empty, accept the proposal
	recoveredOracleDataNoCommitInfo := recoveredOracleData
	recoveredOracleDataNoCommitInfo.ConsensusData = abci.ExtendedCommitInfo{}
	recoveredOracleDataNoCommitInfo.FieldVotePowers = nil
	if reflect.DeepEqual(recoveredOracleDataNoCommitInfo, OracleData{}) {
		k.Logger(ctx).Warn("accepting empty oracle data", "height", req.Height)
		return ACCEPT_PROPOSAL, nil
	}

	validateVEStart := time.Now()
	if err := ValidateVoteExtensions(ctx, k, req.Height, ctx.ChainID(), recoveredOracleData.ConsensusData); err != nil {
		k.Logger(ctx).Error("error validating vote extensions", "height", req.Height, "error", err)
		return REJECT_PROPOSAL, err
	}
	validateVEElapsed := time.Since(validateVEStart)
	k.Logger(ctx).Warn("ValidateVoteExtensions timing", "height", req.Height, "duration_ms", validateVEElapsed.Milliseconds())

	return ACCEPT_PROPOSAL, nil
}

//
// =============================================================================
// PRE-BLOCKER: ORACLE DATA PROCESSING
// =============================================================================
//

// PreBlocker processes oracle data and applies the resulting state updates.
func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("PreBlocker timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	shouldProcessStart := time.Now()
	if !k.shouldProcessOracleData(ctx, req) {
		shouldProcessElapsed := time.Since(shouldProcessStart)
		k.Logger(ctx).Warn("shouldProcessOracleData timing (early exit)", "height", ctx.BlockHeight(), "duration_ms", shouldProcessElapsed.Milliseconds())
		return nil
	}
	shouldProcessElapsed := time.Since(shouldProcessStart)
	k.Logger(ctx).Warn("shouldProcessOracleData timing", "height", ctx.BlockHeight(), "duration_ms", shouldProcessElapsed.Milliseconds())

	unmarshalStart := time.Now()
	oracleData, ok := k.unmarshalOracleData(ctx, req.Txs[0])
	unmarshalElapsed := time.Since(unmarshalStart)
	k.Logger(ctx).Warn("unmarshalOracleData timing", "height", ctx.BlockHeight(), "success", ok, "duration_ms", unmarshalElapsed.Milliseconds())

	if !ok {
		return nil
	}

	validateCanonicalStart := time.Now()
	canonicalVE, ok := k.validateCanonicalVE(ctx, req.Height, oracleData)
	validateCanonicalElapsed := time.Since(validateCanonicalStart)
	k.Logger(ctx).Warn("validateCanonicalVE timing", "height", ctx.BlockHeight(), "success", ok, "duration_ms", validateCanonicalElapsed.Milliseconds())

	if !ok {
		k.Logger(ctx).Error("invalid canonical vote extension")
		return nil
	}

	// Update asset prices if there's consensus on the price fields
	updatePricesStart := time.Now()
	k.updateAssetPrices(ctx, oracleData)
	updatePricesElapsed := time.Since(updatePricesStart)
	k.Logger(ctx).Warn("updateAssetPrices timing", "height", ctx.BlockHeight(), "duration_ms", updatePricesElapsed.Milliseconds())

	// Process different subsystems based on field consensus

	// Validator updates - only if EigenDelegationsHash has consensus
	if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldEigenDelegationsHash) {
		updateValidatorStakesStart := time.Now()
		k.updateValidatorStakes(ctx, oracleData)
		updateValidatorStakesElapsed := time.Since(updateValidatorStakesStart)
		k.Logger(ctx).Warn("updateValidatorStakes timing", "height", ctx.BlockHeight(), "duration_ms", updateValidatorStakesElapsed.Milliseconds())

		updateAVSDelegationStart := time.Now()
		k.updateAVSDelegationStore(ctx, oracleData)
		updateAVSDelegationElapsed := time.Since(updateAVSDelegationStart)
		k.Logger(ctx).Warn("updateAVSDelegationStore timing", "height", ctx.BlockHeight(), "duration_ms", updateAVSDelegationElapsed.Milliseconds())
	}

	// Bitcoin header processing - only if BTC header fields have consensus
	btcHeaderFields := []VoteExtensionField{VEFieldLatestBtcHeaderHash, VEFieldRequestedBtcHeaderHash}
	if anyFieldHasConsensus(oracleData.FieldVotePowers, btcHeaderFields) {
		storeBtcHeadersStart := time.Now()
		if err := k.storeBitcoinBlockHeaders(ctx, oracleData); err != nil {
			k.Logger(ctx).Error("error storing Bitcoin headers", "error", err)
		}
		storeBtcHeadersElapsed := time.Since(storeBtcHeadersStart)
		k.Logger(ctx).Warn("storeBitcoinBlockHeaders timing", "height", ctx.BlockHeight(), "duration_ms", storeBtcHeadersElapsed.Milliseconds())
	}

	if ctx.BlockHeight()%2 == 0 { // TODO: is this needed?
		nonceFields := []VoteExtensionField{
			VEFieldRequestedStakerNonce,
			VEFieldRequestedEthMinterNonce,
			VEFieldRequestedUnstakerNonce,
			VEFieldRequestedCompleterNonce,
		}
		if anyFieldHasConsensus(oracleData.FieldVotePowers, nonceFields) {
			updateNoncesStart := time.Now()
			k.updateNonces(ctx, oracleData)
			updateNoncesElapsed := time.Since(updateNoncesStart)
			k.Logger(ctx).Warn("updateNonces timing", "height", ctx.BlockHeight(), "duration_ms", updateNoncesElapsed.Milliseconds())
		}

		if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldEthBurnEventsHash) {
			storeZenBtcBurnStart := time.Now()
			k.storeNewZenBTCBurnEventsEthereum(ctx, oracleData)
			storeZenBtcBurnElapsed := time.Since(storeZenBtcBurnStart)
			k.Logger(ctx).Warn("storeNewZenBTCBurnEventsEthereum timing", "height", ctx.BlockHeight(), "duration_ms", storeZenBtcBurnElapsed.Milliseconds())
		}
		if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldRedemptionsHash) {
			storeRedemptionsStart := time.Now()
			k.storeNewZenBTCRedemptions(ctx, oracleData)
			storeRedemptionsElapsed := time.Since(storeRedemptionsStart)
			k.Logger(ctx).Warn("storeNewZenBTCRedemptions timing", "height", ctx.BlockHeight(), "duration_ms", storeRedemptionsElapsed.Milliseconds())
		}

		processZenBTCStakingStart := time.Now()
		k.processZenBTCStaking(ctx, oracleData)
		processZenBTCStakingElapsed := time.Since(processZenBTCStakingStart)
		k.Logger(ctx).Warn("processZenBTCStaking timing", "height", ctx.BlockHeight(), "duration_ms", processZenBTCStakingElapsed.Milliseconds())

		processZenBTCMintsStart := time.Now()
		k.processZenBTCMintsEthereum(ctx, oracleData)
		processZenBTCMintsElapsed := time.Since(processZenBTCMintsStart)
		k.Logger(ctx).Warn("processZenBTCMintsEthereum timing", "height", ctx.BlockHeight(), "duration_ms", processZenBTCMintsElapsed.Milliseconds())

		processZenBTCBurnStart := time.Now()
		k.processZenBTCBurnEventsEthereum(ctx, oracleData)
		processZenBTCBurnElapsed := time.Since(processZenBTCBurnStart)
		k.Logger(ctx).Warn("processZenBTCBurnEventsEthereum timing", "height", ctx.BlockHeight(), "duration_ms", processZenBTCBurnElapsed.Milliseconds())

		processRedemptionsStart := time.Now()
		k.processZenBTCRedemptions(ctx, oracleData)
		processRedemptionsElapsed := time.Since(processRedemptionsStart)
		k.Logger(ctx).Warn("processZenBTCRedemptions timing", "height", ctx.BlockHeight(), "duration_ms", processRedemptionsElapsed.Milliseconds())

		checkRedemptionsStart := time.Now()
		k.checkForRedemptionFulfilment(ctx)
		checkRedemptionsElapsed := time.Since(checkRedemptionsStart)
		k.Logger(ctx).Warn("checkForRedemptionFulfilment timing", "height", ctx.BlockHeight(), "duration_ms", checkRedemptionsElapsed.Milliseconds())
	}

	recordNonVotingStart := time.Now()
	k.recordNonVotingValidators(ctx, req)
	recordNonVotingElapsed := time.Since(recordNonVotingStart)
	k.Logger(ctx).Warn("recordNonVotingValidators timing", "height", ctx.BlockHeight(), "duration_ms", recordNonVotingElapsed.Milliseconds())

	recordMismatchedStart := time.Now()
	k.recordMismatchedVoteExtensions(ctx, req.Height, canonicalVE, oracleData.ConsensusData)
	recordMismatchedElapsed := time.Since(recordMismatchedStart)
	k.Logger(ctx).Warn("recordMismatchedVoteExtensions timing", "height", ctx.BlockHeight(), "duration_ms", recordMismatchedElapsed.Milliseconds())

	return nil
}

// shouldProcessOracleData checks if oracle data should be processed for this block.
func (k *Keeper) shouldProcessOracleData(ctx sdk.Context, req *abci.RequestFinalizeBlock) bool {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("shouldProcessOracleData timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	if len(req.Txs) == 0 {
		k.Logger(ctx).Debug("no transactions in block")
		return false
	}

	if req.Height == 1 || !VoteExtensionsEnabled(ctx) {
		k.Logger(ctx).Debug("vote extensions not enabled for this block", "height", req.Height)
		return false
	}

	containsVEStart := time.Now()
	containsVE := ContainsVoteExtension(req.Txs[0], k.txDecoder)
	containsVEElapsed := time.Since(containsVEStart)
	k.Logger(ctx).Warn("ContainsVoteExtension timing", "height", ctx.BlockHeight(), "duration_ms", containsVEElapsed.Milliseconds())

	if !containsVE {
		k.Logger(ctx).Debug("first transaction does not contain vote extension", "height", req.Height)
		return false
	}

	return true
}

// validateCanonicalVE validates the canonical vote extension from oracle data.
func (k *Keeper) validateCanonicalVE(ctx sdk.Context, height int64, oracleData OracleData) (VoteExtension, bool) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("validateCanonicalVE timing", "height", height, "duration_ms", elapsed.Milliseconds())
	}()

	getSuperMajorityStart := time.Now()
	voteExt, fieldVotePowers, err := k.GetSuperMajorityVEData(ctx, height, oracleData.ConsensusData)
	getSuperMajorityElapsed := time.Since(getSuperMajorityStart)
	k.Logger(ctx).Warn("GetSuperMajorityVEData timing", "height", height, "duration_ms", getSuperMajorityElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error getting super majority VE data", "height", height, "error", err)
		return VoteExtension{}, false
	}

	if reflect.DeepEqual(voteExt, VoteExtension{}) {
		k.Logger(ctx).Warn("accepting empty vote extension", "height", height)
		return VoteExtension{}, true
	}

	validateOracleDataStart := time.Now()
	k.validateOracleData(ctx, voteExt, &oracleData, fieldVotePowers)
	validateOracleDataElapsed := time.Since(validateOracleDataStart)
	k.Logger(ctx).Warn("validateOracleData timing", "height", height, "duration_ms", validateOracleDataElapsed.Milliseconds())

	// Log final consensus summary after validation
	k.Logger(ctx).Info("final consensus summary",
		"fields_with_consensus", len(oracleData.FieldVotePowers),
		"stage", "post_validation")

	return voteExt, true
}

// getValidatedOracleData retrieves and validates oracle data based on a vote extension.
// Only validates fields that have reached consensus as indicated in fieldVotePowers.
func (k *Keeper) getValidatedOracleData(ctx sdk.Context, voteExt VoteExtension, fieldVotePowers map[VoteExtensionField]int64) (*OracleData, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("getValidatedOracleData timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	// We only fetch Ethereum state if we have consensus on EthBlockHeight
	var oracleData *OracleData
	var err error

	if _, ok := fieldVotePowers[VEFieldEthBlockHeight]; ok {
		getSidecarByEthHeightStart := time.Now()
		oracleData, err = k.GetSidecarStateByEthHeight(ctx, voteExt.EthBlockHeight)
		getSidecarByEthHeightElapsed := time.Since(getSidecarByEthHeightStart)
		k.Logger(ctx).Warn("GetSidecarStateByEthHeight timing", "eth_height", voteExt.EthBlockHeight, "duration_ms", getSidecarByEthHeightElapsed.Milliseconds())

		if err != nil {
			return nil, fmt.Errorf("error fetching oracle state: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no consensus on eth block height")
	}

	retrieveBtcHeadersStart := time.Now()
	latestHeader, requestedHeader, err := k.retrieveBitcoinHeaders(ctx)
	retrieveBtcHeadersElapsed := time.Since(retrieveBtcHeadersStart)
	k.Logger(ctx).Warn("retrieveBitcoinHeaders timing", "duration_ms", retrieveBtcHeadersElapsed.Milliseconds())

	if err != nil {
		return nil, fmt.Errorf("error fetching bitcoin headers: %w", err)
	}

	// Collect fields that fail validation to revoke consensus
	mismatchedFields := make([]VoteExtensionField, 0)

	// Copy latest Bitcoin header data if we have consensus on both height and hash fields
	if fieldHasConsensus(fieldVotePowers, VEFieldLatestBtcBlockHeight) &&
		fieldHasConsensus(fieldVotePowers, VEFieldLatestBtcHeaderHash) &&
		latestHeader != nil {
		oracleData.LatestBtcBlockHeight = latestHeader.BlockHeight
		oracleData.LatestBtcBlockHeader = *latestHeader.BlockHeader
	}

	// Copy requested Bitcoin header data if we have consensus on both height and hash fields
	if fieldHasConsensus(fieldVotePowers, VEFieldRequestedBtcBlockHeight) &&
		fieldHasConsensus(fieldVotePowers, VEFieldRequestedBtcHeaderHash) &&
		requestedHeader != nil {
		oracleData.RequestedBtcBlockHeight = requestedHeader.BlockHeight
		oracleData.RequestedBtcBlockHeader = *requestedHeader.BlockHeader
	}

	// Verify Solana recent blockhash if there's consensus on it
	// if fieldHasConsensus(fieldVotePowers, VEFieldSolanaRecentBlockhash) {
	// 	currentBlockhash, err := k.GetSolanaRecentBlockhash(ctx)
	// 	if err != nil {
	// 		k.Logger(ctx).Error("error getting Solana recent blockhash for validation", "error", err)
	// 		// Skip the rest of this validation block on error
	// 	} else if currentBlockhash != "" && voteExt.SolanaRecentBlockhash != currentBlockhash {
	// 		// Check for mismatch only if we have a valid current blockhash
	// 		k.Logger(ctx).Warn("solana recent blockhash mismatch",
	// 			"voteExt", voteExt.SolanaRecentBlockhash,
	// 			"current", currentBlockhash)
	// 		mismatchedFields = append(mismatchedFields, VEFieldSolanaRecentBlockhash)
	// 	} else {
	// 		// No mismatch or empty current blockhash, use the consensus value
	// 		oracleData.SolanaRecentBlockhash = voteExt.SolanaRecentBlockhash
	// 	}
	// }

	// Verify nonce fields and copy them if they have consensus
	verifyNonceFieldsStart := time.Now()
	nonceFields := []struct {
		field       VoteExtensionField
		keyID       uint64
		voteExtVal  uint64
		oracleField *uint64
	}{
		{VEFieldRequestedStakerNonce, k.zenBTCKeeper.GetStakerKeyID(ctx), voteExt.RequestedStakerNonce, &oracleData.RequestedStakerNonce},
		{VEFieldRequestedEthMinterNonce, k.zenBTCKeeper.GetEthMinterKeyID(ctx), voteExt.RequestedEthMinterNonce, &oracleData.RequestedEthMinterNonce},
		{VEFieldRequestedUnstakerNonce, k.zenBTCKeeper.GetUnstakerKeyID(ctx), voteExt.RequestedUnstakerNonce, &oracleData.RequestedUnstakerNonce},
		{VEFieldRequestedCompleterNonce, k.zenBTCKeeper.GetCompleterKeyID(ctx), voteExt.RequestedCompleterNonce, &oracleData.RequestedCompleterNonce},
	}

	for _, nf := range nonceFields {
		if fieldHasConsensus(fieldVotePowers, nf.field) {
			// Also verify nonce against what would be fetched
			requested, err := k.EthereumNonceRequested.Get(ctx, nf.keyID)
			if err != nil {
				if !errors.Is(err, collections.ErrNotFound) {
					k.Logger(ctx).Error("error checking nonce request state", "keyID", nf.keyID, "error", err)
				}
			} else if requested {
				currentNonce, err := k.lookupEthereumNonce(ctx, nf.keyID)
				if err != nil {
					k.Logger(ctx).Error("error looking up Ethereum nonce for validation", "keyID", nf.keyID, "error", err)
				} else if currentNonce != nf.voteExtVal && nf.voteExtVal != 0 {
					k.Logger(ctx).Warn("nonce mismatch for key",
						"keyID", nf.keyID,
						"voteExt", nf.voteExtVal,
						"current", currentNonce)
					mismatchedFields = append(mismatchedFields, nf.field)
				}
			}

			*nf.oracleField = nf.voteExtVal
		}
	}
	verifyNonceFieldsElapsed := time.Since(verifyNonceFieldsStart)
	k.Logger(ctx).Warn("verify nonce fields timing", "duration_ms", verifyNonceFieldsElapsed.Milliseconds())

	// Process mismatched fields
	nullifyMismatchedFieldsStart := time.Now()
	k.nullifyMismatchedFields(ctx, mismatchedFields, fieldVotePowers, oracleData)
	nullifyMismatchedFieldsElapsed := time.Since(nullifyMismatchedFieldsStart)
	k.Logger(ctx).Warn("nullifyMismatchedFields timing", "duration_ms", nullifyMismatchedFieldsElapsed.Milliseconds())

	// Call the standard validateOracleData to check other fields
	validateOracleDataStart := time.Now()
	k.validateOracleData(ctx, voteExt, oracleData, fieldVotePowers)
	validateOracleDataElapsed := time.Since(validateOracleDataStart)
	k.Logger(ctx).Warn("validateOracleData timing", "duration_ms", validateOracleDataElapsed.Milliseconds())

	return oracleData, nil
}

//
// =============================================================================
// VALIDATOR & DELEGATION STATE UPDATES
// =============================================================================
//

// updateValidatorStakes updates validator stake values and delegation mappings.
func (k *Keeper) updateValidatorStakes(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("updateValidatorStakes timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

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

		getValidatorStart := time.Now()
		validator, err := k.GetZenrockValidator(ctx, valAddr)
		getValidatorElapsed := time.Since(getValidatorStart)
		k.Logger(ctx).Warn("GetZenrockValidator timing", "validator", delegation.Validator, "duration_ms", getValidatorElapsed.Milliseconds())

		if err != nil || validator.Status != types.Bonded {
			k.Logger(ctx).Debug("invalid delegation for", "validator", delegation.Validator, "error", err)
			continue
		}

		validator.TokensAVS = sdkmath.Int(delegation.Stake)

		setValidatorStart := time.Now()
		if err = k.SetValidator(ctx, validator); err != nil {
			k.Logger(ctx).Error("error setting validator", "validator", delegation.Validator, "error", err)
			continue
		}
		setValidatorElapsed := time.Since(setValidatorStart)
		k.Logger(ctx).Warn("SetValidator timing", "validator", delegation.Validator, "duration_ms", setValidatorElapsed.Milliseconds())

		setDelegationStart := time.Now()
		if err = k.ValidatorDelegations.Set(ctx, valAddr.String(), delegation.Stake); err != nil {
			k.Logger(ctx).Error("error setting validator delegations", "validator", delegation.Validator, "error", err)
			continue
		}
		setDelegationElapsed := time.Since(setDelegationStart)
		k.Logger(ctx).Warn("ValidatorDelegations.Set timing", "validator", delegation.Validator, "duration_ms", setDelegationElapsed.Milliseconds())

		validatorInAVSDelegationSet[valAddr.String()] = true
	}

	removeStaleStart := time.Now()
	k.removeStaleValidatorDelegations(ctx, validatorInAVSDelegationSet)
	removeStaleElapsed := time.Since(removeStaleStart)
	k.Logger(ctx).Warn("removeStaleValidatorDelegations timing", "duration_ms", removeStaleElapsed.Milliseconds())
}

// removeStaleValidatorDelegations removes delegation entries for validators not present in the current AVS data.
func (k *Keeper) removeStaleValidatorDelegations(ctx sdk.Context, validatorInAVSDelegationSet map[string]bool) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("removeStaleValidatorDelegations timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	var validatorsToRemove []string

	validatorWalkStart := time.Now()
	if err := k.ValidatorDelegations.Walk(ctx, nil, func(valAddr string, stake sdkmath.Int) (bool, error) {
		if !validatorInAVSDelegationSet[valAddr] {
			validatorsToRemove = append(validatorsToRemove, valAddr)
		}
		return true, nil
	}); err != nil {
		k.Logger(ctx).Error("error walking validator delegations", "error", err)
	}
	validatorWalkElapsed := time.Since(validatorWalkStart)
	k.Logger(ctx).Warn("ValidatorDelegations.Walk timing", "duration_ms", validatorWalkElapsed.Milliseconds())

	for _, valAddr := range validatorsToRemove {
		removeStart := time.Now()
		if err := k.ValidatorDelegations.Remove(ctx, valAddr); err != nil {
			k.Logger(ctx).Error("error removing validator delegation", "validator", valAddr, "error", err)
			continue
		}
		removeElapsed := time.Since(removeStart)
		k.Logger(ctx).Warn("ValidatorDelegations.Remove timing", "validator", valAddr, "duration_ms", removeElapsed.Milliseconds())

		updateTokensStart := time.Now()
		if err := k.updateValidatorTokensAVS(ctx, valAddr); err != nil {
			k.Logger(ctx).Error("error updating validator TokensAVS", "validator", valAddr, "error", err)
		}
		updateTokensElapsed := time.Since(updateTokensStart)
		k.Logger(ctx).Warn("updateValidatorTokensAVS timing", "validator", valAddr, "duration_ms", updateTokensElapsed.Milliseconds())
	}
}

// updateValidatorTokensAVS resets a validator's AVS tokens to zero.
func (k *Keeper) updateValidatorTokensAVS(ctx sdk.Context, valAddr string) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("updateValidatorTokensAVS timing", "validator", valAddr, "duration_ms", elapsed.Milliseconds())
	}()

	getValidatorStart := time.Now()
	validator, err := k.GetZenrockValidator(ctx, sdk.ValAddress(valAddr))
	getValidatorElapsed := time.Since(getValidatorStart)
	k.Logger(ctx).Warn("GetZenrockValidator timing", "validator", valAddr, "duration_ms", getValidatorElapsed.Milliseconds())

	if err != nil {
		return fmt.Errorf("error retrieving validator for removal: %w", err)
	}

	validator.TokensAVS = sdkmath.ZeroInt()

	setValidatorStart := time.Now()
	if err = k.SetValidator(ctx, validator); err != nil {
		return fmt.Errorf("error updating validator after removal: %w", err)
	}
	setValidatorElapsed := time.Since(setValidatorStart)
	k.Logger(ctx).Warn("SetValidator timing", "validator", valAddr, "duration_ms", setValidatorElapsed.Milliseconds())

	return nil
}

// updateAVSDelegationStore updates the AVS delegation store with new delegation amounts.
func (k *Keeper) updateAVSDelegationStore(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("updateAVSDelegationStore timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	for validatorAddr, delegatorMap := range oracleData.EigenDelegationsMap {
		for delegatorAddr, amount := range delegatorMap {
			setDelegationStart := time.Now()
			if err := k.AVSDelegations.Set(ctx, collections.Join(validatorAddr, delegatorAddr), sdkmath.NewIntFromBigInt(amount)); err != nil {
				k.Logger(ctx).Error("error setting AVS delegations", "error", err)
			}
			setDelegationElapsed := time.Since(setDelegationStart)
			k.Logger(ctx).Warn("AVSDelegations.Set timing", "validator", validatorAddr, "delegator", delegatorAddr, "duration_ms", setDelegationElapsed.Milliseconds())
		}
	}
}

//
// =============================================================================
// BITCOIN HEADER PROCESSING
// =============================================================================
//

// storeBitcoinBlockHeader stores the Bitcoin header and handles historical header requests.
func (k *Keeper) storeBitcoinBlockHeaders(ctx sdk.Context, oracleData OracleData) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("storeBitcoinBlockHeaders timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	// First store the latest Bitcoin header if available
	if oracleData.LatestBtcBlockHeight > 0 && oracleData.LatestBtcBlockHeader.MerkleRoot != "" {
		checkBtcHeaderStart := time.Now()
		latestHeaderExists, err := k.BtcBlockHeaders.Has(ctx, oracleData.LatestBtcBlockHeight)
		checkBtcHeaderElapsed := time.Since(checkBtcHeaderStart)
		k.Logger(ctx).Warn("BtcBlockHeaders.Has timing", "height", oracleData.LatestBtcBlockHeight, "duration_ms", checkBtcHeaderElapsed.Milliseconds())

		if err != nil {
			k.Logger(ctx).Error("error checking if latest Bitcoin header exists", "height", oracleData.LatestBtcBlockHeight, "error", err)
		} else if !latestHeaderExists {
			// Only store if it doesn't already exist
			setBtcHeaderStart := time.Now()
			if err := k.BtcBlockHeaders.Set(ctx, oracleData.LatestBtcBlockHeight, oracleData.LatestBtcBlockHeader); err != nil {
				k.Logger(ctx).Error("error storing latest Bitcoin header", "height", oracleData.LatestBtcBlockHeight, "error", err)
			} else {
				k.Logger(ctx).Info("stored latest Bitcoin header", "height", oracleData.LatestBtcBlockHeight)
			}
			setBtcHeaderElapsed := time.Since(setBtcHeaderStart)
			k.Logger(ctx).Warn("BtcBlockHeaders.Set (latest) timing", "height", oracleData.LatestBtcBlockHeight, "duration_ms", setBtcHeaderElapsed.Milliseconds())
		}
	}

	// Process the requested Bitcoin header
	headerHeight := oracleData.RequestedBtcBlockHeight
	if headerHeight == 0 || oracleData.RequestedBtcBlockHeader.MerkleRoot == "" {
		k.Logger(ctx).Debug("no requested bitcoin header")
		return nil
	}

	// Get requested headers
	getRequestedHeadersStart := time.Now()
	requestedHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	getRequestedHeadersElapsed := time.Since(getRequestedHeadersStart)
	k.Logger(ctx).Warn("RequestedHistoricalBitcoinHeaders.Get timing", "duration_ms", getRequestedHeadersElapsed.Milliseconds())

	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting requested historical Bitcoin headers", "error", err)
			return err
		}
		k.Logger(ctx).Info("requested historical Bitcoin headers store not initialised", "height", headerHeight)
	}

	// Check if the header is historical
	isHistoricalStart := time.Now()
	isHistorical := k.isHistoricalHeader(headerHeight, requestedHeaders.Heights)
	isHistoricalElapsed := time.Since(isHistoricalStart)
	k.Logger(ctx).Warn("isHistoricalHeader timing", "height", headerHeight, "duration_ms", isHistoricalElapsed.Milliseconds())

	// Check if header exists (for logging only)
	hasHeaderStart := time.Now()
	headerExists, err := k.BtcBlockHeaders.Has(ctx, headerHeight)
	hasHeaderElapsed := time.Since(hasHeaderStart)
	k.Logger(ctx).Warn("BtcBlockHeaders.Has timing", "height", headerHeight, "duration_ms", hasHeaderElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error checking if Bitcoin header exists", "height", headerHeight, "error", err)
		return err
	}

	logger := k.Logger(ctx).With(
		"height", headerHeight,
		"is_historical", isHistorical,
		"already_exists", headerExists,
		"requested_headers", requestedHeaders.Heights)

	// Always store the header regardless of whether it exists
	setHeaderStart := time.Now()
	if err := k.BtcBlockHeaders.Set(ctx, headerHeight, oracleData.RequestedBtcBlockHeader); err != nil {
		k.Logger(ctx).Error("error storing Bitcoin header", "height", headerHeight, "error", err)
		return err
	}
	setHeaderElapsed := time.Since(setHeaderStart)
	k.Logger(ctx).Warn("BtcBlockHeaders.Set timing", "height", headerHeight, "duration_ms", setHeaderElapsed.Milliseconds())

	logger.Info("stored Bitcoin header",
		"type", map[bool]string{true: "historical", false: "latest"}[isHistorical])

	// Process according to header type
	if isHistorical {
		// Remove the processed historical header from the requested list
		processHistoricalStart := time.Now()
		requestedHeaders.Heights = slices.DeleteFunc(requestedHeaders.Heights, func(height int64) bool {
			return height == headerHeight
		})

		if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
			k.Logger(ctx).Error("error updating requested historical Bitcoin headers", "error", err)
			return err
		}
		processHistoricalElapsed := time.Since(processHistoricalStart)
		k.Logger(ctx).Warn("Process historical header timing", "height", headerHeight, "duration_ms", processHistoricalElapsed.Milliseconds())

		logger.Debug("removed processed historical header request",
			"remaining_requests", len(requestedHeaders.Heights))
	} else if !headerExists {
		// Only check for reorgs for non-historical headers that weren't already stored
		checkReorgStart := time.Now()
		if err := k.checkForBitcoinReorg(ctx, oracleData, requestedHeaders); err != nil {
			k.Logger(ctx).Error("error handling potential Bitcoin reorg", "height", headerHeight, "error", err)
		}
		checkReorgElapsed := time.Since(checkReorgStart)
		k.Logger(ctx).Warn("checkForBitcoinReorg timing", "height", headerHeight, "duration_ms", checkReorgElapsed.Milliseconds())
	}

	return nil
}

// isHistoricalHeader checks if the given Bitcoin block height is in the list of requested historical headers.
func (k *Keeper) isHistoricalHeader(height int64, requestedHeights []int64) bool {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(nil).Warn("isHistoricalHeader timing", "height", height, "duration_ms", elapsed.Milliseconds())
	}()

	for _, h := range requestedHeights {
		if h == height {
			return true
		}
	}
	return false
}

// checkForBitcoinReorg detects reorgs by requesting previous headers when a new one is received.
func (k *Keeper) checkForBitcoinReorg(ctx sdk.Context, oracleData OracleData, requestedHeaders zenbtctypes.RequestedBitcoinHeaders) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("checkForBitcoinReorg timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	var numHistoricalHeadersToRequest int64 = 20
	if strings.HasPrefix(ctx.ChainID(), "diamond") {
		numHistoricalHeadersToRequest = 6
	}

	prevHeights := make([]int64, 0, numHistoricalHeadersToRequest)
	for i := int64(1); i <= numHistoricalHeadersToRequest; i++ {
		prevHeight := oracleData.RequestedBtcBlockHeight - i
		if prevHeight <= 0 {
			break
		}
		prevHeights = append(prevHeights, prevHeight)
	}

	if len(prevHeights) == 0 {
		k.Logger(ctx).Error("no previous heights to request (this should not happen with a valid VE)", "height", oracleData.RequestedBtcBlockHeight)
		return nil
	}

	setRequestedHeadersStart := time.Now()
	requestedHeaders.Heights = append(requestedHeaders.Heights, prevHeights...)
	if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error setting requested historical Bitcoin headers", "error", err)
		return err
	}
	setRequestedHeadersElapsed := time.Since(setRequestedHeadersStart)
	k.Logger(ctx).Warn("RequestedHistoricalBitcoinHeaders.Set timing", "duration_ms", setRequestedHeadersElapsed.Milliseconds())

	k.Logger(ctx).Info("requested headers after reorg check", "headers", requestedHeaders.Heights)

	return nil
}

//
// =============================================================================
// ZENBTC PROCESSING: STAKING, MINTING, BURN EVENTS & REDEMPTIONS
// =============================================================================
//

// checkForUpdateAndDispatchTx processes nonce updates and transaction dispatch
func checkForUpdateAndDispatchTx[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedNonce uint64,
	pendingTxs []T,
	nonceUpdatedCallback func(tx T) error,
	txDispatchCallback func(tx T) error,
) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("checkForUpdateAndDispatchTx timing", "keyID", keyID, "duration_ms", elapsed.Milliseconds())
	}()

	if len(pendingTxs) == 0 {
		return
	}

	getNonceDataStart := time.Now()
	nonceData, err := k.getNonceDataWithInit(ctx, keyID)
	getNonceDataElapsed := time.Since(getNonceDataStart)
	k.Logger(ctx).Warn("getNonceDataWithInit timing", "keyID", keyID, "duration_ms", getNonceDataElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error getting nonce data", "keyID", keyID, "error", err)
		return
	}
	k.Logger(ctx).Info("Nonce info",
		"nonce", nonceData.Nonce,
		"prev", nonceData.PrevNonce,
		"counter", nonceData.Counter,
		"skip", nonceData.Skip,
		"requested", requestedNonce,
	)

	if nonceData.Nonce != 0 && requestedNonce == 0 {
		return
	}

	handleNonceUpdateStart := time.Now()
	nonceUpdated, err := handleNonceUpdate(k, ctx, keyID, requestedNonce, nonceData, pendingTxs[0], nonceUpdatedCallback)
	handleNonceUpdateElapsed := time.Since(handleNonceUpdateStart)
	k.Logger(ctx).Warn("handleNonceUpdate timing", "keyID", keyID, "duration_ms", handleNonceUpdateElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error handling nonce update", "keyID", keyID, "error", err)
		return
	}

	if len(pendingTxs) == 1 && nonceUpdated {
		clearNonceReqStart := time.Now()
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		clearNonceReqElapsed := time.Since(clearNonceReqStart)
		k.Logger(ctx).Warn("clearEthereumNonceRequest timing", "keyID", keyID, "duration_ms", clearNonceReqElapsed.Milliseconds())
		return
	}

	if nonceData.Skip {
		return
	}

	// If tx[0] confirmed on-chain via nonce increment, dispatch tx[1]. If not then retry dispatching tx[0].
	txIndex := 0
	if nonceUpdated {
		txIndex = 1
	}

	txDispatchStart := time.Now()
	if err := txDispatchCallback(pendingTxs[txIndex]); err != nil {
		k.Logger(ctx).Error("tx dispatch callback error", "keyID", keyID, "error", err)
	}
	txDispatchElapsed := time.Since(txDispatchStart)
	k.Logger(ctx).Warn("txDispatchCallback timing", "keyID", keyID, "duration_ms", txDispatchElapsed.Milliseconds())
}

// processZenBTCTransaction is a generic helper that encapsulates the common logic for nonce update and tx dispatch.
func processZenBTCTransaction[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedNonce uint64,
	pendingGetter func(ctx sdk.Context) ([]T, error),
	nonceUpdatedCallback func(tx T) error,
	txDispatchCallback func(tx T) error,
) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("processZenBTCTransaction timing", "keyID", keyID, "duration_ms", elapsed.Milliseconds())
	}()

	isNonceReqStart := time.Now()
	isRequested, err := k.isNonceRequested(ctx, keyID)
	isNonceReqElapsed := time.Since(isNonceReqStart)
	k.Logger(ctx).Warn("isNonceRequested timing", "keyID", keyID, "duration_ms", isNonceReqElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error checking nonce request state", "keyID", keyID, "error", err)
		return
	}
	if !isRequested {
		return
	}

	getPendingStart := time.Now()
	pendingTxs, err := pendingGetter(ctx)
	getPendingElapsed := time.Since(getPendingStart)
	k.Logger(ctx).Warn("pendingGetter timing", "keyID", keyID, "txCount", len(pendingTxs), "duration_ms", getPendingElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error getting pending transactions", "error", err)
		return
	}

	if len(pendingTxs) == 0 {
		clearNonceReqStart := time.Now()
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		clearNonceReqElapsed := time.Since(clearNonceReqStart)
		k.Logger(ctx).Warn("clearEthereumNonceRequest timing", "keyID", keyID, "duration_ms", clearNonceReqElapsed.Milliseconds())
		return
	}

	dispatchStart := time.Now()
	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedNonce, pendingTxs, nonceUpdatedCallback, txDispatchCallback)
	dispatchElapsed := time.Since(dispatchStart)
	k.Logger(ctx).Warn("checkForUpdateAndDispatchTx timing", "keyID", keyID, "duration_ms", dispatchElapsed.Milliseconds())
}

// getPendingTransactions is a generic helper that walks a collections.Map with key type uint64
// and returns a slice of items of type T that satisfy the provided predicate, up to a given limit.
// If limit is 0, all matching items will be returned.
func getPendingTransactions[T any](ctx sdk.Context, store collections.Map[uint64, T], predicate func(T) bool, firstPendingID uint64, limit int) ([]T, error) {
	// startTime := time.Now()
	var results []T
	queryRange := &collections.Range[uint64]{}
	err := store.Walk(ctx, queryRange.StartInclusive(firstPendingID), func(key uint64, value T) (bool, error) {
		if predicate(value) {
			results = append(results, value)
			if limit > 0 && len(results) >= limit {
				return true, nil
			}
		}
		return false, nil
	})

	// elapsed := time.Since(startTime)
	// k := getKeeper(ctx) // Assuming there's a way to get the keeper from context
	// if k != nil {
	// 	k.Logger(ctx).Warn("getPendingTransactions timing", "count", len(results), "duration_ms", elapsed.Milliseconds())
	// }

	return results, err
}

// getNonceDataWithInit gets the nonce data for a key, initializing it if it doesn't exist
func (k *Keeper) getNonceDataWithInit(ctx sdk.Context, keyID uint64) (zenbtctypes.NonceData, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("getNonceDataWithInit timing", "keyID", keyID, "duration_ms", elapsed.Milliseconds())
	}()

	nonceData, err := k.LastUsedEthereumNonce.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return zenbtctypes.NonceData{}, fmt.Errorf("error getting last used ethereum nonce: %w", err)
		}
		nonceData = zenbtctypes.NonceData{Nonce: 0, PrevNonce: 0, Counter: 0, Skip: true}

		setNonceStart := time.Now()
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			return zenbtctypes.NonceData{}, fmt.Errorf("error setting last used ethereum nonce: %w", err)
		}
		setNonceElapsed := time.Since(setNonceStart)
		k.Logger(ctx).Warn("LastUsedEthereumNonce.Set timing", "keyID", keyID, "duration_ms", setNonceElapsed.Milliseconds())
	}
	return nonceData, nil
}

// isNonceRequested checks if a nonce has been requested for the given key
func (k *Keeper) isNonceRequested(ctx sdk.Context, keyID uint64) (bool, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("isNonceRequested timing", "keyID", keyID, "duration_ms", elapsed.Milliseconds())
	}()

	requested, err := k.EthereumNonceRequested.Get(ctx, keyID)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("error getting ethereum nonce request state: %w", err)
	}
	return requested, nil
}

// handleNonceUpdate handles the nonce update logic and returns whether an update occurred
func handleNonceUpdate[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedNonce uint64,
	nonceData zenbtctypes.NonceData,
	tx T,
	nonceUpdatedCallback func(tx T) error,
) (bool, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("handleNonceUpdate timing", "keyID", keyID, "duration_ms", elapsed.Milliseconds())
	}()

	if requestedNonce != nonceData.PrevNonce {
		callbackStart := time.Now()
		if err := nonceUpdatedCallback(tx); err != nil {
			return false, fmt.Errorf("nonce update callback error: %w", err)
		}
		callbackElapsed := time.Since(callbackStart)
		k.Logger(ctx).Warn("nonceUpdatedCallback timing", "keyID", keyID, "duration_ms", callbackElapsed.Milliseconds())

		k.Logger(ctx).Warn("nonce updated for key",
			"keyID", keyID,
			"requestedNonce", requestedNonce,
			"prevNonce", nonceData.PrevNonce,
			"currentNonce", nonceData.Nonce,
		)
		nonceData.PrevNonce = nonceData.Nonce

		setNonceStart := time.Now()
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			return false, fmt.Errorf("error setting last used Ethereum nonce: %w", err)
		}
		setNonceElapsed := time.Since(setNonceStart)
		k.Logger(ctx).Warn("LastUsedEthereumNonce.Set timing", "keyID", keyID, "duration_ms", setNonceElapsed.Milliseconds())

		return true, nil
	}
	return false, nil
}

// updateNonces handles updating nonce state for keys used for minting and unstaking.
func (k *Keeper) updateNonces(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("updateNonces timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	for _, key := range k.getZenBTCKeyIDs(ctx) {
		isNonceReqStart := time.Now()
		isRequested, err := k.isNonceRequested(ctx, key)
		isNonceReqElapsed := time.Since(isNonceReqStart)
		k.Logger(ctx).Warn("isNonceRequested timing", "keyID", key, "duration_ms", isNonceReqElapsed.Milliseconds())

		if err != nil {
			k.Logger(ctx).Error("error checking nonce request state", "keyID", key, "error", err)
			continue
		}
		if !isRequested {
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
		getNonceDataStart := time.Now()
		nonceData, err := k.getNonceDataWithInit(ctx, key)
		getNonceDataElapsed := time.Since(getNonceDataStart)
		k.Logger(ctx).Warn("getNonceDataWithInit timing", "keyID", key, "duration_ms", getNonceDataElapsed.Milliseconds())

		if err != nil {
			k.Logger(ctx).Error("error getting nonce data", "keyID", key, "error", err)
			continue
		}
		if nonceData.Nonce != 0 && currentNonce == 0 {
			continue
		}

		updateNonceStateStart := time.Now()
		if err := k.updateNonceState(ctx, key, currentNonce); err != nil {
			k.Logger(ctx).Error("error updating nonce state", "keyID", key, "error", err)
		}
		updateNonceStateElapsed := time.Since(updateNonceStateStart)
		k.Logger(ctx).Warn("updateNonceState timing", "keyID", key, "duration_ms", updateNonceStateElapsed.Milliseconds())
	}
}

// processZenBTCStaking processes pending staking transactions.
func (k *Keeper) processZenBTCStaking(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("processZenBTCStaking timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetStakerKeyID(ctx),
		oracleData.RequestedStakerNonce,
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			getPendingTxStart := time.Now()
			txs, err := k.getPendingMintTransactionsByStatus(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED)
			getPendingTxElapsed := time.Since(getPendingTxStart)
			k.Logger(ctx).Warn("getPendingMintTransactionsByStatus timing", "status", "DEPOSITED", "duration_ms", getPendingTxElapsed.Milliseconds())
			return txs, err
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			setPendingTxStart := time.Now()
			tx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED
			if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx); err != nil {
				return err
			}
			setPendingTxElapsed := time.Since(setPendingTxStart)
			k.Logger(ctx).Warn("SetPendingMintTransaction timing", "tx_id", tx.Id, "duration_ms", setPendingTxElapsed.Milliseconds())

			if types.IsSolanaCAIP2(tx.Caip2ChainId) {
				setSolanaBlockhashStart := time.Now()
				err := k.SolanaBlockhashRequested.Set(ctx, true)
				setSolanaBlockhashElapsed := time.Since(setSolanaBlockhashStart)
				k.Logger(ctx).Warn("SolanaBlockhashRequested.Set timing", "duration_ms", setSolanaBlockhashElapsed.Milliseconds())
				return err
			} else if types.IsEthereumCAIP2(tx.Caip2ChainId) {
				setEthNonceStart := time.Now()
				err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), true)
				setEthNonceElapsed := time.Since(setEthNonceStart)
				k.Logger(ctx).Warn("EthereumNonceRequested.Set timing", "duration_ms", setEthNonceElapsed.Milliseconds())
				return err
			}
			return fmt.Errorf("unsupported chain type for chain ID: %s", tx.Caip2ChainId)
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			setFirstPendingStart := time.Now()
			if err := k.zenBTCKeeper.SetFirstPendingStakeTransaction(ctx, tx.Id); err != nil {
				return err
			}
			setFirstPendingElapsed := time.Since(setFirstPendingStart)
			k.Logger(ctx).Warn("SetFirstPendingStakeTransaction timing", "tx_id", tx.Id, "duration_ms", setFirstPendingElapsed.Milliseconds())

			// Check for consensus
			validateConsensusStart := time.Now()
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedStakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC stake", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}
			validateConsensusElapsed := time.Since(validateConsensusStart)
			k.Logger(ctx).Warn("validateConsensusForTxFields timing", "tx_id", tx.Id, "duration_ms", validateConsensusElapsed.Milliseconds())

			k.Logger(ctx).Warn("processing zenBTC stake",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.RequestedStakerNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			constructStakeTxStart := time.Now()
			unsignedTxHash, unsignedTx, err := k.constructStakeTx(
				ctx,
				getChainIDForEigen(ctx),
				tx.Amount,
				oracleData.RequestedStakerNonce,
				oracleData.EthGasLimit,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			constructStakeTxElapsed := time.Since(constructStakeTxStart)
			k.Logger(ctx).Warn("constructStakeTx timing", "tx_id", tx.Id, "duration_ms", constructStakeTxElapsed.Milliseconds())

			if err != nil {
				return err
			}

			submitTxStart := time.Now()
			err = k.submitEthereumTransaction(
				ctx,
				tx.Creator,
				k.zenBTCKeeper.GetStakerKeyID(ctx),
				treasurytypes.WalletType(tx.ChainType),
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
			submitTxElapsed := time.Since(submitTxStart)
			k.Logger(ctx).Warn("submitEthereumTransaction timing", "tx_id", tx.Id, "duration_ms", submitTxElapsed.Milliseconds())

			return err
		},
	)
}

// processZenBTCMintsEthereum processes pending mint transactions.
func (k *Keeper) processZenBTCMintsEthereum(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("processZenBTCMintsEthereum timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetEthMinterKeyID(ctx),
		oracleData.RequestedEthMinterNonce,
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			getPendingTxStart := time.Now()
			txs, err := k.getPendingMintTransactionsByStatus(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED)
			getPendingTxElapsed := time.Since(getPendingTxStart)
			k.Logger(ctx).Warn("getPendingMintTransactionsByStatus timing", "status", "STAKED", "duration_ms", getPendingTxElapsed.Milliseconds())
			return txs, err
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			k.Logger(ctx).Warn("processing zenBTC mint",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.RequestedEthMinterNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			getSupplyStart := time.Now()
			supply, err := k.zenBTCKeeper.GetSupply(ctx)
			getSupplyElapsed := time.Since(getSupplyStart)
			k.Logger(ctx).Warn("GetSupply timing", "duration_ms", getSupplyElapsed.Milliseconds())

			if err != nil {
				return err
			}
			supply.PendingZenBTC -= tx.Amount
			supply.MintedZenBTC += tx.Amount

			setSupplyStart := time.Now()
			if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
				return err
			}
			setSupplyElapsed := time.Since(setSupplyStart)
			k.Logger(ctx).Warn("SetSupply timing", "duration_ms", setSupplyElapsed.Milliseconds())

			k.Logger(ctx).Warn("pending mint supply updated",
				"pending_mint_old", supply.PendingZenBTC+tx.Amount,
				"pending_mint_new", supply.PendingZenBTC,
			)
			k.Logger(ctx).Warn("minted supply updated",
				"minted_old", supply.MintedZenBTC-tx.Amount,
				"minted_new", supply.MintedZenBTC,
			)
			tx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED

			setPendingTxStart := time.Now()
			err = k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
			setPendingTxElapsed := time.Since(setPendingTxStart)
			k.Logger(ctx).Warn("SetPendingMintTransaction timing", "tx_id", tx.Id, "duration_ms", setPendingTxElapsed.Milliseconds())

			return err
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			setFirstPendingStart := time.Now()
			if err := k.zenBTCKeeper.SetFirstPendingMintTransaction(ctx, tx.Id); err != nil {
				return err
			}
			setFirstPendingElapsed := time.Since(setFirstPendingStart)
			k.Logger(ctx).Warn("SetFirstPendingMintTransaction timing", "tx_id", tx.Id, "duration_ms", setFirstPendingElapsed.Milliseconds())

			// Check for consensus
			validateConsensusStart := time.Now()
			requiredFields := []VoteExtensionField{VEFieldRequestedEthMinterNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"zenBTC mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}
			validateConsensusElapsed := time.Since(validateConsensusStart)
			k.Logger(ctx).Warn("validateConsensusForTxFields timing", "tx_id", tx.Id, "duration_ms", validateConsensusElapsed.Milliseconds())
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

			chainID, err := types.ValidateChainID(ctx, tx.Caip2ChainId)
			if err != nil {
				return fmt.Errorf("unsupported chain ID: %w", err)
			}

			constructMintTxStart := time.Now()
			unsignedMintTxHash, unsignedMintTx, err := k.constructMintTx(
				ctx,
				tx.RecipientAddress,
				chainID,
				tx.Amount,
				feeZenBTC,
				oracleData.RequestedEthMinterNonce,
				oracleData.EthGasLimit,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			constructMintTxElapsed := time.Since(constructMintTxStart)
			k.Logger(ctx).Warn("constructMintTx timing", "tx_id", tx.Id, "duration_ms", constructMintTxElapsed.Milliseconds())

			if err != nil {
				return err
			}

			submitTxStart := time.Now()
			err = k.submitEthereumTransaction(
				ctx,
				tx.Creator,
				k.zenBTCKeeper.GetEthMinterKeyID(ctx),
				treasurytypes.WalletType(tx.ChainType),
				chainID,
				unsignedMintTx,
				unsignedMintTxHash,
			)
			submitTxElapsed := time.Since(submitTxStart)
			k.Logger(ctx).Warn("submitEthereumTransaction timing", "tx_id", tx.Id, "duration_ms", submitTxElapsed.Milliseconds())

			return err
		},
	)
}

// storeNewZenBTCBurnEventsEthereum stores new burn events coming from Ethereum.
func (k *Keeper) storeNewZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("storeNewZenBTCBurnEventsEthereum timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	foundNewBurn := false
	// Loop over each burn event from oracle to check for new ones.
	for _, burn := range oracleData.EthBurnEvents {
		// Check if this burn event already exists
		exists := false
		walkBurnStart := time.Now()
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
		walkBurnElapsed := time.Since(walkBurnStart)
		k.Logger(ctx).Warn("WalkBurnEvents timing", "burn_txid", burn.TxID, "duration_ms", walkBurnElapsed.Milliseconds())

		if !exists {
			newBurn := zenbtctypes.BurnEvent{
				TxID:            burn.TxID,
				LogIndex:        burn.LogIndex,
				ChainID:         burn.ChainID,
				DestinationAddr: burn.DestinationAddr,
				Amount:          burn.Amount,
				Status:          zenbtctypes.BurnStatus_BURN_STATUS_BURNED,
			}
			createBurnStart := time.Now()
			id, err := k.zenBTCKeeper.CreateBurnEvent(ctx, &newBurn)
			createBurnElapsed := time.Since(createBurnStart)
			k.Logger(ctx).Warn("CreateBurnEvent timing", "burn_txid", burn.TxID, "duration_ms", createBurnElapsed.Milliseconds())

			if err != nil {
				k.Logger(ctx).Error("error creating burn event", "error", err)
				continue
			}
			k.Logger(ctx).Info("created new burn event", "id", id)
			foundNewBurn = true
		}
	}

	if foundNewBurn {
		setNonceReqStart := time.Now()
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "error", err)
		}
		setNonceReqElapsed := time.Since(setNonceReqStart)
		k.Logger(ctx).Warn("EthereumNonceRequested.Set timing", "key", "unstaker", "duration_ms", setNonceReqElapsed.Milliseconds())
	}
}

// processZenBTCBurnEventsEthereum processes pending burn events by constructing unstake transactions.
func (k *Keeper) processZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("processZenBTCBurnEventsEthereum timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetUnstakerKeyID(ctx),
		oracleData.RequestedUnstakerNonce,
		func(ctx sdk.Context) ([]zenbtctypes.BurnEvent, error) {
			getPendingBurnStart := time.Now()
			events, err := k.getPendingBurnEvents(ctx)
			getPendingBurnElapsed := time.Since(getPendingBurnStart)
			k.Logger(ctx).Warn("getPendingBurnEvents timing", "duration_ms", getPendingBurnElapsed.Milliseconds())
			return events, err
		},
		func(be zenbtctypes.BurnEvent) error {
			be.Status = zenbtctypes.BurnStatus_BURN_STATUS_UNSTAKING
			setBurnStart := time.Now()
			err := k.zenBTCKeeper.SetBurnEvent(ctx, be.Id, be)
			setBurnElapsed := time.Since(setBurnStart)
			k.Logger(ctx).Warn("SetBurnEvent timing", "burn_id", be.Id, "duration_ms", setBurnElapsed.Milliseconds())
			return err
		},
		func(be zenbtctypes.BurnEvent) error {
			setFirstPendingStart := time.Now()
			if err := k.zenBTCKeeper.SetFirstPendingBurnEvent(ctx, be.Id); err != nil {
				return err
			}
			setFirstPendingElapsed := time.Since(setFirstPendingStart)
			k.Logger(ctx).Warn("SetFirstPendingBurnEvent timing", "burn_id", be.Id, "duration_ms", setFirstPendingElapsed.Milliseconds())

			// Check for consensus
			validateConsensusStart := time.Now()
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedUnstakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC burn unstake", fmt.Sprintf("burn_id: %d, destination: %s, amount: %d", be.Id, be.DestinationAddr, be.Amount)); err != nil {
				return err
			}
			validateConsensusElapsed := time.Since(validateConsensusStart)
			k.Logger(ctx).Warn("validateConsensusForTxFields timing", "burn_id", be.Id, "duration_ms", validateConsensusElapsed.Milliseconds())

			k.Logger(ctx).Warn("processing zenBTC burn unstake",
				"burn_event", be,
				"nonce", oracleData.RequestedUnstakerNonce,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			constructUnstakeTxStart := time.Now()
			unsignedTxHash, unsignedTx, err := k.constructUnstakeTx(
				ctx,
				getChainIDForEigen(ctx),
				be.DestinationAddr,
				be.Amount,
				oracleData.RequestedUnstakerNonce,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			constructUnstakeTxElapsed := time.Since(constructUnstakeTxStart)
			k.Logger(ctx).Warn("constructUnstakeTx timing", "burn_id", be.Id, "duration_ms", constructUnstakeTxElapsed.Milliseconds())

			if err != nil {
				return err
			}

			getAddressStart := time.Now()
			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			getAddressElapsed := time.Since(getAddressStart)
			k.Logger(ctx).Warn("getAddressByKeyID timing", "key", "unstaker", "duration_ms", getAddressElapsed.Milliseconds())

			if err != nil {
				return err
			}

			submitTxStart := time.Now()
			err = k.submitEthereumTransaction(
				ctx,
				creator,
				k.zenBTCKeeper.GetUnstakerKeyID(ctx),
				treasurytypes.WalletType_WALLET_TYPE_EVM,
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
			submitTxElapsed := time.Since(submitTxStart)
			k.Logger(ctx).Warn("submitEthereumTransaction timing", "burn_id", be.Id, "duration_ms", submitTxElapsed.Milliseconds())

			return err
		},
	)
}

// Helper function to submit Ethereum transactions
func (k *Keeper) submitEthereumTransaction(ctx sdk.Context, creator string, keyID uint64, walletType treasurytypes.WalletType, chainID uint64, unsignedTx []byte, unsignedTxHash []byte) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("submitEthereumTransaction timing", "keyID", keyID, "duration_ms", elapsed.Milliseconds())
	}()

	createAnyStart := time.Now()
	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: chainID})
	createAnyElapsed := time.Since(createAnyStart)
	k.Logger(ctx).Warn("NewAnyWithValue timing", "duration_ms", createAnyElapsed.Milliseconds())

	if err != nil {
		return err
	}

	handleSignStart := time.Now()
	_, err = k.treasuryKeeper.HandleSignTransactionRequest(
		ctx,
		&treasurytypes.MsgNewSignTransactionRequest{
			Creator:             creator,
			KeyId:               keyID,
			WalletType:          walletType,
			UnsignedTransaction: unsignedTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedTxHash)),
	)
	handleSignElapsed := time.Since(handleSignStart)
	k.Logger(ctx).Warn("HandleSignTransactionRequest timing", "keyID", keyID, "duration_ms", handleSignElapsed.Milliseconds())

	return err
}

// storeNewZenBTCRedemptions processes new redemption events.
func (k *Keeper) storeNewZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("storeNewZenBTCRedemptions timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	// Find the first INITIATED redemption.
	var firstInitiatedRedemption zenbtctypes.Redemption
	var found bool

	findInitiatedStart := time.Now()
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
	findInitiatedElapsed := time.Since(findInitiatedStart)
	k.Logger(ctx).Warn("Find initiated redemption timing", "duration_ms", findInitiatedElapsed.Milliseconds())

	// If an INITIATED redemption is found, check if it exists in oracleData.
	if found {
		checkOracleDataStart := time.Now()
		redemptionExists := false
		for _, redemption := range oracleData.Redemptions {
			if redemption.Id == firstInitiatedRedemption.Data.Id {
				redemptionExists = true
				break
			}
		}
		checkOracleDataElapsed := time.Since(checkOracleDataStart)
		k.Logger(ctx).Warn("Check redemption exists in oracle data timing", "duration_ms", checkOracleDataElapsed.Milliseconds())

		// If not present, mark it as unstaked.
		if !redemptionExists {
			updateRedemptionStart := time.Now()
			firstInitiatedRedemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
			if err := k.zenBTCKeeper.SetRedemption(ctx, firstInitiatedRedemption.Data.Id, firstInitiatedRedemption); err != nil {
				k.Logger(ctx).Error("error updating redemption status to unstaked", "error", err)
				return
			}
			updateRedemptionElapsed := time.Since(updateRedemptionStart)
			k.Logger(ctx).Warn("Update redemption to unstaked timing", "id", firstInitiatedRedemption.Data.Id, "duration_ms", updateRedemptionElapsed.Milliseconds())
		}
	}

	if len(oracleData.Redemptions) == 0 {
		return
	}

	getExchangeRateStart := time.Now()
	exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
	getExchangeRateElapsed := time.Since(getExchangeRateStart)
	k.Logger(ctx).Warn("GetExchangeRate timing", "duration_ms", getExchangeRateElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error getting zenBTC exchange rate", "error", err)
		return
	}

	foundNewRedemption := false

	for _, redemption := range oracleData.Redemptions {
		checkRedemptionStart := time.Now()
		redemptionExists, err := k.zenBTCKeeper.HasRedemption(ctx, redemption.Id)
		checkRedemptionElapsed := time.Since(checkRedemptionStart)
		k.Logger(ctx).Warn("HasRedemption timing", "id", redemption.Id, "duration_ms", checkRedemptionElapsed.Milliseconds())

		if err != nil {
			k.Logger(ctx).Error("error checking redemption existence", "error", err)
			continue
		}
		if redemptionExists {
			k.Logger(ctx).Debug("redemption already stored", "id", redemption.Id)
			continue
		}

		foundNewRedemption = true

		btcAmount := sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Amount)).Mul(exchangeRate).TruncateInt64()
		// Convert zenBTC amount to BTC amount.
		setRedemptionStart := time.Now()
		if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Id, zenbtctypes.Redemption{
			Data: zenbtctypes.RedemptionData{
				Id:                 redemption.Id,
				DestinationAddress: redemption.DestinationAddress,
				Amount:             uint64(btcAmount),
			},
			Status: zenbtctypes.RedemptionStatus_INITIATED,
		}); err != nil {
			k.Logger(ctx).Error("error adding redemption to store", "error", err)
			continue
		}
		setRedemptionElapsed := time.Since(setRedemptionStart)
		k.Logger(ctx).Warn("SetRedemption timing", "id", redemption.Id, "duration_ms", setRedemptionElapsed.Milliseconds())
	}

	if foundNewRedemption {
		setNonceReqStart := time.Now()
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "error", err)
		}
		setNonceReqElapsed := time.Since(setNonceReqStart)
		k.Logger(ctx).Warn("EthereumNonceRequested.Set timing", "key", "completer", "duration_ms", setNonceReqElapsed.Milliseconds())
	}
}

// processZenBTCRedemptions processes pending redemption completions.
func (k *Keeper) processZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("processZenBTCRedemptions timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetCompleterKeyID(ctx),
		oracleData.RequestedCompleterNonce,
		func(ctx sdk.Context) ([]zenbtctypes.Redemption, error) {
			getFirstPendingStart := time.Now()
			firstPendingID, err := k.zenBTCKeeper.GetFirstPendingRedemption(ctx)
			getFirstPendingElapsed := time.Since(getFirstPendingStart)
			k.Logger(ctx).Warn("GetFirstPendingRedemption timing", "duration_ms", getFirstPendingElapsed.Milliseconds())

			if err != nil {
				firstPendingID = 0
			}

			getRedemptionsStart := time.Now()
			redemptions, err := k.getRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_INITIATED, 2, firstPendingID)
			getRedemptionsElapsed := time.Since(getRedemptionsStart)
			k.Logger(ctx).Warn("getRedemptionsByStatus timing", "count", len(redemptions), "duration_ms", getRedemptionsElapsed.Milliseconds())

			return redemptions, err
		},
		func(r zenbtctypes.Redemption) error {
			r.Status = zenbtctypes.RedemptionStatus_UNSTAKED

			setRedemptionStart := time.Now()
			if err := k.zenBTCKeeper.SetRedemption(ctx, r.Data.Id, r); err != nil {
				return err
			}
			setRedemptionElapsed := time.Since(setRedemptionStart)
			k.Logger(ctx).Warn("SetRedemption timing", "id", r.Data.Id, "duration_ms", setRedemptionElapsed.Milliseconds())

			setNonceReqStart := time.Now()
			err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), true)
			setNonceReqElapsed := time.Since(setNonceReqStart)
			k.Logger(ctx).Warn("EthereumNonceRequested.Set timing", "key", "staker", "duration_ms", setNonceReqElapsed.Milliseconds())

			return err
		},
		func(r zenbtctypes.Redemption) error {
			setFirstPendingStart := time.Now()
			if err := k.zenBTCKeeper.SetFirstPendingRedemption(ctx, r.Data.Id); err != nil {
				return err
			}
			setFirstPendingElapsed := time.Since(setFirstPendingStart)
			k.Logger(ctx).Warn("SetFirstPendingRedemption timing", "id", r.Data.Id, "duration_ms", setFirstPendingElapsed.Milliseconds())

			// Check for consensus
			validateConsensusStart := time.Now()
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedCompleterNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC redemption", fmt.Sprintf("redemption_id: %d, amount: %d", r.Data.Id, r.Data.Amount)); err != nil {
				return err
			}
			validateConsensusElapsed := time.Since(validateConsensusStart)
			k.Logger(ctx).Warn("validateConsensusForTxFields timing", "id", r.Data.Id, "duration_ms", validateConsensusElapsed.Milliseconds())

			k.Logger(ctx).Warn("processing zenBTC complete",
				"id", r.Data.Id,
				"nonce", oracleData.RequestedCompleterNonce,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			constructCompleteTxStart := time.Now()
			unsignedTxHash, unsignedTx, err := k.constructCompleteTx(
				ctx,
				getChainIDForEigen(ctx),
				r.Data.Id,
				oracleData.RequestedCompleterNonce,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			constructCompleteTxElapsed := time.Since(constructCompleteTxStart)
			k.Logger(ctx).Warn("constructCompleteTx timing", "id", r.Data.Id, "duration_ms", constructCompleteTxElapsed.Milliseconds())

			if err != nil {
				return err
			}

			getAddressStart := time.Now()
			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			getAddressElapsed := time.Since(getAddressStart)
			k.Logger(ctx).Warn("getAddressByKeyID timing", "key", "completer", "duration_ms", getAddressElapsed.Milliseconds())

			if err != nil {
				return err
			}

			submitTxStart := time.Now()
			err = k.submitEthereumTransaction(
				ctx,
				creator,
				k.zenBTCKeeper.GetCompleterKeyID(ctx),
				treasurytypes.WalletType_WALLET_TYPE_EVM,
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
			submitTxElapsed := time.Since(submitTxStart)
			k.Logger(ctx).Warn("submitEthereumTransaction timing", "id", r.Data.Id, "duration_ms", submitTxElapsed.Milliseconds())

			return err
		},
	)
}

func (k *Keeper) checkForRedemptionFulfilment(ctx sdk.Context) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		k.Logger(ctx).Warn("checkForRedemptionFulfilment timing", "height", ctx.BlockHeight(), "duration_ms", elapsed.Milliseconds())
	}()

	getFirstRedemptionStart := time.Now()
	startingIndex, err := k.zenBTCKeeper.GetFirstRedemptionAwaitingSign(ctx)
	getFirstRedemptionElapsed := time.Since(getFirstRedemptionStart)
	k.Logger(ctx).Warn("GetFirstRedemptionAwaitingSign timing", "duration_ms", getFirstRedemptionElapsed.Milliseconds())

	if err != nil {
		startingIndex = 0
	}

	getRedemptionsStart := time.Now()
	redemptions, err := k.getRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_AWAITING_SIGN, 0, startingIndex)
	getRedemptionsElapsed := time.Since(getRedemptionsStart)
	k.Logger(ctx).Warn("getRedemptionsByStatus timing", "count", len(redemptions), "duration_ms", getRedemptionsElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error getting redemptions", "error", err)
		return
	}

	if len(redemptions) == 0 {
		return
	}

	setFirstRedemptionStart := time.Now()
	if err := k.zenBTCKeeper.SetFirstRedemptionAwaitingSign(ctx, redemptions[0].Data.Id); err != nil {
		k.Logger(ctx).Error("error setting first redemption awaiting sign", "error", err)
	}
	setFirstRedemptionElapsed := time.Since(setFirstRedemptionStart)
	k.Logger(ctx).Warn("SetFirstRedemptionAwaitingSign timing", "id", redemptions[0].Data.Id, "duration_ms", setFirstRedemptionElapsed.Milliseconds())

	getSupplyStart := time.Now()
	supply, err := k.zenBTCKeeper.GetSupply(ctx)
	getSupplyElapsed := time.Since(getSupplyStart)
	k.Logger(ctx).Warn("GetSupply timing", "duration_ms", getSupplyElapsed.Milliseconds())

	if err != nil {
		k.Logger(ctx).Error("error getting zenBTC supply", "error", err)
		return
	}

	for _, redemption := range redemptions {
		getSignReqStart := time.Now()
		signReq, err := k.treasuryKeeper.SignRequestStore.Get(ctx, redemption.Data.SignReqId)
		getSignReqElapsed := time.Since(getSignReqStart)
		k.Logger(ctx).Warn("SignRequestStore.Get timing", "id", redemption.Data.SignReqId, "duration_ms", getSignReqElapsed.Milliseconds())

		if err != nil {
			k.Logger(ctx).Error("error getting sign request for redemption", "id", redemption.Data.Id, "error", err)
			continue
		}

		if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING {
			continue
		}
		if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED {
			// Get current exchange rate
			getExchangeRateStart := time.Now()
			exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
			getExchangeRateElapsed := time.Since(getExchangeRateStart)
			k.Logger(ctx).Warn("GetExchangeRate timing", "duration_ms", getExchangeRateElapsed.Milliseconds())

			if err != nil {
				k.Logger(ctx).Error("error getting zenBTC exchange rate", "error", err)
				continue
			}

			// redemption.Data.Amount is in zenBTC (what user wants to redeem)
			// Calculate how much BTC they should receive based on current exchange rate
			convertAmountStart := time.Now()
			btcToRelease := uint64(sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Data.Amount)).Quo(exchangeRate).TruncateInt64())
			convertAmountElapsed := time.Since(convertAmountStart)
			k.Logger(ctx).Warn("Convert amount timing", "id", redemption.Data.Id, "duration_ms", convertAmountElapsed.Milliseconds())

			// Invariant checks
			if supply.MintedZenBTC < redemption.Data.Amount {
				k.Logger(ctx).Error("insufficient minted zenBTC for redemption", "id", redemption.Data.Id)
				continue
			}
			if supply.CustodiedBTC < btcToRelease {
				k.Logger(ctx).Error("insufficient custodied BTC for redemption", "id", redemption.Data.Id)
				continue
			}

			// Update supplies (zenBTC burned, BTC released)
			supply.MintedZenBTC -= redemption.Data.Amount
			supply.CustodiedBTC -= btcToRelease

			k.Logger(ctx).Warn("minted supply updated", "minted_old", supply.MintedZenBTC+redemption.Data.Amount, "minted_new", supply.MintedZenBTC)
			k.Logger(ctx).Warn("custodied supply updated", "custodied_old", supply.CustodiedBTC+btcToRelease, "custodied_new", supply.CustodiedBTC)

			redemption.Status = zenbtctypes.RedemptionStatus_COMPLETED
		}
		if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
			redemption.Data.SignReqId = 0
			redemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
		}

		setRedemptionStart := time.Now()
		if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Data.Id, redemption); err != nil {
			k.Logger(ctx).Error("error updating redemption status", "error", err)
		}
		setRedemptionElapsed := time.Since(setRedemptionStart)
		k.Logger(ctx).Warn("SetRedemption timing", "id", redemption.Data.Id, "duration_ms", setRedemptionElapsed.Milliseconds())
	}

	setSupplyStart := time.Now()
	if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
		k.Logger(ctx).Error("error updating zenBTC supply", "error", err)
	}
	setSupplyElapsed := time.Since(setSupplyStart)
	k.Logger(ctx).Warn("SetSupply timing", "duration_ms", setSupplyElapsed.Milliseconds())
}
