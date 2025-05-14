package integration

import (
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Integration_Workspace_Create(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	wsCountBefore, err := c.GetWorkspaceCount(ctx)
	require.Nil(t, err)

	wsres, err := c.CreateWorkspace(ctx, 0, 0, []string{})
	require.Nil(t, err)
	assert.NotEmpty(t, wsres.Addr)

	wsCountAfter, err := c.GetWorkspaceCount(ctx)
	require.Nil(t, err)

	assert.Greater(t, wsCountAfter, wsCountBefore)

	ws, err := c.GetWorkspace(ctx, wsres.Addr)
	require.Nil(t, err)
	require.NotNil(t, ws)
	assert.Equal(t, wsres.Addr, ws.Address)
}

func Test_Integration_Workspace_CreateWithAdditionalOwners(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	wsres, err := c.CreateWorkspace(ctx, 0, 0, []string{c.IdentityBob.Address.String()})

	require.Nil(t, err)
	assert.NotEmpty(t, wsres.Addr)

	ws, err := c.GetWorkspace(ctx, wsres.Addr)
	require.Nil(t, err)
	require.NotNil(t, ws)
	assert.Equal(t, 2, len(ws.Owners))
	assert.Contains(t, ws.Owners, c.IdentityBob.Address.String())
}

func Test_Integration_Workspace_CreateWithPolicy(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	polRes, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
	}, 1, 0)
	require.Nil(t, err)

	polRes2, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
	}, 1, 0)
	require.Nil(t, err)

	wsres, err := c.CreateWorkspace(ctx, polRes.Id, polRes2.Id, []string{})
	require.Nil(t, err)

	ws, err := c.GetWorkspace(ctx, wsres.Addr)
	require.Nil(t, err)

	assert.Equal(t, polRes.Id, ws.AdminPolicyId)
	assert.Equal(t, polRes2.Id, ws.SignPolicyId)
}

func Test_Integration_Workspace_CreateWithMissingOwnerInPolicy(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	polRes, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
		c.IdentityBob.Address.String(),
	}, 1, 0)
	require.Nil(t, err)

	_, err = c.CreateWorkspace(ctx, polRes.Id, polRes.Id, []string{})
	require.NotNil(t, err)
}

func Test_Integration_Workspace_CreateWithMultiUserPolicy(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	polRes, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
		c.IdentityBob.Address.String(),
	}, 1, 0)
	require.Nil(t, err)

	wsres, err := c.CreateWorkspace(ctx, polRes.Id, polRes.Id, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err)

	ws, err := c.GetWorkspace(ctx, wsres.Addr)
	require.Nil(t, err)
	assert.Equal(t, polRes.Id, ws.AdminPolicyId)
	assert.Equal(t, polRes.Id, ws.SignPolicyId)
}

func Test_Integration_Workspace_AddOwner(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	wsres, err := c.CreateWorkspace(ctx, 0, 0, []string{})
	require.Nil(t, err)

	_, err = c.AddWorkspaceOwner(ctx, wsres.Addr, c.IdentityBob.Address.String())
	require.Nil(t, err)

	ws, err := c.GetWorkspace(ctx, wsres.Addr)
	require.Nil(t, err)

	require.Equal(t, 2, len(ws.Owners))
	assert.Equal(t, c.IdentityAlice.Address.String(), ws.Owners[0])
	assert.Equal(t, c.IdentityBob.Address.String(), ws.Owners[1])
}

func Test_Integration_Workspace_RemoveOwner(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	wsres, err := c.CreateWorkspace(ctx, 0, 0, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err)

	_, err = c.RemoveWorkspaceOwner(ctx, wsres.Addr, c.IdentityAlice.Address.String())
	require.Nil(t, err)

	ws, err := c.GetWorkspace(ctx, wsres.Addr)
	require.Nil(t, err)

	require.Equal(t, 1, len(ws.Owners))
	assert.Equal(t, c.IdentityBob.Address.String(), ws.Owners[0])
}

func Test_Integration_Workspace_AppendChildWorkspace(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	wsresParent, err := c.CreateWorkspace(ctx, 0, 0, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err)

	wsresChild, err := c.CreateWorkspace(ctx, 0, 0, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err)

	_, err = c.AppendChildWorkspace(ctx, wsresParent.Addr, wsresChild.Addr)
	require.Nil(t, err)

	wsParent, err := c.GetWorkspace(ctx, wsresParent.Addr)
	require.Nil(t, err)

	require.Equal(t, 1, len(wsParent.ChildWorkspaces))
	assert.Equal(t, wsresChild.Addr, wsParent.ChildWorkspaces[0])
}

func Test_Integration_Workspace_NewChildWorkspace(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	wsresParent, err := c.CreateWorkspace(ctx, 0, 0, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err)

	wsresChild, err := c.NewChildWorkspace(ctx, wsresParent.Addr)
	require.Nil(t, err)

	wsParent, err := c.GetWorkspace(ctx, wsresParent.Addr)
	require.Nil(t, err)

	require.Equal(t, 1, len(wsParent.ChildWorkspaces))
	assert.Equal(t, wsresChild.Address, wsParent.ChildWorkspaces[0])
}

func Test_Integration_Workspace_Update(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	polRes1, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
	}, 1, 0)
	require.Nil(t, err)

	polRes2, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
	}, 1, 0)
	require.Nil(t, err)

	wsres, err := c.CreateWorkspace(ctx, 0, 0, []string{})
	require.Nil(t, err)

	_, err = c.UpdateWorkspace(ctx, wsres.Addr, polRes1.Id, polRes2.Id)
	require.Nil(t, err)

	ws, err := c.GetWorkspace(ctx, wsres.Addr)
	require.Nil(t, err)

	assert.Equal(t, polRes1.Id, ws.AdminPolicyId)
	assert.Equal(t, polRes2.Id, ws.SignPolicyId)
}

func Test_Integration_Workspace_UpdateWithMissingOwner(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	polRes, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
		c.IdentityBob.Address.String(),
	}, 1, 0)
	require.Nil(t, err)

	wsres, err := c.CreateWorkspace(ctx, 0, 0, []string{})
	require.Nil(t, err)

	_, err = c.UpdateWorkspace(ctx, wsres.Addr, polRes.Id, polRes.Id)
	require.NotNil(t, err)
}
