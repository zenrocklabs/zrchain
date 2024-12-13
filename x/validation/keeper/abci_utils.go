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

	"cosmossdk.io/collections"
	"cosmossdk.io/core/comet"
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

	"github.com/Zenrock-Foundation/zrchain/v5/sidecar/proto/api"
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
		EthereumRedemptions:        resp.RedemptionsEthereum,
		SolanaRedemptions:          resp.RedemptionsSolana,
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

func deriveAVSContractStateHash(avsDelegations map[string]map[string]*big.Int) ([32]byte, error) {
	return deriveHash(avsDelegations)
}

func deriveRedemptionsHash(redemptions []api.Redemption) ([32]byte, error) {
	return deriveHash(redemptions)
}

func deriveBitcoinHeaderHash(header *sidecar.BTCBlockHeader) ([32]byte, error) {
	return deriveHash(header)
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

func (k *Keeper) lookupEthereumNonce(ctx context.Context) (uint64, error) {
	addr, err := k.getZenBTCMinterAddressEVM(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting ZenBTC minter address: %w", err)
	}

	nonceResp, err := k.sidecarClient.GetLatestEthereumNonceForAccount(ctx, &sidecar.LatestEthereumNonceForAccountRequest{Address: addr})
	if err != nil {
		return 0, fmt.Errorf("error fetching Ethereum nonce: %w", err)
	}

	return nonceResp.Nonce, nil
}

// func (k *Keeper) constructMintTx(ctx context.Context, recipientAddr string, amount, fee, nonce, gasLimit, baseFee, tipCap uint64) ([]byte, []byte, error) {
// 	encodedMintData, err := encodeMintData(common.HexToAddress(recipientAddr), new(big.Int).SetUint64(amount), fee)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	chainID := big.NewInt(17000)
// 	addr := common.HexToAddress(k.GetZenBTCEthContractAddr(ctx))
// 	gasTipCap := new(big.Int).SetUint64(tipCap)
// 	gasFeeCap := new(big.Int).Mul(new(big.Int).SetUint64(baseFee), big.NewInt(2))
// 	gasFeeCap.Add(gasFeeCap, gasTipCap)

// 	unsignedTx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
// 		ChainID:    chainID,
// 		Nonce:      nonce,
// 		GasTipCap:  gasTipCap,
// 		GasFeeCap:  gasFeeCap,
// 		Gas:        gasLimit,
// 		To:         &addr,
// 		Value:      big.NewInt(0), // we shouldn't send any ETH
// 		Data:       encodedMintData,
// 		AccessList: nil,
// 	})

// 	unsignedTxBz, err := unsignedTx.MarshalBinary()
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	signer := ethtypes.LatestSignerForChainID(chainID)

// 	return signer.Hash(unsignedTx).Bytes(), unsignedTxBz, nil
// }

// TODO: use above function instead of this one if possible
func (k *Keeper) constructMintTx(ctx context.Context, recipientAddr string, chainID, amount, fee, nonce, gasLimit, baseFee, tipCap uint64) ([]byte, []byte, error) {
	// if chainID != 17000 && chainID != 11155111 {
	if chainID != 17000 {
		return nil, nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	encodedMintData, err := encodeWrapCallData(common.HexToAddress(recipientAddr), new(big.Int).SetUint64(amount), fee)
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.GetZenBTCEthBatcherAddr(ctx))

	// Convert EIP-1559 fees to legacy gas price
	// gasPrice = baseFee + tipCap
	gasPrice := new(big.Int).Add(
		new(big.Int).SetUint64(baseFee),
		new(big.Int).SetUint64(tipCap),
	)

	// TODO: REMOVE THIS LINE BELOW
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(25))

	unsignedTx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &addr,
		Value:    big.NewInt(0), // we shouldn't send any ETH
		Data:     encodedMintData,
	})

	unsignedTxBz, err := unsignedTx.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	signer := ethtypes.LatestSignerForChainID(new(big.Int).SetUint64(chainID))

	return signer.Hash(unsignedTx).Bytes(), unsignedTxBz, nil
}

func encodeWrapCallData(recipientAddr common.Address, amount *big.Int, fee uint64) ([]byte, error) {
	if !amount.IsUint64() {
		return nil, fmt.Errorf("amount exceeds uint64 max value")
	}

	parsed, err := bindings.ZenbtcbatcherMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %v", err)
	}

	// Pack using the contract binding's ABI for the wrapZenBTC function
	data, err := parsed.Pack(
		"wrapZenBTC",
		recipientAddr,
		amount.Uint64(),
		fee,
		bindings.ISignatureUtilsSignatureWithExpiry{Signature: []byte{}, Expiry: big.NewInt(0)}, [32]byte{},
		// The fields on the line above can be left empty as we don't need them to delegate to our operator
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode wrapZenBTC call data: %v", err)
	}

	return data, nil
}

func (k *Keeper) getZenBTCMinterAddressEVM(ctx context.Context) (string, error) {
	keyID := k.GetZenBTCMinterKeyID(ctx)

	q, err := k.treasuryKeeper.KeyByID(ctx, &treasurytypes.QueryKeyByIDRequest{
		Id:         keyID,
		WalletType: treasurytypes.WalletType_WALLET_TYPE_EVM,
		Prefixes:   make([]string, 0),
	})
	if err != nil {
		return "", err
	}

	return q.Wallets[0].Address, nil
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
		return k.sidecarClient.GetLatestBitcoinBlockHeader(ctx, &sidecar.LatestBitcoinBlockHeaderRequest{ChainName: "testnet4"}) // TODO: use config
	}

	return k.sidecarClient.GetBitcoinBlockHeaderByHeight(ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{ChainName: "testnet4", BlockHeight: requestedBitcoinHeaders.Heights[0]})
}

func (k *Keeper) getNextEthereumNonce(ctx context.Context) (uint64, error) {
	nonce, err := k.lookupEthereumNonce(ctx)
	if err != nil {
		return 0, err
	}

	firstRun := false
	lastUsedNonce, err := k.LastUsedEthereumNonce.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return 0, err
		}
		lastUsedNonce = zenbtctypes.NonceData{Nonce: nonce, Counter: 0}
		firstRun = true
	}

	if !firstRun {
		if nonce == lastUsedNonce.Nonce {
			lastUsedNonce.Counter++
		} else {
			lastUsedNonce.Nonce = nonce
			lastUsedNonce.Counter = 0
		}
	}

	if err = k.LastUsedEthereumNonce.Set(ctx, lastUsedNonce); err != nil {
		return 0, err
	}

	if lastUsedNonce.Counter%8 != 0 { // only retry mint using same nonce every 8 blocks
		return 0, nil
	}

	return nonce, nil
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
