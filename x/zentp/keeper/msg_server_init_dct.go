package keeper

import (
	"context"

	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) InitDct(goCtx context.Context, msg *types.MsgInitDct) (*types.MsgInitDctResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := k.DctStore.Get(ctx, msg.Asset.Denom)
	if err == nil {
		return nil, types.ErrDctAlreadyExists
	}

	if _, err := validationtypes.ValidateSolanaChainID(goCtx, msg.DestinationChain); err != nil {
		return nil, err
	}

	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(msg.Creator),
		types.ZentpDctCollectorName,
		sdk.NewCoins(*msg.Asset),
	); err != nil {
		return nil, err
	}

	keyIds, err := k.treasuryKeeper.CreateSolanaKeys(ctx)
	if err != nil {
		return nil, err
	}

	dct := types.Dct{
		Denom: msg.Asset.Denom,
		Solana: &types.Solana{
			SignerKeyId:       keyIds[0],
			NonceAccountKey:   keyIds[1],
			NonceAuthorityKey: keyIds[2],
			Btl:               20,
		},
		Status: types.DctStatus_DCT_STATUS_KEYS_REQUESTED,
	}

	err = k.DctStore.Set(ctx, msg.Asset.Denom, dct)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeInitDct,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Asset.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Asset.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDestinationChain, msg.DestinationChain),
		),
	})

	return &types.MsgInitDctResponse{Dct: &dct}, nil
}
