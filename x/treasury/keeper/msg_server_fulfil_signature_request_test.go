package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/stretchr/testify/require"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

var defaultKr = idTypes.Keyring{
	Address:        "keyring1pfnq7r04rept47gaf5cpdew2",
	Creator:        "testCreator",
	Description:    "testDescription",
	Admins:         []string{"testCreator"},
	Parties:        []string{"testCreator", "testCreator2"},
	PartyThreshold: 2,
	KeyReqFee:      0,
	SigReqFee:      0,
	IsActive:       true,
}

var defaultWs = idTypes.Workspace{
	Address: "workspace14a2hpadpsy9h4auve2z8lw",
	Creator: "testOwner",
	Owners:  []string{"testOwner"},
}

var defaultECDSAKey = types.Key{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
}

// var defaultBTCKey = types.Key{
// 	Id:            1,
// 	WorkspaceAddr: "testWorkspace",
// 	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
// 	Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
// 	PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
// }

var defaultBitcoinKey = types.Key{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
	PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
}

var defaultECDSAKeyResponse = types.KeyResponse{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1.String(),
	PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
}

var defaultBitcoinKeyResponse = types.KeyResponse{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1.String(),
	PublicKey:     []byte{0x03, 0xca, 0x27, 0xea, 0x7b, 0x06, 0x41, 0x49, 0x7b, 0x19, 0xa7, 0x23, 0xe3, 0xb9, 0x25, 0x90, 0x80, 0x1c, 0x7e, 0x79, 0xb1, 0x14, 0x25, 0x3f, 0xc1, 0xe9, 0x9d, 0xf1, 0xfd, 0x97, 0x30, 0x52, 0x6b},
}

var secondECDSAKey = types.Key{
	Id:            2,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "differentKeyring",
	Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
	PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
}

var defaultEdDSAKey = types.Key{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_EDDSA_ED25519,
	PublicKey:     []byte{0x29, 0xdb, 0xaa, 0xd6, 0x75, 0x87, 0x8b, 0x2f, 0x16, 0xcb, 0x14, 0x82, 0x7a, 0x4c, 0x1a, 0x41, 0xbb, 0xd6, 0x3c, 0x3b, 0x60, 0x1b, 0xdc, 0x5e, 0x27, 0x7d, 0x00, 0xb1, 0x20, 0xff, 0xee, 0xce},
}

var defaultEdDSAKeyResponse = types.KeyResponse{
	Id:            1,
	WorkspaceAddr: "testWorkspace",
	KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
	Type:          types.KeyType_KEY_TYPE_EDDSA_ED25519.String(),
	PublicKey:     []byte{0x29, 0xdb, 0xaa, 0xd6, 0x75, 0x87, 0x8b, 0x2f, 0x16, 0xcb, 0x14, 0x82, 0x7a, 0x4c, 0x1a, 0x41, 0xbb, 0xd6, 0x3c, 0x3b, 0x60, 0x1b, 0xdc, 0x5e, 0x27, 0x7d, 0x00, 0xb1, 0x20, 0xff, 0xee, 0xce},
}

func Test_msgServer_FulfilSignatureRequest(t *testing.T) {
	var sigRequestPayload = [][]byte{{0x77, 0x8f, 0x57, 0x2f, 0x33, 0xaf, 0xab, 0x83, 0x13, 0x65, 0xd5, 0x2e, 0x56, 0x3a, 0x0d, 0xdd, 0x77, 0x8f, 0x57, 0x2f, 0x33, 0xaf, 0xab, 0x83, 0x13, 0x65, 0xd5, 0x2e, 0x56, 0x3a, 0x0d, 0xdd}}
	var sigPayloadECDSA = []*types.SignedDataWithID{{SignRequestId: 1, SignedData: []byte{0xd8, 0x35, 0x01, 0x30, 0xb6, 0x8f, 0xdc, 0x95, 0x12, 0xc3, 0xfa, 0x28, 0xd4, 0xa7, 0xb9, 0xc4, 0xab, 0x50, 0xe0, 0x8e, 0xf1, 0xe2, 0xd8, 0xe4, 0x6c, 0x3e, 0xe7, 0x2f, 0xf7, 0x8e, 0xec, 0xbd, 0x1e, 0x7e, 0xd7, 0x30, 0x41, 0xde, 0x77, 0x48, 0xf8, 0x5d, 0xb7, 0x24, 0x39, 0x5e, 0x76, 0x24, 0xba, 0xaa, 0xa6, 0x20, 0x70, 0x0e, 0x76, 0xae, 0x0e, 0x53, 0x83, 0xbf, 0x8d, 0x2b, 0xa2, 0xdb}}}
	var sigPayloadEdDSA = []*types.SignedDataWithID{{SignRequestId: 1, SignedData: []byte{0xa1, 0xd5, 0x35, 0x71, 0x0d, 0x81, 0x11, 0x2c, 0x06, 0x8a, 0xad, 0xe8, 0x0d, 0x5d, 0x2f, 0xf1, 0x2a, 0xdd, 0xcc, 0xf8, 0x6a, 0x6b, 0x38, 0xd9, 0x2d, 0x80, 0xcb, 0x1d, 0x8e, 0xa6, 0x4c, 0xd3, 0x72, 0xea, 0x75, 0x19, 0x34, 0x0e, 0xc5, 0x49, 0xfd, 0xae, 0xb2, 0x28, 0x11, 0x8e, 0xa4, 0x42, 0xfc, 0xea, 0x64, 0x5b, 0x65, 0x16, 0x59, 0xdc, 0x4b, 0x0e, 0xf9, 0x39, 0x24, 0x18, 0xe7, 0x0d}}}
	// Different valid signature for testing mismatched signatures
	var differentSigPayloadECDSA = []byte{0xd9, 0x35, 0x01, 0x30, 0xb6, 0x8f, 0xdc, 0x95, 0x12, 0xc3, 0xfa, 0x28, 0xd4, 0xa7, 0xb9, 0xc4, 0xab, 0x50, 0xe0, 0x8e, 0xf1, 0xe2, 0xd8, 0xe4, 0x6c, 0x3e, 0xe7, 0x2f, 0xf7, 0x8e, 0xec, 0xbd, 0x1e, 0x7e, 0xd7, 0x30, 0x41, 0xde, 0x77, 0x48, 0xf8, 0x5d, 0xb7, 0x24, 0x39, 0x5e, 0x76, 0x24, 0xba, 0xaa, 0xa6, 0x20, 0x70, 0x0e, 0x76, 0xae, 0x0e, 0x53, 0x83, 0xbf, 0x8d, 0x2b, 0xa2, 0xdb}

	var defaultSigRequest = types.SignRequest{
		Id:             1,
		Creator:        "testCreator",
		KeyIds:         []uint64{1},
		DataForSigning: sigRequestPayload,
		Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL,
		KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		KeyringPartySigs: []*types.PartySignature{
			{
				Creator:   "testCreator",
				Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
			},
		},
	}

	var defaultResponseECDSA = types.MsgFulfilSignatureRequest{
		Creator:               "testRequestCreator",
		RequestId:             1,
		Status:                types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
		SignedData:            sigPayloadECDSA[0].SignedData,
		KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
	}

	type args struct {
		key *types.Key
		req *types.SignRequest
		msg *types.MsgFulfilSignatureRequest
	}
	tests := []struct {
		name       string
		args       args
		want       *types.MsgFulfilSignatureRequestResponse
		wantSigReq *types.SignRequest
		wantErr    bool
	}{
		{
			name: "PASS: return signature request - ECDSA",
			args: args{
				key: &defaultECDSAKey,
				req: &defaultSigRequest,
				msg: &defaultResponseECDSA,
			},
			want: &types.MsgFulfilSignatureRequestResponse{},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				SignedData:     sigPayloadECDSA,
				KeyringPartySigs: []*types.PartySignature{
					{
						Creator:   "testCreator",
						Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
					{
						Creator:   "testRequestCreator",
						Signature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PASS: return signature request - EDDSA",
			args: args{
				key: &defaultEdDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL,
					KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
					KeyringPartySigs: []*types.PartySignature{
						{
							Creator:   "testCreator",
							Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
						},
					},
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:               "testRequestCreator",
					RequestId:             1,
					Status:                types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:            sigPayloadEdDSA[0].SignedData,
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			want: &types.MsgFulfilSignatureRequestResponse{},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
				KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				SignedData:     sigPayloadEdDSA,
				KeyringPartySigs: []*types.PartySignature{
					{
						Creator:   "testCreator",
						Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
					{
						Creator:   "testRequestCreator",
						Signature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PASS: Reject ECDSA Signature Request",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:      "testCreator",
					RequestId:    1,
					Status:       types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
					RejectReason: "test",
				},
			},
			want: &types.MsgFulfilSignatureRequestResponse{},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				RejectReason:   "test",
			},
			wantErr: false,
		},
		{
			name: "PASS: test unsigned transaction",
			args: args{
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					PublicKey:     []byte{0x02, 0xc5, 0x17, 0x97, 0x42, 0x4c, 0x35, 0x22, 0xe6, 0xed, 0x70, 0x44, 0xb8, 0xf5, 0xa2, 0xc9, 0x1d, 0x59, 0xb6, 0x19, 0xc1, 0x2c, 0x22, 0x3f, 0x41, 0xb1, 0xdb, 0xf9, 0x6b, 0x4e, 0x7e, 0x1c, 0x3e},
				},
				req: &types.SignRequest{
					Id:             1,
					Creator:        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
					KeyIds:         []uint64{1},
					DataForSigning: [][]byte{{0xe8, 0xcc, 0x5a, 0x88, 0x37, 0xc8, 0x4e, 0x8f, 0x21, 0xab, 0xcc, 0x1e, 0xa6, 0x66, 0x24, 0xa9, 0x3e, 0x01, 0x12, 0xbe, 0xc6, 0xf3, 0xdb, 0x3a, 0xef, 0x79, 0x64, 0x9d, 0x6c, 0xac, 0xc0, 0xe0}},
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					KeyringPartySigs: []*types.PartySignature{
						{
							Creator:   "testCreator",
							Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
						},
					},
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:               "testCreator2",
					RequestId:             1,
					Status:                types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:            []byte{0x2a, 0x75, 0x7d, 0x0b, 0xca, 0x85, 0xea, 0x2d, 0xba, 0x26, 0x36, 0x66, 0xfb, 0x84, 0xa0, 0x35, 0xe4, 0xaa, 0x3f, 0x0e, 0xe8, 0x30, 0xac, 0x48, 0x59, 0x4b, 0xb8, 0xae, 0x60, 0x77, 0x5c, 0x65, 0x6a, 0x59, 0xfd, 0x38, 0xfe, 0x8c, 0x28, 0x80, 0xc3, 0x0a, 0x64, 0xc9, 0x98, 0x29, 0x6e, 0xc7, 0xaf, 0xde, 0xfa, 0xc3, 0x50, 0x8a, 0x41, 0xab, 0xbe, 0x4c, 0xfb, 0xe7, 0xbd, 0x7f, 0xb3, 0x2a, 0x1b},
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			want: &types.MsgFulfilSignatureRequestResponse{},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
				KeyIds:         []uint64{1},
				DataForSigning: [][]byte{{0xe8, 0xcc, 0x5a, 0x88, 0x37, 0xc8, 0x4e, 0x8f, 0x21, 0xab, 0xcc, 0x1e, 0xa6, 0x66, 0x24, 0xa9, 0x3e, 0x01, 0x12, 0xbe, 0xc6, 0xf3, 0xdb, 0x3a, 0xef, 0x79, 0x64, 0x9d, 0x6c, 0xac, 0xc0, 0xe0}},
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				SignedData:     []*types.SignedDataWithID{{SignRequestId: 1, SignedData: []byte{0x2a, 0x75, 0x7d, 0x0b, 0xca, 0x85, 0xea, 0x2d, 0xba, 0x26, 0x36, 0x66, 0xfb, 0x84, 0xa0, 0x35, 0xe4, 0xaa, 0x3f, 0x0e, 0xe8, 0x30, 0xac, 0x48, 0x59, 0x4b, 0xb8, 0xae, 0x60, 0x77, 0x5c, 0x65, 0x6a, 0x59, 0xfd, 0x38, 0xfe, 0x8c, 0x28, 0x80, 0xc3, 0x0a, 0x64, 0xc9, 0x98, 0x29, 0x6e, 0xc7, 0xaf, 0xde, 0xfa, 0xc3, 0x50, 0x8a, 0x41, 0xab, 0xbe, 0x4c, 0xfb, 0xe7, 0xbd, 0x7f, 0xb3, 0x2a, 0x00}}},
				KeyringPartySigs: []*types.PartySignature{
					{
						Creator:   "testCreator",
						Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
					{
						Creator:   "testCreator2",
						Signature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "FAIL: Empty Status field in rejection",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:      "testCreator",
					RequestId:    1,
					Status:       types.SignRequestStatus_SIGN_REQUEST_STATUS_UNSPECIFIED,
					RejectReason: "test",
				},
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: invalid signature, want eddsa got ecdsa",
			args: args{
				key: &defaultEdDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				},
				msg: &defaultResponseECDSA,
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				RejectReason:   "signature verification failed: verifySignature- invalid signature d8350130b68fdc9512c3fa28d4a7b9c4ab50e08ef1e2d8e46c3ee72ff78eecbd1e7ed73041de7748f85db724395e7624baaaa620700e76ae0e5383bf8d2ba2db from keyring keyring1pfnq7r04rept47gaf5cpdew2",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: signature request status is already fulfilled",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:    "testCreator",
					RequestId:  1,
					Status:     types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData: sigPayloadEdDSA[0].SignedData,
				},
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: signature request not found",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:    "testCreator",
					RequestId:  9999,
					Status:     types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData: sigPayloadEdDSA[0].SignedData,
				},
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: true,
		},
		{
			name: "FAIL: invalid ecdsa signature length",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:   "testCreator",
					RequestId: 1,
					Status:    types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:
					// added 11 on top to max out on length
					[]byte{11, 173, 224, 103, 159, 55, 251, 255, 212, 50, 145, 235, 181, 19, 14, 120, 168, 208, 151, 33, 204, 161, 79, 118, 229, 75, 22, 185, 234, 115, 125, 170, 101, 55, 197, 218, 94, 172, 32, 139, 21, 141, 104, 163, 109, 45, 47, 80, 110, 39, 5, 156, 88, 31, 82, 123, 246, 67, 21, 199, 126, 75, 222, 65, 115, 1},
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				RejectReason:   "signature verification failed: verifySignature- invalid ecdsa signature 0bade0679f37fbffd43291ebb5130e78a8d09721cca14f76e54b16b9ea737daa6537c5da5eac208b158d68a36d2d2f506e27059c581f527bf64315c77e4bde417301 of length 66",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: invalid key type",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_UNSPECIFIED,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:               "testCreator",
					RequestId:             1,
					Status:                types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:            []byte{173, 224, 103, 159, 55, 251, 255, 212, 50, 145, 235, 181, 19, 14, 120, 168, 208, 151, 33, 204, 161, 79, 118, 229, 75, 22, 185, 234, 115, 125, 170, 101, 55, 197, 218, 94, 172, 32, 139, 21, 141, 104, 163, 109, 45, 47, 80, 110, 39, 5, 156, 88, 31, 82, 123, 246, 67, 21, 199, 126, 75, 222, 65, 115, 1},
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_UNSPECIFIED,
				RejectReason:   "signature verification failed: verifySignature- invalid key type: KEY_TYPE_UNSPECIFIED",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: invalid ecdsa signature",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:   "testCreator",
					RequestId: 1,
					Status:    types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:
					// invalid signature data
					[]byte{175, 224, 103, 159, 55, 251, 255, 212, 50, 145, 235, 181, 19, 14, 120, 168, 208, 151, 33, 204, 161, 79, 118, 229, 75, 22, 185, 234, 115, 125, 170, 101, 55, 197, 218, 94, 172, 32, 139, 21, 141, 104, 163, 109, 45, 47, 80, 110, 39, 5, 156, 88, 31, 82, 123, 246, 67, 21, 199, 126, 75, 222, 65, 115, 1},
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				RejectReason:   "signature verification failed: verifySignature- invalid signature afe0679f37fbffd43291ebb5130e78a8d09721cca14f76e54b16b9ea737daa6537c5da5eac208b158d68a36d2d2f506e27059c581f527bf64315c77e4bde417301 from keyring keyring1pfnq7r04rept47gaf5cpdew2",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: invalid eddsa signature",
			args: args{
				key: &defaultEdDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:   "testCreator",
					RequestId: 1,
					Status:    types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:
					// invalid signature data
					[]byte{225, 103, 159, 55, 251, 255, 212, 50, 145, 235, 181, 19, 14, 120, 168, 208, 151, 33, 204, 161, 79, 118, 229, 75, 22, 185, 234, 115, 125, 170, 101, 55, 197, 218, 94, 172, 32, 139, 21, 141, 104, 163, 109, 45, 47, 80, 110, 39, 5, 156, 88, 31, 82, 123, 246, 67, 21, 199, 126, 75, 222, 65, 115, 1},
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				RejectReason:   "signature verification failed: verifySignature- invalid signature e1679f37fbffd43291ebb5130e78a8d09721cca14f76e54b16b9ea737daa6537c5da5eac208b158d68a36d2d2f506e27059c581f527bf64315c77e4bde417301 from keyring keyring1pfnq7r04rept47gaf5cpdew2",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: invalid eddsa signature length",
			args: args{
				key: &defaultEdDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:   "testCreator",
					RequestId: 1,
					Status:    types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:
					// invalid length for eddsa sig
					[]byte{225, 103, 255, 212, 50, 145, 235, 181, 19, 14, 120, 168, 208, 151, 33, 204, 161, 79, 118, 229, 75, 22, 185, 234, 115, 125, 170, 101, 55, 197, 218, 94, 172, 32, 139, 21, 141, 104, 163, 109, 45, 47, 80, 110, 39, 5, 156, 88, 31, 82, 123, 246, 67, 21, 199, 126, 75, 222, 65, 115, 1},
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				RejectReason:   "signature verification failed: verifySignature- invalid eddsa signature e167ffd43291ebb5130e78a8d09721cca14f76e54b16b9ea737daa6537c5da5eac208b158d68a36d2d2f506e27059c581f527bf64315c77e4bde417301 of length 61",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: no mpc party signature",
			args: args{
				key: &defaultEdDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:    "testCreator",
					RequestId:  1,
					Status:     types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData: []byte{173, 224, 103, 159, 55, 251, 255, 212, 50, 145, 235, 181, 19, 14, 120, 168, 208, 151, 33, 204, 161, 79, 118, 229, 75, 22, 185, 234, 115, 125, 170, 101, 55, 197, 218, 94, 172, 32, 139, 21, 141, 104, 163, 109, 45, 47, 80, 110, 39, 5, 156, 88, 31, 82, 123, 246, 67, 21, 199, 126, 75, 222, 65, 115, 1},
					// KeyringPartySignature is missing
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				RejectReason:   "invalid length of mpc party signature",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: invalid mpc party signature",
			args: args{
				key: &defaultEdDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
					KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:               "testCreator",
					RequestId:             1,
					Status:                types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:            []byte{173, 224, 103, 159, 55, 251, 255, 212, 50, 145, 235, 181, 19, 14, 120, 168, 208, 151, 33, 204, 161, 79, 118, 229, 75, 22, 185, 234, 115, 125, 170, 101, 55, 197, 218, 94, 172, 32, 139, 21, 141, 104, 163, 109, 45, 47, 80, 110, 39, 5, 156, 88, 31, 82, 123, 246, 67, 21, 199, 126, 75, 222, 65, 115, 1},
					KeyringPartySignature: []byte("InvalidLengthSignature"), // should be 64 bytes long, will fail
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				RejectReason:   "invalid length of mpc party signature",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: signed data mismatch on partial request",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					SignedData:     sigPayloadECDSA,
					KeyringPartySigs: []*types.PartySignature{
						{
							Creator:   "testCreator",
							Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
						},
					},
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:               "testCreator",
					RequestId:             1,
					Status:                types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:            differentSigPayloadECDSA,
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				SignedData:     sigPayloadECDSA,
				KeyringPartySigs: []*types.PartySignature{
					{
						Creator:   "testCreator",
						Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
				},
				RejectReason: "party testCreator already sent a fulfilment",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
		},
		{
			name: "FAIL: party testCreator already sent a fulfilment",
			args: args{
				key: &defaultECDSAKey,
				req: &types.SignRequest{
					Id:             1,
					Creator:        "testCreator",
					KeyIds:         []uint64{1},
					DataForSigning: sigRequestPayload,
					Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL,
					KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					SignedData:     sigPayloadECDSA,
					KeyringPartySigs: []*types.PartySignature{
						{
							Creator:   "testCreator",
							Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
						},
					},
				},
				msg: &types.MsgFulfilSignatureRequest{
					Creator:               "testCreator",
					RequestId:             1,
					Status:                types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED,
					SignedData:            differentSigPayloadECDSA,
					KeyringPartySignature: []byte("1TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
				},
			},
			wantSigReq: &types.SignRequest{
				Id:             1,
				Creator:        "testCreator",
				KeyIds:         []uint64{1},
				DataForSigning: sigRequestPayload,
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				SignedData:     sigPayloadECDSA,
				KeyringPartySigs: []*types.PartySignature{
					{
						Creator:   "testCreator",
						Signature: []byte("0TestSignatureTestSignatureTestSignatureTestSignatureTestSignatu"),
					},
				},
				RejectReason: "party testCreator already sent a fulfilment",
			},
			want:    &types.MsgFulfilSignatureRequestResponse{},
			wantErr: false,
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
				PortId:     idTypes.PortID,
				Keyrings:   []idTypes.Keyring{defaultKr},
				Workspaces: []idTypes.Workspace{defaultWs},
			}
			identity.InitGenesis(ctx, *ik, idGenesis)

			trGenesis := types.GenesisState{
				PortId:       types.PortID,
				Keys:         []types.Key{*tt.args.key},
				SignRequests: []types.SignRequest{*tt.args.req},
			}
			treasury.InitGenesis(ctx, *tk, trGenesis)

			got, err := msgSer.FulfilSignatureRequest(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)

				gotSigReq, err := tk.SignRequestStore.Get(ctx, tt.args.msg.RequestId)
				require.NoError(t, err)
				require.Equal(t, tt.wantSigReq, &gotSigReq)
			}
		})
	}
}
