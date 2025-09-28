package keeper

import (
    "context"
    "errors"
    "fmt"

    "cosmossdk.io/collections"
    "github.com/Zenrock-Foundation/zrchain/v6/shared"
    "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
    zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
    "golang.org/x/exp/slices"
)

// RequestHeaderBackfill enqueues a specific Bitcoin header height to be fetched
// by sidecars and processed through the normal ABCI flow. Access is gated by
// the same authority check as TriggerEventBackfill.
func (k msgServer) RequestHeaderBackfill(ctx context.Context, msg *types.MsgRequestHeaderBackfill) (*types.MsgRequestHeaderBackfillResponse, error) {
    if msg.Authority != shared.AdminAuthAddr && msg.Authority != k.authority {
        return nil, fmt.Errorf("invalid authority; expected %s or %s, got %s", shared.AdminAuthAddr, k.authority, msg.Authority)
    }

    if msg.Height <= 0 {
        return nil, fmt.Errorf("height must be > 0")
    }

    // Get current requested headers list (initialize if not present)
    requested, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
    if err != nil {
        if !errors.Is(err, collections.ErrNotFound) {
            return nil, err
        }
        requested = zenbtctypes.RequestedBitcoinHeaders{}
    }

    // Avoid duplicates
    if !slices.Contains(requested.Heights, msg.Height) {
        requested.Heights = append(requested.Heights, msg.Height)
    }

    if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requested); err != nil {
        return nil, err
    }

    return &types.MsgRequestHeaderBackfillResponse{}, nil
}
