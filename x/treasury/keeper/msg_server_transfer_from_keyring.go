package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/app/params"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferFromKeyring(goCtx context.Context, msg *types.MsgTransferFromKeyring) (*types.MsgTransferFromKeyringResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keyring, err := k.identityKeeper.KeyringStore.Get(ctx, msg.Keyring)
	if err != nil {
		return nil, fmt.Errorf("keyring %s is nil", msg.Keyring)
	}

	if !keyring.IsAdmin(msg.Creator) {
		return nil, fmt.Errorf("initiator %s should be admin", msg.Creator)
	}

	if !keyring.IsAdmin(msg.Recipient) {
		return nil, fmt.Errorf("recipient %s should be admin", msg.Recipient)
	}

	if err = k.bankKeeper.SendCoins(
		ctx,
		sdk.MustAccAddressFromBech32(keyring.Address),
		sdk.MustAccAddressFromBech32(msg.Recipient),
		sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(msg.Amount))),
	); err != nil {
		return nil, err
	}
	return &types.MsgTransferFromKeyringResponse{}, nil
}
