package types

const (
	// ModuleName defines the module name
	ModuleName = "bedrock"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_bedrock"

    
)

var (
	ParamsKey = []byte("p_bedrock")
)



func KeyPrefix(p string) []byte {
    return []byte(p)
}
