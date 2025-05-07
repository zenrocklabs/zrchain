package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	idTypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	minttypes "github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	treasuryTypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/keeper"
	zentptestutil "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
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
	ctrl := gomock.NewController(s.T())
	accountKeeper := zentptestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := zentptestutil.NewMockBankKeeper(ctrl)
	treasuryKeeper := zentptestutil.NewMockTreasuryKeeper(ctrl)
	identityKeeper := zentptestutil.NewMockIdentityKeeper(ctrl)
	validationKeeper := zentptestutil.NewMockValidationKeeper(ctrl)
	mintKeeper := zentptestutil.NewMockMintKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{})

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
		"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
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

func (s *IntegrationTestSuite) TestGetBridgeFeeAmount() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.zentpKeeper.ParamStore.Set(s.ctx, params)
	s.Require().NoError(err)

	// Call the GetBridgeFeeAmount function
	bridgeFeeAmount, err := s.zentpKeeper.GetBridgeFeeAmount(s.ctx, 1000, params.BridgeFee)
	s.Require().NoError(err)

	// Assert the results
	s.Require().Equal(sdk.NewCoins(sdk.NewCoin("urock", sdkmath.NewIntFromUint64(10))), bridgeFeeAmount)
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

	err = s.zentpKeeper.UpdateMint(s.ctx, 1, &mint)
	s.Require().NoError(err)

	mints, err := s.zentpKeeper.GetMintsWithStatus(s.ctx, types.BridgeStatus_BRIDGE_STATUS_PENDING)
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
