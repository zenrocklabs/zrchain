package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k msgServer) SetSolanaCounters(ctx context.Context, msg *types.MsgSetSolanaCounters) (*types.MsgSetSolanaCountersResponse, error) {
	if msg.Authority != shared.AdminAuthAddr && msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority; expected %s or %s, got %s", shared.AdminAuthAddr, k.authority, msg.Authority)
	}

	if msg.Asset == dcttypes.Asset_ASSET_UNSPECIFIED {
		return nil, fmt.Errorf("asset must be specified")
	}

	// Key by asset name - each asset has one Solana program and one set of counters
	assetKey := msg.Asset.String()

	// Set the counters
	counters := types.SolanaCounters{
		MintCounter:       msg.MintCounter,
		RedemptionCounter: msg.RedemptionCounter,
	}

	if err := k.SolanaCounters.Set(ctx, assetKey, counters); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.Logger(ctx).Info("Set Solana counters",
		"asset", assetKey,
		"mint_counter", msg.MintCounter,
		"redemption_counter", msg.RedemptionCounter,
		"height", sdkCtx.BlockHeight(),
	)

	return &types.MsgSetSolanaCountersResponse{}, nil
}
