package integration

import (
	"context"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/tests/integration/testclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func Test_Integration_Policy_Create(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	polCountBefore, err := c.GetPolicyCount(ctx)
	require.Nil(t, err)

	polRes, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
	}, 1, 123)
	require.Nil(t, err)

	polCountAfter, err := c.GetPolicyCount(ctx)
	require.Nil(t, err)
	assert.Greater(t, polCountAfter, polCountBefore)

	require.NotNil(t, polRes)

	pol, err := c.GetPolicy(ctx, polRes.Id)
	require.Nil(t, err)
	assert.Equal(t, uint64(123), pol.Btl)

	bpp, err := c.GetBoolParsePolicy(ctx, polRes.Id)
	require.Nil(t, err)

	require.Equal(t, 1, len(bpp.Participants))
	assert.Equal(t, c.IdentityAlice.Address.String(), bpp.Participants[0].Address)
	assert.Equal(t, "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 0", bpp.Definition)
}

func Test_Integration_Policy_CreateMultiUser(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	polCountBefore, err := c.GetPolicyCount(ctx)
	require.Nil(t, err)

	polRes, err := c.CreateBoolParsePolicy(ctx, "test-policy", []string{
		c.IdentityAlice.Address.String(),
		c.IdentityBob.Address.String(),
	}, 2, 0)
	require.Nil(t, err)

	polCountAfter, err := c.GetPolicyCount(ctx)
	require.Nil(t, err)
	assert.Greater(t, polCountAfter, polCountBefore)

	require.NotNil(t, polRes)

	bpp, err := c.GetBoolParsePolicy(ctx, polRes.Id)
	require.Nil(t, err)

	require.Equal(t, 2, len(bpp.Participants))
	assert.Equal(t, c.IdentityAlice.Address.String(), bpp.Participants[0].Address)
	assert.Equal(t, c.IdentityBob.Address.String(), bpp.Participants[1].Address)
	assert.Equal(t, "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty+zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq > 1", bpp.Definition)
}

func validateBoolParsePolicy(t *testing.T, ctx context.Context, c *testclient.TestClient, boolPolicy *policytypes.BoolparserPolicy) {
	polCountBefore, err := c.GetPolicyCount(ctx)
	require.Nil(t, err)

	anyPol, err := codectypes.NewAnyWithValue(boolPolicy)
	require.Nil(t, err)

	_, err = c.CreatePolicyRaw(ctx, "test-policy", anyPol, 123)
	require.NotNil(t, err)

	polCountAfter, err := c.GetPolicyCount(ctx)
	require.Nil(t, err)
	assert.Equal(t, polCountAfter, polCountBefore)
}

func Test_Integration_Policy_Create_InvalidAddress(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	boolPolicy := &policytypes.BoolparserPolicy{
		Participants: []*policytypes.PolicyParticipant{
			{
				Address: "invalid_address",
			},
		},
		Definition: "invalid_address > 0",
	}

	validateBoolParsePolicy(t, ctx, c, boolPolicy)
}

func Test_Integration_Policy_Create_MissingPartitipant(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	boolPolicy := &policytypes.BoolparserPolicy{
		Participants: []*policytypes.PolicyParticipant{
			{
				Address: c.IdentityAlice.Address.String(),
			},
		},
		Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + a2 > 0",
	}

	validateBoolParsePolicy(t, ctx, c, boolPolicy)
}

func Test_Integration_Policy_Create_MissingAddress(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	boolPolicy := &policytypes.BoolparserPolicy{
		Participants: []*policytypes.PolicyParticipant{
			{
				Address: "",
			},
		},
		Definition: " > 0",
	}

	validateBoolParsePolicy(t, ctx, c, boolPolicy)
}

func Test_Integration_Policy_Create_DuplicateAddress(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	boolPolicy := &policytypes.BoolparserPolicy{
		Participants: []*policytypes.PolicyParticipant{
			{
				Address: c.IdentityAlice.Address.String(),
			},
			{
				Address: c.IdentityAlice.Address.String(),
			},
		},
		Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty + zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 0",
	}

	validateBoolParsePolicy(t, ctx, c, boolPolicy)
}

func Test_Integration_Policy_Create_UnusedParticipant(t *testing.T) {
	ctx, c := testclient.GetTestClient()

	boolPolicy := &policytypes.BoolparserPolicy{
		Participants: []*policytypes.PolicyParticipant{
			{
				Address: c.IdentityAlice.Address.String(),
			},
			{
				Address: c.IdentityBob.Address.String(),
			},
		},
		Definition: "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty > 0",
	}

	validateBoolParsePolicy(t, ctx, c, boolPolicy)
}
