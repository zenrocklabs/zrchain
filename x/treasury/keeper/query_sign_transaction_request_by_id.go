package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func (k Keeper) SignTransactionRequestByID(goCtx context.Context, req *types.QuerySignTransactionRequestByIDRequest) (*types.QuerySignTransactionRequestByIDResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidArgument, "request is nil")
	}

	signTxReq, err := k.SignTransactionRequestStore.Get(goCtx, req.Id)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "sign transaction request %d not found", req.Id)
	}

	return &types.QuerySignTransactionRequestByIDResponse{
		SignTransactionRequest: &types.SignTxReqResponse{
			Id:                  signTxReq.Id,
			Creator:             signTxReq.Creator,
			KeyIds:              signTxReq.KeyIds,
			WalletType:          signTxReq.WalletType.String(),
			UnsignedTransaction: signTxReq.UnsignedTransaction,
			SignRequestId:       signTxReq.SignRequestId,
			NoBroadcast:         signTxReq.NoBroadcast,
		},
	}, nil
}
