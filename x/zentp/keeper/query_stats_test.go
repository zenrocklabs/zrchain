package keeper_test

import (
	"fmt"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/testutil/nullify"
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
	k, ctx := keepertest.ZentpKeeper(t)
	msgs := make([]types.Bridge, 2)
	burns := make([]types.Bridge, 2)

	addr := "zenrock1vpg325p5a7v73jaqyk02muh6n22p852h60anxx"

	msgs[0] = types.Bridge{
		Creator: addr,
		Amount:  100,
		Denom:   "urock",
		State:   types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
	}
	msgs[1] = types.Bridge{
		Creator: addr,
		Amount:  200,
		Denom:   "urock",
		State:   types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
	}

	burns[0] = types.Bridge{
		Creator:          addr,
		RecipientAddress: addr,
		Amount:           50,
		Denom:            "urock",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
	}
	burns[1] = types.Bridge{
		Creator:          addr,
		RecipientAddress: addr,
		Amount:           25,
		Denom:            "urock",
		State:            types.BridgeStatus_BRIDGE_STATUS_COMPLETED,
	}

	k.UpdateMint(ctx, 0, &msgs[0])
	k.UpdateMint(ctx, 1, &msgs[1])
	k.UpdateBurn(ctx, 0, &burns[0])
	k.UpdateBurn(ctx, 1, &burns[1])

	tests := []struct {
		desc     string
		request  *types.QueryStatsRequest
		response *types.QueryStatsResponse
		err      error
	}{
		{
			desc:    "Total",
			request: &types.QueryStatsRequest{},
			response: &types.QueryStatsResponse{
				TotalMinted: 300,
				MintsCount:  2,
				TotalBurned: 75,
				BurnsCount:  2,
			},
		},
		{
			desc:    "By Address",
			request: &types.QueryStatsRequest{Address: addr},
			response: &types.QueryStatsResponse{
				TotalMinted: 300,
				MintsCount:  2,
				TotalBurned: 75,
				BurnsCount:  2,
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
			response, err := k.Stats(ctx, tc.request)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
