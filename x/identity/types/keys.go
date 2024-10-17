package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "identity"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_identity"

	// Version defines the current version the IBC module supports
	Version = "identity-1"

	// PortID is the default port id that module binds to
	PortID = "identity"
)

var (
	KeyringsKey       = collections.NewPrefix(0)
	KeyringCountKey   = collections.NewPrefix(1)
	WorkspacesKey     = collections.NewPrefix(2)
	WorkspaceCountKey = collections.NewPrefix(3)
	ParamsKey         = collections.NewPrefix(4)

	KeyringsIndex       = "keyrings"
	KeyringCountIndex   = "keyring_count"
	WorkspacesIndex     = "workspaces"
	WorkspaceCountIndex = "workspace_count"
	ParamsIndex         = "params"

	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("identity-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
