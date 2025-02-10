package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) SignatureRequests(goCtx context.Context, req *types.QuerySignatureRequestsRequest) (*types.QuerySignatureRequestsResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	requests, pageRes, err := query.CollectionFilteredPaginate[uint64, types.SignRequest, collections.Map[uint64, types.SignRequest], *types.SignReqResponse](
		goCtx,
		k.SignRequestStore,
		req.Pagination,
		func(key uint64, value types.SignRequest) (bool, error) {
			if req.Status == types.SignRequestStatus_SIGN_REQUEST_STATUS_UNSPECIFIED || value.Status == req.Status {
				keyInfo, err := k.KeyStore.Get(goCtx, value.KeyIds[0])
				if err != nil {
					return false, nil
				}

				if req.KeyringAddr != "" && keyInfo.KeyringAddr != req.KeyringAddr {
					return false, nil
				}

				return true, nil
			}
			return false, nil
		},
		func(key uint64, value types.SignRequest) (*types.SignReqResponse, error) {
			return &types.SignReqResponse{
				Id:                     value.Id,
				Creator:                value.Creator,
				KeyIds:                 value.KeyIds,
				KeyType:                value.KeyType.String(),
				DataForSigning:         value.DataForSigning,
				Status:                 value.Status.String(),
				SignedData:             value.SignedData,
				KeyringPartySignatures: value.KeyringPartySignatures,
				RejectReason:           value.RejectReason,
				Metadata:               value.Metadata,
				CacheId:                value.CacheId,
				MpcBtl:                 value.MpcBtl,
				ParentReqId:            value.ParentReqId,
				ChildReqIds:            value.ChildReqIds,
				ZenbtcTxBytes:          value.ZenbtcTxBytes,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QuerySignatureRequestsResponse{SignRequests: requests, Pagination: pageRes}, nil
}
