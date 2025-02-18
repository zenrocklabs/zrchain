package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/stretchr/testify/require"
)

var defaultKeyRequest = types.KeyRequest{
	Id:            1,
	Creator:       "testCreator",
	WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
}

var keyRequestWithKeyring2 = types.KeyRequest{
	Id:            2,
	Creator:       "testCreator",
	WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
	KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
}

var keyRequestWithKeyring3 = types.KeyRequest{
	Id:            3,
	Creator:       "testCreator",
	WorkspaceAddr: "workspace1xklrytgff7w32j52v34w36",
	KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
}

var partialKeyRequest = types.KeyRequest{
	Id:                     1,
	Creator:                "testCreator",
	WorkspaceAddr:          "workspace14a2hpadpsy9h4auve2z8lw",
	KeyringAddr:            "keyring1pfnq7r04rept47gaf5cpdew2",
	KeyType:                types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	Status:                 types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL,
	KeyringPartySignatures: [][]byte{[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur")},
	PublicKey:              defaultECDSAKey.PublicKey,
}

var defaultKeyReqResponse = types.KeyReqResponse{
	Id:            1,
	Creator:       "testCreator",
	WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
	Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING.String(),
}

var defaultKeyReqResponse2 = types.KeyReqResponse{
	Id:            2,
	Creator:       "testCreator",
	WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
	KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
	Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING.String(),
}

var defaultKeyReqResponse3 = types.KeyReqResponse{
	Id:            3,
	Creator:       "testCreator",
	WorkspaceAddr: "workspace1xklrytgff7w32j52v34w36",
	KeyringAddr:   "keyring1k6vc6vhp6e6l3rxalue9v4ux",
	KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
	Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED.String(),
}

func Test_msgServer_FulfilKeyRequest(t *testing.T) {
	// both of these keys are invalid as they are too long
	invalidECDSAPubKey := []byte{154, 135, 176, 26, 117, 104, 94, 9, 73, 68, 162, 139, 9, 231, 47, 249, 137, 156, 60, 87, 66, 163}
	invalidEdDSAPubkey := []byte{1, 243, 178, 23, 221, 136, 81, 23, 248, 229, 31, 154, 135, 176, 26, 117, 104, 94, 9, 73, 68, 162, 139, 9, 231, 47, 249, 137, 156, 60, 87, 66, 163}

	type args struct {
		keyring    *idTypes.Keyring
		workspace  *idTypes.Workspace
		keyRequest *types.KeyRequest
		msg        *types.MsgFulfilKeyRequest
	}
	tests := []struct {
		name           string
		args           args
		wantKeyRequest *types.KeyRequest
		want           *types.MsgFulfilKeyRequestResponse
		wantErr        bool
	}{
		{
			name: "PASS: return a new ecdsa key",
			args: args{
				keyring:    &defaultKr,
				workspace:  &defaultWs,
				keyRequest: &defaultKeyRequest,
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultECDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			wantKeyRequest: &types.KeyRequest{
				Id:                     1,
				Creator:                "testCreator",
				WorkspaceAddr:          "workspace14a2hpadpsy9h4auve2z8lw",
				KeyringAddr:            "keyring1pfnq7r04rept47gaf5cpdew2",
				KeyType:                types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				Status:                 types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL,
				KeyringPartySignatures: [][]byte{[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur")},
				PublicKey:              defaultECDSAKey.PublicKey,
			},
			want: &types.MsgFulfilKeyRequestResponse{},
		},
		{
			name: "PASS: reject the request",
			args: args{
				keyring:    &defaultKr,
				workspace:  &defaultWs,
				keyRequest: &defaultKeyRequest,
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED,
					types.NewMsgFulfilKeyRequestReject("test"),
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			wantKeyRequest: &types.KeyRequest{
				Id:            1,
				Creator:       "testCreator",
				WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
				KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
				KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED,
				RejectReason:  "test",
			},
			want: &types.MsgFulfilKeyRequestResponse{},
		},
		{
			name: "PASS: return a new eddsa key",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			wantKeyRequest: &types.KeyRequest{
				Id:                     1,
				Creator:                "testCreator",
				WorkspaceAddr:          "workspace14a2hpadpsy9h4auve2z8lw",
				KeyringAddr:            "keyring1pfnq7r04rept47gaf5cpdew2",
				KeyType:                types.KeyType_KEY_TYPE_EDDSA_ED25519,
				Status:                 types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL,
				KeyringPartySignatures: [][]byte{[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur")},
				PublicKey:              defaultEdDSAKey.PublicKey,
			},
			want: &types.MsgFulfilKeyRequestResponse{},
		},
		{
			name: "PASS: reject eddsa key request",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED,
					&types.MsgFulfilKeyRequest_RejectReason{RejectReason: "test"},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			wantKeyRequest: &types.KeyRequest{
				Id:            1,
				Creator:       "testCreator",
				WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
				KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
				KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
				Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED,
				RejectReason:  "test",
			},
			want: &types.MsgFulfilKeyRequestResponse{},
		},
		{
			name: "FAIL: request not found",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					999,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: keyring not found",
			args: args{
				keyring: &idTypes.Keyring{
					Address:     "notAKeyring",
					Creator:     "testCreator",
					Description: "testDescription",
					Admins:      []string{"testCreator"},
					Parties:     []string{"testCreator"},
					KeyReqFee:   0,
					SigReqFee:   0,
					IsActive:    true,
				},
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: keyring inactive",
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
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: creator is no keyring party",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"noKeyringParty",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: keyRequest status is not pending",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: keyRequest status is unspecified",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_UNSPECIFIED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: ecdsa pubkey is invalid",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: invalidECDSAPubKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: eddsa pubkey is too long",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: invalidEdDSAPubkey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: invalid key type",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_UNSPECIFIED,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: no mpc party signature",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte{},
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: invalid mpc party signature",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				keyRequest: &types.KeyRequest{
					Id:            1,
					Creator:       "testCreator",
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					KeyType:       types.KeyType_KEY_TYPE_EDDSA_ED25519,
					Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
				},
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultEdDSAKey.PublicKey}},
					[]byte("InvalidLengthSignature"), // should be 64 bytes long, will fail
				),
			},
			want:    &types.MsgFulfilKeyRequestResponse{},
			wantErr: true,
		},
		{
			name: "PASS: fulfil key with 2 parties",
			args: args{
				keyring:    &defaultKr,
				workspace:  &defaultWs,
				keyRequest: &partialKeyRequest,
				msg: types.NewMsgFulfilKeyRequest(
					"testCreator",
					1,
					types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
					&types.MsgFulfilKeyRequest_Key{Key: &types.MsgNewKey{PublicKey: defaultECDSAKey.PublicKey}},
					[]byte("0000000000000000000000000000000000000000000000000SecondSignature"),
				),
			},
			wantKeyRequest: &types.KeyRequest{
				Id:            1,
				Creator:       "testCreator",
				WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
				KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
				KeyType:       types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				Status:        types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED,
				KeyringPartySignatures: [][]byte{
					[]byte("TestSignatureTestSignatureTestSignatureTestSignatureTestSignatur"),
					[]byte("0000000000000000000000000000000000000000000000000SecondSignature"),
				},
				PublicKey: defaultECDSAKey.PublicKey,
			},
			want: &types.MsgFulfilKeyRequestResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			ik := keepers.IdentityKeeper
			ctx := keepers.Ctx

			genesis := idTypes.GenesisState{
				PortId:     idTypes.PortID,
				Keyrings:   []idTypes.Keyring{*tt.args.keyring},
				Workspaces: []idTypes.Workspace{*tt.args.workspace},
			}
			identity.InitGenesis(ctx, *ik, genesis)

			tk := keepers.TreasuryKeeper
			tGenesis := types.GenesisState{
				PortId:      types.PortID,
				KeyRequests: []types.KeyRequest{*tt.args.keyRequest},
			}
			treasury.InitGenesis(ctx, *tk, tGenesis)

			msgSer := keeper.NewMsgServerImpl(*tk)

			got, err := msgSer.FulfilKeyRequest(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)

				gotKeyRequest, err := tk.KeyRequestStore.Get(ctx, tt.args.keyRequest.Id)
				require.NoError(t, err)
				require.Equal(t, tt.wantKeyRequest, &gotKeyRequest)
			}
		})
	}
}
