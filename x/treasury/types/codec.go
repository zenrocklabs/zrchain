package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*Metadata)(nil),
		&MetadataEthereum{},
	)
	registry.RegisterImplementations((*Metadata)(nil),
		&MetadataSolana{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewKeyRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFulfilKeyRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewSignatureRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFulfilSignatureRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewSignTransactionRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferFromKeyring{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewICATransactionRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFulfilICATransactionRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewZrSignSignatureRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateKeyPolicy{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// ModuleCdc references the global x/ibc-transfer module codec. Note, the codec
// should ONLY be used in certain instances of tests and for JSON encoding.
//
// The actual codec used for serialization should be provided to x/ibc transfer and
// defined at the application level.
var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
