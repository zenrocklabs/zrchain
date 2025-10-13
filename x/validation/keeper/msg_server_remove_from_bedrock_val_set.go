package keeper

import (
	"context"
	"fmt"

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

	// Verify the validator exists
	valAddr, err := k.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %w", err)
	}

	_, err = k.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, fmt.Errorf("validator not found: %w", err)
	}

	// Remove from bedrock validator set
	if err := k.BedrockValidatorSet.Remove(ctx, msg.ValidatorAddress); err != nil {
		return nil, err
	}

	return &types.MsgRemoveFromBedrockValSetResponse{}, nil
}
