package integration

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Integration_Keyring_Create(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	countBefore, err := c.GetKeyringCount(ctx)
	require.Nil(t, err)

	txres, err := c.CreateKeyring(ctx, "test-description", 1, 2, 3)
	require.Nil(t, err)

	countAfter, err := c.GetKeyringCount(ctx)
	require.Nil(t, err)

	keyring, err := c.GetKeyring(ctx, txres.Addr)
	require.Nil(t, err)

	assert.Greater(t, countAfter, countBefore)
	assert.Equal(t, "test-description", keyring.Description)
	assert.Equal(t, uint32(1), keyring.PartyThreshold)
	assert.Equal(t, uint64(2), keyring.KeyReqFee)
	assert.Equal(t, uint64(3), keyring.SigReqFee)
}

func Test_Integration_Keyring_Update(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	txres, err := c.CreateKeyring(ctx, "test-description", 1, 2, 3)
	require.Nil(t, err)

	_, err = c.AddKeyringParty(ctx, txres.Addr, c.IdentityAlice.Address.String(), false)
	require.Nil(t, err)
	_, err = c.AddKeyringParty(ctx, txres.Addr, c.IdentityBob.Address.String(), false)
	require.Nil(t, err)

	_, err = c.UpdateKeyring(ctx, txres.Addr, "another-test-description", 2, 5, 6, false)
	require.Nil(t, err)

	keyring, err := c.GetKeyring(ctx, txres.Addr)
	require.Nil(t, err)

	assert.Equal(t, "another-test-description", keyring.Description)
	assert.Equal(t, uint32(2), keyring.PartyThreshold)
	assert.Equal(t, uint64(5), keyring.KeyReqFee)
	assert.Equal(t, uint64(6), keyring.SigReqFee)
	assert.False(t, keyring.IsActive)

	require.Equal(t, 1, len(keyring.Admins))
	assert.Equal(t, c.IdentitySigner.Address.String(), keyring.Admins[0])
}

func Test_Integration_Keyring_Deactivate(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	txres, err := c.CreateKeyring(ctx, "test-description", 0, 0, 0)
	require.Nil(t, err)

	_, err = c.DeactivateKeyring(ctx, txres.Addr)
	require.Nil(t, err)

	keyring, err := c.GetKeyring(ctx, txres.Addr)
	require.Nil(t, err)

	assert.False(t, keyring.IsActive)
}

func Test_Integration_Keyring_AddAdmin(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	txres, err := c.CreateKeyring(ctx, "test-description", 0, 0, 0)
	require.Nil(t, err)

	_, err = c.AddKeyringAdmin(ctx, txres.Addr, c.IdentityBob.Address.String())
	require.Nil(t, err)

	keyring, err := c.GetKeyring(ctx, txres.Addr)
	require.Nil(t, err)
	require.NotNil(t, keyring)

	require.Equal(t, 2, len(keyring.Admins))
	assert.Equal(t, c.IdentityBob.Address.String(), keyring.Admins[1])
}

func Test_Integration_Keyring_AddParty(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	txres, err := c.CreateKeyring(ctx, "test-description", 0, 0, 0)
	require.Nil(t, err)

	_, err = c.AddKeyringParty(ctx, txres.Addr, c.IdentityBob.Address.String(), true)
	require.Nil(t, err)

	keyring, err := c.GetKeyring(ctx, txres.Addr)
	require.Nil(t, err)
	require.NotNil(t, keyring)

	require.Equal(t, 1, len(keyring.Parties))
	assert.Equal(t, c.IdentityBob.Address.String(), keyring.Parties[0])
	assert.Equal(t, uint32(1), keyring.PartyThreshold)
}

func Test_Integration_Keyring_RemoveAdmin(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	txres, err := c.CreateKeyring(ctx, "test-description", 0, 0, 0)
	require.Nil(t, err)

	_, err = c.AddKeyringAdmin(ctx, txres.Addr, c.IdentityBob.Address.String())
	require.Nil(t, err)

	_, err = c.RemoveKeyringAdmin(ctx, txres.Addr, c.IdentityAlice.Address.String())
	require.Nil(t, err)

	keyring, err := c.GetKeyring(ctx, txres.Addr)
	require.Nil(t, err)
	require.NotNil(t, keyring)

	require.Equal(t, 1, len(keyring.Admins))
	assert.Equal(t, c.IdentityBob.Address.String(), keyring.Admins[0])
}

func Test_Integration_Keyring_RemoveOnlyAdmin(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	txres, err := c.CreateKeyring(ctx, "test-description", 0, 0, 0)
	require.Nil(t, err)

	_, err = c.RemoveKeyringAdmin(ctx, txres.Addr, c.IdentityAlice.Address.String())
	require.NotNil(t, err)
}

func Test_Integration_Keyring_RemoveParty(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	txres, err := c.CreateKeyring(ctx, "test-description", 0, 0, 0)
	require.Nil(t, err)

	_, err = c.AddKeyringParty(ctx, txres.Addr, c.IdentityAlice.Address.String(), true)
	require.Nil(t, err)

	_, err = c.AddKeyringParty(ctx, txres.Addr, c.IdentityBob.Address.String(), true)
	require.Nil(t, err)

	_, err = c.RemoveKeyringParty(ctx, txres.Addr, c.IdentityAlice.Address.String(), true)
	require.Nil(t, err)

	keyring, err := c.GetKeyring(ctx, txres.Addr)
	require.Nil(t, err)
	require.NotNil(t, keyring)

	require.Equal(t, 1, len(keyring.Parties))
	assert.Equal(t, c.IdentityBob.Address.String(), keyring.Parties[0])
	assert.Equal(t, uint32(1), keyring.PartyThreshold)
}

func Test_Integration_Keyring_TransferFrom(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	wsres, err := c.CreateWorkspace(ctx, 0, 0, []string{})
	require.Nil(t, err)

	_, _, err = c.CreateKey(ctx, wsres.Addr, c.Keyring, "ecdsa", 0)
	require.Nil(t, err)

	_, err = c.TransferFromKeyring(ctx, c.IdentitySigner, c.Keyring, c.IdentityAlice.Address.String(), "urock", 1100000000000)
	require.NotNil(t, err) // should fail as balance is greater than available

	_, err = c.TransferFromKeyring(ctx, c.IdentitySigner, c.Keyring, c.IdentityAlice.Address.String(), "urock", 1)
	require.Nil(t, err)
}
