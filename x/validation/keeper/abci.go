package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	abci "github.com/cometbft/cometbft/abci/types"
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
		k.requestMintDispatches(ctx)

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

// requestMintDispatches inspects pending zenBTC mints and requests dispatch for Ethereum and Solana
// by setting the appropriate flags. All errors are logged and do not abort the block.
func (k *Keeper) requestMintDispatches(ctx sdk.Context) {
	// Ethereum
	pendingEVM, err := k.getPendingMintTransactions(ctx,
		zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
		zenbtctypes.WalletType_WALLET_TYPE_EVM,
	)
	if err != nil {
		k.Logger(ctx).Error("error fetching pending EVM zenBTC mints", "error", err)
	} else if len(pendingEVM) > 0 {
		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetEthMinterKeyID(ctx), true); err != nil {
			k.Logger(ctx).Error("error setting Ethereum nonce requested flag for minter", "error", err)
		}
	}

	// Solana
	pendingSol, err := k.getPendingMintTransactions(ctx,
		zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
		zenbtctypes.WalletType_WALLET_TYPE_SOLANA,
	)
	if err != nil {
		k.Logger(ctx).Error("error fetching pending Solana zenBTC mints", "error", err)
		return
	}
	if len(pendingSol) == 0 {
		return
	}
	solParams := k.zenBTCKeeper.GetSolanaParams(ctx)
	if err := k.SolanaNonceRequested.Set(ctx, solParams.NonceAccountKey, true); err != nil {
		k.Logger(ctx).Error("error setting Solana nonce requested flag", "error", err)
	}
	for _, tx := range pendingSol {
		if err := k.SetSolanaZenBTCRequestedAccount(ctx, tx.RecipientAddress, true); err != nil {
			k.Logger(ctx).Error("error setting Solana requested account flag", "recipient", tx.RecipientAddress, "error", err)
		}
	}
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
