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

	redemptionsHash, err := deriveHash(oracleData.Redemptions)
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
		RedemptionsHash:            redemptionsHash[:],
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

	voteExt, fieldVotePowers, totalVotePower, err := k.GetSuperMajorityVEData(ctx, req.Height, req.LocalLastCommit)
	if err != nil {
		k.Logger(ctx).Error("error retrieving field-based supermajority vote extension", "height", req.Height, "error", err)
		return nil, nil
	}

	if len(fieldVotePowers) == 0 { // no field reached consensus
		k.Logger(ctx).Warn("no fields reached consensus in vote extension", "height", req.Height)
		return k.marshalOracleData(req, &OracleData{ConsensusData: req.LocalLastCommit})
	}

	// Log essential fields information
	if !HasAllEssentialFields(fieldVotePowers) {
		missingEssentialFields := []string{}
		for _, field := range EssentialVoteExtensionFields {
			if _, ok := fieldVotePowers[field]; !ok {
				missingEssentialFields = append(missingEssentialFields, VEFieldString(field))
			}
		}

		if len(missingEssentialFields) > 0 {
			k.Logger(ctx).Warn("proceeding with partial consensus - missing essential vote extension fields",
				"height", req.Height,
				"fields_with_consensus", len(fieldVotePowers),
				"total_vote_power", totalVotePower,
				"missing_essential_fields", strings.Join(missingEssentialFields, ", "))
		} else {
			k.Logger(ctx).Info("consensus reached on all essential vote extension fields",
				"height", req.Height,
				"fields_with_consensus", len(fieldVotePowers),
				"total_vote_power", totalVotePower)
		}
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
	// Skip validation if this node is not a validator
	if !k.zrConfig.IsValidator {
		return ACCEPT_PROPOSAL, nil
	}

	// For block height 1, we don't have vote extensions.
	if req.Height == 1 {
		return ACCEPT_PROPOSAL, nil
	}

	if !VoteExtensionsEnabled(ctx) {
		k.Logger(ctx).Debug("vote extensions disabled; skipping oracle data validation", "height", req.Height)
		return ACCEPT_PROPOSAL, nil
	}

	// Check if we have any transactions to process
	if len(req.Txs) == 0 {
		k.Logger(ctx).Info("no transactions in proposal to process", "height", req.Height)
		return ACCEPT_PROPOSAL, nil
	}

	// Verify that the transaction contains a vote extension
	if !ContainsVoteExtension(req.Txs[0], k.txDecoder) {
		k.Logger(ctx).Warn("block does not contain vote extensions, rejecting proposal")
		return REJECT_PROPOSAL, nil
	}

	var oracleData OracleData
	if err := json.Unmarshal(req.Txs[0], &oracleData); err != nil {
		k.Logger(ctx).Error("error unmarshalling oracle data", "height", req.Height, "error", err)
		return REJECT_PROPOSAL, nil
	}

	// Check for empty oracle data - if it's empty, accept the proposal
	recoveredOracleDataNoCommitInfo := oracleData
	recoveredOracleDataNoCommitInfo.ConsensusData = abci.ExtendedCommitInfo{}
	if reflect.DeepEqual(recoveredOracleDataNoCommitInfo, OracleData{}) {
		k.Logger(ctx).Warn("accepting empty oracle data", "height", req.Height)
		return ACCEPT_PROPOSAL, nil
	}

	// Get the vote extension consensus based on fields (our "local consensus")
	ourVoteExt, fieldVotePowers, _, err := k.GetSuperMajorityVEData(ctx, req.Height, oracleData.ConsensusData)
	if err != nil {
		k.Logger(ctx).Error("error retrieving field-based supermajority vote extension", "height", req.Height, "error", err)
		return REJECT_PROPOSAL, nil
	}

	if len(fieldVotePowers) == 0 {
		// If we don't have local consensus, there's nothing to check against.
		// This is acceptable as we can still process blocks without oracle data.
		k.Logger(ctx).Warn("no fields reached local consensus in vote extension", "height", req.Height)
		return ACCEPT_PROPOSAL, nil
	}

	// If we have oracle data AND we're missing essential fields, reject the proposal
	if oracleData.HasAnyOracleData() {
		// List missing essential fields
		missingFields := []string{}
		for _, field := range EssentialVoteExtensionFields {
			if _, ok := fieldVotePowers[field]; !ok {
				missingFields = append(missingFields, VEFieldString(field))
			}
		}

		if len(missingFields) > 0 {
			k.Logger(ctx).Error("proposal contains oracle data but is missing consensus on essential fields",
				"height", req.Height,
				"missing_fields", strings.Join(missingFields, ", "))
			return REJECT_PROPOSAL, nil
		}
	}

	// Validate the oracle data against our local consensus
	if len(fieldVotePowers) > 0 {
		ourOracleData, err := k.getValidatedOracleData(ctx, ourVoteExt, fieldVotePowers)
		if err != nil {
			k.Logger(ctx).Error("error validating our oracle data", "height", req.Height, "error", err)
			return REJECT_PROPOSAL, nil
		}

		// If our consensus resulted in oracle data, then the proposal's oracle data should match
		if !oracleData.MatchesValidatedOracleData(ourOracleData) {
			k.Logger(ctx).Error("proposal's oracle data doesn't match our validated oracle data",
				"height", req.Height,
				"our_eth_height", ourOracleData.EthBlockHeight,
				"proposal_eth_height", oracleData.EthBlockHeight,
				"our_btc_height", ourOracleData.BtcBlockHeight,
				"proposal_btc_height", oracleData.BtcBlockHeight)
			return REJECT_PROPOSAL, nil
		}
	}

	k.Logger(ctx).Debug("proposal's oracle data validated successfully", "height", req.Height)
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
		k.processZenBTCMintsEthereum(ctx, oracleData)
		k.processZenBTCBurnEventsEthereum(ctx, oracleData)
		k.processZenBTCRedemptions(ctx, oracleData)
	}

	k.recordNonVotingValidators(ctx, req)
	k.recordMismatchedVoteExtensions(ctx, req.Height, voteExt, oracleData.ConsensusData)

	return nil
}

// shouldProcessOracleData checks if oracle data should be processed for this block.
func (k *Keeper) shouldProcessOracleData(ctx sdk.Context, req *abci.RequestFinalizeBlock) bool {
	// Skip processing if this node is not a validator
	if !k.zrConfig.IsValidator {
		k.Logger(ctx).Debug("not a validator node; skipping oracle data processing")
		return false
	}

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
	// For block height 1, we don't have vote extensions.
	if height == 1 {
		return VoteExtension{}, false
	}

	voteExt, fieldVotePowers, totalVotePower, err := k.GetSuperMajorityVEData(ctx, height, oracleData.ConsensusData)
	if err != nil {
		k.Logger(ctx).Error("error retrieving supermajority vote extensions data", "height", height, "error", err)
		return VoteExtension{}, false
	}

	if len(fieldVotePowers) == 0 {
		k.Logger(ctx).Warn("no consensus on any vote extension fields", "height", height)
		return VoteExtension{}, false
	}

	// Check if all essential fields have consensus when oracle data is present
	if oracleData.HasAnyOracleData() {
		// List missing essential fields
		missingFields := []string{}
		for _, field := range EssentialVoteExtensionFields {
			if _, ok := fieldVotePowers[field]; !ok {
				missingFields = append(missingFields, VEFieldString(field))
			}
		}

		if len(missingFields) > 0 {
			k.Logger(ctx).Warn("missing consensus on essential vote extension fields",
				"height", height,
				"missing_fields", strings.Join(missingFields, ", "),
				"fields_with_consensus", len(fieldVotePowers),
				"total_vote_power", totalVotePower)
			return VoteExtension{}, false
		}
	}

	if err := k.validateOracleData(voteExt, &oracleData, fieldVotePowers); err != nil {
		k.Logger(ctx).Error("error validating oracle data; won't store VE data", "height", height, "error", err)
		return VoteExtension{}, false
	}

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
		// If we don't have consensus on Ethereum height, get the latest state
		oracleData, err = k.GetSidecarState(ctx, ctx.BlockHeight())
		if err != nil {
			return nil, fmt.Errorf("error fetching latest oracle state: %w", err)
		}
	}

	// We only fetch Bitcoin data if we have consensus on BtcBlockHeight
	if _, ok := fieldVotePowers[VEFieldBtcBlockHeight]; ok {
		bitcoinData, err := k.sidecarClient.GetBitcoinBlockHeaderByHeight(
			ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{
				ChainName:   k.bitcoinNetwork(ctx),
				BlockHeight: voteExt.BtcBlockHeight,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error fetching bitcoin header: %w", err)
		}

		oracleData.BtcBlockHeight = bitcoinData.BlockHeight
		oracleData.BtcBlockHeader = *bitcoinData.BlockHeader
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

	if err := k.validateOracleData(voteExt, oracleData, fieldVotePowers); err != nil {
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
func (k *Keeper) storeBitcoinBlockHeader(ctx sdk.Context, oracleData OracleData) {
	k.Logger(ctx).Info("checking bitcoin header", "height", oracleData.BtcBlockHeight, "merkle", oracleData.BtcBlockHeader.MerkleRoot)

	if oracleData.BtcBlockHeight == 0 || oracleData.BtcBlockHeader.MerkleRoot == "" {
		k.Logger(ctx).Error("invalid bitcoin header data", "height", oracleData.BtcBlockHeight, "merkle", oracleData.BtcBlockHeader.MerkleRoot)
		return
	}

	requestedHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting requested historical Bitcoin headers", "error", err)
			return
		}
		k.Logger(ctx).Info("requested historical Bitcoin headers store not initialised", "height", oracleData.BtcBlockHeight)
	}

	k.Logger(ctx).Info("requested headers", "headers", requestedHeaders.Heights)

	isHistorical := k.isHistoricalHeader(oracleData.BtcBlockHeight, requestedHeaders.Heights)
	headerPreviouslySeen, err := k.BtcBlockHeaders.Has(ctx, oracleData.BtcBlockHeight)
	if err != nil {
		k.Logger(ctx).Error("error checking if Bitcoin header is already stored", "height", oracleData.BtcBlockHeight, "error", err)
		return
	}

	k.Logger(ctx).Info("header previously seen", "seen", headerPreviouslySeen, "isHistorical", isHistorical)

	if err := k.BtcBlockHeaders.Set(ctx, oracleData.BtcBlockHeight, oracleData.BtcBlockHeader); err != nil {
		k.Logger(ctx).Error("error storing Bitcoin header", "height", oracleData.BtcBlockHeight, "error", err)
		return
	}

	k.Logger(ctx).Info("stored header", "height", oracleData.BtcBlockHeight)

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

			// Check for consensus on required gas fields
			if !HasRequiredGasFields(oracleData.FieldVotePowers) {
				k.Logger(ctx).Error("cannot process zenBTC stake: missing consensus on gas fields",
					"tx_id", tx.Id,
					"recipient", tx.RecipientAddress,
					"amount", tx.Amount)
				return fmt.Errorf("missing consensus on gas fields required for transaction construction")
			}

			// Check for consensus on nonce field
			if !HasRequiredField(oracleData.FieldVotePowers, VEFieldRequestedStakerNonce) {
				k.Logger(ctx).Error("cannot process zenBTC stake: missing consensus on staker nonce",
					"tx_id", tx.Id,
					"recipient", tx.RecipientAddress,
					"amount", tx.Amount)
				return fmt.Errorf("missing consensus on staker nonce required for transaction construction")
			}

			k.Logger(ctx).Warn("processing zenBTC stake",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.RequestedStakerNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)
			unsignedStakeTxHash, unsignedStakeTx, err := k.constructStakeTx(
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
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
				ctx,
				&treasurytypes.MsgNewSignTransactionRequest{
					Creator:             tx.Creator,
					KeyId:               k.zenBTCKeeper.GetStakerKeyID(ctx),
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

			// Check for consensus on required gas fields
			if !HasRequiredGasFields(oracleData.FieldVotePowers) {
				k.Logger(ctx).Error("cannot process zenBTC mint: missing consensus on gas fields",
					"tx_id", tx.Id,
					"recipient", tx.RecipientAddress,
					"amount", tx.Amount)
				return fmt.Errorf("missing consensus on gas fields required for transaction construction")
			}

			// Check for consensus on nonce field
			if !HasRequiredField(oracleData.FieldVotePowers, VEFieldRequestedEthMinterNonce) {
				k.Logger(ctx).Error("cannot process zenBTC mint: missing consensus on minter nonce",
					"tx_id", tx.Id,
					"recipient", tx.RecipientAddress,
					"amount", tx.Amount)
				return fmt.Errorf("missing consensus on minter nonce required for transaction construction")
			}

			// Check for consensus on price data for fee calculation
			if !HasRequiredField(oracleData.FieldVotePowers, VEFieldBTCUSDPrice) ||
				!HasRequiredField(oracleData.FieldVotePowers, VEFieldETHUSDPrice) {
				k.Logger(ctx).Error("cannot process zenBTC mint: missing consensus on price data",
					"tx_id", tx.Id,
					"recipient", tx.RecipientAddress,
					"amount", tx.Amount)
				return fmt.Errorf("missing consensus on price data required for fee calculation")
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
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: chainID})
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
				ctx,
				&treasurytypes.MsgNewSignTransactionRequest{
					Creator:             tx.Creator,
					KeyId:               k.zenBTCKeeper.GetEthMinterKeyID(ctx),
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

			// Check for consensus on required gas fields
			if !HasRequiredGasFields(oracleData.FieldVotePowers) {
				k.Logger(ctx).Error("cannot process zenBTC burn unstake: missing consensus on gas fields",
					"burn_id", be.Id,
					"destination", be.DestinationAddr,
					"amount", be.Amount)
				return fmt.Errorf("missing consensus on gas fields required for transaction construction")
			}

			// Check for consensus on nonce field
			if !HasRequiredField(oracleData.FieldVotePowers, VEFieldRequestedUnstakerNonce) {
				k.Logger(ctx).Error("cannot process zenBTC burn unstake: missing consensus on unstaker nonce",
					"burn_id", be.Id,
					"destination", be.DestinationAddr,
					"amount", be.Amount)
				return fmt.Errorf("missing consensus on unstaker nonce required for transaction construction")
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
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
			if err != nil {
				return err
			}
			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
				ctx,
				&treasurytypes.MsgNewSignTransactionRequest{
					Creator:             creator,
					KeyId:               k.zenBTCKeeper.GetUnstakerKeyID(ctx),
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

			// Check for consensus on required gas fields
			if !HasRequiredGasFields(oracleData.FieldVotePowers) {
				k.Logger(ctx).Error("cannot process zenBTC redemption: missing consensus on gas fields",
					"redemption_id", r.Data.Id,
					"amount", r.Data.Amount)
				return fmt.Errorf("missing consensus on gas fields required for transaction construction")
			}

			// Check for consensus on nonce field
			if !HasRequiredField(oracleData.FieldVotePowers, VEFieldRequestedCompleterNonce) {
				k.Logger(ctx).Error("cannot process zenBTC redemption: missing consensus on completer nonce",
					"redemption_id", r.Data.Id,
					"amount", r.Data.Amount)
				return fmt.Errorf("missing consensus on completer nonce required for transaction construction")
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
			metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: getChainIDForEigen(ctx)})
			if err != nil {
				return err
			}
			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}
			_, err = k.treasuryKeeper.HandleSignTransactionRequest(
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
			)
			return err
		},
	)
}
