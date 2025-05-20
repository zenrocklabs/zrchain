package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func Test_Integration_SignatureRequest_New(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	ws, err := c.CreateWorkspace(ctx, 0, 0, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err, "CreateWorkspace failed")

	krr, _, err := c.CreateKey(ctx, ws.Addr, c.Keyring, "ecdsa", 0)
	require.Nil(t, err, "CreateKey failed")

	<-time.After(time.Second * 6)

	kr, err := c.GetKeyRequest(ctx, krr.KeyReqId)
	require.Nil(t, err, "GetKeyRequest failed")
	require.NotNil(t, kr)

	require.Equal(t, "KEY_REQUEST_STATUS_FULFILLED", kr.Status, "KeyRequest was not fulfilled, make sure kms is running")

	keys, err := c.GetKeys(ctx, ws.Addr, types.WalletType_WALLET_TYPE_EVM, []string{"zen"})
	require.Nil(t, err, "GetKeys failed")
	require.Equal(t, 1, len(keys))

	data := sha256.Sum256([]byte("some-data"))
	dataHex := hex.EncodeToString(data[:])

	countBefore, err := c.GetSignatureRequestCount(ctx)
	require.Nil(t, err)

	sigreq, _, err := c.CreateSignatureRequest(ctx, keys[0].Key.Id, dataHex)
	require.Nil(t, err)

	countAfter, err := c.GetSignatureRequestCount(ctx)
	require.Nil(t, err)
	require.Greater(t, countAfter, countBefore)

	<-time.After(time.Second * 6)

	signedSigReq, err := c.GetSignatureRequest(ctx, sigreq.SigReqId)
	require.Nil(t, err)
	require.Equal(t, "SIGN_REQUEST_STATUS_FULFILLED", signedSigReq.Status, "SignRequest was not fulfilled, make sure kms is running")

	assert.NotEmpty(t, signedSigReq.SignedData)
	assert.NotEmpty(t, signedSigReq.KeyringPartySigs)
}

func Test_Integration_SignTransactionRequest_Eth(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	keyCount, err := c.GetKeyCount(ctx)
	require.Nil(t, err)

	ws, err := c.CreateWorkspace(ctx, 0, 0, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err, "CreateWorkspace failed")

	if keyCount < 10 {
		err := c.IncrementKeyCount(ctx, keyCount, ws.Addr)
		require.Nil(t, err)
	}

	krr, _, err := c.CreateKey(ctx, ws.Addr, c.Keyring, "ecdsa", 0)
	require.Nil(t, err, "CreateKey failed")

	<-time.After(time.Second * 6)

	kr, err := c.GetKeyRequest(ctx, krr.KeyReqId)
	require.Nil(t, err, "GetKeyRequest failed")
	require.NotNil(t, kr)

	require.Equal(t, "KEY_REQUEST_STATUS_FULFILLED", kr.Status, "KeyRequest was not fulfilled, make sure kms is running")

	keys, err := c.GetKeys(ctx, ws.Addr, types.WalletType_WALLET_TYPE_EVM, []string{"zen"})
	require.Nil(t, err, "GetKeys failed")
	require.Equal(t, 1, len(keys))

	tx := ethtypes.NewTransaction(0, common.HexToAddress("0xCF4324b9FD2fBC83b98a1483E854D31D3E45944C"), big.NewInt(0), 21000, big.NewInt(75000000000), nil)
	txBz, err := tx.MarshalBinary()
	require.Nil(t, err)

	countBeforeTx, err := c.GetSignTransactionRequestCount(ctx)
	require.Nil(t, err)
	countBeforeSig, err := c.GetSignatureRequestCount(ctx)
	require.Nil(t, err)

	sigreq, _, err := c.CreateSignETHTransactionRequest(ctx, keys[0].Key.Id, txBz, 11155111)
	require.Nil(t, err)

	countAfterTx, err := c.GetSignTransactionRequestCount(ctx)
	require.Nil(t, err)
	countAfterSig, err := c.GetSignatureRequestCount(ctx)
	require.Nil(t, err)

	require.Greater(t, countAfterTx, countBeforeTx)
	require.Greater(t, countAfterSig, countBeforeSig)

	<-time.After(time.Second * 6)

	signedSigReq, err := c.GetSignatureRequest(ctx, sigreq.SignatureRequestId)
	require.Nil(t, err)
	require.Equal(t, "SIGN_REQUEST_STATUS_FULFILLED", signedSigReq.Status, "SignRequest was not fulfilled, make sure kms is running")

	assert.NotEmpty(t, signedSigReq.SignedData)
	assert.NotEmpty(t, signedSigReq.KeyringPartySigs)
}
