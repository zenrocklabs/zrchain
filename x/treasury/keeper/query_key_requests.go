package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	"github.com/Zenrock-Foundation/zrchain/v4/x/treasury/types"

	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) KeyRequests(goCtx context.Context, req *types.QueryKeyRequestsRequest) (*types.QueryKeyRequestsResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	requests, pageRes, err := query.CollectionFilteredPaginate[uint64, types.KeyRequest, collections.Map[uint64, types.KeyRequest], *types.KeyReqResponse](
		goCtx,
		k.KeyRequestStore,
		req.Pagination,
		func(key uint64, value types.KeyRequest) (bool, error) {
			keyringAddrMatches := req.KeyringAddr == "" || value.KeyringAddr == req.KeyringAddr
			statusMatches := req.Status == types.KeyRequestStatus_KEY_REQUEST_STATUS_UNSPECIFIED || value.Status == req.Status
			workspaceAddMatches := req.WorkspaceAddr == "" || value.WorkspaceAddr == req.WorkspaceAddr

			return keyringAddrMatches && statusMatches && workspaceAddMatches, nil
		},
		func(key uint64, value types.KeyRequest) (*types.KeyReqResponse, error) {
			return &types.KeyReqResponse{
				Id:                     value.Id,
				Creator:                value.Creator,
				WorkspaceAddr:          value.WorkspaceAddr,
				KeyringAddr:            value.KeyringAddr,
				KeyType:                value.KeyType.String(),
				Status:                 value.Status.String(),
				KeyringPartySignatures: value.KeyringPartySignatures,
				RejectReason:           value.RejectReason,
				Index:                  value.Index,
				SignPolicyId:           value.SignPolicyId,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryKeyRequestsResponse{
		KeyRequests: requests,
		Pagination:  pageRes,
	}, nil
}
