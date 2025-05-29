package testclient

import (
	"context"

	"github.com/cosmos/gogoproto/proto"

	identitytypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
)

func (c *TestClient) CreateWorkspace(ctx context.Context, adminPol, signPol uint64, additionalOwners []string) (*identitytypes.MsgNewWorkspaceResponse, error) {
	msg := &identitytypes.MsgNewWorkspace{
		Creator:          c.txc.Identity.Address.String(),
		AdminPolicyId:    adminPol,
		SignPolicyId:     signPol,
		AdditionalOwners: additionalOwners,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	wsres := &identitytypes.MsgNewWorkspaceResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, wsres)

	return wsres, err
}

func (c *TestClient) AddWorkspaceOwner(ctx context.Context, workspaceAddress string, ownerAddress string) (*identitytypes.MsgAddWorkspaceOwnerResponse, error) {
	msg := &identitytypes.MsgAddWorkspaceOwner{
		Creator:       c.txc.Identity.Address.String(),
		WorkspaceAddr: workspaceAddress,
		NewOwner:      ownerAddress,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgAddWorkspaceOwnerResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) RemoveWorkspaceOwner(ctx context.Context, workspaceAddress string, ownerAddress string) (*identitytypes.MsgRemoveWorkspaceOwnerResponse, error) {
	msg := &identitytypes.MsgRemoveWorkspaceOwner{
		Creator:       c.txc.Identity.Address.String(),
		WorkspaceAddr: workspaceAddress,
		Owner:         ownerAddress,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgRemoveWorkspaceOwnerResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) AppendChildWorkspace(ctx context.Context, parent string, child string) (*identitytypes.MsgAppendChildWorkspaceResponse, error) {
	msg := &identitytypes.MsgAppendChildWorkspace{
		Creator:             c.txc.Identity.Address.String(),
		ParentWorkspaceAddr: parent,
		ChildWorkspaceAddr:  child,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgAppendChildWorkspaceResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) NewChildWorkspace(ctx context.Context, parent string) (*identitytypes.MsgNewChildWorkspaceResponse, error) {
	msg := &identitytypes.MsgNewChildWorkspace{
		Creator:             c.txc.Identity.Address.String(),
		ParentWorkspaceAddr: parent,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgNewChildWorkspaceResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) UpdateWorkspace(ctx context.Context, workspace string, adminPol, signPol uint64) (*identitytypes.MsgUpdateWorkspaceResponse, error) {
	msg := &identitytypes.MsgUpdateWorkspace{
		Creator:       c.txc.Identity.Address.String(),
		WorkspaceAddr: workspace,
		AdminPolicyId: adminPol,
		SignPolicyId:  signPol,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgUpdateWorkspaceResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) GetWorkspace(ctx context.Context, workspaceAddress string) (*identitytypes.Workspace, error) {
	res, err := c.iqc.WorkspaceByAddress(ctx, &identitytypes.QueryWorkspaceByAddressRequest{
		WorkspaceAddr: workspaceAddress,
	})
	if err != nil {
		return nil, err
	}

	return res.Workspace, nil
}

func (c *TestClient) GetWorkspaceCount(ctx context.Context) (uint64, error) {
	res, err := c.iqc.Workspaces(ctx, &identitytypes.QueryWorkspacesRequest{})
	if err != nil {
		return 0, err
	}

	return res.Pagination.Total, nil
}

func (c *TestClient) CreateKeyring(ctx context.Context, description string, partyThreshold uint32, keyReqFee, sigReqFee uint64) (*identitytypes.MsgNewKeyringResponse, error) {
	msg := &identitytypes.MsgNewKeyring{
		Creator:        c.txc.Identity.Address.String(),
		Description:    description,
		PartyThreshold: partyThreshold,
		KeyReqFee:      keyReqFee,
		SigReqFee:      sigReqFee,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgNewKeyringResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) UpdateKeyring(ctx context.Context, keyring, description string, partyThreshold uint32, keyReqFee, sigReqFee uint64, active bool) (*identitytypes.MsgUpdateKeyringResponse, error) {
	msg := &identitytypes.MsgUpdateKeyring{
		Creator:        c.txc.Identity.Address.String(),
		KeyringAddr:    keyring,
		Description:    description,
		PartyThreshold: partyThreshold,
		KeyReqFee:      keyReqFee,
		SigReqFee:      sigReqFee,
		IsActive:       active,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgUpdateKeyringResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) DeactivateKeyring(ctx context.Context, keyring string) (*identitytypes.MsgDeactivateKeyringResponse, error) {
	msg := &identitytypes.MsgDeactivateKeyring{
		Creator:     c.txc.Identity.Address.String(),
		KeyringAddr: keyring,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgDeactivateKeyringResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) AddKeyringAdmin(ctx context.Context, keyring, admin string) (*identitytypes.MsgAddKeyringAdminResponse, error) {
	msg := &identitytypes.MsgAddKeyringAdmin{
		Creator:     c.txc.Identity.Address.String(),
		KeyringAddr: keyring,
		Admin:       admin,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgAddKeyringAdminResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) AddKeyringParty(ctx context.Context, keyring, party string, incThreshold bool) (*identitytypes.MsgAddKeyringPartyResponse, error) {
	msg := &identitytypes.MsgAddKeyringParty{
		Creator:           c.txc.Identity.Address.String(),
		KeyringAddr:       keyring,
		Party:             party,
		IncreaseThreshold: incThreshold,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgAddKeyringPartyResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) RemoveKeyringAdmin(ctx context.Context, keyring, admin string) (*identitytypes.MsgRemoveKeyringAdminResponse, error) {
	msg := &identitytypes.MsgRemoveKeyringAdmin{
		Creator:     c.txc.Identity.Address.String(),
		KeyringAddr: keyring,
		Admin:       admin,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgRemoveKeyringAdminResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) RemoveKeyringParty(ctx context.Context, keyring, party string, decThreshold bool) (*identitytypes.MsgRemoveKeyringPartyResponse, error) {
	msg := &identitytypes.MsgRemoveKeyringParty{
		Creator:           c.txc.Identity.Address.String(),
		KeyringAddr:       keyring,
		Party:             party,
		DecreaseThreshold: decThreshold,
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &identitytypes.MsgRemoveKeyringPartyResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) GetKeyringCount(ctx context.Context) (uint64, error) {
	res, err := c.iqc.Keyrings(ctx, &identitytypes.QueryKeyringsRequest{})
	if err != nil {
		return 0, err
	}

	return res.Pagination.Total, nil
}

func (c *TestClient) GetKeyring(ctx context.Context, keyring string) (*identitytypes.Keyring, error) {
	res, err := c.iqc.KeyringByAddress(ctx, &identitytypes.QueryKeyringByAddressRequest{
		KeyringAddr: keyring,
	})
	if err != nil {
		return nil, err
	}

	return res.Keyring, nil
}
