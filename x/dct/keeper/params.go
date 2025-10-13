package keeper

import (
    "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

var (
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
    return &types.Params{Assets: []types.AssetParams{defaultZenZECCfg}}
}
