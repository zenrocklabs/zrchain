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
				{
					RpcMethod: "Mints",
					Use:       "mints",
					Short:     "Query mints with optional filters",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"id":      {Usage: "Filter by mint ID"},
						"creator": {Usage: "Filter by creator address"},
						"status":  {Usage: "Filter by status"},
						"denom":   {Usage: "Filter by denom"},
						"tx_id":   {Usage: "Filter by transaction ID"},
					},
				},
				{
					RpcMethod: "Burns",
					Use:       "burns",
					Short:     "Query burns with optional filters",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"id":     {Usage: "Filter by burn ID"},
						"denom":  {Usage: "Filter by denom"},
						"status": {Usage: "Filter by status"},
						"tx_id":  {Usage: "Filter by transaction ID"},
					},
				},
				{
					RpcMethod: "QuerySolanaROCKSupply",
					Use:       "solana-rock-supply",
					Short:     "Query the total ROCK supply on Solana",
				},
				{
					RpcMethod: "Stats",
					Use:       "stats",
					Short:     "Query total mints and burns for an address",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"address": {Usage: "The address to query stats for (optional)"},
						"denom":   {Usage: "The denom to query stats for (optional)"},
					},
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
				{
					RpcMethod: "SetSolanaROCKSupply",
					Use:       `set-solana-rock-supply [amount]`,
					Short:     "Set the total ROCK supply on Solana (gov-gated)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
