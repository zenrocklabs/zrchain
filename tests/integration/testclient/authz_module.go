package testclient

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/gogoproto/proto"
)

func (c *TestClient) GetGrantsCount(ctx context.Context, granter string) (uint64, error) {
	res, err := c.aqc.GranterGrants(ctx, &authztypes.QueryGranterGrantsRequest{
		Granter: granter,
	})
	if err != nil {
		return 0, err
	}

	return res.Pagination.Total, nil
}

func (c *TestClient) ExecuteAs(ctx context.Context, identity client.Identity, msgs ...cosmostypes.Msg) (*authztypes.MsgExecResponse, []abcitypes.Event, error) {
	anymsgs := []*codectypes.Any{}
	for _, msg := range msgs {
		anyMsg, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, nil, err
		}
		anymsgs = append(anymsgs, anyMsg)
	}
	msg := &authztypes.MsgExec{
		Grantee: identity.Address.String(),
		Msgs:    anymsgs,
	}

	txres, evts, err := c.executeTxWithIdentity(ctx, identity, msg)
	if err != nil {
		return nil, nil, err
	}

	res := &authztypes.MsgExecResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, evts, err
}

func (c *TestClient) GetGranteeGrants(ctx context.Context, grantee string) ([]*authztypes.GrantAuthorization, error) {
	res, err := c.aqc.GranteeGrants(ctx, &authztypes.QueryGranteeGrantsRequest{
		Grantee: grantee,
	})
	if err != nil {
		return nil, err
	}
	return res.Grants, nil
}

func (c *TestClient) DecodeGenericGrant(grant codectypes.Any) (*authztypes.GenericAuthorization, error) {
	g := &authztypes.GenericAuthorization{}
	err := proto.Unmarshal(grant.Value, g)
	return g, err
}
