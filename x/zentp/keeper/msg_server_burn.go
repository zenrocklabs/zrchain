package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
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

	// check if there's sufficient balance
	moduleAddr := k.accountKeeper.GetModuleAddress(msg.ModuleAccount)
	balance := k.bankKeeper.GetBalance(ctx, moduleAddr, msg.Denom)
	if balance.IsLT(coins[0]) {
		return nil, errorsmod.Wrapf(types.ErrInsufficientBalance, "insufficient balance; required %s, available %s", coins[0].String(), balance.String())
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

	burnAmount := math.NewIntFromUint64(msg.Amount)
	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, burnAmount))

	if len(coins) == 0 {
		return nil, errorsmod.Wrapf(errors.New("invalid amount"), "invalid amount; expected %d", msg.Amount)
	}

	return coins, nil
}
