package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"

	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	policytypes "github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

func (k msgServer) NewICATransactionRequest(goCtx context.Context, msg *types.MsgNewICATransactionRequest) (*types.MsgNewICATransactionRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	key, err := k.KeyStore.Get(ctx, msg.KeyId)
	if err != nil {
		return nil, fmt.Errorf("key %v not found", msg.KeyId)
	}

	if msg.InputPayload == "" {
		return nil, fmt.Errorf("no payload provided")
	}

	signPolicyId := key.SignPolicyId

	ws, err := k.identityKeeper.WorkspaceStore.Get(ctx, key.WorkspaceAddr)
	if err != nil {
		return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
	}

	if signPolicyId == 0 {
		ws, err := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr)
		if err != nil {
			return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
		}
		signPolicyId = ws.SignPolicyId
	}

	keyring, err := k.identityKeeper.GetKeyring(ctx, key.KeyringAddr)
	if err != nil || !keyring.IsActive {
		return nil, fmt.Errorf("keyring %s is nil or is inactive", keyring.Address)
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, signPolicyId, msg.Btl, nil, ws.Owners)
	if err != nil {
		return nil, err
	}
	return k.NewICATransactionRequestActionHandler(ctx, act)
}

func (k msgServer) NewICATransactionRequestActionHandler(ctx sdk.Context, act *policytypes.Action) (*types.MsgNewICATransactionRequestResponse, error) {
	return policykeeper.TryExecuteAction(
		k.policyKeeper,
		k.cdc,
		ctx,
		act,
		func(ctx sdk.Context, msg *types.MsgNewICATransactionRequest) (*types.MsgNewICATransactionRequestResponse, error) {
			key, err := k.KeyStore.Get(ctx, msg.KeyId)
			if err != nil {
				return nil, fmt.Errorf("key %v not found", msg.KeyId)
			}

			if _, err := k.identityKeeper.GetWorkspace(ctx, key.WorkspaceAddr); err != nil {
				return nil, fmt.Errorf("workspace %s not found", key.WorkspaceAddr)
			}

			cdc := codec.NewProtoCodec(types.ModuleCdc.InterfaceRegistry())

			// TODO - add support for different encodings and memo
			packetData, err := generatePacketData(cdc, []byte(msg.InputPayload), "", "proto3")
			if err != nil {
				return nil, err
			}

			icaMsg := icacontrollertypes.NewMsgSendTx(
				msg.Creator,
				msg.ConnectionId,
				msg.RelativeTimeoutTimestamp,
				*packetData,
			)

			msgBytes, err := cdc.Marshal(icaMsg)
			if err != nil {
				return nil, err
			}

			request := &types.ICATransactionRequest{
				Creator:  msg.Creator,
				KeyId:    msg.KeyId,
				InputMsg: msgBytes,
			}

			id, err := k.CreateICATransactionRequest(ctx, request)
			if err != nil {
				return nil, err
			}

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventNewICATransactionRequest,
					sdk.NewAttribute(types.AttributeRequestId, strconv.FormatUint(id, 10)),
				),
			})

			return &types.MsgNewICATransactionRequestResponse{}, nil
		},
	)
}

// generatePacketData takes in message bytes and a memo and serializes the message into an
// instance of InterchainAccountPacketData.
func generatePacketData(cdc *codec.ProtoCodec, msgBytes []byte, memo string, encoding string) (*icatypes.InterchainAccountPacketData, error) {
	protoMessages, err := convertBytesIntoProtoMessages(cdc, msgBytes)
	if err != nil {
		return nil, err
	}

	return generateIcaPacketDataFromProtoMessages(cdc, protoMessages, memo, encoding)
}

// convertBytesIntoProtoMessages returns a list of proto messages from bytes. The bytes can be in the form of a single
// message, or a json array of messages.
func convertBytesIntoProtoMessages(cdc *codec.ProtoCodec, msgBytes []byte) ([]proto.Message, error) {
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(msgBytes, &rawMessages); err != nil {
		// if we fail to unmarshal a list of messages, we assume we are just dealing with a single message.
		// in this case we return a list of a single item.
		var msg sdk.Msg
		if err := cdc.UnmarshalInterfaceJSON(msgBytes, &msg); err != nil {
			return nil, err
		}

		return []proto.Message{msg}, nil
	}

	sdkMessages := make([]proto.Message, len(rawMessages))
	for i, anyJSON := range rawMessages {
		var msg sdk.Msg
		if err := cdc.UnmarshalInterfaceJSON(anyJSON, &msg); err != nil {
			return nil, err
		}

		sdkMessages[i] = msg
	}

	return sdkMessages, nil
}

// generateIcaPacketDataFromProtoMessages generates ica packet data from a given set of proto encoded sdk messages and a memo.
func generateIcaPacketDataFromProtoMessages(cdc *codec.ProtoCodec, sdkMessages []proto.Message, memo string, encoding string) (*icatypes.InterchainAccountPacketData, error) {
	icaPacketDataBytes, err := icatypes.SerializeCosmosTx(cdc, sdkMessages, encoding)
	if err != nil {
		return nil, err
	}

	icaPacketData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: icaPacketDataBytes,
		Memo: memo,
	}

	if err := icaPacketData.ValidateBasic(); err != nil {
		return nil, err
	}

	// return cdc.MarshalJSON(&icaPacketData)
	return &icaPacketData, nil
}
