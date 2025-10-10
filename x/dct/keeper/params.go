package keeper

import (
    "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

var (
    defaultZenBTCSolana = &types.Solana{
        SignerKeyId:         7,
        ProgramId:           "3jo4mdc6QbGRigia2jvmKShbmz3aWq4Y8bgUXfur5StT",
        NonceAccountKey:     9,
        NonceAuthorityKey:   8,
        MintAddress:         "9oBkgQUkq8jvzK98D7uib6GYSZzmjnZ6QEGJRrAeKnDj",
        FeeWallet:           "FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd",
        Fee:                 0,
        MultisigKeyAddress:  "8cmZY2id22vxpXs2H3YYQNARuPHNuYwa7jipW1q1v9Fy",
        Btl:                 20,
        EventStoreProgramId: "Hsu6LJz42sZhs2GvF9yzD6L9n2AZTeHnjDx6Cp4DvEdf",
    }

    defaultZenBTCCfg = types.AssetParams{
        Asset:               types.Asset_ASSET_ZENBTC,
        DepositKeyringAddr:  "keyring1pfnq7r04rept47gaf5cpdew2",
        StakerKeyId:         6,
        EthMinterKeyId:      2,
        UnstakerKeyId:       4,
        CompleterKeyId:      7,
        RewardsDepositKeyId: 5,
        ChangeAddressKeyIds: []uint64{3},
        ProxyAddress:        "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty",
        EthTokenAddr:        "0xC8CdeDd20cCb4c06884ac4C2fF952A0B7cC230a3",
        ControllerAddr:      "0x5b9Ea8d5486D388a158F026c337DF950866dA5e9",
        Solana:              defaultZenBTCSolana,
    }

    defaultZenZECSolana = &types.Solana{
        SignerKeyId:         17,
        ProgramId:           "7q3u7bL1nKzftYFbwUE8yuzhzbQxFwM7xyMh1cZFzenZ",
        NonceAccountKey:     19,
        NonceAuthorityKey:   18,
        MintAddress:         "ZC3hZPnfYg1y5SP62x9XyLJmMnt4zUpfcu8JZenZEC",
        FeeWallet:           "FvzecFee1YdzWm7Nzvy2Kc7CquBrEPNvkMRPQHrfeee",
        Fee:                 0,
        MultisigKeyAddress:  "7ZeCmxg9APm9wJgZtPy1X6Qt23afVqzVrqVZC9HexbBP",
        Btl:                 20,
        EventStoreProgramId: "7zecEvtStreProg11111111111111111111111111111",
    }

    defaultZenZECCfg = types.AssetParams{
        Asset:               types.Asset_ASSET_ZENZEC,
        DepositKeyringAddr:  "keyring1zenzec000000000000000000000000000000",
        StakerKeyId:         16,
        EthMinterKeyId:      0,
        UnstakerKeyId:       14,
        CompleterKeyId:      15,
        RewardsDepositKeyId: 12,
        ChangeAddressKeyIds: []uint64{13},
        ProxyAddress:        "zen1zecproxytest0000000000000000000",
        EthTokenAddr:        "",
        ControllerAddr:      "",
        Solana:              defaultZenZECSolana,
    }
)

func DefaultParams() *types.Params {
    return &types.Params{Assets: []types.AssetParams{defaultZenBTCCfg, defaultZenZECCfg}}
}
