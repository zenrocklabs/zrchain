package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

// DCTTxClient provides a client interface for interacting with the DCT module.
// It wraps the RawTxClient to provide specialized methods for DCT-related transactions
// across multiple assets.
type DCTTxClient struct {
	c *RawTxClient
}

// NewDCTTxClient returns a new DCTTxClient instance.
//
// Parameters:
//   - c: A RawTxClient for handling low-level transaction operations
//
// Returns:
//   - *DCTTxClient: A new DCT transaction client instance
func NewDCTTxClient(c *RawTxClient) *DCTTxClient {
	return &DCTTxClient{c: c}
}

// NewVerifyDepositBlockInclusion creates a transaction to verify a deposit's inclusion in a block.
// This method is used to prove that a transaction exists and is properly included in the blockchain.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//   - chainName: The name of the chain (e.g., "mainnet", "testnet")
//   - blockHeight: The height of the block containing the deposit
//   - rawTX: The raw transaction data
//   - index: The index of the transaction in the block
//   - vout: The output index in the transaction
//   - proof: Merkle proof demonstrating the transaction's inclusion
//   - depositAddr: The address where the deposit was made
//   - amount: The amount deposited (in the smallest unit, e.g., satoshis)
//
// Returns:
//   - string: The transaction hash if verification is successful
//   - error: An error if verification fails or transaction submission fails
func (c *DCTTxClient) NewVerifyDepositBlockInclusion(
	ctx context.Context,
	asset types.Asset,
	chainName string,
	blockHeight int64,
	rawTX string,
	index int32,
	vout uint64,
	proof []string,
	depositAddr string,
	amount uint64,
) (string, error) {
	msg := &types.MsgVerifyDepositBlockInclusion{
		Creator:     c.c.Identity.Address.String(),
		Asset:       asset,
		ChainName:   chainName,
		BlockHeight: blockHeight,
		RawTx:       rawTX,
		Index:       index,
		Proof:       proof,
		DepositAddr: depositAddr,
		Amount:      amount,
		Vout:        vout,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, DCTGasLimit, DCTDefaultFees, msg)
	if err != nil {
		return "", err
	}

	hash, err := c.c.SendWaitTx(ctx, txBytes)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// NewSubmitUnsignedRedemptionTx submits an unsigned redemption transaction for a specific asset.
//
// Parameters:
//   - ctx: Context for the request
//   - asset: The asset type (e.g., ASSET_ZENBTC, ASSET_ZENZEC)
//   - hashkeys: Array of input hashes with their key IDs
//   - txBytes: The unsigned transaction bytes
//   - cacheID: Cache identifier for the transaction
//   - chainName: The name of the chain
//   - redemptionIndexes: Indexes of the redemptions being processed
//
// Returns:
//   - string: The transaction hash if submission is successful
//   - error: An error if submission fails
func (c *DCTTxClient) NewSubmitUnsignedRedemptionTx(
	ctx context.Context,
	asset types.Asset,
	hashkeys []*types.InputHashes,
	txBytes []byte,
	cacheID []byte,
	chainName string,
	redemptionIndexes []uint64,
) (string, error) {
	msg := &types.MsgSubmitUnsignedRedemptionTx{
		Creator:           c.c.Identity.Address.String(),
		Asset:             asset,
		Inputs:            hashkeys,
		Txbytes:           txBytes,
		CacheId:           cacheID,
		ChainName:         chainName,
		RedemptionIndexes: redemptionIndexes,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, DCTGasLimit, DCTDefaultFees, msg)
	if err != nil {
		return "", err
	}

	hash, err := c.c.SendWaitTx(ctx, txBytes)
	if err != nil {
		return "", err
	}

	return hash, nil
}
