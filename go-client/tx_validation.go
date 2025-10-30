package client

import (
	"context"
	"fmt"
	"strings"

	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
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

	txBytes, err := c.c.BuildAndSignTx(ctx, InjectHeaderGasLimit, DefaultFees, msg)
	if err != nil {
		return "", err
	}

	return c.c.SendWaitTx(ctx, txBytes)
}

// ManuallyInputZcashHeader injects the provided Zcash header directly into consensus state.
func (c *ValidationTxClient) ManuallyInputZcashHeader(ctx context.Context, header types.ZcashHeader) (string, error) {
	if header.BlockHeight <= 0 {
		return "", fmt.Errorf("header block height must be greater than zero")
	}
	if header.MerkleRoot == "" {
		return "", fmt.Errorf("header merkle root must be provided")
	}
	if header.BlockHash == "" {
		return "", fmt.Errorf("header block hash must be provided")
	}

	msg := &types.MsgManuallyInputZcashHeader{
		Authority: c.c.Identity.Address.String(),
		Header:    header,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, InjectHeaderGasLimit, DefaultFees, msg)
	if err != nil {
		return "", err
	}

	return c.c.SendWaitTx(ctx, txBytes)
}

// AdvanceSolanaNonce submits a maintenance transaction that advances a Solana durable
// nonce account using a provided recent blockhash. Either zenBTC must be true or a DCT
// asset must be specified.
func (c *ValidationTxClient) AdvanceSolanaNonce(ctx context.Context, recentBlockhash, caip2ChainID string, zenBTC bool, asset dcttypes.Asset) (string, error) {
	if recentBlockhash == "" {
		return "", fmt.Errorf("recent blockhash must be provided")
	}
	if caip2ChainID == "" {
		return "", fmt.Errorf("caip2 chain id must be provided")
	}
	if !zenBTC {
		if asset == dcttypes.Asset_ASSET_UNSPECIFIED {
			return "", fmt.Errorf("asset must be provided when zenBTC flag is false")
		}
		if _, ok := dcttypes.Asset_name[int32(asset)]; !ok {
			return "", fmt.Errorf("invalid asset: %s", asset.String())
		}
	}

	msg := &types.MsgAdvanceSolanaNonce{
		Authority:       c.c.Identity.Address.String(),
		Zenbtc:          zenBTC,
		Asset:           asset,
		RecentBlockhash: recentBlockhash,
		Caip2ChainId:    caip2ChainID,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, DefaultGasLimit, DefaultFees, msg)
	if err != nil {
		// Provide context if the recent blockhash appears incorrectly formatted.
		if strings.Contains(err.Error(), "blockhash") {
			return "", fmt.Errorf("failed to build advance nonce tx (blockhash=%s): %w", recentBlockhash, err)
		}
		return "", err
	}

	return c.c.SendWaitTx(ctx, txBytes)
}
