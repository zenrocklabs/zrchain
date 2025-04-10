package v4

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func ChangeKeyIdtoKeyIds(ctx sdk.Context, signTxReqCol collections.Map[uint64, types.SignTransactionRequest], codec codec.BinaryCodec) error {
	signTxReqStore, err := signTxReqCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	oldSignTxReqs, err := signTxReqStore.Values()
	if err != nil {
		return err
	}

	for _, signTxReq := range oldSignTxReqs {

		signTxReq.KeyIds = []uint64{signTxReq.KeyId}

		err = signTxReqCol.Set(ctx, signTxReq.Id, signTxReq)
		if err != nil {
			return err
		}
	}

	return nil
}
