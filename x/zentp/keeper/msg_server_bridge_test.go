package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	minttypes "github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
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

	s.msgServer = keeper.NewMsgServerImpl(s.zentpKeeper)
}

func (s *IntegrationTestSuite) TestBridge() {
	// Setup test parameters
	params := types.DefaultParams()
	err := s.zentpKeeper.ParamStore.Set(s.ctx, params)
	s.Require().NoError(err)

	// Mock getting the mint params
	s.mintKeeper.EXPECT().GetParams(s.ctx).Return(minttypes.DefaultParams(), nil)

	// Create test message
	msg := &types.MsgBridge{
		Creator:          "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		SourceAddress:    "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		DestinationChain: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
		Amount:           1000,
		Denom:            "urock",
		RecipientAddress: "1BbzosnmC3EVe7XcMgHYd6fUtcfdzUvfeaVZxaZ2QsE",
	}

	// Mock bank keeper GetBalance
	s.bankKeeper.EXPECT().GetBalance(
		s.ctx,
		sdk.MustAccAddressFromBech32(msg.Creator),
		msg.Denom,
	).Return(sdk.NewCoin(msg.Denom, math.NewIntFromUint64(msg.Amount*20)))

	burnAmount := math.LegacyNewDecFromInt(math.NewIntFromUint64(msg.Amount)).Mul(params.BridgeFee).TruncateInt()

	// Mock bank keeper SendCoinsFromAccountToModule
	s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(
		s.ctx,
		sdk.MustAccAddressFromBech32(msg.SourceAddress),
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin("urock", math.NewIntFromUint64(msg.Amount).Add(burnAmount))),
	).Return(nil)

	// Mock validation keeper SetSolanaRequestedNonce
	s.validationKeeper.EXPECT().SetSolanaRequestedNonce(
		s.ctx,
		params.Solana.NonceAccountKey, // Default nonce account key
		true,
	).Return(nil)

	// Mock validation keeper SetSolanaRequestedAccount
	s.validationKeeper.EXPECT().SetSolanaRequestedAccount(
		s.ctx,
		msg.RecipientAddress,
		true,
	).Return(nil)

	// Call the Bridge handler
	response, err := s.msgServer.Bridge(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(uint64(1), response.Id)
}
