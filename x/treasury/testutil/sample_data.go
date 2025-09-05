package testutil

import (
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

var DefaultKeys = []types.Key{
	{
		Id:            1,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
		Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
	},
	{
		Id:            2,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
		Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     []byte{3, 49, 216, 1, 30, 80, 59, 231, 205, 4, 237, 36, 101, 79, 9, 119, 30, 114, 19, 70, 242, 44, 91, 95, 238, 20, 245, 197, 193, 86, 22, 120, 35},
	},
}
