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

// GetSuperMajorityVEData tallies votes for individual fields of the VoteExtension instead of requiring
// consensus on the entire object. This makes the system more resilient by allowing fields
// that have supermajority consensus to be accepted even if other fields don't reach consensus.
func (k Keeper) GetSuperMajorityVEData(ctx context.Context, currentHeight int64, extCommit abci.ExtendedCommitInfo) (VoteExtension, map[VoteExtensionField]int64, int64, error) {
	// Maps to store votes for each field
	eigenDelegationsHashVotes := make(map[string]fieldVote)
	ethBurnEventsHashVotes := make(map[string]fieldVote)
	redemptionsHashVotes := make(map[string]fieldVote)
	requestedBtcHeaderHashVotes := make(map[string]fieldVote)
	requestedBtcBlockHeightVotes := make(map[int64]fieldVote)
	ethBlockHeightVotes := make(map[uint64]fieldVote)
	ethGasLimitVotes := make(map[uint64]fieldVote)
	ethBaseFeeVotes := make(map[uint64]fieldVote)
	ethTipCapVotes := make(map[uint64]fieldVote)
	solanaLamportsPerSignatureVotes := make(map[uint64]fieldVote)
	requestedStakerNonceVotes := make(map[uint64]fieldVote)
	requestedEthMinterNonceVotes := make(map[uint64]fieldVote)
	requestedUnstakerNonceVotes := make(map[uint64]fieldVote)
	requestedCompleterNonceVotes := make(map[uint64]fieldVote)
	rockUSDPriceVotes := make(map[string]fieldVote)
	btcUSDPriceVotes := make(map[string]fieldVote)
	ethUSDPriceVotes := make(map[string]fieldVote)
	latestBtcBlockHeightVotes := make(map[int64]fieldVote)
	latestBtcHeaderHashVotes := make(map[string]fieldVote)

	var totalVotePower int64
	fieldVotePowers := make(map[VoteExtensionField]int64)

	// Process all votes
	for _, vote := range extCommit.Votes {
		totalVotePower += vote.Validator.Power

		voteExt, err := k.validateVote(ctx, vote, currentHeight)
		if err != nil {
			continue
		}

		// Tally votes for each field
		tallyFieldVote(eigenDelegationsHashVotes, bytesToString(voteExt.EigenDelegationsHash), voteExt.EigenDelegationsHash, vote.Validator.Power)
		tallyFieldVote(ethBurnEventsHashVotes, bytesToString(voteExt.EthBurnEventsHash), voteExt.EthBurnEventsHash, vote.Validator.Power)
		tallyFieldVote(redemptionsHashVotes, bytesToString(voteExt.RedemptionsHash), voteExt.RedemptionsHash, vote.Validator.Power)
		tallyFieldVote(requestedBtcHeaderHashVotes, bytesToString(voteExt.RequestedBtcHeaderHash), voteExt.RequestedBtcHeaderHash, vote.Validator.Power)
		tallyFieldVote(requestedBtcBlockHeightVotes, voteExt.RequestedBtcBlockHeight, voteExt.RequestedBtcBlockHeight, vote.Validator.Power)
		tallyFieldVote(ethBlockHeightVotes, voteExt.EthBlockHeight, voteExt.EthBlockHeight, vote.Validator.Power)
		tallyFieldVote(ethGasLimitVotes, voteExt.EthGasLimit, voteExt.EthGasLimit, vote.Validator.Power)
		tallyFieldVote(ethBaseFeeVotes, voteExt.EthBaseFee, voteExt.EthBaseFee, vote.Validator.Power)
		tallyFieldVote(ethTipCapVotes, voteExt.EthTipCap, voteExt.EthTipCap, vote.Validator.Power)
		tallyFieldVote(solanaLamportsPerSignatureVotes, voteExt.SolanaLamportsPerSignature, voteExt.SolanaLamportsPerSignature, vote.Validator.Power)
		tallyFieldVote(requestedStakerNonceVotes, voteExt.RequestedStakerNonce, voteExt.RequestedStakerNonce, vote.Validator.Power)
		tallyFieldVote(requestedEthMinterNonceVotes, voteExt.RequestedEthMinterNonce, voteExt.RequestedEthMinterNonce, vote.Validator.Power)
		tallyFieldVote(requestedUnstakerNonceVotes, voteExt.RequestedUnstakerNonce, voteExt.RequestedUnstakerNonce, vote.Validator.Power)
		tallyFieldVote(requestedCompleterNonceVotes, voteExt.RequestedCompleterNonce, voteExt.RequestedCompleterNonce, vote.Validator.Power)
		tallyFieldVote(rockUSDPriceVotes, voteExt.ROCKUSDPrice.String(), voteExt.ROCKUSDPrice, vote.Validator.Power)
		tallyFieldVote(btcUSDPriceVotes, voteExt.BTCUSDPrice.String(), voteExt.BTCUSDPrice, vote.Validator.Power)
		tallyFieldVote(ethUSDPriceVotes, voteExt.ETHUSDPrice.String(), voteExt.ETHUSDPrice, vote.Validator.Power)
		tallyFieldVote(latestBtcBlockHeightVotes, voteExt.LatestBtcBlockHeight, voteExt.LatestBtcBlockHeight, vote.Validator.Power)
		tallyFieldVote(latestBtcHeaderHashVotes, bytesToString(voteExt.LatestBtcHeaderHash), voteExt.LatestBtcHeaderHash, vote.Validator.Power)
	}

	// Create consensus VoteExtension with fields that have supermajority
	var consensusVE VoteExtension
	consensusVE.ZRChainBlockHeight = currentHeight - 1

	// Check for consensus on each field and use the most voted value if it has supermajority
	var requiredVotePower = requisiteVotePower(totalVotePower)

	// Handle EigenDelegationsHash
	if mostVoted, votePower := getMostVotedField(eigenDelegationsHashVotes); votePower >= requiredVotePower {
		consensusVE.EigenDelegationsHash = mostVoted.([]byte)
		fieldVotePowers[VEFieldEigenDelegationsHash] = votePower
	}

	// Handle EthBurnEventsHash
	if mostVoted, votePower := getMostVotedField(ethBurnEventsHashVotes); votePower >= requiredVotePower {
		consensusVE.EthBurnEventsHash = mostVoted.([]byte)
		fieldVotePowers[VEFieldEthBurnEventsHash] = votePower
	}

	// Handle RedemptionsHash
	if mostVoted, votePower := getMostVotedField(redemptionsHashVotes); votePower >= requiredVotePower {
		consensusVE.RedemptionsHash = mostVoted.([]byte)
		fieldVotePowers[VEFieldRedemptionsHash] = votePower
	}

	// Handle RequestedBtcHeaderHash
	if mostVoted, votePower := getMostVotedField(requestedBtcHeaderHashVotes); votePower >= requiredVotePower {
		consensusVE.RequestedBtcHeaderHash = mostVoted.([]byte)
		fieldVotePowers[VEFieldRequestedBtcHeaderHash] = votePower
	}

	// Handle RequestedBtcBlockHeight
	if mostVoted, votePower := getMostVotedField(requestedBtcBlockHeightVotes); votePower >= requiredVotePower {
		consensusVE.RequestedBtcBlockHeight = mostVoted.(int64)
		fieldVotePowers[VEFieldRequestedBtcBlockHeight] = votePower
	}

	// Handle LatestBtcBlockHeight
	if mostVoted, votePower := getMostVotedField(latestBtcBlockHeightVotes); votePower >= requiredVotePower {
		consensusVE.LatestBtcBlockHeight = mostVoted.(int64)
		fieldVotePowers[VEFieldLatestBtcBlockHeight] = votePower
	}

	// Handle LatestBtcHeaderHash
	if mostVoted, votePower := getMostVotedField(latestBtcHeaderHashVotes); votePower >= requiredVotePower {
		consensusVE.LatestBtcHeaderHash = mostVoted.([]byte)
		fieldVotePowers[VEFieldLatestBtcHeaderHash] = votePower
	}

	// Handle EthBlockHeight
	if mostVoted, votePower := getMostVotedField(ethBlockHeightVotes); votePower >= requiredVotePower {
		consensusVE.EthBlockHeight = mostVoted.(uint64)
		fieldVotePowers[VEFieldEthBlockHeight] = votePower
	}

	// Handle EthGasLimit
	if mostVoted, votePower := getMostVotedField(ethGasLimitVotes); votePower >= requiredVotePower {
		consensusVE.EthGasLimit = mostVoted.(uint64)
		fieldVotePowers[VEFieldEthGasLimit] = votePower
	}

	// Handle EthBaseFee
	if mostVoted, votePower := getMostVotedField(ethBaseFeeVotes); votePower >= requiredVotePower {
		consensusVE.EthBaseFee = mostVoted.(uint64)
		fieldVotePowers[VEFieldEthBaseFee] = votePower
	}

	// Handle EthTipCap
	if mostVoted, votePower := getMostVotedField(ethTipCapVotes); votePower >= requiredVotePower {
		consensusVE.EthTipCap = mostVoted.(uint64)
		fieldVotePowers[VEFieldEthTipCap] = votePower
	}

	// Handle SolanaLamportsPerSignature
	if mostVoted, votePower := getMostVotedField(solanaLamportsPerSignatureVotes); votePower >= requiredVotePower {
		consensusVE.SolanaLamportsPerSignature = mostVoted.(uint64)
		fieldVotePowers[VEFieldSolanaLamportsPerSignature] = votePower
	}

	// Handle RequestedStakerNonce
	if mostVoted, votePower := getMostVotedField(requestedStakerNonceVotes); votePower >= requiredVotePower {
		consensusVE.RequestedStakerNonce = mostVoted.(uint64)
		fieldVotePowers[VEFieldRequestedStakerNonce] = votePower
	}

	// Handle RequestedEthMinterNonce
	if mostVoted, votePower := getMostVotedField(requestedEthMinterNonceVotes); votePower >= requiredVotePower {
		consensusVE.RequestedEthMinterNonce = mostVoted.(uint64)
		fieldVotePowers[VEFieldRequestedEthMinterNonce] = votePower
	}

	// Handle RequestedUnstakerNonce
	if mostVoted, votePower := getMostVotedField(requestedUnstakerNonceVotes); votePower >= requiredVotePower {
		consensusVE.RequestedUnstakerNonce = mostVoted.(uint64)
		fieldVotePowers[VEFieldRequestedUnstakerNonce] = votePower
	}

	// Handle RequestedCompleterNonce
	if mostVoted, votePower := getMostVotedField(requestedCompleterNonceVotes); votePower >= requiredVotePower {
		consensusVE.RequestedCompleterNonce = mostVoted.(uint64)
		fieldVotePowers[VEFieldRequestedCompleterNonce] = votePower
	}

	// Handle ROCKUSDPrice
	if mostVoted, votePower := getMostVotedField(rockUSDPriceVotes); votePower >= requiredVotePower {
		consensusVE.ROCKUSDPrice = mostVoted.(math.LegacyDec)
		fieldVotePowers[VEFieldROCKUSDPrice] = votePower
	}

	// Handle BTCUSDPrice
	if mostVoted, votePower := getMostVotedField(btcUSDPriceVotes); votePower >= requiredVotePower {
		consensusVE.BTCUSDPrice = mostVoted.(math.LegacyDec)
		fieldVotePowers[VEFieldBTCUSDPrice] = votePower
	}

	// Handle ETHUSDPrice
	if mostVoted, votePower := getMostVotedField(ethUSDPriceVotes); votePower >= requiredVotePower {
		consensusVE.ETHUSDPrice = mostVoted.(math.LegacyDec)
		fieldVotePowers[VEFieldETHUSDPrice] = votePower
	}

	// Check if all essential fields have consensus
	hasEssentialFields := HasAllEssentialFields(fieldVotePowers)

	// Log which fields have reached consensus
	if len(fieldVotePowers) > 0 {
		essentialFieldsMsg := "all essential fields have consensus"
		if !hasEssentialFields {
			essentialFieldsMsg = "missing consensus on some essential fields"

			// Log which essential fields are missing
			missingFields := []string{}
			for _, field := range EssentialVoteExtensionFields {
				if _, ok := fieldVotePowers[field]; !ok {
					missingFields = append(missingFields, field.String())
				}
			}

			k.Logger(ctx).Warn("missing consensus on essential vote extension fields",
				"missing_fields", strings.Join(missingFields, ", "))
		}

		k.Logger(ctx).Info("consensus reached on vote extension fields",
			"fields_with_consensus", len(fieldVotePowers),
			"total_fields", 17,
			"essential_fields_status", essentialFieldsMsg)

		for field, power := range fieldVotePowers {
			k.Logger(ctx).Debug("field consensus details",
				"field", field.String(),
				"vote_power", power,
				"required_power", requiredVotePower,
				"is_essential", IsEssentialField(field))
		}
	} else {
		k.Logger(ctx).Warn("no consensus reached on any vote extension fields")
	}

	return consensusVE, fieldVotePowers, totalVotePower, nil
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

// BitcoinHeadersResponse combines the latest and requested Bitcoin headers
type BitcoinHeadersResponse struct {
	Latest       *sidecar.BitcoinBlockHeaderResponse
	Requested    *sidecar.BitcoinBlockHeaderResponse
	HasRequested bool
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

	if err := k.AssetPrices.Set(ctx, types.Asset_ROCK, oracleData.ROCKUSDPrice); err != nil {
		k.Logger(ctx).Error("error setting ROCK price", "height", ctx.BlockHeight(), "err", err)
	}

	if err := k.AssetPrices.Set(ctx, types.Asset_BTC, oracleData.BTCUSDPrice); err != nil {
		k.Logger(ctx).Error("error setting BTC price", "height", ctx.BlockHeight(), "err", err)
	}

	if err := k.AssetPrices.Set(ctx, types.Asset_ETH, oracleData.ETHUSDPrice); err != nil {
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
// Only fields that have reached consensus (present in fieldVotePowers) are validated.
func (k *Keeper) validateOracleData(voteExt VoteExtension, oracleData *OracleData, fieldVotePowers map[VoteExtensionField]int64) error {
	var validationErrors []string

	// Validate hashes only if fields have consensus
	if _, ok := fieldVotePowers[VEFieldEigenDelegationsHash]; ok {
		if err := validateHashField(VEFieldEigenDelegationsHash.String(), voteExt.EigenDelegationsHash, oracleData.EigenDelegationsMap); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}
	}

	if _, ok := fieldVotePowers[VEFieldEthBurnEventsHash]; ok {
		if err := validateHashField(VEFieldEthBurnEventsHash.String(), voteExt.EthBurnEventsHash, oracleData.EthBurnEvents); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}
	}

	if _, ok := fieldVotePowers[VEFieldRedemptionsHash]; ok {
		if err := validateHashField(VEFieldRedemptionsHash.String(), voteExt.RedemptionsHash, oracleData.Redemptions); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}
	}

	if _, ok := fieldVotePowers[VEFieldRequestedBtcHeaderHash]; ok {
		if err := validateHashField(VEFieldRequestedBtcHeaderHash.String(), voteExt.RequestedBtcHeaderHash, &oracleData.RequestedBtcBlockHeader); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}
	}

	// Check Ethereum-related fields
	if _, ok := fieldVotePowers[VEFieldEthBlockHeight]; ok {
		if voteExt.EthBlockHeight != oracleData.EthBlockHeight {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldEthBlockHeight.String(), voteExt.EthBlockHeight, oracleData.EthBlockHeight))
		}
	}

	if _, ok := fieldVotePowers[VEFieldEthGasLimit]; ok {
		if voteExt.EthGasLimit != oracleData.EthGasLimit {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldEthGasLimit.String(), voteExt.EthGasLimit, oracleData.EthGasLimit))
		}
	}

	if _, ok := fieldVotePowers[VEFieldEthBaseFee]; ok {
		if voteExt.EthBaseFee != oracleData.EthBaseFee {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldEthBaseFee.String(), voteExt.EthBaseFee, oracleData.EthBaseFee))
		}
	}

	if _, ok := fieldVotePowers[VEFieldEthTipCap]; ok {
		if voteExt.EthTipCap != oracleData.EthTipCap {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldEthTipCap.String(), voteExt.EthTipCap, oracleData.EthTipCap))
		}
	}

	// Check Bitcoin height
	if _, ok := fieldVotePowers[VEFieldRequestedBtcBlockHeight]; ok {
		if voteExt.RequestedBtcBlockHeight != oracleData.RequestedBtcBlockHeight {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldRequestedBtcBlockHeight.String(), voteExt.RequestedBtcBlockHeight, oracleData.RequestedBtcBlockHeight))
		}
	}

	// Check nonce-related fields
	if _, ok := fieldVotePowers[VEFieldRequestedStakerNonce]; ok {
		if voteExt.RequestedStakerNonce != oracleData.RequestedStakerNonce {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldRequestedStakerNonce.String(), voteExt.RequestedStakerNonce, oracleData.RequestedStakerNonce))
		}
	}

	if _, ok := fieldVotePowers[VEFieldRequestedEthMinterNonce]; ok {
		if voteExt.RequestedEthMinterNonce != oracleData.RequestedEthMinterNonce {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldRequestedEthMinterNonce.String(), voteExt.RequestedEthMinterNonce, oracleData.RequestedEthMinterNonce))
		}
	}

	if _, ok := fieldVotePowers[VEFieldRequestedUnstakerNonce]; ok {
		if voteExt.RequestedUnstakerNonce != oracleData.RequestedUnstakerNonce {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldRequestedUnstakerNonce.String(), voteExt.RequestedUnstakerNonce, oracleData.RequestedUnstakerNonce))
		}
	}

	if _, ok := fieldVotePowers[VEFieldRequestedCompleterNonce]; ok {
		if voteExt.RequestedCompleterNonce != oracleData.RequestedCompleterNonce {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldRequestedCompleterNonce.String(), voteExt.RequestedCompleterNonce, oracleData.RequestedCompleterNonce))
		}
	}

	// Check price fields
	if _, ok := fieldVotePowers[VEFieldROCKUSDPrice]; ok {
		if !voteExt.ROCKUSDPrice.Equal(oracleData.ROCKUSDPrice) {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %s, got %s",
					VEFieldROCKUSDPrice.String(), voteExt.ROCKUSDPrice, oracleData.ROCKUSDPrice))
		}
	}

	if _, ok := fieldVotePowers[VEFieldBTCUSDPrice]; ok {
		if !voteExt.BTCUSDPrice.Equal(oracleData.BTCUSDPrice) {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %s, got %s",
					VEFieldBTCUSDPrice.String(), voteExt.BTCUSDPrice, oracleData.BTCUSDPrice))
		}
	}

	if _, ok := fieldVotePowers[VEFieldETHUSDPrice]; ok {
		if !voteExt.ETHUSDPrice.Equal(oracleData.ETHUSDPrice) {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %s, got %s",
					VEFieldETHUSDPrice.String(), voteExt.ETHUSDPrice, oracleData.ETHUSDPrice))
		}
	}

	// Check Latest Bitcoin height and hash fields
	if _, ok := fieldVotePowers[VEFieldLatestBtcBlockHeight]; ok {
		if voteExt.LatestBtcBlockHeight != oracleData.LatestBtcBlockHeight {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s mismatch, expected %d, got %d",
					VEFieldLatestBtcBlockHeight.String(), voteExt.LatestBtcBlockHeight, oracleData.LatestBtcBlockHeight))
		}
	}

	if _, ok := fieldVotePowers[VEFieldLatestBtcHeaderHash]; ok {
		if err := validateHashField(VEFieldLatestBtcHeaderHash.String(), voteExt.LatestBtcHeaderHash, &oracleData.LatestBtcBlockHeader); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("oracle data validation failed: %s", strings.Join(validationErrors, "; "))
	}

	return nil
}

// bytesToString converts a byte slice to a string for use as a map key
func bytesToString(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// tallyFieldVote adds a vote for a field to the appropriate map
func tallyFieldVote[K comparable, V any](votes map[K]fieldVote, key K, value V, votePower int64) {
	if existingVote, ok := votes[key]; ok {
		existingVote.votePower += votePower
		votes[key] = existingVote
	} else {
		votes[key] = fieldVote{
			value:     value,
			votePower: votePower,
		}
	}
}

// getMostVotedField returns the most voted value and its vote power for a field
func getMostVotedField[K comparable](votes map[K]fieldVote) (interface{}, int64) {
	var mostVotedValue interface{}
	var maxVotePower int64

	for _, vote := range votes {
		if vote.votePower > maxVotePower {
			maxVotePower = vote.votePower
			mostVotedValue = vote.value
		}
	}

	return mostVotedValue, maxVotePower
}

// HasRequiredGasFields checks if all essential gas/fee related fields have reached consensus
func HasRequiredGasFields(fieldVotePowers map[VoteExtensionField]int64) bool {
	requiredFields := []VoteExtensionField{
		VEFieldEthGasLimit,
		VEFieldEthBaseFee,
		VEFieldEthTipCap,
	}

	for _, field := range requiredFields {
		if _, ok := fieldVotePowers[field]; !ok {
			return false
		}
	}
	return true
}

// HasRequiredField checks if the specific field has reached consensus
func HasRequiredField(fieldVotePowers map[VoteExtensionField]int64, field VoteExtensionField) bool {
	_, ok := fieldVotePowers[field]
	return ok
}
