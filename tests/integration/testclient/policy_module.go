package testclient

import (
	"context"
	"fmt"
	"strings"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	client "github.com/Zenrock-Foundation/zrchain/v6/go-client"
	policytypes "github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
)

func (c *TestClient) CreateBoolParsePolicy(ctx context.Context, name string, participants []string, threshold uint, btl uint64) (*policytypes.MsgNewPolicyResponse, error) {
	boolPolicy := policytypes.BoolparserPolicy{
		Participants: []*policytypes.PolicyParticipant{},
	}

	participantAddresses := []string{}
	for _, p := range participants {
		boolPolicy.Participants = append(boolPolicy.Participants, &policytypes.PolicyParticipant{
			Address: p,
		})
		participantAddresses = append(participantAddresses, p)
	}

	parts := strings.Join(participantAddresses, "+")
	definition := fmt.Sprintf("%s > %d", parts, threshold-1)
	boolPolicy.Definition = definition

	anyPol, err := codectypes.NewAnyWithValue(&boolPolicy)
	if err != nil {
		return nil, err
	}

	msg := &policytypes.MsgNewPolicy{
		Creator: c.txc.Identity.Address.String(),
		Name:    name,
		Policy:  anyPol,
		Btl:     btl,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	pres := &policytypes.MsgNewPolicyResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, pres)

	return pres, err
}

func (c *TestClient) CreatePolicyRaw(ctx context.Context, name string, anyPol *codectypes.Any, btl uint64) (*policytypes.MsgNewPolicyResponse, error) {
	msg := &policytypes.MsgNewPolicy{
		Creator: c.txc.Identity.Address.String(),
		Name:    name,
		Policy:  anyPol,
		Btl:     btl,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	pres := &policytypes.MsgNewPolicyResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, pres)

	return pres, err
}

func (c *TestClient) GetPolicyCount(ctx context.Context) (uint64, error) {
	res, err := c.pqc.Policies(ctx, &policytypes.QueryPoliciesRequest{})
	if err != nil {
		return 0, err
	}

	return res.Pagination.Total, nil
}

func (c *TestClient) GetPolicy(ctx context.Context, policyId uint64) (*policytypes.Policy, error) {
	res, err := c.pqc.PolicyById(ctx, &policytypes.QueryPolicyByIdRequest{
		Id: policyId,
	})
	if err != nil {
		return nil, err
	}

	return res.Policy.Policy, err
}

func (c *TestClient) GetBoolParsePolicy(ctx context.Context, policyId uint64) (*policytypes.BoolparserPolicy, error) {
	res, err := c.pqc.PolicyById(ctx, &policytypes.QueryPolicyByIdRequest{
		Id: policyId,
	})
	if err != nil {
		return nil, err
	}

	bpp := &policytypes.BoolparserPolicy{}
	err = proto.Unmarshal(res.Policy.Policy.Policy.Value, bpp)
	return bpp, err
}

func (c *TestClient) AddMultiGrant(ctx context.Context, grantee string, msgs []string) (*policytypes.MsgAddMultiGrantResponse, error) {
	res, _, err := c.AddMultiGrantWithIdentity(ctx, c.IdentitySigner, grantee, msgs)
	return res, err
}

func (c *TestClient) AddMultiGrantWithIdentity(ctx context.Context, identity client.Identity, grantee string, msgs []string) (*policytypes.MsgAddMultiGrantResponse, []abcitypes.Event, error) {
	msg := &policytypes.MsgAddMultiGrant{
		Creator: identity.Address.String(),
		Grantee: grantee,
		Msgs:    msgs,
	}

	txres, evts, err := c.executeTxWithIdentity(ctx, identity, msg)
	if err != nil {
		return nil, nil, err
	}

	res := &policytypes.MsgAddMultiGrantResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, evts, err
}

func (c *TestClient) RemoveMultiGrant(ctx context.Context, grantee string, msgs []string) (*policytypes.MsgRemoveMultiGrantResponse, error) {
	msg := &policytypes.MsgRemoveMultiGrant{
		Creator: c.txc.Identity.Address.String(),
		Grantee: grantee,
		Msgs:    msgs,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &policytypes.MsgRemoveMultiGrantResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) AddSignMethod(ctx context.Context) (*policytypes.MsgAddSignMethodResponse, error) {
	signMethod := &policytypes.SignMethodPasskey{
		RawId:             nil,
		AttestationObject: nil,
		ClientDataJson:    nil,
		Active:            true,
	}
	sigConfigAny, err := codectypes.NewAnyWithValue(signMethod)
	if err != nil {
		return nil, err
	}

	msg := &policytypes.MsgAddSignMethod{
		Creator: c.txc.Identity.Address.String(),
		Config:  sigConfigAny,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &policytypes.MsgAddSignMethodResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) RemoveSignMethod(ctx context.Context) (*policytypes.MsgRemoveSignMethodResponse, error) {
	msg := &policytypes.MsgRemoveSignMethod{
		Creator: c.txc.Identity.Address.String(),
		Id:      "",
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &policytypes.MsgRemoveSignMethodResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) ApproveAction(ctx context.Context, identity client.Identity, actionType string, actionId uint64) (*policytypes.MsgApproveActionResponse, []abcitypes.Event, error) {
	msg := &policytypes.MsgApproveAction{
		Creator:              identity.Address.String(),
		ActionType:           actionType,
		ActionId:             actionId,
		AdditionalSignatures: nil,
	}

	txres, evts, err := c.executeTxWithIdentity(ctx, identity, msg)
	if err != nil {
		return nil, nil, err
	}

	res := &policytypes.MsgApproveActionResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, evts, err
}

func (c *TestClient) RevokeAction(ctx context.Context) (*policytypes.MsgRevokeActionResponse, error) {
	msg := &policytypes.MsgRevokeAction{
		Creator:  c.txc.Identity.Address.String(),
		ActionId: 0,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &policytypes.MsgRevokeActionResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) GetSignMethods(ctx context.Context, address string) ([]*policytypes.SignMethodPasskey, error) {
	res, err := c.pqc.SignMethodsByAddress(ctx, &policytypes.QuerySignMethodsByAddressRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	sigMethods := []*policytypes.SignMethodPasskey{}
	for _, s := range res.Config {
		if s.TypeUrl != "/zrchain.policy.SignMethodPasskey" {
			return nil, fmt.Errorf("unsupported signmethod found: %s", s.TypeUrl)
		}
		sm := &policytypes.SignMethodPasskey{}
		err := proto.Unmarshal(s.Value, sm)
		if err != nil {
			return nil, err
		}
		sigMethods = append(sigMethods, sm)
	}
	return sigMethods, nil
}

func (c *TestClient) GetSignMethodCount(ctx context.Context, address string) (uint64, error) {
	res, err := c.pqc.SignMethodsByAddress(ctx, &policytypes.QuerySignMethodsByAddressRequest{
		Address: address,
	})
	if err != nil {
		return 0, err
	}
	return res.Pagination.Total, nil
}

func (c *TestClient) GetActionDetails(ctx context.Context, actionId uint64) (*policytypes.QueryActionDetailsByIdResponse, error) {
	res, err := c.pqc.ActionDetailsById(ctx, &policytypes.QueryActionDetailsByIdRequest{
		Id: actionId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *TestClient) GetActions(ctx context.Context, address string, status policytypes.ActionStatus) ([]policytypes.ActionResponse, error) {
	res, err := c.pqc.Actions(ctx, &policytypes.QueryActionsRequest{
		Address: address,
		Status:  status,
	})
	if err != nil {
		return nil, err
	}

	return res.Actions, nil
}

func (c *TestClient) GetActionCount(ctx context.Context) (uint64, error) {
	res, err := c.pqc.Actions(ctx, &policytypes.QueryActionsRequest{})
	if err != nil {
		return 0, err
	}
	return res.Pagination.Total, nil
}

func (c *TestClient) GetPoliciesByCreator(ctx context.Context, creators []string) ([]*policytypes.Policy, error) {
	res, err := c.pqc.PoliciesByCreator(ctx, &policytypes.QueryPoliciesByCreatorRequest{
		Creators: creators,
	})
	if err != nil {
		return nil, err
	}
	return res.Policies, nil
}
