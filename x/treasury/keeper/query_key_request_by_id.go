package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) KeyRequestByID(goCtx context.Context, req *types.QueryKeyRequestByIDRequest) (*types.QueryKeyRequestByIDResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	keyReq, err := k.KeyRequestStore.Get(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("key request %d not found", req.Id)
	}

	return &types.QueryKeyRequestByIDResponse{
		KeyRequest: &types.KeyReqResponse{
			Id:                     keyReq.Id,
			Creator:                keyReq.Creator,
			WorkspaceAddr:          keyReq.WorkspaceAddr,
			KeyringAddr:            keyReq.KeyringAddr,
			KeyType:                keyReq.KeyType.String(),
			Status:                 keyReq.Status.String(),
			KeyringPartySignatures: keyReq.KeyringPartySignatures,
			RejectReason:           keyReq.RejectReason,
			Index:                  keyReq.Index,
			SignPolicyId:           keyReq.SignPolicyId,
			ZenbtcMetadata:         keyReq.ZenbtcMetadata,
			MpcBtl:                 keyReq.MpcBtl,
			PublicKey:              keyReq.PublicKey,
		},
	}, nil
}
