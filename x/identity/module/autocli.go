package identity

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/Zenrock-Foundation/zrchain/v5/api/zrchain/identity"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "Workspaces",
					Use:       "workspaces --owner [owner]",
					Short:     "Query workspaces, optionally filtering by owner",
				},
				{
					RpcMethod: "WorkspaceByAddress",
					Use:       "workspace-by-address [workspace-addr]",
					Short:     "Query workspaceByAddress",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "workspace_addr"},
					},
				},
				{
					RpcMethod: "Keyrings",
					Use:       "keyrings",
					Short:     "Query keyrings",
				},
				{
					RpcMethod: "KeyringByAddress",
					Use:       "keyring-by-address [keyring-addr]",
					Short:     "Query keyring by address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring_addr"},
					},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "NewWorkspace",
					Use:       "new-workspace --admin-policy-id [admin-policy-id] --sign-policy-id [sign-policy-id] --additional-owners [additional-owners]",
					Short:     "Send a new-workspace tx",
				},
				{
					RpcMethod: "AddWorkspaceOwner",
					Use:       "add-workspace-owner [workspace-addr] [owner-address] --btl [btl]",
					Short:     "Send a add-workspace-owner tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "workspace_addr"},
						{ProtoField: "new_owner"},
					},
				},
				{
					RpcMethod: "AppendChildWorkspace",
					Use:       "append-child-workspace [parent-workspace-addr] [child-workspace-addr] --btl [btl]",
					Short:     "Send a append-child-workspace tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "parent_workspace_addr"},
						{ProtoField: "child_workspace_addr"},
					},
				},
				{
					RpcMethod: "NewChildWorkspace",
					Use:       "new-child-workspace [parent-workspace-addr] --btl [btl]",
					Short:     "Send a new-child-workspace tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "parent_workspace_addr"},
					},
				},
				{
					RpcMethod: "RemoveWorkspaceOwner",
					Use:       "remove-workspace-owner [workspace-address] [owner-address] --btl [btl]",
					Short:     "Send a remove-workspace-owner tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "workspace_addr"},
						{ProtoField: "owner"},
					},
				},
				{
					RpcMethod: "UpdateWorkspace",
					Use:       "update-workspace [workspace-address] [admin-policy-id] [sign-policy-id] --btl [btl]",
					Short:     "Send a update-workspace tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "workspace_addr"},
						{ProtoField: "admin_policy_id"},
						{ProtoField: "sign_policy_id"},
					},
				},
				{
					RpcMethod: "AddKeyringParty",
					Use:       "add-keyring-party [keyring-addr] [party] --increase-threshold [true]",
					Short:     "Send a add-keyring-party tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring_addr"},
						{ProtoField: "party"},
					},
				},
				{
					RpcMethod: "UpdateKeyring",
					Use:       "update-keyring [keyring-addr] [is-active:true|false] [party-threshold] [key-req-fee] [sig-req-fee] [description]",
					Short:     "Send a update-keyring tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring_addr"},
						{ProtoField: "is_active"},
						{ProtoField: "party_threshold"},
						{ProtoField: "key_req_fee"},
						{ProtoField: "sig_req_fee"},
						{ProtoField: "description"},
					},
				},
				{
					RpcMethod: "RemoveKeyringParty",
					Use:       "remove-keyring-party [keyring-addr] [party] --decrease-threshold [true]",
					Short:     "Send a remove-keyring-party tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring_addr"},
						{ProtoField: "party"},
					},
				},
				{
					RpcMethod: "NewKeyring",
					Use:       "new-keyring [description] [key-request-fee] [sign-request-fee]",
					Short:     "Send a new-keyring tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "description"},
						{ProtoField: "key_req_fee"},
						{ProtoField: "sig_req_fee"},
					},
				},
				{
					RpcMethod: "AddKeyringAdmin",
					Use:       "add-keyring-admin [keyring-addr] [admin]",
					Short:     "Send a add-keyring-admin tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring_addr"},
						{ProtoField: "admin"},
					},
				},
				{
					RpcMethod: "RemoveKeyringAdmin",
					Use:       "remove-keyring-admin [keyring-addr] [admin]",
					Short:     "Send a remove-keyring-admin tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring_addr"},
						{ProtoField: "admin"},
					},
				},
				{
					RpcMethod: "DeactivateKeyring",
					Use:       "deactivate-keyring [keyring-addr]",
					Short:     "Send a deactivate-keyring tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring_addr"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
