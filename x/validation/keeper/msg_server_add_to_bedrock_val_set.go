package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k msgServer) AddToBedrockValSet(ctx context.Context, msg *types.MsgAddToBedrockValSet) (*types.MsgAddToBedrockValSetResponse, error) {
	if msg.Authority != shared.AdminAuthAddr && msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority; expected %s or %s, got %s", shared.AdminAuthAddr, k.authority, msg.Authority)
	}

	// Validate that the validator address is valid
	if msg.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address must be provided")
	}

	// Verify the validator exists
	if _, err := k.GetZenrockValidatorFromBech32(ctx, msg.ValidatorAddress); err != nil {
		return nil, fmt.Errorf("validator not found: %w", err)
	}

	// Add to bedrock validator set
	if err := k.BedrockValidatorSet.Set(ctx, msg.ValidatorAddress, true); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.rebalanceAllBedrockAssets(sdkCtx); err != nil {
		return nil, fmt.Errorf("failed to rebalance bedrock assets: %w", err)
	}

	return &types.MsgAddToBedrockValSetResponse{}, nil
}
