package zenex

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/Zenrock-Foundation/zrchain/v6/api/zrchain/zenex"
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
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "Swap",
					Use:       "swap [workspace] [rockbtc|btcrock] [amount_in] [yield] [sender_key] [recipient_key]",
					Short:     "Send a swap tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "workspace"},
						{ProtoField: "pair"},
						{ProtoField: "amount_in"},
						{ProtoField: "yield"},
						{ProtoField: "sender_key"},
						{ProtoField: "recipient_key"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
