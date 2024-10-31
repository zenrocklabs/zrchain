package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	if voteExt.IsInvalid() {
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
	avsDelegationsHash, err := deriveAVSContractStateHash(oracleData.AVSDelegationsMap)
	if err != nil {
		return VoteExtension{}, fmt.Errorf("error deriving AVS contract delegation state hash: %w", err)
	}

	bitcoinData, err := k.sidecarClient.GetLatestBitcoinBlockHeader(ctx, &sidecar.LatestBitcoinBlockHeaderRequest{ChainName: "testnet4"}) // TODO: use config
	if err != nil {
		return VoteExtension{}, err
	}

	nonce, err := k.lookupEthereumNonce(ctx)
	if err != nil {
		return VoteExtension{}, err
	}

	voteExt := VoteExtension{
		ZRChainBlockHeight: height,
		ROCKUSDPrice:       oracleData.ROCKUSDPrice,
		ETHUSDPrice:        oracleData.ETHUSDPrice,
		AVSDelegationsHash: avsDelegationsHash[:],
		BtcBlockHeight:     bitcoinData.BlockHeight,
		BtcMerkleRoot:      bitcoinData.BlockHeader.MerkleRoot,
		EthBlockHeight:     oracleData.EthBlockHeight,
		EthBlockHash:       oracleData.EthBlockHash,
		EthGasLimit:        oracleData.EthGasLimit,
		EthBaseFee:         oracleData.EthBaseFee,
		EthTipCap:          oracleData.EthTipCap,
		RequestedEthNonce:  nonce,
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

	if voteExt.IsInvalid() {
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
		k.Logger(ctx).Error("error in getValidatedOracleData; injecting empty oracle data", "height", req.Height, "error", err)
		return nil, nil
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

	voteExt, err := k.GetSuperMajorityVE(ctx, req.Height, recoveredOracleData.ConsensusData)
	if err != nil {
		return REJECT_PROPOSAL, fmt.Errorf("error retrieving supermajority vote extensions: %w", err)
	}
	if reflect.DeepEqual(voteExt, VoteExtension{}) {
		k.Logger(ctx).Warn("accepting empty vote extension", "height", req.Height)
		return ACCEPT_PROPOSAL, nil
	}

	if err := k.validateOracleData(voteExt, &recoveredOracleData); err != nil {
		return REJECT_PROPOSAL, err
	}

	k.recordMismatchedVoteExtensions(ctx, req.Height, voteExt, recoveredOracleData.ConsensusData)

	return ACCEPT_PROPOSAL, nil
}

// PreBlocker is called before each block to process oracle data and update state.
// We don't return errors in the PreBlocker as this would halt the chain. Instead, we log errors and continue.
func (k *Keeper) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) error {
	if len(req.Txs) == 0 {
		return nil
	}

	if req.Height == 1 || !VoteExtensionsEnabled(ctx) {
		return nil
	}

	voteExtTx := req.Txs[0] // vote extension is always the first transaction in the block

	if !ContainsVoteExtension(voteExtTx, k.txDecoder) {
		return nil
	}

	oracleData, err := k.unmarshalOracleData(voteExtTx)
	if err != nil {
		k.Logger(ctx).Warn("error getting oracle data from tx", "height", req.Height, "error", err)
		return nil
	}

	k.updateAssetPrices(ctx, oracleData)

	k.updateValidatorStakes(ctx, oracleData)

	k.updateAVSDelegationStore(ctx, oracleData)

	k.storeBitcoinBlockHeader(ctx, oracleData)

	k.createMintTransaction(ctx, oracleData)

	k.recordNonVotingValidators(ctx, req)

	return nil
}

func (k *Keeper) getValidatedOracleData(ctx context.Context, voteExt VoteExtension) (*OracleData, *VoteExtension, error) {
	oracleData, err := k.GetSidecarStateByEthHeight(ctx, voteExt.EthBlockHeight)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching oracle state: %w", err)
	}

	bitcoinData, err := k.sidecarClient.GetBitcoinBlockHeaderByHeight(
		ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{ChainName: "testnet4", BlockHeight: voteExt.BtcBlockHeight}, // TODO: use config
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching bitcoin header: %w", err)
	}

	oracleData.BtcBlockHeight = bitcoinData.BlockHeight
	oracleData.BtcBlockHeader = *bitcoinData.BlockHeader
	oracleData.RequestedEthNonce = voteExt.RequestedEthNonce

	if err := k.validateOracleData(voteExt, oracleData); err != nil {
		return nil, nil, err
	}

	return oracleData, &voteExt, nil
}

func (k *Keeper) marshalOracleData(req *abci.RequestPrepareProposal, oracleData *OracleData) ([]byte, error) {
	oracleDataBz, err := json.Marshal(oracleData)
	if err != nil {
		return nil, fmt.Errorf("error encoding oracle data: %w", err)
	}

	if int64(len(oracleDataBz)) > req.MaxTxBytes {
		return nil, fmt.Errorf("oracle data too large: %d > %d", len(oracleDataBz), req.MaxTxBytes)
	}

	return oracleDataBz, nil
}

func (k *Keeper) unmarshalOracleData(tx []byte) (OracleData, error) {
	if len(tx) == 0 {
		return OracleData{}, fmt.Errorf("no transactions in block")
	}

	var oracleData OracleData
	if err := json.Unmarshal(tx, &oracleData); err != nil {
		return OracleData{}, err
	}

	return oracleData, nil
}

func (k *Keeper) updateAssetPrices(ctx sdk.Context, oracleData OracleData) {
	if err := k.AssetPrices.Set(ctx, "rock", types.AssetPrice{PriceUSD: oracleData.ROCKUSDPrice}); err != nil {
		k.Logger(ctx).Error("error setting ROCK price", "height", ctx.BlockHeight(), "err", err)
	}

	if err := k.AssetPrices.Set(ctx, "eth", types.AssetPrice{PriceUSD: oracleData.ETHUSDPrice}); err != nil {
		k.Logger(ctx).Error("error setting ETH price", "height", ctx.BlockHeight(), "err", err)
	}
}

func (k *Keeper) updateValidatorStakes(ctx sdk.Context, oracleData OracleData) {
	validatorInAVSDelegationSet := make(map[string]bool)

	for _, delegation := range oracleData.ValidatorDelegations {
		valAddr, err := sdk.ValAddressFromBech32(delegation.Validator)
		if err != nil {
			k.Logger(ctx).Error("invalid validator address: "+delegation.Validator, "err", err)
			continue
		}

		validator, err := k.GetZenrockValidator(ctx, valAddr)
		if err != nil {
			k.Logger(ctx).Debug(
				"error retrieving validator "+delegation.Validator, "err", err, "reason", "incorrect address entered in delegation",
			)
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
	for validatorAddr, delegatorMap := range oracleData.AVSDelegationsMap {
		for delegatorAddr, amount := range delegatorMap {
			if err := k.AVSDelegations.Set(ctx, collections.Join(validatorAddr, delegatorAddr), math.NewIntFromBigInt(amount)); err != nil {
				k.Logger(ctx).Error("error setting AVS delegations", "err", err)
			}
		}
	}
}

func (k *Keeper) storeBitcoinBlockHeader(ctx sdk.Context, oracleData OracleData) {
	if oracleData.BtcBlockHeight != 0 && oracleData.BtcBlockHeader.MerkleRoot != "" {
		// TODO(sasha): check if entry exists for this height, if it does, we need to add last 6 block heights to a slice and
		// store all of those merkle roots from the Bitcoin node one at a time per block to not add too much latency to VEs
		if err := k.BtcBlockHeaders.Set(ctx, oracleData.BtcBlockHeight, oracleData.BtcBlockHeader); err != nil {
			k.Logger(ctx).Error("error setting Bitcoin block header", "height", oracleData.BtcBlockHeight, "err", err)
		}
	}
}

func (k *Keeper) createMintTransaction(ctx sdk.Context, oracleData OracleData) error {
	requested, err := k.EthereumNonceRequested.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return fmt.Errorf("error getting EthereumNonceRequested state: %w", err)
		}
		requested = false
		if err := k.EthereumNonceRequested.Set(ctx, requested); err != nil {
			return fmt.Errorf("error setting EthereumNonceRequested state: %w", err)
		}
	}
	if !requested {
		return nil
	}

	pendingMints, err := k.PendingMintTransactions.Get(ctx)
	if err != nil {
		return fmt.Errorf("error getting pending mint transactions: %w", err)
	}
	if len(pendingMints.Txs) == 0 {
		return fmt.Errorf("no pending mint transactions")
	}
	tx := pendingMints.Txs[0]

	unsignedMintTxHash, unsignedMintTx, err := k.constructMintTx(
		ctx,
		tx.RecipientAddress,
		tx.ChainId,
		tx.Amount,
		0, // TODO: update fee (currently hardcoded to 0)
		oracleData.RequestedEthNonce,
		oracleData.EthGasLimit,
		oracleData.EthBaseFee,
		oracleData.EthTipCap,
	)
	if err != nil {
		return fmt.Errorf("error constructing mint transaction: %w", err)
	}

	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataEthereum{ChainId: tx.ChainId})
	if err != nil {
		return fmt.Errorf("error creating metadata: %w", err)
	}

	if _, err := k.treasuryKeeper.HandleSignTransactionRequest(ctx, &treasurytypes.MsgNewSignTransactionRequest{
		Creator:             tx.Creator,
		KeyId:               tx.KeyId,
		WalletType:          tx.ChainType,
		UnsignedTransaction: unsignedMintTx,
		Metadata:            metadata,
		NoBroadcast:         false,
	}, unsignedMintTxHash); err != nil {
		return fmt.Errorf("error creating sign transaction request for zenBTC mint: %w", err)
	}

	pendingMints.Txs = pendingMints.Txs[1:]
	if err := k.PendingMintTransactions.Set(ctx, pendingMints); err != nil {
		return fmt.Errorf("error setting pending mint transactions: %w", err)
	}

	if len(pendingMints.Txs) == 0 {
		if err := k.EthereumNonceRequested.Set(ctx, false); err != nil {
			return fmt.Errorf("error setting EthereumNonceRequested state: %w", err)
		}
	}

	return nil
}

func (k *Keeper) recordMismatchedVoteExtensions(ctx sdk.Context, height int64, canonicalVoteExt VoteExtension, consensusData abci.ExtendedCommitInfo) {
	canonicalVoteExtBz, err := json.Marshal(canonicalVoteExt)
	if err != nil {
		k.Logger(ctx).Error("error marshalling canonical vote extension", "height", height, "error", err)
		return
	}

	for _, v := range consensusData.Votes {
		if !bytes.Equal(v.VoteExtension, canonicalVoteExtBz) {
			info, err := k.ValidationInfos.Get(ctx, height)
			if err != nil {
				info = types.ValidationInfo{}
			}
			info.MismatchedVoteExtensions = append(info.MismatchedVoteExtensions, hex.EncodeToString(v.Validator.Address))
			if err := k.ValidationInfos.Set(ctx, height, info); err != nil {
				k.Logger(ctx).Error("error setting validation info", "height", height, "error", err)
			}
		}
	}
}

func (k *Keeper) recordNonVotingValidators(ctx sdk.Context, req *abci.RequestFinalizeBlock) {
	for _, v := range req.DecidedLastCommit.Votes {
		if v.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			info, err := k.ValidationInfos.Get(ctx, req.Height)
			if err != nil {
				info = types.ValidationInfo{}
			}
			info.NonVotingValidators = append(info.NonVotingValidators, hex.EncodeToString(v.Validator.Address))
			if err := k.ValidationInfos.Set(ctx, req.Height, info); err != nil {
				k.Logger(ctx).Error("error setting validation info", "height", req.Height, "error", err)
			}
		}
	}
}

func (k *Keeper) validateOracleData(voteExt VoteExtension, oracleData *OracleData) error {
	avsDelegationsHash, err := deriveAVSContractStateHash(oracleData.AVSDelegationsMap)
	if err != nil {
		return fmt.Errorf("error deriving AVS contract delegation state hash: %w", err)
	}

	if !bytes.Equal(voteExt.AVSDelegationsHash, avsDelegationsHash[:]) {
		return fmt.Errorf("AVS contract delegation state hash mismatch, expected %x, got %x", voteExt.AVSDelegationsHash, avsDelegationsHash)
	}

	if !voteExt.ROCKUSDPrice.Equal(oracleData.ROCKUSDPrice) {
		return fmt.Errorf("ROCK/USD price mismatch, expected %s, got %s", voteExt.ROCKUSDPrice, oracleData.ROCKUSDPrice)
	}

	if !voteExt.ETHUSDPrice.Equal(oracleData.ETHUSDPrice) {
		return fmt.Errorf("ETH/USD price mismatch, expected %s, got %s", voteExt.ETHUSDPrice, oracleData.ETHUSDPrice)
	}

	if voteExt.EthBlockHeight != oracleData.EthBlockHeight {
		return fmt.Errorf("ethereum block height mismatch, expected %d, got %d", voteExt.EthBlockHeight, oracleData.EthBlockHeight)
	}

	if voteExt.EthBlockHash != oracleData.EthBlockHash {
		return fmt.Errorf("ethereum block hash mismatch, expected %s, got %s", voteExt.EthBlockHash, oracleData.EthBlockHash)
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

	if voteExt.BtcMerkleRoot != oracleData.BtcBlockHeader.MerkleRoot {
		return fmt.Errorf("bitcoin merkle root mismatch, expected %s, got %s - height %d", voteExt.BtcMerkleRoot, oracleData.BtcBlockHeader.MerkleRoot, voteExt.BtcBlockHeight)
	}

	if voteExt.RequestedEthNonce != oracleData.RequestedEthNonce {
		return fmt.Errorf("requested Ethereum nonce mismatch, expected %d, got %d", voteExt.RequestedEthNonce, oracleData.RequestedEthNonce)
	}

	return nil
}
