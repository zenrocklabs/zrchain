package client

import (
	"context"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankTxClient handles bank-related transactions
type BankTxClient struct {
	c *RawTxClient
}

// NewBankTxClient creates a new bank transaction client
func NewBankTxClient(c *RawTxClient) *BankTxClient {
	return &BankTxClient{c: c}
}

// SendCoins sends coins from one address to another
//
// Parameters:
//   - ctx: Context for the operation
//   - fromAddress: The sender's address (must be the identity's address)
//   - toAddress: The recipient's address
//   - amount: The amount to send (in urock)
//   - denom: The denomination (default: "urock")
//
// Returns:
//   - string: Transaction hash
//   - error: An error if the transaction fails
//
// Example:
//
//	hash, err := bankClient.SendCoins(ctx, "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty", "zen1qwnafe2s9eawhah5x6v4593v3tljdntl9zcqpn", 100000000000, "urock")
//	if err != nil {
//	    // Handle error
//	}
func (c *BankTxClient) SendCoins(ctx context.Context, fromAddress, toAddress string, amount uint64, denom string) (string, error) {
	// Convert string addresses to AccAddress
	fromAccAddr, err := types.AccAddressFromBech32(fromAddress)
	if err != nil {
		return "", err
	}
	toAccAddr, err := types.AccAddressFromBech32(toAddress)
	if err != nil {
		return "", err
	}
	// Create the MsgSend message
	msg := banktypes.NewMsgSend(
		fromAccAddr,
		toAccAddr,
		types.NewCoins(types.NewCoin(denom, math.NewIntFromUint64(amount))),
	)
	// Build and sign the transaction
	txBytes, err := c.c.BuildAndSignTx(ctx, DefaultGasLimit, DefaultFees, msg)
	if err != nil {
		return "", err
	}
	// Send the transaction and wait for it to be included
	hash, err := c.c.SendWaitTx(ctx, txBytes)
	if err != nil {
		return "", err
	}
	return hash, nil
}
