package keeper_test

import (
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v5/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v5/x/identity/module"
	idTypes "github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	policymodule "github.com/Zenrock-Foundation/zrchain/v5/x/policy/module"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	treasurymodule "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/module"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_NewSignatureRequest(t *testing.T) {
	type args struct {
		keyring   *idTypes.Keyring
		workspace *idTypes.Workspace
		key       *types.Key
		msg       *types.MsgNewSignatureRequest
	}
	tests := []struct {
		name            string
		args            args
		wantSignRequest *types.SignRequest
		want            *types.MsgNewSignatureRequestResponse
		wantErr         bool
	}{
		{
			name: "PASS: valid signature request",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
					SignPolicyId:  1,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 0),
			},
			wantSignRequest: &types.SignRequest{
				Id:             1,
				Creator:        "testOwner",
				KeyIds:         []uint64{1},
				DataForSigning: [][]byte{{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170}},
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "PASS: valid signature request with btl",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
					SignPolicyId:  1,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 10),
			},
			wantSignRequest: &types.SignRequest{
				Id:             1,
				Creator:        "testOwner",
				KeyIds:         []uint64{1},
				DataForSigning: [][]byte{{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170}},
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:        types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				MpcBtl:         10,
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "PASS: valid eddsa signature request",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_EDDSA_ED25519,
					PublicKey:     defaultEdDSAKey.PublicKey,
					SignPolicyId:  0,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000", 1000, 10),
			},
			wantSignRequest: &types.SignRequest{
				Id:             1,
				Creator:        "testOwner",
				KeyIds:         []uint64{1},
				DataForSigning: [][]byte{{2, 0, 1, 4, 106, 200, 242, 15, 199, 110, 53, 217, 9, 248, 87, 170, 1, 160, 131, 175, 227, 23, 196, 241, 229, 7, 103, 11, 49, 129, 233, 182, 168, 226, 195, 66, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 250, 6, 197, 13, 192, 115, 210, 0, 122, 254, 173, 146, 207, 164, 137, 238, 126, 190, 159, 207, 68, 37, 53, 21, 179, 48, 16, 64, 147, 141, 186, 10, 6, 167, 213, 23, 25, 44, 86, 142, 224, 138, 132, 95, 115, 210, 151, 136, 207, 3, 92, 49, 69, 178, 26, 179, 68, 216, 6, 46, 169, 64, 0, 0, 28, 204, 29, 159, 246, 25, 216, 183, 96, 247, 77, 103, 207, 0, 11, 71, 247, 29, 50, 100, 17, 181, 2, 186, 253, 95, 91, 237, 122, 60, 132, 244, 2, 1, 3, 1, 3, 1, 4, 4, 0, 0, 0, 1, 2, 0, 2, 12, 2, 0, 0, 0, 159, 61, 53, 191, 1, 0, 0, 0}},
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:        types.KeyType_KEY_TYPE_EDDSA_ED25519,
				MpcBtl:         10,
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "PASS: multiple eddsa signature payloads",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_EDDSA_ED25519,
					PublicKey:     defaultEdDSAKey.PublicKey,
					SignPolicyId:  0,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1, 1}, "020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000,020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000", 1000, 10),
			},
			wantSignRequest: &types.SignRequest{
				Id:      1,
				Creator: "testOwner",
				KeyIds:  []uint64{1, 1},
				DataForSigning: [][]byte{
					{2, 0, 1, 4, 106, 200, 242, 15, 199, 110, 53, 217, 9, 248, 87, 170, 1, 160, 131, 175, 227, 23, 196, 241, 229, 7, 103, 11, 49, 129, 233, 182, 168, 226, 195, 66, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 250, 6, 197, 13, 192, 115, 210, 0, 122, 254, 173, 146, 207, 164, 137, 238, 126, 190, 159, 207, 68, 37, 53, 21, 179, 48, 16, 64, 147, 141, 186, 10, 6, 167, 213, 23, 25, 44, 86, 142, 224, 138, 132, 95, 115, 210, 151, 136, 207, 3, 92, 49, 69, 178, 26, 179, 68, 216, 6, 46, 169, 64, 0, 0, 28, 204, 29, 159, 246, 25, 216, 183, 96, 247, 77, 103, 207, 0, 11, 71, 247, 29, 50, 100, 17, 181, 2, 186, 253, 95, 91, 237, 122, 60, 132, 244, 2, 1, 3, 1, 3, 1, 4, 4, 0, 0, 0, 1, 2, 0, 2, 12, 2, 0, 0, 0, 159, 61, 53, 191, 1, 0, 0, 0},
					{2, 0, 1, 4, 106, 200, 242, 15, 199, 110, 53, 217, 9, 248, 87, 170, 1, 160, 131, 175, 227, 23, 196, 241, 229, 7, 103, 11, 49, 129, 233, 182, 168, 226, 195, 66, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 250, 6, 197, 13, 192, 115, 210, 0, 122, 254, 173, 146, 207, 164, 137, 238, 126, 190, 159, 207, 68, 37, 53, 21, 179, 48, 16, 64, 147, 141, 186, 10, 6, 167, 213, 23, 25, 44, 86, 142, 224, 138, 132, 95, 115, 210, 151, 136, 207, 3, 92, 49, 69, 178, 26, 179, 68, 216, 6, 46, 169, 64, 0, 0, 28, 204, 29, 159, 246, 25, 216, 183, 96, 247, 77, 103, 207, 0, 11, 71, 247, 29, 50, 100, 17, 181, 2, 186, 253, 95, 91, 237, 122, 60, 132, 244, 2, 1, 3, 1, 3, 1, 4, 4, 0, 0, 0, 1, 2, 0, 2, 12, 2, 0, 0, 0, 159, 61, 53, 191, 1, 0, 0, 0},
				},
				Status:      types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:     types.KeyType_KEY_TYPE_EDDSA_ED25519,
				MpcBtl:      10,
				ChildReqIds: []uint64{2, 3},
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "PASS: multiple ecdsa signature payloads",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
					PublicKey:     defaultECDSAKey.PublicKey,
					SignPolicyId:  0,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1, 1}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa,778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 10),
			},
			wantSignRequest: &types.SignRequest{
				Id:      1,
				Creator: "testOwner",
				KeyIds:  []uint64{1, 1},
				DataForSigning: [][]byte{
					{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170},
					{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170},
				},
				Status:      types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:     types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
				MpcBtl:      10,
				ChildReqIds: []uint64{2, 3},
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "PASS: multiple btc ecdsa signature payloads",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
					PublicKey:     defaultBitcoinKey.PublicKey,
					SignPolicyId:  0,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1, 1}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa,778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 10),
			},
			wantSignRequest: &types.SignRequest{
				Id:      1,
				Creator: "testOwner",
				KeyIds:  []uint64{1, 1},
				DataForSigning: [][]byte{
					{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170},
					{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170},
				},
				Status:      types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:     types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
				MpcBtl:      10, // TODO: why is it not 10?
				ChildReqIds: []uint64{2, 3},
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "FAIL: empty payload",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_EDDSA_ED25519,
					PublicKey:     defaultEdDSAKey.PublicKey,
					SignPolicyId:  0,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "", 1000, 10),
			},
			wantErr: true,
		},
		{
			name: "FAIL: eddsa dataforsigning too long",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_EDDSA_ED25519,
					PublicKey:     defaultEdDSAKey.PublicKey,
					SignPolicyId:  0,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000020001046ac8f20fc76e35d909f857aa01a083afe317c4f1e507670b3181e9b6a8e2c3420000000000000000000000000000000000000000000000000000000000000000fa06c50dc073d2007afead92cfa489ee7ebe9fcf44253515b3301040938dba0a06a7d517192c568ee08a845f73d29788cf035c3145b21ab344d8062ea94000001ccc1d9ff619d8b760f74d67cf000b47f71d326411b502bafd5f5bed7a3c84f40201030103010404000000010200020c020000009f3d35bf01000000", 1000, 10),
			},
			wantErr: true,
		},
		{
			name: "PASS: valid bitcoin signature request with btl",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key: &types.Key{
					Id:            1,
					WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
					KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
					Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
					PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
					SignPolicyId:  1,
				},
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 10),
			},
			wantSignRequest: &types.SignRequest{
				Id:             1,
				Creator:        "testOwner",
				KeyIds:         []uint64{1},
				DataForSigning: [][]byte{{119, 143, 87, 47, 51, 175, 171, 131, 19, 101, 213, 46, 86, 58, 13, 221, 153, 105, 165, 53, 246, 34, 98, 18, 91, 176, 109, 67, 104, 28, 89, 170}},
				Status:         types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING,
				KeyType:        types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
				MpcBtl:         0,
			},
			want: &types.MsgNewSignatureRequestResponse{SigReqId: 1},
		},
		{
			name: "FAIL: key not found",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignatureRequest("testOwner", []uint64{5}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 0),
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
				msg: types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 0),
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
				msg:       types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "778f572f33afab831365d52e563a0ddd9969a535f62262125bb06d43681c59aa", 1000, 0),
			},
			wantErr: true,
		},
		{
			name: "FAIL: invalid length for data for signing",
			args: args{
				keyring:   &defaultKr,
				workspace: &defaultWs,
				key:       &defaultKey,
				msg:       types.NewMsgNewSignatureRequest("testOwner", []uint64{1}, "778f572f", 1000, 0),
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
			treasurymodule.InitGenesis(ctx, *tk, tGenesis)

			pGenesis := policytypes.GenesisState{
				Policies: []policytypes.Policy{
					policy1,
				},
			}
			policymodule.InitGenesis(ctx, *pk, pGenesis)

			msgSer := keeper.NewMsgServerImpl(*tk)
			got, err := msgSer.NewSignatureRequest(ctx, tt.args.msg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)

				gotSigReq, err := tk.SignRequestStore.Get(ctx, got.SigReqId)
				require.NoError(t, err)
				require.Equal(t, tt.wantSignRequest, &gotSigReq)

				act, err := pk.ActionStore.Get(ctx, 1)
				require.NoError(t, err)
				require.Equal(t, uint64(1000), act.Btl)
				require.Equal(t, tt.args.key.SignPolicyId, act.PolicyId)
			}
		})
	}
}
