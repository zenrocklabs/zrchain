package client

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	cosmos_types "github.com/cosmos/cosmos-sdk/codec/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// TreasuryTxClient contains a raw tx client.
type TreasuryTxClient struct {
	c *RawTxClient
}

// NewTreasuryTxClient returns a TreasuryTxClient.
func NewTreasuryTxClient(c *RawTxClient) *TreasuryTxClient {
	return &TreasuryTxClient{c: c}
}

// NewKeyRequest builds a key request, signs the transaction, sends to the chain, and waits for inclusion in a block.
// CLI equivalent example:
//
// zenrockd tx treasury new-key-request workspace14a2hpadpsy9h4auve2z8lw keyring1pfnq7r04rept47gaf5cpdew2 bitcoin  --from alice --chain-id zenrock -y
func (c *TreasuryTxClient) NewKeyRequest(ctx context.Context, workspace string, keyring string, keyType string) (string, error) {
	msg := types.NewMsgNewKeyRequest(c.c.Identity.Address.String(), workspace, keyring, keyType, 0, 0)
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

func (c *TreasuryTxClient) NewZenBTCKeyRequest(ctx context.Context, workspace string, keyring string, keyType string,
	recipient_addr string, chain_type types.WalletType, chain_id uint64, return_address string) (string, error) {

	metadata := &types.ZenBTCMetadata{
		RecipientAddr: recipient_addr,
		ChainType:     chain_type,
		ChainId:       chain_id,
		ReturnAddress: return_address,
	}

	msg := &types.MsgNewKeyRequest{
		Creator:        c.c.Identity.Address.String(),
		WorkspaceAddr:  workspace,
		KeyringAddr:    keyring,
		KeyType:        keyType,
		Btl:            0,
		SignPolicyId:   0,
		ZenbtcMetadata: metadata,
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

func (c *TreasuryTxClient) NewZrSignKeyRequest(ctx context.Context, creator, address string, walletType uint64) (string, error) {
	/// zenrockd tx treasury new-key-request "" ""  "" --ext-requester eth_address --ext-key-type 60 --from ZrSignConnector --chain-id zenrock --fees 20urock
	msg := &types.MsgNewKeyRequest{
		Creator:      creator,
		ExtRequester: address,
		ExtKeyType:   walletType,
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

func (c *TreasuryTxClient) NewSignTransactionRequest(ctx context.Context, keyID uint64, unsignedTransaction []byte, walletType types.WalletType, metadata *cosmos_types.Any) (string, error) {
	msg := &types.MsgNewSignTransactionRequest{
		Creator:             c.c.Identity.Address.String(),
		KeyId:               keyID,
		WalletType:          walletType,
		UnsignedTransaction: unsignedTransaction,
		Metadata:            metadata,
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

// NewSignatureRequest creates a new signature request for treasury transactions.
// It converts the data for signing into a comma-separated list, creates a new signature request message,
// builds and signs the transaction, sends it to the blockchain & waits for inclusion in a block.
// CLI Equivalent example:
//
// zenrockd tx treasury new-signature-request 1 50081cf6e000400018985834e7ead66fc0a0ce7fbdb220ad88b5f9052bf6814f --yes --from alice --chain-id zenrock
func (c *TreasuryTxClient) NewSignatureRequest(ctx context.Context, keyID uint64, dataForSigning [][]byte, cacheID []byte, unsignedPlusTX []byte) (string, error) {
	// convert data for signing into a comma separated list

	dataForSigningCSV := ""
	separator := ""
	for i, data := range dataForSigning {
		if i > 0 {
			separator = ","
		}
		dataForSigningCSV += separator + hex.EncodeToString(data)
	}

	msg := &types.MsgNewSignatureRequest{
		Creator:        c.c.Identity.Address.String(),
		KeyId:          keyID,
		DataForSigning: dataForSigningCSV,
		Btl:            0,
		CacheId:        cacheID,
		// VerifySigningData:        unsignedPlusTX,
		// VerifySigningDataVersion: types.VerificationVersion_BITCOIN_PLUS,
	}

	txBytes, err := c.c.BuildAndSignTx(ctx, DefaultGasLimit, DefaultFees, msg)
	if err != nil {
		return "", err
	}
	tx, err := c.c.SendWaitTx(ctx, txBytes)
	if err != nil {
		return "", err
	}
	return tx, nil
}

func (c *TreasuryTxClient) NewZrSignSignatureRequest(ctx context.Context, address string, keyType, walletIndex, walletType uint64, data string, cacheID []byte, metadata *cosmos_types.Any, tx, bCast bool, btl uint64) (string, error) {
	// zenrockd tx treasury new-zr-sign-signature-request eth_address 0 0 --data-for-signing 746573742064617461 --from alice --chain-id zenrock --fees 20urock
	msg := &types.MsgNewZrSignSignatureRequest{
		Creator:     c.c.Identity.Address.String(),
		Address:     address,
		KeyType:     keyType,
		WalletIndex: walletIndex,
		WalletType:  types.WalletType(walletType),
		Data:        data,
		Btl:         btl,
		Metadata:    metadata,
		Tx:          tx,
		CacheId:     cacheID,
		NoBroadcast: bCast,
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

// FulfilKeyRequest completes a key request writing the public key bytes to zenrockd. Note that the sender must be authorized to submit transactions
// for the keyring corresponding to the requestID. The transaction will be rejected if the TreasuryTxClient does not have the correct identity address.
func (c *TreasuryTxClient) FulfilKeyRequest(ctx context.Context, requestID uint64, publicKey []byte, partySignature []byte) (string, error) {
	status := types.KeyRequestStatus_KEY_REQUEST_STATUS_FULFILLED
	result := types.NewMsgFulfilKeyRequestKey(publicKey)

	msg := types.NewMsgFulfilKeyRequest(
		c.c.Identity.Address.String(),
		requestID,
		status,
		result,
		partySignature,
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

// FulfilSignatureRequest completes a signature request writing the signature bytes to zenrockd. The sender must be authorized to submit transactions
// for the keyring corresponding to the requestID. The transaction will be rejected if the TreasuryTxClient does not have the correct identity address.
func (c *TreasuryTxClient) FulfilSignatureRequest(ctx context.Context, requestID uint64, sig []byte, partySignature []byte) (string, error) {
	status := types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED

	msg := types.NewMsgFulfilSignatureRequest(
		c.c.Identity.Address.String(),
		requestID,
		status,
		partySignature,
		sig,
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

	// Get the TX after broadcast
	tx, err := c.c.client.GetTx(ctx, &txtypes.GetTxRequest{Hash: hash})
	if err != nil {
		return "", err
	}

	// Errors coming from the Zenrock chain
	if tx.TxResponse.Code != 0 {
		return "", fmt.Errorf("raw_log: %s", tx.TxResponse.RawLog)
	}

	return hash, nil
}

// RejectSignatureRequest notifies zenrockd that a signature request has been rejected. The sender must be authorized to submit transactions
// for the keyring corresponding to the requestID. The transaction will be rejected if the TreasuryTxClient does not have the correct identity address.
func (c *TreasuryTxClient) RejectSignatureRequest(ctx context.Context, requestID uint64, reason string, partySignature []byte) (string, error) {
	status := types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED

	msg := types.NewMsgFulfilSignatureRequest(
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

func (c *TreasuryTxClient) PartySignature(requestID uint64) ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, requestID)
	sig, err := c.c.Identity.PrivKey.Sign(buf)
	if err != nil {
		return nil, err
	}
	return sig, nil
}
