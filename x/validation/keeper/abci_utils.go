package keeper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/comet"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cryptoenc "github.com/cometbft/cometbft/crypto/encoding"
	"github.com/cometbft/cometbft/libs/protoio"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"

	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	bindings "github.com/zenrocklabs/zenbtc/bindings"
)

func (k Keeper) GetSidecarState(ctx context.Context, height int64) (*OracleData, error) {
	resp, err := k.sidecarClient.GetSidecarState(ctx, &sidecar.SidecarStateRequest{})
	if err != nil {
		k.Logger(ctx).Error("error fetching operator stakes (GetSidecarState)", "height", height, "error", err)
		return nil, ErrOracleSidecar
	}
	return k.processOracleResponse(ctx, resp)
}

func (k Keeper) GetSidecarStateByEthHeight(ctx context.Context, height uint64) (*OracleData, error) {
	resp, err := k.sidecarClient.GetSidecarStateByEthHeight(ctx, &sidecar.SidecarStateByEthHeightRequest{EthBlockHeight: height})
	if err != nil {
		k.Logger(ctx).Error("error fetching operator stakes (GetSidecarStateByEthHeight)", "height", height, "error", err)
		return nil, ErrOracleSidecar
	}
	return k.processOracleResponse(ctx, resp)
}

func (k Keeper) processOracleResponse(ctx context.Context, resp *sidecar.SidecarStateResponse) (*OracleData, error) {
	var delegations map[string]map[string]*big.Int

	if err := json.Unmarshal(resp.EigenDelegations, &delegations); err != nil {
		return nil, err
	}

	validatorDelegations, err := k.processDelegations(delegations)
	if err != nil {
		k.Logger(ctx).Error("error processing delegations", "error", err)
		return nil, ErrOracleSidecar
	}

	ROCKUSDPrice, err := sdkmath.LegacyNewDecFromStr(resp.ROCKUSDPrice)
	if err != nil {
		k.Logger(ctx).Error("error parsing rock price", "error", err)
		return nil, ErrOracleSidecar
	}

	BTCUSDPrice, err := sdkmath.LegacyNewDecFromStr(resp.BTCUSDPrice)
	if err != nil {
		k.Logger(ctx).Error("error parsing btc price", "error", err)
		return nil, ErrOracleSidecar
	}

	ETHUSDPrice, err := sdkmath.LegacyNewDecFromStr(resp.ETHUSDPrice)
	if err != nil {
		k.Logger(ctx).Error("error parsing eth price", "error", err)
		return nil, ErrOracleSidecar
	}

	return &OracleData{
		EigenDelegationsMap:        delegations,
		ValidatorDelegations:       validatorDelegations,
		EthBlockHeight:             resp.EthBlockHeight,
		EthGasLimit:                resp.EthGasLimit,
		EthBaseFee:                 resp.EthBaseFee,
		EthTipCap:                  resp.EthTipCap,
		SolanaLamportsPerSignature: resp.SolanaLamportsPerSignature,
		EthBurnEvents:              resp.EthBurnEvents,
		Redemptions:                resp.Redemptions,
		ROCKUSDPrice:               ROCKUSDPrice,
		BTCUSDPrice:                BTCUSDPrice,
		ETHUSDPrice:                ETHUSDPrice,
		ConsensusData:              abci.ExtendedCommitInfo{},
	}, nil
}

func (k Keeper) processDelegations(delegations map[string]map[string]*big.Int) ([]ValidatorDelegations, error) {
	validatorTotals := make(map[string]*big.Int)
	for validator, delegatorMap := range delegations {
		total := new(big.Int)
		for _, amount := range delegatorMap {
			total.Add(total, amount)
		}
		validatorTotals[validator] = total
	}

	validatorDelegations := make([]ValidatorDelegations, 0, len(validatorTotals))
	for validator, totalStake := range validatorTotals {
		validatorDelegations = append(validatorDelegations, ValidatorDelegations{
			Validator: validator,
			Stake:     sdkmath.NewIntFromBigInt(totalStake),
		})
	}

	return validatorDelegations, nil
}

func (k Keeper) GetSuperMajorityVE(ctx context.Context, currentHeight int64, extCommit abci.ExtendedCommitInfo) (VoteExtension, error) {
	votesPerVoteExt := make(map[string]*VEWithVotePower)
	var totalVotePower int64

	for _, vote := range extCommit.Votes {
		totalVotePower += vote.Validator.Power

		voteExt, err := k.validateVote(ctx, vote, currentHeight)
		if err != nil {
			continue
		}

		updateVotesPerVE(votesPerVoteExt, voteExt, vote.Validator.Power)
	}

	if len(votesPerVoteExt) == 0 {
		return VoteExtension{}, nil
	}

	mostVotedVE := getMostVotedVE(votesPerVoteExt)

	finalVoteExt, err := unmarshalVE(mostVotedVE.VoteExtension)
	if err != nil {
		return VoteExtension{}, err
	}

	if !hasReachedSupermajority(totalVotePower, mostVotedVE.VotePower) {
		k.Logger(ctx).Warn("consensus not reached on vote extension",
			"required_vote_power", requisiteVotePower(totalVotePower),
			"actual_vote_power", mostVotedVE.VotePower)
		return VoteExtension{}, nil
	}

	return finalVoteExt, nil
}

func (k Keeper) validateVote(ctx context.Context, vote abci.ExtendedVoteInfo, currentHeight int64) (VoteExtension, error) {
	if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit || len(vote.VoteExtension) == 0 {
		return VoteExtension{}, fmt.Errorf("invalid vote")
	}

	var voteExt VoteExtension
	if err := json.Unmarshal(vote.VoteExtension, &voteExt); err != nil {
		return VoteExtension{}, err
	}

	voteExt.ZRChainBlockHeight = currentHeight - 1
	if voteExt.IsInvalid(k.Logger(ctx)) {
		return VoteExtension{}, fmt.Errorf("invalid vote extension")
	}

	return voteExt, nil
}

func getVESubset(ve VoteExtension) VoteExtension {
	ve.EthBlockHeight = 0
	ve.EthBaseFee = 0
	ve.EthTipCap = 0
	ve.EthGasLimit = 0
	return ve
}

func updateVotesPerVE(votesPerVoteExt map[string]*VEWithVotePower, voteExt VoteExtension, votePower int64) {
	// Use the subset for the key
	marshaledSubsetVE, err := json.Marshal(getVESubset(voteExt))
	if err != nil {
		return
	}

	key := hex.EncodeToString(marshaledSubsetVE)
	if existingVE, ok := votesPerVoteExt[key]; ok {
		existingVE.VotePower += votePower
	} else {
		// Store the full VE
		fullMarshaledVE, err := json.Marshal(voteExt)
		if err != nil {
			return
		}
		votesPerVoteExt[key] = &VEWithVotePower{
			VoteExtension: fullMarshaledVE,
			VotePower:     votePower,
		}
	}
}

func getMostVotedVE(votesPerVoteExt map[string]*VEWithVotePower) *VEWithVotePower {
	var mostVotedVE *VEWithVotePower
	for _, voteExt := range votesPerVoteExt {
		if mostVotedVE == nil || voteExt.VotePower > mostVotedVE.VotePower {
			mostVotedVE = voteExt
		}
	}
	return mostVotedVE
}

func unmarshalVE(voteExtensionBytes []byte) (VoteExtension, error) {
	var voteExt VoteExtension
	err := json.Unmarshal(voteExtensionBytes, &voteExt)
	return voteExt, err
}

func hasReachedSupermajority(totalVotePower, mostVotedVEVotePower int64) bool {
	return mostVotedVEVotePower >= requisiteVotePower(totalVotePower)
}

func requisiteVotePower(totalVotePower int64) int64 {
	return ((totalVotePower * 2) / 3) + 1
}

func deriveHash[T any](data T) ([32]byte, error) {
	dataBz, err := json.Marshal(data)
	if err != nil {
		return [32]byte{}, fmt.Errorf("error encoding data: %w", err)
	}
	return sha256.Sum256(dataBz), nil
}

// ref: https://github.com/cosmos/cosmos-sdk/blob/c64d1010800d60677cc25e2fca5b3d8c37b683cc/baseapp/abci_utils.go#L44
func ValidateVoteExtensions(ctx sdk.Context, validationKeeper baseapp.ValidatorStore, currentHeight int64, chainID string, extCommit abci.ExtendedCommitInfo) error {
	marshalDelimitedFn := func(msg proto.Message) ([]byte, error) {
		var buf bytes.Buffer
		if _, err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	totalVotePower, voteExtVotePower := int64(0), int64(0)
	consParams := ctx.ConsensusParams()

	// Check that both extCommit + commit are ordered in accordance with vp/address.
	if err := validateExtendedCommitAgainstLastCommit(extCommit, ctx.CometInfo().GetLastCommit()); err != nil {
		return err
	}

	for _, vote := range extCommit.Votes {
		totalVotePower += vote.Validator.Power

		if len(vote.ExtensionSignature) == 0 {
			continue
		}

		if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		if consParams.Abci == nil || consParams.Abci.VoteExtensionsEnableHeight == 0 ||
			currentHeight <= consParams.Abci.VoteExtensionsEnableHeight {
			if len(vote.VoteExtension) > 0 || len(vote.ExtensionSignature) > 0 {
				return fmt.Errorf("received VE, but VEs are disabled at current height %d", currentHeight)
			}
			continue
		}

		pubKeyProto, err := validationKeeper.GetPubKeyByConsAddr(ctx, vote.Validator.Address)
		if err != nil {
			continue
		}

		cmtPubKey, err := cryptoenc.PubKeyFromProto(pubKeyProto)
		if err != nil {
			continue
		}

		voteExt := cmtproto.CanonicalVoteExtension{
			Extension: vote.VoteExtension,
			Height:    currentHeight - 1,
			Round:     int64(extCommit.Round),
			ChainId:   chainID,
		}

		voteExtSignature, err := marshalDelimitedFn(&voteExt)
		if err != nil {
			continue
		}

		if !cmtPubKey.VerifySignature(voteExtSignature, vote.ExtensionSignature) {
			continue
		}

		voteExtVotePower += vote.Validator.Power
	}

	if totalVotePower > 0 {
		requiredVotePower := ((totalVotePower * 2) / 3) + 1 // for supermajority
		if voteExtVotePower < requiredVotePower {
			return fmt.Errorf("consensus not reached on vote extension at height %d", currentHeight)
		}
	}

	return nil
}

func validateExtendedCommitAgainstLastCommit(ec abci.ExtendedCommitInfo, lc comet.CommitInfo) error {
	// check that the rounds are the same
	if ec.Round != lc.Round() {
		return fmt.Errorf("extended commit round %d does not match last commit round %d", ec.Round, lc.Round())
	}

	// check that the # of votes are the same
	if len(ec.Votes) != lc.Votes().Len() {
		return fmt.Errorf("extended commit votes length %d does not match last commit votes length %d", len(ec.Votes), lc.Votes().Len())
	}

	// check sort order of extended commit votes
	if !slices.IsSortedFunc(ec.Votes, func(vote1, vote2 abci.ExtendedVoteInfo) int {
		if vote1.Validator.Power == vote2.Validator.Power {
			return bytes.Compare(vote1.Validator.Address, vote2.Validator.Address) // addresses sorted in ascending order (used to break vp conflicts)
		}
		return -int(vote1.Validator.Power - vote2.Validator.Power) // vp sorted in descending order
	}) {
		return fmt.Errorf("extended commit votes are not sorted by voting power")
	}

	addressCache := make(map[string]struct{}, len(ec.Votes))
	// check that consistency between LastCommit and ExtendedCommit
	for i, vote := range ec.Votes {
		// cache addresses to check for duplicates
		if _, ok := addressCache[string(vote.Validator.Address)]; ok {
			return fmt.Errorf("extended commit vote address %X is duplicated", vote.Validator.Address)
		}
		addressCache[string(vote.Validator.Address)] = struct{}{}

		if !bytes.Equal(vote.Validator.Address, lc.Votes().Get(i).Validator().Address()) {
			return fmt.Errorf("extended commit vote address %X does not match last commit vote address %X", vote.Validator.Address, lc.Votes().Get(i).Validator().Address())
		}
		if vote.Validator.Power != lc.Votes().Get(i).Validator().Power() {
			return fmt.Errorf("extended commit vote power %d does not match last commit vote power %d", vote.Validator.Power, lc.Votes().Get(i).Validator().Power())
		}
	}

	return nil
}

func (k *Keeper) lookupEthereumNonce(ctx context.Context, keyID uint64) (uint64, error) {
	addr, err := k.getAddressByKeyID(ctx, keyID, treasurytypes.WalletType_WALLET_TYPE_EVM)
	if err != nil {
		return 0, fmt.Errorf("error getting address for key ID %d: %w", keyID, err)
	}

	nonceResp, err := k.sidecarClient.GetLatestEthereumNonceForAccount(ctx, &sidecar.LatestEthereumNonceForAccountRequest{Address: addr})
	if err != nil {
		return 0, fmt.Errorf("error fetching Ethereum nonce: %w", err)
	}

	return nonceResp.Nonce, nil
}

func (k *Keeper) constructEthereumTx(addr common.Address, chainID uint64, data []byte, nonce, gasLimit, baseFee, tipCap uint64) ([]byte, []byte, error) {
	// TODO: whitelist more chain IDs before mainnet upgrade
	if chainID != 17000 {
		return nil, nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
	chainIDBigInt := new(big.Int).SetUint64(chainID)

	// Set minimum priority fee of 0.05 Gwei
	minTipCap := new(big.Int).SetUint64(50000000)
	gasTipCap := new(big.Int).SetUint64(tipCap)
	if gasTipCap.Cmp(minTipCap) < 0 {
		gasTipCap = minTipCap
	}

	// Add 10% buffer to base fee to account for increases
	baseFeeWithBuffer := new(big.Int).Mul(
		new(big.Int).SetUint64(baseFee),
		new(big.Int).SetUint64(11),
	)
	baseFeeWithBuffer.Div(baseFeeWithBuffer, big.NewInt(10))

	// Set max fee to 2.5x buffered base fee + tip
	gasFeeCap := new(big.Int).Mul(baseFeeWithBuffer, big.NewInt(25))
	gasFeeCap.Div(gasFeeCap, big.NewInt(10))
	gasFeeCap.Add(gasFeeCap, gasTipCap)

	unsignedTx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:    chainIDBigInt,
		Nonce:      nonce,
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasFeeCap,
		Gas:        gasLimit,
		To:         &addr,
		Value:      big.NewInt(0), // we shouldn't send any ETH
		Data:       data,
		AccessList: nil,
		V:          big.NewInt(0),
		R:          big.NewInt(0),
		S:          big.NewInt(0),
	})

	unsignedTxBz, err := unsignedTx.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	signer := ethtypes.LatestSignerForChainID(chainIDBigInt)

	return signer.Hash(unsignedTx).Bytes(), unsignedTxBz, nil
}

func (k *Keeper) constructStakeTx(ctx context.Context, chainID, amount, nonce, gasLimit, baseFee, tipCap uint64) ([]byte, []byte, error) {
	encodedMintData, err := EncodeStakeCallData(new(big.Int).SetUint64(amount))
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.zenBTCKeeper.GetEthBatcherAddr(ctx))
	return k.constructEthereumTx(addr, chainID, encodedMintData, nonce, gasLimit, baseFee, tipCap)
}

func (k *Keeper) constructMintTx(ctx context.Context, recipientAddr string, chainID, amount, fee, nonce, gasLimit, baseFee, tipCap uint64) ([]byte, []byte, error) {
	encodedMintData, err := EncodeWrapCallData(common.HexToAddress(recipientAddr), new(big.Int).SetUint64(amount), fee)
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.zenBTCKeeper.GetEthTokenAddr(ctx))
	return k.constructEthereumTx(addr, chainID, encodedMintData, nonce, gasLimit, baseFee, tipCap)
}

func (k *Keeper) constructUnstakeTx(ctx context.Context, chainID uint64, destinationAddr []byte, amount, ethNonce, baseFee, tipCap uint64) ([]byte, []byte, error) {
	encodedUnstakeData, err := k.EncodeUnstakeCallData(ctx, destinationAddr, amount)
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.zenBTCKeeper.GetEthBatcherAddr(ctx))
	return k.constructEthereumTx(addr, chainID, encodedUnstakeData, ethNonce, 700000, baseFee, tipCap)
}

func (k *Keeper) constructCompleteTx(ctx context.Context, chainID, redemptionID, ethNonce, baseFee, tipCap uint64) ([]byte, []byte, error) {
	encodedCompleteData, err := k.EncodeCompleteCallData(ctx, redemptionID)
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.zenBTCKeeper.GetEthBatcherAddr(ctx))
	return k.constructEthereumTx(addr, chainID, encodedCompleteData, ethNonce, 300000, baseFee, tipCap)
}

func EncodeStakeCallData(amount *big.Int) ([]byte, error) {
	if !amount.IsUint64() {
		return nil, fmt.Errorf("amount exceeds uint64 max value")
	}

	parsed, err := bindings.ZenBTControllerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %v", err)
	}

	// Pack using the contract binding's ABI for the wrapZenBTC function
	data, err := parsed.Pack(
		"stakeRockBTC",
		amount.Uint64(),
		bindings.ISignatureUtilsSignatureWithExpiry{Signature: []byte{}, Expiry: big.NewInt(0)}, [32]byte{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode stake call data: %v", err)
	}

	return data, nil
}

func EncodeWrapCallData(recipientAddr common.Address, amount *big.Int, fee uint64) ([]byte, error) {
	if !amount.IsUint64() {
		return nil, fmt.Errorf("amount exceeds uint64 max value")
	}

	parsed, err := bindings.ZenBTCMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %v", err)
	}

	// Pack using the contract binding's ABI for the wrap function
	data, err := parsed.Pack(
		"wrap",
		recipientAddr,
		amount.Uint64(),
		fee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode wrapZenBTC call data: %v", err)
	}

	return data, nil
}

func (k *Keeper) EncodeUnstakeCallData(ctx context.Context, destinationAddr []byte, amount uint64) ([]byte, error) {
	parsed, err := bindings.ZenBTControllerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %v", err)
	}

	data, err := parsed.Pack(
		"unstakeRockBTCInit",
		new(big.Int).SetUint64(amount),
		destinationAddr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode unstakeRockBTCInit call data: %v", err)
	}

	return data, nil
}

func (k *Keeper) EncodeCompleteCallData(ctx context.Context, redemptionID uint64) ([]byte, error) {
	parsed, err := bindings.ZenBTControllerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %v", err)
	}

	data, err := parsed.Pack(
		"unstakeRockBTComplete",
		new(big.Int).SetUint64(redemptionID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode unstakeRockBTComplete call data: %v", err)
	}

	return data, nil
}

func (k *Keeper) getAddressByKeyID(ctx context.Context, keyID uint64, walletType treasurytypes.WalletType) (string, error) {
	q, err := k.treasuryKeeper.KeyByID(ctx, &treasurytypes.QueryKeyByIDRequest{
		Id:         keyID,
		WalletType: walletType,
		Prefixes:   make([]string, 0),
	})
	if err != nil {
		return "", err
	}

	if len(q.Wallets) == 0 {
		return "", fmt.Errorf("no wallets found for key ID %d", keyID)
	}

	return q.Wallets[0].Address, nil
}

func (k *Keeper) bitcoinNetwork(ctx context.Context) string {
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		return "mainnet"
	}
	return "testnet4"
}

func (k *Keeper) retrieveBitcoinHeader(ctx context.Context) (*sidecar.BitcoinBlockHeaderResponse, error) {
	requestedBitcoinHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		requestedBitcoinHeaders = zenbtctypes.RequestedBitcoinHeaders{}
		if err = k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedBitcoinHeaders); err != nil {
			return nil, err
		}
	}

	if len(requestedBitcoinHeaders.Heights) == 0 {
		return k.sidecarClient.GetLatestBitcoinBlockHeader(ctx, &sidecar.LatestBitcoinBlockHeaderRequest{ChainName: k.bitcoinNetwork(ctx)})
	}

	return k.sidecarClient.GetBitcoinBlockHeaderByHeight(ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{ChainName: k.bitcoinNetwork(ctx), BlockHeight: requestedBitcoinHeaders.Heights[0]})
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

func (k *Keeper) unmarshalOracleData(ctx sdk.Context, tx []byte) (OracleData, bool) {
	if len(tx) == 0 {
		k.Logger(ctx).Error("no transactions or vote extension in block")
		return OracleData{}, false
	}

	var oracleData OracleData
	if err := json.Unmarshal(tx, &oracleData); err != nil {
		k.Logger(ctx).Error("error unmarshalling oracle data JSON", "err", err)
		return OracleData{}, false
	}

	return oracleData, true
}

func (k *Keeper) updateAssetPrices(ctx sdk.Context, oracleData OracleData) {
	if err := k.AssetPrices.Set(ctx, types.Asset_ROCK, oracleData.ROCKUSDPrice); err != nil {
		k.Logger(ctx).Error("error setting ROCK price", "height", ctx.BlockHeight(), "err", err)
	}

	if err := k.AssetPrices.Set(ctx, types.Asset_zenBTC, oracleData.BTCUSDPrice); err != nil {
		k.Logger(ctx).Error("error setting BTC price", "height", ctx.BlockHeight(), "err", err)
	}

	if err := k.AssetPrices.Set(ctx, types.Asset_stETH, oracleData.ETHUSDPrice); err != nil {
		k.Logger(ctx).Error("error setting ETH price", "height", ctx.BlockHeight(), "err", err)
	}
}

func (k *Keeper) getZenBTCKeyIDs(ctx context.Context) []uint64 {
	return []uint64{
		k.zenBTCKeeper.GetStakerKeyID(ctx),
		k.zenBTCKeeper.GetEthMinterKeyID(ctx),
		k.zenBTCKeeper.GetUnstakerKeyID(ctx),
		k.zenBTCKeeper.GetCompleterKeyID(ctx),
	}
}

func (k *Keeper) updateNonceState(ctx sdk.Context, keyID uint64, currentNonce uint64) error {
	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx, keyID)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return err
		}
		lastUsedNonce = zenbtctypes.NonceData{
			Nonce:     currentNonce,
			PrevNonce: currentNonce,
			Counter:   0,
			Skip:      true,
		}
	} else {
		if currentNonce == lastUsedNonce.Nonce {
			lastUsedNonce.Counter++
		} else {
			lastUsedNonce.PrevNonce = lastUsedNonce.Nonce
			lastUsedNonce.Nonce = currentNonce
			lastUsedNonce.Counter = 0
		}
		lastUsedNonce.Skip = lastUsedNonce.Counter%8 != 0
	}

	return k.LastUsedEthereumNonce.Set(ctx, keyID, lastUsedNonce)
}

// CalculateZenBTCMintFee calculates the zenBTC fee required for minting
// Returns 0 if BTCUSDPrice is zero
func (k Keeper) CalculateZenBTCMintFee(
	ethBaseFee uint64,
	ethTipCap uint64,
	ethGasLimit uint64,
	btcUSDPrice sdkmath.LegacyDec,
	ethUSDPrice sdkmath.LegacyDec,
	exchangeRate sdkmath.LegacyDec,
) uint64 {
	if btcUSDPrice.IsZero() {
		return 0
	}

	// Calculate total ETH gas fee in wei (base fee + tip)
	baseFeePlusTip := new(big.Int).Add(
		new(big.Int).SetUint64(ethBaseFee),
		new(big.Int).SetUint64(ethTipCap),
	)
	feeInWei := new(big.Int).Mul(
		baseFeePlusTip,
		new(big.Int).SetUint64(ethGasLimit),
	)

	// Convert wei to ETH (divide by 1e18)
	feeInETH := new(big.Float).Quo(
		new(big.Float).SetInt(feeInWei),
		new(big.Float).SetInt64(1e18),
	)

	// Convert ETH fee to BTC
	ethToBTC := ethUSDPrice.Quo(btcUSDPrice)
	feeBTCFloat := new(big.Float).Mul(
		feeInETH,
		new(big.Float).SetFloat64(ethToBTC.MustFloat64()),
	)

	// Convert to satoshis (multiply by 1e8)
	satoshisFloat := new(big.Float).Mul(
		feeBTCFloat,
		new(big.Float).SetInt64(1e8),
	)

	satoshisInt, _ := satoshisFloat.Int(nil)
	satoshis := satoshisInt.Uint64()

	// Convert BTC fee to zenBTC using exchange rate
	feeZenBTC := math.LegacyNewDecFromInt(math.NewIntFromUint64(satoshis)).Quo(exchangeRate).TruncateInt().Uint64()

	return feeZenBTC
}

// clearEthereumNonceRequest resets the nonce-request flag for a given key.
func (k *Keeper) clearEthereumNonceRequest(ctx sdk.Context, keyID uint64) error {
	k.Logger(ctx).Warn("set EthereumNonceRequested state to false", "keyID", keyID)
	return k.EthereumNonceRequested.Set(ctx, keyID, false)
}

// getPendingMintTransactionsByStatus retrieves up to 2 pending mint transactions matching the given status.
func (k *Keeper) getPendingMintTransactionsByStatus(ctx sdk.Context, status zenbtctypes.MintTransactionStatus) ([]zenbtctypes.PendingMintTransaction, error) {
	firstPendingID := uint64(0)
	var err error
	if status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED {
		firstPendingID, err = k.zenBTCKeeper.GetFirstPendingStakeTransaction(ctx)
	} else if status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED {
		firstPendingID, err = k.zenBTCKeeper.GetFirstPendingMintTransaction(ctx)
	}
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		firstPendingID = 0
	}
	return getPendingTransactions(
		ctx,
		k.zenBTCKeeper.GetPendingMintTransactionsStore(),
		func(tx zenbtctypes.PendingMintTransaction) bool {
			return tx.Status == status
		},
		firstPendingID,
		2,
	)
}

// getPendingBurnEvents retrieves up to 2 pending burn events with status BURNED.
func (k *Keeper) getPendingBurnEvents(ctx sdk.Context) ([]zenbtctypes.BurnEvent, error) {
	firstPendingID, err := k.zenBTCKeeper.GetFirstPendingBurnEvent(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		firstPendingID = 0
	}
	return getPendingTransactions(
		ctx,
		k.zenBTCKeeper.GetBurnEventsStore(),
		func(event zenbtctypes.BurnEvent) bool {
			return event.Status == zenbtctypes.BurnStatus_BURN_STATUS_BURNED
		},
		firstPendingID,
		2,
	)
}

// getPendingRedemptions retrieves up to 2 pending redemptions with status INITIATED.
func (k *Keeper) getPendingRedemptions(ctx sdk.Context) ([]zenbtctypes.Redemption, error) {
	firstPendingID, err := k.zenBTCKeeper.GetFirstPendingRedemption(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		firstPendingID = 0
	}
	return getPendingTransactions(
		ctx,
		k.zenBTCKeeper.GetRedemptionsStore(),
		func(r zenbtctypes.Redemption) bool {
			return r.Status == zenbtctypes.RedemptionStatus_INITIATED
		},
		firstPendingID,
		2,
	)
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

func getChainIDForEigen(ctx sdk.Context) uint64 {
	var chainID uint64 = 17000
	if strings.HasPrefix(ctx.ChainID(), "diamond") {
		chainID = 1
	}
	return chainID
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
