package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"github.com/Zenrock-Foundation/zrchain/v5/app/params"
	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/testutil/sample"
	treasuryModule "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type bankKeeperMock struct {
	transactions []struct {
		fromAddr sdk.AccAddress
		toAddr   sdk.AccAddress
		toModule string
		amount   sdk.Coins
	}
}

func newBankKeeperMock() *bankKeeperMock {
	return &bankKeeperMock{
		transactions: make([]struct {
			fromAddr sdk.AccAddress
			toAddr   sdk.AccAddress
			toModule string
			amount   sdk.Coins
		}, 0),
	}
}

func (b *bankKeeperMock) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	b.transactions = append(b.transactions, struct {
		fromAddr sdk.AccAddress
		toAddr   sdk.AccAddress
		toModule string
		amount   sdk.Coins
	}{fromAddr, toAddr, "", amt})
	return nil
}

func (b *bankKeeperMock) SendCoinsFromAccountToModule(ctx context.Context, fromAddr sdk.AccAddress, toModule string, amt sdk.Coins) error {
	b.transactions = append(b.transactions, struct {
		fromAddr sdk.AccAddress
		toAddr   sdk.AccAddress
		toModule string
		amount   sdk.Coins
	}{fromAddr, sdk.AccAddress{}, toModule, amt})
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

	// Generate addresses once and reuse them
	addrFrom := sample.AccAddress()
	addrTo := sample.AccAddress()

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
				feeAddr:    addrTo,
				fee:        1000,
				commission: 10,
			},
			want: want{
				fee:        900,
				commission: 100,
			},
		},
		{
			name:    "PASS: Send fee to module",
			wantErr: false,
			args: args{
				feeAddr:    types.KeyringCollectorName,
				fee:        1000,
				commission: 10,
			},
			want: want{
				fee:        900,
				commission: 100,
			},
		},
		{
			name:    "PASS: Zero commission",
			wantErr: false,
			args: args{
				feeAddr:    addrTo,
				fee:        1000,
				commission: 0,
			},
			want: want{
				fee:        1000,
				commission: 0,
			},
		},
		{
			name:    "PASS: 100% commission",
			wantErr: false,
			args: args{
				feeAddr:    addrTo,
				fee:        1000,
				commission: 100,
			},
			want: want{
				fee:        0,
				commission: 1000,
			},
		},
		{
			name:    "PASS: Small fee amount",
			wantErr: false,
			args: args{
				feeAddr:    addrTo,
				fee:        10,
				commission: 10,
			},
			want: want{
				fee:        9,
				commission: 1,
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

			tkGenesis.Params.KeyringCommission = tt.args.commission
			tkGenesis.Params.KeyringCommissionDestination = tt.args.feeAddr

			treasuryModule.InitGenesis(ctx, treasuryKeeper, tkGenesis)

			err := treasuryKeeper.SplitKeyringFee(ctx, addrFrom, tt.args.feeAddr, tt.args.fee)
			require.NoError(t, err)

			// Verify transactions
			require.Equal(t, 2, len(bkmock.transactions), "Expected exactly two transactions")

			// First transaction (commission)
			require.Equal(t,
				addrFrom,
				sdk.MustBech32ifyAddressBytes(sdk.Bech32MainPrefix, bkmock.transactions[0].fromAddr),
				"First transaction from address should be sender",
			)
			require.Equal(t,
				types.KeyringCollectorName,
				bkmock.transactions[0].toModule,
				"First transaction should go to KeyringCollector module",
			)
			require.Equal(t,
				tt.want.commission,
				bkmock.transactions[0].amount.AmountOf(params.BondDenom).Uint64(),
				"First transaction amount should be commission",
			)

			// Second transaction
			require.Equal(t,
				addrFrom,
				sdk.MustBech32ifyAddressBytes(sdk.Bech32MainPrefix, bkmock.transactions[1].fromAddr),
				"Second transaction from address should be sender",
			)
			if tt.args.feeAddr == types.KeyringCollectorName {
				require.Equal(t,
					types.KeyringCollectorName,
					bkmock.transactions[1].toModule,
					"Second transaction should go to KeyringCollector module",
				)
			} else {
				require.Equal(t,
					tt.args.feeAddr,
					sdk.MustBech32ifyAddressBytes(sdk.Bech32MainPrefix, bkmock.transactions[1].toAddr),
					"Second transaction should go to recipient",
				)
			}
			require.Equal(t, tt.want.fee, bkmock.transactions[1].amount.AmountOf(params.BondDenom).Uint64(), "Second transaction amount should be remaining fee")
		})
	}
}
