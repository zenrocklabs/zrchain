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

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
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
	redemptionsHash, err := deriveHash(oracleData.Redemptions)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving redemptions hash: %w", err)
	}

	latestHeader, requestedHeader, err := k.retrieveBitcoinHeaders(ctx)
	if err != nil {
		return VoteExtension{}, err
	}
	latestBitcoinHeaderHash, err := deriveHash(latestHeader.BlockHeader)
	if err != nil {
		return VoteExtension{}, err
	}
	// Only set requested header fields if there's a requested header
	requestedBtcBlockHeight := int64(0)
	var requestedBtcHeaderHash []byte
	if requestedHeader != nil {
		// Get requested Bitcoin header hash
		requestedBitcoinHeaderHash, err := deriveHash(requestedHeader.BlockHeader)
		if err != nil {
			return VoteExtension{}, err
		}
		requestedBtcBlockHeight = requestedHeader.BlockHeight
		requestedBtcHeaderHash = requestedBitcoinHeaderHash[:]
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

	voteExt, fieldVotePowers, err := k.GetSuperMajorityVEData(ctx, req.Height, req.LocalLastCommit)
	if err != nil {
		k.Logger(ctx).Error("error retrieving supermajority vote extension data", "height", req.Height, "error", err)
		return nil, nil
	}

	if len(fieldVotePowers) == 0 { // no field reached consensus
		k.Logger(ctx).Warn("no fields reached consensus in vote extension", "height", req.Height)
		return k.marshalOracleData(req, &OracleData{ConsensusData: req.LocalLastCommit, FieldVotePowers: fieldVotePowers})
	}

	if voteExt.ZRChainBlockHeight != req.Height-1 { // vote extension is from previous block
		k.Logger(ctx).Error("mismatched height for vote extension", "height", req.Height, "voteExt.ZRChainBlockHeight", voteExt.ZRChainBlockHeight)
		return nil, nil
	}

	oracleData, err := k.getValidatedOracleData(ctx, voteExt, fieldVotePowers)
	if err != nil {
		k.Logger(ctx).Warn("error in getValidatedOracleData; injecting empty oracle data", "height", req.Height, "error", err)
		oracleData = &OracleData{}
	}

	oracleData.ConsensusData = req.LocalLastCommit

	return k.marshalOracleData(req, oracleData)
}

// ProcessProposal is executed by all validators to check whether the proposer prepared valid data.
func (k *Keeper) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
	// Return early if this node is not a validator so non-validators don't need to be running a sidecar
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

	// Check for empty oracle data - if it's empty, accept the proposal
	recoveredOracleDataNoCommitInfo := recoveredOracleData
	recoveredOracleDataNoCommitInfo.ConsensusData = abci.ExtendedCommitInfo{}
	recoveredOracleDataNoCommitInfo.FieldVotePowers = nil
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

// PreBlocker processes oracle data and applies the resulting state updates.
func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	if !k.shouldProcessOracleData(ctx, req) {
		return nil
	}

	oracleData, ok := k.unmarshalOracleData(ctx, req.Txs[0])
	if !ok {
		return nil
	}

	canonicalVE, ok := k.validateCanonicalVE(ctx, req.Height, oracleData)
	if !ok {
		k.Logger(ctx).Error("invalid canonical vote extension")
		return nil
	}

	// Update asset prices if there's consensus on the price fields
	k.updateAssetPrices(ctx, oracleData)

	// Process different subsystems based on field consensus

	// Validator updates - only if EigenDelegationsHash has consensus
	if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldEigenDelegationsHash) {
		k.updateValidatorStakes(ctx, oracleData)
		k.updateAVSDelegationStore(ctx, oracleData)
	}

	// Bitcoin header processing - only if BTC header fields have consensus
	btcHeaderFields := []VoteExtensionField{VEFieldLatestBtcHeaderHash, VEFieldRequestedBtcHeaderHash}
	if anyFieldHasConsensus(oracleData.FieldVotePowers, btcHeaderFields) {
		if err := k.storeBitcoinBlockHeaders(ctx, oracleData); err != nil {
			k.Logger(ctx).Error("error storing Bitcoin headers", "error", err)
		}
	}

	if ctx.BlockHeight()%2 == 0 { // TODO: is this needed?

		nonceFields := []VoteExtensionField{
			VEFieldRequestedStakerNonce,
			VEFieldRequestedEthMinterNonce,
			VEFieldRequestedUnstakerNonce,
			VEFieldRequestedCompleterNonce,
		}
		if anyFieldHasConsensus(oracleData.FieldVotePowers, nonceFields) {
			k.updateNonces(ctx, oracleData)
		}

		if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldEthBurnEventsHash) {
			k.storeNewZenBTCBurnEventsEthereum(ctx, oracleData)
		}
		if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldRedemptionsHash) {
			k.storeNewZenBTCRedemptions(ctx, oracleData)
		}

		k.processZenBTCStaking(ctx, oracleData)
		k.processZenBTCMintsEthereum(ctx, oracleData)
		k.processZenBTCBurnEventsEthereum(ctx, oracleData)
		k.processZenBTCRedemptions(ctx, oracleData)
	}

	k.recordNonVotingValidators(ctx, req)
	k.recordMismatchedVoteExtensions(ctx, req.Height, canonicalVE, oracleData.ConsensusData)

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

// validateCanonicalVE validates the canonical vote extension from oracle data.
func (k *Keeper) validateCanonicalVE(ctx sdk.Context, height int64, oracleData OracleData) (VoteExtension, bool) {
	voteExt, fieldVotePowers, err := k.GetSuperMajorityVEData(ctx, height, oracleData.ConsensusData)
	if err != nil {
		k.Logger(ctx).Error("error getting super majority VE data", "height", height, "error", err)
		return VoteExtension{}, false
	}

	if reflect.DeepEqual(voteExt, VoteExtension{}) {
		k.Logger(ctx).Warn("accepting empty vote extension", "height", height)
		return VoteExtension{}, true
	}

	if err := k.validateOracleData(ctx, voteExt, &oracleData, fieldVotePowers); err != nil {
		k.Logger(ctx).Error("error validating oracle data; won't store VE data", "height", height, "error", err)
		return VoteExtension{}, false
	}

	// Log final consensus summary after validation
	k.Logger(ctx).Info("final consensus summary",
		"fields_with_consensus", len(fieldVotePowers),
		"stage", "post_validation")

	return voteExt, true
}

// getValidatedOracleData retrieves and validates oracle data based on a vote extension.
// Only validates fields that have reached consensus as indicated in fieldVotePowers.
func (k *Keeper) getValidatedOracleData(ctx sdk.Context, voteExt VoteExtension, fieldVotePowers map[VoteExtensionField]int64) (*OracleData, error) {
	// We only fetch Ethereum state if we have consensus on EthBlockHeight
	var oracleData *OracleData
	var err error

	if _, ok := fieldVotePowers[VEFieldEthBlockHeight]; ok {
		oracleData, err = k.GetSidecarStateByEthHeight(ctx, voteExt.EthBlockHeight)
		if err != nil {
			return nil, fmt.Errorf("error fetching oracle state: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no consensus on eth block height")
	}

	latestHeader, requestedHeader, err := k.retrieveBitcoinHeaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching bitcoin headers: %w", err)
	}

	// Copy latest Bitcoin header data if we have consensus
	if _, ok := fieldVotePowers[VEFieldLatestBtcBlockHeight]; ok && latestHeader != nil {
		oracleData.LatestBtcBlockHeight = latestHeader.BlockHeight
		oracleData.LatestBtcBlockHeader = *latestHeader.BlockHeader
	}
	// Copy requested Bitcoin header data if we have consensus and the header exists
	if _, ok := fieldVotePowers[VEFieldRequestedBtcBlockHeight]; ok && requestedHeader != nil {
		oracleData.RequestedBtcBlockHeight = requestedHeader.BlockHeight
		oracleData.RequestedBtcBlockHeader = *requestedHeader.BlockHeader
	}

	// Copy over nonce data if we have consensus on those fields
	if _, ok := fieldVotePowers[VEFieldRequestedStakerNonce]; ok {
		oracleData.RequestedStakerNonce = voteExt.RequestedStakerNonce
	}
	if _, ok := fieldVotePowers[VEFieldRequestedEthMinterNonce]; ok {
		oracleData.RequestedEthMinterNonce = voteExt.RequestedEthMinterNonce
	}
	if _, ok := fieldVotePowers[VEFieldRequestedUnstakerNonce]; ok {
		oracleData.RequestedUnstakerNonce = voteExt.RequestedUnstakerNonce
	}
	if _, ok := fieldVotePowers[VEFieldRequestedCompleterNonce]; ok {
		oracleData.RequestedCompleterNonce = voteExt.RequestedCompleterNonce
	}

	// Store the field vote powers for later use in transaction dispatch callbacks
	oracleData.FieldVotePowers = fieldVotePowers

	if err := k.validateOracleData(ctx, voteExt, oracleData, fieldVotePowers); err != nil {
		return nil, err
	}

	return oracleData, nil
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
func (k *Keeper) storeBitcoinBlockHeaders(ctx sdk.Context, oracleData OracleData) error {
	// First store the latest Bitcoin header if available
	if oracleData.LatestBtcBlockHeight > 0 && oracleData.LatestBtcBlockHeader.MerkleRoot != "" {
		latestHeaderExists, err := k.BtcBlockHeaders.Has(ctx, oracleData.LatestBtcBlockHeight)
		if err != nil {
			k.Logger(ctx).Error("error checking if latest Bitcoin header exists", "height", oracleData.LatestBtcBlockHeight, "error", err)
		} else if !latestHeaderExists {
			// Only store if it doesn't already exist
			if err := k.BtcBlockHeaders.Set(ctx, oracleData.LatestBtcBlockHeight, oracleData.LatestBtcBlockHeader); err != nil {
				k.Logger(ctx).Error("error storing latest Bitcoin header", "height", oracleData.LatestBtcBlockHeight, "error", err)
			} else {
				k.Logger(ctx).Info("stored latest Bitcoin header", "height", oracleData.LatestBtcBlockHeight)
			}
		}
	}

	// Process the requested Bitcoin header
	headerHeight := oracleData.RequestedBtcBlockHeight
	if headerHeight == 0 || oracleData.RequestedBtcBlockHeader.MerkleRoot == "" {
		k.Logger(ctx).Debug("no requested bitcoin header")
		return nil
	}

	// Get requested headers
	requestedHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting requested historical Bitcoin headers", "error", err)
			return err
		}
		k.Logger(ctx).Info("requested historical Bitcoin headers store not initialised", "height", headerHeight)
	}

	// Check if the header is historical
	isHistorical := k.isHistoricalHeader(headerHeight, requestedHeaders.Heights)

	// Check if header exists (for logging only)
	headerExists, err := k.BtcBlockHeaders.Has(ctx, headerHeight)
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
	if err := k.BtcBlockHeaders.Set(ctx, headerHeight, oracleData.RequestedBtcBlockHeader); err != nil {
		k.Logger(ctx).Error("error storing Bitcoin header", "height", headerHeight, "error", err)
		return err
	}
	logger.Info("stored Bitcoin header",
		"type", map[bool]string{true: "historical", false: "latest"}[isHistorical])

	// Process according to header type
	if isHistorical {
		// Remove the processed historical header from the requested list
		requestedHeaders.Heights = slices.DeleteFunc(requestedHeaders.Heights, func(height int64) bool {
			return height == headerHeight
		})

		if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
			k.Logger(ctx).Error("error updating requested historical Bitcoin headers", "error", err)
			return err
		}

		logger.Debug("removed processed historical header request",
			"remaining_requests", len(requestedHeaders.Heights))
	} else if !headerExists {
		// Only check for reorgs for non-historical headers that weren't already stored
		if err := k.checkForBitcoinReorg(ctx, oracleData, requestedHeaders); err != nil {
			k.Logger(ctx).Error("error handling potential Bitcoin reorg", "height", headerHeight, "error", err)
		}
	}

	return nil
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

	requestedHeaders.Heights = append(requestedHeaders.Heights, prevHeights...)
	if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error setting requested historical Bitcoin headers", "error", err)
		return err
	}

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
	if len(pendingTxs) == 0 {
		return
	}

	nonceData, err := k.getNonceDataWithInit(ctx, keyID)
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

	nonceUpdated, err := handleNonceUpdate(k, ctx, keyID, requestedNonce, nonceData, pendingTxs[0], nonceUpdatedCallback)
	if err != nil {
		k.Logger(ctx).Error("error handling nonce update", "keyID", keyID, "error", err)
		return
	}

	if len(pendingTxs) == 1 && nonceUpdated {
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
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

	if err := txDispatchCallback(pendingTxs[txIndex]); err != nil {
		k.Logger(ctx).Error("tx dispatch callback error", "keyID", keyID, "error", err)
	}
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
	isRequested, err := k.isNonceRequested(ctx, keyID)
	if err != nil {
		k.Logger(ctx).Error("error checking nonce request state", "keyID", keyID, "error", err)
		return
	}
	if !isRequested {
		return
	}

	pendingTxs, err := pendingGetter(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting pending transactions", "error", err)
		return
	}

	if len(pendingTxs) == 0 {
		if err := k.clearEthereumNonceRequest(ctx, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		return
	}
	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedNonce, pendingTxs, nonceUpdatedCallback, txDispatchCallback)
}

// getPendingTransactions is a generic helper that walks a collections.Map with key type uint64
// and returns a slice of items of type T that satisfy the provided predicate, up to a given limit.
func getPendingTransactions[T any](ctx sdk.Context, store collections.Map[uint64, T], predicate func(T) bool, firstPendingID uint64, limit int) ([]T, error) {
	var results []T
	queryRange := &collections.Range[uint64]{}
	err := store.Walk(ctx, queryRange.StartInclusive(firstPendingID), func(key uint64, value T) (bool, error) {
		if predicate(value) {
			results = append(results, value)
			if len(results) >= limit {
				return true, nil
			}
		}
		return false, nil
	})
	return results, err
}

// getNonceDataWithInit gets the nonce data for a key, initializing it if it doesn't exist
func (k *Keeper) getNonceDataWithInit(ctx sdk.Context, keyID uint64) (zenbtctypes.NonceData, error) {
	nonceData, err := k.LastUsedEthereumNonce.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return zenbtctypes.NonceData{}, fmt.Errorf("error getting last used ethereum nonce: %w", err)
		}
		nonceData = zenbtctypes.NonceData{Nonce: 0, PrevNonce: 0, Counter: 0, Skip: true}
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			return zenbtctypes.NonceData{}, fmt.Errorf("error setting last used ethereum nonce: %w", err)
		}
	}
	return nonceData, nil
}

// isNonceRequested checks if a nonce has been requested for the given key
func (k *Keeper) isNonceRequested(ctx sdk.Context, keyID uint64) (bool, error) {
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
	if requestedNonce != nonceData.PrevNonce {
		if err := nonceUpdatedCallback(tx); err != nil {
			return false, fmt.Errorf("nonce update callback error: %w", err)
		}
		k.Logger(ctx).Warn("nonce updated for key",
			"keyID", keyID,
			"requestedNonce", requestedNonce,
			"prevNonce", nonceData.PrevNonce,
			"currentNonce", nonceData.Nonce,
		)
		nonceData.PrevNonce = nonceData.Nonce
		if err := k.LastUsedEthereumNonce.Set(ctx, keyID, nonceData); err != nil {
			return false, fmt.Errorf("error setting last used Ethereum nonce: %w", err)
		}
		return true, nil
	}
	return false, nil
}

// updateNonces handles updating nonce state for keys used for minting and unstaking.
func (k *Keeper) updateNonces(ctx sdk.Context, oracleData OracleData) {
	for _, key := range k.getZenBTCKeyIDs(ctx) {
		isRequested, err := k.isNonceRequested(ctx, key)
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
		nonceData, err := k.getNonceDataWithInit(ctx, key)
		if err != nil {
			k.Logger(ctx).Error("error getting nonce data", "keyID", key, "error", err)
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

// processZenBTCStaking processes pending staking transactions.
func (k *Keeper) processZenBTCStaking(ctx sdk.Context, oracleData OracleData) {
	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetStakerKeyID(ctx),
		oracleData.RequestedStakerNonce,
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			return k.getPendingMintTransactionsByStatus(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED)
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			tx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED
			if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx); err != nil {
				return err
			}
			return k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), true)
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			if err := k.zenBTCKeeper.SetFirstPendingStakeTransaction(ctx, tx.Id); err != nil {
				return err
			}

			// Check for consensus
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedStakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC stake", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}

			k.Logger(ctx).Warn("processing zenBTC stake",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.RequestedStakerNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			unsignedTxHash, unsignedTx, err := k.constructStakeTx(
				ctx,
				getChainIDForEigen(ctx),
				tx.Amount,
				oracleData.RequestedStakerNonce,
				oracleData.EthGasLimit,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}

			return k.submitEthereumTransaction(
				ctx,
				tx.Creator,
				k.zenBTCKeeper.GetStakerKeyID(ctx),
				treasurytypes.WalletType(tx.ChainType),
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
		},
	)
}

// processZenBTCMintsEthereum processes pending mint transactions.
func (k *Keeper) processZenBTCMintsEthereum(ctx sdk.Context, oracleData OracleData) {
	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetEthMinterKeyID(ctx),
		oracleData.RequestedEthMinterNonce,
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			return k.getPendingMintTransactionsByStatus(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED)
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
			supply, err := k.zenBTCKeeper.GetSupply(ctx)
			if err != nil {
				return err
			}
			supply.PendingZenBTC -= tx.Amount
			supply.MintedZenBTC += tx.Amount
			if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
				return err
			}
			k.Logger(ctx).Warn("pending mint supply updated",
				"pending_mint_old", supply.PendingZenBTC+tx.Amount,
				"pending_mint_new", supply.PendingZenBTC,
			)
			k.Logger(ctx).Warn("minted supply updated",
				"minted_old", supply.MintedZenBTC-tx.Amount,
				"minted_new", supply.MintedZenBTC,
			)
			tx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
			return k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			if err := k.zenBTCKeeper.SetFirstPendingMintTransaction(ctx, tx.Id); err != nil {
				return err
			}

			// Check for consensus
			requiredFields := []VoteExtensionField{VEFieldRequestedEthMinterNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"zenBTC mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}

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
				oracleData.RequestedEthMinterNonce,
				oracleData.EthGasLimit,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}

			return k.submitEthereumTransaction(
				ctx,
				tx.Creator,
				k.zenBTCKeeper.GetEthMinterKeyID(ctx),
				treasurytypes.WalletType(tx.ChainType),
				chainID,
				unsignedMintTx,
				unsignedMintTxHash,
			)
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
	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetUnstakerKeyID(ctx),
		oracleData.RequestedUnstakerNonce,
		func(ctx sdk.Context) ([]zenbtctypes.BurnEvent, error) {
			return k.getPendingBurnEvents(ctx)
		},
		func(be zenbtctypes.BurnEvent) error {
			be.Status = zenbtctypes.BurnStatus_BURN_STATUS_UNSTAKING
			return k.zenBTCKeeper.SetBurnEvent(ctx, be.Id, be)
		},
		func(be zenbtctypes.BurnEvent) error {
			if err := k.zenBTCKeeper.SetFirstPendingBurnEvent(ctx, be.Id); err != nil {
				return err
			}

			// Check for consensus
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedUnstakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC burn unstake", fmt.Sprintf("burn_id: %d, destination: %s, amount: %d", be.Id, be.DestinationAddr, be.Amount)); err != nil {
				return err
			}

			k.Logger(ctx).Warn("processing zenBTC burn unstake",
				"burn_event", be,
				"nonce", oracleData.RequestedUnstakerNonce,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)
			unsignedTxHash, unsignedTx, err := k.constructUnstakeTx(
				ctx,
				getChainIDForEigen(ctx),
				be.DestinationAddr,
				be.Amount,
				oracleData.RequestedUnstakerNonce,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}

			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}

			return k.submitEthereumTransaction(
				ctx,
				creator,
				k.zenBTCKeeper.GetUnstakerKeyID(ctx),
				treasurytypes.WalletType_WALLET_TYPE_EVM,
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
		},
	)
}

// Helper function to submit Ethereum transactions
func (k *Keeper) submitEthereumTransaction(ctx sdk.Context, creator string, keyID uint64, walletType treasurytypes.WalletType, chainID uint64, unsignedTx []byte, unsignedTxHash []byte) error {
	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: chainID})
	if err != nil {
		return err
	}
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
	return err
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

		btcAmount := math.LegacyNewDecFromInt(math.NewIntFromUint64(redemption.Amount)).Mul(exchangeRate).TruncateInt64()
		// Convert zenBTC amount to BTC amount.
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
	}

	if foundNewRedemption {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting EthereumNonceRequested state", "error", err)
		}
	}
}

// processZenBTCRedemptions processes pending redemption completions.
func (k *Keeper) processZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	processZenBTCTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetCompleterKeyID(ctx),
		oracleData.RequestedCompleterNonce,
		func(ctx sdk.Context) ([]zenbtctypes.Redemption, error) {
			return k.getPendingRedemptions(ctx)
		},
		func(r zenbtctypes.Redemption) error {
			r.Status = zenbtctypes.RedemptionStatus_UNSTAKED
			if err := k.zenBTCKeeper.SetRedemption(ctx, r.Data.Id, r); err != nil {
				return err
			}
			return k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), true)
		},
		func(r zenbtctypes.Redemption) error {
			if err := k.zenBTCKeeper.SetFirstPendingRedemption(ctx, r.Data.Id); err != nil {
				return err
			}

			// Check for consensus
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedCompleterNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC redemption", fmt.Sprintf("redemption_id: %d, amount: %d", r.Data.Id, r.Data.Amount)); err != nil {
				return err
			}

			k.Logger(ctx).Warn("processing zenBTC complete",
				"id", r.Data.Id,
				"nonce", oracleData.RequestedCompleterNonce,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)
			unsignedTxHash, unsignedTx, err := k.constructCompleteTx(
				ctx,
				getChainIDForEigen(ctx),
				r.Data.Id,
				oracleData.RequestedCompleterNonce,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}

			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}

			return k.submitEthereumTransaction(
				ctx,
				creator,
				k.zenBTCKeeper.GetCompleterKeyID(ctx),
				treasurytypes.WalletType_WALLET_TYPE_EVM,
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
		},
	)
}
