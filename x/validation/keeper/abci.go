package keeper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	solSystem "github.com/gagliardetto/solana-go/programs/system"
	solToken "github.com/gagliardetto/solana-go/programs/token"
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

	solNonce, err := k.collectSolanaNonces(ctx)
	if err != nil {
		return VoteExtension{}, err
	}
	solNonceHash, err := deriveHash(solNonce)
	if err != nil {
		return VoteExtension{}, err
	}
	solanaMintEventsHash, err := deriveHash(oracleData.SolanaMintEvents)

	solAccs, err := k.collectSolanaAccounts(ctx)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error collecting solana accounts: %w", err)
	}
	solAccsHash, err := deriveHash(solAccs)
	if err != nil {
		return VoteExtension{}, err
	}

	solanaBurnEventsHash, err := deriveHash(oracleData.SolanaBurnEvents)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving solana burn events hash: %w", err)
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
		SolanaMintNonceHashes:      solNonceHash[:],
		SolanaAccountsHash:         solAccsHash[:],
		SolanaMintEventsHash:       solanaMintEventsHash[:],
		SolanaBurnEventsHash:       solanaBurnEventsHash[:],
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
		// if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaBurnEventsHash) {
		// 	k.storeNewZenBTCBurnEventsSolana(ctx, oracleData)
		// }
		if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldRedemptionsHash) {
			k.storeNewZenBTCRedemptions(ctx, oracleData)
		}

		k.processZenBTCStaking(ctx, oracleData)
		k.processZenBTCMintsEthereum(ctx, oracleData)
		k.processZenBTCMintsSolana(ctx, oracleData)
		k.processSolanaZenBTCMintEvents(ctx, oracleData)
		k.processZenBTCBurnEvents(ctx, oracleData)
		k.processZenBTCRedemptions(ctx, oracleData)
		k.checkForRedemptionFulfilment(ctx)
		k.processSolanaZenBTCMintEvents(ctx, oracleData)
		k.processSolanaROCKMints(ctx, oracleData)
		k.processSolanaROCKMintEvents(ctx, oracleData)
		k.processSolanaROCKBurnEvents(ctx, oracleData)
		k.clearSolanaAccounts(ctx)
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

	k.validateOracleData(ctx, voteExt, &oracleData, fieldVotePowers)

	// Log final consensus summary after validation
	k.Logger(ctx).Info("final consensus summary",
		"fields_with_consensus", len(oracleData.FieldVotePowers),
		"stage", "post_validation")

	return voteExt, true
}

// getValidatedOracleData retrieves and validates oracle data based on a vote extension.
// Only validates fields that have reached consensus as indicated in fieldVotePowers.
func (k *Keeper) getValidatedOracleData(ctx sdk.Context, voteExt VoteExtension, fieldVotePowers map[VoteExtensionField]int64) (*OracleData, error) {
	// We only fetch Ethereum state if we have consensus on EthBlockHeight
	var oracleData *OracleData
	var err error

	if fieldHasConsensus(fieldVotePowers, VEFieldEthBlockHeight) {
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

	// Verify nonce fields and copy them if they have consensus
	nonceFields := []struct {
		field       VoteExtensionField
		keyID       uint64
		oracleField *uint64
	}{
		{VEFieldRequestedStakerNonce, k.zenBTCKeeper.GetStakerKeyID(ctx), &oracleData.RequestedStakerNonce},
		{VEFieldRequestedEthMinterNonce, k.zenBTCKeeper.GetEthMinterKeyID(ctx), &oracleData.RequestedEthMinterNonce},
		{VEFieldRequestedUnstakerNonce, k.zenBTCKeeper.GetUnstakerKeyID(ctx), &oracleData.RequestedUnstakerNonce},
		{VEFieldRequestedCompleterNonce, k.zenBTCKeeper.GetCompleterKeyID(ctx), &oracleData.RequestedCompleterNonce},
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
				}
				*nf.oracleField = currentNonce
			}
		}
	}

	if fieldHasConsensus(fieldVotePowers, VEFieldSolanaMintNoncesHash) {
		oracleData.SolanaMintNonces, err = k.collectSolanaNonces(ctx)
		if err != nil {
			return nil, fmt.Errorf("error collecting solana nonces: %w", err)
		}

	}

	if fieldHasConsensus(fieldVotePowers, VEFieldSolanaAccountsHash) {
		solAccs, err := k.collectSolanaAccounts(ctx)
		if err != nil {
			return nil, fmt.Errorf("error collecting solana accounts: %w", err)
		}
		oracleData.SolanaAccounts = solAccs
	}
	// Store the field vote powers for later use in transaction dispatch callbacks
	oracleData.FieldVotePowers = fieldVotePowers

	// Call the standard validateOracleData to check other fields
	k.validateOracleData(ctx, voteExt, oracleData, fieldVotePowers)

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

		validator.TokensAVS = sdkmath.Int(delegation.Stake)

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

	if err := k.ValidatorDelegations.Walk(ctx, nil, func(valAddr string, stake sdkmath.Int) (bool, error) {
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

	validator.TokensAVS = sdkmath.ZeroInt()

	if err = k.SetValidator(ctx, validator); err != nil {
		return fmt.Errorf("error updating validator after removal: %w", err)
	}

	return nil
}

// updateAVSDelegationStore updates the AVS delegation store with new delegation amounts.
func (k *Keeper) updateAVSDelegationStore(ctx sdk.Context, oracleData OracleData) {
	for validatorAddr, delegatorMap := range oracleData.EigenDelegationsMap {
		for delegatorAddr, amount := range delegatorMap {
			if err := k.AVSDelegations.Set(ctx, collections.Join(validatorAddr, delegatorAddr), sdkmath.NewIntFromBigInt(amount)); err != nil {
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
	requestedEthNonce *uint64,
	requestedSolNonce *solSystem.NonceAccount,
	nonceReqStore collections.Map[uint64, bool],
	pendingTxs []T,
	txDispatchCallback func(tx T) error,
	nonceUpdatedCallback func(tx T) error,
) {
	if len(pendingTxs) == 0 {
		return
	}

	nonceUpdated := false

	if requestedEthNonce != nil {
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
			"requested", requestedEthNonce,
		)
		if nonceData.Nonce != 0 && *requestedEthNonce == 0 {
			return
		}

		nonceUpdated, err = handleNonceUpdate(k, ctx, keyID, *requestedEthNonce, nonceData, pendingTxs[0], nonceUpdatedCallback)
		if err != nil {
			k.Logger(ctx).Error("error handling nonce update", "keyID", keyID, "error", err)
			return
		}

		if len(pendingTxs) == 1 && nonceUpdated {
			if err := k.clearNonceRequest(ctx, nonceReqStore, keyID); err != nil {
				k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
			}
			return
		}

		if nonceData.Skip {
			return
		}
	} else if requestedSolNonce != nil {
		k.Logger(ctx).Error("requested solana nonce", "nonce", requestedSolNonce.Nonce)

		if requestedSolNonce.Nonce.IsZero() {
			k.Logger(ctx).Error("solana nonce is zero")
			return
		}

		if len(pendingTxs) == 0 {
			if err := k.clearNonceRequest(ctx, nonceReqStore, keyID); err != nil {
				k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
			}
			return
		}

		if err := nonceUpdatedCallback(pendingTxs[0]); err != nil {
			k.Logger(ctx).Error("error handling nonce update", "keyID", keyID, "error", err)
			return
		}

		k.Logger(ctx).Error("solana nonce updated", "keyID", keyID, "nonce", requestedSolNonce.Nonce)
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

// processTransaction is a generic helper that encapsulates the common logic for nonce update and tx dispatch.
func processTransaction[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedEthNonce *uint64,
	requestedSolNonce *solSystem.NonceAccount,
	pendingGetter func(ctx sdk.Context) ([]T, error),
	txDispatchCallback func(tx T) error,
	nonceUpdatedCallback func(tx T) error,
) {
	nonceReqStore := k.EthereumNonceRequested
	if requestedEthNonce == nil {
		nonceReqStore = k.SolanaNonceRequested
	}

	isRequested, err := isNonceRequested(ctx, nonceReqStore, keyID)
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
		if err := k.clearNonceRequest(ctx, nonceReqStore, keyID); err != nil {
			k.Logger(ctx).Error("error clearing ethereum nonce request", "keyID", keyID, "error", err)
		}
		return
	}
	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedEthNonce, requestedSolNonce, nonceReqStore, pendingTxs, txDispatchCallback, nonceUpdatedCallback)
}

// getPendingTransactions is a generic helper that walks a collections.Map with key type uint64
// and returns a slice of items of type T that satisfy the provided predicate, up to a given limit.
// If limit is 0, all matching items will be returned.
func getPendingTransactions[T any](ctx sdk.Context, store collections.Map[uint64, T], predicate func(T) bool, firstPendingID uint64, limit int) ([]T, error) {
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
func isNonceRequested(ctx sdk.Context, store collections.Map[uint64, bool], keyID uint64) (bool, error) {
	requested, err := store.Get(ctx, keyID)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("error getting nonce request state: %w", err)
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
		isRequested, err := isNonceRequested(ctx, k.EthereumNonceRequested, key)
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

// clearNonceRequest resets the nonce-request flag for a given key.
func (k *Keeper) clearNonceRequest(ctx sdk.Context, store collections.Map[uint64, bool], keyID uint64) error {
	k.Logger(ctx).Warn("set requested nonce state to false", "keyID", keyID)
	return store.Set(ctx, keyID, false)
}

// processZenBTCStaking processes pending staking transactions.
func (k *Keeper) processZenBTCStaking(ctx sdk.Context, oracleData OracleData) {
	processTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetStakerKeyID(ctx),
		&oracleData.RequestedStakerNonce,
		nil,
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			return k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				zenbtctypes.WalletType_WALLET_TYPE_UNSPECIFIED,
			)
		},
		// Dispatch stake transaction
		func(tx zenbtctypes.PendingMintTransaction) error {
			if err := k.zenBTCKeeper.SetFirstPendingStakeTransaction(ctx, tx.Id); err != nil {
				return err
			}

			// Check for consensus
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedStakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC stake", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}

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

			k.Logger(ctx).Warn("processing zenBTC stake",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.RequestedStakerNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			return k.submitEthereumTransaction(
				ctx,
				tx.Creator,
				k.zenBTCKeeper.GetStakerKeyID(ctx),
				treasurytypes.WalletType_WALLET_TYPE_EVM, //treasurytypes.WalletType(tx.ChainType),
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
		},
		// Successfully processed stake transaction
		func(tx zenbtctypes.PendingMintTransaction) error {
			tx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED
			if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx); err != nil {
				return err
			}
			if types.IsSolanaCAIP2(tx.Caip2ChainId) {
				solParams := k.zenBTCKeeper.GetSolanaParams(ctx)
				if err := k.SolanaNonceRequested.Set(ctx, solParams.NonceAccountKey, true); err != nil {
					return err
				}
				if err := k.SetSolanaRequestedAccount(ctx, tx.RecipientAddress, true); err != nil {
					return err
				}
				k.Logger(ctx).Error("processed zenbtc stake", "tx_id", tx.Id, "recipient", tx.RecipientAddress, "amount", tx.Amount)
				return nil
			} else if types.IsEthereumCAIP2(tx.Caip2ChainId) {
				return k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), true)
			}
			return fmt.Errorf("unsupported chain type for chain ID: %s", tx.Caip2ChainId)
		},
	)
}

// processZenBTCMintsEthereum processes pending mint transactions.
func (k *Keeper) processZenBTCMintsEthereum(ctx sdk.Context, oracleData OracleData) {
	processTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetEthMinterKeyID(ctx),
		&oracleData.RequestedEthMinterNonce,
		nil,
		// Get pending mint transactions
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			return k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED,
				zenbtctypes.WalletType_WALLET_TYPE_EVM,
			)
		},
		// Dispatch mint transaction
		func(tx zenbtctypes.PendingMintTransaction) error {
			if err := k.zenBTCKeeper.SetFirstPendingEthMintTransaction(ctx, tx.Id); err != nil {
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

			// Get decimal values from string representations
			btcUSDPrice, err := sdkmath.LegacyNewDecFromStr(oracleData.BTCUSDPrice)
			if err != nil || btcUSDPrice.IsNil() || btcUSDPrice.IsZero() {
				k.Logger(ctx).Error("invalid BTC/USD price", "error", err)
				return nil
			}
			ethUSDPrice, err := sdkmath.LegacyNewDecFromStr(oracleData.ETHUSDPrice)
			if err != nil || ethUSDPrice.IsNil() || ethUSDPrice.IsZero() {
				k.Logger(ctx).Error("invalid ETH/USD price", "error", err)
				return nil
			}

			feeZenBTC := k.CalculateZenBTCMintFee(
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
				oracleData.EthGasLimit,
				btcUSDPrice,
				ethUSDPrice,
				exchangeRate,
			)

			chainID, err := types.ValidateChainID(ctx, tx.Caip2ChainId)
			if err != nil {
				return fmt.Errorf("unsupported chain ID: %w", err)
			}

			unsignedMintTxHash, unsignedMintTx, err := k.constructMintTx(
				ctx,
				tx.RecipientAddress,
				chainID.Uint64(),
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

			k.Logger(ctx).Warn("processing zenBTC mint",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.RequestedEthMinterNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			return k.submitEthereumTransaction(
				ctx,
				tx.Creator,
				k.zenBTCKeeper.GetEthMinterKeyID(ctx),
				treasurytypes.WalletType(tx.ChainType),
				chainID.Uint64(),
				unsignedMintTx,
				unsignedMintTxHash,
			)
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
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
	)
}

// processZenBTCMintsSolana processes pending mint transactions.
func (k *Keeper) processZenBTCMintsSolana(ctx sdk.Context, oracleData OracleData) {
	processTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey,
		nil,
		oracleData.SolanaMintNonces[k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey],
		// Get pending mint transactions
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			pendingMints, err := k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED,
				zenbtctypes.WalletType_WALLET_TYPE_SOLANA,
			)
			k.Logger(ctx).Warn("pending zenbtc solana mints", "mints", fmt.Sprintf("%v", pendingMints), "count", len(pendingMints))
			return pendingMints, err
		},
		// Dispatch mint transaction
		func(tx zenbtctypes.PendingMintTransaction) error {
			k.Logger(ctx).Error("dispatch handler triggered", "tx_id", tx.Id, "recipient", tx.RecipientAddress, "amount", tx.Amount)
			if tx.BlockHeight > 0 {
				k.Logger(ctx).Info("waiting for pending zenbtc solana mint tx", "tx_id", tx.Id, "block_height", tx.BlockHeight)
				return nil
			}
			if err := k.zenBTCKeeper.SetFirstPendingSolMintTransaction(ctx, tx.Id); err != nil {
				return err
			}

			if len(oracleData.SolanaMintNonces) == 0 {
				return fmt.Errorf("no nonce available for zenbtc solana mint")

			}
			// Check for consensus
			requiredFields := []VoteExtensionField{
				VEFieldSolanaMintNoncesHash,
				VEFieldBTCUSDPrice,
				VEFieldETHUSDPrice,
				VEFieldSolanaAccountsHash,
			}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"zenBTC mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}

			// Get decimal values from string representations
			btcUSDPrice, err := sdkmath.LegacyNewDecFromStr(oracleData.BTCUSDPrice)
			if err != nil || btcUSDPrice.IsNil() || btcUSDPrice.IsZero() {
				k.Logger(ctx).Error("invalid BTC/USD price", "error", err)
				return nil
			}
			ethUSDPrice, err := sdkmath.LegacyNewDecFromStr(oracleData.ETHUSDPrice)
			if err != nil || ethUSDPrice.IsNil() || ethUSDPrice.IsZero() {
				k.Logger(ctx).Error("invalid ETH/USD price", "error", err)
				return nil
			}

			solParams := k.zenBTCKeeper.GetSolanaParams(ctx)
			txPrepReq := &solanaMintTxRequest{}
			// add a solana instruction if we need to fund the ata
			ata, ok := oracleData.SolanaAccounts[tx.RecipientAddress]
			if !ok {
				return fmt.Errorf("ata account not retrieved for address: %s", tx.RecipientAddress)
			}
			if ata.State == solToken.Uninitialized {
				txPrepReq.fundReceiver = true
			}

			txPrepReq.amount = tx.Amount
			txPrepReq.fee = solParams.Fee
			txPrepReq.recipient = tx.RecipientAddress
			txPrepReq.nonce = oracleData.SolanaMintNonces[solParams.NonceAccountKey]
			txPrepReq.programID = solParams.ProgramId
			txPrepReq.mintAddress = solParams.MintAddress
			txPrepReq.feeWallet = solParams.FeeWallet
			txPrepReq.nonceAccountKey = solParams.NonceAccountKey
			txPrepReq.nonceAuthorityKey = solParams.NonceAuthorityKey
			txPrepReq.signerKey = solParams.SignerKeyId
			txPrepReq.zenbtc = true
			k.Logger(ctx).Error("processing zenbtc solana mint", "tx_id", tx.Id, "recipient", tx.RecipientAddress, "amount", tx.Amount)
			transaction, err := k.PrepareSolanaMintTx(ctx, txPrepReq)
			if err != nil {
				return fmt.Errorf("PrepareSolRockMintTx: %w", err)
			}

			k.Logger(ctx).Warn("processing zenBTC mint",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.SolanaMintNonces[solParams.NonceAccountKey],
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)

			txID, err := k.submitSolanaTransaction(
				ctx,
				tx.Creator,
				[]uint64{solParams.SignerKeyId, solParams.NonceAuthorityKey},
				treasurytypes.WalletType(tx.ChainType),
				tx.Caip2ChainId,
				transaction,
			)
			if err != nil {
				return err
			}
			tx.ZrchainTxId = txID
			tx.BlockHeight = ctx.BlockHeight()
			k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
			nonce := types.SolanaNonce{Nonce: oracleData.SolanaMintNonces[solParams.NonceAccountKey].Nonce[:]}
			k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, nonce)
			return nil
		},
		func(tx zenbtctypes.PendingMintTransaction) error {
			// If BlockHeight is 0, this transaction was either just dispatched in the current block
			// by the txDispatchCallback, or it has been reset for a full retry by prior logic in this callback.
			// It doesn't need BTL/event checks *yet*.
			if tx.BlockHeight == 0 {
				k.Logger(ctx).Debug("Solana Mint Nonce Update: tx.BlockHeight is 0. No BTL/event check in this invocation.", "tx_id", tx.Id)
				return k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
			}

			solParams := k.zenBTCKeeper.GetSolanaParams(ctx)
			k.Logger(ctx).Info("Solana Mint Status Check Begin", "tx_id", tx.Id, "recipient", tx.RecipientAddress, "amount", tx.Amount, "tx_block_height", tx.BlockHeight, "btl", solParams.Btl, "current_chain_height", ctx.BlockHeight(), "awaiting_event_since", tx.AwaitingEventSince)

			// --- Primary BTL Timeout Check (Blocks To Live) ---
			if ctx.BlockHeight() > tx.BlockHeight+solParams.Btl {
				tx = k.processBtlSolanaMint(ctx, tx, oracleData, *solParams)
			}

			// --- Secondary Event Arrival Timeout Check ---
			// This check applies if tx.AwaitingEventSince was set (indicating nonce advanced)
			// and was not subsequently cleared by the BTL check itself.
			if tx.AwaitingEventSince > 0 {
				tx = k.processSecondaryTimeoutSolanaMint(ctx, tx, oracleData, *solParams)
			}

			k.Logger(ctx).Info("Solana Mint Status Check End", "tx_id", tx.Id, "tx_block_height_after_checks", tx.BlockHeight, "awaiting_event_since_after_checks", tx.AwaitingEventSince)
			// Persist any modifications to the transaction state.
			return k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
		},
	)
}

// processSolanaROCKMints processes pending mint transactions.
func (k *Keeper) processSolanaROCKMints(ctx sdk.Context, oracleData OracleData) {
	processTransaction(
		k,
		ctx,
		k.zentpKeeper.GetSolanaParams(ctx).NonceAccountKey,
		nil,
		oracleData.SolanaMintNonces[k.zentpKeeper.GetSolanaParams(ctx).NonceAccountKey],
		func(ctx sdk.Context) ([]*zentptypes.Bridge, error) {
			mints, err := k.zentpKeeper.GetMintsWithStatus(ctx, zentptypes.BridgeStatus_BRIDGE_STATUS_PENDING)
			return mints, err
		},
		func(tx *zentptypes.Bridge) error {
			// Check whether this tx has already been processed, if it has been - we wait for it to complete (or timeout)
			if tx.BlockHeight > 0 {
				k.Logger(ctx).Info("waiting for pending zentp solana mint tx", "tx_id", tx.Id, "block_height", tx.BlockHeight)
				return nil
			}

			// Check for consensus
			requiredFields := []VoteExtensionField{VEFieldSolanaAccountsHash}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"solROCK mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return fmt.Errorf("validateConsensusForTxFields: %w", err)
			}
			val, err := k.SolanaAccountsRequested.Get(ctx, tx.RecipientAddress)
			if err == nil {
				_ = val // TODO: fix!
			}

			// add a solana instruction if we need to fund the ata
			fundReceiver := false
			ata, ok := oracleData.SolanaAccounts[tx.RecipientAddress]
			if !ok {
				return fmt.Errorf("ata account not retrieved for address: %s", tx.RecipientAddress)
			}
			if ata.State == solToken.Uninitialized {
				fundReceiver = true
			}
			solParams := k.zentpKeeper.GetSolanaParams(ctx)
			transaction, err := k.PrepareSolanaMintTx(ctx, &solanaMintTxRequest{
				amount:            tx.Amount,
				fee:               solParams.Fee,
				recipient:         tx.RecipientAddress,
				nonce:             oracleData.SolanaMintNonces[solParams.NonceAccountKey],
				fundReceiver:      fundReceiver,
				programID:         solParams.ProgramId,
				mintAddress:       solParams.MintAddress,
				feeWallet:         solParams.FeeWallet,
				nonceAccountKey:   solParams.NonceAccountKey,
				nonceAuthorityKey: solParams.NonceAuthorityKey,
				signerKey:         solParams.SignerKeyId,
				rock:              true,
			})
			if err != nil {
				return fmt.Errorf("PrepareSolRockMintTx: %w", err)
			}

			id, err := k.submitSolanaTransaction(
				ctx,
				tx.Creator,
				[]uint64{solParams.SignerKeyId, solParams.NonceAuthorityKey},
				treasurytypes.WalletType_WALLET_TYPE_SOLANA,
				tx.DestinationChain,
				transaction,
			)
			if err != nil {
				return fmt.Errorf("submitSolanaTransaction: %w", err)
			}
			tx.State = zentptypes.BridgeStatus_BRIDGE_STATUS_PENDING
			tx.TxId = id
			tx.BlockHeight = ctx.BlockHeight()
			k.zentpKeeper.UpdateMint(ctx, tx.Id, tx)
			nonce := types.SolanaNonce{Nonce: oracleData.SolanaMintNonces[solParams.NonceAccountKey].Nonce[:]}
			k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, nonce)
			return k.zentpKeeper.UpdateMint(ctx, tx.Id, tx)
		},
		func(tx *zentptypes.Bridge) error {
			if tx.BlockHeight == 0 {
				return nil
			}
			solParams := k.zentpKeeper.GetSolanaParams(ctx)
			if ctx.BlockHeight() > tx.BlockHeight+solParams.Btl {
				currentNonce := oracleData.SolanaMintNonces[solParams.NonceAccountKey].Nonce
				lastNonce, err := k.LastUsedSolanaNonce.Get(ctx, solParams.NonceAccountKey)
				if err != nil {
					k.Logger(ctx).Error("error getting last used solana nonce", "error", err)
					return err
				}
				k.Logger(ctx).Error("nonces", "last_hex", hex.EncodeToString(lastNonce.Nonce), "current_hex", hex.EncodeToString(currentNonce[:]))
				if bytes.Equal(currentNonce[:], lastNonce.Nonce[:]) {
					tx.BlockHeight = 0 // this will trigger the tx to get retried
				}
				// else the transaction has been included in a block, and we should wait for the mint event
			}
			return nil
		},
	)

}

// processROCKBurns processes pending mint transactions.
func (k *Keeper) processSolanaROCKMintEvents(ctx sdk.Context, oracleData OracleData) {

	protocolWalletAddress, bridgeFee, err := k.zentpKeeper.GetBridgeFeeParams(ctx)
	if err != nil {
		k.Logger(ctx).Error("GetBridgeFeeParams: ", err.Error())
		return
	}

	pendingMints, err := k.zentpKeeper.GetMintsWithStatus(ctx, zentptypes.BridgeStatus_BRIDGE_STATUS_PENDING)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return
		}
		k.Logger(ctx).Error("GetMintsWithStatus: ", err.Error())
		return
	}

	if len(pendingMints) == 0 {
		return
	}

	for _, pendingMint := range pendingMints {
		tx, err := k.treasuryKeeper.SignTransactionRequestStore.Get(ctx, pendingMint.TxId)
		if err != nil {
			k.Logger(ctx).Error("SignTransactionRequestStore.Get: ", err.Error())
			return
		}
		sigReq, err := k.treasuryKeeper.SignRequestStore.Get(ctx, tx.SignRequestId)
		if err != nil {
			k.Logger(ctx).Error("SignRequestStore.Get: ", err.Error())
		}

		bridgeFeeCoins, err := k.zentpKeeper.GetBridgeFeeAmount(ctx, pendingMint.Amount, bridgeFee)
		if err != nil {
			k.Logger(ctx).Error("GetBridgeFeeAmount: ", err.Error())
			return
		}

		bridgeAmount := sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(pendingMint.Amount)))

		var (
			signatures []byte
			sigHash    [32]byte
		)

		for _, id := range sigReq.ChildReqIds {
			childReq, err := k.treasuryKeeper.SignRequestStore.Get(ctx, id)
			if err != nil {
				k.Logger(ctx).Error("SignRequestStore.Get: ", err.Error())
				return
			}
			if len(childReq.SignedData) != 1 {
				continue
			}
			signatures = append(signatures, childReq.SignedData[0].SignedData...)
		}
		sigHash = sha256.Sum256(signatures)
		for _, event := range oracleData.SolanaMintEvents {
			if bytes.Equal(event.SigHash, sigHash[:]) {

				err = k.bankKeeper.BurnCoins(ctx, zentptypes.ModuleName, bridgeAmount)
				if err != nil {
					k.Logger(ctx).Error("Failed to burn coins", "denom", pendingMint.Denom, "error", err.Error())
					return
				}

				err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zentptypes.ModuleName, protocolWalletAddress, bridgeFeeCoins)
				if err != nil {
					k.Logger(ctx).Error("SendCoinsFromModuleToAccount: ", err.Error())
					return
				}
				pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED
				err = k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint)
				if err != nil {
					k.Logger(ctx).Error("UpdateMint: ", err.Error())
				}

				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeValidation,
						sdk.NewAttribute(types.AttributeKeyBurnAmount, bridgeAmount.String()),
						sdk.NewAttribute(types.AttributeKeyBridgeFee, bridgeFeeCoins.String()),
					),
				)
			}
		}
	}
}

// processROCKBurns processes pending mint transactions.
func (k *Keeper) processSolanaZenBTCMintEvents(ctx sdk.Context, oracleData OracleData) {
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Started.", "oracle_event_count", len(oracleData.SolanaMintEvents))

	firstPendingID, err := k.zenBTCKeeper.GetFirstPendingSolMintTransaction(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: No first pending Solana mint transaction ID found. Nothing to process.")
			return
		}
		k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting first pending Solana mint transaction ID.", "error", err)
		return
	}

	if firstPendingID == 0 {
		k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: First pending Solana mint transaction ID is 0. Nothing to process.")
		return
	}
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Processing with first pending ID.", "first_pending_id", firstPendingID)

	pendingMint, err := k.zenBTCKeeper.GetPendingMintTransactionsStore().Get(ctx, firstPendingID)
	if err != nil {
		k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting pending mint transaction from store.", "id", firstPendingID, "error", err)
		return
	}
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Retrieved pending mint transaction.", "pending_mint_id", pendingMint.Id, "zrchain_tx_id", pendingMint.ZrchainTxId, "status", pendingMint.Status, "recipient", pendingMint.RecipientAddress, "amount", pendingMint.Amount)

	if pendingMint.ZrchainTxId == 0 {
		k.Logger(ctx).Warn("ProcessSolanaZenBTCMintEvents: PendingMint has ZrchainTxId == 0. Cannot match with treasury sign requests. Skipping.", "pending_mint_id", pendingMint.Id)
		return
	}

	signTxReq, err := k.treasuryKeeper.SignTransactionRequestStore.Get(ctx, pendingMint.ZrchainTxId)
	if err != nil {
		k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting SignTransactionRequest from treasury.", "zrchain_tx_id_searched", pendingMint.ZrchainTxId, "error", err)
		return
	}
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Retrieved SignTransactionRequest from treasury.", "zrchain_tx_id", signTxReq.Id, "sign_request_id", signTxReq.SignRequestId)

	mainSignReq, err := k.treasuryKeeper.SignRequestStore.Get(ctx, signTxReq.SignRequestId)
	if err != nil {
		k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting main SignRequest from treasury.", "sign_request_id_searched", signTxReq.SignRequestId, "error", err)
		return
	}
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Retrieved main SignRequest from treasury.", "main_sign_request_id", mainSignReq.Id, "child_req_count", len(mainSignReq.ChildReqIds), "status", mainSignReq.Status)

	var signatures [][]byte
	foundAllChildSignatures := true
	for i, childReqID := range mainSignReq.ChildReqIds {
		childReq, err := k.treasuryKeeper.SignRequestStore.Get(ctx, childReqID)
		if err != nil {
			k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting child SignRequest.", "child_req_id", childReqID, "error", err)
			foundAllChildSignatures = false
			break
		}
		if len(childReq.SignedData) == 0 || len(childReq.SignedData[0].SignedData) == 0 {
			k.Logger(ctx).Warn("ProcessSolanaZenBTCMintEvents: Child SignRequest has no signed data or empty signature.", "child_req_id", childReqID, "signed_data_count", len(childReq.SignedData))
			foundAllChildSignatures = false
			break
		}
		signatures = append(signatures, childReq.SignedData[0].SignedData)
		k.Logger(ctx).Debug("ProcessSolanaZenBTCMintEvents: Appended signature from child request.", "child_idx", i, "child_req_id", childReqID, "signature_hex", hex.EncodeToString(childReq.SignedData[0].SignedData))
	}

	if !foundAllChildSignatures {
		k.Logger(ctx).Warn("ProcessSolanaZenBTCMintEvents: Did not find all child signatures or some were empty. Cannot compute sigHash.", "main_sign_request_id", mainSignReq.Id)
		return
	}

	if len(signatures) == 0 {
		k.Logger(ctx).Warn("ProcessSolanaZenBTCMintEvents: No signatures collected from child requests. Cannot compute sigHash.", "main_sign_request_id", mainSignReq.Id)
		return
	}

	concatenatedSignatures := bytes.Join(signatures, []byte{})
	sigHash := sha256.Sum256(concatenatedSignatures)
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Computed local sigHash.", "concatenated_signature_len", len(concatenatedSignatures), "local_sig_hash_hex", hex.EncodeToString(sigHash[:]))

	for i, event := range oracleData.SolanaMintEvents {
		k.Logger(ctx).Debug("ProcessSolanaZenBTCMintEvents: Comparing with oracle event.", "oracle_event_idx", i, "oracle_sig_hash_hex", hex.EncodeToString(event.SigHash))
		if bytes.Equal(event.SigHash, sigHash[:]) {
			k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: MATCH FOUND! Oracle event sigHash matches local sigHash.", "oracle_event_idx", i, "pending_mint_id", pendingMint.Id)

			supply, err := k.zenBTCKeeper.GetSupply(ctx)
			if err != nil {
				k.Logger(ctx).Error("zenBTCKeeper.GetSupply: ", err.Error())
				return
			}
			supply.PendingZenBTC -= pendingMint.Amount
			supply.MintedZenBTC += pendingMint.Amount
			if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
				k.Logger(ctx).Error("zenBTCKeeper.SetSupply: ", err.Error())
				return
			}
			k.Logger(ctx).Warn("pending mint supply updated",
				"pending_mint_old", supply.PendingZenBTC+pendingMint.Amount,
				"pending_mint_new", supply.PendingZenBTC,
			)
			k.Logger(ctx).Warn("minted supply updated",
				"minted_old", supply.MintedZenBTC-pendingMint.Amount,
				"minted_new", supply.MintedZenBTC,
			)
			pendingMint.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
			if err = k.zenBTCKeeper.SetPendingMintTransaction(ctx, pendingMint); err != nil {
				k.Logger(ctx).Error("zenBTCKeeper.SetPendingMintTransaction: ", err.Error())
			}
			if err = k.zenBTCKeeper.SetFirstPendingSolMintTransaction(ctx, 0); err != nil {
				k.Logger(ctx).Error("zenBTCKeeper.SetFirstPendingSolMintTransaction: ", err.Error())
			}
		}
	}
}

// storeNewZenBTCBurnEventsEthereum stores new burn events coming from Ethereum.
func (k *Keeper) storeNewZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	k.storeNewZenBTCBurnEvents(ctx, oracleData.EthBurnEvents, "ethereum", "error setting EthereumNonceRequested state")
}

// storeNewZenBTCBurnEventsSolana stores new burn events coming from Solana.
func (k *Keeper) storeNewZenBTCBurnEventsSolana(ctx sdk.Context, oracleData OracleData) {
	k.storeNewZenBTCBurnEvents(ctx, oracleData.SolanaBurnEvents, "solana", "error setting EthereumNonceRequested state for unstaker")
}

// storeNewZenBTCBurnEvents is a helper function to store new burn events from a given source.
func (k *Keeper) storeNewZenBTCBurnEvents(ctx sdk.Context, burnEvents []sidecarapitypes.BurnEvent, source string, nonceErrorMsg string) {
	k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: Started.", "source", source, "incoming_event_count", len(burnEvents))
	if source == "solana" && len(burnEvents) > 0 {
		for i, dbgEvent := range burnEvents {
			k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: Solana event details from oracle.",
				"source", source, "idx", i, "tx_id", dbgEvent.TxID, "log_idx", dbgEvent.LogIndex,
				"chain_id", dbgEvent.ChainID, "amount", dbgEvent.Amount, "destination_addr_hex", hex.EncodeToString(dbgEvent.DestinationAddr))
		}
	}

	foundNewBurn := false
	// Loop over each burn event from oracle to check for new ones.
	for _, burn := range burnEvents {
		// Check if this burn event already exists
		exists := false
		if err := k.zenBTCKeeper.WalkBurnEvents(ctx, func(id uint64, existingBurn zenbtctypes.BurnEvent) (bool, error) {
			// Compare fields from the input burn event data with the stored BurnEvent
			if existingBurn.TxID == burn.TxID &&
				existingBurn.LogIndex == burn.LogIndex &&
				existingBurn.ChainID == burn.ChainID {
				k.Logger(ctx).Debug("StoreNewZenBTCBurnEvents: Event already exists in store.", "source", source, "tx_id", burn.TxID, "log_idx", burn.LogIndex, "chain_id", burn.ChainID, "existing_burn_id", id)
				exists = true
				return true, nil
			}
			return false, nil // Continue walking
		}); err != nil {
			k.Logger(ctx).Error("StoreNewZenBTCBurnEvents: Error walking burn events. Skipping event.", "source", source, "tx_id", burn.TxID, "error", err)
			continue // Process next event
		}

		if !exists {
			k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: New event, creating BurnEvent.", "source", source, "tx_id", burn.TxID, "log_idx", burn.LogIndex, "chain_id", burn.ChainID, "amount", burn.Amount, "destination_addr_hex", hex.EncodeToString(burn.DestinationAddr))
			// Create a new BurnEvent using data from the input struct
			newBurn := zenbtctypes.BurnEvent{
				TxID:            burn.TxID,
				LogIndex:        burn.LogIndex,
				ChainID:         burn.ChainID,
				DestinationAddr: burn.DestinationAddr,
				Amount:          burn.Amount,
				Status:          zenbtctypes.BurnStatus_BURN_STATUS_BURNED,
			}
			createdID, createErr := k.zenBTCKeeper.CreateBurnEvent(ctx, &newBurn)
			if createErr != nil {
				k.Logger(ctx).Error("StoreNewZenBTCBurnEvents: Error creating burn event in store.", "source", source, "tx_id", burn.TxID, "error", createErr)
				continue // Process next event
			}
			k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: Successfully created new burn event in store.", "source", source, "new_burn_id", createdID, "tx_id", burn.TxID, "log_idx", burn.LogIndex)
			foundNewBurn = true
		} else {
			k.Logger(ctx).Debug("StoreNewZenBTCBurnEvents: Skipping pre-existing event.", "source", source, "tx_id", burn.TxID, "log_idx", burn.LogIndex)
		}
	}

	// If a new burn event is found, we need to request the unstaker's Ethereum nonce
	// because the unstaking transaction happens on Ethereum, regardless of the burn source.
	if foundNewBurn {
		unstakerKeyID := k.zenBTCKeeper.GetUnstakerKeyID(ctx)
		k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: New burn events found. Setting EthereumNonceRequested for unstaker.", "source", source, "unstaker_key_id", unstakerKeyID)
		if err := k.EthereumNonceRequested.Set(ctx, unstakerKeyID, true); err != nil {
			k.Logger(ctx).Warn(fmt.Sprintf("storeNewZenBTCBurnEvents: %s", nonceErrorMsg), "source", source, "error", err)
		}
	} else {
		k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: No new burn events found to store.", "source", source)
	}
}

// processZenBTCBurnEvents processes pending burn events by constructing unstake transactions.
func (k *Keeper) processZenBTCBurnEvents(ctx sdk.Context, oracleData OracleData) {
	k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Started.")
	processTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetUnstakerKeyID(ctx),
		&oracleData.RequestedUnstakerNonce,
		nil,
		func(ctx sdk.Context) ([]zenbtctypes.BurnEvent, error) {
			pendingBurns, err := k.getPendingBurnEvents(ctx)
			if err != nil {
				k.Logger(ctx).Error("ProcessZenBTCBurnEvents: Error in getPendingBurnEvents.", "error", err)
				return nil, err
			}
			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Pending burn events fetched.", "count", len(pendingBurns))
			if len(pendingBurns) > 0 {
				for i, pb := range pendingBurns {
					k.Logger(ctx).Debug("ProcessZenBTCBurnEvents: Pending burn event details.",
						"idx", i, "burn_id", pb.Id, "tx_id", pb.TxID, "log_idx", pb.LogIndex, "chain_id", pb.ChainID,
						"amount", pb.Amount, "destination_addr_hex", hex.EncodeToString(pb.DestinationAddr), "status", pb.Status)
				}
			}
			return pendingBurns, nil
		},
		// Dispatch unstake transaction
		func(be zenbtctypes.BurnEvent) error {
			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Dispatching unstake for burn event.",
				"burn_id", be.Id, "origin_tx_id", be.TxID, "origin_chain_id", be.ChainID,
				"amount", be.Amount, "destination_addr_hex", hex.EncodeToString(be.DestinationAddr))

			if err := k.zenBTCKeeper.SetFirstPendingBurnEvent(ctx, be.Id); err != nil {
				k.Logger(ctx).Error("ProcessZenBTCBurnEvents: Failed to set first pending burn event.", "burn_id", be.Id, "error", err)
				return err // Return error to potentially halt or indicate failure
			}

			// Check for consensus
			requiredFields := []VoteExtensionField{VEFieldRequestedUnstakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice}
			consensusCheckDetails := fmt.Sprintf("burn_id: %d, origin_chain: %s, destination_addr_hex: %s, amount: %d", be.Id, be.ChainID, hex.EncodeToString(be.DestinationAddr), be.Amount)
			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Validating consensus for unstake.", "burn_id", be.Id, "required_fields", fmt.Sprintf("%v", requiredFields), "details_for_log", consensusCheckDetails)
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields, "zenBTC burn unstake", consensusCheckDetails); err != nil {
				k.Logger(ctx).Error("ProcessZenBTCBurnEvents: Consensus validation failed for unstake.", "burn_id", be.Id, "error", err)
				// If consensus fails, we don't proceed with this burn event in this block.
				// It will be retried in the next block if it's still the first pending.
				// To ensure it *is* retried and doesn't block others if it's a persistent consensus issue for *this* event,
				// we might need a mechanism to advance FirstPendingBurnEvent or mark this event as unprocessable.
				// For now, returning nil to let processTransaction handle it as a non-dispatch,
				// but this could lead to a stuck event if consensus is permanently missing for its required fields.
				// A more robust solution might involve a temporary "unprocessable" status.
				return nil // Returning nil, as per current validateConsensusForTxFields behavior (logs error, returns nil if missing consensus)
			}
			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Consensus validated for unstake.", "burn_id", be.Id)

			// Ensure DestinationAddr is not empty, as it's critical for the unstake transaction
			if len(be.DestinationAddr) == 0 {
				k.Logger(ctx).Error("ProcessZenBTCBurnEvents: Burn event has empty DestinationAddr. Cannot construct unstake tx.", "burn_id", be.Id)
				// This is a critical data issue. This burn event cannot be processed.
				// Consider moving it to a failed/error status and advancing the FirstPendingBurnEvent.
				// For now, returning an error to signify failure for this specific event.
				return fmt.Errorf("burn event %d has empty DestinationAddr", be.Id)
			}

			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Constructing unstake Ethereum transaction.", "burn_id", be.Id, "destination_addr_hex", hex.EncodeToString(be.DestinationAddr), "amount", be.Amount, "unstaker_nonce", oracleData.RequestedUnstakerNonce)
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
				k.Logger(ctx).Error("ProcessZenBTCBurnEvents: Failed to construct unstake Ethereum transaction.", "burn_id", be.Id, "error", err)
				return err
			}

			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				k.Logger(ctx).Error("ProcessZenBTCBurnEvents: Failed to get creator address for unstake tx.", "burn_id", be.Id, "unstaker_key_id", k.zenBTCKeeper.GetUnstakerKeyID(ctx), "error", err)
				return err
			}
			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Creator address for unstake tx.", "burn_id", be.Id, "creator", creator)

			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Submitting Ethereum transaction for unstake.",
				"burn_id", be.Id, "creator", creator, "unstaker_key_id", k.zenBTCKeeper.GetUnstakerKeyID(ctx),
				"eigen_chain_id", getChainIDForEigen(ctx))

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
		// Successfully processed unstake transaction
		func(be zenbtctypes.BurnEvent) error {
			k.Logger(ctx).Info("ProcessZenBTCBurnEvents: Nonce advanced for unstake. Updating burn event status.", "burn_id", be.Id, "old_status", be.Status, "new_status", zenbtctypes.BurnStatus_BURN_STATUS_UNSTAKING)
			be.Status = zenbtctypes.BurnStatus_BURN_STATUS_UNSTAKING
			return k.zenBTCKeeper.SetBurnEvent(ctx, be.Id, be)
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

		btcAmount := sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Amount)).Mul(exchangeRate).TruncateInt64()
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
	processTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetCompleterKeyID(ctx),
		&oracleData.RequestedCompleterNonce,
		nil,
		func(ctx sdk.Context) ([]zenbtctypes.Redemption, error) {
			firstPendingID, err := k.zenBTCKeeper.GetFirstPendingRedemption(ctx)
			if err != nil {
				firstPendingID = 0
			}
			return k.getRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_INITIATED, 2, firstPendingID)
		},
		// Dispatch unstake completer transaction
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
		// Successfully processed redemption, set to unstaked.
		func(r zenbtctypes.Redemption) error {
			r.Status = zenbtctypes.RedemptionStatus_UNSTAKED
			if err := k.zenBTCKeeper.SetRedemption(ctx, r.Data.Id, r); err != nil {
				return err
			}
			return k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), true)
		},
	)
}

func (k *Keeper) checkForRedemptionFulfilment(ctx sdk.Context) {
	startingIndex, err := k.zenBTCKeeper.GetFirstRedemptionAwaitingSign(ctx)
	if err != nil {
		startingIndex = 0
	}

	redemptions, err := k.getRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_AWAITING_SIGN, 0, startingIndex)
	if err != nil {
		k.Logger(ctx).Error("error getting redemptions", "error", err)
		return
	}

	if len(redemptions) == 0 {
		return
	}

	if err := k.zenBTCKeeper.SetFirstRedemptionAwaitingSign(ctx, redemptions[0].Data.Id); err != nil {
		k.Logger(ctx).Error("error setting first redemption awaiting sign", "error", err)
	}

	supply, err := k.zenBTCKeeper.GetSupply(ctx)
	if err != nil {
		k.Logger(ctx).Error("error getting zenBTC supply", "error", err)
		return
	}

	for _, redemption := range redemptions {
		signReq, err := k.treasuryKeeper.SignRequestStore.Get(ctx, redemption.Data.SignReqId)
		if err != nil {
			k.Logger(ctx).Error("error getting sign request for redemption", "id", redemption.Data.Id, "error", err)
			continue
		}

		if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING {
			continue
		}
		if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED {
			// Get current exchange rate
			exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
			if err != nil {
				k.Logger(ctx).Error("error getting zenBTC exchange rate", "error", err)
				continue
			}

			// redemption.Data.Amount is in zenBTC (what user wants to redeem)
			// Calculate how much BTC they should receive based on current exchange rate
			btcToRelease := uint64(sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Data.Amount)).Quo(exchangeRate).TruncateInt64())

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

		if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Data.Id, redemption); err != nil {
			k.Logger(ctx).Error("error updating redemption status", "error", err)
		}
	}

	if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
		k.Logger(ctx).Error("error updating zenBTC supply", "error", err)
	}

}

func (k Keeper) processSolanaROCKBurnEvents(ctx sdk.Context, oracleData OracleData) {
	var toProcess []*sidecarapitypes.BurnEvent
	for _, e := range oracleData.SolanaBurnEvents {

		addr, err := sdk.Bech32ifyAddressBytes("zen", e.DestinationAddr[:20])
		if err != nil {
			k.Logger(ctx).Error(fmt.Errorf("Bech32ifyAddressBytes: %w", err).Error())
			continue
		}
		burns, err := k.zentpKeeper.GetBurns(ctx, addr, e.ChainID, e.TxID)
		if err != nil {
			k.Logger(ctx).Error(err.Error())
			continue
		}
		if len(burns) > 0 {
			continue // burn already processed
		} else {
			toProcess = append(toProcess, &e)
		}
	}

	// TODO do cleanup on error. e.g. burn minted funds if there is an error sendig them to the recipient, or adding of the bridge fails
	for _, burn := range toProcess {
		coins := sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(burn.Amount)))
		if err := k.bankKeeper.MintCoins(ctx, zentptypes.ModuleName, coins); err != nil {
			k.Logger(ctx).Error(fmt.Errorf("MintCoins: %w", err).Error())
			continue
		}
		addr, err := sdk.Bech32ifyAddressBytes("zen", burn.DestinationAddr[:20])
		if err != nil {
			k.Logger(ctx).Error(fmt.Errorf("Bech32ifyAddressBytes: %w", err).Error())
			continue
		}
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			k.Logger(ctx).Error(fmt.Errorf("AccAddressFromBech32: %w", err).Error())
			continue
		}

		protocolWalletAddress, bridgeFee, err := k.zentpKeeper.GetBridgeFeeParams(ctx)
		if err != nil {
			k.Logger(ctx).Error("GetBridgeFeeParams: ", err.Error())
			return
		}

		bridgeFeeCoins, err := k.zentpKeeper.GetBridgeFeeAmount(ctx, burn.Amount, bridgeFee)
		if err != nil {
			k.Logger(ctx).Error("GetBridgeFeeAmount: ", err.Error())
			return
		}

		bridgeAmount := sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(burn.Amount).Sub(bridgeFeeCoins.AmountOf(params.BondDenom))))

		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zentptypes.ModuleName, accAddr, bridgeAmount); err != nil {
			k.Logger(ctx).Error(fmt.Errorf("SendCoinsFromModuleToAccount: %w", err).Error())
		}

		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zentptypes.ModuleName, protocolWalletAddress, bridgeFeeCoins); err != nil {
			k.Logger(ctx).Error(fmt.Errorf("SendCoinsFromModuleToAccount: %w", err).Error())
		}
		err = k.zentpKeeper.AddBurn(ctx, &zentptypes.Bridge{
			Denom:            params.BondDenom,
			Amount:           burn.Amount,
			RecipientAddress: accAddr.String(),
			SourceChain:      burn.ChainID,
			TxHash:           burn.TxID,
			State:            zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED,
			BlockHeight:      ctx.BlockHeight(),
		})
		if err != nil {
			k.Logger(ctx).Error(err.Error())
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeValidation,
				sdk.NewAttribute(types.AttributeKeyBridgeAmount, bridgeAmount.String()),
				sdk.NewAttribute(types.AttributeKeyBridgeFee, bridgeFeeCoins.String()),
				sdk.NewAttribute(types.AttributeKeyBurnDestination, addr),
			),
		)
	}
}
