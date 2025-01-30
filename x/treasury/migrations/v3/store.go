package v3

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func RejectBadTestnetRequests(ctx sdk.Context, signRequestStore collections.Map[uint64, types.SignRequest], codec codec.BinaryCodec) error {
	if err := signRequestStore.Walk(ctx, nil, func(id uint64, signReq types.SignRequest) (bool, error) {
		if signReq.DataForSigning == nil || signReq.DataForSigning[0] == nil {
			signReq.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			signReq.RejectReason = "data for signing is empty"
			if err := signRequestStore.Set(ctx, id, signReq); err != nil {
				return true, err
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	return nil
}
