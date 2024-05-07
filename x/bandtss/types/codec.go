package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers the necessary x/bandtss interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateGroup{}, "bandtss/CreateGroup")
	legacy.RegisterAminoMsg(cdc, &MsgReplaceGroup{}, "bandtss/ReplaceGroup")
	legacy.RegisterAminoMsg(cdc, &MsgRequestSignature{}, "bandtss/RequestSignature")
	legacy.RegisterAminoMsg(cdc, &MsgActivate{}, "bandtss/Activate")
	legacy.RegisterAminoMsg(cdc, &MsgHealthCheck{}, "bandtss/HealthCheck")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "bandtss/UpdateParams")
}

// RegisterInterfaces register the bandtss module interfaces to protobuf Any.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateGroup{},
		&MsgReplaceGroup{},
		&MsgRequestSignature{},
		&MsgActivate{},
		&MsgHealthCheck{},
		&MsgUpdateParams{},
	)
}

// RegisterRequestSignatureTypeCodec registers an external signature request type defined
// in another module for the internal ModuleCdc. This allows the MsgRequestSignature
// to be correctly Amino encoded and decoded.
//
// NOTE: This should only be used for applications that are still using a concrete
// Amino codec for serialization.
func RegisterSignatureOrderTypeCodec(o interface{}, name string) {
	amino.RegisterConcrete(o, name, nil)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
