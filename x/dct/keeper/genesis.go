package keeper

import (
    "context"
    "sort"

    "cosmossdk.io/collections"

    dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

// DefaultGenesis returns the default genesis state for the DCT module.
func DefaultGenesis() *dcttypes.GenesisState {
    params := DefaultParams()
    assetStates := make([]dcttypes.AssetGenesisState, len(params.Assets))
    for i, ap := range params.Assets {
        assetStates[i] = dcttypes.AssetGenesisState{Asset: ap.Asset}
    }
    return &dcttypes.GenesisState{
        Params: *params,
        Assets: assetStates,
    }
}

func (k Keeper) ExportState(ctx context.Context, genState *dcttypes.GenesisState) error {
    params, err := k.GetParams(ctx)
    if err != nil {
        return err
    }
    genState.Params = params

    assetIndex := make(map[dcttypes.Asset]*dcttypes.AssetGenesisState)
    for _, ap := range params.Assets {
        assetIndex[ap.Asset] = &dcttypes.AssetGenesisState{Asset: ap.Asset}
    }

    ensureAsset := func(asset dcttypes.Asset) *dcttypes.AssetGenesisState {
        if asset == dcttypes.Asset_ASSET_UNSPECIFIED {
            return nil
        }
        if _, ok := assetIndex[asset]; !ok {
            assetIndex[asset] = &dcttypes.AssetGenesisState{Asset: asset}
        }
        return assetIndex[asset]
    }

    // Lock transactions
    if err := k.LockTransactions.Walk(ctx, nil, func(key collections.Pair[string, string], value dcttypes.LockTransaction) (bool, error) {
        asset, err := k.assetFromKey(key.K1())
        if err != nil {
            return true, err
        }
        state := ensureAsset(asset)
        if state == nil {
            return false, nil
        }
        if state.LockTransactions == nil {
            state.LockTransactions = make(map[string]dcttypes.LockTransaction)
        }
        state.LockTransactions[key.K2()] = value
        return false, nil
    }); err != nil {
        return err
    }

    // Pending mint transactions
    if err := k.PendingMintTransactions.Walk(ctx, nil, func(key collections.Pair[string, uint64], value dcttypes.PendingMintTransaction) (bool, error) {
        asset, err := k.assetFromKey(key.K1())
        if err != nil {
            return true, err
        }
        state := ensureAsset(asset)
        if state == nil {
            return false, nil
        }
        if state.PendingMintTransactions == nil {
            state.PendingMintTransactions = make(map[uint64]dcttypes.PendingMintTransaction)
        }
        value.Asset = asset
        state.PendingMintTransactions[key.K2()] = value
        return false, nil
    }); err != nil {
        return err
    }

    // First pending stake transactions
    if err := k.FirstPendingStakeTransaction.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.FirstPendingStakeTransaction = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    // First pending Solana mint transactions
    if err := k.FirstPendingSolMintTransaction.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.FirstPendingSolMintTransaction = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    // First pending ETH mint
    if err := k.FirstPendingEthMintTransaction.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.FirstPendingEthMintTransaction = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    // Pending mint counts
    if err := k.PendingMintTransactionCount.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.PendingMintTransactionCount = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    // Burn events and counts
    if err := k.BurnEvents.Walk(ctx, nil, func(key collections.Pair[string, uint64], value dcttypes.BurnEvent) (bool, error) {
        asset, err := k.assetFromKey(key.K1())
        if err != nil {
            return true, err
        }
        state := ensureAsset(asset)
        if state == nil {
            return false, nil
        }
        if state.BurnEvents == nil {
            state.BurnEvents = make(map[uint64]dcttypes.BurnEvent)
        }
        value.Asset = asset
        state.BurnEvents[key.K2()] = value
        return false, nil
    }); err != nil {
        return err
    }

    if err := k.BurnEventCount.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.BurnEventCount = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    if err := k.FirstPendingBurnEvent.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.FirstPendingBurnEvent = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    // Redemptions
    if err := k.Redemptions.Walk(ctx, nil, func(key collections.Pair[string, uint64], value dcttypes.Redemption) (bool, error) {
        asset, err := k.assetFromKey(key.K1())
        if err != nil {
            return true, err
        }
        state := ensureAsset(asset)
        if state == nil {
            return false, nil
        }
        if state.Redemptions == nil {
            state.Redemptions = make(map[uint64]dcttypes.Redemption)
        }
        value.Data.Asset = asset
        state.Redemptions[key.K2()] = value
        return false, nil
    }); err != nil {
        return err
    }

    if err := k.FirstPendingRedemption.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.FirstPendingRedemption = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    if err := k.FirstRedemptionAwaitingSign.Walk(ctx, nil, func(assetKey string, value uint64) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            state.FirstRedemptionAwaitingSign = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    // Supply snapshots
    if err := k.Supply.Walk(ctx, nil, func(assetKey string, value dcttypes.Supply) (bool, error) {
        asset, err := k.assetFromKey(assetKey)
        if err != nil {
            return true, err
        }
        if state := ensureAsset(asset); state != nil {
            value.Asset = asset
        state.Supply = value
        }
        return false, nil
    }); err != nil {
        return err
    }

    assetStates := make([]dcttypes.AssetGenesisState, 0, len(assetIndex))
    for _, state := range assetIndex {
        assetStates = append(assetStates, *state)
    }

    sort.Slice(assetStates, func(i, j int) bool {
        return assetStates[i].Asset < assetStates[j].Asset
    })

    genState.Assets = assetStates
    return nil
}

func (k Keeper) InitGenesis(ctx context.Context, genState *dcttypes.GenesisState) error {
    if err := k.SetParams(ctx, genState.Params); err != nil {
        return err
    }

    for _, assetState := range genState.Assets {
        asset := assetState.Asset
        if asset == dcttypes.Asset_ASSET_UNSPECIFIED {
            continue
        }

        assetKey, err := k.getAssetKey(asset)
        if err != nil {
            return err
        }

        if assetState.LockTransactions != nil {
            for key, value := range assetState.LockTransactions {
                pairKey, err := k.lockKey(asset, key)
                if err != nil {
                    return err
                }
                if err := k.LockTransactions.Set(ctx, pairKey, value); err != nil {
                    return err
                }
            }
        }

        if assetState.PendingMintTransactions != nil {
            for id, tx := range assetState.PendingMintTransactions {
                tx.Asset = asset
                key, err := k.pendingMintKey(asset, id)
                if err != nil {
                    return err
                }
                if err := k.PendingMintTransactions.Set(ctx, key, tx); err != nil {
                    return err
                }
            }
        }

        if assetState.PendingMintTransactionCount != 0 {
            if err := k.PendingMintTransactionCount.Set(ctx, assetKey, assetState.PendingMintTransactionCount); err != nil {
                return err
            }
        }

        if assetState.FirstPendingStakeTransaction != 0 {
            if err := k.FirstPendingStakeTransaction.Set(ctx, assetKey, assetState.FirstPendingStakeTransaction); err != nil {
                return err
            }
        }

        if assetState.FirstPendingSolMintTransaction != 0 {
            if err := k.FirstPendingSolMintTransaction.Set(ctx, assetKey, assetState.FirstPendingSolMintTransaction); err != nil {
                return err
            }
        }

        if assetState.FirstPendingEthMintTransaction != 0 {
            if err := k.FirstPendingEthMintTransaction.Set(ctx, assetKey, assetState.FirstPendingEthMintTransaction); err != nil {
                return err
            }
        }

        if assetState.BurnEvents != nil {
            for id, burn := range assetState.BurnEvents {
                burn.Asset = asset
                key, err := k.burnKey(asset, id)
                if err != nil {
                    return err
                }
                if err := k.BurnEvents.Set(ctx, key, burn); err != nil {
                    return err
                }
            }
        }

        if assetState.BurnEventCount != 0 {
            if err := k.BurnEventCount.Set(ctx, assetKey, assetState.BurnEventCount); err != nil {
                return err
            }
        }

        if assetState.FirstPendingBurnEvent != 0 {
            if err := k.FirstPendingBurnEvent.Set(ctx, assetKey, assetState.FirstPendingBurnEvent); err != nil {
                return err
            }
        }

        if assetState.Redemptions != nil {
            for id, red := range assetState.Redemptions {
                red.Data.Asset = asset
                key, err := k.redemptionKey(asset, id)
                if err != nil {
                    return err
                }
                if err := k.Redemptions.Set(ctx, key, red); err != nil {
                    return err
                }
            }
        }

        if assetState.FirstPendingRedemption != 0 {
            if err := k.FirstPendingRedemption.Set(ctx, assetKey, assetState.FirstPendingRedemption); err != nil {
                return err
            }
        }

        if assetState.FirstRedemptionAwaitingSign != 0 {
            if err := k.FirstRedemptionAwaitingSign.Set(ctx, assetKey, assetState.FirstRedemptionAwaitingSign); err != nil {
                return err
            }
        }

        if assetState.Supply.Asset != dcttypes.Asset_ASSET_UNSPECIFIED {
            if err := k.Supply.Set(ctx, assetKey, assetState.Supply); err != nil {
                return err
            }
        }
    }

    return nil
}
