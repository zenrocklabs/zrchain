package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"golang.org/x/exp/slices"
)

func (k msgServer) TriggerEventBackfill(ctx context.Context, msg *types.MsgTriggerEventBackfill) (*types.MsgTriggerEventBackfillResponse, error) {
	if msg.Authority != shared.AdminAuthAddr && msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority; expected %s or %s, got %s", shared.AdminAuthAddr, k.authority, msg.Authority)
	}

	// Validate the request based on its event type to ensure it's actionable.
	switch msg.EventType {
	case types.EventType_EVENT_TYPE_ZENTP_BURN:
		// First, check if a burn event with this transaction hash already exists in zentp.
		existingBurns, err := k.zentpKeeper.GetBurns(ctx, "", "", msg.TxHash)
		if err != nil {
			return nil, err
		}
		if len(existingBurns) > 0 {
			return nil, fmt.Errorf("burn event with tx hash '%s' already exists", msg.TxHash)
		}

		// ZenTP burns are always on Solana. Validate the chain ID is a supported Solana network.
		if _, err := types.ValidateSolanaChainID(ctx, msg.Caip2ChainId); err != nil {
			return nil, err
		}
		// Basic validation for TxHash - ensure it's not empty.
		if msg.TxHash == "" {
			return nil, fmt.Errorf("transaction hash cannot be empty for zentp burn backfill")
		}
	default:
		return nil, fmt.Errorf("currently unsupported backfill request type: %s", msg.EventType)
	}

	backfillRequests, err := k.BackfillRequests.Get(ctx)
	if err != nil {
		// If no backfill requests exist yet, it's not an error. Initialize it.
		if errors.Is(err, collections.ErrNotFound) {
			backfillRequests = types.BackfillRequests{}
		} else {
			return nil, err
		}
	}

	// Prevent duplicate requests by checking the content of the requests.
	if slices.ContainsFunc(backfillRequests.Requests, func(req *types.MsgTriggerEventBackfill) bool {
		return req.Caip2ChainId == msg.Caip2ChainId && req.TxHash == msg.TxHash
	}) {
		return nil, fmt.Errorf("backfill request already exists")
	}

	backfillRequests.Requests = append(backfillRequests.Requests, msg)

	if err = k.BackfillRequests.Set(ctx, backfillRequests); err != nil {
		return nil, err
	}

	return &types.MsgTriggerEventBackfillResponse{}, nil
}
