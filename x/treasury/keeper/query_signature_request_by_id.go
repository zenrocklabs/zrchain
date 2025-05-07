package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func (k Keeper) SignatureRequestByID(goCtx context.Context, req *types.QuerySignatureRequestByIDRequest) (*types.QuerySignatureRequestByIDResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidArgument, "request is nil")
	}

	signReq, err := k.SignRequestStore.Get(goCtx, req.Id)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrNotFound, "signature request %d not found", req.Id)
	}

	return &types.QuerySignatureRequestByIDResponse{
		SignRequest: &types.SignReqResponse{
			Id:               signReq.Id,
			Creator:          signReq.Creator,
			KeyIds:           signReq.KeyIds,
			KeyType:          signReq.KeyType.String(),
			DataForSigning:   signReq.DataForSigning,
			Status:           signReq.Status.String(),
			SignedData:       signReq.SignedData,
			KeyringPartySigs: signReq.KeyringPartySigs,
			RejectReason:     signReq.RejectReason,
			Metadata:         signReq.Metadata,
			CacheId:          signReq.CacheId,
			MpcBtl:           signReq.MpcBtl,
			ParentReqId:      signReq.ParentReqId,
			ChildReqIds:      signReq.ChildReqIds,
			ZenbtcTxBytes:    signReq.ZenbtcTxBytes,
		},
	}, nil
}
