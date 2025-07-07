package testutil

import (
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

var (
	DefaultControllerAddr                    = "0x5b9Ea8d5486D388a158F026c337DF950866dA5e9"
	DefaultEthTokenAddr                      = "0xC8CdeDd20cCb4c06884ac4C2fF952A0B7cC230a3"
	DefaultDepositKeyringAddr                = "keyring1pfnq7r04rept47gaf5cpdew2"
	DefaultEthMinterKeyID             uint64 = 2
	DefaultChangeAddressKeyIDs               = []uint64{3}
	DefaultUnstakerKeyID              uint64 = 4
	DefaultRewardsDepositKeyID        uint64 = 5
	DefaultStakerKeyID                uint64 = 6
	DefaultCompleterKeyID             uint64 = 7
	DefaultTestnetBitcoinProxyAddress        = "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
	DefaultMainnetBitcoinProxyAddress        = "zen1mgl98jt30nemuqtt5asldk49ju9lnx0pfke79q"
	// DefaultStrategyAddr               = "0x0000000000000000000000000000000000000000"
	// DefaultStakerKeyID = 0
	// DefaultBurnerKeyID = 0
	DefaultSolana = &zenbtctypes.Solana{
		SignerKeyId:        7,
		ProgramId:          "3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT",
		NonceAuthorityKey:  8,
		NonceAccountKey:    9,
		MintAddress:        "9oBkgQUkq8jvzK98D7Uib6GYSZZmjnZ6QEGJRrAeKnDj",
		FeeWallet:          "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
		Fee:                0,
		MultisigKeyAddress: "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
		Btl:                20,
	}
)
