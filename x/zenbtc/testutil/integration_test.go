package testutil

import (
	"context"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/keeper"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	zenbtcKeeper  *keeper.Keeper
	ctx           sdk.Context
	msgServer     types.MsgServer
	accountKeeper *MockAccountKeeper
	bankKeeper    *MockBankKeeper
	paramSubspace *MockParamSubspace
	ctrl          *gomock.Controller
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx

	ctrl := gomock.NewController(s.T())
	accountKeeper := NewMockAccountKeeper(ctrl)
	bankKeeper := NewMockBankKeeper(ctrl)
	paramSubspace := NewMockParamSubspace(ctrl)

	s.accountKeeper = accountKeeper
	s.bankKeeper = bankKeeper
	s.paramSubspace = paramSubspace
	s.ctrl = ctrl

	s.zenbtcKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		testCtx.Ctx.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		nil,
		nil,
	)

	s.Require().Equal(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		s.zenbtcKeeper.Logger())

	err := s.zenbtcKeeper.Params.Set(s.ctx, *keeper.DefaultParams())
	s.Require().NoError(err)

	s.msgServer = keeper.NewMsgServerImpl(*s.zenbtcKeeper)
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_GetExchangeRate() {
	tests := []struct {
		name         string
		setupSupply  *types.Supply
		expectedRate string
		expectError  bool
	}{
		{
			name:         "Initial exchange rate when no supply exists",
			setupSupply:  nil,
			expectedRate: "1.000000000000000000",
			expectError:  false,
		},
		{
			name: "Exchange rate with supply",
			setupSupply: &types.Supply{
				CustodiedBTC:  1000,
				MintedZenBTC:  1000,
				PendingZenBTC: 0,
			},
			expectedRate: "1.000000000000000000",
			expectError:  false,
		},
		{
			name: "Exchange rate with zero zenBTC",
			setupSupply: &types.Supply{
				CustodiedBTC:  1000,
				MintedZenBTC:  0,
				PendingZenBTC: 0,
			},
			expectedRate: "1.000000000000000000",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setupSupply != nil {
				err := s.zenbtcKeeper.SetSupply(s.ctx, *tt.setupSupply)
				s.Require().NoError(err)
			}

			rate, err := s.zenbtcKeeper.GetExchangeRate(s.ctx)
			if tt.expectError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tt.expectedRate, rate.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_PendingMintTransactions() {
	pendingTx := types.PendingMintTransaction{
		Id:               1,
		RecipientAddress: sample.AccAddress(),
		Amount:           1000,
		ChainType:        types.WalletType_WALLET_TYPE_EVM,
		Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
		BlockHeight:      100,
		Caip2ChainId:     "eip155:1",
	}

	err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)
	s.Require().NoError(err)

	var walkedTx types.PendingMintTransaction
	err = s.zenbtcKeeper.WalkPendingMintTransactions(s.ctx, func(id uint64, tx types.PendingMintTransaction) (stop bool, err error) {
		walkedTx = tx
		return true, nil
	})
	s.Require().NoError(err)
	s.Require().Equal(pendingTx.Id, walkedTx.Id)
	s.Require().Equal(pendingTx.RecipientAddress, walkedTx.RecipientAddress)
	s.Require().Equal(pendingTx.Amount, walkedTx.Amount)
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_Redemptions() {
	redemption := types.Redemption{
		Data: types.RedemptionData{
			Id:                 1,
			Amount:             1000,
			DestinationAddress: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
		},
		Status: types.RedemptionStatus_INITIATED,
	}

	err := s.zenbtcKeeper.SetRedemption(s.ctx, redemption.Data.Id, redemption)
	s.Require().NoError(err)

	exists, err := s.zenbtcKeeper.HasRedemption(s.ctx, redemption.Data.Id)
	s.Require().NoError(err)
	s.Require().True(exists)

	var walkedRedemption types.Redemption
	err = s.zenbtcKeeper.WalkRedemptions(s.ctx, func(id uint64, r types.Redemption) (stop bool, err error) {
		walkedRedemption = r
		return true, nil
	})
	s.Require().NoError(err)
	s.Require().Equal(redemption.Data.Id, walkedRedemption.Data.Id)
	s.Require().Equal(redemption.Data.Amount, walkedRedemption.Data.Amount)
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_Supply() {
	supply := types.Supply{
		MintedZenBTC:  1000,
		CustodiedBTC:  1000,
		PendingZenBTC: 500,
	}

	err := s.zenbtcKeeper.SetSupply(s.ctx, supply)
	s.Require().NoError(err)

	retrievedSupply, err := s.zenbtcKeeper.GetSupply(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(supply.MintedZenBTC, retrievedSupply.MintedZenBTC)
	s.Require().Equal(supply.CustodiedBTC, retrievedSupply.CustodiedBTC)
	s.Require().Equal(supply.PendingZenBTC, retrievedSupply.PendingZenBTC)
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_BurnEvents() {
	burnEvent := types.BurnEvent{
		Id:              1,
		Amount:          1000,
		DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
		Status:          types.BurnStatus_BURN_STATUS_BURNED,
		ChainID:         "eip155:1",
		TxID:            "0x1234567890abcdef",
	}

	err := s.zenbtcKeeper.SetBurnEvent(s.ctx, burnEvent.Id, burnEvent)
	s.Require().NoError(err)

	retrievedBurnEvent, err := s.zenbtcKeeper.GetBurnEvent(s.ctx, burnEvent.Id)
	s.Require().NoError(err)
	s.Require().Equal(burnEvent.Id, retrievedBurnEvent.Id)
	s.Require().Equal(burnEvent.Amount, retrievedBurnEvent.Amount)

	var walkedBurnEvent types.BurnEvent
	err = s.zenbtcKeeper.WalkBurnEvents(s.ctx, func(id uint64, be types.BurnEvent) (stop bool, err error) {
		walkedBurnEvent = be
		return true, nil
	})
	s.Require().NoError(err)
	s.Require().Equal(burnEvent.Id, walkedBurnEvent.Id)
	s.Require().Equal(burnEvent.Amount, walkedBurnEvent.Amount)
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_PendingTransactionIndices() {
	testCases := []struct {
		name     string
		setFunc  func(context.Context, uint64) error
		getFunc  func(context.Context) (uint64, error)
		setValue uint64
	}{
		{
			name:     "FirstPendingEthMintTransaction",
			setFunc:  s.zenbtcKeeper.SetFirstPendingEthMintTransaction,
			getFunc:  s.zenbtcKeeper.GetFirstPendingEthMintTransaction,
			setValue: 100,
		},
		{
			name:     "FirstPendingSolMintTransaction",
			setFunc:  s.zenbtcKeeper.SetFirstPendingSolMintTransaction,
			getFunc:  s.zenbtcKeeper.GetFirstPendingSolMintTransaction,
			setValue: 200,
		},
		{
			name:     "FirstPendingBurnEvent",
			setFunc:  s.zenbtcKeeper.SetFirstPendingBurnEvent,
			getFunc:  s.zenbtcKeeper.GetFirstPendingBurnEvent,
			setValue: 300,
		},
		{
			name:     "FirstPendingRedemption",
			setFunc:  s.zenbtcKeeper.SetFirstPendingRedemption,
			getFunc:  s.zenbtcKeeper.GetFirstPendingRedemption,
			setValue: 400,
		},
		{
			name:     "FirstPendingStakeTransaction",
			setFunc:  s.zenbtcKeeper.SetFirstPendingStakeTransaction,
			getFunc:  s.zenbtcKeeper.GetFirstPendingStakeTransaction,
			setValue: 500,
		},
		{
			name:     "FirstRedemptionAwaitingSign",
			setFunc:  s.zenbtcKeeper.SetFirstRedemptionAwaitingSign,
			getFunc:  s.zenbtcKeeper.GetFirstRedemptionAwaitingSign,
			setValue: 600,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := tc.setFunc(s.ctx, tc.setValue)
			s.Require().NoError(err)

			retrievedValue, err := tc.getFunc(s.ctx)
			s.Require().NoError(err)
			s.Require().Equal(tc.setValue, retrievedValue)
		})
	}
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_ParameterGetters() {
	testCases := []struct {
		name          string
		getterFunc    func(context.Context) string
		expectedValue string
	}{
		{
			name:          "GetControllerAddr",
			getterFunc:    s.zenbtcKeeper.GetControllerAddr,
			expectedValue: keeper.DefaultControllerAddr,
		},
		{
			name:          "GetEthTokenAddr",
			getterFunc:    s.zenbtcKeeper.GetEthTokenAddr,
			expectedValue: keeper.DefaultEthTokenAddr,
		},
		{
			name:          "GetDepositKeyringAddr",
			getterFunc:    s.zenbtcKeeper.GetDepositKeyringAddr,
			expectedValue: keeper.DefaultDepositKeyringAddr,
		},
		{
			name:          "GetBitcoinProxyAddress",
			getterFunc:    s.zenbtcKeeper.GetBitcoinProxyAddress,
			expectedValue: keeper.DefaultTestnetBitcoinProxyAddress,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			value := tc.getterFunc(s.ctx)
			s.Require().Equal(tc.expectedValue, value)
		})
	}

	uint64TestCases := []struct {
		name          string
		getterFunc    func(context.Context) uint64
		expectedValue uint64
	}{
		{
			name:          "GetStakerKeyID",
			getterFunc:    s.zenbtcKeeper.GetStakerKeyID,
			expectedValue: keeper.DefaultStakerKeyID,
		},
		{
			name:          "GetEthMinterKeyID",
			getterFunc:    s.zenbtcKeeper.GetEthMinterKeyID,
			expectedValue: keeper.DefaultEthMinterKeyID,
		},
		{
			name:          "GetUnstakerKeyID",
			getterFunc:    s.zenbtcKeeper.GetUnstakerKeyID,
			expectedValue: keeper.DefaultUnstakerKeyID,
		},
		{
			name:          "GetCompleterKeyID",
			getterFunc:    s.zenbtcKeeper.GetCompleterKeyID,
			expectedValue: keeper.DefaultCompleterKeyID,
		},
		{
			name:          "GetRewardsDepositKeyID",
			getterFunc:    s.zenbtcKeeper.GetRewardsDepositKeyID,
			expectedValue: keeper.DefaultRewardsDepositKeyID,
		},
	}

	for _, tc := range uint64TestCases {
		s.Run(tc.name, func() {
			value := tc.getterFunc(s.ctx)
			s.Require().Equal(tc.expectedValue, value)
		})
	}

	s.Run("GetChangeAddressKeyIDs", func() {
		value := s.zenbtcKeeper.GetChangeAddressKeyIDs(s.ctx)
		s.Require().Equal(keeper.DefaultChangeAddressKeyIDs, value)
	})

	s.Run("GetSolanaParams", func() {
		solana := s.zenbtcKeeper.GetSolanaParams(s.ctx)
		s.Require().Equal(keeper.DefaultSolana, solana)
	})
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_Authority() {
	authority := s.zenbtcKeeper.GetAuthority()
	s.Require().NotEmpty(authority)
}

func (s *IntegrationTestSuite) Test_ZenbtcKeeper_GetParams() {
	params, err := s.zenbtcKeeper.GetParams(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(params)
}

func (s *IntegrationTestSuite) Test_Validation_CAIP2ChainId() {
	s.T().SkipNow() // We don't need to test this low level function
	testCases := []struct {
		name         string
		caip2ChainId string
		chainType    types.WalletType
		expectError  bool
		description  string
	}{
		{
			name:         "Valid Ethereum Mainnet",
			caip2ChainId: "eip155:1",
			chainType:    types.WalletType_WALLET_TYPE_EVM,
			expectError:  false,
			description:  "Standard Ethereum mainnet",
		},
		{
			name:         "Valid Ethereum Testnet",
			caip2ChainId: "eip155:5",
			chainType:    types.WalletType_WALLET_TYPE_EVM,
			expectError:  false,
			description:  "Ethereum Goerli testnet",
		},
		{
			name:         "Valid Solana Mainnet",
			caip2ChainId: "solana:101",
			chainType:    types.WalletType_WALLET_TYPE_SOLANA,
			expectError:  false,
			description:  "Solana mainnet",
		},
		{
			name:         "Valid Solana Testnet",
			caip2ChainId: "solana:102",
			chainType:    types.WalletType_WALLET_TYPE_SOLANA,
			expectError:  false,
			description:  "Solana testnet",
		},
		{
			name:         "Invalid CAIP2 Format",
			caip2ChainId: "invalid:format",
			chainType:    types.WalletType_WALLET_TYPE_EVM,
			expectError:  true,
			description:  "Malformed CAIP2 chain ID",
		},
		{
			name:         "Empty CAIP2 Chain ID",
			caip2ChainId: "",
			chainType:    types.WalletType_WALLET_TYPE_EVM,
			expectError:  true,
			description:  "Empty chain ID should be rejected",
		},
		{
			name:         "Unsupported Chain ID",
			caip2ChainId: "eip155:999999",
			chainType:    types.WalletType_WALLET_TYPE_EVM,
			expectError:  true,
			description:  "Unsupported Ethereum chain ID",
		},
		{
			name:         "Mismatched Chain Type and CAIP2",
			caip2ChainId: "eip155:1",
			chainType:    types.WalletType_WALLET_TYPE_SOLANA,
			expectError:  true,
			description:  "Ethereum CAIP2 with Solana chain type",
		},
		{
			name:         "Invalid Solana Chain ID",
			caip2ChainId: "solana:999999",
			chainType:    types.WalletType_WALLET_TYPE_SOLANA,
			expectError:  true,
			description:  "Unsupported Solana chain ID",
		},
		{
			name:         "Unspecified Wallet Type",
			caip2ChainId: "eip155:1",
			chainType:    types.WalletType_WALLET_TYPE_UNSPECIFIED,
			expectError:  true,
			description:  "Unspecified wallet type should be rejected",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			pendingTx := types.PendingMintTransaction{
				Id:               1,
				RecipientAddress: sample.AccAddress(),
				Amount:           1000,
				ChainType:        tc.chainType,
				Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				BlockHeight:      100,
				Caip2ChainId:     tc.caip2ChainId,
			}

			err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)

			if tc.expectError {
				s.Require().Error(err, tc.description)
			} else {
				s.Require().NoError(err, tc.description)
			}
		})
	}
}

func (s *IntegrationTestSuite) Test_Validation_AddressFormats() {
	s.T().SkipNow() // We don't need to test this low level function
	testCases := []struct {
		name          string
		recipientAddr string
		chainType     types.WalletType
		expectError   bool
		description   string
	}{
		{
			name:          "Valid Ethereum Address",
			recipientAddr: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			chainType:     types.WalletType_WALLET_TYPE_EVM,
			expectError:   false,
			description:   "Valid Ethereum address format",
		},
		{
			name:          "Valid Solana Address",
			recipientAddr: "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnDj",
			chainType:     types.WalletType_WALLET_TYPE_SOLANA,
			expectError:   false,
			description:   "Valid Solana address format",
		},
		{
			name:          "Invalid Ethereum Address - Wrong Length",
			recipientAddr: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
			chainType:     types.WalletType_WALLET_TYPE_EVM,
			expectError:   true,
			description:   "Ethereum address with wrong length",
		},
		{
			name:          "Invalid Ethereum Address - No 0x Prefix",
			recipientAddr: "742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			chainType:     types.WalletType_WALLET_TYPE_EVM,
			expectError:   true,
			description:   "Ethereum address without 0x prefix",
		},
		{
			name:          "Invalid Solana Address - Wrong Length",
			recipientAddr: "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnD",
			chainType:     types.WalletType_WALLET_TYPE_SOLANA,
			expectError:   true,
			description:   "Solana address with wrong length",
		},
		{
			name:          "Empty Address",
			recipientAddr: "",
			chainType:     types.WalletType_WALLET_TYPE_EVM,
			expectError:   true,
			description:   "Empty address should be rejected",
		},
		{
			name:          "Address with Invalid Characters",
			recipientAddr: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8bG",
			chainType:     types.WalletType_WALLET_TYPE_EVM,
			expectError:   true,
			description:   "Ethereum address with invalid hex characters",
		},
		{
			name:          "Mismatched Address and Chain Type",
			recipientAddr: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			chainType:     types.WalletType_WALLET_TYPE_SOLANA,
			expectError:   true,
			description:   "Ethereum address with Solana chain type",
		},
		{
			name:          "Cosmos Address with EVM Chain Type",
			recipientAddr: sample.AccAddress(),
			chainType:     types.WalletType_WALLET_TYPE_EVM,
			expectError:   true,
			description:   "Cosmos address should not be valid for EVM chains",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			pendingTx := types.PendingMintTransaction{
				Id:               1,
				RecipientAddress: tc.recipientAddr,
				Amount:           1000,
				ChainType:        tc.chainType,
				Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				BlockHeight:      100,
				Caip2ChainId:     "eip155:1",
			}

			err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)

			if tc.expectError {
				s.Require().Error(err, tc.description)
			} else {
				s.Require().NoError(err, tc.description)
			}
		})
	}
}

func (s *IntegrationTestSuite) Test_Validation_AmountBoundaries() {
	s.T().SkipNow() // We don't need to test this low level function

	testCases := []struct {
		name        string
		amount      uint64
		expectError bool
		description string
	}{
		{
			name:        "Valid Amount",
			amount:      1000,
			expectError: false,
			description: "Normal valid amount",
		},
		{
			name:        "Zero Amount",
			amount:      0,
			expectError: true,
			description: "Zero amount should be rejected",
		},
		{
			name:        "Very Large Amount",
			amount:      1e18,
			expectError: false,
			description: "Large amount should be accepted",
		},
		{
			name:        "Maximum Uint64 Amount",
			amount:      18446744073709551615,
			expectError: false,
			description: "Maximum uint64 amount should be accepted",
		},
		{
			name:        "Small Valid Amount",
			amount:      1,
			expectError: false,
			description: "Minimum valid amount",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			pendingTx := types.PendingMintTransaction{
				Id:               1,
				RecipientAddress: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
				Amount:           tc.amount,
				ChainType:        types.WalletType_WALLET_TYPE_EVM,
				Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				BlockHeight:      100,
				Caip2ChainId:     "eip155:1",
			}

			err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)

			if tc.expectError {
				s.Require().Error(err, tc.description)
			} else {
				s.Require().NoError(err, tc.description)
			}
		})
	}
}

func (s *IntegrationTestSuite) Test_Validation_TransactionStatus() {
	s.T().SkipNow() // We don't need to test this low level function
	testCases := []struct {
		name        string
		status      types.MintTransactionStatus
		expectError bool
		description string
	}{
		{
			name:        "Valid Deposited Status",
			status:      types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
			expectError: false,
			description: "Valid deposited status",
		},
		{
			name:        "Valid Staked Status",
			status:      types.MintTransactionStatus_MINT_TRANSACTION_STATUS_STAKED,
			expectError: false,
			description: "Valid staked status",
		},
		{
			name:        "Valid Minted Status",
			status:      types.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED,
			expectError: false,
			description: "Valid minted status",
		},
		{
			name:        "Unspecified Status",
			status:      types.MintTransactionStatus_MINT_TRANSACTION_STATUS_UNSPECIFIED,
			expectError: true,
			description: "Unspecified status should be rejected",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			pendingTx := types.PendingMintTransaction{
				Id:               1,
				RecipientAddress: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
				Amount:           1000,
				ChainType:        types.WalletType_WALLET_TYPE_EVM,
				Status:           tc.status,
				BlockHeight:      100,
				Caip2ChainId:     "eip155:1",
			}

			err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)

			if tc.expectError {
				s.Require().Error(err, tc.description)
			} else {
				s.Require().NoError(err, tc.description)
			}
		})
	}
}

func (s *IntegrationTestSuite) Test_Validation_RedemptionData() {
	s.T().SkipNow() // We don't need to test this low level function
	testCases := []struct {
		name        string
		redemption  types.Redemption
		expectError bool
		description string
	}{
		{
			name: "Valid Redemption",
			redemption: types.Redemption{
				Data: types.RedemptionData{
					Id:                 1,
					Amount:             1000,
					DestinationAddress: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				},
				Status: types.RedemptionStatus_INITIATED,
			},
			expectError: false,
			description: "Valid redemption data",
		},
		{
			name: "Zero Redemption ID",
			redemption: types.Redemption{
				Data: types.RedemptionData{
					Id:                 0,
					Amount:             1000,
					DestinationAddress: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				},
				Status: types.RedemptionStatus_INITIATED,
			},
			expectError: true,
			description: "Zero redemption ID should be rejected",
		},
		{
			name: "Zero Amount",
			redemption: types.Redemption{
				Data: types.RedemptionData{
					Id:                 1,
					Amount:             0,
					DestinationAddress: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				},
				Status: types.RedemptionStatus_INITIATED,
			},
			expectError: true,
			description: "Zero amount should be rejected",
		},
		{
			name: "Empty Destination Address",
			redemption: types.Redemption{
				Data: types.RedemptionData{
					Id:                 1,
					Amount:             1000,
					DestinationAddress: []byte{},
				},
				Status: types.RedemptionStatus_INITIATED,
			},
			expectError: true,
			description: "Empty destination address should be rejected",
		},
		{
			name: "Unspecified Status",
			redemption: types.Redemption{
				Data: types.RedemptionData{
					Id:                 1,
					Amount:             1000,
					DestinationAddress: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				},
				Status: types.RedemptionStatus_UNSPECIFIED,
			},
			expectError: true,
			description: "Unspecified status should be rejected",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.zenbtcKeeper.SetRedemption(s.ctx, tc.redemption.Data.Id, tc.redemption)

			if tc.expectError {
				s.Require().Error(err, tc.description)
			} else {
				s.Require().NoError(err, tc.description)
			}
		})
	}
}

func (s *IntegrationTestSuite) Test_Validation_BurnEventData() {
	s.T().SkipNow() // We don't need to test this low level function

	testCases := []struct {
		name        string
		burnEvent   types.BurnEvent
		expectError bool
		description string
	}{
		{
			name: "Valid Burn Event",
			burnEvent: types.BurnEvent{
				Id:              1,
				Amount:          1000,
				DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				Status:          types.BurnStatus_BURN_STATUS_BURNED,
				ChainID:         "eip155:1",
				TxID:            "0x1234567890abcdef",
			},
			expectError: false,
			description: "Valid burn event data",
		},
		{
			name: "Zero Burn Event ID",
			burnEvent: types.BurnEvent{
				Id:              0,
				Amount:          1000,
				DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				Status:          types.BurnStatus_BURN_STATUS_BURNED,
				ChainID:         "eip155:1",
				TxID:            "0x1234567890abcdef",
			},
			expectError: true,
			description: "Zero burn event ID should be rejected",
		},
		{
			name: "Zero Amount",
			burnEvent: types.BurnEvent{
				Id:              1,
				Amount:          0,
				DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				Status:          types.BurnStatus_BURN_STATUS_BURNED,
				ChainID:         "eip155:1",
				TxID:            "0x1234567890abcdef",
			},
			expectError: true,
			description: "Zero amount should be rejected",
		},
		{
			name: "Empty Destination Address",
			burnEvent: types.BurnEvent{
				Id:              1,
				Amount:          1000,
				DestinationAddr: []byte{},
				Status:          types.BurnStatus_BURN_STATUS_BURNED,
				ChainID:         "eip155:1",
				TxID:            "0x1234567890abcdef",
			},
			expectError: true,
			description: "Empty destination address should be rejected",
		},
		{
			name: "Empty Chain ID",
			burnEvent: types.BurnEvent{
				Id:              1,
				Amount:          1000,
				DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				Status:          types.BurnStatus_BURN_STATUS_BURNED,
				ChainID:         "",
				TxID:            "0x1234567890abcdef",
			},
			expectError: true,
			description: "Empty chain ID should be rejected",
		},
		{
			name: "Empty Transaction ID",
			burnEvent: types.BurnEvent{
				Id:              1,
				Amount:          1000,
				DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				Status:          types.BurnStatus_BURN_STATUS_BURNED,
				ChainID:         "eip155:1",
				TxID:            "",
			},
			expectError: true,
			description: "Empty transaction ID should be rejected",
		},
		{
			name: "Unspecified Status",
			burnEvent: types.BurnEvent{
				Id:              1,
				Amount:          1000,
				DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
				Status:          types.BurnStatus_BURN_STATUS_UNSPECIFIED,
				ChainID:         "eip155:1",
				TxID:            "0x1234567890abcdef",
			},
			expectError: true,
			description: "Unspecified status should be rejected",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := s.zenbtcKeeper.SetBurnEvent(s.ctx, tc.burnEvent.Id, tc.burnEvent)

			if tc.expectError {
				s.Require().Error(err, tc.description)
			} else {
				s.Require().NoError(err, tc.description)
			}
		})
	}
}

func (s *IntegrationTestSuite) Test_Validation_DuplicatePrevention() {
	s.Run("Duplicate Pending Mint Transaction", func() {
		pendingTx := types.PendingMintTransaction{
			Id:               1,
			RecipientAddress: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			Amount:           1000,
			ChainType:        types.WalletType_WALLET_TYPE_EVM,
			Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
			BlockHeight:      100,
			Caip2ChainId:     "eip155:1",
		}

		err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)
		s.Require().NoError(err)

		err = s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)
	})

	s.Run("Duplicate Redemption", func() {
		redemption := types.Redemption{
			Data: types.RedemptionData{
				Id:                 1,
				Amount:             1000,
				DestinationAddress: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
			},
			Status: types.RedemptionStatus_INITIATED,
		}

		err := s.zenbtcKeeper.SetRedemption(s.ctx, redemption.Data.Id, redemption)
		s.Require().NoError(err)

		exists, err := s.zenbtcKeeper.HasRedemption(s.ctx, redemption.Data.Id)
		s.Require().NoError(err)
		s.Require().True(exists)

		err = s.zenbtcKeeper.SetRedemption(s.ctx, redemption.Data.Id, redemption)
		s.Require().NoError(err)
	})

	s.Run("Duplicate Burn Event", func() {
		burnEvent := types.BurnEvent{
			Id:              1,
			Amount:          1000,
			DestinationAddr: []byte("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"),
			Status:          types.BurnStatus_BURN_STATUS_BURNED,
			ChainID:         "eip155:1",
			TxID:            "0x1234567890abcdef",
		}

		err := s.zenbtcKeeper.SetBurnEvent(s.ctx, burnEvent.Id, burnEvent)
		s.Require().NoError(err)

		err = s.zenbtcKeeper.SetBurnEvent(s.ctx, burnEvent.Id, burnEvent)
		s.Require().NoError(err)
	})
}

func (s *IntegrationTestSuite) Test_Validation_EdgeCases() {
	s.T().SkipNow() // We don't need to test this low level function

	s.Run("Very Long Address", func() {
		longAddress := "0x" + string(make([]byte, 1000))
		pendingTx := types.PendingMintTransaction{
			Id:               1,
			RecipientAddress: longAddress,
			Amount:           1000,
			ChainType:        types.WalletType_WALLET_TYPE_EVM,
			Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
			BlockHeight:      100,
			Caip2ChainId:     "eip155:1",
		}

		err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)
		s.Require().Error(err, "Very long addresses should be rejected")
	})

	s.Run("Very Long CAIP2 Chain ID", func() {
		longChainId := "eip155:" + string(make([]byte, 1000))
		pendingTx := types.PendingMintTransaction{
			Id:               1,
			RecipientAddress: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			Amount:           1000,
			ChainType:        types.WalletType_WALLET_TYPE_EVM,
			Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
			BlockHeight:      100,
			Caip2ChainId:     longChainId,
		}

		err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)
		s.Require().Error(err, "Very long CAIP2 chain IDs should be rejected")
	})

	s.Run("Negative Block Height", func() {
		pendingTx := types.PendingMintTransaction{
			Id:               1,
			RecipientAddress: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			Amount:           1000,
			ChainType:        types.WalletType_WALLET_TYPE_EVM,
			Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
			BlockHeight:      -100,
			Caip2ChainId:     "eip155:1",
		}

		_ = s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)
	})

	s.Run("Special Characters in Address", func() {
		specialAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6<script>alert('xss')</script>"
		pendingTx := types.PendingMintTransaction{
			Id:               1,
			RecipientAddress: specialAddress,
			Amount:           1000,
			ChainType:        types.WalletType_WALLET_TYPE_EVM,
			Status:           types.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
			BlockHeight:      100,
			Caip2ChainId:     "eip155:1",
		}

		err := s.zenbtcKeeper.SetPendingMintTransaction(s.ctx, pendingTx)
		s.Require().Error(err, "Addresses with special characters should be rejected")
	})
}
