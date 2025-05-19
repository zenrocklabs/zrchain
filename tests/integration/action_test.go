package integration

import (
	"strconv"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Integration_Action_Approve(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	// create a 2 user policy
	pol, err := c.CreateBoolParsePolicy(ctx, "test-approve-action", []string{
		c.IdentitySigner.Address.String(),
		c.IdentityBob.Address.String(),
	}, 2, 1000)
	require.Nil(t, err, "CreateBoolParsePolicy failed")

	// create a workspace
	ws, err := c.CreateWorkspace(ctx, pol.Id, pol.Id, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err, "CreateWorkspace failed")

	// request a new key
	kr, evts, err := c.CreateKey(ctx, ws.Addr, c.Keyring, "ecdsa", 0)
	require.Nil(t, err, "CreateKey failed")
	require.NotNil(t, kr)
	require.NotNil(t, evts)

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

	// get action details
	action, err := c.GetActionDetails(ctx, actionId)
	require.Nil(t, err, "GetActionDetails failed")
	require.NotNil(t, action)

	// sign the action with 2nd user
	approveAction, evts, err := c.ApproveAction(ctx, c.IdentityBob, action.Action.Msg.TypeUrl, actionId)
	require.Nil(t, err, "ApproveAction failed")
	require.NotNil(t, approveAction)

	// verify action has been approved
	assert.Equal(t, "ACTION_STATUS_COMPLETED", approveAction.Status)

	// verify keyrequest has been created
	keyreqIdValue := ""
	for _, evt := range evts {
		if evt.Type == "new_key_request" {
		outer2:
			for _, attr := range evt.Attributes {
				switch attr.Key {
				case "request_id":
					keyreqIdValue = attr.Value
					break outer2
				}
			}
		}
	}

	require.NotEmpty(t, keyreqIdValue, "KeyRequestId not found")

	keyreqId, err := strconv.ParseUint(keyreqIdValue, 10, 64)
	require.Nil(t, err)

	keyreq, err := c.GetKeyRequest(ctx, keyreqId)
	require.Nil(t, err, "GetKeyRequest failed")
	require.NotNil(t, keyreq)
}
