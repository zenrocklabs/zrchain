//go:build test
// +build test

package keeper_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	ubermock "go.uber.org/mock/gomock"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtestutil "github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

// TestMismatchDiagnostic provides detailed diagnostics to understand the root cause
// of why single validators are being flagged as mismatched
func TestMismatchDiagnostic(t *testing.T) {
	// Set up test environment
	key := storetypes.NewKVStoreKey(validationtypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{}).WithLogger(log.NewNopLogger())

	ctrl := ubermock.NewController(t)
	defer ctrl.Finish()

	// Set up mock keepers with minimal expectations
	accountKeeper := validationtestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := validationtestutil.NewMockBankKeeper(ctrl)
	treasuryKeeper := validationtestutil.NewMockTreasuryKeeper(ctrl)
	zentpKeeper := validationtestutil.NewMockZentpKeeper(ctrl)
	zenBTCKeeper := validationtestutil.NewMockZenBTCKeeper(ctrl)

	bondedAcc := authtypes.NewEmptyModuleAccount(validationtypes.BondedPoolName)
	notBondedAcc := authtypes.NewEmptyModuleAccount(validationtypes.NotBondedPoolName)
	accountKeeper.EXPECT().GetModuleAddress(validationtypes.BondedPoolName).Return(bondedAcc.GetAddress()).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(validationtypes.NotBondedPoolName).Return(notBondedAcc.GetAddress()).AnyTimes()
	accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("zen")).AnyTimes()
	bankKeeper.EXPECT().GetAllBalances(gomock.Any(), gomock.Any()).Return(sdk.NewCoins()).AnyTimes()
	zentpKeeper.EXPECT().GetTotalROCKSupply(gomock.Any()).Return(math.NewInt(1000000), nil).AnyTimes()
	zentpKeeper.EXPECT().CheckROCKSupplyCap(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetStakerKeyID(gomock.Any()).Return(uint64(1)).AnyTimes()
	zenBTCKeeper.EXPECT().GetEthMinterKeyID(gomock.Any()).Return(uint64(2)).AnyTimes()
	zenBTCKeeper.EXPECT().GetUnstakerKeyID(gomock.Any()).Return(uint64(3)).AnyTimes()
	zenBTCKeeper.EXPECT().GetCompleterKeyID(gomock.Any()).Return(uint64(4)).AnyTimes()

	// Create keeper
	validationKeeper := keeper.NewKeeper(
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

	require.NoError(t, validationKeeper.SetParams(ctx, validationtypes.DefaultParams()))

	// Set up sidecar client
	mockSidecarClient := validationtestutil.NewMocksidecarClient(ctrl)
	validationKeeper.SetSidecarClient(mockSidecarClient)

	t.Log("=== DIAGNOSTIC TEST: Single Validator Mismatch Bug ===")

	// Step 1: Create a realistic vote extension
	originalVoteExt := keeper.VoteExtension{
		EigenDelegationsHash:    []byte("diagnostic_delegations_hash"),
		RequestedBtcBlockHeight: 850000,
		RequestedBtcHeaderHash:  []byte("diagnostic_btc_header"),
		EthBlockHeight:          1000,
		EthGasLimit:             30000000,
		EthBaseFee:              20000000000,
		EthTipCap:               2000000000,
		RequestedStakerNonce:    10,
		RequestedEthMinterNonce: 20,
		RequestedUnstakerNonce:  30,
		RequestedCompleterNonce: 40,
		SolanaMintNoncesHash:    []byte("diagnostic_solana_mint_nonces"),
		SolanaAccountsHash:      []byte("diagnostic_solana_accounts"),
		EthBurnEventsHash:       []byte("diagnostic_eth_burn_events"),
		SolanaBurnEventsHash:    []byte("diagnostic_solana_burn_events"),
		SolanaMintEventsHash:    []byte("diagnostic_solana_mint_events"),
		RedemptionsHash:         []byte("diagnostic_redemptions"),
		ROCKUSDPrice:            "0.01801",
		BTCUSDPrice:             "102855.235",
		ETHUSDPrice:             "3000.00",
		LatestBtcBlockHeight:    850001,
		LatestBtcHeaderHash:     []byte("diagnostic_latest_btc_header"),
		SidecarVersionName:      "diagnostic_sidecar_v1",
	}

	// Step 2: Serialize the vote extension (this is what the validator submits)
	originalVoteExtBytes, err := json.Marshal(originalVoteExt)
	require.NoError(t, err)

	t.Logf("Original Vote Extension Size: %d bytes", len(originalVoteExtBytes))
	t.Logf("Original Vote Extension JSON: %s", string(originalVoteExtBytes))

	// Step 3: Create single validator consensus data
	validatorAddr := []byte("diagnostic_validator_addr")
	validatorHexAddr := hex.EncodeToString(validatorAddr)

	consensusData := abci.ExtendedCommitInfo{
		Round: 1,
		Votes: []abci.ExtendedVoteInfo{
			{
				Validator: abci.Validator{
					Address: validatorAddr,
					Power:   1000000, // Single validator with all voting power
				},
				VoteExtension: originalVoteExtBytes,
				BlockIdFlag:   cmtproto.BlockIDFlagCommit,
			},
		},
	}

	t.Logf("Validator Address: %s", validatorHexAddr)
	t.Logf("Validator Power: %d", consensusData.Votes[0].Validator.Power)

	// Step 4: Process vote extension through GetConsensusAndPluralityVEData
	t.Log("\n=== STEP 4: Processing through GetConsensusAndPluralityVEData ===")

	canonicalVE, pluralityVE, fieldVotePowers, err := validationKeeper.GetConsensusAndPluralityVEData(ctx, 100, consensusData)
	require.NoError(t, err)

	t.Logf("Canonical VE obtained successfully")
	t.Logf("Fields with consensus: %d", len(fieldVotePowers))

	// Log some key field values to verify processing
	t.Logf("Canonical EthBlockHeight: %d", canonicalVE.EthBlockHeight)
	t.Logf("Canonical ROCKUSDPrice: %s", canonicalVE.ROCKUSDPrice)
	t.Logf("Canonical EigenDelegationsHash: %x", canonicalVE.EigenDelegationsHash)

	// Step 5: Serialize the canonical vote extension (this is what recordMismatchedVoteExtensions creates)
	canonicalVoteExtBytes, err := json.Marshal(canonicalVE)
	require.NoError(t, err)

	t.Logf("\nCanonical Vote Extension Size: %d bytes", len(canonicalVoteExtBytes))
	t.Logf("Canonical Vote Extension JSON: %s", string(canonicalVoteExtBytes))

	// Step 6: Compare the byte arrays (this is the core comparison in recordMismatchedVoteExtensions)
	t.Log("\n=== STEP 6: Byte Comparison Analysis ===")

	areEqual := bytes.Equal(originalVoteExtBytes, canonicalVoteExtBytes)
	t.Logf("Bytes Equal: %t", areEqual)

	if !areEqual {
		t.Log("*** MISMATCH DETECTED ***")
		t.Log("This explains why single validators are flagged as mismatched!")

		// Find differences
		minLen := len(originalVoteExtBytes)
		if len(canonicalVoteExtBytes) < minLen {
			minLen = len(canonicalVoteExtBytes)
		}

		firstDiff := -1
		for i := 0; i < minLen; i++ {
			if originalVoteExtBytes[i] != canonicalVoteExtBytes[i] {
				firstDiff = i
				break
			}
		}

		if firstDiff >= 0 {
			t.Logf("First difference at byte %d", firstDiff)
			start := max(0, firstDiff-10)
			end := min(len(originalVoteExtBytes), firstDiff+10)
			t.Logf("Original around diff:  %s", string(originalVoteExtBytes[start:end]))
			end = min(len(canonicalVoteExtBytes), firstDiff+10)
			t.Logf("Canonical around diff: %s", string(canonicalVoteExtBytes[start:end]))
		}

		if len(originalVoteExtBytes) != len(canonicalVoteExtBytes) {
			t.Logf("Length difference: original=%d, canonical=%d",
				len(originalVoteExtBytes), len(canonicalVoteExtBytes))
		}
	} else {
		t.Log("Bytes are identical - this is unexpected given the failing test!")
	}

	// Step 7: Test the actual recordMismatchedVoteExtensions logic
	t.Log("\n=== STEP 7: Testing recordMismatchedVoteExtensions Logic ===")

	// Store initial state
	initialInfo, _ := validationKeeper.ValidationInfos.Get(ctx, 100)
	initialMismatchCount := len(initialInfo.MismatchedVoteExtensions)

	// Call the method that's causing the problem
	validationKeeper.RecordMismatchedVoteExtensions(ctx, 100, canonicalVE, consensusData)

	// Check if mismatch was recorded
	afterInfo, err := validationKeeper.ValidationInfos.Get(ctx, 100)
	if err != nil {
		t.Log("No ValidationInfo found - no mismatches recorded")
	} else {
		t.Logf("Mismatches before: %d", initialMismatchCount)
		t.Logf("Mismatches after: %d", len(afterInfo.MismatchedVoteExtensions))

		for i, mismatchedAddr := range afterInfo.MismatchedVoteExtensions {
			t.Logf("Mismatched validator %d: %s", i, mismatchedAddr)
			if mismatchedAddr == validatorHexAddr {
				t.Log("*** BUG CONFIRMED: Single validator flagged as mismatched! ***")
			}
		}
	}

	// Step 8: Check mismatch counts
	mismatchCount, err := validationKeeper.ValidatorMismatchCounts.Get(ctx, validatorHexAddr)
	if err != nil {
		t.Log("No mismatch count recorded")
	} else {
		t.Logf("Mismatch count for validator: %d", mismatchCount.TotalCount)
		if mismatchCount.TotalCount > 0 {
			t.Log("*** BUG CONFIRMED: Mismatch count incremented for single validator! ***")
		}
	}

	// Step 9: Deep dive into the vote processing
	t.Log("\n=== STEP 9: Vote Processing Deep Dive ===")

	// Check if the vote validation is working correctly
	for i, vote := range consensusData.Votes {
		t.Logf("Processing vote %d:", i)
		t.Logf("  Validator: %s", hex.EncodeToString(vote.Validator.Address))
		t.Logf("  Power: %d", vote.Validator.Power)
		t.Logf("  BlockIdFlag: %v", vote.BlockIdFlag)
		t.Logf("  VoteExtension size: %d", len(vote.VoteExtension))

		// Test vote validation
		validatedVE, err := validationKeeper.ValidateVote(ctx, vote, 100)
		if err != nil {
			t.Logf("  Vote validation failed: %v", err)
		} else {
			t.Logf("  Vote validation passed")
			t.Logf("  Validated EthBlockHeight: %d", validatedVE.EthBlockHeight)
		}
	}

	t.Log("\n=== DIAGNOSTIC COMPLETE ===")
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
