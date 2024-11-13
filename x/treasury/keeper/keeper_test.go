package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	treasuryModule "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bankKeeperMock struct {
	SendCoinsCalled                    int
	SendCoinsFromAccountToModuleCalled int

	FromAddrs []sdk.AccAddress
	ToAddrs   []sdk.AccAddress
	Amounts   []sdk.Coins
	Modules   []string
}

func newBankKeeperMock() *bankKeeperMock {
	return &bankKeeperMock{
		SendCoinsCalled:                    0,
		SendCoinsFromAccountToModuleCalled: 0,
		FromAddrs:                          []sdk.AccAddress{},
		ToAddrs:                            []sdk.AccAddress{},
		Amounts:                            []sdk.Coins{},
		Modules:                            []string{},
	}
}

func (m *bankKeeperMock) SendCoins(_ context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	m.SendCoinsCalled++

	m.FromAddrs = append(m.FromAddrs, fromAddr)
	m.ToAddrs = append(m.ToAddrs, toAddr)
	m.Amounts = append(m.Amounts, amt)

	return nil
}

func (m *bankKeeperMock) SendCoinsFromAccountToModule(_ context.Context, fromAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	m.SendCoinsFromAccountToModuleCalled++

	m.FromAddrs = append(m.FromAddrs, fromAddr)
	m.Modules = append(m.Modules, recipientModule)
	m.Amounts = append(m.Amounts, amt)

	return nil
}

func Test_TreasuryKeeper_splitKeyringFee(t *testing.T) {
	type args struct {
		feeAddr    string
		fee        uint64
		commission uint64
	}
	type want struct {
		fee        uint64
		commission uint64
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "PASS: Send fee to address",
			wantErr: false,
			args: args{
				feeAddr:    sample.AccAddress(),
				fee:        1000,
				commission: 10,
			},
			want: want{
				fee:        900,
				commission: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := dbm.NewMemDB()
			stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

			bkmock := newBankKeeperMock()
			policyKeeper, ctx := keepertest.PolicyKeeper(t, db, stateStore, nil)
			identityKeeper, _ := keepertest.IdentityKeeper(t, &policyKeeper, db, stateStore)
			treasuryKeeper, _ := keepertest.TreasuryKeeper(t, &policyKeeper, &identityKeeper, bkmock, db, stateStore)

			tkGenesis := types.GenesisState{
				Params: types.DefaultParams(),
			}

			addrFrom := sample.AccAddress()
			addrTo := sample.AccAddress()

			tkGenesis.Params.KeyringCommission = tt.args.commission
			tkGenesis.Params.KeyringCommissionDestination = tt.args.feeAddr

			treasuryModule.InitGenesis(ctx, treasuryKeeper, tkGenesis)

			err := treasuryKeeper.SplitKeyringFee(ctx, addrFrom, addrTo, 1000)
			require.Nil(t, err)

			assert.Equal(t, 2, bkmock.SendCoinsCalled)
			assert.Equal(t, sdk.MustAccAddressFromBech32(tt.args.feeAddr), bkmock.ToAddrs[0])
			assert.Equal(t, sdk.MustAccAddressFromBech32(addrTo), bkmock.ToAddrs[1])

			assert.Equal(t, sdk.MustAccAddressFromBech32(addrFrom), bkmock.FromAddrs[0])
			assert.Equal(t, sdk.MustAccAddressFromBech32(addrFrom), bkmock.FromAddrs[1])

			assert.Equal(t, tt.want.commission, bkmock.Amounts[0].AmountOf("urock").Uint64())
			assert.Equal(t, tt.want.fee, bkmock.Amounts[1].AmountOf("urock").Uint64())
		})
	}
}
