package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
)

type ZenTPTxClient struct {
	c *RawTxClient
}

func NewZenTPTxClient(c *RawTxClient) *ZenTPTxClient {
	return &ZenTPTxClient{c: c}
}

func (c *ZenTPTxClient) NewMsgBridge(
	ctx context.Context,
	sourceAddress string,
	amount uint64,
	denom string,
	destinationChain string,
	recipientAddress string,
) (string, error) {
	msg := &types.MsgBridge{
		Creator:          c.c.Identity.Address.String(),
		SourceAddress:    sourceAddress,
		Amount:           amount,
		Denom:            denom,
		DestinationChain: destinationChain,
		RecipientAddress: recipientAddress,
	}
	txBytes, err := c.c.BuildAndSignTx(ctx, ZenBTCGasLimit, ZenBTCDefaultFees, msg)
	if err != nil {
		return "", err
	}

	hash, err := c.c.SendWaitTx(ctx, txBytes)
	if err != nil {
		return "", err
	}

	return hash, nil
}
