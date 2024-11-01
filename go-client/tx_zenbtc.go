package client

import (
	"context"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

type ZenBTCTxClient struct {
	c *RawTxClient
}

func NewZenBTCTxClient(c *RawTxClient) *ZenBTCTxClient {
	return &ZenBTCTxClient{c: c}
}

func (c *ZenBTCTxClient) NewVerifyDepositBlockInclusion(
	ctx context.Context, chainName string, blockHeight int64, rawTX string, index int32, proof []string, depositAddr string, amount uint64,
) (string, error) {
	msg := &types.MsgVerifyDepositBlockInclusion{
		Creator:     c.c.Identity.Address.String(),
		ChainName:   chainName,
		BlockHeight: blockHeight,
		RawTx:       rawTX,
		Index:       index,
		Proof:       proof,
		DepositAddr: depositAddr,
		Amount:      amount,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, DefaultGasLimit, DefaultFees, msg)
	if err != nil {
		return "", err
	}

	hash, err := c.c.SendWaitTx(ctx, txBytes)
	if err != nil {
		return "", err
	}

	return hash, nil
}
