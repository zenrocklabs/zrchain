package treasury

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/Zenrock-Foundation/zrchain/v5/api/zrchain/treasury"
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
					RpcMethod: "KeyRequests",
					Use:       "key-requests --keyring-addr [keyring-addr] --status [pending|partial|fulfilled|rejected] --workspace-addr [workspace-addr]",
					Short:     "Query KeyRequests, optionally filtering by their keyring address, current status, and workspace address",
				},
				{
					RpcMethod: "Keys",
					Use:       "keys --wallet-type [wallet-type] --workspace-addr [workspace-addr] --prefixes [prefixes]",
					Short:     "Query Keys, optionally by workspace address, deriving wallets for specified type/prefixes",
				},
				{
					RpcMethod: "KeyByID",
					Use:       "key-by-id [id] --wallet-type [wallet-type] --prefixes [prefixes]",
					Short:     "Query Key by ID, deriving wallets for specified type/prefixes",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "SignatureRequests",
					Use:       "signature-requests --keyring-addr [keyring-addr] --status [pending|partial|fulfilled|rejected]",
					Short:     "Query SignatureRequests, optionally filtering by keyring address, request status",
				},
				{
					RpcMethod: "SignatureRequestByID",
					Use:       "signature-request-by-id [id]",
					Short:     "Query SignatureRequests by ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "SignTransactionRequests",
					Use:       "sign-transaction-requests --key-id [key-id] --wallet_type [wallet-type] --status [status]",
					Short:     "Query SignTransactionRequests, optionally filtering key id, wallet type, or request status",
				},
				{
					RpcMethod: "SignTransactionRequestByID",
					Use:       "sign-transaction-request-by-id [id]",
					Short:     "Query SignTransactionRequests by ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "ZrSignKeys",
					Use:       "zr-sign-keys [address] --wallet_type [wallet-type]",
					Short:     "Query ZrSignKeys, filtering by address and optionally by wallet type",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "KeyByAddress",
					Use:       "key-by-address [address] --keyring-addr [keyring-addr] --key-type [key-type] --wallet-type [wallet-type] --prefixes [prefixes]",
					Short:     "Query KeyByAddress, optionally filtering by keyring address, key type, wallet type, and deriving wallets for specified type/prefixes",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "ZenbtcWallets",
					Use:       "zenbtc-wallets --recipient-addr [recipient-addr] --chain-type [chain-type] --chain-id [chain-id] --return-addr [return-addr]",
					Short:     "query-zenbtc-wallets optionally filtering by recipient address, chain type, chain id, and return address",
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
					RpcMethod: "NewKeyRequest",
					Use:       `new-key-request [workspace-addr] [keyring-addr] [secp256k1|ed25519|bitcoin] --btl [btl] --sign-policy-id [policy_id] --zenbtc-metadata [zenbtc-metadata]`,
					Short:     "Broadcast message NewKeyRequest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "workspace_addr"},
						{ProtoField: "keyring_addr"},
						{ProtoField: "key_type"},
					},
				},
				{
					RpcMethod: "FulfilKeyRequest",
					Use:       `fulfil-key-request [id] [pending|partial|fulfilled|rejected] [keyring-party-sig]`,
					Short:     "Broadcast message FulfilKeyRequest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "request_id"},
						{ProtoField: "status"},
						{ProtoField: "keyring_party_signature"},
					},
				},
				{
					RpcMethod: "NewSignatureRequest",
					Use:       "new-signature-request [key-ids] [data-for-signing] --btl [btl] --cache-id [cache-id]",
					Short:     "Broadcast message NewSignatureRequest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "key_ids"},
						{ProtoField: "data_for_signing"},
					},
				},
				{
					RpcMethod: "FulfilSignatureRequest",
					Use:       "fulfil-signature-request [request-id] [partial|fulfilled|rejected] [keyring-party-sig] --signed-data [signed-data] --reject-reason [reject-reason]",
					Short:     "Broadcast message FulfilSignatureRequest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "request_id"},
						{ProtoField: "status"},
						{ProtoField: "keyring_party_signature"},
					},
				},
				{
					RpcMethod: "NewSignTransactionRequest",
					Use:       "new-sign-transaction-request [key-id] [wallet-type] [unsigned-tx] --metadata [metadata] --btl [btl] --cache-id [cache-id]",
					Short:     "Broadcast message NewSignTransactionRequest",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "key_id"},
						{ProtoField: "wallet_type"},
						{ProtoField: "unsigned_transaction"},
					},
				},
				{
					RpcMethod: "TransferFromKeyring",
					Use:       "transfer-from-keyring [keyring] [recipient] [amount] [denom]",
					Short:     "Send a transfer-from-keyring tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "keyring"},
						{ProtoField: "recipient"},
						{ProtoField: "amount"},
						{ProtoField: "denom"},
					},
				},
				{
					RpcMethod:      "NewICATransactionRequest",
					Use:            "new-ica-transaction-request",
					Short:          "Send a new-ica-transaction-request tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod:      "FulfilICATransactionRequest",
					Use:            "fulfil-ica-transaction-request",
					Short:          "Send a fulfil-ica-transaction-request tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod: "NewZrSignSignatureRequest",
					Use:       "new-zr-sign-signature-request [address] [wallet-type] [wallet-index]",
					Short:     "Send a new-zr-sign-signature-request tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
						{ProtoField: "wallet_type"},
						{ProtoField: "wallet_index"},
					},
				},
				{
					RpcMethod: "UpdateKeyPolicy",
					Use:       "update-key-policy [key-id] [sign_policy_id]",
					Short:     "Send a update-key-policy tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "key_id"},
						{ProtoField: "sign_policy_id"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
