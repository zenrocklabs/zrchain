package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

type ZenexTxClient struct {
	c *RawTxClient
}

func NewZenexTxClient(c *RawTxClient) *ZenexTxClient {
	return &ZenexTxClient{c: c}
}

func (c *ZenTPTxClient) NewMsgSwap(
	ctx context.Context,
	creator string,
	amountIn uint64,
	rockKeyID uint64,
	btcKeyID uint64,
	pair string,
	workspace string,
) (string, error) {
	msg := &types.MsgSwapRequest{
		Creator:   creator,
		Pair:      pair,
		Workspace: workspace,
		AmountIn:  amountIn,
		RockKeyId: rockKeyID,
		BtcKeyId:  btcKeyID,
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
