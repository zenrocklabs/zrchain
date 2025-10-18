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
	"reflect"
	"slices"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/comet"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock/generated/rock_spl_token"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solzenbtc/generated/zenbtc_spl_token"
	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	zenbtctypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cryptoenc "github.com/cometbft/cometbft/crypto/encoding"
	"github.com/cometbft/cometbft/libs/protoio"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/zenrocklabs/goem/ethereum"

	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	bindings "github.com/Zenrock-Foundation/zrchain/v6/zenbtc/bindings"
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

	return &OracleData{
		EigenDelegationsMap:  delegations,
		ValidatorDelegations: validatorDelegations,
		EthBlockHeight:       resp.EthBlockHeight,
		EthGasLimit:          resp.EthGasLimit,
		EthBaseFee:           resp.EthBaseFee,
		EthTipCap:            resp.EthTipCap,
		EthBurnEvents:        resp.EthBurnEvents,
		Redemptions:          resp.Redemptions,
		ROCKUSDPrice:         resp.ROCKUSDPrice,
		BTCUSDPrice:          resp.BTCUSDPrice,
		ETHUSDPrice:          resp.ETHUSDPrice,
		ConsensusData:        abci.ExtendedCommitInfo{},
		SolanaBurnEvents:     resp.SolanaBurnEvents,
		SolanaMintEvents:     resp.SolanaMintEvents,
		SidecarVersionName:   resp.SidecarVersionName,
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
			Stake:     math.NewIntFromBigInt(totalStake),
		})
	}

	return validatorDelegations, nil
}

// GetConsensusAndPluralityVEData tallies votes for individual fields of the VoteExtension instead of requiring
// consensus on the entire object. This makes the system more resilient by allowing fields
// that have supermajority consensus to be accepted even if other fields don't reach consensus.
// It returns both a consensus vote extension (fields with supermajority/simple majority) and
// a plurality vote extension (fields with the most votes regardless of threshold).
// A deterministic tie-breaking mechanism is used based on the lexicographic ordering of the string representation
// of the values. This ensures all validators will select the same consensus value regardless
// of iteration order.
func (k Keeper) GetConsensusAndPluralityVEData(ctx context.Context, currentHeight int64, extCommit abci.ExtendedCommitInfo) (VoteExtension, VoteExtension, map[VoteExtensionField]int64, error) {
	// Get field handlers first
	fieldHandlers := initializeFieldHandlers()

	// Use a generic map to store votes for all fields
	fieldVotes := make(map[VoteExtensionField]map[string]fieldVote)

	// Initialize maps for each field handler to ensure all fields have a map
	for _, handler := range fieldHandlers {
		fieldVotes[handler.Field] = make(map[string]fieldVote)
	}

	var totalVotePower int64
	fieldVotePowers := make(map[VoteExtensionField]int64)

	// Track informational fields that don't require consensus
	var firstSidecarVersionName string

	// Process all votes
	for _, vote := range extCommit.Votes {
		totalVotePower += vote.Validator.Power

		voteExt, err := k.validateVote(ctx, vote, currentHeight)
		if err != nil {
			continue
		}

		// Capture informational fields from first valid vote (no consensus required)
		if firstSidecarVersionName == "" && voteExt.SidecarVersionName != "" {
			firstSidecarVersionName = voteExt.SidecarVersionName
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
	// Create plurality VoteExtension with fields that have the most votes (regardless of threshold)
	var pluralityVE VoteExtension

	// Set informational fields that don't require consensus
	consensusVE.SidecarVersionName = firstSidecarVersionName
	pluralityVE.SidecarVersionName = firstSidecarVersionName

	superMajorityThreshold := superMajorityVotePower(totalVotePower)
	simpleMajorityThreshold := simpleMajorityVotePower(totalVotePower)

	// Apply consensus for each field
	for _, handler := range fieldHandlers {
		votes := fieldVotes[handler.Field]
		if len(votes) == 0 {
			continue
		}

		// Find the maximum vote power for this field
		var maxVotePower int64
		for _, vote := range votes {
			if vote.votePower > maxVotePower {
				maxVotePower = vote.votePower
			}
		}

		// Collect all values that have the maximum vote power (for both consensus and plurality)
		var tiedValues []struct {
			key   string
			value any
		}

		for key, vote := range votes {
			if vote.votePower == maxVotePower {
				tiedValues = append(tiedValues, struct {
					key   string
					value any
				}{key, vote.value})
			}
		}

		// If there are multiple values with the same vote power, use deterministic tie-breaking
		if len(tiedValues) > 1 {
			// Log the tie-breaking event with field information
			k.Logger(ctx).Info("performing deterministic tie-breaking",
				"field", handler.Field.String(),
				"tied_values_count", len(tiedValues),
				"max_vote_power", maxVotePower)

			// Sort by the hash of their serialized representation for deterministic selection
			slices.SortFunc(tiedValues, func(a, b struct {
				key   string
				value any
			}) int {
				// Use the key, which is already a deterministic string representation
				return strings.Compare(a.key, b.key)
			})
		}

		// Always select the first value after deterministic sorting for plurality
		mostVotedValue := tiedValues[0].value
		handler.SetValue(mostVotedValue, &pluralityVE)

		// Use simple majority for gas-related fields, supermajority for others
		requiredPower := superMajorityThreshold
		if isGasField(handler.Field) {
			requiredPower = simpleMajorityThreshold
		}

		// Check if any value has sufficient votes for consensus
		if maxVotePower >= requiredPower {
			handler.SetValue(mostVotedValue, &consensusVE)
			fieldVotePowers[handler.Field] = maxVotePower
		}
	}

	// Log consensus results
	k.logConsensusResults(ctx, fieldVotePowers)

	return consensusVE, pluralityVE, fieldVotePowers, nil
}

// logConsensusResults logs information about which fields reached consensus
func (k Keeper) logConsensusResults(ctx context.Context, fieldVotePowers map[VoteExtensionField]int64) {
	if len(fieldVotePowers) == 0 {
		k.Logger(ctx).Error("no consensus reached on any vote extension fields")
		return
	}

	totalVotePower := int64(0)
	for _, votePower := range fieldVotePowers {
		if votePower > totalVotePower {
			totalVotePower = votePower
		}
	}

	// Log the count of fields with consensus
	k.Logger(ctx).Info("consensus summary", "fields_with_consensus", len(fieldVotePowers), "stage", "initial")

	// Collect fields with and without consensus
	fieldsWithConsensus := make([]string, 0)
	fieldsWithoutConsensus := make([]string, 0)

	for field := VEFieldEigenDelegationsHash; field <= VEFieldLatestZcashHeaderHash; field++ {
		_, hasConsensus := fieldVotePowers[field]
		if hasConsensus {
			fieldsWithConsensus = append(fieldsWithConsensus, field.String())
		} else {
			fieldsWithoutConsensus = append(fieldsWithoutConsensus, field.String())
		}
	}

	// Log consolidated field status
	if len(fieldsWithConsensus) > 0 {
		k.Logger(ctx).Debug("fields with consensus",
			"fields", strings.Join(fieldsWithConsensus, ","),
			"stage", "initial")
	}

	if len(fieldsWithoutConsensus) > 0 {
		k.Logger(ctx).Warn("fields without consensus",
			"fields", strings.Join(fieldsWithoutConsensus, ","),
			"stage", "initial")
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
	caip2ChainID := fmt.Sprintf("eip155:%d", chainID)
	if _, err := types.ValidateEVMChainID(ctx, caip2ChainID); err != nil {
		return nil, nil, err
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
	address, err := k.treasuryKeeper.GetAddressByWalletType(sdk.UnwrapSDKContext(ctx), keyID, walletType, make([]string, 0))
	if err != nil {
		return "", fmt.Errorf("error getting address for key ID %d: %w", keyID, err)
	}

	return address, nil
}

func (k *Keeper) bitcoinNetwork(ctx context.Context) string {
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "diamond") {
		return "mainnet"
	}
	// This is the chainID needed in zrchain so it uses the bitcoin regnet node
	if strings.HasPrefix(sdk.UnwrapSDKContext(ctx).ChainID(), "docker") {
		return "regnet"
	}
	return "testnet4"
}

func (k *Keeper) retrieveBitcoinHeaders(ctx context.Context, requestedBtcHeaderHeight int64) (*sidecar.BitcoinBlockHeaderResponse, *sidecar.BitcoinBlockHeaderResponse, error) {
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

	// If a consensus-selected requested height is provided, fetch exactly that
	if requestedBtcHeaderHeight > 0 {
		requested, err := k.sidecarClient.GetBitcoinBlockHeaderByHeight(ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{
			ChainName:   k.bitcoinNetwork(ctx),
			BlockHeight: requestedBtcHeaderHeight,
		})
		if err != nil {
			return latest, nil, fmt.Errorf("failed to get requested Bitcoin header at height %d: %w", requestedBtcHeaderHeight, err)
		}
		return latest, requested, nil
	}

	// Otherwise fall back to the first pending requested height if any
	if len(requestedBitcoinHeaders.Heights) > 0 {
		requested, err := k.sidecarClient.GetBitcoinBlockHeaderByHeight(ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{
			ChainName:   k.bitcoinNetwork(ctx),
			BlockHeight: requestedBitcoinHeaders.Heights[0],
		})
		if err != nil {
			return latest, nil, fmt.Errorf("failed to get requested Bitcoin header at height %d: %w", requestedBitcoinHeaders.Heights[0], err)
		}
		return latest, requested, nil
	}

	return latest, nil, nil
}

// retrieveZcashHeaders retrieves the latest and optionally a requested ZCash block header from the sidecar
func (k *Keeper) retrieveZcashHeaders(ctx context.Context, requestedZcashHeaderHeight int64) (*sidecar.BitcoinBlockHeaderResponse, *sidecar.BitcoinBlockHeaderResponse, error) {
	// Always get the latest ZCash header
	latest, err := k.sidecarClient.GetLatestZcashBlockHeader(ctx, &sidecar.LatestBitcoinBlockHeaderRequest{
		ChainName: "", // ZCash doesn't use chain name like Bitcoin (no regnet/testnet/mainnet distinction in same way)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get latest ZCash header: %w", err)
	}

	// Check if there are requested historical headers
	requestedZcashHeaders, err := k.RequestedHistoricalZcashHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, nil, err
		}
		requestedZcashHeaders = dcttypes.RequestedZcashHeaders{}
		if err = k.RequestedHistoricalZcashHeaders.Set(ctx, requestedZcashHeaders); err != nil {
			return nil, nil, err
		}
	}

	// If a consensus-selected requested height is provided, fetch exactly that
	if requestedZcashHeaderHeight > 0 {
		requested, err := k.sidecarClient.GetZcashBlockHeaderByHeight(ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{
			ChainName:   "",
			BlockHeight: requestedZcashHeaderHeight,
		})
		if err != nil {
			return latest, nil, fmt.Errorf("failed to get requested ZCash header at height %d: %w", requestedZcashHeaderHeight, err)
		}
		return latest, requested, nil
	}

	// Otherwise fall back to the first pending requested height if any
	if len(requestedZcashHeaders.Heights) > 0 {
		requested, err := k.sidecarClient.GetZcashBlockHeaderByHeight(ctx, &sidecar.BitcoinBlockHeaderByHeightRequest{
			ChainName:   "",
			BlockHeight: requestedZcashHeaders.Heights[0],
		})
		if err != nil {
			return latest, nil, fmt.Errorf("failed to get requested ZCash header at height %d: %w", requestedZcashHeaders.Heights[0], err)
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
	// Parse prices once at the beginning
	rockPrice, rockPriceErr := math.LegacyNewDecFromStr(oracleData.ROCKUSDPrice)
	btcPrice, btcPriceErr := math.LegacyNewDecFromStr(oracleData.BTCUSDPrice)
	ethPrice, ethPriceErr := math.LegacyNewDecFromStr(oracleData.ETHUSDPrice)

	// Check if prices are valid
	pricesAreValid := true
	if oracleData.ROCKUSDPrice == "" || oracleData.BTCUSDPrice == "" || oracleData.ETHUSDPrice == "" ||
		rockPriceErr != nil || btcPriceErr != nil || ethPriceErr != nil ||
		rockPrice.IsNil() || rockPrice.IsZero() ||
		btcPrice.IsNil() || btcPrice.IsZero() ||
		ethPrice.IsNil() || ethPrice.IsZero() {
		pricesAreValid = false
	}

	// Update the last valid VE height if prices are valid
	if pricesAreValid {
		if err := k.LastValidVEHeight.Set(ctx, ctx.BlockHeight()); err != nil {
			k.Logger(ctx).Error("error setting last valid VE height", "height", ctx.BlockHeight(), "err", err)
		}
	} else {
		// Handle invalid prices
		lastValidVEHeight, err := k.LastValidVEHeight.Get(ctx)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error("error getting last valid VE height", "height", ctx.BlockHeight(), "err", err)
			}
			lastValidVEHeight = 0
		}

		retentionRange := k.GetPriceRetentionBlockRange(ctx)

		// Safety check for when block height is less than retention range
		if ctx.BlockHeight() < retentionRange {
			k.Logger(ctx).Warn("current block height is less than retention range; not zeroing asset prices",
				"block_height", ctx.BlockHeight(),
				"retention_range", retentionRange)
			return
		}

		// Keep using existing prices if we're within retention range
		if ctx.BlockHeight()-lastValidVEHeight < retentionRange {
			k.Logger(ctx).Warn("last valid VE height is within price retention range; not zeroing asset prices",
				"retention_range", retentionRange)
			return
		}
	}

	// Outside retention range or prices are valid
	// Invalid prices will be zero/nil values
	if err := k.AssetPrices.Set(ctx, types.Asset_ROCK, rockPrice); err != nil {
		k.Logger(ctx).Error("error setting ROCK price", "height", ctx.BlockHeight(), "err", err)
	}

	if err := k.AssetPrices.Set(ctx, types.Asset_BTC, btcPrice); err != nil {
		k.Logger(ctx).Error("error setting BTC price", "height", ctx.BlockHeight(), "err", err)
	}

	if err := k.AssetPrices.Set(ctx, types.Asset_ETH, ethPrice); err != nil {
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
	btcUSDPrice math.LegacyDec,
	ethUSDPrice math.LegacyDec,
	exchangeRate math.LegacyDec,
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

	// Convert ETH fee to USD
	feeInUSD := new(big.Float).Mul(
		feeInETH,
		new(big.Float).SetFloat64(ethUSDPrice.MustFloat64()),
	)

	// Convert USD fee to BTC
	feeInBTC := new(big.Float).Quo(
		feeInUSD,
		new(big.Float).SetFloat64(btcUSDPrice.MustFloat64()),
	)

	// Convert to satoshis (multiply by 1e8)
	satoshisFloat := new(big.Float).Mul(
		feeInBTC,
		new(big.Float).SetInt64(1e8),
	)

	satoshisInt, _ := satoshisFloat.Int(nil)
	satoshis := satoshisInt.Uint64()

	// Convert BTC fee to zenBTC using exchange rate
	feeZenBTC := math.LegacyNewDecFromInt(math.NewIntFromUint64(satoshis)).Quo(exchangeRate).TruncateInt().Uint64()

	return feeZenBTC
}

// CalculateFlatZenBTCMintFee calculates a flat $5 fee in zenBTC
// Returns 0 if BTCUSDPrice is zero or exchangeRate is zero
func (k Keeper) CalculateFlatZenBTCMintFee(
	btcUSDPrice math.LegacyDec,
	exchangeRate math.LegacyDec,
) uint64 {
	if btcUSDPrice.IsZero() || exchangeRate.IsZero() {
		return 0
	}

	// Flat $5 fee
	feeUSD := math.LegacyNewDec(5)

	// Convert USD fee to BTC
	feeInBTC := feeUSD.Quo(btcUSDPrice)

	// Convert to satoshis (multiply by 1e8)
	satoshis := feeInBTC.Mul(math.LegacyNewDec(1e8)).TruncateInt().Uint64()

	// Convert BTC fee to zenBTC using exchange rate
	feeZenBTC := math.LegacyNewDecFromInt(math.NewIntFromUint64(satoshis)).Quo(exchangeRate).TruncateInt().Uint64()

	return feeZenBTC
}

// getPendingMintTransactionsByStatus retrieves up to 2 pending mint transactions matching the given status.
func (k *Keeper) getPendingMintTransactions(ctx sdk.Context, status zenbtctypes.MintTransactionStatus, walletType zenbtctypes.WalletType) ([]zenbtctypes.PendingMintTransaction, error) {
	firstPendingID := uint64(0)
	var err error
	if status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED {
		firstPendingID, err = k.zenBTCKeeper.GetFirstPendingStakeTransaction(ctx)
	} else if status == zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED {
		if walletType == zenbtctypes.WalletType_WALLET_TYPE_SOLANA {
			firstPendingID, err = k.zenBTCKeeper.GetFirstPendingSolMintTransaction(ctx)
		} else if walletType == zenbtctypes.WalletType_WALLET_TYPE_EVM {
			firstPendingID, err = k.zenBTCKeeper.GetFirstPendingEthMintTransaction(ctx)
		}
	}
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		firstPendingID = 0
	}
	var results []zenbtctypes.PendingMintTransaction
	err = k.zenBTCKeeper.WalkPendingMintTransactions(ctx, func(id uint64, tx zenbtctypes.PendingMintTransaction) (bool, error) {
		if id < firstPendingID {
			return false, nil // continue walking
		}

		isMatchingStatus := tx.Status == status
		if walletType == zenbtctypes.WalletType_WALLET_TYPE_UNSPECIFIED {
			if isMatchingStatus {
				results = append(results, tx)
				if len(results) >= 2 {
					return true, nil // stop walking
				}
			}
			return false, nil // continue walking
		}

		isMatchingNetwork := false
		if walletType == zenbtctypes.WalletType_WALLET_TYPE_SOLANA {
			isMatchingNetwork = types.IsSolanaCAIP2(ctx, tx.Caip2ChainId)
		} else if walletType == zenbtctypes.WalletType_WALLET_TYPE_EVM {
			isMatchingNetwork = types.IsEthereumCAIP2(ctx, tx.Caip2ChainId)
		}

		if isMatchingStatus && isMatchingNetwork {
			results = append(results, tx)
			if len(results) >= 2 {
				return true, nil // stop walking
			}
		}
		return false, nil // continue walking
	})

	return results, err
}

func (k *Keeper) getPendingDCTMintTransactions(ctx sdk.Context, asset dcttypes.Asset, status dcttypes.MintTransactionStatus, walletType dcttypes.WalletType) ([]dcttypes.PendingMintTransaction, error) {
	firstPendingID := uint64(0)
	var err error
	if status == dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED {
		firstPendingID, err = k.dctKeeper.GetFirstPendingStakeTransaction(ctx, asset)
	} else if status == dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED {
		if walletType == dcttypes.WalletType_WALLET_TYPE_SOLANA {
			firstPendingID, err = k.dctKeeper.GetFirstPendingSolMintTransaction(ctx, asset)
		} else if walletType == dcttypes.WalletType_WALLET_TYPE_EVM {
			firstPendingID, err = k.dctKeeper.GetFirstPendingEthMintTransaction(ctx, asset)
		}
	}
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		firstPendingID = 0
	}

	results := make([]dcttypes.PendingMintTransaction, 0, 2)
	err = k.dctKeeper.WalkPendingMintTransactions(ctx, asset, func(id uint64, tx dcttypes.PendingMintTransaction) (bool, error) {
		if id < firstPendingID {
			return false, nil
		}
		if tx.Status == status && tx.ChainType == walletType {
			results = append(results, tx)
			if len(results) >= 2 {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
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
	var results []zenbtctypes.BurnEvent
	err = k.zenBTCKeeper.WalkBurnEvents(ctx, func(id uint64, event zenbtctypes.BurnEvent) (bool, error) {
		if id < firstPendingID {
			return false, nil // continue walking
		}

		if event.Status == zenbtctypes.BurnStatus_BURN_STATUS_BURNED {
			results = append(results, event)
			if len(results) >= 2 {
				return true, nil // stop walking
			}
		}
		return false, nil // continue walking
	})

	return results, err
}

// getPendingRedemptions retrieves pending redemptions with the specified status.
// If limit is 0, all matching redemptions will be returned.
func (k *Keeper) GetRedemptionsByStatus(ctx sdk.Context, status zenbtctypes.RedemptionStatus, limit int, startingIndex uint64) ([]zenbtctypes.Redemption, error) {
	var results []zenbtctypes.Redemption
	err := k.zenBTCKeeper.WalkRedemptions(ctx, func(id uint64, r zenbtctypes.Redemption) (bool, error) {
		if id < startingIndex {
			return false, nil // continue walking
		}

		if r.Status == status {
			results = append(results, r)
			if limit > 0 && len(results) >= limit {
				return true, nil // stop walking
			}
		}
		return false, nil // continue walking
	})

	return results, err
}

// GetDCTRedemptionsByStatus retrieves DCT redemptions for a specific asset with the specified status.
// If limit is 0, all matching redemptions will be returned.
func (k *Keeper) GetDCTRedemptionsByStatus(ctx sdk.Context, asset dcttypes.Asset, status dcttypes.RedemptionStatus, limit int, startingIndex uint64) ([]dcttypes.Redemption, error) {
	var results []dcttypes.Redemption
	err := k.dctKeeper.WalkRedemptions(ctx, asset, func(id uint64, r dcttypes.Redemption) (bool, error) {
		if id < startingIndex {
			return false, nil // continue walking
		}

		if r.Status == status {
			results = append(results, r)
			if limit > 0 && len(results) >= limit {
				return true, nil // stop walking
			}
		}
		return false, nil // continue walking
	})

	return results, err
}

func (k *Keeper) recordMismatchedVoteExtensions(ctx sdk.Context, height int64, pluralityVoteExt VoteExtension, consensusData abci.ExtendedCommitInfo) {
	// Compare against the plurality vote extension (most voted values) rather than consensus.
	// This prevents jailing the entire validator set during prolonged periods without supermajority consensus.
	// If the plurality VE is empty, it means no votes were cast.
	// In this case, a validator that submitted an empty vote extension should not be penalized.
	isCanonicalEmpty := reflect.DeepEqual(pluralityVoteExt, VoteExtension{})

	canonicalVoteExtBz, err := json.Marshal(pluralityVoteExt)
	if err != nil {
		k.Logger(ctx).Error("error marshalling canonical vote extension", "height", height, "error", err)
		return
	}

	// Keep track of unique validators that had mismatches in this block - O(N)
	mismatchedValidators := make(map[string]struct{})

	for _, v := range consensusData.Votes {
		// If the canonical VE is empty and the validator also submitted an empty VE, it's not a mismatch.
		if isCanonicalEmpty && len(v.VoteExtension) == 0 {
			continue
		}

		if !bytes.Equal(v.VoteExtension, canonicalVoteExtBz) {
			validatorHexAddr := hex.EncodeToString(v.Validator.Address)

			// Skip recording mismatches for jailed or unbonded validators
			consAddr := sdk.ConsAddress(v.Validator.Address)
			validator, err := k.GetValidatorByConsAddr(ctx, consAddr)
			if err != nil {
				k.Logger(ctx).Error("Failed to get validator by consensus address", "consAddr", consAddr.String(), "error", err)
				continue
			}

			// Check if validator should be skipped
			if validator.Jailed || validator.Status == types.Unbonded {
				continue
			}

			mismatchedValidators[validatorHexAddr] = struct{}{}

			// Still record in ValidationInfo for backward compatibility
			info, err := k.ValidationInfos.Get(ctx, height)
			if err != nil {
				info = types.ValidationInfo{}
			}
			info.MismatchedVoteExtensions = append(info.MismatchedVoteExtensions, validatorHexAddr)
			if err := k.ValidationInfos.Set(ctx, height, info); err != nil {
				k.Logger(ctx).Error("error setting validation info", "height", height, "error", err)
			}
		}
	}

	// Update the sliding window counters for each mismatched validator - O(M) where M = unique mismatched validators
	for validatorHexAddr := range mismatchedValidators {
		k.updateValidatorMismatchCount(ctx, validatorHexAddr, height)
	}
}

// updateValidatorMismatchCount updates the sliding window counter for a validator's mismatches
// Time complexity: O(1) amortized per call
func (k *Keeper) updateValidatorMismatchCount(ctx sdk.Context, validatorHexAddr string, blockHeight int64) {
	// Get existing count or create new one
	mismatchCount, err := k.ValidatorMismatchCounts.Get(ctx, validatorHexAddr)
	if err != nil {
		// No existing record, create new
		mismatchCount = types.ValidatorMismatchCount{
			ValidatorAddress: validatorHexAddr,
			MismatchBlocks:   []int64{blockHeight},
			TotalCount:       1,
		}
		if err := k.ValidatorMismatchCounts.Set(ctx, validatorHexAddr, mismatchCount); err != nil {
			k.Logger(ctx).Error("error setting validator mismatch count", "validator", validatorHexAddr, "error", err)
		}
		return
	}

	// Remove blocks that are outside the sliding window (older than configured window size)
	windowStart := blockHeight - k.GetVEWindowSize(ctx) + 1
	newMismatchBlocks := make([]int64, 0, len(mismatchCount.MismatchBlocks)+1)

	// Keep only blocks within the window - O(W) where W is configurable window size
	for _, block := range mismatchCount.MismatchBlocks {
		if block >= windowStart {
			newMismatchBlocks = append(newMismatchBlocks, block)
		}
	}

	// Add the new block (maintaining sorted order since we always append increasing heights)
	newMismatchBlocks = append(newMismatchBlocks, blockHeight)

	// Update the count
	mismatchCount.MismatchBlocks = newMismatchBlocks
	mismatchCount.TotalCount = uint32(len(newMismatchBlocks))

	if err := k.ValidatorMismatchCounts.Set(ctx, validatorHexAddr, mismatchCount); err != nil {
		k.Logger(ctx).Error("error updating validator mismatch count", "validator", validatorHexAddr, "error", err)
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
	var chainID uint64 = ethereum.HoodiChainId.Uint64()
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

// validationMismatch holds details about a validation failure.
type validationMismatch struct {
	field    VoteExtensionField
	expected string
	actual   string
	err      error // Optional: Used for hash validation errors
}

// handleValidationMismatches processes fields that failed validation.
// It logs detailed mismatch information and removes the fields from consensus consideration.
func (k *Keeper) handleValidationMismatches(ctx context.Context, mismatches []validationMismatch, fieldVotePowers map[VoteExtensionField]int64, oracleData *OracleData) {
	if len(mismatches) == 0 {
		return
	}

	for _, mismatch := range mismatches {
		logCtx := []any{
			"field", mismatch.field.String(),
			"consensus_status", "revoked",
		}
		if mismatch.err != nil {
			// Primarily for hash mismatches where the error contains details
			logCtx = append(logCtx, "error", mismatch.err)
		} else {
			// For direct value mismatches
			logCtx = append(logCtx, "expected", mismatch.expected, "actual", mismatch.actual)
		}

		k.Logger(ctx).Warn("field had consensus but failed data validation", logCtx...)

		delete(fieldVotePowers, mismatch.field)
	}

	// Update FieldVotePowers in oracleData to reflect the validated fields
	if oracleData != nil {
		oracleData.FieldVotePowers = fieldVotePowers
	}
}

// validateHashField derives a hash from the given data and compares it with the expected value.
// Returns a detailed error on mismatch.
func validateHashField(fieldName string, expectedHash []byte, data any) error {
	derivedHash, err := deriveHash(data)
	if err != nil {
		return fmt.Errorf("error deriving %s hash: %w", fieldName, err)
	}
	if !bytes.Equal(expectedHash, derivedHash[:]) {
		// Ensure error message clearly shows expected vs actual
		return fmt.Errorf("%s hash mismatch, expected %x, got %x", fieldName, expectedHash, derivedHash[:])
	}
	return nil
}

// validateOracleData verifies that the vote extension and oracle data match.
// Only fields that have reached consensus (present in fieldVotePowers) are validated.
// Fields that fail validation are collected and handled by handleValidationMismatches.
func (k *Keeper) validateOracleData(ctx context.Context, voteExt VoteExtension, oracleData *OracleData, fieldVotePowers map[VoteExtensionField]int64) {
	mismatches := make([]validationMismatch, 0)

	// Helper function to add mismatch details
	recordMismatch := func(field VoteExtensionField, expected, actual any, err ...error) {
		mismatch := validationMismatch{field: field}
		if len(err) > 0 && err[0] != nil {
			mismatch.err = err[0]
		} else {
			// Format expected value
			if b, ok := expected.([]byte); ok {
				mismatch.expected = "0x" + hex.EncodeToString(b)
			} else {
				mismatch.expected = fmt.Sprintf("%v", expected)
			}

			// Format actual value
			if b, ok := actual.([]byte); ok {
				mismatch.actual = "0x" + hex.EncodeToString(b)
			} else {
				mismatch.actual = fmt.Sprintf("%v", actual)
			}
		}
		mismatches = append(mismatches, mismatch)
	}

	// Validate whether hash fields have consensus
	if fieldHasConsensus(fieldVotePowers, VEFieldEigenDelegationsHash) {
		if err := validateHashField(VEFieldEigenDelegationsHash.String(), voteExt.EigenDelegationsHash, oracleData.EigenDelegationsMap); err != nil {
			recordMismatch(VEFieldEigenDelegationsHash, voteExt.EigenDelegationsHash, "derived_hash", err)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldEthBurnEventsHash) {
		if err := validateHashField(VEFieldEthBurnEventsHash.String(), voteExt.EthBurnEventsHash, oracleData.EthBurnEvents); err != nil {
			recordMismatch(VEFieldEthBurnEventsHash, voteExt.EthBurnEventsHash, "derived_hash", err)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldRedemptionsHash) {
		if err := validateHashField(VEFieldRedemptionsHash.String(), voteExt.RedemptionsHash, oracleData.Redemptions); err != nil {
			recordMismatch(VEFieldRedemptionsHash, voteExt.RedemptionsHash, "derived_hash", err)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldSolanaAccountsHash) {
		if err := validateHashField(VEFieldSolanaAccountsHash.String(), voteExt.SolanaAccountsHash, oracleData.SolanaAccounts); err != nil {
			recordMismatch(VEFieldSolanaAccountsHash, voteExt.SolanaAccountsHash, "derived_hash", err)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldSolanaMintNoncesHash) {
		if err := validateHashField(VEFieldSolanaMintNoncesHash.String(), voteExt.SolanaMintNoncesHash, oracleData.SolanaMintNonces); err != nil {
			recordMismatch(VEFieldSolanaMintNoncesHash, voteExt.SolanaMintNoncesHash, "derived_hash", err)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldSolanaMintEventsHash) {
		if err := validateHashField(VEFieldSolanaMintEventsHash.String(), voteExt.SolanaMintEventsHash, oracleData.SolanaMintEvents); err != nil {
			recordMismatch(VEFieldSolanaMintEventsHash, voteExt.SolanaMintEventsHash, "derived_hash", err)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldSolanaBurnEventsHash) {
		if err := validateHashField(VEFieldSolanaBurnEventsHash.String(), voteExt.SolanaBurnEventsHash, oracleData.SolanaBurnEvents); err != nil {
			recordMismatch(VEFieldSolanaBurnEventsHash, voteExt.SolanaBurnEventsHash, "derived_hash", err)
		}
	}
	// Skip RequestedBtcHeaderHash validation when there are no requested headers (indicated by RequestedBtcBlockHeight == 0)
	if fieldHasConsensus(fieldVotePowers, VEFieldRequestedBtcHeaderHash) {
		if oracleData.RequestedBtcBlockHeight != 0 {
			if err := validateHashField(VEFieldRequestedBtcHeaderHash.String(), voteExt.RequestedBtcHeaderHash, &oracleData.RequestedBtcBlockHeader); err != nil {
				recordMismatch(VEFieldRequestedBtcHeaderHash, voteExt.RequestedBtcHeaderHash, "derived_hash", err)
			}
		} else if len(voteExt.RequestedBtcHeaderHash) > 0 {
			// Mismatch if oracle has no requested height but VE has a hash
			recordMismatch(VEFieldRequestedBtcHeaderHash, hex.EncodeToString(voteExt.RequestedBtcHeaderHash), "nil (no requested height in oracleData)")
		}
	}

	// Check Ethereum-related fields
	if fieldHasConsensus(fieldVotePowers, VEFieldEthBlockHeight) {
		if voteExt.EthBlockHeight != oracleData.EthBlockHeight {
			recordMismatch(VEFieldEthBlockHeight, voteExt.EthBlockHeight, oracleData.EthBlockHeight)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldEthGasLimit) {
		if voteExt.EthGasLimit != oracleData.EthGasLimit {
			recordMismatch(VEFieldEthGasLimit, voteExt.EthGasLimit, oracleData.EthGasLimit)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldEthBaseFee) {
		if voteExt.EthBaseFee != oracleData.EthBaseFee {
			recordMismatch(VEFieldEthBaseFee, voteExt.EthBaseFee, oracleData.EthBaseFee)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldEthTipCap) {
		if voteExt.EthTipCap != oracleData.EthTipCap {
			recordMismatch(VEFieldEthTipCap, voteExt.EthTipCap, oracleData.EthTipCap)
		}
	}

	// Check Bitcoin height
	if fieldHasConsensus(fieldVotePowers, VEFieldRequestedBtcBlockHeight) {
		if voteExt.RequestedBtcBlockHeight != oracleData.RequestedBtcBlockHeight {
			recordMismatch(VEFieldRequestedBtcBlockHeight, voteExt.RequestedBtcBlockHeight, oracleData.RequestedBtcBlockHeight)
		}
	}

	// Check nonce-related fields
	if fieldHasConsensus(fieldVotePowers, VEFieldRequestedStakerNonce) {
		if voteExt.RequestedStakerNonce != oracleData.RequestedStakerNonce {
			recordMismatch(VEFieldRequestedStakerNonce, voteExt.RequestedStakerNonce, oracleData.RequestedStakerNonce)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldRequestedEthMinterNonce) {
		if voteExt.RequestedEthMinterNonce != oracleData.RequestedEthMinterNonce {
			recordMismatch(VEFieldRequestedEthMinterNonce, voteExt.RequestedEthMinterNonce, oracleData.RequestedEthMinterNonce)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldRequestedUnstakerNonce) {
		if voteExt.RequestedUnstakerNonce != oracleData.RequestedUnstakerNonce {
			recordMismatch(VEFieldRequestedUnstakerNonce, voteExt.RequestedUnstakerNonce, oracleData.RequestedUnstakerNonce)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldRequestedCompleterNonce) {
		if voteExt.RequestedCompleterNonce != oracleData.RequestedCompleterNonce {
			recordMismatch(VEFieldRequestedCompleterNonce, voteExt.RequestedCompleterNonce, oracleData.RequestedCompleterNonce)
		}
	}

	// Check price fields - compare as strings first for simplicity
	if fieldHasConsensus(fieldVotePowers, VEFieldROCKUSDPrice) {
		if voteExt.ROCKUSDPrice != oracleData.ROCKUSDPrice {
			// Log the string mismatch; deeper decimal comparison could be added if needed
			recordMismatch(VEFieldROCKUSDPrice, voteExt.ROCKUSDPrice, oracleData.ROCKUSDPrice)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldBTCUSDPrice) {
		if voteExt.BTCUSDPrice != oracleData.BTCUSDPrice {
			recordMismatch(VEFieldBTCUSDPrice, voteExt.BTCUSDPrice, oracleData.BTCUSDPrice)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldETHUSDPrice) {
		if voteExt.ETHUSDPrice != oracleData.ETHUSDPrice {
			recordMismatch(VEFieldETHUSDPrice, voteExt.ETHUSDPrice, oracleData.ETHUSDPrice)
		}
	}

	// Check Latest Bitcoin height and hash fields
	if fieldHasConsensus(fieldVotePowers, VEFieldLatestBtcBlockHeight) {
		if voteExt.LatestBtcBlockHeight != oracleData.LatestBtcBlockHeight {
			recordMismatch(VEFieldLatestBtcBlockHeight, voteExt.LatestBtcBlockHeight, oracleData.LatestBtcBlockHeight)
		}
	}
	if fieldHasConsensus(fieldVotePowers, VEFieldLatestBtcHeaderHash) {
		if err := validateHashField(VEFieldLatestBtcHeaderHash.String(), voteExt.LatestBtcHeaderHash, oracleData.LatestBtcBlockHeader); err != nil {
			recordMismatch(VEFieldLatestBtcHeaderHash, voteExt.LatestBtcHeaderHash, oracleData.LatestBtcBlockHeader, err)
		}
	}

	// Handle all collected mismatches
	k.handleValidationMismatches(ctx, mismatches, fieldVotePowers, oracleData)
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
			KeyIds:              []uint64{keyID},
			WalletType:          walletType,
			UnsignedTransaction: unsignedTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedTxHash)),
	)
	return err
}

// Helper function to submit Ethereum transactions
func (k *Keeper) submitSolanaTransaction(ctx sdk.Context, creator string, keyIDs []uint64, walletType treasurytypes.WalletType, chainID string, unsignedTx []byte) (uint64, error) {
	metadata, err := codectypes.NewAnyWithValue(&treasurytypes.MetadataSolana{
		Network: zentptypes.Caip2ToSolananNetwork(chainID)})
	if err != nil {
		return 0, err
	}
	resp, err := k.treasuryKeeper.HandleSignTransactionRequest(
		ctx,
		&treasurytypes.MsgNewSignTransactionRequest{
			Creator:             creator,
			KeyIds:              keyIDs,
			WalletType:          walletType,
			UnsignedTransaction: unsignedTx,
			Metadata:            metadata,
			NoBroadcast:         false,
		},
		[]byte(hex.EncodeToString(unsignedTx)),
	)
	if err != nil {
		return 0, err
	}
	return resp.Id, nil
}

func (k Keeper) GetSolanaNonceAccount(goCtx context.Context, keyID uint64) (system.NonceAccount, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	key, err := k.treasuryKeeper.GetKey(ctx, keyID)
	if err != nil {
		return system.NonceAccount{}, err
	}
	publicKey, err := treasurytypes.SolanaPubkey(key)
	if err != nil {
		return system.NonceAccount{}, err
	}
	resp, err := k.sidecarClient.GetSolanaAccountInfo(goCtx, &sidecar.SolanaAccountInfoRequest{
		PubKey: publicKey.String(),
	})
	if err != nil {
		k.Logger(ctx).Error("failed to get solana account info from sidecar: ", err)
		return system.NonceAccount{}, err
	}
	nonceAccount := system.NonceAccount{}

	if len(resp.Account) == 0 {
		return nonceAccount, fmt.Errorf("nonce account %s is likely not a valid nonce account.", publicKey.String())
	}

	decoder := bin.NewBorshDecoder(resp.Account)

	if err = nonceAccount.UnmarshalWithDecoder(decoder); err != nil {
		return nonceAccount, fmt.Errorf("failed to unmarshal nonce account: %w (data length: %d)", err, len(resp.Account))
	}
	return nonceAccount, err
}

func (k Keeper) GetSolanaTokenAccount(goCtx context.Context, address, mint string) (token.Account, error) {
	recipientPubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return token.Account{}, err
	}
	mintPubKey, err := solana.PublicKeyFromBase58(mint)
	if err != nil {
		return token.Account{}, err
	}
	receiverAta, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintPubKey)
	if err != nil {
		return token.Account{}, err
	}
	resp, err := k.sidecarClient.GetSolanaAccountInfo(goCtx, &sidecar.SolanaAccountInfoRequest{
		PubKey: receiverAta.String(),
	})
	if err != nil {
		if err.Error() == "rpc error: code = Unknown desc = not found" {
			return token.Account{}, nil
		}
		return token.Account{}, err
	}

	tokenAccount := new(token.Account)

	if resp.Account == nil {
		return *tokenAccount, nil
	}
	decoder := bin.NewBorshDecoder(resp.Account)

	err = tokenAccount.UnmarshalWithDecoder(decoder)
	if err != nil {
		return token.Account{}, err
	}

	return *tokenAccount, nil
}

// newCreateATAIdempotentInstruction creates an idempotent ATA creation instruction.
// This instruction will succeed even if the ATA already exists, preventing race conditions.
// The idempotent version uses instruction data [1] instead of [] (empty bytes).
func newCreateATAIdempotentInstruction(payer, wallet, mint solana.PublicKey) solana.Instruction {
	createATAInstruction := ata.NewCreateInstruction(
		payer,
		wallet,
		mint,
	).Build()

	// Override instruction data with [1] to make it idempotent
	return solana.NewInstruction(
		createATAInstruction.ProgramID(),
		createATAInstruction.Accounts(),
		[]byte{1},
	)
}

type solanaMintTxRequest struct {
	amount            uint64
	fee               uint64
	recipient         string
	nonce             *system.NonceAccount
	fundReceiver      bool
	programID         string
	mintAddress       string
	feeWallet         string
	nonceAccountKey   uint64
	nonceAuthorityKey uint64
	signerKey         uint64
	multisigKey       string
	rock              bool
	zenbtc            bool
}

func (k Keeper) PrepareSolanaMintTx(goCtx context.Context, req *solanaMintTxRequest) ([]byte, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	programID, err := solana.PublicKeyFromBase58(req.programID)
	if err != nil {
		return nil, err
	}

	nonceAccKey, err := k.treasuryKeeper.GetKey(ctx, req.nonceAccountKey)
	if err != nil {
		return nil, err
	}

	nonceAccPubKey, err := treasurytypes.SolanaPubkey(nonceAccKey)
	if err != nil {
		return nil, err
	}

	nonceAuthKey, err := k.treasuryKeeper.GetKey(ctx, req.nonceAuthorityKey)
	if err != nil {
		return nil, err
	}
	nonceAuthPubKey, err := treasurytypes.SolanaPubkey(nonceAuthKey)
	if err != nil {
		return nil, err
	}

	signerKey, err := k.treasuryKeeper.GetKey(ctx, req.signerKey)
	if err != nil {
		return nil, err
	}
	signerPubKey, err := treasurytypes.SolanaPubkey(signerKey)
	if err != nil {
		return nil, err
	}

	mintKey, err := solana.PublicKeyFromBase58(req.mintAddress)
	if err != nil {
		return nil, err
	}

	feeKey, err := solana.PublicKeyFromBase58(req.feeWallet)
	if err != nil {
		return nil, err
	}

	recipientPubKey, err := solana.PublicKeyFromBase58(req.recipient)
	if err != nil {
		return nil, err
	}

	var instructions []solana.Instruction

	instructions = append(instructions, system.NewAdvanceNonceAccountInstruction(
		*nonceAccPubKey,
		solana.SysVarRecentBlockHashesPubkey,
		*nonceAuthPubKey,
	).Build())

	feeWalletAta, _, err := solana.FindAssociatedTokenAddress(feeKey, mintKey)
	if err != nil {
		return nil, err
	}

	receiverAta, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintKey)
	if err != nil {
		return nil, err
	}

	if req.fundReceiver {
		instructions = append(
			instructions,
			newCreateATAIdempotentInstruction(
				*signerPubKey,
				recipientPubKey,
				mintKey,
			),
		)
		k.Logger(ctx).Info(
			"Added idempotent ATA creation instruction to tx",
			"signer", *signerPubKey,
			"recipient", recipientPubKey.String(),
			"mint", mintKey.String(),
		)
	}
	if req.rock {
		instructions = append(instructions, solrock.Wrap(
			programID,
			rock_spl_token.WrapArgs{
				Value: req.amount,
				Fee:   req.fee,
			},
			*signerPubKey,
			mintKey,
			feeKey,
			feeWalletAta,
			recipientPubKey,
			receiverAta,
		))
	} else if req.zenbtc {
		var multiSigKey solana.PublicKey
		if req.multisigKey != "" {
			multiSigKey, err = solana.PublicKeyFromBase58(req.multisigKey)
			if err != nil {
				return nil, err
			}
		} else {
			multiSigKeyAddress := k.zenBTCKeeper.GetSolanaParams(ctx).MultisigKeyAddress
			multiSigKey, err = solana.PublicKeyFromBase58(multiSigKeyAddress)
			if err != nil {
				return nil, err
			}
		}

		instructions = append(instructions, solzenbtc.Wrap(
			programID,
			zenbtc_spl_token.WrapArgs{Value: req.amount, Fee: req.fee},
			*signerPubKey,
			mintKey,
			multiSigKey,
			feeKey,
			feeWalletAta,
			recipientPubKey,
			receiverAta,
		))

		k.Logger(ctx).Info("Added wrap instruction to tx",
			"programID", programID.String(),
			"amount", req.amount,
			"fee", req.fee,
			"signerPubKey", *signerPubKey,
			"mintKey", mintKey.String(),
			"multisigKey", multiSigKey.String(),
			"feeKey", feeKey.String(),
			"feeWalletAta", feeWalletAta.String(),
			"recipientWalletPubKey", recipientPubKey.String(),
			"receiverAta", receiverAta.String(),
		)
	} else {
		return nil, fmt.Errorf("neither rock nor zenbtc flag is set")
	}

	tx, err := solana.NewTransaction(
		instructions,
		solana.Hash(req.nonce.Nonce),
		solana.TransactionPayer(*signerPubKey),
	)
	if err != nil {
		return nil, err
	}
	txBytes, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}

func (k Keeper) retrieveSolanaNonces(goCtx context.Context) (map[uint64]*system.NonceAccount, error) {
	nonces := map[uint64]*system.NonceAccount{}
	//pendingSolROCKMints, err := k.zentpKeeper.GetMintsWithStatus(goCtx, zentptypes.BridgeStatus_BRIDGE_STATUS_PENDING)
	//if err != nil {
	//	return nil, err
	//}
	//if len(pendingSolROCKMints) == 0 {
	solParams := k.zentpKeeper.GetSolanaParams(goCtx)
	solNonceRequested, err := k.SolanaNonceRequested.Get(goCtx, solParams.NonceAccountKey)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
	}
	if solNonceRequested {
		n, err := k.GetSolanaNonceAccount(goCtx, solParams.NonceAccountKey)
		nonces[solParams.NonceAccountKey] = &n
		if err != nil {
			return nil, err
		}
	}
	//}

	//ctx := sdk.UnwrapSDKContext(goCtx)
	//pendingZenBTCMints, err := k.getPendingMintTransactions(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, zenbtctypes.WalletType_WALLET_TYPE_SOLANA)
	//if err != nil {
	//	return nil, err
	//}
	//if len(pendingZenBTCMints) == 0 {
	zenBTCsolParams := k.zenBTCKeeper.GetSolanaParams(goCtx)
	zenBTCsolNonceRequested, err := k.SolanaNonceRequested.Get(goCtx, zenBTCsolParams.NonceAccountKey)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		zenBTCsolNonceRequested = false
	}

	if zenBTCsolNonceRequested {
		nonceAcc, err := k.GetSolanaNonceAccount(goCtx, zenBTCsolParams.NonceAccountKey)
		if err != nil {
			return nil, err
		}
		nonces[zenBTCsolParams.NonceAccountKey] = &nonceAcc
	}
	//}

	// Retrieve nonce accounts for all DCT assets (zenZEC, etc.)
	if k.dctKeeper != nil {
		assets, err := k.dctKeeper.ListSupportedAssets(goCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to list DCT assets: %w", err)
		}

		for _, asset := range assets {
			dctSolParams, err := k.dctKeeper.GetSolanaParams(goCtx, asset)
			if err != nil {
				// Skip assets without Solana params
				continue
			}
			if dctSolParams == nil {
				continue
			}

			dctNonceRequested, err := k.SolanaNonceRequested.Get(goCtx, dctSolParams.NonceAccountKey)
			if err != nil {
				if !errors.Is(err, collections.ErrNotFound) {
					return nil, fmt.Errorf("failed to check nonce request for asset %s: %w", asset.String(), err)
				}
				dctNonceRequested = false
			}

			if dctNonceRequested {
				nonceAcc, err := k.GetSolanaNonceAccount(goCtx, dctSolParams.NonceAccountKey)
				if err != nil {
					return nil, fmt.Errorf("failed to get nonce account for asset %s (key %d): %w", asset.String(), dctSolParams.NonceAccountKey, err)
				}
				nonces[dctSolParams.NonceAccountKey] = &nonceAcc
			}
		}
	}

	return nonces, nil
}

// populateAccountsForSolanaMints is a helper to collect Solana token accounts for a specific mint and request store.
// It populates the provided solAccs map with accounts keyed by their ATA addresses.
func (k Keeper) populateAccountsForSolanaMints(
	ctx context.Context,
	requestStore collections.Map[string, bool],
	mintAddress string,
	flowDescription string, // For logging purposes, e.g., "ZenBTC" or "ZenTP"
	solAccs map[string]token.Account, // Map to populate
) error {
	storeIterator, err := requestStore.Iterate(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to iterate %s Solana account requests: %w", flowDescription, err)
	}
	defer storeIterator.Close()

	ownerKeys, err := storeIterator.Keys()
	if err != nil {
		return fmt.Errorf("failed to get keys from %s Solana account request store: %w", flowDescription, err)
	}

	for _, ownerAddressStr := range ownerKeys {
		requested, err := requestStore.Get(ctx, ownerAddressStr)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error(fmt.Sprintf("Owner address key not found during %s Solana account collection, skipping", flowDescription), "key", ownerAddressStr)
				continue
			}
			return fmt.Errorf("failed to check if %s account %s is requested: %w", flowDescription, ownerAddressStr, err)
		}

		if requested {
			ownerPubKey, err := solana.PublicKeyFromBase58(ownerAddressStr)
			if err != nil {
				k.Logger(ctx).Error(fmt.Sprintf("Invalid owner address for %s, cannot derive ATA", flowDescription), "owner", ownerAddressStr, "error", err)
				continue
			}
			mintPubKey, err := solana.PublicKeyFromBase58(mintAddress)
			if err != nil {
				// This should ideally not happen if mintAddress was pre-validated
				k.Logger(ctx).Error(fmt.Sprintf("Invalid %s mint address, cannot derive ATA", flowDescription), "mint", mintAddress, "error", err)
				continue // Or return error, as this is a config issue
			}
			ataAddress, _, err := solana.FindAssociatedTokenAddress(ownerPubKey, mintPubKey)
			if err != nil {
				k.Logger(ctx).Error(fmt.Sprintf("Failed to derive ATA for %s account", flowDescription), "owner", ownerAddressStr, "mint", mintAddress, "error", err)
				continue
			}

			acc, err := k.GetSolanaTokenAccount(ctx, ownerAddressStr, mintAddress)
			if err != nil {
				// Log error but continue, and store the (potentially zero-value) account
				// This preserves existing behavior for hashing consistency.
				k.Logger(ctx).Error(fmt.Sprintf("Failed to get Solana token account for %s", flowDescription), "owner", ownerAddressStr, "mint", mintAddress, "ata", ataAddress.String(), "error", err)
			}
			solAccs[ataAddress.String()] = acc
		}
	}
	return nil
}

// populateDCTAccountsForAsset collects Solana token accounts for DCT assets scoped by asset identifier.
func (k Keeper) populateDCTAccountsForAsset(
	ctx context.Context,
	asset dcttypes.Asset,
	mintAddress string,
	flowDescription string,
	solAccs map[string]token.Account,
) error {
	assetKey := asset.String()
	rangeForAsset := collections.NewPrefixedPairRange[string, string](assetKey)
	return k.SolanaDCTAccountsRequested.Walk(ctx, rangeForAsset, func(key collections.Pair[string, string], requested bool) (bool, error) {
		if !requested {
			return false, nil
		}
		ownerAddressStr := key.K2()
		ownerPubKey, err := solana.PublicKeyFromBase58(ownerAddressStr)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Invalid owner address for %s", flowDescription), "owner", ownerAddressStr, "error", err)
			return false, nil
		}
		mintPubKey, err := solana.PublicKeyFromBase58(mintAddress)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Invalid %s mint address", flowDescription), "mint", mintAddress, "error", err)
			return false, nil
		}
		ataAddress, _, err := solana.FindAssociatedTokenAddress(ownerPubKey, mintPubKey)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Failed to derive ATA for %s account", flowDescription), "owner", ownerAddressStr, "mint", mintAddress, "error", err)
			return false, nil
		}

		acc, err := k.GetSolanaTokenAccount(ctx, ownerAddressStr, mintAddress)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Failed to get Solana token account for %s", flowDescription), "owner", ownerAddressStr, "mint", mintAddress, "ata", ataAddress.String(), "error", err)
		}
		solAccs[ataAddress.String()] = acc
		return false, nil
	})
}

// retrieveSolanaAccounts retrieves all requested Solana token accounts.
// The returned map keys are ATA addresses.
func (k Keeper) retrieveSolanaAccounts(ctx context.Context) (map[string]token.Account, error) {
	solAccs := make(map[string]token.Account) // Key: ATA Address string

	// 1. Process ZenBTC related accounts
	zenBTCMintAddress := ""
	if k.zenBTCKeeper != nil && k.zenBTCKeeper.GetSolanaParams(ctx) != nil {
		zenBTCMintAddress = k.zenBTCKeeper.GetSolanaParams(ctx).MintAddress
	}

	if zenBTCMintAddress == "" {
		k.Logger(ctx).Warn("ZenBTC Solana mint address is not configured. Skipping ZenBTC account collection.")
	} else {
		if err := k.populateAccountsForSolanaMints(ctx, k.SolanaAccountsRequested, zenBTCMintAddress, "ZenBTC", solAccs); err != nil {
			// The helper function already logs specifics, so we just bubble up a general error here if needed.
			return nil, fmt.Errorf("error processing ZenBTC Solana account requests: %w", err)
		}
	}

	// 2. Process ZenTP related accounts
	zenTPMintAddress := ""
	if k.zentpKeeper != nil && k.zentpKeeper.GetSolanaParams(ctx) != nil { // Assuming zentpKeeper has GetSolanaParams
		zenTPMintAddress = k.zentpKeeper.GetSolanaParams(ctx).MintAddress
	}

	if zenTPMintAddress == "" {
		k.Logger(ctx).Warn("ZenTP Solana mint address is not configured. Skipping ZenTP account collection.")
	} else {
		if err := k.populateAccountsForSolanaMints(ctx, k.SolanaZenTPAccountsRequested, zenTPMintAddress, "ZenTP", solAccs); err != nil {
			return nil, fmt.Errorf("error processing ZenTP Solana account requests: %w", err)
		}
	}

	if k.dctKeeper != nil {
		assets, err := k.dctKeeper.ListSupportedAssets(ctx)
		if err != nil {
			k.Logger(ctx).Error("error listing DCT assets for Solana account collection", "error", err)
		} else {
			for _, asset := range assets {
				solParams, err := k.dctKeeper.GetSolanaParams(ctx, asset)
				if err != nil {
					k.Logger(ctx).Error("error fetching DCT Solana params", "asset", asset.String(), "error", err)
					continue
				}
				if solParams == nil || solParams.MintAddress == "" {
					continue
				}
				if err := k.populateDCTAccountsForAsset(ctx, asset, solParams.MintAddress, fmt.Sprintf("DCT-%s", asset.String()), solAccs); err != nil {
					return nil, fmt.Errorf("error processing DCT Solana account requests for asset %s: %w", asset.String(), err)
				}
			}
		}
	}

	return solAccs, nil
}

func (k Keeper) clearSolanaAccounts(ctx sdk.Context) {
	pendingsROCK, err := k.zentpKeeper.GetMintsWithStatusPending(ctx)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	pendingsZenBTC, err := k.getPendingMintTransactions(ctx, zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, zenbtctypes.WalletType_WALLET_TYPE_SOLANA)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	// Clear ZenBTC related requests if no pending ZenBTC Solana mints
	if len(pendingsZenBTC) == 0 {
		if err = k.SolanaAccountsRequested.Clear(ctx, nil); err != nil {
			k.Logger(ctx).Error("Error clearing SolanaAccountsRequested (ZenBTC): " + err.Error())
		}
	}

	// Clear ZenTP related requests if no pending ROCK Solana mints (assuming pendingsROCK is for ZenTP)
	if len(pendingsROCK) == 0 {
		if err = k.SolanaZenTPAccountsRequested.Clear(ctx, nil); err != nil {
			k.Logger(ctx).Error("Error clearing SolanaZenTPAccountsRequested: " + err.Error())
		}
	}

	if k.dctKeeper != nil {
		assets, err := k.dctKeeper.ListSupportedAssets(ctx)
		if err != nil {
			k.Logger(ctx).Error("error listing DCT assets for clearing Solana accounts", "error", err)
		} else {
			for _, asset := range assets {
				pendingDCT, err := k.getPendingDCTMintTransactions(ctx, asset, dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED, dcttypes.WalletType_WALLET_TYPE_SOLANA)
				if err != nil {
					k.Logger(ctx).Error("error fetching pending DCT Solana mints during cleanup", "asset", asset.String(), "error", err)
					continue
				}
				if len(pendingDCT) == 0 {
					rangeForAsset := collections.NewPrefixedPairRange[string, string](asset.String())
					if err := k.SolanaDCTAccountsRequested.Clear(ctx, rangeForAsset); err != nil {
						k.Logger(ctx).Error("error clearing DCT Solana account requests", "asset", asset.String(), "error", err)
					}
				}
			}
		}
	}
}

const (
	// solanaEventConfirmationWindowBlocks defines the window in blocks to wait for a Solana event
	// after a nonce has been observed to advance, before considering the transaction timed out.
	// TODO: This should ideally be a configurable module parameter.
	solanaEventConfirmationWindowBlocks = 100
)

// handleSolanaTransactionBTLTimeout contains the core logic for BTL timeout processing.
// It checks if the on-chain nonce has advanced compared to the last used nonce.
// If not, it resets the transaction's block height and awaiting event status for a full retry.
// If the nonce has advanced, it sets the awaiting event status if it wasn't already set.
func (k Keeper) handleSolanaTransactionBTLTimeout(
	ctx sdk.Context,
	txID uint64,
	txBlockHeight int64,
	txAwaitingEventSince int64,
	solParamsNonceAccountKey uint64,
	solParamsBTL int64,
	oracleData OracleData,
) (newBlockHeight int64, newAwaitingEventSince int64) {
	k.Logger(ctx).Debug("BTL Timeout Logic: Initiating checks.", "tx_id", txID, "tx_block_height", txBlockHeight, "btl", solParamsBTL)

	// Initialize return values to current tx values
	newBlockHeight = txBlockHeight
	newAwaitingEventSince = txAwaitingEventSince

	currentLiveNonceAccount := oracleData.SolanaMintNonces[solParamsNonceAccountKey]

	if currentLiveNonceAccount == nil || currentLiveNonceAccount.Nonce.IsZero() {
		k.Logger(ctx).Warn("BTL Logic: Live Solana nonce is zero or unavailable in oracleData. Resetting transaction.", "tx_id", txID)
		newBlockHeight = 0
		newAwaitingEventSince = 0
		return
	}

	lastUsedNonceStored, err := k.LastUsedSolanaNonce.Get(ctx, solParamsNonceAccountKey)
	if err != nil {
		k.Logger(ctx).Error("BTL Logic: Failed to get LastUsedSolanaNonce from store. Resetting transaction.", "tx_id", txID, "error", err)
		newBlockHeight = 0
		newAwaitingEventSince = 0
		return
	}

	currentLiveNonceBytes := currentLiveNonceAccount.Nonce[:]
	k.Logger(ctx).Info("BTL Logic: Comparing nonces.", "tx_id", txID, "last_used_hex", hex.EncodeToString(lastUsedNonceStored.Nonce), "current_on_chain_hex", hex.EncodeToString(currentLiveNonceBytes))

	if bytes.Equal(currentLiveNonceBytes, lastUsedNonceStored.Nonce) {
		k.Logger(ctx).Info("BTL Logic: On-chain nonce matches LastUsedSolanaNonce (no advancement). Resetting transaction for full retry.", "tx_id", txID)
		newBlockHeight = 0
		newAwaitingEventSince = 0
	} else {
		k.Logger(ctx).Info("BTL Logic: On-chain nonce differs from LastUsedSolanaNonce (advanced). Setting AwaitingEventSince if not already set.", "tx_id", txID)
		// If the transaction was not already awaiting an event, mark it now.
		if txAwaitingEventSince == 0 {
			newAwaitingEventSince = ctx.BlockHeight()
			k.Logger(ctx).Info("BTL Logic: Set AwaitingEventSince.", "tx_id", txID, "awaiting_event_since_block", newAwaitingEventSince)
		}
		// CRITICAL: newBlockHeight is NOT reset here if it was > 0. Nonce advanced, so the original dispatch "consumed" its chance.
	}
	return
}

// handleSolanaEventArrivalTimeout contains the core logic for secondary event timeout processing.
// This is called if a transaction has been marked as awaiting an event (txAwaitingEventSince > 0).
// If the event confirmation window has passed without the event being processed,
// it resets the transaction for a full retry and attempts to update the LastUsedSolanaNonce.
func (k Keeper) handleSolanaEventArrivalTimeout(
	ctx sdk.Context,
	txID uint64,
	txRecipientAddress string,
	txAmount uint64,
	txBlockHeight int64,
	txAwaitingEventSince int64,
	solParamsNonceAccountKey uint64,
	oracleData OracleData,
) (newBlockHeight int64, newAwaitingEventSince int64) {
	// Initialize return values to current tx values
	newBlockHeight = txBlockHeight
	newAwaitingEventSince = txAwaitingEventSince

	// Only proceed if the transaction was actually awaiting an event
	if txAwaitingEventSince == 0 {
		return
	}

	k.Logger(ctx).Info("Secondary Timeout Logic: Checking for event arrival.", "tx_id", txID, "awaiting_event_since", txAwaitingEventSince, "current_height", ctx.BlockHeight(), "confirmation_window", solanaEventConfirmationWindowBlocks)

	if ctx.BlockHeight() > txAwaitingEventSince+solanaEventConfirmationWindowBlocks {
		k.Logger(ctx).Warn("Secondary Timeout Logic: SolanaMintEvent not received within window. Resetting transaction for retry and attempting to update LastUsedSolanaNonce.",
			"tx_id", txID, "recipient", txRecipientAddress, "amount", txAmount,
			"awaiting_since_block", txAwaitingEventSince, "timeout_window", solanaEventConfirmationWindowBlocks)

		currentLiveNonceForRetryUpdate := oracleData.SolanaMintNonces[solParamsNonceAccountKey]
		if currentLiveNonceForRetryUpdate == nil || currentLiveNonceForRetryUpdate.Nonce.IsZero() {
			k.Logger(ctx).Warn("Secondary Timeout Logic: Current on-chain Solana nonce is zero or unavailable in oracleData. Retry will use previously stored LastUsedSolanaNonce.", "tx_id", txID)
		} else {
			newLastNonceToStore := types.SolanaNonce{Nonce: currentLiveNonceForRetryUpdate.Nonce[:]}
			if err := k.LastUsedSolanaNonce.Set(ctx, solParamsNonceAccountKey, newLastNonceToStore); err != nil {
				k.Logger(ctx).Error("Secondary Timeout Logic: Failed to update LastUsedSolanaNonce for retry. Next retry will use older LastUsedSolanaNonce.", "tx_id", txID, "error", err)
			} else {
				k.Logger(ctx).Info("Secondary Timeout Logic: Successfully updated LastUsedSolanaNonce before retry.", "tx_id", txID, "new_last_used_nonce_hex", hex.EncodeToString(newLastNonceToStore.Nonce))
			}
		}

		newBlockHeight = 0
		newAwaitingEventSince = 0
		k.Logger(ctx).Info("Secondary Timeout Logic: Transaction has been reset for a full retry.", "tx_id", txID)
	}
	return
}

// Helper function to process BTL timeout for Solana mints
func (k Keeper) processBtlSolanaMint(ctx sdk.Context, tx zenbtctypes.PendingMintTransaction, oracleData OracleData, solParams zenbtctypes.Solana) zenbtctypes.PendingMintTransaction {
	newBlockHeight, newAwaitingEventSince := k.handleSolanaTransactionBTLTimeout(
		ctx,
		tx.Id,
		tx.BlockHeight,
		tx.AwaitingEventSince,
		solParams.NonceAccountKey,
		solParams.Btl,
		oracleData,
	)
	tx.BlockHeight = newBlockHeight
	tx.AwaitingEventSince = newAwaitingEventSince
	return tx
}

// Helper function to process secondary event timeout for Solana mints
func (k Keeper) processSecondaryTimeoutSolanaMint(ctx sdk.Context, tx zenbtctypes.PendingMintTransaction, oracleData OracleData, solParams zenbtctypes.Solana) zenbtctypes.PendingMintTransaction {
	newBlockHeight, newAwaitingEventSince := k.handleSolanaEventArrivalTimeout(
		ctx,
		tx.Id,
		tx.RecipientAddress,
		tx.Amount,
		tx.BlockHeight,
		tx.AwaitingEventSince,
		solParams.NonceAccountKey,
		oracleData,
	)
	tx.BlockHeight = newBlockHeight
	tx.AwaitingEventSince = newAwaitingEventSince
	return tx
}

func (k Keeper) processSecondaryTimeoutSolanaROCKMint(ctx sdk.Context, tx zentptypes.Bridge, oracleData OracleData, solParams zentptypes.Solana) zentptypes.Bridge {
	newBlockHeight, newAwaitingEventSince := k.handleSolanaEventArrivalTimeout(
		ctx,
		tx.Id,
		tx.RecipientAddress,
		tx.Amount,
		tx.BlockHeight,
		tx.AwaitingEventSince,
		solParams.NonceAccountKey,
		oracleData,
	)
	tx.BlockHeight = newBlockHeight
	tx.AwaitingEventSince = newAwaitingEventSince
	return tx
}

// Helper function to process BTL timeout for Solana ROCK mints
func (k Keeper) processBtlSolanaROCKMint(ctx sdk.Context, tx zentptypes.Bridge, oracleData OracleData, solParams zentptypes.Solana) zentptypes.Bridge {
	newBlockHeight, newAwaitingEventSince := k.handleSolanaTransactionBTLTimeout(
		ctx,
		tx.Id,
		tx.BlockHeight,
		tx.AwaitingEventSince,
		solParams.NonceAccountKey,
		solParams.Btl,
		oracleData,
	)
	tx.BlockHeight = newBlockHeight
	tx.AwaitingEventSince = newAwaitingEventSince
	return tx
}

// ClearProcessedBackfillRequests removes backfill requests from the queue that have been successfully processed.
// It identifies processed requests by matching the eventType and checking if the transaction hash
// is present in the provided processedTxHashes map.
func (k *Keeper) ClearProcessedBackfillRequests(ctx context.Context, eventType types.EventType, processedTxHashes map[string]bool) {
	if len(processedTxHashes) == 0 {
		return // Nothing was processed, so nothing to clear.
	}

	backfillRequests, err := k.BackfillRequests.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("error getting backfill requests for cleanup", "error", err)
		}
		return // Can't clean up if we can't get the list.
	}

	if len(backfillRequests.Requests) == 0 {
		return
	}

	initialRequestCount := len(backfillRequests.Requests)

	// Filter out the requests that have been processed.
	remainingRequests := slices.DeleteFunc(backfillRequests.Requests, func(req *types.MsgTriggerEventBackfill) bool {
		// Check if the request is the correct type and if its hash is in the processed map.
		if req.EventType == eventType && processedTxHashes[req.TxHash] {
			k.Logger(ctx).Info("clearing processed backfill request", "type", eventType.String(), "tx_hash", req.TxHash)
			return true // Delete this item.
		}
		return false
	})
	backfillRequests.Requests = remainingRequests

	// If we removed any requests, update the store.
	if len(backfillRequests.Requests) < initialRequestCount {
		if err := k.SetBackfillRequests(ctx, backfillRequests); err != nil {
			k.Logger(ctx).Error("error updating backfill requests after processing events", "type", eventType.String(), "error", err)
		}
	}
}
