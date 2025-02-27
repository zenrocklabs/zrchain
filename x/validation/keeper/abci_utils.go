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

	// Get sorted list of validators for deterministic iteration
	validators := make([]string, 0, len(delegations))
	for validator := range delegations {
		validators = append(validators, validator)
	}
	slices.Sort(validators)

	for _, validator := range validators {
		delegatorMap := delegations[validator]
		total := new(big.Int)

		// Get sorted list of delegators for deterministic iteration
		delegators := make([]string, 0, len(delegatorMap))
		for delegator := range delegatorMap {
			delegators = append(delegators, delegator)
		}
		slices.Sort(delegators)

		for _, delegator := range delegators {
			amount := delegatorMap[delegator]
			total.Add(total, amount)
		}
		validatorTotals[validator] = total
	}

	// Get sorted list of validators again for deterministic output
	validators = make([]string, 0, len(validatorTotals))
	for validator := range validatorTotals {
		validators = append(validators, validator)
	}
	slices.Sort(validators)

	validatorDelegations := make([]ValidatorDelegations, 0, len(validatorTotals))
	for _, validator := range validators {
		totalStake := validatorTotals[validator]
		validatorDelegations = append(validatorDelegations, ValidatorDelegations{
			Validator: validator,
			Stake:     sdkmath.NewIntFromBigInt(totalStake),
		})
	}

	return validatorDelegations, nil
}

// GetSuperMajorityVEData tallies votes for individual fields of the VoteExtension instead of requiring
// consensus on the entire object. This makes the system more resilient by allowing fields
// that have supermajority consensus to be accepted even if other fields don't reach consensus.
func (k Keeper) GetSuperMajorityVEData(ctx context.Context, currentHeight int64, extCommit abci.ExtendedCommitInfo) (VoteExtension, map[VoteExtensionField]int64, int64, error) {
	// Use a generic map to store votes for all fields
	fieldVotes := make(map[VoteExtensionField]map[string]fieldVote)

	// Initialize maps for each field type
	for i := VEFieldZRChainBlockHeight; i <= VEFieldLatestBtcHeaderHash; i++ {
		fieldVotes[i] = make(map[string]fieldVote)
	}

	var totalVotePower int64
	fieldVotePowers := make(map[VoteExtensionField]int64)

	// Get field handlers
	fieldHandlers := initializeFieldHandlers()
	// Sort handlers by field for deterministic processing
	slices.SortFunc(fieldHandlers, func(a, b FieldHandler) int {
		return int(a.Field) - int(b.Field)
	})

	// Process all votes
	for _, vote := range extCommit.Votes {
		totalVotePower += vote.Validator.Power

		voteExt, err := k.validateVote(ctx, vote, currentHeight)
		if err != nil {
			continue
		}

		// Process each field using the field handlers
		for _, handler := range fieldHandlers {
			key := genericGetKey(handler.GetValue(voteExt))
			value := handler.GetValue(voteExt)

			votes := fieldVotes[handler.Field]
			if existingVote, ok := votes[key]; ok {
				existingVote.votePower += vote.Validator.Power
				votes[key] = existingVote
			} else {
				votes[key] = fieldVote{
					value:     value,
					votePower: vote.Validator.Power,
				}
			}
		}
	}

	// Create consensus VoteExtension with fields that have supermajority
	var consensusVE VoteExtension
	consensusVE.ZRChainBlockHeight = currentHeight - 1

	superMajorityThreshold := superMajorityVotePower(totalVotePower)
	simpleMajorityThreshold := simpleMajorityVotePower(totalVotePower)

	// Apply consensus for each field - ensure deterministic order
	for _, handler := range fieldHandlers {
		votes := fieldVotes[handler.Field]
		if len(votes) == 0 {
			continue
		}

		var mostVotedValue any
		var maxVotePower int64

		// Get keys in deterministic order
		keys := make([]string, 0, len(votes))
		for k := range votes {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		// Find most voted value
		for _, key := range keys {
			vote := votes[key]
			if vote.votePower > maxVotePower {
				maxVotePower = vote.votePower
				mostVotedValue = vote.value
			}
		}

		// Use simple majority for less critical i.e. gas-related fields, supermajority for others
		requiredPower := superMajorityThreshold
		if isGasField(handler.Field) {
			requiredPower = simpleMajorityThreshold
		}

		if maxVotePower >= requiredPower {
			handler.SetValue(mostVotedValue, &consensusVE)
			fieldVotePowers[handler.Field] = maxVotePower
		}
	}

	// Log consensus results
	k.logConsensusResults(ctx, fieldVotePowers, superMajorityThreshold, simpleMajorityThreshold)

	return consensusVE, fieldVotePowers, totalVotePower, nil
}

// logConsensusResults logs information about which fields reached consensus
func (k Keeper) logConsensusResults(ctx context.Context, fieldVotePowers map[VoteExtensionField]int64, superMajorityThreshold, simpleMajorityThreshold int64) {
	if len(fieldVotePowers) == 0 {
		k.Logger(ctx).Warn("no consensus reached on any vote extension fields")
		return
	}

	totalVotePower := int64(0)
	for _, votePower := range fieldVotePowers {
		if votePower > totalVotePower {
			totalVotePower = votePower
		}
	}

	// Log the count of fields with consensus
	k.Logger(ctx).Info("consensus summary", "fields_with_consensus", len(fieldVotePowers))

	// Use a deterministic field order for logging
	fields := make([]VoteExtensionField, 0, int(VEFieldLatestBtcHeaderHash)+1)
	for field := VEFieldZRChainBlockHeight; field <= VEFieldLatestBtcHeaderHash; field++ {
		fields = append(fields, field)
	}
	slices.Sort(fields)

	// Loop through all possible fields and log their consensus status
	for _, field := range fields {
		// Skip logging the ZRChainBlockHeight field
		if field == VEFieldZRChainBlockHeight {
			continue
		}

		_, hasConsensus := fieldVotePowers[field]
		k.Logger(ctx).Info("field consensus status",
			"field", field.String(),
			"has_consensus", hasConsensus,
			"vote_power", func() int64 {
				if power, ok := fieldVotePowers[field]; ok {
					return power
				}
				return 0
			}(),
			"required_power", func() int64 {
				if isGasField(field) {
					return simpleMajorityThreshold
				}
				return superMajorityThreshold
			}(),
			"threshold", func() string {
				if isGasField(field) {
					return "simple_majority"
				}
				return "supermajority"
			}(),
		)
	}
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

// superMajorityVotePower calculates the required vote power for a supermajority (2/3+)
func superMajorityVotePower(totalVotePower int64) int64 {
	return ((totalVotePower * 2) / 3) + 1
}

// simpleMajorityVotePower calculates the required vote power for a simple majority (>50%)
func simpleMajorityVotePower(totalVotePower int64) int64 {
	return (totalVotePower / 2) + 1
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
		requiredVotePower := superMajorityVotePower(totalVotePower)
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

func (k *Keeper) constructEthereumTx(ctx context.Context, addr common.Address, chainID uint64, data []byte, nonce, gasLimit, baseFee, tipCap uint64) ([]byte, []byte, error) {
	validatedChainID, err := types.ValidateChainID(ctx, chainID)
	if err != nil {
		return nil, nil, err
	}
	chainIDBigInt := new(big.Int).SetUint64(validatedChainID)

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

	addr := common.HexToAddress(k.zenBTCKeeper.GetControllerAddr(ctx))
	return k.constructEthereumTx(ctx, addr, chainID, encodedMintData, nonce, gasLimit, baseFee, tipCap)
}

func (k *Keeper) constructMintTx(ctx context.Context, recipientAddr string, chainID, amount, fee, nonce, gasLimit, baseFee, tipCap uint64) ([]byte, []byte, error) {
	encodedMintData, err := EncodeWrapCallData(common.HexToAddress(recipientAddr), new(big.Int).SetUint64(amount), fee)
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.zenBTCKeeper.GetEthTokenAddr(ctx))
	return k.constructEthereumTx(ctx, addr, chainID, encodedMintData, nonce, gasLimit, baseFee, tipCap)
}

func (k *Keeper) constructUnstakeTx(ctx context.Context, chainID uint64, destinationAddr []byte, amount, ethNonce, baseFee, tipCap uint64) ([]byte, []byte, error) {
	encodedUnstakeData, err := k.EncodeUnstakeCallData(ctx, destinationAddr, amount)
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.zenBTCKeeper.GetControllerAddr(ctx))
	return k.constructEthereumTx(ctx, addr, chainID, encodedUnstakeData, ethNonce, 700000, baseFee, tipCap)
}

func (k *Keeper) constructCompleteTx(ctx context.Context, chainID, redemptionID, ethNonce, baseFee, tipCap uint64) ([]byte, []byte, error) {
	encodedCompleteData, err := k.EncodeCompleteCallData(ctx, redemptionID)
	if err != nil {
		return nil, nil, err
	}

	addr := common.HexToAddress(k.zenBTCKeeper.GetControllerAddr(ctx))
	return k.constructEthereumTx(ctx, addr, chainID, encodedCompleteData, ethNonce, 300000, baseFee, tipCap)
}

func EncodeStakeCallData(amount *big.Int) ([]byte, error) {
	if !amount.IsUint64() {
		return nil, fmt.Errorf("amount exceeds uint64 max value")
	}

	parsed, err := bindings.ZenBTControllerMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %v", err)
	}

	// Pack using the contract binding's ABI for the stakeRockBTC function
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

func (k *Keeper) retrieveBitcoinHeaders(ctx context.Context) (*sidecar.BitcoinBlockHeaderResponse, *sidecar.BitcoinBlockHeaderResponse, error) {
	// Always get the latest Bitcoin header
	latest, err := k.sidecarClient.GetLatestBitcoinBlockHeader(ctx, &sidecar.LatestBitcoinBlockHeaderRequest{
		ChainName: k.bitcoinNetwork(ctx),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get latest Bitcoin header: %w", err)
	}

	// Check if there are requested historical headers
	requestedBitcoinHeaders, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, nil, err
		}
		requestedBitcoinHeaders = zenbtctypes.RequestedBitcoinHeaders{}
		if err = k.RequestedHistoricalBitcoinHeaders.Set(ctx, requestedBitcoinHeaders); err != nil {
			return nil, nil, err
		}
	}

	// Get requested historical headers if any
	if len(requestedBitcoinHeaders.Heights) > 0 {
		requested, err := k.sidecarClient.GetBitcoinBlockHeaderByHeight(ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{
			ChainName:   k.bitcoinNetwork(ctx),
			BlockHeight: requestedBitcoinHeaders.Heights[0],
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get requested Bitcoin header at height %d: %w", requestedBitcoinHeaders.Heights[0], err)
		}
		return latest, requested, nil
	}

	return latest, nil, nil
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
	pricesAreValid := true
	if oracleData.ROCKUSDPrice.IsZero() || oracleData.BTCUSDPrice.IsZero() || oracleData.ETHUSDPrice.IsZero() {
		pricesAreValid = false
	}

	if pricesAreValid {
		if err := k.LastValidVEHeight.Set(ctx, ctx.BlockHeight()); err != nil {
			k.Logger(ctx).Error("error setting last valid VE height", "height", ctx.BlockHeight(), "err", err)
		}
	} else {
		lastValidVEHeight, err := k.LastValidVEHeight.Get(ctx)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error("error getting last valid VE height", "height", ctx.BlockHeight(), "err", err)
			}
			lastValidVEHeight = 0
		}

		retentionRange := k.GetPriceRetentionBlockRange(ctx)
		// Add safety check for when block height is less than retention range
		if ctx.BlockHeight() < retentionRange {
			k.Logger(ctx).Warn("current block height is less than retention range; not zeroing asset prices",
				"block_height", ctx.BlockHeight(),
				"retention_range", retentionRange)
			return
		}
		// Calculate number of blocks since last valid VE
		if ctx.BlockHeight()-lastValidVEHeight < retentionRange {
			k.Logger(ctx).Warn("last valid VE height is within price retention range; not zeroing asset prices",
				"retention_range", retentionRange)
			return
		}
	}

	// Ensure deterministic order of asset price updates
	assets := []types.Asset{types.Asset_ROCK, types.Asset_BTC, types.Asset_ETH}
	assetPrices := []math.LegacyDec{oracleData.ROCKUSDPrice, oracleData.BTCUSDPrice, oracleData.ETHUSDPrice}

	for i, asset := range assets {
		if err := k.AssetPrices.Set(ctx, asset, assetPrices[i]); err != nil {
			k.Logger(ctx).Error("error setting asset price", "asset", asset.String(), "height", ctx.BlockHeight(), "err", err)
		}
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
// Only fields that have reached consensus (present in fieldVotePowers) are validated.
// Fields that fail validation are removed from fieldVotePowers to prevent them from being used downstream.
func (k *Keeper) validateOracleData(ctx context.Context, voteExt VoteExtension, oracleData *OracleData, fieldVotePowers map[VoteExtensionField]int64) error {
	invalidFields := make([]VoteExtensionField, 0)

	// Get sorted list of fields for deterministic validation
	fields := make([]VoteExtensionField, 0, len(fieldVotePowers))
	for field := range fieldVotePowers {
		fields = append(fields, field)
	}
	slices.Sort(fields)

	// Validate hashes only if fields have consensus - using deterministic order
	for _, field := range fields {
		switch field {
		case VEFieldEigenDelegationsHash:
			if err := validateHashField(field.String(), voteExt.EigenDelegationsHash, oracleData.EigenDelegationsMap); err != nil {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldEthBurnEventsHash:
			if err := validateHashField(field.String(), voteExt.EthBurnEventsHash, oracleData.EthBurnEvents); err != nil {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldRedemptionsHash:
			if err := validateHashField(field.String(), voteExt.RedemptionsHash, oracleData.Redemptions); err != nil {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldRequestedBtcHeaderHash:
			if oracleData.RequestedBtcBlockHeight != 0 {
				if err := validateHashField(field.String(), voteExt.RequestedBtcHeaderHash, &oracleData.RequestedBtcBlockHeader); err != nil {
					invalidFields = append(invalidFields, field)
				}
			}
		case VEFieldEthBlockHeight:
			if voteExt.EthBlockHeight != oracleData.EthBlockHeight {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldEthGasLimit:
			if voteExt.EthGasLimit != oracleData.EthGasLimit {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldEthBaseFee:
			if voteExt.EthBaseFee != oracleData.EthBaseFee {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldEthTipCap:
			if voteExt.EthTipCap != oracleData.EthTipCap {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldRequestedBtcBlockHeight:
			if voteExt.RequestedBtcBlockHeight != oracleData.RequestedBtcBlockHeight {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldRequestedStakerNonce:
			if voteExt.RequestedStakerNonce != oracleData.RequestedStakerNonce {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldRequestedEthMinterNonce:
			if voteExt.RequestedEthMinterNonce != oracleData.RequestedEthMinterNonce {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldRequestedUnstakerNonce:
			if voteExt.RequestedUnstakerNonce != oracleData.RequestedUnstakerNonce {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldRequestedCompleterNonce:
			if voteExt.RequestedCompleterNonce != oracleData.RequestedCompleterNonce {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldROCKUSDPrice:
			if !voteExt.ROCKUSDPrice.Equal(oracleData.ROCKUSDPrice) {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldBTCUSDPrice:
			if !voteExt.BTCUSDPrice.Equal(oracleData.BTCUSDPrice) {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldETHUSDPrice:
			if !voteExt.ETHUSDPrice.Equal(oracleData.ETHUSDPrice) {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldLatestBtcBlockHeight:
			if voteExt.LatestBtcBlockHeight != oracleData.LatestBtcBlockHeight {
				invalidFields = append(invalidFields, field)
			}
		case VEFieldLatestBtcHeaderHash:
			if err := validateHashField(field.String(), voteExt.LatestBtcHeaderHash, &oracleData.LatestBtcBlockHeader); err != nil {
				invalidFields = append(invalidFields, field)
			}
		}
	}

	// Sort invalidFields for deterministic ordering
	slices.Sort(invalidFields)

	// Remove invalid fields from fieldVotePowers - already in sorted order
	for _, field := range invalidFields {
		delete(fieldVotePowers, field)
		k.Logger(ctx).Info("removed mismatched field from fieldVotePowers", "field", field.String())
	}

	// Update FieldVotePowers in oracleData to reflect the validated fields
	oracleData.FieldVotePowers = fieldVotePowers

	return nil
}

// Helper function to validate consensus on multiple required fields for transactions
func (k *Keeper) validateConsensusForTxFields(ctx sdk.Context, oracleData OracleData, requiredFields []VoteExtensionField, txType, txDetails string) error {
	// Always check for gas fields consensus first
	if !HasRequiredGasFields(oracleData.FieldVotePowers) {
		k.Logger(ctx).Error(fmt.Sprintf("cannot process %s: missing consensus on gas fields", txType),
			"details", txDetails)
		return fmt.Errorf("missing consensus on gas fields required for transaction construction")
	}

	// Check if all required fields have consensus
	missingFields := allFieldsHaveConsensus(oracleData.FieldVotePowers, requiredFields)
	if len(missingFields) > 0 {
		// Sort missing fields for deterministic error messages
		slices.Sort(missingFields)

		fieldNames := make([]string, 0, len(missingFields))
		for _, field := range missingFields {
			fieldNames = append(fieldNames, field.String())
		}
		k.Logger(ctx).Error(fmt.Sprintf("cannot process %s: missing consensus on fields: %s", txType, strings.Join(fieldNames, ", ")),
			"details", txDetails)
		return fmt.Errorf("missing consensus on fields required for transaction construction: %s", strings.Join(fieldNames, ", "))
	}

	return nil
}

// fieldsHaveConsensus checks if all specified fields have consensus and returns any fields that don't
func allFieldsHaveConsensus(fieldVotePowers map[VoteExtensionField]int64, fields []VoteExtensionField) []VoteExtensionField {
	var missingConsensus []VoteExtensionField

	// Sort fields for deterministic checking
	sortedFields := make([]VoteExtensionField, len(fields))
	copy(sortedFields, fields)
	slices.Sort(sortedFields)

	for _, field := range sortedFields {
		if !fieldHasConsensus(fieldVotePowers, field) {
			missingConsensus = append(missingConsensus, field)
		}
	}
	return missingConsensus
}

// anyFieldHasConsensus checks if at least one of the specified fields has consensus
func anyFieldHasConsensus(fieldVotePowers map[VoteExtensionField]int64, fields []VoteExtensionField) bool {
	// Sort fields for deterministic checking
	sortedFields := make([]VoteExtensionField, len(fields))
	copy(sortedFields, fields)
	slices.Sort(sortedFields)

	for _, field := range sortedFields {
		if fieldHasConsensus(fieldVotePowers, field) {
			return true
		}
	}
	return false
}
