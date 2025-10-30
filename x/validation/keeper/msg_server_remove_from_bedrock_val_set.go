package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

func (k msgServer) RemoveFromBedrockValSet(ctx context.Context, msg *types.MsgRemoveFromBedrockValSet) (*types.MsgRemoveFromBedrockValSetResponse, error) {
	if msg.Authority != shared.AdminAuthAddr && msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority; expected %s or %s, got %s", shared.AdminAuthAddr, k.authority, msg.Authority)
	}

	// Validate that the validator address is valid
	if msg.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address must be provided")
	}

	// Verify the validator exists and capture current bedrock balances
	validator, err := k.GetZenrockValidatorFromBech32(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("validator not found: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	removedBalances := make(map[types.Asset]sdkmath.Int)
	for _, tokenData := range validator.TokensBedrock {
		if tokenData == nil || tokenData.Amount.IsZero() {
			continue
		}
		amt := tokenData.Amount
		if existing, ok := removedBalances[tokenData.Asset]; ok {
			removedBalances[tokenData.Asset] = existing.Add(amt)
		} else {
			removedBalances[tokenData.Asset] = amt
		}
		tokenData.Amount = sdkmath.ZeroInt()
	}

	if err := k.SetValidator(ctx, validator); err != nil {
		return nil, fmt.Errorf("failed to clear validator bedrock balances: %w", err)
	}

	// Remove from bedrock validator set
	if err := k.BedrockValidatorSet.Remove(ctx, msg.ValidatorAddress); err != nil {
		return nil, err
	}

	// Redistribute removed stake across the remaining bedrock validators
	for asset, amount := range removedBalances {
		if amount.IsZero() {
			continue
		}
		if err := k.distributeBedrockTokens(sdkCtx, asset, amount); err != nil {
			return nil, fmt.Errorf("failed to redistribute %s bedrock stake: %w", asset.String(), err)
		}
	}

	if err := k.rebalanceAllBedrockAssets(sdkCtx); err != nil {
		return nil, fmt.Errorf("failed to rebalance bedrock assets: %w", err)
	}

	return &types.MsgRemoveFromBedrockValSetResponse{}, nil
}
