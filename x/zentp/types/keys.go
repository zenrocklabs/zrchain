package types

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
)



func KeyPrefix(p string) []byte {
    return []byte(p)
}
