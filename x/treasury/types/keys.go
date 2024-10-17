package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "treasury"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_treasury"

	// Version defines the current version the IBC module supports
	Version = "treasury-1"

	// PortID is the default port id that module binds to
	PortID = "treasury"

	TxValueKey = "TXVALUE"

	DataForSigningKey = "DataForSigning"
)

var (
	KeysKey                        = collections.NewPrefix(0)
	KeyCountKey                    = collections.NewPrefix(1)
	WalletsKey                     = collections.NewPrefix(2)
	WalletCountKey                 = collections.NewPrefix(3)
	KeyRequestsKey                 = collections.NewPrefix(4)
	KeyRequestCountKey             = collections.NewPrefix(5)
	SignRequestsKey                = collections.NewPrefix(6)
	SignRequestCountKey            = collections.NewPrefix(7)
	SignTransactionRequestsKey     = collections.NewPrefix(8)
	SignTransactionRequestCountKey = collections.NewPrefix(9)
	ICATransactionRequestsKey      = collections.NewPrefix(10)
	ICATransactionRequestCountKey  = collections.NewPrefix(11)
	ParamsKey                      = collections.NewPrefix(12)

	KeysIndex                        = "keys"
	KeyCountIndex                    = "key_count"
	WalletsIndex                     = "wallets"
	WalletCountIndex                 = "wallet_count"
	KeyRequestsIndex                 = "key_requests"
	KeyRequestCountIndex             = "key_request_count"
	SignRequestsIndex                = "sign_requests"
	SignRequestCountIndex            = "sign_request_count"
	SignTransactionRequestsIndex     = "sign_transaction_requests"
	SignTransactionRequestCountIndex = "sign_transaction_request_count"
	ICATransactionRequestsIndex      = "ica_transaction_requests"
	ICATransactionRequestCountIndex  = "ica_transaction_request_count"
	ParamsIndex                      = "params"

	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("treasury-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
