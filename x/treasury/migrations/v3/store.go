package v3

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func RejectBadTestnetRequests(ctx sdk.Context, signRequestStore collections.Map[uint64, types.SignRequest], keyRequestStore collections.Map[uint64, types.KeyRequest], codec codec.BinaryCodec) error {
	// Process sign requests
	if err := signRequestStore.Walk(ctx, nil, func(id uint64, signReq types.SignRequest) (bool, error) {
		if signReq.Status == types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING ||
			signReq.Status == types.SignRequestStatus_SIGN_REQUEST_STATUS_PARTIAL {

			signReq.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
			signReq.RejectReason = "rejected by migration"
			if err := signRequestStore.Set(ctx, id, signReq); err != nil {
				return true, err
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	// Process key requests
	if err := keyRequestStore.Walk(ctx, nil, func(id uint64, keyReq types.KeyRequest) (bool, error) {
		if keyReq.Status == types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING ||
			keyReq.Status == types.KeyRequestStatus_KEY_REQUEST_STATUS_PARTIAL {

			keyReq.Status = types.KeyRequestStatus_KEY_REQUEST_STATUS_REJECTED
			keyReq.RejectReason = "rejected by migration"
			if err := keyRequestStore.Set(ctx, id, keyReq); err != nil {
				return true, err
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	return nil
}
