package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "zenex"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_zenex"
)

var (
	ParamsKey     = collections.NewPrefix(0)
	SwapsKey      = collections.NewPrefix(1)
	SwapsCountKey = collections.NewPrefix(2)

	ParamsIndex     = "params"
	SwapsIndex      = "swaps"
	SwapsCountIndex = "swaps_count"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
