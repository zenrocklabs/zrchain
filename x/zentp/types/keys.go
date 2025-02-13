package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "zentp"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_zentp"
)

var (
	ParamsKey = []byte("p_zentp")

	BurnsKey = collections.NewPrefix(0)
	MintsKey = collections.NewPrefix(1)

	BurnsIndex = "burns"
	MintsIndex = "mints"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
