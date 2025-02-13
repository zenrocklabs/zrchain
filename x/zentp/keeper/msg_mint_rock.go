package keeper

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	var signerKey *treasurytypes.Key
	signerID := k.GetParams(goCtx).ZrchainRelayerKeyId
	mints, err := k.GetMints(ctx, req.RecipientKeyId, req.DestinationChain)
	if err != nil {
		return nil, err
	}
	if len(mints) > 0 {
		signerID = req.RecipientKeyId
		signerKey = dstKey
	} else {
		signerKey, err = k.treasuryKeeper.GetKey(ctx, signerID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find key for recipient key id %d", signerID)
		}
	}

	var hash string
	switch dstKey.Type {
	case treasurytypes.KeyType_KEY_TYPE_EDDSA_ED25519:
		hash, err = k.PrepareSolRockMintTx(goCtx, req.Amount, signerKey, dstKey)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", dstKey.Type)
	}
	if err != nil {
		return nil, err
	}

	k.mintStore.Set(ctx, []byte(hash), &types.Mint{
		Id:               0,
		WorkspaceAddr:    "",
		SourceKeyId:      0,
		DestinationChain: "",
		Amount:           0,
		RecipientKeyId:   0,
		TxHash:           "",
		State:            0,
	})
	return &types.MsgMintRockResponse{
		TxHash: hash,
	}, nil
}
