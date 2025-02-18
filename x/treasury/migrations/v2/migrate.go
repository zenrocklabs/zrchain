package v2

import (
	"fmt"

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

func RejectBadTestnetRequests(ctx sdk.Context, signRequestStore collections.Map[uint64, types.SignRequest], codec codec.BinaryCodec) error {
	if err := signRequestStore.Walk(ctx, nil, func(id uint64, signReq types.SignRequest) (bool, error) {
		if len(signReq.DataForSigning) == 0 || len(signReq.DataForSigning[0]) == 0 {
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

func ChangeZenBtcMetadataChainIdtoCaip2Id(ctx sdk.Context, keyCol collections.Map[uint64, types.Key], keyReqCol collections.Map[uint64, types.KeyRequest], codec codec.BinaryCodec) error {
	keyReqStore, err := keyReqCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}

	oldKeyReqs, err := keyReqStore.Values()
	if err != nil {
		return err
	}

	for _, keyReq := range oldKeyReqs {

		if keyReq.ZenbtcMetadata != nil {

			switch keyReq.ZenbtcMetadata.ChainType {
			case types.WalletType_WALLET_TYPE_EVM:
				keyReq.ZenbtcMetadata.Caip2ChainId = fmt.Sprintf("eip155:%d", keyReq.ZenbtcMetadata.ChainId)
			}

			err = keyReqCol.Set(ctx, keyReq.Id, keyReq)
			if err != nil {
				return err
			}
		}

	}

	keyStore, err := keyCol.Iterate(ctx, nil)
	if err != nil {
		return err
	}

	oldKeys, err := keyStore.Values()
	if err != nil {
		return err
	}

	for _, key := range oldKeys {

		if key.ZenbtcMetadata != nil {

			switch key.ZenbtcMetadata.ChainType {
			case types.WalletType_WALLET_TYPE_EVM:
				key.ZenbtcMetadata.Caip2ChainId = fmt.Sprintf("eip155:%d", key.ZenbtcMetadata.ChainId)
			}

			err = keyCol.Set(ctx, key.Id, key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
