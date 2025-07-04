package keeper_test

import (
	"fmt"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	zentp "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/testutil"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/stretchr/testify/require"
)

func TestStatsQueryWithPagination(t *testing.T) {
	k, ctx := keepertest.ZentpKeeper(t)
	msgs := make([]types.Bridge, 205)
	burns := make([]types.Bridge, 205)

	addr := "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx"

	var totalMintAmount uint64
	var totalBurnAmount uint64

	for i := range msgs {
		amount := uint64(i + 1)
		msgs[i] = types.Bridge{
			Creator: fmt.Sprintf("%s%d", addr, i),
			Amount:  amount,
			Denom:   "urock",
			State:   types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		}
		totalMintAmount += amount
		k.UpdateMint(ctx, uint64(i), &msgs[i])
	}

	for i := range burns {
		amount := uint64(i + 1)
		burns[i] = types.Bridge{
			Creator:          fmt.Sprintf("%s%d", addr, i),
			RecipientAddress: fmt.Sprintf("%s%d", addr, i+len(msgs)),
			Amount:           amount,
			Denom:            "urock",
			State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
		}
		totalBurnAmount += amount
		k.UpdateBurn(ctx, uint64(i), &burns[i])
	}

	stats, err := k.Stats(ctx, &types.QueryStatsRequest{})
	require.NoError(t, err)

	require.Equal(t, totalMintAmount, stats.TotalMinted)
	require.Equal(t, uint64(len(msgs)), stats.MintsCount)
	require.Equal(t, totalBurnAmount, stats.TotalBurned)
	require.Equal(t, uint64(len(burns)), stats.BurnsCount)

	// Test with address
	stats, err = k.Stats(ctx, &types.QueryStatsRequest{Address: msgs[0].Creator})
	require.NoError(t, err)
	require.Equal(t, msgs[0].Amount, stats.TotalMinted)
	require.Equal(t, uint64(1), stats.MintsCount)

	stats, err = k.Stats(ctx, &types.QueryStatsRequest{Address: burns[0].RecipientAddress})
	require.NoError(t, err)
	require.Equal(t, burns[0].Amount, stats.TotalBurned)
	require.Equal(t, uint64(1), stats.BurnsCount)
}

func TestStatsQuery(t *testing.T) {

	tests := []struct {
		desc     string
		mints    []types.Bridge
		burns    []types.Bridge
		request  *types.QueryStatsRequest
		response *types.QueryStatsResponse
		err      error
	}{
		{
			desc:    "Total",
			mints:   testutil.DefaultMints,
			burns:   testutil.DefaultBurns,
			request: &types.QueryStatsRequest{},
			response: &types.QueryStatsResponse{
				TotalMinted: 4000100,
				MintsCount:  2,
				TotalBurned: 3000050,
				BurnsCount:  6,
			},
		},
		{
			desc:    "By Address",
			mints:   testutil.DefaultMints,
			burns:   testutil.DefaultBurns,
			request: &types.QueryStatsRequest{Address: testutil.DefaultMints[0].Creator},
			response: &types.QueryStatsResponse{
				TotalMinted: 100,
				MintsCount:  1,
				TotalBurned: 1100000,
				BurnsCount:  2,
			},
		},
		{
			desc:    "By Denom",
			mints:   testutil.DefaultMints,
			burns:   testutil.DefaultBurns,
			request: &types.QueryStatsRequest{Denom: "urock"},
			response: &types.QueryStatsResponse{
				TotalMinted: 4000100,
				MintsCount:  2,
				TotalBurned: 3000050,
				BurnsCount:  6,
			},
		},
		{
			desc:    "Nil request",
			request: nil,
			err:     fmt.Errorf("request is nil"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			zk, ctx := keepertest.ZentpKeeper(t)

			genesis := types.GenesisState{
				Params: types.DefaultParams(),
				Burns:  tc.burns,
				Mints:  tc.mints,
			}

			zentp.InitGenesis(ctx, zk, genesis)
			response, err := zk.Stats(ctx, tc.request)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
