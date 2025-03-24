package types

import (
	"github.com/Zenrock-Foundation/zrchain/v6/policy"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgApproveAction{}, "policy/ApproveAction", nil)
	cdc.RegisterConcrete(&MsgNewPolicy{}, "policy/MsgNewPolicy", nil)
	cdc.RegisterConcrete(&MsgRevokeAction{}, "policy/MsgRevokeAction", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*SignMethod)(nil),
		&SignMethodPasskey{},
	)
	registry.RegisterImplementations((*AdditionalSignature)(nil),
		&AdditionalSignaturePasskey{},
	)
	registry.RegisterImplementations((*policy.Policy)(nil),
		&BoolparserPolicy{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewPolicy{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRevokeAction{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgApproveAction{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddSignMethod{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveSignMethod{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddMultiGrant{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveMultiGrant{},
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
