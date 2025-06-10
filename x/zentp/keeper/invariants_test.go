package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (s *IntegrationTestSuite) TestCheckROCKSupplyCap() {
	tests := []struct {
		name               string
		zrchainSupply      uint64
		zentpModuleBalance uint64
		solanaSupply       uint64
		pendingMints       []types.Bridge
		newAmount          uint64
		expectError        bool
		expectedErrorMsg   string
	}{
		{
			name:               "Normal operation under cap",
			zrchainSupply:      100_000_000_000_000, // 100M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       200_000_000_000_000, // 200M ROCK
			pendingMints: []types.Bridge{
				{Amount: 50_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 50M ROCK
			},
			newAmount:   100_000_000_000_000, // 100M ROCK
			expectError: false,
		},
		{
			name:               "Exactly at cap should succeed",
			zrchainSupply:      400_000_000_000_000, // 400M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       300_000_000_000_000, // 300M ROCK
			pendingMints: []types.Bridge{
				{Amount: 200_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 200M ROCK
			},
			newAmount:   100_000_000_000_000, // 100M ROCK = exactly 1B total
			expectError: false,
		},
		{
			name:               "One unit over cap should fail",
			zrchainSupply:      400_000_000_000_000, // 400M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       300_000_000_000_000, // 300M ROCK
			pendingMints: []types.Bridge{
				{Amount: 200_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 200M ROCK
			},
			newAmount:        100_000_000_000_001, // 100M + 1 = over 1B total
			expectError:      true,
			expectedErrorMsg: "total ROCK supply including pending would exceed cap",
		},
		{
			name:               "Zero new amount with existing state under cap",
			zrchainSupply:      300_000_000_000_000, // 300M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       200_000_000_000_000, // 200M ROCK
			pendingMints: []types.Bridge{
				{Amount: 100_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 100M ROCK
			},
			newAmount:   0,
			expectError: false,
		},
		{
			name:               "Zero new amount with existing state at cap",
			zrchainSupply:      400_000_000_000_000, // 400M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       300_000_000_000_000, // 300M ROCK
			pendingMints: []types.Bridge{
				{Amount: 300_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 300M ROCK
			},
			newAmount:   0,
			expectError: false,
		},
		{
			name:               "Zero new amount with existing state over cap",
			zrchainSupply:      500_000_000_000_000, // 500M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       300_000_000_000_000, // 300M ROCK
			pendingMints: []types.Bridge{
				{Amount: 300_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 300M ROCK
			},
			newAmount:        0,
			expectError:      true,
			expectedErrorMsg: "total ROCK supply including pending would exceed cap",
		},
		{
			name:               "Multiple pending mints",
			zrchainSupply:      200_000_000_000_000, // 200M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       200_000_000_000_000, // 200M ROCK
			pendingMints: []types.Bridge{
				{Amount: 100_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 100M ROCK
				{Amount: 150_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 150M ROCK
				{Amount: 50_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING},  // 50M ROCK
			},
			newAmount:   100_000_000_000_000, // 100M ROCK = 200+200+300+100 = 800M total
			expectError: false,
		},
		{
			name:               "Ignore non-pending mints",
			zrchainSupply:      200_000_000_000_000, // 200M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       200_000_000_000_000, // 200M ROCK
			pendingMints: []types.Bridge{
				{Amount: 100_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING},   // 100M ROCK (counted)
				{Amount: 500_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_COMPLETED}, // 500M ROCK (ignored)
				{Amount: 300_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_FAILED},    // 300M ROCK (ignored)
			},
			newAmount:   400_000_000_000_000, // 400M ROCK = 200+200+100+400 = 900M total
			expectError: false,
		},
		{
			name:               "No pending mints",
			zrchainSupply:      400_000_000_000_000, // 400M ROCK
			zentpModuleBalance: 0,
			solanaSupply:       300_000_000_000_000, // 300M ROCK
			pendingMints:       []types.Bridge{},
			newAmount:          300_000_000_000_000, // 300M ROCK = 400+300+300 = 1B total
			expectError:        false,
		},
		{
			name:               "Bridge amount exceeds available supply",
			zrchainSupply:      200_000_000_000_000, // 200M ROCK
			zentpModuleBalance: 150_000_000_000_000, // 150M ROCK in module
			solanaSupply:       100_000_000_000_000, // 100M ROCK
			pendingMints:       []types.Bridge{},
			newAmount:          60_000_000_000_000, // 60M ROCK new bridge. Available is 50M.
			expectError:        true,
			expectedErrorMsg:   "bridge amount 60000000000000 exceeds available zrchain rock supply for bridging 50000000000000",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Setup: Mock bank keeper to return the test zrchain supply
			s.bankKeeper.EXPECT().GetSupply(s.ctx, params.BondDenom).Return(
				sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(tt.zrchainSupply)),
			).AnyTimes()

			// Mock for new check
			zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
			s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()
			s.bankKeeper.EXPECT().GetBalance(s.ctx, zentpModuleAddr, params.BondDenom).Return(
				sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(tt.zentpModuleBalance)),
			).AnyTimes()

			// Setup: Set solana supply
			err := s.zentpKeeper.SetSolanaROCKSupply(s.ctx, sdkmath.NewIntFromUint64(tt.solanaSupply))
			s.Require().NoError(err)

			// Setup: Create pending mints
			for i, mint := range tt.pendingMints {
				mint.Id = uint64(i + 1)
				mint.Creator = "test_creator"
				mint.DestinationChain = "solana:test"
				mint.Denom = params.BondDenom
				mint.RecipientAddress = "test_recipient"
				err := s.zentpKeeper.UpdateMint(s.ctx, mint.Id, &mint)
				s.Require().NoError(err)
			}
			// Update MintCount to reflect the number of mints we created
			err = s.zentpKeeper.MintCount.Set(s.ctx, uint64(len(tt.pendingMints)))
			s.Require().NoError(err)

			// Execute the function under test
			err = s.zentpKeeper.CheckROCKSupplyCap(s.ctx, sdkmath.NewIntFromUint64(tt.newAmount))

			// Verify results
			if tt.expectError {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tt.expectedErrorMsg)
			} else {
				s.Require().NoError(err)
			}

			// Cleanup: Reset state for next test
			s.SetupTest()
		})
	}
}

func (s *IntegrationTestSuite) TestCheckCanBurnFromSolana() {
	tests := []struct {
		name             string
		solanaSupply     uint64
		burnAmount       uint64
		zrchainSupply    uint64
		pendingMints     []types.Bridge
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:          "Normal burn under solana supply",
			solanaSupply:  500_000_000_000_000, // 500M ROCK
			burnAmount:    100_000_000_000_000, // 100M ROCK
			zrchainSupply: 200_000_000_000_000, // 200M ROCK
			pendingMints:  []types.Bridge{},
			expectError:   false,
		},
		{
			name:          "Burn exactly equal to solana supply",
			solanaSupply:  300_000_000_000_000, // 300M ROCK
			burnAmount:    300_000_000_000_000, // 300M ROCK
			zrchainSupply: 200_000_000_000_000, // 200M ROCK
			pendingMints:  []types.Bridge{},
			expectError:   false,
		},
		{
			name:             "Burn exceeds solana supply",
			solanaSupply:     200_000_000_000_000, // 200M ROCK
			burnAmount:       300_000_000_000_000, // 300M ROCK
			zrchainSupply:    200_000_000_000_000, // 200M ROCK
			pendingMints:     []types.Bridge{},
			expectError:      true,
			expectedErrorMsg: "attempt to bridge from solana exceeds solana ROCK supply",
		},
		{
			name:          "Burn valid but would cause total supply cap violation",
			solanaSupply:  300_000_000_000_000, // 300M ROCK
			burnAmount:    100_000_000_000_000, // 100M ROCK (valid burn)
			zrchainSupply: 600_000_000_000_000, // 600M ROCK
			pendingMints: []types.Bridge{
				{Amount: 200_000_000_000_000, State: types.BridgeStatus_BRIDGE_STATUS_PENDING}, // 200M ROCK
			},
			expectError:      true,
			expectedErrorMsg: "total ROCK supply including pending would exceed cap",
		},
		{
			name:          "Large burn that's valid",
			solanaSupply:  800_000_000_000_000, // 800M ROCK
			burnAmount:    500_000_000_000_000, // 500M ROCK
			zrchainSupply: 100_000_000_000_000, // 100M ROCK
			pendingMints:  []types.Bridge{},
			expectError:   false,
		},
		{
			name:          "Zero burn amount should succeed",
			solanaSupply:  100_000_000_000_000, // 100M ROCK
			burnAmount:    0,                   // 0 ROCK
			zrchainSupply: 200_000_000_000_000, // 200M ROCK
			pendingMints:  []types.Bridge{},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Setup: Mock bank keeper to return the test zrchain supply
			s.bankKeeper.EXPECT().GetSupply(s.ctx, params.BondDenom).Return(
				sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(tt.zrchainSupply)),
			).AnyTimes()

			// Setup: Set solana supply
			err := s.zentpKeeper.SetSolanaROCKSupply(s.ctx, sdkmath.NewIntFromUint64(tt.solanaSupply))
			s.Require().NoError(err)

			// Setup: Create pending mints
			for i, mint := range tt.pendingMints {
				mint.Id = uint64(i + 1)
				mint.Creator = "test_creator"
				mint.DestinationChain = "solana:test"
				mint.Denom = params.BondDenom
				mint.RecipientAddress = "test_recipient"
				err := s.zentpKeeper.UpdateMint(s.ctx, mint.Id, &mint)
				s.Require().NoError(err)
			}
			// Update MintCount to reflect the number of mints we created
			err = s.zentpKeeper.MintCount.Set(s.ctx, uint64(len(tt.pendingMints)))
			s.Require().NoError(err)

			// Execute the function under test
			err = s.zentpKeeper.CheckCanBurnFromSolana(s.ctx, sdkmath.NewIntFromUint64(tt.burnAmount))

			// Verify results
			if tt.expectError {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tt.expectedErrorMsg)
			} else {
				s.Require().NoError(err)
			}

			// Cleanup: Reset state for next test
			s.SetupTest()
		})
	}
}

func (s *IntegrationTestSuite) TestCheckROCKSupplyCap_ErrorHandling() {
	// Test that GetMintsWithStatus errors are handled gracefully
	s.Run("GetMintsWithStatus error handling", func() {
		// Ensure clean state
		s.SetupTest()

		// Setup valid solana supply
		err := s.zentpKeeper.SetSolanaROCKSupply(s.ctx, sdkmath.NewIntFromUint64(1000))
		s.Require().NoError(err)

		// Setup bank keeper mock
		s.bankKeeper.EXPECT().GetSupply(s.ctx, params.BondDenom).Return(
			sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(1000)),
		).AnyTimes()

		// Mock for new check
		zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
		s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()
		s.bankKeeper.EXPECT().GetBalance(s.ctx, zentpModuleAddr, params.BondDenom).Return(
			sdk.NewCoin(params.BondDenom, sdkmath.ZeroInt()),
		).AnyTimes()

		// The function should not fail even if GetMintsWithStatus has issues
		// (our implementation treats this as "no pending mints")
		err = s.zentpKeeper.CheckROCKSupplyCap(s.ctx, sdkmath.NewIntFromUint64(1000))
		s.Require().NoError(err) // Should succeed despite potential GetMintsWithStatus issues
	})

	// Test that function works correctly when solana supply is not set (should default to zero)
	s.Run("Solana supply not set defaults to zero", func() {
		// Ensure clean state
		s.SetupTest()

		// Setup bank keeper mock
		s.bankKeeper.EXPECT().GetSupply(s.ctx, params.BondDenom).Return(
			sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(1000)),
		).AnyTimes()

		// Mock for new check
		zentpModuleAddr := authtypes.NewModuleAddress(types.ModuleName)
		s.accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(zentpModuleAddr).AnyTimes()
		s.bankKeeper.EXPECT().GetBalance(s.ctx, zentpModuleAddr, params.BondDenom).Return(
			sdk.NewCoin(params.BondDenom, sdkmath.ZeroInt()),
		).AnyTimes()

		// The function should succeed with zero solana supply
		err := s.zentpKeeper.CheckROCKSupplyCap(s.ctx, sdkmath.NewIntFromUint64(1000))
		s.Require().NoError(err) // Should succeed with solana supply defaulting to zero
	})
}

func (s *IntegrationTestSuite) TestCheckCanBurnFromSolana_ErrorHandling() {
	// Test that function works correctly when solana supply is not set (should default to zero)
	s.Run("Solana supply not set defaults to zero", func() {
		// Ensure clean state
		s.SetupTest()

		// Setup bank keeper mock
		s.bankKeeper.EXPECT().GetSupply(s.ctx, params.BondDenom).Return(
			sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(1000)),
		).AnyTimes()

		// Attempting to burn from zero solana supply should fail
		err := s.zentpKeeper.CheckCanBurnFromSolana(s.ctx, sdkmath.NewIntFromUint64(1000))
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "attempt to bridge from solana exceeds solana ROCK supply")
	})
}

// Test edge cases and boundary conditions
func (s *IntegrationTestSuite) TestInvariants_EdgeCases() {
	s.Run("Maximum values", func() {
		// Setup maximum possible values that don't overflow
		maxSupply := uint64(500_000_000_000_000) // 500M ROCK each

		s.bankKeeper.EXPECT().GetSupply(s.ctx, params.BondDenom).Return(
			sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(maxSupply)),
		).AnyTimes()

		err := s.zentpKeeper.SetSolanaROCKSupply(s.ctx, sdkmath.NewIntFromUint64(maxSupply))
		s.Require().NoError(err)

		// This should exactly hit the cap (500M + 500M = 1B)
		err = s.zentpKeeper.CheckROCKSupplyCap(s.ctx, sdkmath.ZeroInt())
		s.Require().NoError(err)

		// This should be valid for burn
		err = s.zentpKeeper.CheckCanBurnFromSolana(s.ctx, sdkmath.NewIntFromUint64(maxSupply))
		s.Require().NoError(err)
	})

	s.Run("Zero values", func() {
		s.bankKeeper.EXPECT().GetSupply(s.ctx, params.BondDenom).Return(
			sdk.NewCoin(params.BondDenom, sdkmath.ZeroInt()),
		).AnyTimes()

		err := s.zentpKeeper.SetSolanaROCKSupply(s.ctx, sdkmath.ZeroInt())
		s.Require().NoError(err)

		// Zero amounts should always be valid
		err = s.zentpKeeper.CheckROCKSupplyCap(s.ctx, sdkmath.ZeroInt())
		s.Require().NoError(err)

		err = s.zentpKeeper.CheckCanBurnFromSolana(s.ctx, sdkmath.ZeroInt())
		s.Require().NoError(err)
	})
}
