package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

// Package client provides transaction handling capabilities for Inter-Chain Account (ICA) operations
// in the Zenrock blockchain. This file implements methods for fulfilling and rejecting
// ICA transaction requests through the treasury module.

// FulfilICATransactionRequest completes a signature request by writing the signature bytes to zenrockd.
// The sender must be authorized to submit transactions for the keyring corresponding to the requestID.
// The transaction will be rejected if the TreasuryTxClient does not have the correct identity address.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - requestID: The unique identifier of the ICA transaction request
//   - signedData: The signed transaction bytes
//   - partySignature: The signature of the authorized party
//
// Returns:
//   - string: The transaction hash if successful
//   - error: An error if the transaction fails or authorization is invalid
//
// Example:
//
//	hash, err := client.FulfilICATransactionRequest(
//	    context.Background(),
//	    123,
//	    signedBytes,
//	    partySignature,
//	)
//	if err != nil {
//	    // Handle error
//	}
func (c *TreasuryTxClient) FulfilICATransactionRequest(ctx context.Context, requestID uint64, signedData []byte, partySignature []byte) (string, error) {
	status := types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED

	msg := types.NewMsgFulfilICATransactionRequest(
		c.c.Identity.Address.String(),
		requestID,
		status,
		partySignature,
		signedData,
		"",
	)

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

// RejectICATransactionRequest notifies zenrockd that a signature request has been rejected.
// The sender must be authorized to submit transactions for the keyring corresponding to the requestID.
// The transaction will be rejected if the TreasuryTxClient does not have the correct identity address.
//
// Parameters:
//   - ctx: Context for the request, can be used for timeouts and cancellation
//   - requestID: The unique identifier of the ICA transaction request to reject
//   - reason: A string explaining why the request was rejected
//   - partySignature: The signature of the authorized party
//
// Returns:
//   - string: The transaction hash if successful
//   - error: An error if the transaction fails or authorization is invalid
//
// Example:
//
//	hash, err := client.RejectICATransactionRequest(
//	    context.Background(),
//	    123,
//	    "Invalid transaction parameters",
//	    partySignature,
//	)
//	if err != nil {
//	    // Handle error
//	}
func (c *TreasuryTxClient) RejectICATransactionRequest(ctx context.Context, requestID uint64, reason string, partySignature []byte) (string, error) {
	status := types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED

	msg := types.NewMsgFulfilICATransactionRequest(
		c.c.Identity.Address.String(),
		requestID,
		status,
		partySignature,
		nil,
		reason,
	)

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
