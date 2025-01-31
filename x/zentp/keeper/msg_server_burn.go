package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// checks if message is called by the x/gov module
	if k.GetAuthority() != msg.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	// checks if the module account exists
	if !k.accountKeeper.HasAccount(ctx, k.accountKeeper.GetModuleAddress(msg.ModuleAccount)) {
		return nil, errorsmod.Wrapf(types.ErrInvalidModuleAccount, "module account %s does not exist", msg.ModuleAccount)
	}

	// creates sdk.Coins to burn from message
	coins, err := prepareBurn(msg)
	if err != nil {
		return nil, err
	}

	// burns the specified amount of tokens from the specified address
	err = k.bankKeeper.BurnCoins(ctx, msg.ModuleAccount, coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}

// returns sdk.Coins to burn from message
func prepareBurn(msg *types.MsgBurn) (sdk.Coins, error) {

	burnAmount := math.NewInt(msg.Amount)
	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, burnAmount))

	return coins, nil
}
