package keeper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"slices"
	"sort"
	"strings"

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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sidecar "github.com/Zenrock-Foundation/zrchain/v4/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v4/x/validation/types"
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
	if err := json.Unmarshal(resp.Delegations, &delegations); err != nil {
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

	ETHUSDPrice, err := sdkmath.LegacyNewDecFromStr(resp.ETHUSDPrice)
	if err != nil {
		k.Logger(ctx).Error("error parsing eth price", "error", err)
		return nil, ErrOracleSidecar
	}

	return &OracleData{
		ROCKUSDPrice:         ROCKUSDPrice,
		ETHUSDPrice:          ETHUSDPrice,
		AVSDelegationsMap:    delegations,
		ValidatorDelegations: validatorDelegations,
		EthBlockHeight:       resp.EthBlockHeight,
		EthBlockHash:         common.HexToHash(resp.EthBlockHash),
		EthGasPrice:          resp.EthGasPrice,
		EthGasLimit:          resp.EthGasLimit,
		ConsensusData:        abci.ExtendedCommitInfo{},
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

func deriveAVSContractStateHash(avsDelegations map[string]map[string]*big.Int) ([32]byte, error) {
	sortDelegationsMap(avsDelegations) // Sort AVS delegations for determinism

	avsDelegationsBz, err := json.Marshal(avsDelegations)
	if err != nil {
		return [32]byte{}, fmt.Errorf("error encoding AVS delegations: %w", err)
	}

	return sha256.Sum256(avsDelegationsBz), nil
}

func sortDelegationsMap(m map[string]map[string]*big.Int) {
	outerKeys := make([]string, 0, len(m))
	for k := range m {
		outerKeys = append(outerKeys, k)
	}
	sort.Strings(outerKeys)

	for _, outerKey := range outerKeys {
		innerMap := m[outerKey]

		innerKeys := make([]string, 0, len(innerMap))
		for k := range innerMap {
			innerKeys = append(innerKeys, k)
		}
		sort.Strings(innerKeys)
	}
}

func (k Keeper) GetSuperMajorityVE(ctx context.Context, currentHeight int64, extCommit abci.ExtendedCommitInfo) (VoteExtension, error) {
	votesPerVoteExt := make(map[string]*VEWithVotePower)
	var totalVotePower int64

	for _, vote := range extCommit.Votes {
		totalVotePower += vote.Validator.Power

		voteExt, err := validateVote(vote, currentHeight)
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

func validateVote(vote abci.ExtendedVoteInfo, currentHeight int64) (VoteExtension, error) {
	if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit || len(vote.VoteExtension) == 0 {
		return VoteExtension{}, fmt.Errorf("invalid vote")
	}

	var voteExt VoteExtension
	if err := json.Unmarshal(vote.VoteExtension, &voteExt); err != nil {
		return VoteExtension{}, err
	}

	voteExt.ZRChainBlockHeight = currentHeight - 1
	if voteExt.IsInvalid() {
		return VoteExtension{}, fmt.Errorf("invalid vote extension")
	}

	return voteExt, nil
}

func getVESubset(ve VoteExtension) VoteExtension {
	ve.EthBlockHeight = 0
	ve.EthBlockHash = [32]byte{}
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
			VoteExtension: fullMarshaledVE, // Store the full VE here
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

func (k *Keeper) nextRequestedEthereumNonceHeight(ctx context.Context) (uint64, error) {
	var requestedHeight uint64 = 0
	if err := k.RequestedEthereumNonceHeights.Walk(ctx, nil, func(height uint64) (bool, error) {
		requestedHeight = height
		return true, nil
	}); err != nil {
		k.Logger(ctx).Error("error walking through RequestedEthereumNonceHeights", "err", err)
		return 0, err
	}
	return requestedHeight, nil
}

func (k *Keeper) lookupEthereumNonce(ctx context.Context) (uint64, error) {
	nonceHeight, err := k.nextRequestedEthereumNonceHeight(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting next requested Ethereum nonce height: %w", err)
	}

	if nonceHeight == 0 {
		return 0, nil
	}

	addr, err := k.getZenBTCMinterAddressEVM(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting ZenBTC minter address: %w", err)
	}

	nonceResp, err := k.sidecarClient.GetEthereumNonceAtHeight(ctx, &sidecar.EthereumNonceAtHeightRequest{
		Address: addr,
		Height:  nonceHeight,
	})
	if err != nil {
		return 0, fmt.Errorf("error fetching Ethereum nonce: %w", err)
	}

	return nonceResp.Nonce, nil
}

func (k *Keeper) lookupNextUnlockTx(ctx context.Context) (string, uint64, error) {
	chain, txID, withdrawalInfo, totalTxs, err := k.nextUnconfirmedUnlockTx(ctx)
	if err != nil {
		return "", 0, err
	}

	if totalTxs == 0 || chain == "" || txID == "" {
		return "", 0, nil
	}

	return k.querySidecarForUnlockTx(ctx, chain, txID, withdrawalInfo)
}

func (k *Keeper) nextUnconfirmedUnlockTx(ctx context.Context) (string, string, types.WithdrawalInfo, uint64, error) {
	var lowestRetryCount uint32 = math.MaxUint32
	chain := ""
	txID := ""
	withdrawalInfo := types.WithdrawalInfo{}
	var totalTxs uint64 = 0
	if err := k.UnconfirmedUnlockTxs.Walk(ctx, nil, func(key collections.Pair[string, string], withdrawInfo types.WithdrawalInfo) (bool, error) {
		totalTxs++
		if withdrawalInfo.RetryCount < lowestRetryCount {
			lowestRetryCount = withdrawalInfo.RetryCount
			chain = key.K1()
			txID = key.K2()
			withdrawalInfo = withdrawInfo
		}
		return false, nil
	}); err != nil {
		k.Logger(ctx).Error("error walking through UnconfirmedSolanaUnlockTxs", "err", err)
		return "", "", types.WithdrawalInfo{}, 0, err
	}
	return chain, txID, withdrawalInfo, totalTxs, nil
}

func (k *Keeper) querySidecarForUnlockTx(ctx context.Context, chain, txID string, withdrawalInfo types.WithdrawalInfo) (string, uint64, error) {
	txHeight, err := k.getTransactionHeight(ctx, chain, txID)
	if err != nil {
		k.handleUnlockTxError(ctx, chain, txID, withdrawalInfo, err)
	}
	return chain, txHeight, err
}

func (k *Keeper) getTransactionHeight(ctx context.Context, chain, txID string) (uint64, error) {
	switch chain {
	case "eth":
		resp, err := k.sidecarClient.GetEthereumTransaction(ctx, &sidecar.EthereumTransactionRequest{TxHash: txID})
		if err != nil {
			return 0, err
		}
		return resp.TxHeight, nil
	case "sol":
		resp, err := k.sidecarClient.GetSolanaTransaction(ctx, &sidecar.SolanaTransactionRequest{TxSignature: txID})
		if err != nil {
			return 0, err
		}
		return resp.TxSlot, nil
	default:
		return 0, fmt.Errorf("unsupported chain: %s", chain)
	}
}

func (k *Keeper) handleUnlockTxError(ctx context.Context, chain, txID string, withdrawalInfo types.WithdrawalInfo, err error) {
	k.Logger(ctx).Warn("error retrieving unlock tx", "error", err)
	withdrawalInfo.RetryCount++
	key := collections.Join(chain, txID)
	if withdrawalInfo.RetryCount >= 25 {
		k.Logger(ctx).Warn("unlock tx retry count exceeded; removing from store", "txSignature", txID)
		if removeErr := k.UnconfirmedUnlockTxs.Remove(ctx, key); removeErr != nil {
			k.Logger(ctx).Error("error removing unlock tx", "txID", txID, "error", removeErr)
		}
	} else {
		k.UnconfirmedUnlockTxs.Set(ctx, key, withdrawalInfo)
	}
}

func (k *Keeper) constructMintTx(ctx context.Context, recipientAddr string, amount *big.Int, fee uint64, nonce, gasPrice, gasLimit uint64) ([]byte, error) {
	encodedMintData, err := encodeMintData(common.HexToAddress(recipientAddr), amount, fee)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTransaction(
		nonce,
		common.HexToAddress(k.GetZenBTCEthContractAddr(ctx)),
		big.NewInt(0), // we shouldn't send any ETH
		gasLimit,
		big.NewInt(int64(gasPrice)),
		encodedMintData,
	)
	// unsignedRlpTx, err := tx.MarshalBinary()
	chainId := big.NewInt(17000)
	signer := ethtypes.LatestSignerForChainID(chainId)
	// after signature
	// signedTx,err := tx.WithSignature(signer, signature) TODO: Do we need to do this?
	// err := client.Broadcast(signedTx)
	return signer.Hash(tx).Bytes(), nil
}

func encodeMintData(recipientAddr common.Address, amount *big.Int, fee uint64) ([]byte, error) {
	const mintFunctionABI = `[{"name":"mint","type":"function","inputs":[{"type":"address","name":"account"},{"type":"uint256","name":"value"},{"type":"uint256","name":"fee"}],"outputs":[]}]`

	parsedABI, err := abi.JSON(strings.NewReader(mintFunctionABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}
	feeAmount := new(big.Int).SetUint64(fee)

	data, err := parsedABI.Pack("mint", recipientAddr, amount, feeAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode mint function data: %v", err)
	}
	return data, nil
}

func (k *Keeper) getZenBTCMinterAddressEVM(ctx context.Context) (string, error) {
	q, err := k.treasuryKeeper.KeyByID(ctx, &treasurytypes.QueryKeyByIDRequest{
		Id:         k.GetZenBTCMinterKeyID(ctx),
		WalletType: treasurytypes.WalletType_WALLET_TYPE_EVM,
	})
	if err != nil {
		return "", err
	}
	return q.Wallets[0].Address, nil
}
