package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	identity "github.com/Zenrock-Foundation/zrchain/v6/x/identity/module"
	identitytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/identity/testutil"
	identitytypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	treasury "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/module"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

func TestMsgSwap(t *testing.T) {
	t.SkipNow() // missing expected keeper calls
	k, _, ctx := setupMsgServer(t)
	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))
	// _ := sdk.UnwrapSDKContext(ctx)

	// default params
	testCases := []struct {
		name      string
		input     *types.MsgSwap
		expErr    bool
		expErrMsg string
		want      *types.MsgSwapResponse
		wantSwap  *types.Swap
	}{
		{
			name: "Pass: Happy Path",
			input: &types.MsgSwap{
				Creator:      "testCreator",
				Pair:         "rockbtc",
				Workspace:    "workspace14a2hpadpsy9h4auve2z8lw",
				AmountIn:     sdkmath.LegacyNewDec(100000),
				Yield:        false,
				SenderKey:    1,
				RecipientKey: 2,
			},
			expErr: false,
			want: &types.MsgSwapResponse{
				Id: 1,
			},
			wantSwap: &types.Swap{
				SwapId: 1,
				Status: types.SwapStatus_SWAP_STATUS_REQUESTED,
				Pair:   "rockbtc",
				Data: &types.SwapData{
					BaseToken: &validationtypes.AssetData{
						Asset: validationtypes.Asset_ROCK,
					},
					QuoteToken: &validationtypes.AssetData{
						Asset: validationtypes.Asset_BTC,
					},
					Price:     sdkmath.LegacyNewDec(100000),
					AmountIn:  sdkmath.LegacyNewDec(100000),
					AmountOut: sdkmath.LegacyNewDec(100000),
				},
				SenderKeyId:    1,
				RecipientKeyId: 2,
				Workspace:      "workspace14a2hpadpsy9h4auve2z8lw",
				ZenbtcYield:    false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			zk := keepers.ZenexKeeper
			ik := keepers.IdentityKeeper
			tk := keepers.TreasuryKeeper
			ctx := keepers.Ctx
			msgSer := keeper.NewMsgServerImpl(*zk)

			identityGenesis := identitytypes.GenesisState{
				Workspaces: []identitytypes.Workspace{identitytestutil.DefaultWs},
			}
			identity.InitGenesis(ctx, *ik, identityGenesis)

			treasuryGenesis := treasurytypes.GenesisState{
				Keys: treasurytestutil.DefaultKeys,
			}

			treasury.InitGenesis(ctx, *tk, treasuryGenesis)

			got, err := msgSer.Swap(ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)

				gotSwap, err := zk.SwapsStore.Get(ctx, got.Id)
				require.NoError(t, err)
				require.Equal(t, tc.wantSwap, &gotSwap)
			}
		})
	}
}
