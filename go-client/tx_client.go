package client

import (
	"fmt"

	"google.golang.org/grpc"
)

// Package client provides transaction handling capabilities for the Zenrock blockchain.
// This file specifically implements the transaction client interfaces for sending
// transactions to zenrockd and interacting with various modules.

// TxClient provides a unified interface for sending transactions to the Zenrock blockchain.
// It aggregates multiple specialized transaction clients to handle different types of
// transactions across various modules.
type TxClient struct {
	*RawTxClient      // Handles low-level transaction construction and signing
	*TreasuryTxClient // Handles treasury-specific transactions
	*ZenBTCTxClient   // Handles ZenBTC-specific transactions
	*ZenTPTxClient    // Handles ZenTP-specific transactions
	*ZenexTxClient    // Handles ZenEX-specific transactions
}

// NewTxClient creates a new transaction client instance with all necessary sub-clients initialized.
//
// Parameters:
//   - id: The Identity used for signing transactions
//   - chainID: The blockchain network identifier
//   - c: A gRPC connection to a Zenrock node
//   - accountFetcher: Interface for retrieving account information
//
// Returns:
//   - *TxClient: A new transaction client instance
//   - error: An error if initialization fails
//
// Example:
//
//	id := NewIdentity(...)
//	conn, _ := grpc.Dial(...)
//	client, err := NewTxClient(id, "zenrock_123-1", conn, accountFetcher)
//	if err != nil {
//	    // Handle error
//	}
func NewTxClient(id Identity, chainID string, c *grpc.ClientConn, accountFetcher AccountFetcher) (*TxClient, error) {
	raw, err := NewRawTxClient(id, chainID, c, accountFetcher)
	if err != nil {
		return nil, fmt.Errorf("can't create raw tx client: %w", err)
	}
	return &TxClient{
		RawTxClient:      raw,
		TreasuryTxClient: NewTreasuryTxClient(raw),
		ZenBTCTxClient:   NewZenBTCTxClient(raw),
		ZenTPTxClient:    NewZenTPTxClient(raw),
		ZenexTxClient:    NewZenexTxClient(raw),
	}, nil
}
