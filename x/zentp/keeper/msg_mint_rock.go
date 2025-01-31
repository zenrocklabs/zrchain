package keeper

import (
	"context"
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
)

func (k msgServer) MintRock(goCtx context.Context, req *types.MsgMintRock) (*types.MsgMintRockResponse, error) {
	if k.GetAuthority() != req.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	key, err := k.treasuryKeeper.GetKey(ctx, req.SourceKeyId)
	if err != nil {
		return nil, err
	}

	zenAddr, err := treasurytypes.NativeAddress(key, "zen")
	if err != nil {
		return nil, err
	}

	// TODO: use bond denom
	bal := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(zenAddr), "urock")
	if bal.IsLT(sdk.NewCoin("urock", sdkmath.NewIntFromUint64(req.Amount))) {
		return nil, errors.New("not enough balance")
	}
	return &types.MsgMintRockResponse{}, nil
}
