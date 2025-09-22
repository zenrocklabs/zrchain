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
	{
		Id:            3,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
		Type:          types.KeyType_KEY_TYPE_EDDSA_ED25519,
		PublicKey:     []byte{0x29, 0xdb, 0xaa, 0xd6, 0x75, 0x87, 0x8b, 0x2f, 0x16, 0xcb, 0x14, 0x82, 0x7a, 0x4c, 0x1a, 0x41, 0xbb, 0xd6, 0x3c, 0x3b, 0x60, 0x1b, 0xdc, 0x5e, 0x27, 0x7d, 0x00, 0xb1, 0x20, 0xff, 0xee, 0xce},
	},
	{
		Id:            4,
		WorkspaceAddr: "workspace14a2hpadpsy9h4auve2z8lw",
		KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
		Type:          types.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		PublicKey:     []byte{0x29, 0xdb, 0xaa, 0xd6, 0x75, 0x87, 0x8b, 0x2f, 0x16, 0xcb, 0x14, 0x82, 0x7a, 0x4c, 0x1a, 0x41, 0xbb, 0xd6, 0x3c, 0x3b, 0x60, 0x1b, 0xdc, 0x5e, 0x27, 0x7d, 0x00, 0xb1, 0x20, 0xff, 0xee, 0xce},
	},
	{
		Id:            5,
		WorkspaceAddr: "workspace1hadsaz03etk33jxkvvr9k5",
		KeyringAddr:   "keyring1pfnq7r04rept47gaf5cpdew2",
		Type:          types.KeyType_KEY_TYPE_ECDSA_SECP256K1,
		PublicKey:     []byte{0x29, 0xdb, 0xaa, 0xd6, 0x75, 0x87, 0x8b, 0x2f, 0x16, 0xcb, 0x14, 0x82, 0x7a, 0x4c, 0x1a, 0x41, 0xbb, 0xd6, 0x3c, 0x3b, 0x60, 0x1b, 0xdc, 0x5e, 0x27, 0x7d, 0x00, 0xb1, 0x20, 0xff, 0xee, 0xce},
	},
}
