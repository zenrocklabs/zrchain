package keeper_test

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"github.com/stretchr/testify/require"
)

func TestListSupportedAssets(t *testing.T) {
	t.Run("returns zenZEC when params has zenZEC", func(t *testing.T) {
		k, ctx := DctKeeper(t)

		params := keeper.DefaultParams()
		require.NoError(t, k.SetParams(ctx, *params))

		assets, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)
		require.Contains(t, assets, types.Asset_ASSET_ZENZEC, "Should contain ASSET_ZENZEC from params")
		require.Len(t, assets, 1, "Should have exactly one asset (zenZEC)")
	})

	t.Run("always returns zenZEC even with empty params", func(t *testing.T) {
		k, ctx := DctKeeper(t)

		// Set empty params (no assets configured)
		emptyParams := &types.Params{Assets: []types.AssetParams{}}
		require.NoError(t, k.SetParams(ctx, *emptyParams))

		assets, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)
		require.Contains(t, assets, types.Asset_ASSET_ZENZEC, "Should always contain ASSET_ZENZEC even with empty params")
		require.Len(t, assets, 1, "Should have exactly one asset (zenZEC)")
	})

	t.Run("returns zenZEC and other configured assets", func(t *testing.T) {
		k, ctx := DctKeeper(t)

		// Configure zenBTC as well (even though it's not normally in DCT)
		paramsWithMultiple := &types.Params{
			Assets: []types.AssetParams{
				{Asset: types.Asset_ASSET_ZENBTC}, // Some other asset
			},
		}
		require.NoError(t, k.SetParams(ctx, *paramsWithMultiple))

		assets, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)
		require.Contains(t, assets, types.Asset_ASSET_ZENZEC, "Should always contain ASSET_ZENZEC")
		require.Contains(t, assets, types.Asset_ASSET_ZENBTC, "Should contain ASSET_ZENBTC from params")
		require.Len(t, assets, 2, "Should have two assets")
	})

	t.Run("filters out unspecified assets", func(t *testing.T) {
		k, ctx := DctKeeper(t)

		paramsWithUnspecified := &types.Params{
			Assets: []types.AssetParams{
				{Asset: types.Asset_ASSET_UNSPECIFIED}, // Should be filtered out
			},
		}
		require.NoError(t, k.SetParams(ctx, *paramsWithUnspecified))

		assets, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)
		require.Contains(t, assets, types.Asset_ASSET_ZENZEC, "Should always contain ASSET_ZENZEC")
		require.NotContains(t, assets, types.Asset_ASSET_UNSPECIFIED, "Should not contain ASSET_UNSPECIFIED")
		require.Len(t, assets, 1, "Should have exactly one asset (zenZEC)")
	})

	t.Run("deduplicates zenZEC if present in params", func(t *testing.T) {
		k, ctx := DctKeeper(t)

		// zenZEC appears in params (normal case)
		params := keeper.DefaultParams()
		require.NoError(t, k.SetParams(ctx, *params))

		assets, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)

		// Count how many times zenZEC appears
		zenZECCount := 0
		for _, asset := range assets {
			if asset == types.Asset_ASSET_ZENZEC {
				zenZECCount++
			}
		}
		require.Equal(t, 1, zenZECCount, "ASSET_ZENZEC should appear exactly once, not duplicated")
	})

	t.Run("returns assets in deterministic sorted order", func(t *testing.T) {
		k, ctx := DctKeeper(t)

		// Configure multiple assets
		paramsWithMultiple := &types.Params{
			Assets: []types.AssetParams{
				{Asset: types.Asset_ASSET_ZENBTC},  // Value: 1
				{Asset: types.Asset_ASSET_ZENZEC},  // Value: 2
			},
		}
		require.NoError(t, k.SetParams(ctx, *paramsWithMultiple))

		// Call multiple times and verify same order
		assets1, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)

		assets2, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)

		assets3, err := k.ListSupportedAssets(ctx)
		require.NoError(t, err)

		// Verify all calls return same order
		require.Equal(t, assets1, assets2, "Multiple calls should return same order")
		require.Equal(t, assets2, assets3, "Multiple calls should return same order")

		// Verify they're sorted by enum value (ZENBTC=1, ZENZEC=2)
		require.Len(t, assets1, 2, "Should have 2 assets")
		require.Equal(t, types.Asset_ASSET_ZENBTC, assets1[0], "ZENBTC (1) should come before ZENZEC (2)")
		require.Equal(t, types.Asset_ASSET_ZENZEC, assets1[1], "ZENZEC (2) should come after ZENBTC (1)")
	})
}
