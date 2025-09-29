package client

import (
	"context"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
)

// ValidationTxClient provides helper methods for validation module transactions.
type ValidationTxClient struct {
	c *RawTxClient
}

// NewValidationTxClient creates a new ValidationTxClient instance.
func NewValidationTxClient(c *RawTxClient) *ValidationTxClient {
	return &ValidationTxClient{c: c}
}

// RequestHeaderBackfill enqueues a Bitcoin header height for backfill through the normal ABCI flow.
func (c *ValidationTxClient) RequestHeaderBackfill(ctx context.Context, height int64) (string, error) {
	if height <= 0 {
		return "", fmt.Errorf("height must be greater than zero")
	}

	msg := &types.MsgRequestHeaderBackfill{
		Authority: c.c.Identity.Address.String(),
		Height:    height,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, DefaultGasLimit, DefaultFees, msg)
	if err != nil {
		return "", err
	}

	return c.c.SendWaitTx(ctx, txBytes)
}

// ManuallyInputBitcoinHeader injects the provided Bitcoin header directly into consensus state.
func (c *ValidationTxClient) ManuallyInputBitcoinHeader(ctx context.Context, header types.BitcoinHeader) (string, error) {
	if header.BlockHeight <= 0 {
		return "", fmt.Errorf("header block height must be greater than zero")
	}
	if header.MerkleRoot == "" {
		return "", fmt.Errorf("header merkle root must be provided")
	}
	if header.BlockHash == "" {
		return "", fmt.Errorf("header block hash must be provided")
	}

	msg := &types.MsgManuallyInputBitcoinHeader{
		Authority: c.c.Identity.Address.String(),
		Header:    header,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, DefaultGasLimit, DefaultFees, msg)
	if err != nil {
		return "", err
	}

	return c.c.SendWaitTx(ctx, txBytes)
}
