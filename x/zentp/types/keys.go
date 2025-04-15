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

	BurnsKey     = collections.NewPrefix(0)
	MintsKey     = collections.NewPrefix(1)
	MintCountKey = collections.NewPrefix(2)
	BurnCountKey = collections.NewPrefix(3)

	BurnsIndex     = "burns"
	MintsIndex     = "mints"
	MintCountIndex = "mint_count"
	BurnCountIndex = "burn_count"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
