package keeper_test

import (
	"encoding/hex"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
)

var zrSignDefaultKey = types.Key{
	Id:            1,
	WorkspaceAddr: "workspace1xklrytgff7w32j52v34w36",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	PublicKey:     []byte{0x02, 0x3b, 0x4d, 0x34, 0x2a, 0x6a, 0x52, 0x7c, 0xca, 0x8c, 0x2d, 0x92, 0x6a, 0xce, 0x76, 0x36, 0xa1, 0x9f, 0x09, 0x25, 0x28, 0x74, 0xe2, 0x36, 0x8f, 0xe7, 0xa3, 0x25, 0xd5, 0x0e, 0xff, 0xed, 0xb3},
	Index:         0,
}

var zrSignWorkspace = idTypes.Workspace{
	Address:         "workspace1mphgzyhncnzyggfxmv4nmh",
	Creator:         "0xbe27E519A220D4336A5F9d08C7fa542175E30dA0",
	Owners:          []string{"0xbe27E519A220D4336A5F9d08C7fa542175E30dA0"},
	ChildWorkspaces: []string{"workspace1xklrytgff7w32j52v34w36"},
	AdminPolicyId:   0,
	SignPolicyId:    0,
}
var zrSignChildWorkspace = idTypes.Workspace{
	Address: "workspace1xklrytgff7w32j52v34w36",
	Creator: "0xbe27E519A220D4336A5F9d08C7fa542175E30dA0",
	Owners:  []string{"0xbe27E519A220D4336A5F9d08C7fa542175E30dA0"},
	Alias:   "60",
}

func Test_msgServer_NewZrSignSignatureRequest_Hash_OrData(t *testing.T) {
	metadata, err := codectypes.NewAnyWithValue(&types.MetadataEthereum{})
	require.NoError(t, err)
	rawData, err := hex.DecodeString("805dc8973ad04df3f1195bda4e1a00894dbedce0288eb8c6a3bf8d228616489e")
	require.NoError(t, err)
	type args struct {
		keyring    *idTypes.Keyring
		workspaces []idTypes.Workspace
		key        *types.Key
		msg        *types.MsgNewZrSignSignatureRequest
	}
	tests := []struct {
		name     string
		args     args
		wantReq  *types.SignRequest
		wantResp *types.MsgNewZrSignSignatureRequestResponse
		wantErr  bool
	}{
		{
			name: "PASS: valid zrSign Sign data request",
			args: args{
				keyring:    &defaultKr,
				workspaces: []idTypes.Workspace{zrSignChildWorkspace, zrSignWorkspace},
				key:        &zrSignDefaultKey,
				msg: &types.MsgNewZrSignSignatureRequest{
					Creator:     "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
					Address:     "0xbe27E519A220D4336A5F9d08C7fa542175E30dA0",
					KeyType:     60,
					WalletIndex: 0,
					Data:        "805dc8973ad04df3f1195bda4e1a00894dbedce0288eb8c6a3bf8d228616489e",
					WalletType:  types.WalletType_WALLET_TYPE_EVM,
					NoBroadcast: true,
					Tx:          false,
					Metadata:    metadata,
				},
			},
			wantReq: &types.SignRequest{
				Id:             1,
				Creator:        "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
				KeyIds:         []uint64{1},
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				DataForSigning: [][]byte{rawData},
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				MpcBtl:         20,
			},
			wantResp: &types.MsgNewZrSignSignatureRequestResponse{
				ReqId: 1,
			},
		},
		{
			name: "FAIL: unauthorised zrSign request",
			args: args{
				keyring:    &defaultKr,
				workspaces: []idTypes.Workspace{zrSignChildWorkspace, zrSignWorkspace},
				key:        &zrSignDefaultKey,
				msg: &types.MsgNewZrSignSignatureRequest{
					Creator: "invalid",
				},
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

			idGenesis := idTypes.GenesisState{
				PortId:     idTypes.PortID,
				Keyrings:   []idTypes.Keyring{*tt.args.keyring},
				Workspaces: tt.args.workspaces,
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			tGenesis := types.GenesisState{
				PortId: types.PortID,
				Keys:   []types.Key{*tt.args.key},
				Params: types.Params{
					MpcKeyring:    "keyring1pfnq7r04rept47gaf5cpdew2",
					ZrSignAddress: "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
				},
			}
			treasury.InitGenesis(ctx, *tk, tGenesis)
			msgSer := keeper.NewMsgServerImpl(*tk, true)
			got, err := msgSer.NewZrSignSignatureRequest(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, got)

				gotSignatureReq, err := tk.SignRequestStore.Get(ctx, got.ReqId)
				require.NoError(t, err)
				require.Equal(t, tt.wantReq, &gotSignatureReq)
			}
		})
	}
}

func Test_msgServer_NewZrSignSignatureRequest_Transaction(t *testing.T) {
	metadata, err := codectypes.NewAnyWithValue(&types.MetadataEthereum{
		ChainId: 111,
	})
	require.NoError(t, err)
	rawData, err := hex.DecodeString("e679850fd8f4bd2382520894cf4324b9fd2fbc83b98a1483e854d31d3e45944c8203e880808080")
	require.NoError(t, err)
	type args struct {
		keyring    *idTypes.Keyring
		workspaces []idTypes.Workspace
		key        *types.Key
		msg        *types.MsgNewZrSignSignatureRequest
	}
	tests := []struct {
		name     string
		args     args
		wantReq  *types.SignTransactionRequest
		wantResp *types.MsgNewZrSignSignatureRequestResponse
		wantErr  bool
	}{
		{
			name: "PASS: valid signTransactionRequest",
			args: args{
				keyring:    &defaultKr,
				workspaces: []idTypes.Workspace{zrSignChildWorkspace, zrSignWorkspace},
				key:        &zrSignDefaultKey,
				msg: &types.MsgNewZrSignSignatureRequest{
					Creator:     "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
					Address:     "0xbe27E519A220D4336A5F9d08C7fa542175E30dA0",
					KeyType:     60,
					WalletIndex: 0,
					Data:        "e679850fd8f4bd2382520894cf4324b9fd2fbc83b98a1483e854d31d3e45944c8203e880808080",
					WalletType:  types.WalletType_WALLET_TYPE_EVM,
					NoBroadcast: true,
					Tx:          true,
					Metadata:    metadata,
				},
			},
			wantReq: &types.SignTransactionRequest{
				Id:                  1,
				Creator:             "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
				KeyIds:              []uint64{1},
				WalletType:          types.WalletType_WALLET_TYPE_EVM,
				UnsignedTransaction: rawData,
				SignRequestId:       1,
				NoBroadcast:         true,
			},
			wantResp: &types.MsgNewZrSignSignatureRequestResponse{
				ReqId: 1,
			},
		},
		{
			name: "FAIL: invalid signTransactionRequest",
			args: args{
				keyring:    &defaultKr,
				workspaces: []idTypes.Workspace{zrSignChildWorkspace, zrSignWorkspace},
				key:        &zrSignDefaultKey,
				msg: &types.MsgNewZrSignSignatureRequest{
					Creator:     "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
					Address:     "0xbe27E519A220D4336A5F9d08C7fa542175E30dA0",
					KeyType:     60,
					WalletIndex: 0,
					Data:        "invalid",
					WalletType:  types.WalletType_WALLET_TYPE_EVM,
					NoBroadcast: true,
					Tx:          true,
					Metadata:    metadata,
				},
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

			idGenesis := idTypes.GenesisState{
				PortId:     idTypes.PortID,
				Keyrings:   []idTypes.Keyring{*tt.args.keyring},
				Workspaces: tt.args.workspaces,
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			tGenesis := types.GenesisState{
				PortId: types.PortID,
				Keys:   []types.Key{*tt.args.key},
				Params: types.Params{
					MpcKeyring:    "keyring1pfnq7r04rept47gaf5cpdew2",
					ZrSignAddress: "zen1zpmqphp46nsn097ysltk4j5wmpjn9gd5gwyfnq",
				},
			}
			treasury.InitGenesis(ctx, *tk, tGenesis)
			msgSer := keeper.NewMsgServerImpl(*tk, true)
			got, err := msgSer.NewZrSignSignatureRequest(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, got)

				gotSignatureReq, err := tk.SignTransactionRequestStore.Get(ctx, got.ReqId)
				require.NoError(t, err)
				require.Equal(t, tt.wantReq, &gotSignatureReq)
			}
		})
	}
}
