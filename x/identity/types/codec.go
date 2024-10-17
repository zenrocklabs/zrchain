package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewWorkspace{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddWorkspaceOwner{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAppendChildWorkspace{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewChildWorkspace{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveWorkspaceOwner{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddKeyringParty{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateKeyring{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveKeyringParty{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewKeyring{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddKeyringAdmin{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveKeyringAdmin{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateWorkspace{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeactivateKeyring{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// ModuleCdc references the global x/ibc-transfer module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/ibc transfer and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
