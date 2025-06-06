package zentp

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/Zenrock-Foundation/zrchain/v6/api/zrchain/zentp"
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
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Use:       `update-params [params]`,
					Short:     "Update the parameters of the zentp module",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "params"},
					},
				},
				{
					RpcMethod: "Bridge",
					Use:       `bridge [amount] [denom] [src-address] [dst-chain] [recipient-address]`,
					Short:     "Bridge tokens from zrchain to a destination chain",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
						{ProtoField: "denom"},
						{ProtoField: "source_address"},
						{ProtoField: "destination_chain"},
						{ProtoField: "recipient_address"},
					},
				},
				{
					RpcMethod: "Burn",
					Use:       `burn [module-account] [denom] [amount]`,
					Short:     "Burn tokens from a zrchain module account",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "module_account"},
						{ProtoField: "denom"},
						{ProtoField: "amount"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
