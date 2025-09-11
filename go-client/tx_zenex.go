package client

import (
	"context"

	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
)

type ZenexTxClient struct {
	c *RawTxClient
}

func NewZenexTxClient(c *RawTxClient) *ZenexTxClient {
	return &ZenexTxClient{c: c}
}

func (c *ZenexTxClient) NewMsgSwapRequest(
	ctx context.Context,
	amountIn uint64,
	rockKeyID uint64,
	btcKeyID uint64,
	pair types.TradePair,
	workspace string,
) (string, error) {
	msg := &types.MsgSwapRequest{
		Creator:   c.c.Identity.Address.String(),
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

func (c *ZenexTxClient) NewMsgZenexTransferRequest(
	ctx context.Context,
	creator string,
	swapId uint64,
	unsignedTx []byte,
	walletType treasurytypes.WalletType,
) (string, error) {
	msg := &types.MsgZenexTransferRequest{
		Creator:    c.c.Identity.Address.String(),
		SwapId:     swapId,
		UnsignedTx: unsignedTx,
		WalletType: walletType,
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

func (c *ZenexTxClient) NewMsgAcknowledgePoolTransfer(
	ctx context.Context,
	swapId uint64,
	sourceTxHash string,
	status types.SwapStatus,
) (string, error) {
	msg := &types.MsgAcknowledgePoolTransfer{
		Creator:      c.c.Identity.Address.String(),
		SwapId:       swapId,
		SourceTxHash: sourceTxHash,
		Status:       status,
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
