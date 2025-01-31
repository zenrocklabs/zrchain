package zentp

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/Zenrock-Foundation/zrchain/v5/api/zrchain/zentp"
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
					RpcMethod: "MintRock",
					Use:       `mint-rock [amount] [destination] [recipient]`,
					Short:     "Mint new ROCK tokens on destination chain",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
						{ProtoField: "destination"},
						{ProtoField: "recipient"},
					},
				},
				{
					RpcMethod:      "Burn",
					Use:            "burn [address] [denom] [amount]",
					Short:          "Send a burn tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "denom"}, {ProtoField: "amount"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
