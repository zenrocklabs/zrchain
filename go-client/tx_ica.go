package client

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
)

// FulfilICATransactionRequest completes a signature request writing the signature bytes to zenrockd. The sender must be authorized to submit transactions
// for the keyring corresponding to the requestID. The transaction will be rejected if the TreasuryTxClient does not have the correct identity address.
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

// RejectICATransactionRequest notifies zenrockd that a signature request has been rejected. The sender must be authorized to submit transactions
// for the keyring corresponding to the requestID. The transaction will be rejected if the TreasuryTxClient does not have the correct identity address.
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
