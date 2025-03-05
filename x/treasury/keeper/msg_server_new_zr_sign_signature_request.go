package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"cosmossdk.io/errors"
	cosmos_types "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewZrSignSignatureRequest(goCtx context.Context, msg *types.MsgNewZrSignSignatureRequest) (*types.MsgNewZrSignSignatureRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store, err := k.ParamStore.Get(ctx)
	if err != nil {
		return nil, err
	}
	if msg.Creator != store.ZrSignAddress {
		return nil, fmt.Errorf("creator is not zrsign")
	}

	walletTypeStr := strconv.FormatUint(msg.KeyType, 10)
	ws, err := k.identityKeeper.GetZrSignWorkspaces(goCtx, msg.Address, walletTypeStr)
	if err != nil {
		return nil, errors.Wrap(err, "GetZrSignWorkspaces")
	}
	if len(ws) != 1 {
		return nil, fmt.Errorf("wallet type workspace not found: %s", walletTypeStr)
	}
	wsID, _ := ws[walletTypeStr]

	keys, _, err := query.CollectionFilteredPaginate(
		goCtx,
		k.KeyStore,
		nil,
		func(key uint64, value types.Key) (bool, error) {
			return value.WorkspaceAddr == wsID && value.Index == msg.WalletIndex, nil
		},
		func(key uint64, value types.Key) (types.Key, error) {
			return value, nil
		})
	if err != nil {
		return nil, err
	}

	if len(keys) != 1 {
		return nil, fmt.Errorf("key index not found: %s", walletTypeStr)
	}
	key := keys[0]

	response := &types.MsgNewZrSignSignatureRequestResponse{}

	if msg.Tx {
		tx, err := k.zrSignParseTx(msg.WalletType, &key, msg.Data, msg.Metadata)
		if err != nil {
			return nil, err
		}
		byteData, err := hex.DecodeString(msg.Data)
		if err != nil {
			return nil, err
		}
		resp, err := k.HandleSignTransactionRequest(ctx, &types.MsgNewSignTransactionRequest{
			Creator:             msg.Creator,
			KeyIds:              []uint64{key.Id},
			WalletType:          msg.WalletType,
			UnsignedTransaction: byteData,
			Metadata:            msg.Metadata,
			Btl:                 msg.Btl,
			CacheId:             msg.CacheId,
			NoBroadcast:         msg.NoBroadcast,
		}, tx.DataForSigning)
		if err != nil {
			return nil, errors.Wrap(err, "signatureRequest")
		}
		response.ReqId = resp.Id
	} else {
		resp, err := k.HandleSignatureRequest(ctx, &types.MsgNewSignatureRequest{
			Creator:                  msg.Creator,
			KeyIds:                   []uint64{key.Id},
			DataForSigning:           msg.Data,
			Btl:                      msg.Btl,
			CacheId:                  msg.CacheId,
			VerifySigningData:        msg.VerifySigningData,
			VerifySigningDataVersion: msg.VerifySigningDataVersion,
		})
		if err != nil {
			return nil, errors.Wrap(err, "signatureRequest")
		}
		response.ReqId = resp.SigReqId
	}

	return response, nil
}

func (k msgServer) zrSignParseTx(walletType types.WalletType, key *types.Key, data string, metadata *cosmos_types.Any) (*types.Transfer, error) {
	// use wallet to parse unsigned transaction
	w, err := types.NewWallet(key, walletType)
	if err != nil {
		return nil, err
	}

	parser, ok := w.(types.TxParser)
	if !ok {
		return nil, fmt.Errorf("wallet %v does not implement TxParser", walletType)
	}

	var meta types.Metadata
	if err := k.cdc.UnpackAny(metadata, &meta); err != nil {
		return nil, fmt.Errorf("failed to unpack metadata: %w", err)
	}

	byteData, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	tx, err := parser.ParseTx(byteData, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tx: %w", err)
	}

	return &tx, nil
}
