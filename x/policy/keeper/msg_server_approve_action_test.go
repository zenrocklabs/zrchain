package keeper_test

import (
	"encoding/base64"
	"testing"

	keepertest "github.com/Zenrock-Foundation/zrchain/v6/testutil/keeper"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	idtypes "github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/keeper"
	policy "github.com/Zenrock-Foundation/zrchain/v6/x/policy/module"
	"github.com/Zenrock-Foundation/zrchain/v6/x/policy/types"
)

type handlers struct {
	addWorkspaceOwnerCalled   bool
	addWorkspaceOwnerErr      error
	addWorkspaceOwnerAction   *types.Action
	addWorkspaceOwnerResponse *idtypes.MsgAddWorkspaceOwnerResponse
}

func (h *handlers) AddOwnerActionHandler(ctx sdk.Context, act *types.Action) (*idtypes.MsgAddWorkspaceOwnerResponse, error) {
	h.addWorkspaceOwnerCalled = true
	h.addWorkspaceOwnerAction = act
	return h.addWorkspaceOwnerResponse, h.addWorkspaceOwnerErr
}

func Test_msgServer_ApproveAction(t *testing.T) {

	h := &handlers{}
	addToWorkspaceMsg := idtypes.NewMsgAddWorkspaceOwner(
		"creator",
		"workspaceaddr",
		"newOwner",
		1000)
	addToWorkspaceMsgAny, _ := codectypes.NewAnyWithValue(addToWorkspaceMsg)

	policyData := types.BoolparserPolicy{
		Definition: "testApprover > 0",
		Participants: []*types.PolicyParticipant{
			{Address: "testApprover"},
		},
	}

	policyDataAny, err := codectypes.NewAnyWithValue(&policyData)
	require.NoError(t, err, "error decoding policyData")

	var defaultPolicy = types.Policy{
		Id:     1,
		Name:   "defaultPolicy",
		Policy: policyDataAny,
	}

	var defaultApproval = types.MsgApproveAction{
		Creator:    "testApprover",
		ActionType: addToWorkspaceMsgAny.TypeUrl,
		ActionId:   1,
	}

	var defaultAction = types.Action{
		Id:         1,
		Approvers:  []string{},
		Status:     types.ActionStatus_ACTION_STATUS_PENDING,
		PolicyId:   1,
		Msg:        addToWorkspaceMsgAny,
		Creator:    "testCreator",
		Btl:        1000,
		PolicyData: nil,
	}

	type args struct {
		action *types.Action
		policy *types.Policy
		msg    *types.MsgApproveAction
	}

	tests := []struct {
		name       string
		args       args
		want       *types.MsgApproveActionResponse
		wantAction *types.Action
		wantErr    bool
	}{
		{
			name: "PASS: approve action",
			args: args{
				action: &defaultAction,
				policy: &defaultPolicy,
				msg:    &defaultApproval,
			},
			want:    &types.MsgApproveActionResponse{Status: types.ActionStatus_ACTION_STATUS_COMPLETED.String()},
			wantErr: false,
		},
		{
			name: "FAIL: invalid action",
			args: args{
				action: &defaultAction,
				policy: &defaultPolicy,
				msg: &types.MsgApproveAction{
					Creator:    "testApprover",
					ActionType: addToWorkspaceMsgAny.TypeUrl,
					ActionId:   2,
				},
			},
			wantErr: true,
		},
		{
			name: "FAIL: invalid action status",
			args: args{
				action: &types.Action{
					Id:         1,
					Approvers:  []string{},
					Status:     types.ActionStatus_ACTION_STATUS_COMPLETED,
					PolicyId:   1,
					Msg:        addToWorkspaceMsgAny,
					Creator:    "testCreator",
					Btl:        1000,
					PolicyData: nil,
				},
				policy: &defaultPolicy,
				msg:    &defaultApproval,
			},
			wantErr: true,
		},
		{
			name: "FAIL: approver outside policy",
			args: args{
				action: &defaultAction,
				policy: &defaultPolicy,
				msg: &types.MsgApproveAction{
					Creator:    "invalidApprover",
					ActionType: addToWorkspaceMsgAny.TypeUrl,
					ActionId:   1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			pk := keepers.PolicyKeeper
			msgSer := keeper.NewMsgServerImpl(*pk)
			keeper.RegisterActionHandler(
				pk,
				"/zrchain.identity.MsgAddWorkspaceOwner",
				h.AddOwnerActionHandler,
			)
			polGenesis := types.GenesisState{
				Params:   types.Params{},
				PortId:   "42",
				Policies: []types.Policy{*tt.args.policy},
				Actions:  []types.Action{*tt.args.action},
			}
			policy.InitGenesis(keepers.Ctx, *pk, polGenesis)

			_, err := msgSer.ApproveAction(keepers.Ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, h.addWorkspaceOwnerCalled)
			}
		})
	}
}

func Test_msgServer_ApproveAction_Passkeys(t *testing.T) {

	h := &handlers{}
	addToWorkspaceMsg := idtypes.NewMsgAddWorkspaceOwner(
		"creator",
		"workspaceaddr",
		"newOwner",
		1000)
	addToWorkspaceMsgAny, _ := codectypes.NewAnyWithValue(addToWorkspaceMsg)

	policyData := types.BoolparserPolicy{
		Definition: "testApprover + passkey{c29tZS1wYXNza2V5LWlk} > 1",
		Participants: []*types.PolicyParticipant{
			{Address: "testApprover"},
			{Address: "passkey{c29tZS1wYXNza2V5LWlk}"}, // base64 urlencoded "some-passkey-id"
		},
	}

	policyDataAny, err := codectypes.NewAnyWithValue(&policyData)
	require.NoError(t, err, "error decoding policyData")

	var defaultPolicy = types.Policy{
		Id:     1,
		Name:   "defaultPolicy",
		Policy: policyDataAny,
	}

	valid_sig_challenge, _ := base64.StdEncoding.DecodeString("4SmZBUFhQRfG/jjWNxCbfdZyxstdOclxkGkT/3nmWOk=")
	valid_sig_clientdata, _ := base64.StdEncoding.DecodeString("eyJ0eXBlIjoid2ViYXV0aG4uZ2V0IiwiY2hhbGxlbmdlIjoiNFNtWkJVRmhRUmZHX2pqV054Q2JmZFp5eHN0ZE9jbHhrR2tUXzNubVdPayIsIm9yaWdpbiI6Imh0dHBzOi8vZGV2ZWxvcGVyLm1vemlsbGEub3JnIiwiY3Jvc3NPcmlnaW4iOmZhbHNlfQ==")
	valid_sig_authdata, _ := base64.StdEncoding.DecodeString("EX2C26HZXr74XcytYnlABgSk0ZjcNvHZSXM6Ejfm4y0dAAAAAA==")
	valid_sig_signature, _ := base64.StdEncoding.DecodeString("MEYCIQDoCuL2wydyiRjbc6KCwjvowIVPxjksaGlarIf+E4c8BgIhAJrECkq4tOJ6E4/QrJjfNIARzLMa6KUwnbPsd3Kta1YQ")

	additionalSig := &types.AdditionalSignaturePasskey{
		RawId:             []byte("some-passkey-id"),
		AuthenticatorData: valid_sig_authdata,
		ClientDataJson:    valid_sig_clientdata,
		Signature:         valid_sig_signature,
	}
	additionalSigAny, err := codectypes.NewAnyWithValue(additionalSig)
	require.NoError(t, err, "error encoding additionalSig")

	invalid_sig_signature, _ := base64.StdEncoding.DecodeString("MEUCIByGu1Fad/tI+5+AMNnb9mttdIaL1GyI4t1sufVPBtkSAiEAgfaeQ0RjRxABrG6qbADEaVmrgiQ/8UCMZtHlhqmNseM=")
	additionalSig_invalid_sig := &types.AdditionalSignaturePasskey{
		RawId:             []byte("some-passkey-id"),
		AuthenticatorData: valid_sig_authdata,
		ClientDataJson:    valid_sig_clientdata,
		Signature:         invalid_sig_signature,
	}
	additionalSigAny_invalid_sig, err := codectypes.NewAnyWithValue(additionalSig_invalid_sig)
	require.NoError(t, err, "error encoding additionalSigAny_invalid_sig")

	invalid_sig_clientdata, _ := base64.StdEncoding.DecodeString("eyJ0eXBlIjoid2ViYXV0aG4uZ2V0IiwiY2hhbGxlbmdlIjoiQ2lFRFIwTEhVVVJaWXpuSWgwbXZqM3F1dDVXUlRmczBhdEM5TThTY3dSVzZETFkiLCJvcmlnaW4iOiJodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZyIsImNyb3NzT3JpZ2luIjpmYWxzZX0=")
	additionalSig_invalid_challenge := &types.AdditionalSignaturePasskey{
		RawId:             []byte("some-passkey-id"),
		AuthenticatorData: valid_sig_authdata,
		ClientDataJson:    invalid_sig_clientdata,
		Signature:         valid_sig_signature,
	}
	additionalSigAny_invalid_challenge, err := codectypes.NewAnyWithValue(additionalSig_invalid_challenge)
	require.NoError(t, err, "error encoding additionalSig_invalid_challenge")

	valid_tx_bytes, _ := base64.StdEncoding.DecodeString("CoEECv4DCiAvemVucm9jay5wb2xpY3kuTXNnQWRkU2lnbk1ldGhvZBLZAwoqemVuMTN5M3RtNjhnbXU5a250Y3h3dm11ZTgycDZha2FjbnB0MnY3bnR5EqoDCiEvemVucm9jay5wb2xpY3kuU2lnbk1ldGhvZFBhc3NrZXkShAMKG1RiZlJGUDUzYTh5SlBrU0hJeXh3eWN2VFhYWRIUTbfRFP53a8yJPkSHIyxwycvTXXYatgGjY2ZtdGRub25lZ2F0dFN0bXSgaGF1dGhEYXRhWJgRfYLbodlevvhdzK1ieUAGBKTRmNw28dlJczoSN+bjLV0AAAAA+/wwBxVOTsyMC24CBVfXvQAUTbfRFP53a8yJPkSHIyxwycvTXXalAQIDJiABIVggf6JpOflG9El/S3+/YvEBT69317zGkfjG2XcvBluSMBYiWCDqkyO8H1ULfUL6lLUz7iQTC0Ilqu7Kwa8yS4nnQDsv5CKVAXsidHlwZSI6IndlYmF1dGhuLmNyZWF0ZSIsImNoYWxsZW5nZSI6IkNpRURSMExIVVVSWll6bkloMG12ajNxdXQ1V1JUZnMwYXRDOU04U2N3Ulc2RExVIiwib3JpZ2luIjoiaHR0cHM6Ly9kZXZlbG9wZXIubW96aWxsYS5vcmciLCJjcm9zc09yaWdpbiI6ZmFsc2V9ElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQNHQsdRRFljOciHSa+Peq63lZFN+zRq0L0zxJzBFboMtRIECgIIARgBEgQQwJoMGkD3sYFTIYM+izWXmSKQ+6ouuUwHfqPT0IWqG1mFnlmzhmiyk0MMGEhvlDDX2kMvCey/B7NS1+f5V1z9uQLN1S5R")
	valid_clientdata, _ := base64.StdEncoding.DecodeString("eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiQ2lFRFIwTEhVVVJaWXpuSWgwbXZqM3F1dDVXUlRmczBhdEM5TThTY3dSVzZETFUiLCJvcmlnaW4iOiJodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZyIsImNyb3NzT3JpZ2luIjpmYWxzZX0=")
	valid_attestation, _ := base64.StdEncoding.DecodeString("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYEX2C26HZXr74XcytYnlABgSk0ZjcNvHZSXM6Ejfm4y1dAAAAAPv8MAcVTk7MjAtuAgVX170AFIGqzDVypWqJMqUgI7gVXesVD/ffpQECAyYgASFYIPMPOLDaUaN8TL4B7BDt0dH1vnUFNk9P5uqej1QYOKiEIlggg27Llwf7cPImDTUUtz9zXUkrrPgOoAgDl5uQzYq2zkA=")
	valid_sign_method := &types.SignMethodPasskey{
		RawId:             []byte("some-passkey-id"),
		AttestationObject: valid_attestation,
		ClientDataJson:    valid_clientdata,
	}
	valid_config, err := codectypes.NewAnyWithValue(valid_sign_method)
	require.NoError(t, err, "error encoding valid config")

	var defaultApproval = types.MsgApproveAction{
		Creator:              "testApprover",
		ActionType:           addToWorkspaceMsgAny.TypeUrl,
		ActionId:             1,
		AdditionalSignatures: []*codectypes.Any{additionalSigAny},
	}
	var approval_invalid_sig = types.MsgApproveAction{
		Creator:              "testApprover",
		ActionType:           addToWorkspaceMsgAny.TypeUrl,
		ActionId:             1,
		AdditionalSignatures: []*codectypes.Any{additionalSigAny_invalid_sig},
	}
	var approvals_invalid_challenge = types.MsgApproveAction{
		Creator:              "testApprover",
		ActionType:           addToWorkspaceMsgAny.TypeUrl,
		ActionId:             1,
		AdditionalSignatures: []*codectypes.Any{additionalSigAny_invalid_challenge},
	}

	var defaultAction = types.Action{
		Id:        1,
		Approvers: []string{},
		Status:    types.ActionStatus_ACTION_STATUS_PENDING,
		PolicyId:  1,
		Msg:       addToWorkspaceMsgAny,
		Creator:   "testCreator",
		Btl:       1000,
		PolicyData: []*types.KeyValue{
			{
				Key:   "challenge-passkey{c29tZS1wYXNza2V5LWlk}",
				Value: valid_sig_challenge,
			},
		},
	}

	type args struct {
		action  *types.Action
		policy  *types.Policy
		msg     *types.MsgApproveAction
		creator string
		config  *codectypes.Any
	}

	tests := []struct {
		name       string
		args       args
		want       []string
		wantAction *types.Action
		wantErr    bool
	}{
		{
			name: "FAIL: config missing",
			args: args{
				action:  &defaultAction,
				policy:  &defaultPolicy,
				msg:     &defaultApproval,
				creator: "testApprover",
			},
			want:    []string{"testApprover"},
			wantErr: true,
		},
		{
			name: "FAIL: invalid signature",
			args: args{
				action:  &defaultAction,
				policy:  &defaultPolicy,
				msg:     &approval_invalid_sig,
				creator: "testApprover",
				config:  valid_config,
			},
			want:    []string{"testApprover"},
			wantErr: false,
		},
		{
			name: "FAIL: invalid challenge",
			args: args{
				action:  &defaultAction,
				policy:  &defaultPolicy,
				msg:     &approvals_invalid_challenge,
				creator: "testApprover",
				config:  valid_config,
			},
			want:    []string{"testApprover"},
			wantErr: false,
		},
		{
			name: "PASS: approve action",
			args: args{
				action:  &defaultAction,
				policy:  &defaultPolicy,
				msg:     &defaultApproval,
				creator: "testApprover",
				config:  valid_config,
			},
			want:    []string{"testApprover", "passkey{c29tZS1wYXNza2V5LWlk}"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keepers := keepertest.NewTest(t)
			pk := keepers.PolicyKeeper
			ctx := keepers.Ctx.WithTxBytes(valid_tx_bytes)
			msgSer := keeper.NewMsgServerImpl(*pk)
			keeper.RegisterActionHandler(
				pk,
				"/zrchain.identity.MsgAddWorkspaceOwner",
				h.AddOwnerActionHandler,
			)
			polGenesis := types.GenesisState{
				Params:   types.Params{},
				PortId:   "42",
				Policies: []types.Policy{*tt.args.policy},
				Actions:  []types.Action{*tt.args.action},
			}
			policy.InitGenesis(ctx, *pk, polGenesis)

			if tt.args.config != nil {
				_, err := msgSer.AddSignMethod(ctx, &types.MsgAddSignMethod{
					Creator: tt.args.creator,
					Config:  tt.args.config,
				})
				require.NoError(t, err)
			}

			_, err := msgSer.ApproveAction(ctx, tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if h.addWorkspaceOwnerAction != nil {
				require.Equal(t, tt.want, h.addWorkspaceOwnerAction.Approvers)
			}
		})
	}
}
