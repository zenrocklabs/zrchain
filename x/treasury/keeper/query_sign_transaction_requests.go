package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func (k Keeper) SignTransactionRequests(goCtx context.Context, req *types.QuerySignTransactionRequestsRequest) (*types.QuerySignTransactionRequestsResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	requests, pageRes, err := query.CollectionFilteredPaginate[uint64, types.SignTransactionRequest, collections.Map[uint64, types.SignTransactionRequest], *types.SignTransactionRequestsResponse](
		goCtx,
		k.SignTransactionRequestStore,
		req.Pagination,
		func(key uint64, value types.SignTransactionRequest) (bool, error) {
			keyIDMatch := false
			if req.KeyId == 0 {
				keyIDMatch = true
			} else {
				for _, keyID := range value.KeyIds {
					if keyID == req.KeyId {
						keyIDMatch = true
						break
					}
				}
			}
			signReq, err := k.SignRequestStore.Get(goCtx, value.SignRequestId)
			if err != nil {
				return false, nil
			}
			statusMatch := req.Status == types.SignRequestStatus_SIGN_REQUEST_STATUS_UNSPECIFIED || signReq.Status == req.Status
			walletMatch := req.WalletType != types.WalletType_WALLET_TYPE_UNSPECIFIED

			return keyIDMatch && statusMatch && walletMatch, nil
		},
		func(key uint64, value types.SignTransactionRequest) (*types.SignTransactionRequestsResponse, error) {
			signReq, err := k.SignRequestStore.Get(goCtx, value.SignRequestId)
			if err != nil {
				return nil, err
			}
			return &types.SignTransactionRequestsResponse{
				SignTransactionRequests: &types.SignTxReqResponse{
					Id:                  value.Id,
					Creator:             value.Creator,
					KeyIds:              value.KeyIds,
					WalletType:          value.WalletType.String(),
					UnsignedTransaction: value.UnsignedTransaction,
					SignRequestId:       value.SignRequestId,
					NoBroadcast:         value.NoBroadcast,
				},
				SignRequest: &types.SignReqResponse{
					Id:                     signReq.Id,
					Creator:                signReq.Creator,
					KeyIds:                 signReq.KeyIds,
					KeyType:                signReq.KeyType.String(),
					DataForSigning:         signReq.DataForSigning,
					Status:                 signReq.Status.String(),
					SignedData:             signReq.SignedData,
					KeyringPartySignatures: signReq.KeyringPartySignatures,
					RejectReason:           signReq.RejectReason,
					Metadata:               signReq.Metadata,
					CacheId:                signReq.CacheId,
					ParentReqId:            signReq.ParentReqId,
					ChildReqIds:            signReq.ChildReqIds,
				},
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QuerySignTransactionRequestsResponse{SignTransactionRequests: requests, Pagination: pageRes}, nil
}
