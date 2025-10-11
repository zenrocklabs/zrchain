package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	treasuryKeeper keeper.Keeper
	ctx            sdk.Context
	msgServer      types.MsgServer
	bankKeeper     *treasurytestutil.MockBankKeeper
	identityKeeper *treasurytestutil.MockIdentityKeeper
	policyKeeper   *treasurytestutil.MockPolicyKeeper
	zentpKeeper    *treasurytestutil.MockZentpKeeper
	ctrl           *gomock.Controller
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {

	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx

	ctrl := gomock.NewController(s.T())
	bankKeeper := treasurytestutil.NewMockBankKeeper(ctrl)
	identityKeeper := treasurytestutil.NewMockIdentityKeeper(ctrl)
	policyKeeper := treasurytestutil.NewMockPolicyKeeper(ctrl)
	zentpKeeper := treasurytestutil.NewMockZentpKeeper(ctrl)

	s.bankKeeper = bankKeeper
	s.identityKeeper = identityKeeper
	s.policyKeeper = policyKeeper
	s.zentpKeeper = zentpKeeper
	s.ctrl = ctrl

	// Set up mock expectations before creating the keeper
	s.policyKeeper.EXPECT().ActionHandler(gomock.Any()).Return(nil, false).AnyTimes()
	s.policyKeeper.EXPECT().RegisterActionHandler(gomock.Any(), gomock.Any()).AnyTimes()
	s.policyKeeper.EXPECT().GeneratorHandler(gomock.Any()).Return(nil, false).AnyTimes()
	s.policyKeeper.EXPECT().RegisterPolicyGeneratorHandler(gomock.Any(), gomock.Any()).AnyTimes()

	s.treasuryKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		testCtx.Ctx.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		bankKeeper,
		identityKeeper,
		policyKeeper,
		nil, // zenBTCKeeper
		zentpKeeper,
		nil, // dctKeeper
		nil, // zenexKeeper
	)

	s.Require().Equal(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		s.treasuryKeeper.Logger())

	err := s.treasuryKeeper.ParamStore.Set(s.ctx, types.DefaultParams())
	s.Require().NoError(err)

	s.msgServer = keeper.NewMsgServerImpl(s.treasuryKeeper, false)
}

func (s *IntegrationTestSuite) Test_TreasuryKeeper_splitKeyringFee() {
	type args struct {
		feeAddr    string
		fee        uint64
		commission uint64
	}

	// Generate addresses once and reuse them
	addrFrom := sample.AccAddress()
	addrTo := sample.AccAddress()

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Send fee to address",
			args: args{
				feeAddr:    addrTo,
				fee:        1000,
				commission: 10,
			},
		},
		{
			name: "Send fee to module",
			args: args{
				feeAddr:    types.KeyringCollectorName,
				fee:        1000,
				commission: 10,
			},
		},
		{
			name: "Zero commission",
			args: args{
				feeAddr:    addrTo,
				fee:        1000,
				commission: 0,
			},
		},
		{
			name: "100% commission",
			args: args{
				feeAddr:    addrTo,
				fee:        1000,
				commission: 100,
			},
		},
		{
			name: "Small fee amount",
			args: args{
				feeAddr:    addrTo,
				fee:        10,
				commission: 10,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Calculate expected commission and fee
			commission := tt.args.commission
			fee := tt.args.fee
			feeAddr := tt.args.feeAddr

			// Update the params with the test case's commission value
			params, err := s.treasuryKeeper.ParamStore.Get(s.ctx)
			s.Require().NoError(err)
			params.KeyringCommission = commission
			params.KeyringCommissionDestination = feeAddr
			err = s.treasuryKeeper.ParamStore.Set(s.ctx, params)
			s.Require().NoError(err)

			// Calculate expected commission and fee split
			commissionAmount := uint64((fee * commission) / 100)

			if commissionAmount > 0 {
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(
					s.ctx,
					sdk.MustAccAddressFromBech32(addrFrom),
					types.KeyringCollectorName,
					gomock.Any(), // coin
				).Return(nil)
			}

			if feeAddr == types.KeyringCollectorName {
				s.bankKeeper.EXPECT().SendCoinsFromAccountToModule(
					s.ctx,
					sdk.MustAccAddressFromBech32(addrFrom),
					types.KeyringCollectorName,
					gomock.Any(), // coin
				).Return(nil)
			} else {
				s.bankKeeper.EXPECT().SendCoins(
					s.ctx,
					sdk.MustAccAddressFromBech32(addrFrom),
					sdk.MustAccAddressFromBech32(feeAddr),
					gomock.Any(), // coin
				).Return(nil)
			}

			err = s.treasuryKeeper.SplitKeyringFee(s.ctx, addrFrom, feeAddr, fee)
			s.Require().NoError(err)
		})
	}
}
