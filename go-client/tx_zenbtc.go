package client

import (
	"context"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// ZenBTCTxClient provides a client interface for interacting with the ZenBTC module.
// It wraps the RawTxClient to provide specialized methods for ZenBTC-related transactions.
type ZenBTCTxClient struct {
	c *RawTxClient
}

// NewZenBTCTxClient returns a new ZenBTCTxClient instance.
//
// Parameters:
//   - c: A RawTxClient for handling low-level transaction operations
//
// Returns:
//   - *ZenBTCTxClient: A new ZenBTC transaction client instance
func NewZenBTCTxClient(c *RawTxClient) *ZenBTCTxClient {
	return &ZenBTCTxClient{c: c}
}

// NewVerifyDepositBlockInclusion creates a transaction to verify a Bitcoin deposit's inclusion in a block.
// This method is used to prove that a Bitcoin transaction exists and is properly included in the Bitcoin blockchain.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - chainName: The name of the Bitcoin chain (e.g., "mainnet", "testnet")
//   - blockHeight: The height of the block containing the deposit
//   - rawTX: The raw Bitcoin transaction data
//   - index: The index of the transaction in the block
//   - proof: Merkle proof demonstrating the transaction's inclusion
//   - depositAddr: The Bitcoin address where the deposit was made
//   - amount: The amount of Bitcoin deposited (in satoshis)
//
// Returns:
//   - string: The transaction hash if verification is successful
//   - error: An error if verification fails or transaction submission fails
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

func (c *ZenBTCTxClient) NewSubmitUnsignedRedemptionTx(ctx context.Context, hashkeys []*types.InputHashes, txBytes []byte, cacheID []byte, chainName string) (string, error) {
	msg := &types.MsgSubmitUnsignedRedemptionTx{
		Creator:   c.c.Identity.Address.String(),
		Inputs:    hashkeys,
		Txbytes:   txBytes,
		CacheId:   cacheID,
		ChainName: chainName,
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
