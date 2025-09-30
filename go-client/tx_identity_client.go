package client

import (
	"context"

	identitytypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
)

// Package client provides transaction handling capabilities for the Zenrock identity module.
// This file implements methods for managing key requests, signature requests, and other
// treasury-related transactions.

// IdentityTxClient provides methods for interacting with the identity module.
// It wraps a RawTxClient to handle treasury-specific transactions.
type IdentityTxClient struct {
	c *RawTxClient
}

// NewIdentityTxClient creates a new instance of IdentityTxClient.
//
// Parameters:
//   - c: A RawTxClient for handling low-level transaction operations
//
// Returns:
//   - *IdentityTxClient: A new identity transaction client instance
func NewIdentityTxClient(c *RawTxClient) *IdentityTxClient {
	return &IdentityTxClient{c: c}
}

func (c *IdentityTxClient) NewWorkspace(ctx context.Context) (string, error) {
	msg := identitytypes.NewMsgNewWorkspace(c.c.Identity.Address.String(), 0, 0)
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
