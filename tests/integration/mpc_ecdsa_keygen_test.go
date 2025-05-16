//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Integration_Request_ECDSA_Key(t *testing.T) {
	t.Log("Running integration test: Request ECDSA Key")
	ctx, c := testclient.GetTestClient()

	ws, err := c.CreateWorkspace(ctx, 0, 0, []string{})
	require.Nil(t, err, "CreateWorkspace failed")

	krr, _, err := c.CreateKey(ctx, ws.Addr, "keyring1pfnq7r04rept47gaf5cpdew2", "ecdsa", 0)
	require.Nil(t, err, "CreateKey failed")
	require.NotNil(t, krr)
	assert.Greater(t, krr.KeyReqId, uint64(0))

	kr, err := c.GetKeyRequest(ctx, krr.KeyReqId)
	require.Nil(t, err)

	assert.Equal(t, keyring.Addr, kr.KeyringAddr)
	assert.Equal(t, "KEY_TYPE_ECDSA_SECP256K1", kr.KeyType)
	assert.Equal(t, ws.Addr, kr.WorkspaceAddr)
	assert.Equal(t, "KEY_REQUEST_STATUS_PENDING", kr.Status)
}
