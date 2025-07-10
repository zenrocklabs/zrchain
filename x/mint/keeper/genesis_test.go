package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	ubermock "go.uber.org/mock/gomock"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/keeper"
	minttestutil "github.com/Zenrock-Foundation/zrchain/v6/x/mint/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var minterAcc = authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)

type GenesisTestSuite struct {
	suite.Suite

	sdkCtx        sdk.Context
	keeper        keeper.Keeper
	cdc           codec.BinaryCodec
	accountKeeper types.AccountKeeper
	key           *storetypes.KVStoreKey
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (s *GenesisTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})

	// gomock initializations
	ctrl := ubermock.NewController(s.T())
	s.cdc = codec.NewProtoCodec(encCfg.InterfaceRegistry)
	s.sdkCtx = testCtx.Ctx
	s.key = key

	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	zentpKeeper := minttestutil.NewMockZentpKeeper(ctrl)
	s.accountKeeper = accountKeeper
	accountKeeper.EXPECT().GetModuleAddress(minterAcc.Name).Return(minterAcc.GetAddress())
	accountKeeper.EXPECT().GetModuleAccount(s.sdkCtx, minterAcc.Name).Return(minterAcc)

	s.keeper = keeper.NewKeeper(s.cdc, runtime.NewKVStoreService(key), stakingKeeper, accountKeeper, bankKeeper, zentpKeeper, authtypes.FeeCollectorName, "")
}

func (s *GenesisTestSuite) TestImportExportGenesis() {
	genesisState := types.DefaultGenesisState()
	genesisState.Minter = types.NewMinter(math.LegacyNewDecWithPrec(20, 2), math.LegacyNewDec(1))
	genesisState.Params = types.NewParams(
		"testDenom",
		"zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn",
		math.LegacyNewDecWithPrec(15, 2),
		math.LegacyNewDecWithPrec(22, 2),
		math.LegacyNewDecWithPrec(9, 2),
		math.LegacyNewDecWithPrec(69, 2),
		math.LegacyNewDecWithPrec(10, 2),
		math.LegacyNewDecWithPrec(10, 2),
		math.LegacyNewDecWithPrec(10, 2),
		math.LegacyNewDecWithPrec(10, 2),
		math.LegacyNewDecWithPrec(10, 2),
		math.LegacyNewDecWithPrec(10, 2),
		uint64(60*60*8766/5),
	)

	s.keeper.InitGenesis(s.sdkCtx, s.accountKeeper, genesisState)

	minter, err := s.keeper.Minter.Get(s.sdkCtx)
	s.Require().Equal(genesisState.Minter, minter)
	s.Require().NoError(err)

	invalidCtx := testutil.DefaultContextWithDB(s.T(), s.key, storetypes.NewTransientStoreKey("transient_test"))
	_, err = s.keeper.Minter.Get(invalidCtx.Ctx)
	s.Require().ErrorIs(err, collections.ErrNotFound)

	params, err := s.keeper.Params.Get(s.sdkCtx)
	s.Require().Equal(genesisState.Params, params)
	s.Require().NoError(err)

	genesisState2 := s.keeper.ExportGenesis(s.sdkCtx)
	s.Require().Equal(genesisState, genesisState2)
}
