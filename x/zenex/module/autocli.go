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
				{
					RpcMethod: "Swaps",
					Use:       "swaps",
					Short:     "Returns swap objects",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"creator":        {Usage: "Filter by creator address"},
						"swap_id":        {Usage: "Filter by swap ID"},
						"status":         {Usage: "Filter by status (initiated|requested|completed|rejected)"},
						"workspace":      {Usage: "Filter by workspace"},
						"pair":           {Usage: "Filter by pair (rock-btc|btc-rock)"},
						"source_tx_hash": {Usage: "Filter by source transaction hash"},
					},
				},

				{
					RpcMethod:      "RockPool",
					Use:            "rock-pool",
					Short:          "Query rock-pool",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
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
					Use:       "update-params",
					Short:     "Update the parameters of the module",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "params"},
					},
				},
				{
					RpcMethod: "SwapRequest",
					Use:       "swap [workspace] [rock-btc|btc-rock] [amount_in] [rock_key_id] [btc_key_id]",
					Short:     "Send a swap tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "workspace"},
						{ProtoField: "pair"},
						{ProtoField: "amount_in"},
						{ProtoField: "rock_key_id"},
						{ProtoField: "btc_key_id"},
					},
				},
				{
					RpcMethod: "ZenexTransferRequest",
					Use:       "zenex-transfer-request [swap-id] [unsigned-plus-tx] [wallet-type] [cache-id] [data-for-signing]",
					Short:     "Send a zenex-transfer-request tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "swap_id"},
						{ProtoField: "unsigned_plus_tx"},
						{ProtoField: "wallet_type"},
						{ProtoField: "cache_id"},
						{ProtoField: "data_for_signing"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"reject_reason": {Usage: "Optional reason for rejection (only used when status is rejected)"},
					},
				},
				{
					RpcMethod: "AcknowledgePoolTransfer",
					Use:       "acknowledge-pool-transfer [swap-id] [source-tx-hash] [status]",
					Short:     "Send a AcknowledgePoolTransfer tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "swap_id"},
						{ProtoField: "source_tx_hash"},
						{ProtoField: "status"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"reject_reason": {Usage: "Optional reason for rejection (only used when status is rejected)"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
