package keeper_test

import (
	"encoding/base64"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	policymodule "github.com/Zenrock-Foundation/zrchain/v5/x/policy/module"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultKey = types.Key{
	Id:            1,
	WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
}

func Test_msgServer_NewSignTransactionRequest(t *testing.T) {
	unsignedTx1, err := hexutil.Decode("0xe9028501596f95f1825208945ff137d4b0fdcd49dca30c7cf57e578a026d278985e8d4a5100080808080")
	assert.NoError(t, err)

	unsignedCosmosTx, err := base64.RawStdEncoding.DecodeString("CogBCoUBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmUKKnplbjEzeTN0bTY4Z211OWtudGN4d3ZtdWU4MnA2YWthY25wdDJ2N250eRIqemVuMXM1Mm1lOXA3ZGFrOXU3OGR0MDZwMHNnZ3g0a3d0eGdwbDA4d3M0GgsKBW5yb2NrEgIxMBIGEgQQwJoM")
	assert.NoError(t, err)

	defaultKey := types.Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
		Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
	}
	keyWithPolicy := types.Key{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
		Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
		SignPolicyId:  1,
	}

	defaultKey2 := types.Key{
		Id:            2,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		KeyringAddr:   "notAkeyring",
		Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
	}

	var metadataAny *cdctypes.Any

	metadata := types.MetadataEthereum{
		ChainId: 11155111,
	}

	metadataAny, _ = cdctypes.NewAnyWithValue(&metadata)

	failMetadataAny, _ := cdctypes.NewAnyWithValue(nil)

	type args struct {
		keyring   *idTypes.Keyring
		workspace *idTypes.Workspace
		key       *types.Key
		msg       *types.MsgNewSignTransactionRequest
	}
	tests := []struct {
		name              string
		args              args
		wantSignTxRequest *types.SignTransactionRequest
		want              *types.MsgNewSignTransactionRequestResponse
		wantErr           bool
	}{
		{
			name: "PASS: valid signTransactionRequest",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &keyWithPolicy,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_EVM, unsignedTx1, metadataAny, 1000),
			},
			wantSignTxRequest: &types.SignTransactionRequest{
				Id:                  1,
				Creator:             "testOwner",
				KeyId:               1,
				WalletType:          types.WalletType_WALLET_TYPE_EVM,
				UnsignedTransaction: hexutil.MustDecode("0xe9028501596f95f1825208945ff137d4b0fdcd49dca30c7cf57e578a026d278985e8d4a5100080808080"), // hash of the unsigned tx
				SignRequestId:       1,
			},
			want: &types.MsgNewSignTransactionRequestResponse{
				Id:                 1,
				SignatureRequestId: 1,
			},
		},
		{
			name: "FAIL: key not found",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 5, types.WalletType_WALLET_TYPE_EVM, unsignedTx1, metadataAny, 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: workspace not found",
			args: args{
				keyring: &defaultKr,
				workspace: &idTypes.Workspace{
					Address: "otherWorkspace",
					Creator: "testOwner",
					Owners:  []string{"testOwner"},
				},
				key: &defaultKey,
				msg: types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_EVM, unsignedTx1, metadataAny, 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: keyring is inactive",
			args: args{
				keyring: &idTypes.Keyring{
					Address:     "keyring1pfnq7r04rept47gaf5cpdew2",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{"testCreator"},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    false,
				},
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_EVM, unsignedTx1, metadataAny, 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey2,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_EVM, unsignedTx1, metadataAny, 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: invalid wallet type",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_UNSPECIFIED, unsignedTx1, metadataAny, 1000),
			},
			wantErr: true,
		},
		{
			name: "PASS: valid signTransactionRequest for native wallet",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_NATIVE, unsignedCosmosTx, nil, 1000),
			},
			wantSignTxRequest: &types.SignTransactionRequest{
				Id:                  1,
				Creator:             "testOwner",
				KeyId:               1,
				WalletType:          types.WalletType_WALLET_TYPE_NATIVE,
				UnsignedTransaction: unsignedCosmosTx,
				SignRequestId:       1,
			},
			want: &types.MsgNewSignTransactionRequestResponse{
				Id:                 1,
				SignatureRequestId: 1,
			},
		},
		{
			name: "FAIL: wallet does not implement tx parser",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_EVM, unsignedTx1, nil, 1000),
			},
			wantErr: true,
		},
		{
			name: "FAIL: invalid metadata field",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignTransactionRequest("testOwner", 1, types.WalletType_WALLET_TYPE_EVM, unsignedTx1, failMetadataAny, 1000),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			pk := keepers.PolicyKeeper
			tk := keepers.TreasuryKeeper
			ctx := keepers.Ctx

			idGenesis := idTypes.GenesisState{
				PortId:     idTypes.PortID,
				Keyrings:   []idTypes.Keyring{*tt.args.keyring},
				Workspaces: []idTypes.Workspace{*tt.args.workspace},
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			tGenesis := types.GenesisState{
				PortId: types.PortID,
				Keys:   []types.Key{*tt.args.key},
			}
			treasury.InitGenesis(ctx, *tk, tGenesis)
			pGenesis := policytypes.GenesisState{
				Policies: []policytypes.Policy{
					policy1,
				},
			}
			policymodule.InitGenesis(ctx, *pk, pGenesis)

			msgSer := keeper.NewMsgServerImpl(*tk)

			got, err := msgSer.NewSignTransactionRequest(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)

				gotSignTxReq, err := tk.SignTransactionRequestStore.Get(ctx, got.Id)
				require.NoError(t, err)
				require.Equal(t, tt.wantSignTxRequest, &gotSignTxReq)

				act, err := pk.ActionStore.Get(ctx, 1)
				require.NoError(t, err)
				require.Equal(t, uint64(1000), act.Btl)
				require.Equal(t, tt.args.key.SignPolicyId, act.PolicyId)
			}
		})
	}
}
