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
	BurnsKeyDeprecated  = collections.NewPrefix(0)
	MintsKeyDeprecated  = collections.NewPrefix(1)
	MintCountKey        = collections.NewPrefix(2)
	BurnCountKey        = collections.NewPrefix(3)
	ParamsKey           = collections.NewPrefix(4)
	SolanaROCKSupplyKey = collections.NewPrefix(5)
	ZentpFeesKey        = collections.NewPrefix(7)
	BurnsKey            = collections.NewPrefix(8)
	MintsKey            = collections.NewPrefix(9)
	DctStoreKey         = collections.NewPrefix(10)
	DctMintStoreKey     = collections.NewPrefix(11)
	DctBurnStoreKey     = collections.NewPrefix(12)
	DctMintCountKey     = collections.NewPrefix(13)
	DctBurnCountKey     = collections.NewPrefix(14)

	BurnsIndexDeprecated  = "burns"
	MintsIndexDeprecated  = "mints"
	MintCountIndex        = "mint_count"
	BurnCountIndex        = "burn_count"
	ParamsIndex           = "params"
	SolanaROCKSupplyIndex = "solana_rock_supply"
	ZentpFeesIndex        = "zentp_fees"
	BurnsIndex            = "burns_v2"
	MintsIndex            = "mints_v2"
	DctStoreIndex         = "dct_store"
	DctMintStoreIndex     = "dct_mint_store"
	DctBurnStoreIndex     = "dct_burn_store"
	DctMintCountIndex     = "dct_mint_count"
	DctBurnCountIndex     = "dct_burn_count"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
