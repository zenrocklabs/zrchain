package keeper

import (
	"context"
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
)

func (k msgServer) MintRock(goCtx context.Context, req *types.MsgMintRock) (*types.MsgMintRockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	srcKey, err := k.treasuryKeeper.GetKey(ctx, req.SourceKeyId)
	if err != nil {
		return nil, err
	}

	dstKey, err := k.treasuryKeeper.GetKey(ctx, req.RecipientKeyId)
	if err != nil {
		return nil, err
	}

	chainType, err := treasurytypes.Caip2ToKeyType(req.DestinationChain)
	if err != nil {
		return nil, err
	}

	if dstKey.Type != chainType {
		return nil, fmt.Errorf("destination chain key type (%s) does not match recipient key chain type (%s)", chainType.String(), dstKey.PublicKey)
	}

	zenAddr, err := treasurytypes.NativeAddress(srcKey, "zen")
	if err != nil {
		return nil, err
	}

	bal := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(zenAddr), params.BondDenom)
	if bal.IsLT(sdk.NewCoin("urock", sdkmath.NewIntFromUint64(req.Amount))) {
		return nil, errors.New("not enough balance")
	}

	if !k.UserOwnsKey(goCtx, req.Creator, srcKey) {
		return nil, errors.New("creator does not own src key")
	}
	if !k.UserOwnsKey(goCtx, req.Creator, dstKey) {
		return nil, errors.New("creator does not own dst key")
	}

	return &types.MsgMintRockResponse{}, nil
}
