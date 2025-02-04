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
					Use:       `mint-rock [amount] [src-key-id] [dst-chain] [recipient-key-id]`,
					Short:     "Mint new ROCK tokens on destination chain",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
						{ProtoField: "source_key_id"},
						{ProtoField: "destination_chain"},
						{ProtoField: "recipient_key_id"},
					},
				},
				{
					RpcMethod: "Burn",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "BurnRock",
					Use:            "burn-rock [chain-id] [key-id] [amount] [recipient]",
					Short:          "Send a burn_rock tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "chainId"}, {ProtoField: "keyId"}, {ProtoField: "amount"}, {ProtoField: "recipient"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
