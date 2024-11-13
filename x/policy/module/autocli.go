package policy

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/Zenrock-Foundation/zrchain/v5/api/zrchain/policy"
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
					RpcMethod: "Actions",
					Use:       "actions --address [address] --status [status]",
					Short:     "Query actions, optionally filtering by address and status",
				},
				{
					RpcMethod: "Policies",
					Use:       "policies",
					Short:     "Query policies",
				},
				{
					RpcMethod: "PolicyById",
					Use:       "policy-by-id [id]",
					Short:     "Query policy-by-id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "SignMethodsByAddress",
					Use:       "sign-methods-by-address [address]",
					Short:     "Query signature methods assigned to a specified address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "PoliciesByCreator",
					Use:       "policies-by-creator [creators]",
					Short:     "Query policiesByCreator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "creators"},
					},
				},

				{
					RpcMethod:      "ActionDetailsById",
					Use:            "action-details-by-id [id]",
					Short:          "Query action_details_by_id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
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
					RpcMethod: "NewPolicy",
					Use:       "new-policy [name] [policy] --btl [btl]",
					Short:     "Send a new-policy tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
						{ProtoField: "policy"},
					},
				},
				{
					RpcMethod: "RevokeAction",
					Use:       "revoke-action [action-id]",
					Short:     "Send a revoke-action tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "action_id"},
					},
				},
				{
					RpcMethod: "ApproveAction",
					Use:       "approve-action [action-type] [action-id]",
					Short:     "Send a approve-action tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "action_type"},
						{ProtoField: "action_id"},
					},
				},
				{
					RpcMethod: "AddSignMethod",
					Use:       "add-sign-method [id] [config]",
					Short:     "Send a add-sign-method tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "config"},
					},
				},
				{
					RpcMethod: "RemoveSignMethod",
					Use:       "remove-sign-method [id]",
					Short:     "Send a remove-sign-method tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod:      "AddMultiGrant",
					Use:            "add-multi-grant [grantee] [msgs]",
					Short:          "Send a MsgAddMultiGrant tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "grantee"}, {ProtoField: "msgs"}},
				},
				{
					RpcMethod:      "RemoveMultiGrant",
					Use:            "remove-multi-grant [grantee] [msgs]",
					Short:          "Send a MsgRemoveMultiGrant tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "grantee"}, {ProtoField: "msgs"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
