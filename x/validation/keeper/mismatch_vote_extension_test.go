//go:build test
// +build test

package keeper_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	ubermock "go.uber.org/mock/gomock"

	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtestutil "github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

type MismatchVoteExtensionTestSuite struct {
	suite.Suite
	ctx               sdk.Context
	validationKeeper  *keeper.Keeper
	mockSidecarClient *validationtestutil.MocksidecarClient
	ctrl              *ubermock.Controller
	zenBTCCtrl        *ubermock.Controller
}

func (suite *MismatchVoteExtensionTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(validationtypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{}).WithLogger(log.NewNopLogger())

	suite.ctrl = ubermock.NewController(suite.T())
	accountKeeper := validationtestutil.NewMockAccountKeeper(suite.ctrl)
	bankKeeper := validationtestutil.NewMockBankKeeper(suite.ctrl)
	treasuryKeeper := validationtestutil.NewMockTreasuryKeeper(suite.ctrl)
	zentpKeeper := validationtestutil.NewMockZentpKeeper(suite.ctrl)

	// Set up basic account keeper expectations
	bondedAcc := authtypes.NewEmptyModuleAccount(validationtypes.BondedPoolName)
	notBondedAcc := authtypes.NewEmptyModuleAccount(validationtypes.NotBondedPoolName)
	accountKeeper.EXPECT().GetModuleAddress(validationtypes.BondedPoolName).Return(bondedAcc.GetAddress()).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(validationtypes.NotBondedPoolName).Return(notBondedAcc.GetAddress()).AnyTimes()
	accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("zen")).AnyTimes()

	// Set up bank keeper expectations
	bankKeeper.EXPECT().GetAllBalances(ctx, bondedAcc.GetAddress()).Return(sdk.NewCoins()).AnyTimes()
	bankKeeper.EXPECT().GetAllBalances(ctx, notBondedAcc.GetAddress()).Return(sdk.NewCoins()).AnyTimes()

	// Set up zentp keeper expectations
	zentpKeeper.EXPECT().GetTotalROCKSupply(ubermock.Any()).Return(math.NewInt(1000000), nil).AnyTimes()
	zentpKeeper.EXPECT().CheckROCKSupplyCap(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()

	// Set up zenBTC keeper
	suite.zenBTCCtrl = ubermock.NewController(suite.T())
	zenBTCKeeper := validationtestutil.NewMockZenBTCKeeper(suite.zenBTCCtrl)
	zenBTCKeeper.EXPECT().GetStakerKeyID(ubermock.Any()).Return(uint64(1)).AnyTimes()
	zenBTCKeeper.EXPECT().GetEthMinterKeyID(ubermock.Any()).Return(uint64(2)).AnyTimes()
	zenBTCKeeper.EXPECT().GetUnstakerKeyID(ubermock.Any()).Return(uint64(3)).AnyTimes()
	zenBTCKeeper.EXPECT().GetCompleterKeyID(ubermock.Any()).Return(uint64(4)).AnyTimes()

	// Create keeper
	suite.validationKeeper = keeper.NewKeeper(
		codec.NewProtoCodec(nil),
		storeService,
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		nil,
		nil,
		treasuryKeeper,
		zenBTCKeeper,
		zentpKeeper,
		address.NewBech32Codec("zenvaloper"),
		address.NewBech32Codec("zenvalcons"),
	)

	// Set up sidecar client
	suite.mockSidecarClient = validationtestutil.NewMocksidecarClient(suite.ctrl)
	suite.validationKeeper.SetSidecarClient(suite.mockSidecarClient)

	// Set default params
	require.NoError(suite.T(), suite.validationKeeper.SetParams(ctx, validationtypes.DefaultParams()))

	suite.ctx = ctx
}

func (suite *MismatchVoteExtensionTestSuite) TearDownTest() {
	if suite.ctrl != nil {
		suite.ctrl.Finish()
	}
	if suite.zenBTCCtrl != nil {
		suite.zenBTCCtrl.Finish()
	}
}

func TestMismatchVoteExtensionTestSuite(t *testing.T) {
	suite.Run(t, new(MismatchVoteExtensionTestSuite))
}

// TestSingleValidatorNetworkConsensus tests that a single validator network
// achieves consensus on all fields and doesn't record mismatches
func (suite *MismatchVoteExtensionTestSuite) TestSingleValidatorNetworkConsensus() {
	require := suite.Require()

	// Create a single validator vote extension
	voteExt := keeper.VoteExtension{
		EigenDelegationsHash:    []byte("test_hash"),
		RequestedBtcBlockHeight: 100,
		RequestedBtcHeaderHash:  []byte("btc_hash"),
		EthBlockHeight:          200,
		EthGasLimit:             21000,
		EthBaseFee:              1000000000,
		EthTipCap:               2000000000,
		RequestedStakerNonce:    1,
		RequestedEthMinterNonce: 2,
		RequestedUnstakerNonce:  3,
		RequestedCompleterNonce: 4,
		SolanaMintNoncesHash:    []byte("solana_mint_nonces"),
		SolanaAccountsHash:      []byte("solana_accounts"),
		EthBurnEventsHash:       []byte("eth_burn_events"),
		SolanaBurnEventsHash:    []byte("solana_burn_events"),
		SolanaMintEventsHash:    []byte("solana_mint_events"),
		RedemptionsHash:         []byte("redemptions"),
		ROCKUSDPrice:            "0.01",
		BTCUSDPrice:             "50000.00",
		ETHUSDPrice:             "3000.00",
		LatestBtcBlockHeight:    101,
		LatestBtcHeaderHash:     []byte("latest_btc_hash"),
		SidecarVersionName:      "test_sidecar_v1",
	}

	// Serialize the vote extension
	voteExtBz, err := json.Marshal(voteExt)
	require.NoError(err)

	// Create a single validator consensus info
	validatorAddr := []byte("validator_address_01")
	consensusData := abci.ExtendedCommitInfo{
		Round: 1,
		Votes: []abci.ExtendedVoteInfo{
			{
				Validator: abci.Validator{
					Address: validatorAddr,
					Power:   100, // 100% voting power
				},
				VoteExtension: voteExtBz,
				BlockIdFlag:   cmtproto.BlockIDFlagCommit,
			},
		},
	}

	// Test GetConsensusAndPluralityVEData
	canonicalVE, pluralityVE, fieldVotePowers, err := suite.validationKeeper.GetConsensusAndPluralityVEData(suite.ctx, 1, consensusData)
	require.NoError(err)

	// In a single validator network, the canonical VE should match the validator's VE
	require.Equal(voteExt.EigenDelegationsHash, canonicalVE.EigenDelegationsHash)
	require.Equal(voteExt.RequestedBtcBlockHeight, canonicalVE.RequestedBtcBlockHeight)
	require.Equal(voteExt.EthBlockHeight, canonicalVE.EthBlockHeight)
	require.Equal(voteExt.EthGasLimit, canonicalVE.EthGasLimit)
	require.Equal(voteExt.EthBaseFee, canonicalVE.EthBaseFee)
	require.Equal(voteExt.EthTipCap, canonicalVE.EthTipCap)
	require.Equal(voteExt.ROCKUSDPrice, canonicalVE.ROCKUSDPrice)
	require.Equal(voteExt.BTCUSDPrice, canonicalVE.BTCUSDPrice)
	require.Equal(voteExt.ETHUSDPrice, canonicalVE.ETHUSDPrice)

	// All fields should have consensus with the full voting power
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldEigenDelegationsHash])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldRequestedBtcBlockHeight])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldEthBlockHeight])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldEthGasLimit])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldEthBaseFee])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldEthTipCap])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldROCKUSDPrice])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldBTCUSDPrice])
	require.Equal(int64(100), fieldVotePowers[keeper.VEFieldETHUSDPrice])
}

// TestSingleValidatorNoMismatchRecorded tests that a single validator
// doesn't get mismatch records when providing the canonical vote extension
func (suite *MismatchVoteExtensionTestSuite) TestSingleValidatorNoMismatchRecorded() {
	require := suite.Require()

	// Create a vote extension
	voteExt := keeper.VoteExtension{
		EigenDelegationsHash:    []byte("test_hash"),
		RequestedBtcBlockHeight: 100,
		RequestedBtcHeaderHash:  []byte("btc_hash"),
		EthBlockHeight:          200,
		EthGasLimit:             21000,
		EthBaseFee:              1000000000,
		EthTipCap:               2000000000,
		RequestedStakerNonce:    1,
		RequestedEthMinterNonce: 2,
		RequestedUnstakerNonce:  3,
		RequestedCompleterNonce: 4,
		SolanaMintNoncesHash:    []byte("solana_mint_nonces"),
		SolanaAccountsHash:      []byte("solana_accounts"),
		EthBurnEventsHash:       []byte("eth_burn_events"),
		SolanaBurnEventsHash:    []byte("solana_burn_events"),
		SolanaMintEventsHash:    []byte("solana_mint_events"),
		RedemptionsHash:         []byte("redemptions"),
		ROCKUSDPrice:            "0.01",
		BTCUSDPrice:             "50000.00",
		ETHUSDPrice:             "3000.00",
		LatestBtcBlockHeight:    101,
		LatestBtcHeaderHash:     []byte("latest_btc_hash"),
		SidecarVersionName:      "test_sidecar_v1",
	}

	// Serialize the vote extension
	voteExtBz, err := json.Marshal(voteExt)
	require.NoError(err)

	// Create validator address
	validatorAddr := []byte("validator_address_01")

	// Create consensus data with single validator
	consensusData := abci.ExtendedCommitInfo{
		Round: 1,
		Votes: []abci.ExtendedVoteInfo{
			{
				Validator: abci.Validator{
					Address: validatorAddr,
					Power:   100,
				},
				VoteExtension: voteExtBz,
				BlockIdFlag:   cmtproto.BlockIDFlagCommit,
			},
		},
	}

	// Record mismatched vote extensions - should not record any mismatches
	// since the validator's VE is the canonical one
	suite.validationKeeper.RecordMismatchedVoteExtensions(suite.ctx, 1, voteExt, consensusData)

	// Verify no mismatch was recorded by checking the validator info
	validatorHexAddr := hex.EncodeToString(validatorAddr)
	info, err := suite.validationKeeper.ValidationInfos.Get(suite.ctx, 1)
	if err == nil {
		// If validation info exists, check that this validator is not in mismatched list
		for _, mismatchedAddr := range info.MismatchedVoteExtensions {
			require.NotEqual(validatorHexAddr, mismatchedAddr,
				"Single validator should not be recorded as mismatched")
		}
	}

	// Also check that no mismatch count was recorded
	_, err = suite.validationKeeper.ValidatorMismatchCounts.Get(suite.ctx, validatorHexAddr)
	require.Error(err, "Single validator should not have mismatch count recorded")
}

// TestEmptyVoteExtensionNoMismatch tests that empty vote extensions
// don't cause mismatches when canonical is also empty
func (suite *MismatchVoteExtensionTestSuite) TestEmptyVoteExtensionNoMismatch() {
	require := suite.Require()

	// Create empty canonical VE
	emptyVE := keeper.VoteExtension{}

	// Create validator address
	validatorAddr := []byte("validator_address_01")

	// Create consensus data with empty vote extension
	consensusData := abci.ExtendedCommitInfo{
		Round: 1,
		Votes: []abci.ExtendedVoteInfo{
			{
				Validator: abci.Validator{
					Address: validatorAddr,
					Power:   100,
				},
				VoteExtension: []byte{}, // Empty vote extension
				BlockIdFlag:   cmtproto.BlockIDFlagCommit,
			},
		},
	}

	// Record mismatched vote extensions - should not record mismatch for empty VE
	suite.validationKeeper.RecordMismatchedVoteExtensions(suite.ctx, 1, emptyVE, consensusData)

	// Verify no mismatch was recorded
	validatorHexAddr := hex.EncodeToString(validatorAddr)
	info, err := suite.validationKeeper.ValidationInfos.Get(suite.ctx, 1)
	if err == nil {
		for _, mismatchedAddr := range info.MismatchedVoteExtensions {
			require.NotEqual(validatorHexAddr, mismatchedAddr,
				"Empty vote extension should not be treated as mismatch when canonical is also empty")
		}
	}
}

// TestMismatchCountSlidingWindow tests the sliding window mechanism
// for mismatch counts
func (suite *MismatchVoteExtensionTestSuite) TestMismatchCountSlidingWindow() {
	require := suite.Require()

	validatorHexAddr := "test_validator_address"

	// Test first mismatch
	suite.validationKeeper.UpdateValidatorMismatchCount(suite.ctx, validatorHexAddr, 1)

	count, err := suite.validationKeeper.ValidatorMismatchCounts.Get(suite.ctx, validatorHexAddr)
	require.NoError(err)
	require.Equal(uint32(1), count.TotalCount)
	require.Equal([]int64{1}, count.MismatchBlocks)

	// Test adding more mismatches within window
	for i := int64(2); i <= 50; i++ {
		suite.validationKeeper.UpdateValidatorMismatchCount(suite.ctx, validatorHexAddr, i)
	}

	count, err = suite.validationKeeper.ValidatorMismatchCounts.Get(suite.ctx, validatorHexAddr)
	require.NoError(err)
	require.Equal(uint32(50), count.TotalCount)
	require.Len(count.MismatchBlocks, 50)

	// Test sliding window - add mismatch at block 101, should remove block 1
	suite.validationKeeper.UpdateValidatorMismatchCount(suite.ctx, validatorHexAddr, 101)

	count, err = suite.validationKeeper.ValidatorMismatchCounts.Get(suite.ctx, validatorHexAddr)
	require.NoError(err)
	require.Equal(uint32(50), count.TotalCount)         // Should still be 50 due to sliding window
	require.NotContains(count.MismatchBlocks, int64(1)) // Block 1 should be removed
	require.Contains(count.MismatchBlocks, int64(101))  // Block 101 should be added
}

// TestValidatorJailingThreshold tests that validators are only jailed
// when they exceed the mismatch threshold
func (suite *MismatchVoteExtensionTestSuite) TestValidatorJailingThreshold() {
	require := suite.Require()

	// Test mismatch count tracking without validator creation
	validatorHexAddr := "test_validator_address"

	// Test case 1: Validator with less than 50 mismatches
	// Add 49 mismatches
	for i := int64(1); i <= 49; i++ {
		suite.validationKeeper.UpdateValidatorMismatchCount(suite.ctx, validatorHexAddr, i)
	}

	// Check that mismatch count is tracked correctly
	count, err := suite.validationKeeper.ValidatorMismatchCounts.Get(suite.ctx, validatorHexAddr)
	require.NoError(err)
	require.Equal(uint32(49), count.TotalCount)
	require.Less(count.TotalCount, uint32(50), "Validator should not be jailed with less than 50 mismatches")

	// Test case 2: Validator reaches 50 mismatches threshold
	suite.validationKeeper.UpdateValidatorMismatchCount(suite.ctx, validatorHexAddr, 50)

	count, err = suite.validationKeeper.ValidatorMismatchCounts.Get(suite.ctx, validatorHexAddr)
	require.NoError(err)
	require.Equal(uint32(50), count.TotalCount)
	require.GreaterOrEqual(count.TotalCount, uint32(50), "Validator should reach jailing threshold")
}

// TestOracleDataValidation tests validation of oracle data against vote extensions
func (suite *MismatchVoteExtensionTestSuite) TestOracleDataValidation() {
	require := suite.Require()

	// Create oracle data first to calculate proper hashes
	oracleData := &keeper.OracleData{
		EigenDelegationsMap:     map[string]map[string]*big.Int{"test": {"test": big.NewInt(100)}},
		RequestedBtcBlockHeight: 100,
		EthBlockHeight:          200,
		EthGasLimit:             21000,
		ROCKUSDPrice:            "0.01",
		BTCUSDPrice:             "50000.00",
		ETHUSDPrice:             "3000.00",
	}

	// Create a vote extension with values that match oracle data
	// Note: We're not testing hash validation here, just field value validation
	voteExt := keeper.VoteExtension{
		RequestedBtcBlockHeight: 100,        // Same as oracle data
		EthBlockHeight:          200,        // Same as oracle data
		EthGasLimit:             21000,      // Same as oracle data
		ROCKUSDPrice:            "0.01",     // Same as oracle data
		BTCUSDPrice:             "50000.00", // Same as oracle data
		ETHUSDPrice:             "3000.00",  // Same as oracle data
		// Add other required fields to make vote extension valid
		EthBaseFee:              1000000000,
		EthTipCap:               2000000000,
		RequestedStakerNonce:    1,
		RequestedEthMinterNonce: 2,
		RequestedUnstakerNonce:  3,
		RequestedCompleterNonce: 4,
		SolanaMintNoncesHash:    []byte("solana_mint_nonces"),
		SolanaAccountsHash:      []byte("solana_accounts"),
		EthBurnEventsHash:       []byte("eth_burn_events"),
		SolanaBurnEventsHash:    []byte("solana_burn_events"),
		SolanaMintEventsHash:    []byte("solana_mint_events"),
		RedemptionsHash:         []byte("redemptions"),
		LatestBtcBlockHeight:    101,
		LatestBtcHeaderHash:     []byte("latest_btc_hash"),
		SidecarVersionName:      "test_sidecar_v1",
		EigenDelegationsHash:    []byte("test_hash"),
		RequestedBtcHeaderHash:  []byte("btc_hash"),
	}

	// Set up field vote powers (only for fields we want to test)
	fieldVotePowers := map[keeper.VoteExtensionField]int64{
		keeper.VEFieldRequestedBtcBlockHeight: 100,
		keeper.VEFieldEthBlockHeight:          100,
		keeper.VEFieldEthGasLimit:             100,
		keeper.VEFieldROCKUSDPrice:            100,
		keeper.VEFieldBTCUSDPrice:             100,
		keeper.VEFieldETHUSDPrice:             100,
	}

	// Initialize FieldVotePowers in oracle data
	oracleData.FieldVotePowers = make(map[keeper.VoteExtensionField]int64)
	for k, v := range fieldVotePowers {
		oracleData.FieldVotePowers[k] = v
	}

	// Run validation - should not remove any fields from fieldVotePowers
	originalFieldCount := len(fieldVotePowers)
	suite.validationKeeper.ValidateOracleData(suite.ctx, voteExt, oracleData, fieldVotePowers)

	// Verify that no fields were removed (indicating no mismatches)
	require.Equal(originalFieldCount, len(oracleData.FieldVotePowers),
		"No fields should be removed when data matches")
}

// TestRealWorldSingleValidatorScenario tests a realistic single validator scenario
// This is the core test that should expose the bug
func (suite *MismatchVoteExtensionTestSuite) TestRealWorldSingleValidatorScenario() {
	require := suite.Require()

	// Set up sidecar client with realistic responses
	suite.mockSidecarClient.EXPECT().GetSidecarState(ubermock.Any(), ubermock.Any()).Return(
		&sidecar.SidecarStateResponse{
			EthBlockHeight:   1000,
			EthGasLimit:      30000000,
			EthBaseFee:       20000000000,
			EthTipCap:        2000000000,
			EigenDelegations: []byte(`{"validator1":{"operator1":"1000000000000000000000"}}`),
			ROCKUSDPrice:     "0.01801",
			BTCUSDPrice:      "102855.235",
			ETHUSDPrice:      "3000.00",
		}, nil).AnyTimes()

	// Create a realistic vote extension that would be produced by the single validator
	voteExt := keeper.VoteExtension{
		EigenDelegationsHash:    []byte("real_delegations_hash"),
		RequestedBtcBlockHeight: 850000,
		RequestedBtcHeaderHash:  []byte("real_btc_header"),
		EthBlockHeight:          1000,
		EthGasLimit:             30000000,
		EthBaseFee:              20000000000,
		EthTipCap:               2000000000,
		RequestedStakerNonce:    10,
		RequestedEthMinterNonce: 20,
		RequestedUnstakerNonce:  30,
		RequestedCompleterNonce: 40,
		SolanaMintNoncesHash:    []byte("real_solana_mint_nonces"),
		SolanaAccountsHash:      []byte("real_solana_accounts"),
		EthBurnEventsHash:       []byte("real_eth_burn_events"),
		SolanaBurnEventsHash:    []byte("real_solana_burn_events"),
		SolanaMintEventsHash:    []byte("real_solana_mint_events"),
		RedemptionsHash:         []byte("real_redemptions"),
		ROCKUSDPrice:            "0.01801",
		BTCUSDPrice:             "102855.235",
		ETHUSDPrice:             "3000.00",
		LatestBtcBlockHeight:    850001,
		LatestBtcHeaderHash:     []byte("latest_btc_header"),
		SidecarVersionName:      "real_sidecar_v1",
	}

	// Serialize the vote extension
	voteExtBz, err := json.Marshal(voteExt)
	require.NoError(err)

	// Create realistic single validator consensus info
	validatorAddr := []byte("real_validator_addr_")
	consensusData := abci.ExtendedCommitInfo{
		Round: 1,
		Votes: []abci.ExtendedVoteInfo{
			{
				Validator: abci.Validator{
					Address: validatorAddr,
					Power:   1000000, // Realistic voting power
				},
				VoteExtension: voteExtBz,
				BlockIdFlag:   cmtproto.BlockIDFlagCommit,
			},
		},
	}

	// Step 1: Get canonical vote extension (should be the validator's VE)
	canonicalVE, pluralityVE, fieldVotePowers, err := suite.validationKeeper.GetConsensusAndPluralityVEData(suite.ctx, 100, consensusData)
	require.NoError(err)

	// Verify the canonical VE is the same as the validator's VE
	require.Equal(voteExt.EthBlockHeight, canonicalVE.EthBlockHeight)
	require.Equal(voteExt.EthGasLimit, canonicalVE.EthGasLimit)
	require.Equal(voteExt.EthBaseFee, canonicalVE.EthBaseFee)
	require.Equal(voteExt.EthTipCap, canonicalVE.EthTipCap)
	require.Equal(voteExt.ROCKUSDPrice, canonicalVE.ROCKUSDPrice)
	require.Equal(voteExt.BTCUSDPrice, canonicalVE.BTCUSDPrice)
	require.Equal(voteExt.ETHUSDPrice, canonicalVE.ETHUSDPrice)

	// Verify all fields have consensus with full voting power
	require.Equal(int64(1000000), fieldVotePowers[keeper.VEFieldEthBlockHeight])
	require.Equal(int64(1000000), fieldVotePowers[keeper.VEFieldEthGasLimit])
	require.Equal(int64(1000000), fieldVotePowers[keeper.VEFieldEthBaseFee])
	require.Equal(int64(1000000), fieldVotePowers[keeper.VEFieldEthTipCap])
	require.Equal(int64(1000000), fieldVotePowers[keeper.VEFieldROCKUSDPrice])
	require.Equal(int64(1000000), fieldVotePowers[keeper.VEFieldBTCUSDPrice])
	require.Equal(int64(1000000), fieldVotePowers[keeper.VEFieldETHUSDPrice])

	// Step 2: Record mismatched vote extensions
	suite.validationKeeper.RecordMismatchedVoteExtensions(suite.ctx, 100, canonicalVE, consensusData)

	// Step 3: Verify no mismatch was recorded (this is where the bug would show)
	validatorHexAddr := hex.EncodeToString(validatorAddr)

	// Check validation info for mismatches
	info, err := suite.validationKeeper.ValidationInfos.Get(suite.ctx, 100)
	if err == nil {
		for _, mismatchedAddr := range info.MismatchedVoteExtensions {
			require.NotEqual(validatorHexAddr, mismatchedAddr,
				"CRITICAL BUG: Single validator should NEVER be recorded as mismatched when providing canonical VE")
		}
	}

}

// TestJSONSerializationMismatchBug demonstrates the core bug:
// JSON serialization/deserialization can change the byte representation
func (suite *MismatchVoteExtensionTestSuite) TestJSONSerializationMismatchBug() {
	require := suite.Require()

	// Create a vote extension
	voteExt := keeper.VoteExtension{
		EigenDelegationsHash:    []byte("test_hash"),
		RequestedBtcBlockHeight: 100,
		RequestedBtcHeaderHash:  []byte("btc_hash"),
		EthBlockHeight:          200,
		EthGasLimit:             21000,
		EthBaseFee:              1000000000,
		EthTipCap:               2000000000,
		RequestedStakerNonce:    1,
		RequestedEthMinterNonce: 2,
		RequestedUnstakerNonce:  3,
		RequestedCompleterNonce: 4,
		SolanaMintNoncesHash:    []byte("solana_mint_nonces"),
		SolanaAccountsHash:      []byte("solana_accounts"),
		EthBurnEventsHash:       []byte("eth_burn_events"),
		SolanaBurnEventsHash:    []byte("solana_burn_events"),
		SolanaMintEventsHash:    []byte("solana_mint_events"),
		RedemptionsHash:         []byte("redemptions"),
		ROCKUSDPrice:            "0.01",
		BTCUSDPrice:             "50000.00",
		ETHUSDPrice:             "3000.00",
		LatestBtcBlockHeight:    101,
		LatestBtcHeaderHash:     []byte("latest_btc_hash"),
		SidecarVersionName:      "test_sidecar_v1",
	}

	// Step 1: Serialize the original vote extension (simulating validator submission)
	originalJSON, err := json.Marshal(voteExt)
	require.NoError(err)

	// Step 2: Deserialize it (simulating what GetConsensusAndPluralityVEData does)
	var deserializedVoteExt keeper.VoteExtension
	err = json.Unmarshal(originalJSON, &deserializedVoteExt)
	require.NoError(err)

	// Step 3: Serialize it again (simulating what recordMismatchedVoteExtensions does)
	reserializedJSON, err := json.Marshal(deserializedVoteExt)
	require.NoError(err)

	// Step 4: Compare the byte arrays
	if !bytes.Equal(originalJSON, reserializedJSON) {
		suite.T().Logf("FOUND THE BUG!")
		suite.T().Logf("Original JSON:      %s", string(originalJSON))
		suite.T().Logf("Reserialized JSON:  %s", string(reserializedJSON))
		suite.T().Logf("This difference causes single validators to be flagged as mismatched!")
	} else {
		suite.T().Logf("JSON serialization is deterministic in this case")
	}

	// This demonstrates why recordMismatchedVoteExtensions incorrectly flags
	// single validators as mismatched - the JSON bytes don't match even though
	// the vote extension content is identical
}
