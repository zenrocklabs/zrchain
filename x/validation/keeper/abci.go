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
	solana "github.com/gagliardetto/solana-go"
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
	oracleData, err := k.gatherOracleDataForVoteExtension(ctx, req.Height)
	if err != nil {
		k.Logger(ctx).Error("error gathering oracle data for vote extension", "height", req.Height, "error", err)
		return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
	}

	voteExt, err := ConstructVoteExtension(oracleData)
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

// gatherOracleDataForVoteExtension fetches all necessary on-chain and sidecar data to build a vote extension.
func (k *Keeper) gatherOracleDataForVoteExtension(ctx context.Context, height int64) (*OracleData, error) {
	oracleData, err := k.GetSidecarState(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("error retrieving sidecar state: %w", err)
	}

	latestHeader, requestedHeader, err := k.retrieveBitcoinHeaders(ctx)
	if err != nil {
		return nil, err
	}
	if latestHeader != nil {
		oracleData.LatestBtcBlockHeight = latestHeader.BlockHeight
		if latestHeader.BlockHeader != nil {
			oracleData.LatestBtcBlockHeader = *latestHeader.BlockHeader
		}
	}
	if requestedHeader != nil {
		oracleData.RequestedBtcBlockHeight = requestedHeader.BlockHeight
		if requestedHeader.BlockHeader != nil {
			oracleData.RequestedBtcBlockHeader = *requestedHeader.BlockHeader
		}
	}

	nonces := make(map[uint64]uint64)
	for _, key := range k.getZenBTCKeyIDs(ctx) {
		requested, err := k.EthereumNonceRequested.Get(ctx, key)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				return nil, err
			}
			requested = false
		}
		if requested {
			nonce, err := k.lookupEthereumNonce(ctx, key)
			if err != nil {
				return nil, err
			}
			nonces[key] = nonce
		}
	}
	oracleData.RequestedStakerNonce = nonces[k.zenBTCKeeper.GetStakerKeyID(ctx)]
	oracleData.RequestedEthMinterNonce = nonces[k.zenBTCKeeper.GetEthMinterKeyID(ctx)]
	oracleData.RequestedUnstakerNonce = nonces[k.zenBTCKeeper.GetUnstakerKeyID(ctx)]
	oracleData.RequestedCompleterNonce = nonces[k.zenBTCKeeper.GetCompleterKeyID(ctx)]

	solNonce, err := k.retrieveSolanaNonces(ctx)
	if err != nil {
		return nil, err
	}
	oracleData.SolanaMintNonces = solNonce

	solAccs, err := k.retrieveSolanaAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving solana accounts: %w", err)
	}
	oracleData.SolanaAccounts = solAccs

	return oracleData, nil
}

// ConstructVoteExtension builds the vote extension based on oracle data.
func ConstructVoteExtension(oracleData *OracleData) (VoteExtension, error) {
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

	latestBitcoinHeaderHash, err := deriveHash(oracleData.LatestBtcBlockHeader)
	if err != nil {
		return VoteExtension{}, err
	}

	// Only set requested header fields if there's a requested header
	var requestedBtcHeaderHash []byte
	if oracleData.RequestedBtcBlockHeight > 0 {
		requestedBitcoinHeaderHash, err := deriveHash(oracleData.RequestedBtcBlockHeader)
		if err != nil {
			return VoteExtension{}, err
		}
		requestedBtcHeaderHash = requestedBitcoinHeaderHash[:]
	}

	solNonceHash, err := deriveHash(oracleData.SolanaMintNonces)
	if err != nil {
		return VoteExtension{}, err
	}

	solAccsHash, err := deriveHash(oracleData.SolanaAccounts)
	if err != nil {
		return VoteExtension{}, err
	}

	solanaMintEventsHash, err := deriveHash(oracleData.SolanaMintEvents)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving solana mint events hash: %w", err)
	}
	solanaBurnEventsHash, err := deriveHash(oracleData.SolanaBurnEvents)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving solana burn events hash: %w", err)
	}

	voteExt := VoteExtension{
		ROCKUSDPrice:            oracleData.ROCKUSDPrice,
		BTCUSDPrice:             oracleData.BTCUSDPrice,
		ETHUSDPrice:             oracleData.ETHUSDPrice,
		EigenDelegationsHash:    avsDelegationsHash[:],
		EthBurnEventsHash:       ethBurnEventsHash[:],
		RedemptionsHash:         redemptionsHash[:],
		RequestedBtcBlockHeight: oracleData.RequestedBtcBlockHeight,
		RequestedBtcHeaderHash:  requestedBtcHeaderHash,
		LatestBtcBlockHeight:    oracleData.LatestBtcBlockHeight,
		LatestBtcHeaderHash:     latestBitcoinHeaderHash[:],
		EthBlockHeight:          oracleData.EthBlockHeight,
		EthGasLimit:             oracleData.EthGasLimit,
		EthBaseFee:              oracleData.EthBaseFee,
		EthTipCap:               oracleData.EthTipCap,
		RequestedStakerNonce:    oracleData.RequestedStakerNonce,
		RequestedEthMinterNonce: oracleData.RequestedEthMinterNonce,
		RequestedUnstakerNonce:  oracleData.RequestedUnstakerNonce,
		RequestedCompleterNonce: oracleData.RequestedCompleterNonce,
		SolanaMintNoncesHash:    solNonceHash[:],
		SolanaAccountsHash:      solAccsHash[:],
		SolanaMintEventsHash:    solanaMintEventsHash[:],
		SolanaBurnEventsHash:    solanaBurnEventsHash[:],
		SidecarVersionName:      oracleData.SidecarVersionName,
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
		k.Logger(ctx).Error("error unmarshalling vote extension", "height", req.Height, "error", err)
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

	consensusVE, _, fieldVotePowers, err := k.GetConsensusAndPluralityVEData(ctx, req.Height, req.LocalLastCommit)
	if err != nil {
		k.Logger(ctx).Error("error retrieving supermajority vote extension data", "height", req.Height, "error", err)
		return nil, nil
	}

	if len(fieldVotePowers) == 0 { // no field reached consensus
		k.Logger(ctx).Warn("no fields reached consensus in vote extension", "height", req.Height)
		return k.marshalOracleData(req, &OracleData{ConsensusData: req.LocalLastCommit, FieldVotePowers: fieldVotePowers})
	}

	oracleData, err := k.GetValidatedOracleData(ctx, consensusVE, fieldVotePowers)
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

	_, pluralityVE, ok := k.validateCanonicalVE(ctx, req.Height, oracleData)
	if !ok {
		k.Logger(ctx).Error("invalid canonical vote extension")
		return nil
	}

	// Update asset prices if there's consensus on the price fields
	k.updateAssetPrices(ctx, oracleData)

	// Validator updates - only if EigenDelegationsHash has consensus
	if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldEigenDelegationsHash) {
		k.UpdateValidatorStakes(ctx, oracleData)
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

		// 1. Update on-chain state based on oracle data that has reached consensus
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

		// Skipping EigenLayer-based redemptions ingestion; redemptions will be initiated directly on burn
		// if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldRedemptionsHash) {
		// 	k.storeNewZenBTCRedemptions(ctx, oracleData)
		// }

		if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
			k.processSolanaZenBTCMintEvents(ctx, oracleData)
			k.processSolanaROCKMintEvents(ctx, oracleData)
		}

		if fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaBurnEventsHash) {
			k.storeNewZenBTCBurnEventsSolana(ctx, oracleData)
			k.processSolanaROCKBurnEvents(ctx, oracleData)
		}

		// 2. Process pending transaction queues based on the latest state
		// Request nonces/accounts for direct minting (no EigenLayer staking)
		if pendingEVM, err := k.getPendingMintTransactions(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED, zenbtctypes.WalletType_WALLET_TYPE_EVM); err == nil && len(pendingEVM) > 0 {
			_ = k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), true)
		}
		if pendingSol, err := k.getPendingMintTransactions(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED, zenbtctypes.WalletType_WALLET_TYPE_SOLANA); err == nil && len(pendingSol) > 0 {
			solParams := k.zenBTCKeeper.GetSolanaParams(ctx)
			_ = k.SolanaNonceRequested.Set(ctx, solParams.NonceAccountKey, true)
			for _, tx := range pendingSol {
				_ = k.SetSolanaZenBTCRequestedAccount(ctx, tx.RecipientAddress, true)
			}
		}

		// Mint directly on destination chains
		k.processZenBTCMintsEthereum(ctx, oracleData)
		k.processZenBTCMintsSolana(ctx, oracleData)

		// Skip EigenLayer unstake and completion; proceed directly to BTC redemption monitoring
		k.checkForRedemptionFulfilment(ctx)
		k.processSolanaROCKMints(ctx, oracleData)

		// 3. Final cleanup steps for the block
		k.clearSolanaAccounts(ctx)
	}

	k.recordNonVotingValidators(ctx, req)
	k.recordMismatchedVoteExtensions(ctx, req.Height, pluralityVE, oracleData.ConsensusData)

	// Perform final invariant checks for the block.
	if err := k.zentpKeeper.CheckROCKSupplyCap(ctx, sdkmath.ZeroInt()); err != nil {
		// This is a critical failure. In a real-world scenario, this should halt the chain.
		k.Logger(ctx).Error("CRITICAL INVARIANT VIOLATION: ROCK supply cap check failed at end of block.", "error", err.Error())
		// For now, we will log a critical error. The chain will continue, but this indicates a serious issue.
	}

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
// Returns both the consensus vote extension and the plurality vote extension.
func (k *Keeper) validateCanonicalVE(ctx sdk.Context, height int64, oracleData OracleData) (VoteExtension, VoteExtension, bool) {
	consensusVE, pluralityVE, fieldVotePowers, err := k.GetConsensusAndPluralityVEData(ctx, height, oracleData.ConsensusData)
	if err != nil {
		k.Logger(ctx).Error("error getting super majority VE data", "height", height, "error", err)
		return VoteExtension{}, VoteExtension{}, false
	}

	if reflect.DeepEqual(consensusVE, VoteExtension{}) {
		k.Logger(ctx).Warn("accepting empty vote extension", "height", height)
		return VoteExtension{}, pluralityVE, true
	}

	k.validateOracleData(ctx, consensusVE, &oracleData, fieldVotePowers)

	// Log final consensus summary after validation
	k.Logger(ctx).Info("final consensus summary",
		"fields_with_consensus", len(oracleData.FieldVotePowers),
		"stage", "post_validation")

	return consensusVE, pluralityVE, true
}

// getValidatedOracleData retrieves and validates oracle data based on a vote extension.
// Only validates fields that have reached consensus as indicated in fieldVotePowers.
func (k *Keeper) GetValidatedOracleData(ctx sdk.Context, voteExt VoteExtension, fieldVotePowers map[VoteExtensionField]int64) (*OracleData, error) {
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
		oracleData.SolanaMintNonces, err = k.retrieveSolanaNonces(ctx)
		if err != nil {
			return nil, fmt.Errorf("error collecting solana nonces: %w", err)
		}

	}

	if fieldHasConsensus(fieldVotePowers, VEFieldSolanaAccountsHash) {
		solAccs, err := k.retrieveSolanaAccounts(ctx)
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
func (k *Keeper) UpdateValidatorStakes(ctx sdk.Context, oracleData OracleData) {
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

	k.RemoveStaleValidatorDelegations(ctx, validatorInAVSDelegationSet)
}

// removeStaleValidatorDelegations removes delegation entries for validators not present in the current AVS data.
func (k *Keeper) RemoveStaleValidatorDelegations(ctx sdk.Context, validatorInAVSDelegationSet map[string]bool) {
	var validatorsToRemove []string

	if err := k.ValidatorDelegations.Walk(ctx, nil, func(valAddr string, stake sdkmath.Int) (bool, error) {
		if !validatorInAVSDelegationSet[valAddr] {
			validatorsToRemove = append(validatorsToRemove, valAddr)
		}
		return false, nil
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

// storeBitcoinBlockHeaders stores the Bitcoin header and handles historical header requests.
func (k *Keeper) storeBitcoinBlockHeaders(ctx sdk.Context, oracleData OracleData) error {
	// Get requested headers and latest stored height early
	requestedHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting requested historical Bitcoin headers", "error", err)
		return err
	}

	latestBtcHeaderHeight, err := k.LatestBtcHeaderHeight.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error getting latest BTC header height", "error", err)
		// Do not return, as we can still proceed, but log the error.
	}

	// Process the latest and requested headers. The latest header is the current Bitcoin tip,
	// while requested headers are typically historical blocks needed for reorg detection or gap filling.
	// Most blocks have no requested headers, and when they do exist, they're usually older than the latest.
	if newHeight, updated := k.processAndStoreBtcHeader(ctx, oracleData.LatestBtcBlockHeight, &oracleData.LatestBtcBlockHeader, latestBtcHeaderHeight, &requestedHeaders, "latest"); updated {
		latestBtcHeaderHeight = newHeight
	}

	// Process requested header only if it's different from the latest one to avoid redundant processing
	if oracleData.RequestedBtcBlockHeight != oracleData.LatestBtcBlockHeight {
		k.processAndStoreBtcHeader(ctx, oracleData.RequestedBtcBlockHeight, &oracleData.RequestedBtcBlockHeader, latestBtcHeaderHeight, &requestedHeaders, "requested")
	}

	// Clean up the list of requested headers by removing any that have now been stored.
	if len(requestedHeaders.Heights) > 0 {
		requestedHeaders.Heights = slices.DeleteFunc(requestedHeaders.Heights, func(height int64) bool {
			has, _ := k.BtcBlockHeaders.Has(ctx, height)
			return has
		})
	}

	// Persist the updated list of requested headers
	if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedHeaders); err != nil {
		k.Logger(ctx).Error("error updating requested historical Bitcoin headers", "error", err)
	}

	return nil
}

// processAndStoreBtcHeader checks if a given Bitcoin header is new, stores it,
// and triggers a reorg check if it's a new high-water mark.
// It returns the updated latestBtcHeaderHeight and a boolean indicating if it was updated.
func (k *Keeper) processAndStoreBtcHeader(
	ctx sdk.Context,
	headerHeight int64,
	header *sidecarapitypes.BTCBlockHeader,
	latestBtcHeaderHeight int64,
	requestedHeaders *zenbtctypes.RequestedBitcoinHeaders,
	headerType string,
) (int64, bool) {
	if headerHeight <= 0 || header == nil || header.MerkleRoot == "" {
		return latestBtcHeaderHeight, false
	}

	// Check if header already exists by comparing hashes
	existingHeader, err := k.BtcBlockHeaders.Get(ctx, headerHeight)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error checking if bitcoin header exists", "type", headerType, "height", headerHeight, "error", err)
		return latestBtcHeaderHeight, false
	}

	// If header exists, compare hashes to see if it's different
	if err == nil {
		existingHash, err := deriveHash(existingHeader)
		if err != nil {
			k.Logger(ctx).Error("error deriving hash for existing header", "type", headerType, "height", headerHeight, "error", err)
			return latestBtcHeaderHeight, false
		}

		newHash, err := deriveHash(*header)
		if err != nil {
			k.Logger(ctx).Error("error deriving hash for new header", "type", headerType, "height", headerHeight, "error", err)
			return latestBtcHeaderHeight, false
		}

		if bytes.Equal(existingHash[:], newHash[:]) {
			return latestBtcHeaderHeight, false
		}
	}

	// Store the new header (either no header existed or hash is different)
	if err := k.BtcBlockHeaders.Set(ctx, headerHeight, *header); err != nil {
		k.Logger(ctx).Error("error storing bitcoin header", "type", headerType, "height", headerHeight, "error", err)
		return latestBtcHeaderHeight, false
	}

	k.Logger(ctx).Info("stored new bitcoin header", "type", headerType, "height", headerHeight)

	// Always perform reorg/gap check when a new header is stored
	k.checkForBitcoinReorg(ctx, headerHeight, latestBtcHeaderHeight, requestedHeaders)

	// Update the latest height if this header is newer than what we had before
	if headerHeight > latestBtcHeaderHeight {
		if err := k.LatestBtcHeaderHeight.Set(ctx, headerHeight); err != nil {
			k.Logger(ctx).Error("error setting latest BTC header height", "error", err)
		}
		return headerHeight, true
	}

	// Header was stored but it's not newer than our current latest height
	return latestBtcHeaderHeight, false
}

// checkForBitcoinReorg checks for gaps and requests previous headers for reorg detection.
// This function does NOT modify the LatestBtcHeaderHeight state.
func (k *Keeper) checkForBitcoinReorg(ctx sdk.Context, newHeaderHeight, latestStoredHeight int64, requestedHeaders *zenbtctypes.RequestedBitcoinHeaders) {
	var numHistoricalHeadersToRequest int64 = 20
	if strings.HasPrefix(ctx.ChainID(), "diamond") {
		numHistoricalHeadersToRequest = 6
	}

	// Check for gaps between the latest stored header and the new header.
	// Only run the check if we have a previously stored height (i.e., it's not the first run).
	if latestStoredHeight > 0 {
		for i := latestStoredHeight + 1; i < newHeaderHeight; i++ {
			requestedHeaders.Heights = append(requestedHeaders.Heights, i)
		}
	}

	// Request N previous headers from the new tip for reorg validation.
	prevHeights := make([]int64, 0, numHistoricalHeadersToRequest)
	for i := int64(1); i <= numHistoricalHeadersToRequest; i++ {
		prevHeight := newHeaderHeight - i
		if prevHeight <= 0 {
			break
		}
		prevHeights = append(prevHeights, prevHeight)
	}

	requestedHeaders.Heights = append(requestedHeaders.Heights, prevHeights...)
	k.Logger(ctx).Info("requested headers after reorg check", "new_requests", requestedHeaders.Heights)
}

//
// =============================================================================
// ZENBTC PROCESSING: STAKING, MINTING, BURN EVENTS & REDEMPTIONS
// =============================================================================
//

// checkForUpdateAndDispatchTx processes nonce updates and transaction dispatch. It contains separate logic
// for Ethereum and Solana based transactions due to their different nonce mechanisms and transaction lifecycles.
func checkForUpdateAndDispatchTx[T any](
	k *Keeper,
	ctx sdk.Context,
	keyID uint64,
	requestedEthNonce *uint64,
	requestedSolNonce *solSystem.NonceAccount,
	nonceReqStore collections.Map[uint64, bool],
	pendingTxs []T,
	txDispatchCallback func(tx T) error,
	txContinuationCallback func(tx T) error,
) {
	if len(pendingTxs) == 0 {
		return
	}

	// Ethereum transaction processing flow.
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

		nonceUpdated, err := handleNonceUpdate(k, ctx, keyID, *requestedEthNonce, nonceData, pendingTxs[0], txContinuationCallback)
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

		// If tx[0] confirmed on-chain via nonce increment, dispatch tx[1]. If not then retry dispatching tx[0].
		txIndex := 0
		if nonceUpdated {
			txIndex = 1
		}

		if len(pendingTxs) <= txIndex {
			return
		}

		if err := txDispatchCallback(pendingTxs[txIndex]); err != nil {
			k.Logger(ctx).Error("tx dispatch callback error", "keyID", keyID, "error", err)
		}
		return
	}

	// Solana transaction processing flow.
	if requestedSolNonce != nil {
		k.Logger(ctx).Info("processing solana transaction with nonce", "nonce", requestedSolNonce.Nonce)

		if requestedSolNonce.Nonce.IsZero() {
			k.Logger(ctx).Error("solana nonce is zero")
			return
		}

		// For Solana, `txContinuationCallback` is a misnomer. It's a status/timeout checker for the head of the queue.
		// We call it, and then attempt to dispatch the same transaction. The dispatch is idempotent.
		if err := txContinuationCallback(pendingTxs[0]); err != nil {
			k.Logger(ctx).Error("error handling solana transaction status check", "keyID", keyID, "error", err)
			return
		}

		// If tx[0] is still pending, dispatch it. The dispatch callback is idempotent.
		if err := txDispatchCallback(pendingTxs[0]); err != nil {
			k.Logger(ctx).Error("tx dispatch callback error", "keyID", keyID, "error", err)
		}
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
	txContinuationCallback func(tx T) error,
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
	checkForUpdateAndDispatchTx(k, ctx, keyID, requestedEthNonce, requestedSolNonce, nonceReqStore, pendingTxs, txDispatchCallback, txContinuationCallback)
}

// getPendingTransactions is a generic helper that walks a store with key type uint64
// and returns a slice of items of type T that satisfy the provided predicate, up to a given limit.
// If limit is 0, all matching items will be returned.
func getPendingTransactions[T any](ctx sdk.Context, store interface {
	Walk(ctx sdk.Context, rng *collections.Range[uint64], fn func(key uint64, value T) (bool, error)) error
}, predicate func(T) bool, firstPendingID uint64, limit int) ([]T, error) {
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
	txContinuationCallback func(tx T) error,
) (bool, error) {
	if requestedNonce != nonceData.PrevNonce {
		if err := txContinuationCallback(tx); err != nil {
			return false, fmt.Errorf("tx continuation callback error: %w", err)
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
		// pendingGetter: Fetches pending zenBTC mints that are in the DEPOSITED state.
		// These are transactions that have received a BTC deposit and are ready to be staked on EigenLayer.
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			return k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				zenbtctypes.WalletType_WALLET_TYPE_UNSPECIFIED,
			)
		},
		// txDispatchCallback: Constructs and submits an Ethereum transaction to stake the deposited assets on EigenLayer.
		// It uses the current nonce and gas parameters from the oracle data.
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
				treasurytypes.WalletType_WALLET_TYPE_EVM,
				getChainIDForEigen(ctx),
				unsignedTx,
				unsignedTxHash,
			)
		},
		// txContinuationCallback: This is called when the stake transaction's nonce has been confirmed on-chain.
		// It updates the transaction's status to STAKED and sets up the system for the next step: minting zenBTC.
		func(tx zenbtctypes.PendingMintTransaction) error {
			tx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED
			if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx); err != nil {
				return err
			}
			if types.IsSolanaCAIP2(ctx, tx.Caip2ChainId) {
				solParams := k.zenBTCKeeper.GetSolanaParams(ctx)
				if err := k.SolanaNonceRequested.Set(ctx, solParams.NonceAccountKey, true); err != nil {
					return err
				}
				if err := k.SetSolanaZenBTCRequestedAccount(ctx, tx.RecipientAddress, true); err != nil {
					return err
				}
				k.Logger(ctx).Warn("processed zenbtc stake", "tx_id", tx.Id, "recipient", tx.RecipientAddress, "amount", tx.Amount)
				return nil
			} else if types.IsEthereumCAIP2(ctx, tx.Caip2ChainId) {
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
		// pendingGetter: Fetches pending zenBTC mints for EVM chains that are in the DEPOSITED state.
		// We skip EigenLayer staking and mint directly on the destination chain.
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			return k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				zenbtctypes.WalletType_WALLET_TYPE_EVM,
			)
		},
		// txDispatchCallback: Constructs and submits an Ethereum transaction to mint zenBTC on the destination EVM chain.
		// It calculates the required fees and constructs the minting transaction.
		func(tx zenbtctypes.PendingMintTransaction) error {
			if err := k.zenBTCKeeper.SetFirstPendingEthMintTransaction(ctx, tx.Id); err != nil {
				return err
			}

			// Check for consensus
			requiredFields := []VoteExtensionField{VEFieldRequestedEthMinterNonce, VEFieldBTCUSDPrice}
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

			feeZenBTC := k.CalculateFlatZenBTCMintFee(
				btcUSDPrice,
				exchangeRate,
			)
			feeZenBTC = min(feeZenBTC, tx.Amount)

			chainID, err := types.ValidateEVMChainID(ctx, tx.Caip2ChainId)
			if err != nil {
				return fmt.Errorf("unsupported chain ID: %w", err)
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
				chainID,
				unsignedMintTx,
				unsignedMintTxHash,
			)
		},
		// txContinuationCallback: This is called when the mint transaction's nonce has been confirmed on-chain.
		// It updates the zenBTC supply and marks the mint transaction as MINTED.
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
			// Update local bedrock accounting for default validator (BTC in sats)
// NOTE: For TokensBedrock accounting we do not use exchangeRate.
			// We record the event amount directly for v1 (units must correspond to BTC-sats in future update).
			if err := k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(tx.Amount)); err != nil {
				k.Logger(ctx).Error("error adjusting bedrock BTC on mint", "error", err)
			}
			return k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
		},
	)
}

// adjustDefaultValidatorBedrockBTC adds (positive) or subtracts (negative via Int sign) BTC sats to the default validator's TokensBedrock (Asset_BTC)
func (k *Keeper) adjustDefaultValidatorBedrockBTC(ctx sdk.Context, delta sdkmath.Int) error {
	oper := k.GetBedrockDefaultValOperAddr(ctx)
	v, err := k.GetZenrockValidatorFromBech32(ctx, oper)
	if err != nil {
		return err
	}
	// Find existing BTC entry
	idx := -1
	for i, td := range v.TokensBedrock {
		if td != nil && td.Asset == types.Asset_BTC {
			idx = i
			break
		}
	}
	if idx >= 0 {
		newAmt := v.TokensBedrock[idx].Amount.Add(delta)
		if newAmt.IsNegative() {
			newAmt = sdkmath.ZeroInt()
		}
		v.TokensBedrock[idx].Amount = newAmt
	} else {
		amt := delta
		if amt.IsNegative() {
			amt = sdkmath.ZeroInt()
		}
		v.TokensBedrock = append(v.TokensBedrock, &types.TokenData{Asset: types.Asset_BTC, Amount: amt})
	}
	return k.SetValidator(ctx, v)
}

// processZenBTCMintsSolana processes pending mint transactions.
func (k *Keeper) processZenBTCMintsSolana(ctx sdk.Context, oracleData OracleData) {
	processTransaction(
		k,
		ctx,
		k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey,
		nil,
		oracleData.SolanaMintNonces[k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey],
		// pendingGetter: Fetches pending zenBTC mints for Solana that are in the DEPOSITED state.
		// We skip EigenLayer staking and mint directly on Solana.
		func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			pendingMints, err := k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				zenbtctypes.WalletType_WALLET_TYPE_SOLANA,
			)
			k.Logger(ctx).Warn("pending zenbtc solana mints", "mints", fmt.Sprintf("%v", pendingMints), "count", len(pendingMints))
			return pendingMints, err
		},
		// txDispatchCallback: Constructs and dispatches a Solana transaction to mint zenBTC.
		// This function is idempotent and will only send a new transaction if one is not already in-flight.
		func(tx zenbtctypes.PendingMintTransaction) error {
			k.Logger(ctx).Warn("dispatch handler triggered", "tx_id", tx.Id, "recipient", tx.RecipientAddress, "amount", tx.Amount)
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

			exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
			if err != nil {
				return err
			}

			feeZenBTC := k.CalculateFlatZenBTCMintFee(
				btcUSDPrice,
				exchangeRate,
			)
			feeZenBTC = min(feeZenBTC, tx.Amount)

			solParams := k.zenBTCKeeper.GetSolanaParams(ctx)

			// Derive the ATA for the recipient and check its state
			recipientPubKey, err := solana.PublicKeyFromBase58(tx.RecipientAddress)
			if err != nil {
				return fmt.Errorf("invalid recipient address %s for ZenBTC mint: %w", tx.RecipientAddress, err)
			}
			mintPubKey, err := solana.PublicKeyFromBase58(solParams.MintAddress)
			if err != nil {
				return fmt.Errorf("invalid ZenBTC mint address %s: %w", solParams.MintAddress, err)
			}
			expectedATA, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintPubKey)
			if err != nil {
				return fmt.Errorf("failed to derive ATA for recipient %s, mint %s: %w", tx.RecipientAddress, solParams.MintAddress, err)
			}

			fundReceiver := false
			ata, ok := oracleData.SolanaAccounts[expectedATA.String()]
			if !ok {
				// This means the ATA was not requested or not found by collectSolanaAccounts.
				// Depending on strictness, this could be an error or imply it needs funding.
				// For safety, if it's not in oracleData, we can assume it needs creation/funding.
				// However, collectSolanaAccounts should have fetched it if it was SetSolanaZenBTCRequestedAccount.
				k.Logger(ctx).Warn("ATA not found in oracleData.SolanaAccounts, assuming it might need funding", "ata", expectedATA.String(), "recipient", tx.RecipientAddress)
				fundReceiver = true // If not found in map, assume it needs to be created.
			} else if ata.State == solToken.Uninitialized {
				fundReceiver = true
			}

			nonce, ok := oracleData.SolanaMintNonces[solParams.NonceAccountKey]
			if !ok {
				return fmt.Errorf("nonce not found in oracleData.SolanaMintNonces for solParams.NonceAccountKey: %d", solParams.NonceAccountKey)
			}

			txPrepReq := &solanaMintTxRequest{
				amount:            tx.Amount,
				fee:               feeZenBTC, // TODO: currently we are not using solParams.Fee
				recipient:         tx.RecipientAddress,
				nonce:             nonce,
				fundReceiver:      fundReceiver,
				programID:         solParams.ProgramId,
				mintAddress:       solParams.MintAddress,
				feeWallet:         solParams.FeeWallet,
				nonceAccountKey:   solParams.NonceAccountKey,
				nonceAuthorityKey: solParams.NonceAuthorityKey,
				signerKey:         solParams.SignerKeyId,
				eventID:           tx.Id,
				zenbtc:            true,
			}
			k.Logger(ctx).Warn("processing zenbtc solana mint", "tx_id", tx.Id, "recipient", tx.RecipientAddress, "amount", tx.Amount, "fee", feeZenBTC)
			transaction, err := k.PrepareSolanaMintTx(ctx, txPrepReq)
			if err != nil {
				return fmt.Errorf("PrepareSolRockMintTx: %w", err)
			}

			k.Logger(ctx).Warn("processing zenBTC mint",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.SolanaMintNonces[solParams.NonceAccountKey],
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
			if err = k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx); err != nil {
				return err
			}
			solNonce := types.SolanaNonce{Nonce: oracleData.SolanaMintNonces[solParams.NonceAccountKey].Nonce[:]}
			return k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, solNonce)
		},
		// txContinuationCallback: For Solana, this function acts as a status and timeout checker for in-flight transactions.
		// It is called on every block to check if a pending transaction has been confirmed, has timed out (BTL),
		// or requires a retry. It manages the lifecycle of the pending Solana transaction.
		func(tx zenbtctypes.PendingMintTransaction) error {
			// If we don't have consensus on SolanaMintEventsHash we cannot reliably determine
			// whether the associated event has arrived on-chain. In that case we should *not* run
			// any retry or timeout logic, otherwise we risk redelivering the same transaction over
			// and over again without a reliable confirmation signal.

			if !fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
				k.Logger(ctx).Debug("Skipping Solana mint retry/timeout checks – no consensus on SolanaMintEventsHash", "tx_id", tx.Id)
				return nil
			}

			// If BlockHeight is 0, this transaction was either just dispatched in the current block
			// by the txDispatchCallback, or it has been reset for a full retry by prior logic in this callback.
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
		// pendingGetter: Fetches pending solROCK mints that are in the PENDING state.
		// These transactions have completed the EigenLayer staking step and are ready for zenBTC to be minted on the destination chain.
		func(ctx sdk.Context) ([]*zentptypes.Bridge, error) {
			mints, err := k.zentpKeeper.GetMintsWithStatusPending(ctx)
			return mints, err
		},
		// txDispatchCallback: Constructs and dispatches a Solana transaction to mint ROCK tokens.
		// This function is idempotent and handles the preparation of the Solana transaction.
		func(tx *zentptypes.Bridge) error {
			// Check whether this tx has already been processed, if it has been - we wait for it to complete (or timeout)
			if tx.BlockHeight > 0 {
				k.Logger(ctx).Info("waiting for pending zentp solana mint tx", "tx_id", tx.Id, "block_height", tx.BlockHeight)
				return nil
			}

			// Check for consensus - VEFieldSolanaMintNoncesHash is also needed for the nonce in PrepareSolanaMintTx
			requiredFields := []VoteExtensionField{VEFieldSolanaAccountsHash, VEFieldSolanaMintNoncesHash}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"solROCK mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return fmt.Errorf("validateConsensusForTxFields: %w", err)
			}

			solParams := k.zentpKeeper.GetSolanaParams(ctx)
			// Derive the ATA for the recipient and check its state
			recipientPubKey, err := solana.PublicKeyFromBase58(tx.RecipientAddress)
			if err != nil {
				return fmt.Errorf("invalid recipient address %s for ZenTP mint: %w", tx.RecipientAddress, err)
			}
			mintPubKey, err := solana.PublicKeyFromBase58(solParams.MintAddress)
			if err != nil {
				return fmt.Errorf("invalid ZenTP mint address %s: %w", solParams.MintAddress, err)
			}
			expectedATA, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintPubKey)
			if err != nil {
				return fmt.Errorf("failed to derive ATA for ZenTP recipient %s, mint %s: %w", tx.RecipientAddress, solParams.MintAddress, err)
			}

			fundReceiver := false
			ata, ok := oracleData.SolanaAccounts[expectedATA.String()]
			if !ok {
				// If the ATA is not in oracleData.SolanaAccounts, it means it wasn't requested via SetSolanaZenTPRequestedAccount
				// or collectSolanaAccounts failed to fetch it. This is a state mismatch if a transaction is being prepared for it.
				// For robustness, one might assume it needs funding, but it could also indicate an issue.
				k.Logger(ctx).Warn("ATA not found in oracleData.SolanaAccounts for ZenTP, tx will proceed assuming it needs funding or creation", "ata", expectedATA.String(), "recipient", tx.RecipientAddress)
				fundReceiver = true
			} else if ata.State == solToken.Uninitialized {
				fundReceiver = true
			}

			nonce, ok := oracleData.SolanaMintNonces[solParams.NonceAccountKey]
			if !ok {
				return fmt.Errorf("nonce not found in oracleData.SolanaMintNonces for solParams.NonceAccountKey: %d", solParams.NonceAccountKey)
			}

			transaction, err := k.PrepareSolanaMintTx(ctx, &solanaMintTxRequest{
				amount:            tx.Amount,
				fee:               min(solParams.Fee, tx.Amount),
				recipient:         tx.RecipientAddress,
				nonce:             nonce,
				fundReceiver:      fundReceiver,
				programID:         solParams.ProgramId,
				mintAddress:       solParams.MintAddress,
				feeWallet:         solParams.FeeWallet,
				nonceAccountKey:   solParams.NonceAccountKey,
				nonceAuthorityKey: solParams.NonceAuthorityKey,
				signerKey:         solParams.SignerKeyId,
				eventID:           tx.Id,
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
			solNonce := types.SolanaNonce{Nonce: nonce.Nonce[:]}
			if err = k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, solNonce); err != nil {
				return fmt.Errorf("LastUsedSolanaNonce.Set: %w", err)
			}
			return k.zentpKeeper.UpdateMint(ctx, tx.Id, tx)
		},
		// txContinuationCallback: For Solana, this function acts as a status and timeout checker for in-flight transactions.
		// It is called on every block to check if a pending transaction has been confirmed, has timed out (BTL),
		// or requires a retry. It manages the lifecycle of the pending Solana transaction.
		func(tx *zentptypes.Bridge) error {
			// If we don't have consensus on SolanaMintEventsHash we cannot reliably determine
			// whether the associated event has arrived on-chain. In that case we should *not* run
			// any retry or timeout logic, otherwise we risk redelivering the same transaction over
			// and over again without a reliable confirmation signal.

			if !fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
				k.Logger(ctx).Debug("Skipping Solana ROCK mint retry/timeout checks – no consensus on SolanaMintEventsHash", "tx_id", tx.Id)
				return nil
			}

			// If BlockHeight is 0, this transaction was either just dispatched in the current block
			// by the txDispatchCallback, or it has been reset for a full retry by prior logic in this callback.
			if tx.BlockHeight == 0 {
				k.Logger(ctx).Debug("Solana ROCK Mint Nonce Update: tx.BlockHeight is 0. No BTL/event check in this invocation.", "tx_id", tx.Id)
				return nil
			}

			solParams := k.zentpKeeper.GetSolanaParams(ctx)
			k.Logger(ctx).Info("Solana ROCK Mint Status Check Begin", "tx_id", tx.Id, "tx_block_height", tx.BlockHeight, "btl", solParams.Btl, "current_chain_height", ctx.BlockHeight(), "awaiting_event_since", tx.AwaitingEventSince)

			if ctx.BlockHeight() > tx.BlockHeight+solParams.Btl {
				*tx = k.processBtlSolanaROCKMint(ctx, *tx, oracleData, *solParams)
			}

			// --- Secondary Event Arrival Timeout Check ---
			// This check applies if tx.AwaitingEventSince was set (indicating nonce advanced)
			// and was not subsequently cleared by the BTL check itself.
			if tx.AwaitingEventSince > 0 {
				*tx = k.processSecondaryTimeoutSolanaROCKMint(ctx, *tx, oracleData, *solParams)
			}

			k.Logger(ctx).Info("Solana ROCK Mint Status Check End", "tx_id", tx.Id, "tx_block_height_after_checks", tx.BlockHeight, "awaiting_event_since_after_checks", tx.AwaitingEventSince)
			return k.zentpKeeper.UpdateMint(ctx, tx.Id, tx)
		},
	)

}

// processROCKBurns processes pending mint transactions.
func (k *Keeper) processSolanaROCKMintEvents(ctx sdk.Context, oracleData OracleData) {

	pendingMints, err := k.zentpKeeper.GetMintsWithStatusPending(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return
		}
		k.Logger(ctx).Error("GetMintsWithPendingStatus: ", err.Error())
		return
	}

	if len(pendingMints) == 0 {
		return
	}

	for _, pendingMint := range pendingMints {
		tx, err := k.treasuryKeeper.GetSignTransactionRequest(ctx, pendingMint.TxId)
		if err != nil {
			k.Logger(ctx).Error("GetSignTransactionRequest: ", err.Error())
			return
		}
		sigReq, err := k.treasuryKeeper.GetSignRequest(ctx, tx.SignRequestId)
		if err != nil {
			k.Logger(ctx).Error("GetSignRequest: ", err.Error())
		}

		var signatures []byte
		var sigHash [32]byte

		for _, id := range sigReq.ChildReqIds {
			childReq, err := k.treasuryKeeper.GetSignRequest(ctx, id)
			if err != nil {
				k.Logger(ctx).Error("GetSignRequest: ", err.Error())
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

				// Perform the paramount check *before* any state change.
				if err := k.zentpKeeper.CheckROCKSupplyCap(ctx, sdkmath.ZeroInt()); err != nil {
					// ABORT. The invariant would be violated. Do not complete the bridge.
					k.Logger(ctx).Error("CRITICAL INVARIANT VIOLATION DETECTED: A mint on Solana would breach the 1bn cap. Aborting bridge completion.", "bridge_id", pendingMint.Id, "error", err.Error())
					pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_FAILED
					if err := k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint); err != nil {
						k.Logger(ctx).Error("CRITICAL: Failed to update mint status to FAILED after invariant violation.", "error", err, "bridge_id", pendingMint.Id)
					}
					continue // Move to the next event, this one is terminated.
				}

				// --- Invariant holds, proceed with bridge completion ---

				// Capture total supply before the state change for conservation check.
				totalSupplyBefore, err := k.zentpKeeper.GetTotalROCKSupply(ctx)
				if err != nil {
					k.Logger(ctx).Error("CRITICAL: Failed to get total rock supply before bridge completion.", "error", err.Error(), "bridge_id", pendingMint.Id)
					continue
				}

				err = k.bankKeeper.BurnCoins(ctx, zentptypes.ModuleName, sdk.NewCoins(sdk.NewCoin(pendingMint.Denom, sdkmath.NewIntFromUint64(pendingMint.Amount))))
				if err != nil {
					k.Logger(ctx).Error("CRITICAL: Failed to burn coins for completed Solana bridge AFTER invariant check. State is now inconsistent.", "denom", pendingMint.Denom, "error", err.Error(), "bridge_id", pendingMint.Id)
					continue
				}

				// Re-fetch solana supply to avoid race conditions within the same block processing multiple events
				solanaSupply, err := k.zentpKeeper.GetSolanaROCKSupply(ctx)
				if err != nil {
					k.Logger(ctx).Error("CRITICAL: Failed to get solana rock supply after burning coins. State is now inconsistent.", "error", err.Error(), "bridge_id", pendingMint.Id)
					continue
				}

				newSolanaSupply := solanaSupply.Add(sdkmath.NewIntFromUint64(pendingMint.Amount))
				if err := k.zentpKeeper.SetSolanaROCKSupply(ctx, newSolanaSupply); err != nil {
					k.Logger(ctx).Error("CRITICAL: Failed to set solana rock supply after burning coins. State is now inconsistent.", "error", err.Error(), "bridge_id", pendingMint.Id)
					continue
				}

				if err := k.LastCompletedZentpMintID.Set(ctx, pendingMint.Id); err != nil {
					k.Logger(ctx).Error("Failed to set last completed zentp mint.", "error", err.Error(), "bridge_id", pendingMint.Id)
					continue
				}

				// Perform the supply conservation check.
				totalSupplyAfter, err := k.zentpKeeper.GetTotalROCKSupply(ctx)
				if err != nil {
					k.Logger(ctx).Error("CRITICAL: Failed to get total rock supply after bridge completion for conservation check.", "error", err.Error(), "bridge_id", pendingMint.Id)
					continue
				}

				if !totalSupplyBefore.Equal(totalSupplyAfter) {
					k.Logger(ctx).Error("CRITICAL INVARIANT VIOLATION: Total ROCK supply changed during bridge operation.", "before", totalSupplyBefore.String(), "after", totalSupplyAfter.String(), "bridge_id", pendingMint.Id)
					// Here we should ideally halt the chain or take other drastic measures.
					// For now, we will fail the mint and log a critical error.
					pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_FAILED
					if err := k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint); err != nil {
						k.Logger(ctx).Error("CRITICAL: Failed to update mint status to FAILED after supply conservation violation.", "error", err, "bridge_id", pendingMint.Id)
					}
					continue
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
						sdk.NewAttribute(types.AttributeKeyBridgeAmount, fmt.Sprintf("%d", pendingMint.Amount)),
						sdk.NewAttribute(types.AttributeKeyBurnDestination, pendingMint.RecipientAddress),
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

	pendingMint, err := k.zenBTCKeeper.GetPendingMintTransaction(ctx, firstPendingID)
	if err != nil {
		k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting pending mint transaction from store.", "id", firstPendingID, "error", err)
		return
	}
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Retrieved pending mint transaction.", "pending_mint_id", pendingMint.Id, "zrchain_tx_id", pendingMint.ZrchainTxId, "status", pendingMint.Status, "recipient", pendingMint.RecipientAddress, "amount", pendingMint.Amount)

	if pendingMint.ZrchainTxId == 0 {
		k.Logger(ctx).Warn("ProcessSolanaZenBTCMintEvents: PendingMint has ZrchainTxId == 0. Cannot match with treasury sign requests. Skipping.", "pending_mint_id", pendingMint.Id)
		return
	}

	signTxReq, err := k.treasuryKeeper.GetSignTransactionRequest(ctx, pendingMint.ZrchainTxId)
	if err != nil {
		k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting SignTransactionRequest from treasury.", "zrchain_tx_id_searched", pendingMint.ZrchainTxId, "error", err)
		return
	}
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Retrieved SignTransactionRequest from treasury.", "zrchain_tx_id", signTxReq.Id, "sign_request_id", signTxReq.SignRequestId)

	mainSignReq, err := k.treasuryKeeper.GetSignRequest(ctx, signTxReq.SignRequestId)
	if err != nil {
		k.Logger(ctx).Error("ProcessSolanaZenBTCMintEvents: Error getting main SignRequest from treasury.", "sign_request_id_searched", signTxReq.SignRequestId, "error", err)
		return
	}
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Retrieved main SignRequest from treasury.", "main_sign_request_id", mainSignReq.Id, "child_req_count", len(mainSignReq.ChildReqIds), "status", mainSignReq.Status)

	var signatures [][]byte
	foundAllChildSignatures := true
	for i, childReqID := range mainSignReq.ChildReqIds {
		childReq, err := k.treasuryKeeper.GetSignRequest(ctx, childReqID)
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
			pendingMint.TxHash = event.TxSig
			pendingMint.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
			if err = k.zenBTCKeeper.SetPendingMintTransaction(ctx, pendingMint); err != nil {
				k.Logger(ctx).Error("zenBTCKeeper.SetPendingMintTransaction: ", err.Error())
			}

			// Adjust bedrock BTC for default validator (convert zenBTC minted -> BTC sats)
// NOTE: For TokensBedrock accounting we do not use exchangeRate.
			if err := k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(pendingMint.Amount)); err != nil {
				k.Logger(ctx).Error("error adjusting bedrock BTC on solana mint", "error", err)
			}

			break // Found and processed, no need to check other events for this pending mint.
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
	// We skip EigenLayer unstake/completion, so nonceErrorMsg is not used in this mode.
	_ = nonceErrorMsg
	k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: Started.", "source", source, "incoming_event_count", len(burnEvents))
	if source == "solana" && len(burnEvents) > 0 {
		for i, dbgEvent := range burnEvents {
			k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: Solana event details from oracle.",
				"source", source, "idx", i, "tx_id", dbgEvent.TxID, "log_idx", dbgEvent.LogIndex,
				"chain_id", dbgEvent.ChainID, "amount", dbgEvent.Amount, "destination_addr_hex", hex.EncodeToString(dbgEvent.DestinationAddr), "is_zenbtc", dbgEvent.IsZenBTC)
		}
	}

	processedInThisRun := make(map[string]bool)
	processedTxHashes := make(map[string]bool)
	// Loop over each burn event from oracle to check for new ones.
	for _, burn := range burnEvents {
		eventKey := fmt.Sprintf("%s-%d-%s", burn.TxID, burn.LogIndex, burn.ChainID)
		if processedInThisRun[eventKey] {
			continue
		}
		processedInThisRun[eventKey] = true
		// For Solana events, we now use the explicit flag to distinguish burn types.
		// We skip ROCK burns here. zenBTC burns will have IsZenBTC = true.
		if !burn.IsZenBTC {
			k.Logger(ctx).Debug("StoreNewZenBTCBurnEvents: Skipping event explicitly marked as not a zenBTC burn.", "tx_id", burn.TxID, "log_idx", burn.LogIndex)
			continue
		}

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
			// The explicit check at the top of the loop replaces the need for the zentp keeper check.

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

			// Direct redemption initiation (skip EigenLayer unstake/completion)
			// Use the newly created burn event ID as the redemption ID. Amount is kept in zenBTC units.
			if has, err := k.zenBTCKeeper.HasRedemption(ctx, createdID); err != nil {
				k.Logger(ctx).Error("StoreNewZenBTCBurnEvents: Error checking redemption existence.", "redemption_id", createdID, "error", err)
			} else if !has {
				red := zenbtctypes.Redemption{
					Data: zenbtctypes.RedemptionData{
						Id:                 createdID,
						DestinationAddress: burn.DestinationAddr,
						Amount:             burn.Amount, // zenBTC amount burned
					},
					Status: zenbtctypes.RedemptionStatus_UNSTAKED,
				}
				if err := k.zenBTCKeeper.SetRedemption(ctx, createdID, red); err != nil {
					k.Logger(ctx).Error("StoreNewZenBTCBurnEvents: Error creating redemption for burn.", "redemption_id", createdID, "error", err)
				} else {
					k.Logger(ctx).Info("StoreNewZenBTCBurnEvents: Created redemption for burn (direct mode).", "redemption_id", createdID)
				}
			}

			processedTxHashes[burn.TxID] = true
		} else {
			k.Logger(ctx).Debug("StoreNewZenBTCBurnEvents: Skipping pre-existing event.", "source", source, "tx_id", burn.TxID, "log_idx", burn.LogIndex)
		}
	}

	// Clear any corresponding backfill requests for successfully processed events.
	k.ClearProcessedBackfillRequests(ctx, types.EventType_EVENT_TYPE_ZENBTC_BURN, processedTxHashes)
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
		// pendingGetter: Fetches burn events that have been observed on-chain (Ethereum or Solana)
		// and are in the BURNED state, ready to be processed for unstaking from EigenLayer.
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
		// txDispatchCallback: Constructs and submits an Ethereum transaction to unstake assets from EigenLayer.
		// This corresponds to the amount of zenBTC that was burned.
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
		// txContinuationCallback: This is called when the unstake transaction's nonce has been confirmed on-chain.
		// It updates the burn event's status to UNSTAKING to reflect that the unstaking process has started.
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
		// pendingGetter: Fetches redemptions that are in the INITIATED state.
		// These are unstaking requests from EigenLayer that are ready to be completed.
		func(ctx sdk.Context) ([]zenbtctypes.Redemption, error) {
			firstPendingID, err := k.zenBTCKeeper.GetFirstPendingRedemption(ctx)
			if err != nil {
				firstPendingID = 0
			}
			return k.GetRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_INITIATED, 2, firstPendingID)
		},
		// txDispatchCallback: Constructs and submits an Ethereum transaction to call the 'complete' function
		// on the EigenLayer contracts, which finalizes the unstaking process.
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
		// txContinuationCallback: This is called when the 'complete' transaction's nonce has been confirmed on-chain.
		// It updates the redemption status to UNSTAKED and requests a nonce for the staker key,
		// anticipating that the released funds might be re-staked.
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

	redemptions, err := k.GetRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_AWAITING_SIGN, 0, startingIndex)
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
		signReq, err := k.treasuryKeeper.GetSignRequest(ctx, redemption.Data.SignReqId)
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

			// Adjust bedrock BTC for default validator (subtract released BTC)
// NOTE: For TokensBedrock accounting we do not use exchangeRate. Use redemption data amount for v1.
			if err := k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(redemption.Data.Amount).Neg()); err != nil {
				k.Logger(ctx).Error("error adjusting bedrock BTC on redemption fulfilment", "error", err)
			}

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
	processedInThisRun := make(map[string]bool)
	for _, e := range oracleData.SolanaBurnEvents {
		// Only process events that are explicitly marked as ROCK burns.
		if e.IsZenBTC {
			continue // This is a zenBTC burn, skip it.
		}

		if _, ok := processedInThisRun[e.TxID]; ok {
			continue
		}

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
			processedInThisRun[e.TxID] = true
		}
	}

	if len(toProcess) == 0 {
		return
	}

	processedTxHashes := make(map[string]bool)

	// TODO do cleanup on error. e.g. burn minted funds if there is an error sendig them to the recipient, or adding of the bridge fails
	for _, burn := range toProcess {
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

		if err := k.zentpKeeper.CheckCanBurnFromSolana(ctx, sdkmath.NewIntFromUint64(burn.Amount)); err != nil {
			k.Logger(ctx).Error("invariant check failed for solana rock burn", "error", err.Error(), "amount", burn.Amount)
			continue
		}

		_, bridgeFee, err := k.zentpKeeper.GetBridgeFeeParams(ctx)
		if err != nil {
			k.Logger(ctx).Error("GetBridgeFeeParams: ", err.Error())
			continue
		}

		bridgeFeeCoins, err := k.zentpKeeper.GetBridgeFeeAmount(ctx, burn.Amount, bridgeFee)
		if err != nil {
			k.Logger(ctx).Error("GetBridgeFeeAmount: ", err.Error())
			continue
		}

		bridgeAmount := sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(burn.Amount).Sub(bridgeFeeCoins.AmountOf(params.BondDenom))))

		coins := sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(burn.Amount)))
		if err := k.bankKeeper.MintCoins(ctx, zentptypes.ModuleName, coins); err != nil {
			k.Logger(ctx).Error(fmt.Errorf("MintCoins: %w", err).Error())
			continue
		}

		solanaSupply, err := k.zentpKeeper.GetSolanaROCKSupply(ctx)
		if err != nil {
			k.Logger(ctx).Error("GetSolanaROCKSupply: " + err.Error())
			continue
		}
		newSolanaSupply := solanaSupply.Sub(sdkmath.NewIntFromUint64(burn.Amount))
		if newSolanaSupply.IsNegative() {
			k.Logger(ctx).Error("solana rock supply underflow", "new_supply", newSolanaSupply.String())
			continue
		}
		if err := k.zentpKeeper.SetSolanaROCKSupply(ctx, newSolanaSupply); err != nil {
			k.Logger(ctx).Error("SetSolanaROCKSupply: ", err.Error())
			continue
		}

		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zentptypes.ModuleName, accAddr, bridgeAmount); err != nil {
			k.Logger(ctx).Error(fmt.Errorf("SendCoinsFromModuleToAccount: %w", err).Error())
		}

		if bridgeFeeCoins.AmountOf(params.BondDenom).IsPositive() {
			if err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, zentptypes.ModuleName, zentptypes.ZentpCollectorName, bridgeFeeCoins); err != nil {
				k.Logger(ctx).Error(fmt.Errorf("SendCoinsFromModuleToModule: %w", err).Error())
			}
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
			continue
		}
		processedTxHashes[burn.TxID] = true

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeValidation,
				sdk.NewAttribute(types.AttributeKeyBridgeAmount, fmt.Sprintf("%d", burn.Amount)),
				sdk.NewAttribute(types.AttributeKeyBridgeFee, bridgeFeeCoins.AmountOf(params.BondDenom).String()),
				sdk.NewAttribute(types.AttributeKeyBurnDestination, addr),
			),
		)
	}

	// Now that events are processed, clear any corresponding backfill requests.
	k.ClearProcessedBackfillRequests(ctx, types.EventType_EVENT_TYPE_ZENTP_BURN, processedTxHashes)
}
