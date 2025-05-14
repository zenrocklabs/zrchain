//go:build integration
// +build integration

package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"
	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Integration_KeyRequest_New(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	keyCount, err := c.GetKeyCount(ctx)
	require.Nil(t, err)

	ws, err := c.CreateWorkspace(ctx, 0, 0, []string{})
	require.Nil(t, err, "CreateWorkspace failed")

	if keyCount < 10 {
		err := c.IncrementKeyCount(ctx, keyCount, ws.Addr)
		require.Nil(t, err)
	}

	krr, _, err := c.CreateKey(ctx, ws.Addr, c.Keyring, "ecdsa", 0)
	require.Nil(t, err, "CreateKey failed")
	require.NotNil(t, krr)
	assert.Greater(t, krr.KeyReqId, uint64(0))

	time.Sleep(time.Second * 10)
	kr, err := c.GetKeyRequest(ctx, krr.KeyReqId)
	require.Nil(t, err)

	assert.Equal(t, c.Keyring, kr.KeyringAddr)
	assert.Equal(t, "KEY_TYPE_ECDSA_SECP256K1", kr.KeyType)
	assert.Equal(t, ws.Addr, kr.WorkspaceAddr)
	assert.Equal(t, "KEY_REQUEST_STATUS_FULFILLED", kr.Status)
}

func Test_Integration_KeyRequest_Policy(t *testing.T) {
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

	polRes, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
		c.IdentityBob.Address.String(),
	}, 2, 0)
	require.Nil(t, err)

	krr, _, err := c.CreateKey(ctx, ws.Addr, c.Keyring, "ecdsa", polRes.Id)
	require.Nil(t, err, "CreateKey failed")
	require.NotNil(t, krr)
	assert.Greater(t, krr.KeyReqId, uint64(0))

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

	_, evts, err := c.CreateSignatureRequest(ctx, keys[0].Key.Id, dataHex)
	require.Nil(t, err)

	actionIdValue := ""
	for _, evt := range evts {
		if evt.Type == "new_action" {
		outer:
			for _, attr := range evt.Attributes {
				switch attr.Key {
				case "action_id":
					actionIdValue = attr.Value
					break outer
				}
			}
		}
	}

	require.NotEmpty(t, actionIdValue, "ActionId not found")
	actionId, err := strconv.ParseUint(actionIdValue, 10, 64)
	require.Nil(t, err)

	action, err := c.GetActionDetails(ctx, actionId)
	require.Nil(t, err, "GetActionDetails failed")
	require.NotNil(t, action)

	assert.Equal(t, action.Policy.Id, polRes.Id)
}
