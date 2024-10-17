package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "policy"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_policy"

	// Version defines the current version the IBC module supports
	Version = "policy-1"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// PortID is the default port id that module binds to
	PortID = "policy"

	SignMethodKey = "signmethod/value/%s/"
)

var (
	ActionsKey     = collections.NewPrefix(0)
	ActionCountKey = collections.NewPrefix(1)
	PoliciesKey    = collections.NewPrefix(2)
	PolicyCountKey = collections.NewPrefix(3)
	SignMethodsKey = collections.NewPrefix(4)
	ParamsKey      = collections.NewPrefix(5)

	ActionsIndex     = "actions"
	ActionCountIndex = "action_count"
	PoliciesIndex    = "policies"
	PolicyCountIndex = "policy_count"
	SignMethodsIndex = "sign_methods"
	ParamsIndex      = "params"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("policy-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
