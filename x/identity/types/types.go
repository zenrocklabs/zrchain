package types

const (
	EventNewWorkspace          = "new_workspace"
	EventAddOwnerToWorkspace   = "add_owner_to_workspace"
	EventOwnerAddedToWorkspace = "owner_added_to_workspace"
)

const (
	AttributeWorkspaceAddr = "workspace_addr"
	AttributeOwnerAddr     = "owner_addr"
	AttributeActionId      = "action_id"
)

const (
	PrefixWorkspaceAddress = "workspace"
	PrefixKeyringAddress   = "keyring"
	WorkspaceAddressLength = 10
	KeyringAddressLength   = 11
)
