package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) InitDctKeys(goCtx context.Context, msg *types.MsgInitDctKeys) (*types.MsgInitDctKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dct, err := k.DctStore.Get(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}
	if dct.Status != types.DctStatus_DCT_STATUS_KEYS_CREATED {
		return nil, types.ErrDctNotInRequestedState
	}

	signReqIds := make([]uint64, 0)

	// takes first unsigned tx for the nonce account
	nonceAccSignReqId, err := k.treasuryKeeper.InitDctNonceAccount(ctx, []uint64{dct.Solana.NonceAccountKey}, msg.UnsignedTx[0])
	if err != nil {
		return nil, err
	}

	signReqIds = append(signReqIds, nonceAccSignReqId)

	params, err := k.ParamStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	paramsSigner := params.Solana.SignerKeyId

	// takes second unsigned tx for the asset SPL creation
	// TODO: Ensure that the dct signer key is given mint authority
	assetSplSignReqId, err := k.treasuryKeeper.CreateAssetSpl(ctx, paramsSigner, msg.UnsignedTx[1])
	if err != nil {
		return nil, err
	}

	signReqIds = append(signReqIds, assetSplSignReqId)

	dct.Status = types.DctStatus_DCT_STATUS_SPL_REQUESTED
	if err := k.DctStore.Set(ctx, dct.Denom, dct); err != nil {
		return nil, err
	}

	// TODO: Check if SPL was actually created and update status to DCT_STATUS_COMPLETED
	// Needs to be done in a separate function based on Oracle Event reporting
	// That function also needs to make a first Bridge call to the destination chain with the escrowed dct funds

	return &types.MsgInitDctKeysResponse{
		SignReqIds: signReqIds,
	}, nil
}
