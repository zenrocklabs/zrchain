package integration

import (
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"

	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/proto"
)

func Test_Integration_MultiGrant_Add(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	countBefore, err := c.GetGrantsCount(ctx, c.IdentitySigner.Address.String())
	require.Nil(t, err)

	grantee := AccAddress()
	// make sure the account is known on chain by sending it 1 urock
	_, err = c.BankSend(ctx, c.IdentityAlice.Address.String(), grantee, 1)
	require.Nil(t, err)

	_, err = c.AddMultiGrant(ctx, grantee, []string{
		"/zrchain.policy.MsgApproveAction",
		"/zrchain.policy.MsgRevokeAction",
	})
	require.Nil(t, err)

	countAfter, err := c.GetGrantsCount(ctx, c.IdentitySigner.Address.String())
	require.Nil(t, err)

	assert.Greater(t, countAfter, countBefore)

	grants, err := c.GetGranteeGrants(ctx, grantee)
	require.Nil(t, err)

	require.Equal(t, 2, len(grants))

	assert.Equal(t, c.IdentityAlice.Address.String(), grants[0].Granter)
	assert.Equal(t, c.IdentityAlice.Address.String(), grants[1].Granter)

	assert.Equal(t, grantee, grants[0].Grantee)
	assert.Equal(t, grantee, grants[1].Grantee)

	g1, err := c.DecodeGenericGrant(*grants[0].Authorization)
	require.Nil(t, err)
	g2, err := c.DecodeGenericGrant(*grants[1].Authorization)
	require.Nil(t, err)

	assert.Equal(t, "/zrchain.policy.MsgApproveAction", g1.Msg)
	assert.Equal(t, "/zrchain.policy.MsgRevokeAction", g2.Msg)
}

func Test_Integration_MultiGrant_Remove(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	grantee := AccAddress()
	// make sure the account is known on chain by sending it 1 urock
	_, err := c.BankSend(ctx, c.IdentityAlice.Address.String(), grantee, 1)
	require.Nil(t, err)

	_, err = c.AddMultiGrant(ctx, grantee, []string{
		"/zrchain.policy.MsgApproveAction",
		"/zrchain.policy.MsgRevokeAction",
	})
	require.Nil(t, err)

	countBefore, err := c.GetGrantsCount(ctx, c.IdentitySigner.Address.String())
	require.Nil(t, err)
	assert.Greater(t, countBefore, uint64(0))

	_, err = c.RemoveMultiGrant(ctx, grantee, []string{
		"/zrchain.policy.MsgRevokeAction",
		"/zrchain.policy.MsgApproveAction",
	})
	require.Nil(t, err)

	countAfter, err := c.GetGrantsCount(ctx, c.IdentitySigner.Address.String())
	require.Nil(t, err)

	assert.Less(t, countAfter, countBefore)

	grants, err := c.GetGranteeGrants(ctx, grantee)
	require.Nil(t, err)

	require.Equal(t, 0, len(grants))
}

func Test_Integration_MultiGrant_Send(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	grantee := c.IdentityBob.Address.String()

	_, err := c.AddMultiGrant(ctx, grantee, []string{
		"/cosmos.bank.v1beta1.MsgSend",
	})
	require.Nil(t, err)

	receiver := AccAddress()

	msg := &banktypes.MsgSend{
		FromAddress: c.IdentityAlice.Address.String(),
		ToAddress:   receiver,
		Amount:      types.NewCoins(types.NewCoin("urock", math.NewInt(1))),
	}

	res, evts, err := c.ExecuteAs(ctx, c.IdentityBob, msg)
	require.Nil(t, err)
	require.NotNil(t, evts)
	require.NotNil(t, res)
}

func Test_Integration_MultiGrant_ApproveAction(t *testing.T) {
	t.SkipNow()
	ctx, c := testclient.GetTestClient()

	grantee := c.IdentityCharlie.Address.String()
	// make sure the account is known on chain and has enough urock to pay the gas
	_, err := c.BankSend(ctx, c.IdentityAlice.Address.String(), grantee, 10000000)
	require.Nil(t, err)

	_, _, err = c.AddMultiGrantWithIdentity(ctx, c.IdentityBob, grantee, []string{
		"/zrchain.policy.MsgApproveAction",
	})
	require.Nil(t, err)

	pol, err := c.CreateBoolParsePolicy(ctx, "test-approve-action", []string{
		c.IdentitySigner.Address.String(),
		c.IdentityBob.Address.String(),
	}, 2, 1000)
	require.Nil(t, err, "CreateBoolParsePolicy failed")

	ws, err := c.CreateWorkspace(ctx, pol.Id, pol.Id, []string{
		c.IdentityBob.Address.String(),
	})
	require.Nil(t, err, "CreateWorkspace failed")

	_, evts, err := c.CreateKey(ctx, ws.Addr, c.Keyring, "ecdsa", 0)
	require.Nil(t, err, "CreateKey failed")

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

	msg := &policytypes.MsgApproveAction{
		Creator:    c.IdentityBob.Address.String(),
		ActionType: action.Action.Msg.TypeUrl,
		ActionId:   action.Id,
	}

	res, evts, err := c.ExecuteAs(ctx, c.IdentityCharlie, msg)
	require.Nil(t, err)
	require.NotNil(t, evts)
	require.NotNil(t, res)

	actres := &policytypes.MsgApproveActionResponse{}
	err = proto.Unmarshal(res.Results[0], actres)
	require.Nil(t, err)
	assert.Equal(t, "ACTION_STATUS_COMPLETED", actres.Status)
}
