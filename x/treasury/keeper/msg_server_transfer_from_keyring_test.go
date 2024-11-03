package keeper_test

import (
	"fmt"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_TransferFromKeyring(t *testing.T) {

	var defaultKr2 = idTypes.Keyring{
		Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
		Creator:     "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
		Description: "testDescription",
		Admins:      []string{"zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"},
		Parties:     []string{"zen10kmgv5gzygnecf46x092ecfe5xcvvv9rdaxmts"},
		KeyReqFee:   0,
		SigReqFee:   0,
		IsActive:    true,
	}

	type args struct {
		keyring *idTypes.Keyring
		msg     *types.MsgTransferFromKeyring
	}
	tests := []struct {
		name                    string
		args                    args
		wantTransferFromKeyring *types.MsgTransferFromKeyring
		want                    *types.MsgTransferFromKeyringResponse
		wantErr                 bool
	}{
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring: &defaultKr2,
				msg:     types.NewMsgTransferFromKeyring("zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty", "notthekeyring", "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty", uint64(1), "urock"),
			},
			wantErr: true,
		},
		{
			name: "FAIL: creator is not admin",
			args: args{
				keyring: &defaultKr2,
				msg:     types.NewMsgTransferFromKeyring("noAdmin", "keyring1pfnq7r04rept47gaf5cpdew2", "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty", uint64(1), "urock"),
			},
			wantErr: true,
		},
		{
			name: "FAIL: recipient is not admin",
			args: args{
				keyring: &defaultKr2,
				msg:     types.NewMsgTransferFromKeyring("zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty", "keyring1pfnq7r04rept47gaf5cpdew2", "noadmin", uint64(1), "urock"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			tk := keepers.TreasuryKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*tk)

			idGenesis := idTypes.GenesisState{
				PortId:   idTypes.PortID,
				Keyrings: []idTypes.Keyring{*tt.args.keyring},
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			trGenesis := types.GenesisState{
				PortId: types.PortID,
			}
			treasury.InitGenesis(ctx, *tk, trGenesis)

			got, err := msgSer.TransferFromKeyring(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
				fmt.Println(ctx.EventManager().Events())
			}
		})
	}
}
