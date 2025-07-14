package keeper_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	ubermock "go.uber.org/mock/gomock"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	validationkeeper "github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	validationtestutil "github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	cmtcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

var (
	bondedAcc    = authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName)
	notBondedAcc = authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName)
	PKs          = simtestutil.CreateTestPubKeys(500)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	stakingKeeper *stakingkeeper.Keeper
	bankKeeper    *stakingtestutil.MockBankKeeper
	accountKeeper *stakingtestutil.MockAccountKeeper
	queryClient   stakingtypes.QueryClient
	msgServer     stakingtypes.MsgServer
}

type ValidationKeeperTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	validationKeeper *validationkeeper.Keeper
	bankKeeper       *validationtestutil.MockBankKeeper
	accountKeeper    *validationtestutil.MockAccountKeeper
	queryClient      validationtypes.QueryClient
	msgServer        validationtypes.MsgServer
	zenBTCCtrl       *ubermock.Controller
}

func (s *KeeperTestSuite) SetupTest() {
	require := s.Require()
	key := storetypes.NewKVStoreKey(stakingtypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	ctrl := gomock.NewController(s.T())
	accountKeeper := stakingtestutil.NewMockAccountKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(stakingtypes.BondedPoolName).Return(bondedAcc.GetAddress())
	accountKeeper.EXPECT().GetModuleAddress(stakingtypes.NotBondedPoolName).Return(notBondedAcc.GetAddress())
	accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("zen")).AnyTimes()

	bankKeeper := stakingtestutil.NewMockBankKeeper(ctrl)

	keeper := stakingkeeper.NewKeeper(
		encCfg.Codec,
		storeService,
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		address.NewBech32Codec("zenvaloper"),
		address.NewBech32Codec("zenvalcons"),
	)
	require.NoError(keeper.SetParams(ctx, stakingtypes.DefaultParams()))

	s.ctx = ctx
	s.stakingKeeper = keeper
	s.bankKeeper = bankKeeper
	s.accountKeeper = accountKeeper

	stakingtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	stakingtypes.RegisterQueryServer(queryHelper, stakingkeeper.Querier{Keeper: keeper})
	s.queryClient = stakingtypes.NewQueryClient(queryHelper)
	s.msgServer = stakingkeeper.NewMsgServerImpl(keeper)
}

func (s *ValidationKeeperTestSuite) ValidationKeeperSetupTest() (*validationkeeper.Keeper, *ubermock.Controller) {
	require := s.Require()
	key := storetypes.NewKVStoreKey(validationtypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	ctrl := ubermock.NewController(s.T())
	accountKeeper := validationtestutil.NewMockAccountKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(validationtypes.BondedPoolName).Return(bondedAcc.GetAddress())
	accountKeeper.EXPECT().GetModuleAddress(validationtypes.NotBondedPoolName).Return(notBondedAcc.GetAddress())
	accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("zen")).AnyTimes()
	accountKeeper.EXPECT().GetModuleAccount(ctx, "bonded_tokens_pool").Return(bondedAcc).AnyTimes()
	accountKeeper.EXPECT().GetModuleAccount(ctx, "not_bonded_tokens_pool").Return(notBondedAcc).AnyTimes()
	accountKeeper.EXPECT().SetModuleAccount(ctx, bondedAcc).AnyTimes()
	accountKeeper.EXPECT().SetModuleAccount(ctx, notBondedAcc).AnyTimes()

	bankKeeper := validationtestutil.NewMockBankKeeper(ctrl)
	// Return the correct bonded pool balance that matches the validator tokens in TestGenesis
	bondDenom := validationtypes.DefaultParams().BondDenom
	bankKeeper.EXPECT().GetAllBalances(ctx, bondedAcc.GetAddress()).Return(
		sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(125000000000000))),
	).AnyTimes()
	bankKeeper.EXPECT().GetAllBalances(ctx, notBondedAcc.GetAddress()).Return(sdk.NewCoins()).AnyTimes()

	zentpKeeper := validationtestutil.NewMockZentpKeeper(ctrl)
	zentpKeeper.EXPECT().GetSolanaParams(ubermock.Any()).Return(&zentptypes.Solana{
		SignerKeyId:       10,
		ProgramId:         "DXREJumiQhNejXa1b5EFPUxtSYdyJXBdiHeu6uX1ribA",
		NonceAuthorityKey: 11,
		NonceAccountKey:   12,
		MintAddress:       "StVNdHNSFK3uVTL5apWHysgze4M8zrsqwjEAH1JM87i",
		FeeWallet:         "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
		Fee:               0,
		Btl:               20,
	}).AnyTimes()
	zentpKeeper.EXPECT().GetMintsWithStatus(ubermock.Any(), ubermock.Any()).Return([]*zentptypes.Bridge{}, nil).AnyTimes()
	zentpKeeper.EXPECT().GetMintsWithStatusPending(ubermock.Any()).Return([]*zentptypes.Bridge{}, nil).AnyTimes()
	zentpKeeper.EXPECT().GetTotalROCKSupply(ubermock.Any()).Return(math.NewInt(1000000), nil).AnyTimes()
	zentpKeeper.EXPECT().CheckROCKSupplyCap(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zentpKeeper.EXPECT().UpdateMint(ubermock.Any(), ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zentpKeeper.EXPECT().GetSolanaROCKSupply(ubermock.Any()).Return(math.NewInt(500000), nil).AnyTimes()
	zentpKeeper.EXPECT().SetSolanaROCKSupply(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zentpKeeper.EXPECT().GetBurns(ubermock.Any(), ubermock.Any(), ubermock.Any(), ubermock.Any()).Return([]*zentptypes.Bridge{}, nil).AnyTimes()
	zentpKeeper.EXPECT().CheckCanBurnFromSolana(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zentpKeeper.EXPECT().GetBridgeFeeParams(ubermock.Any()).Return(sdk.AccAddress{}, math.LegacyNewDec(0), nil).AnyTimes()
	zentpKeeper.EXPECT().GetBridgeFeeAmount(ubermock.Any(), ubermock.Any(), ubermock.Any()).Return(sdk.NewCoins(), nil).AnyTimes()
	zentpKeeper.EXPECT().AddBurn(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	treasuryKeeper := validationtestutil.NewMockTreasuryKeeper(ctrl)
	slashingKeeper := validationtestutil.NewMockSlashingKeeper(ctrl)

	newctrl := ubermock.NewController(s.T())
	zenBTCKeeper := validationtestutil.NewMockZenBTCKeeper(newctrl)
	// Set up expectations for key ID methods that are called by getZenBTCKeyIDs
	zenBTCKeeper.EXPECT().GetStakerKeyID(ubermock.Any()).Return(uint64(1)).AnyTimes()
	zenBTCKeeper.EXPECT().GetEthMinterKeyID(ubermock.Any()).Return(uint64(2)).AnyTimes()
	zenBTCKeeper.EXPECT().GetUnstakerKeyID(ubermock.Any()).Return(uint64(3)).AnyTimes()
	zenBTCKeeper.EXPECT().GetCompleterKeyID(ubermock.Any()).Return(uint64(4)).AnyTimes()
	// Set up expectation for GetSolanaParams which is called by retrieveSolanaNonces
	zenBTCKeeper.EXPECT().GetSolanaParams(ubermock.Any()).Return(validationtestutil.DefaultSolana).AnyTimes()
	zenBTCKeeper.EXPECT().GetFirstRedemptionAwaitingSign(ubermock.Any()).Return(uint64(1), nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetFirstRedemptionAwaitingSign(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetSupply(ubermock.Any()).Return(zenbtctypes.Supply{CustodiedBTC: 1, MintedZenBTC: 1, PendingZenBTC: 1}, nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetSupply(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetExchangeRate(ubermock.Any()).Return(math.LegacyNewDec(1), nil).AnyTimes()
	zenBTCKeeper.EXPECT().WalkRedemptions(ubermock.Any(), ubermock.Any()).DoAndReturn(func(ctx sdk.Context, fn func(id uint64, r zenbtctypes.Redemption) (bool, error)) error { return nil }).AnyTimes()
	zenBTCKeeper.EXPECT().WalkBurnEvents(ubermock.Any(), ubermock.Any()).DoAndReturn(func(ctx sdk.Context, fn func(id uint64, event zenbtctypes.BurnEvent) (bool, error)) error { return nil }).AnyTimes()
	zenBTCKeeper.EXPECT().WalkPendingMintTransactions(ubermock.Any(), ubermock.Any()).DoAndReturn(func(ctx sdk.Context, fn func(id uint64, tx zenbtctypes.PendingMintTransaction) (bool, error)) error {
		return nil
	}).AnyTimes()
	zenBTCKeeper.EXPECT().GetPendingMintTransaction(ubermock.Any(), ubermock.Any()).Return(zenbtctypes.PendingMintTransaction{}, nil).AnyTimes()
	zenBTCKeeper.EXPECT().HasPendingMintTransaction(ubermock.Any(), ubermock.Any()).Return(false, nil).AnyTimes()
	zenBTCKeeper.EXPECT().HasRedemption(ubermock.Any(), ubermock.Any()).Return(false, nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetRedemption(ubermock.Any(), ubermock.Any()).Return(zenbtctypes.Redemption{}, nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetBurnEvent(ubermock.Any(), ubermock.Any()).Return(zenbtctypes.BurnEvent{}, nil).AnyTimes()
	zenBTCKeeper.EXPECT().CreateBurnEvent(ubermock.Any(), ubermock.Any()).Return(uint64(1), nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetFirstPendingSolMintTransaction(ubermock.Any()).Return(uint64(1), nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetFirstPendingStakeTransaction(ubermock.Any()).Return(uint64(0), nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetFirstPendingEthMintTransaction(ubermock.Any()).Return(uint64(0), nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetFirstPendingBurnEvent(ubermock.Any()).Return(uint64(0), nil).AnyTimes()
	zenBTCKeeper.EXPECT().GetFirstPendingRedemption(ubermock.Any()).Return(uint64(0), nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetPendingMintTransaction(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetBurnEvent(ubermock.Any(), ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetRedemption(ubermock.Any(), ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetFirstPendingStakeTransaction(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetFirstPendingEthMintTransaction(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetFirstPendingSolMintTransaction(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetFirstPendingBurnEvent(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()
	zenBTCKeeper.EXPECT().SetFirstPendingRedemption(ubermock.Any(), ubermock.Any()).Return(nil).AnyTimes()

	// Create a proper ZRConfig for testing
	zrConfig := &params.ZRConfig{
		IsValidator: true, // Set to true for testing validator behavior
		SidecarAddr: "localhost:8080",
	}

	// Create a proper tx decoder for testing
	txDecoder := encCfg.TxConfig.TxDecoder()

	keeper := validationkeeper.NewKeeper(
		encCfg.Codec,
		storeService,
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		txDecoder,
		zrConfig,
		treasuryKeeper,
		zenBTCKeeper,
		zentpKeeper,
		slashingKeeper,
		address.NewBech32Codec("zenvaloper"),
		address.NewBech32Codec("zenvalcons"),
	)
	require.NoError(keeper.SetParams(ctx, validationtypes.DefaultParams()))

	// Set up mock sidecar client
	mockSidecarClient := validationtestutil.NewMocksidecarClient(newctrl)
	keeper.SetSidecarClient(mockSidecarClient)

	mockSidecarClient.EXPECT().GetSidecarState(ubermock.Any(), ubermock.Any()).Return(validationtestutil.SampleSidecarState, nil).AnyTimes()
	mockSidecarClient.EXPECT().GetSidecarStateByEthHeight(ubermock.Any(), ubermock.Any()).Return(validationtestutil.SampleSidecarState, nil).AnyTimes()
	mockSidecarClient.EXPECT().GetBitcoinBlockHeaderByHeight(ubermock.Any(), ubermock.Any()).Return(validationtestutil.SampleBtcHeader, nil).AnyTimes()
	mockSidecarClient.EXPECT().GetLatestBitcoinBlockHeader(ubermock.Any(), ubermock.Any()).Return(validationtestutil.SampleBtcHeader, nil).AnyTimes()

	mockSidecarClient.EXPECT().GetLatestEthereumNonceForAccount(ubermock.Any(), ubermock.Any()).Return(validationtestutil.SampleNonceResponse, nil).AnyTimes()

	mockSidecarClient.EXPECT().GetSolanaAccountInfo(ubermock.Any(), ubermock.Any()).Return(validationtestutil.SampleSolanaAccount, nil).AnyTimes()

	s.ctx = ctx
	s.validationKeeper = keeper
	s.bankKeeper = bankKeeper
	s.accountKeeper = accountKeeper

	stakingtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	validationtypes.RegisterQueryServer(queryHelper, validationkeeper.Querier{Keeper: keeper})
	s.queryClient = validationtypes.NewQueryClient(queryHelper)
	s.msgServer = validationkeeper.NewMsgServerImpl(keeper)

	// Store the zenBTC controller in the suite for cleanup
	s.zenBTCCtrl = newctrl

	return keeper, newctrl
}

func (s *ValidationKeeperTestSuite) SetupTest() {
	s.validationKeeper, _ = s.ValidationKeeperSetupTest()
}

func (s *KeeperTestSuite) TestParams() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	expParams := stakingtypes.DefaultParams()
	// check that the empty keeper loads the default
	resParams, err := keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(expParams, resParams)

	expParams.MaxValidators = 555
	expParams.MaxEntries = 111
	require.NoError(keeper.SetParams(ctx, expParams))
	resParams, err = keeper.GetParams(ctx)
	require.NoError(err)
	require.True(expParams.Equal(resParams))
}

func (s *KeeperTestSuite) TestLastTotalPower() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	expTotalPower := math.NewInt(10 ^ 9)
	require.NoError(keeper.SetLastTotalPower(ctx, expTotalPower))
	resTotalPower, err := keeper.GetLastTotalPower(ctx)
	require.NoError(err)
	require.True(expTotalPower.Equal(resTotalPower))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func TestValidationKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationKeeperTestSuite))
}

func ValidationTestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationKeeperTestSuite))
}

func (s *ValidationKeeperTestSuite) TestSetBackfillRequests() {
	ctx, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expBackfillRequests := validationtypes.BackfillRequests{
		Requests: []*validationtypes.MsgTriggerEventBackfill{
			{
				Authority:    keeper.GetAuthority(),
				TxHash:       "someHash",
				Caip2ChainId: "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1",
				EventType:    validationtypes.EventType_EVENT_TYPE_ZENBTC_BURN,
			},
		},
	}
	require.NoError(keeper.SetBackfillRequests(ctx, expBackfillRequests))
	resBackfillRequests, err := keeper.BackfillRequests.Get(ctx)
	require.NoError(err)
	require.Equal(expBackfillRequests, resBackfillRequests)
}

func (s *ValidationKeeperTestSuite) TestSetSolanaRequestedNonce() {
	ctx, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expNonce := uint64(1)
	require.NoError(keeper.SetSolanaRequestedNonce(ctx, expNonce, true))
	resNonce, err := keeper.SolanaNonceRequested.Get(ctx, expNonce)
	require.NoError(err)
	require.Equal(true, resNonce)
}

func (s *ValidationKeeperTestSuite) TestSetSolanaZenTPRequestedAccount() {
	ctx, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expAccount := "someAccount"
	require.NoError(keeper.SetSolanaZenTPRequestedAccount(ctx, expAccount, true))
	resNonce, err := keeper.SolanaZenTPAccountsRequested.Get(ctx, expAccount)
	require.NoError(err)
	require.Equal(true, resNonce)
}

func (s *ValidationKeeperTestSuite) TestSetSolanaZenBTCRequestedAccount() {
	ctx, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expAccount := "someAccount"
	require.NoError(keeper.SetSolanaZenBTCRequestedAccount(ctx, expAccount, true))
	resNonce, err := keeper.SolanaAccountsRequested.Get(ctx, expAccount)
	require.NoError(err)
	require.Equal(true, resNonce)
}

func (s *ValidationKeeperTestSuite) TestSetValidatorUpdates() {
	ctx, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expValidatorUpdates := []abci.ValidatorUpdate{
		{
			PubKey: cmtcrypto.PublicKey{
				Sum: &cmtcrypto.PublicKey_Ed25519{
					Ed25519: []byte("test_public_key_32_bytes_long_here"),
				},
			},
			Power: 100,
		},
	}
	require.NoError(keeper.SetValidatorUpdates(ctx, expValidatorUpdates))
	resValidatorUpdates, err := keeper.GetValidatorUpdates(ctx)
	require.NoError(err)
	require.Equal(expValidatorUpdates, resValidatorUpdates)
}

func (s *ValidationKeeperTestSuite) TestValidatorAddressCodec() {
	_, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expValidatorAddress := "zenvaloper138a4gyfjyghrd4pvuhuezxa6cl0wd5cde3s8rd"
	resValidatorAddressBytes, err := keeper.ValidatorAddressCodec().StringToBytes(expValidatorAddress)
	require.NoError(err)
	require.NotEmpty(resValidatorAddressBytes)

	resValidatorAddressStr, err := keeper.ValidatorAddressCodec().BytesToString(resValidatorAddressBytes)
	require.NoError(err)
	require.Equal(expValidatorAddress, resValidatorAddressStr)
}

func (s *ValidationKeeperTestSuite) TestConsensusAddressCodec() {
	_, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expConsensusAddress := "zenvalcons1jpnwkh0k75u2cyph8sn20s5tzkvt2n7csuhhlg"
	resConsensusAddressBytes, err := keeper.ConsensusAddressCodec().StringToBytes(expConsensusAddress)
	require.NoError(err)
	require.NotEmpty(resConsensusAddressBytes)

	resConsensusAddressStr, err := keeper.ConsensusAddressCodec().BytesToString(resConsensusAddressBytes)
	require.NoError(err)
	require.Equal(expConsensusAddress, resConsensusAddressStr)
}

func (s *ValidationKeeperTestSuite) TestGetLastTotalPower() {
	ctx, keeper := s.ctx, s.validationKeeper
	require := s.Require()

	expTotalPower := math.NewInt(10 ^ 9)
	require.NoError(keeper.SetLastTotalPower(ctx, expTotalPower))
	resTotalPower, err := keeper.GetLastTotalPower(ctx)
	require.NoError(err)
	require.True(expTotalPower.Equal(resTotalPower))
}

func (s *ValidationKeeperTestSuite) TestHooks() {
	ctrl := ubermock.NewController(s.T())

	tests := []struct {
		name     string
		hooks    validationtypes.StakingHooks
		expHooks validationtypes.StakingHooks
	}{
		{
			name:     "nil",
			hooks:    nil,
			expHooks: validationtypes.MultiStakingHooks{},
		},
		{
			name:     "mock",
			hooks:    validationtestutil.NewMockStakingHooks(ctrl),
			expHooks: validationtestutil.NewMockStakingHooks(ctrl),
		},
	}

	for _, test := range tests {

		_, keeper := s.ctx, s.validationKeeper
		require := s.Require()

		keeper.SetHooks(test.hooks)

		hooks := keeper.Hooks()
		require.NotNil(hooks)
		require.Equal(test.expHooks, hooks)
	}
}
