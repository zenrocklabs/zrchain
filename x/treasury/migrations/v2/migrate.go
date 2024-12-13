package v3

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func ChangeKeyIdtoKeyIds(ctx sdk.Context, signReqCol collections.Map[uint64, types.SignRequest], codec codec.BinaryCodec) error {
	signReqStore, err := signReqCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	oldSignReqs, err := signReqStore.Values()
	if err != nil {
		return err
	}
	for _, signReq := range oldSignReqs {

		signReq.KeyIds = []uint64{signReq.KeyId}

		err = signReqCol.Set(ctx, signReq.Id, signReq)
		if err != nil {
			return err
		}
	}

	return nil
}
