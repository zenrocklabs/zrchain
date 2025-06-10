package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetSolanaROCKSupply(goCtx context.Context, msg *types.MsgSetSolanaROCKSupply) (*types.MsgSetSolanaROCKSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetAuthority() != msg.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	supply := math.NewIntFromUint64(msg.Amount)

	if err := k.Keeper.SetSolanaROCKSupply(ctx, supply); err != nil {
		return nil, err
	}

	return &types.MsgSetSolanaROCKSupplyResponse{}, nil
}
