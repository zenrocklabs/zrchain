package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/suite"
	ubermock "go.uber.org/mock/gomock"

	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	idTypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	minttypes "github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	treasuryTypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/keeper"
	zentptestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	zentpKeeper      keeper.Keeper
	ctx              sdk.Context
	msgServer        types.MsgServer
	bankKeeper       *zentptestutil.MockBankKeeper
	accountKeeper    *zentptestutil.MockAccountKeeper
	treasuryKeeper   *zentptestutil.MockTreasuryKeeper
	identityKeeper   *zentptestutil.MockIdentityKeeper
	validationKeeper *zentptestutil.MockValidationKeeper
	mintKeeper       *zentptestutil.MockMintKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	storeService := runtime.NewKVStoreService(key)
	memStoreService := runtime.NewMemStoreService(memKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx

	// gomock initializations
	ctrl := ubermock.NewController(s.T())
	accountKeeper := zentptestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := zentptestutil.NewMockBankKeeper(ctrl)
	treasuryKeeper := zentptestutil.NewMockTreasuryKeeper(ctrl)
	identityKeeper := zentptestutil.NewMockIdentityKeeper(ctrl)
	validationKeeper := zentptestutil.NewMockValidationKeeper(ctrl)
	mintKeeper := zentptestutil.NewMockMintKeeper(ctrl)

	// Set up the initial module address mock expectation for NewKeeper
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(moduleAddr).AnyTimes()

	// Assign the mock keepers to the suite fields
	s.accountKeeper = accountKeeper
	s.bankKeeper = bankKeeper
	s.treasuryKeeper = treasuryKeeper
	s.identityKeeper = identityKeeper
	s.validationKeeper = validationKeeper
	s.mintKeeper = mintKeeper
	s.zentpKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		testCtx.Ctx.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		treasuryKeeper,
		bankKeeper,
		accountKeeper,
		identityKeeper,
		validationKeeper,
		mintKeeper,
		memStoreService,
		false,
	)

	s.Require().Equal(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		s.zentpKeeper.Logger())

	err := s.zentpKeeper.ParamStore.Set(s.ctx, types.DefaultParams())
	s.Require().NoError(err)

	// Initialize MintCount to 0
	err = s.zentpKeeper.MintCount.Set(s.ctx, 0)
	s.Require().NoError(err)

	// Initialize BurnCount to 0
	err = s.zentpKeeper.BurnCount.Set(s.ctx, 0)
	s.Require().NoError(err)

	s.msgServer = keeper.NewMsgServerImpl(s.zentpKeeper)
}

func (s *IntegrationTestSuite) TestGetBridgeFeeParams() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.zentpKeeper.ParamStore.Set(s.ctx, params)
	s.Require().NoError(err)

	// Mock mint keeper GetParams
	mintParams := minttypes.DefaultParams()
	s.mintKeeper.EXPECT().GetParams(s.ctx).Return(mintParams, nil)

	// Call the GetBridgeFeeParams function
	protocolWalletAddress, bridgeFee, err := s.zentpKeeper.GetBridgeFeeParams(s.ctx)
	s.Require().NoError(err)

	// Assert the results
	s.Require().Equal(sdk.MustAccAddressFromBech32(mintParams.ProtocolWalletAddress), protocolWalletAddress)
	s.Require().Equal(bridgeFee, params.BridgeFee)
}

func (s *IntegrationTestSuite) TestAddBurn() {

	// Initialize BurnCount to 0
	err := s.zentpKeeper.BurnCount.Set(s.ctx, 0)
	s.Require().NoError(err)

	burns1, err := s.zentpKeeper.BurnCount.Get(s.ctx)
	s.Require().NoError(err)

	burn := types.Bridge{
		Id:               1,
		Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceAddress:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceChain:      "zrchain",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Amount:           1000,
		Denom:            "urock",
		RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
		TxId:             1,
		TxHash:           "123",
		State:            types.BridgeStatus_BRIDGE_STATUS_PENDING,
	}

	err = s.zentpKeeper.AddBurn(s.ctx, &burn)
	s.Require().NoError(err)

	burns2, err := s.zentpKeeper.BurnCount.Get(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(burns2, burns1+1)

	// Get the burn and verify its contents
	burns, err := s.zentpKeeper.GetBurns(s.ctx, "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE", "zrchain", "123")
	s.Require().NoError(err)
	s.Require().Len(burns, 1)
	s.Require().Equal(burn, *burns[0])
}

func (s *IntegrationTestSuite) TestAddMint() {
	// Initialize MintCount to 1
	err := s.zentpKeeper.MintCount.Set(s.ctx, 1)
	s.Require().NoError(err)

	// Verify mint store is empty at index 1
	mints, err := s.zentpKeeper.GetMints(s.ctx, "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty", "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1")
	s.Require().NoError(err)
	s.Require().Empty(mints) // Should be empty initially

	mint := types.Bridge{
		Id:               1,
		Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceAddress:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceChain:      "zen",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Amount:           1000,
		Denom:            "urock",
		RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
		TxId:             1,
		TxHash:           "123",
		State:            types.BridgeStatus_BRIDGE_STATUS_PENDING,
	}

	// Update the mint at index 1
	err = s.zentpKeeper.UpdateMint(s.ctx, 1, &mint)
	s.Require().NoError(err)

	s.validationKeeper.EXPECT().GetLastCompletedZentpMintID(s.ctx).Return(uint64(1), nil).AnyTimes()

	// Verify the mint was stored correctly
	mints, err = s.zentpKeeper.GetMints(s.ctx, "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty", "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1")
	s.Require().NoError(err)
	s.Require().Len(mints, 1)
	s.Require().Equal(mint, *mints[0])
}

func (s *IntegrationTestSuite) TestGetMintsWithStatus() {
	// Initialize MintCount to 1
	err := s.zentpKeeper.MintCount.Set(s.ctx, 1)
	s.Require().NoError(err)

	mint := types.Bridge{
		Id:               1,
		Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceAddress:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceChain:      "zen",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Amount:           1000,
		Denom:            "urock",
		RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
		TxId:             1,
		TxHash:           "123",
		State:            types.BridgeStatus_BRIDGE_STATUS_PENDING,
	}

	s.validationKeeper.EXPECT().GetLastCompletedZentpMintID(s.ctx).Return(uint64(0), nil).AnyTimes()

	err = s.zentpKeeper.UpdateMint(s.ctx, 1, &mint)
	s.Require().NoError(err)

	mints, err := s.zentpKeeper.GetMintsWithStatusPending(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(mints, 1)
	s.Require().Equal(mint, *mints[0])
}

func (s *IntegrationTestSuite) TestUserOwnsKey() {
	user := "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
	key := &treasuryTypes.Key{
		WorkspaceAddr: "workspace123",
	}

	// Test case 1: User owns the key
	s.identityKeeper.EXPECT().
		Workspaces(s.ctx, &idTypes.QueryWorkspacesRequest{
			Creator: user,
			Owner:   user,
		}).
		Return(&idTypes.QueryWorkspacesResponse{
			Workspaces: []*idTypes.Workspace{
				{
					Address: "workspace123",
				},
			},
		}, nil)

	ownsKey := s.zentpKeeper.UserOwnsKey(s.ctx, user, key)
	s.Require().True(ownsKey)

	// Test case 2: User doesn't own the key
	s.identityKeeper.EXPECT().
		Workspaces(s.ctx, &idTypes.QueryWorkspacesRequest{
			Creator: user,
			Owner:   user,
		}).
		Return(&idTypes.QueryWorkspacesResponse{
			Workspaces: []*idTypes.Workspace{
				{
					Address: "different_workspace",
				},
			},
		}, nil)

	ownsKey = s.zentpKeeper.UserOwnsKey(s.ctx, user, key)
	s.Require().False(ownsKey)

	// Test case 3: Error from identity keeper
	s.identityKeeper.EXPECT().
		Workspaces(s.ctx, &idTypes.QueryWorkspacesRequest{
			Creator: user,
			Owner:   user,
		}).
		Return(nil, fmt.Errorf("some error"))

	ownsKey = s.zentpKeeper.UserOwnsKey(s.ctx, user, key)
	s.Require().False(ownsKey)
}

func (s *IntegrationTestSuite) TestGetBridgeFeeAmount() {

	tests := []struct {
		name                   string
		amount                 uint64
		bridgeFee              sdkmath.LegacyDec
		expectedBridgeFeeCoins sdk.Coins
		expectedError          error
	}{
		{
			name:                   "valid test",
			amount:                 100,
			bridgeFee:              sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedBridgeFeeCoins: sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(1))),
			expectedError:          nil,
		},
		{
			name:                   "low amount",
			amount:                 1,
			bridgeFee:              sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedBridgeFeeCoins: sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(0))),
			expectedError:          nil,
		},
		{
			name:                   "low amount",
			amount:                 50,
			bridgeFee:              sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedBridgeFeeCoins: sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(0))),
			expectedError:          nil,
		},
		{
			name:                   "high amount",
			amount:                 1000000000000000000,
			bridgeFee:              sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedBridgeFeeCoins: sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(10000000000000000))),
			expectedError:          nil,
		},
		{
			name:                   "normal amount",
			amount:                 5000000000,
			bridgeFee:              sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedBridgeFeeCoins: sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(50000000))),
			expectedError:          nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			bridgeFeeCoins, err := s.zentpKeeper.GetBridgeFeeAmount(s.ctx, tc.amount, tc.bridgeFee)
			s.Require().NoError(err)

			expectedFee := sdkmath.NewIntFromUint64(tc.amount).ToLegacyDec().Mul(tc.bridgeFee).TruncateInt()
			s.Require().Equal(bridgeFeeCoins, sdk.NewCoins(sdk.NewCoin(params.BondDenom, expectedFee)))
		})
	}
}

func (s *IntegrationTestSuite) TestAddFeeToBridgeAmount() {

	tests := []struct {
		name                string
		amount              uint64
		bridgeFee           sdkmath.LegacyDec
		expectedTotalAmount uint64
		expectedError       error
	}{
		{
			name:                "valid test",
			amount:              100,
			bridgeFee:           sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedTotalAmount: 101,
			expectedError:       nil,
		},
		{
			name:                "low amount",
			amount:              1,
			bridgeFee:           sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedTotalAmount: 1,
			expectedError:       nil,
		},
		{
			name:                "low amount",
			amount:              50,
			bridgeFee:           sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedTotalAmount: 50,
			expectedError:       nil,
		},
		{
			name:                "high amount",
			amount:              1000000000000000000,
			bridgeFee:           sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedTotalAmount: 1010000000000000000,
			expectedError:       nil,
		},
		{
			name:                "normal amount",
			amount:              5000000000,
			bridgeFee:           sdkmath.LegacyNewDecWithPrec(1, 2),
			expectedTotalAmount: 5050000000,
			expectedError:       nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			// Setup test parameters
			params := types.DefaultParams()
			params.BridgeFee = tc.bridgeFee
			err := s.zentpKeeper.ParamStore.Set(s.ctx, params)
			s.Require().NoError(err)

			// Mock mint keeper GetParams
			mintParams := minttypes.DefaultParams()
			s.mintKeeper.EXPECT().GetParams(s.ctx).Return(mintParams, nil)

			totalAmount, err := s.zentpKeeper.AddFeeToBridgeAmount(s.ctx, tc.amount)
			s.Require().NoError(err)
			s.Require().Equal(totalAmount, tc.expectedTotalAmount)
		})
	}
}

func (s *IntegrationTestSuite) TestCalculateZentpMintFee() {

	tests := []struct {
		name                string
		amount              uint64
		expectedTotalAmount sdkmath.Int
		expectedTotalFee    sdkmath.Int
		expectedError       error
	}{
		{
			name:                "1 ROCK + flat fee",
			amount:              1000000,
			expectedTotalAmount: sdkmath.NewInt(201005000),
			expectedTotalFee:    sdkmath.NewInt(200005000),
			expectedError:       nil,
		},
		{
			name:                "0.0001 ROCK + flat fee",
			amount:              100,
			expectedTotalAmount: sdkmath.NewInt(200000100),
			expectedTotalFee:    sdkmath.NewInt(200000000),
			expectedError:       nil,
		},
		{
			name:                "large amount",
			amount:              1000000000000000000,
			expectedTotalAmount: sdkmath.NewInt(1005000000200000000),
			expectedTotalFee:    sdkmath.NewInt(5000000200000000),
			expectedError:       nil,
		},
		{
			name:                "int overflow - very large amount",
			amount:              18400000000000000000,
			expectedTotalAmount: sdkmath.Int{},
			expectedTotalFee:    sdkmath.Int{},
			expectedError:       fmt.Errorf("total amount %s exceeds max uint64", "18492000000200000000"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()

			totalAmount, totalFee, err := s.zentpKeeper.CalculateZentpMintFee(s.ctx, tc.amount)

			if tc.expectedError != nil {
				s.Require().Error(err)
				s.Require().Equal(tc.expectedError.Error(), err.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedTotalAmount, totalAmount)
				s.Require().Equal(tc.expectedTotalFee, totalFee)
			}
		})
	}
}
