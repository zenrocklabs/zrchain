package keeper_test

import (
	"fmt"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	zentp "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/stretchr/testify/require"
)

func TestQuerySolanaRockSupply(t *testing.T) {

	tests := []struct {
		desc             string
		solanaRockSupply uint64
		request          *types.QuerySolanaROCKSupplyRequest
		response         *types.QuerySolanaROCKSupplyResponse
		err              error
	}{
		{
			desc:             "Total",
			solanaRockSupply: 1000000,
			request:          &types.QuerySolanaROCKSupplyRequest{},
			response: &types.QuerySolanaROCKSupplyResponse{
				Amount: 1000000,
			},
		},
		{
			desc:    "Nil request",
			request: nil,
			err:     fmt.Errorf("request is nil"),
		},
		{
			desc:    "Zero supply",
			request: &types.QuerySolanaROCKSupplyRequest{},
			response: &types.QuerySolanaROCKSupplyResponse{
				Amount: 0,
			},
		},
		{
			desc:             "Large supply",
			solanaRockSupply: 999999999999999,
			request:          &types.QuerySolanaROCKSupplyRequest{},
			response: &types.QuerySolanaROCKSupplyResponse{
				Amount: 999999999999999,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			zk, ctx := keepertest.ZentpKeeper(t)

			genesis := types.GenesisState{
				SolanaRockSupply: tc.solanaRockSupply,
			}

			zentp.InitGenesis(ctx, zk, genesis)
			response, err := zk.QuerySolanaROCKSupply(ctx, tc.request)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
