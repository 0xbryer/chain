package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&TextRequestingSignature{}, "tss/TextRequestingSignature", nil)

	legacy.RegisterAminoMsg(cdc, &MsgCreateGroup{}, "tss/CreateGroup")
	legacy.RegisterAminoMsg(cdc, &MsgReplaceGroup{}, "tss/ReplaceGroup")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateGroupFee{}, "tss/UpdateGroupFee")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitDKGRound1{}, "tss/SubmitDKGRound1")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitDKGRound2{}, "tss/SubmitDKGRound2")
	legacy.RegisterAminoMsg(cdc, &MsgComplain{}, "tss/Complaint")
	legacy.RegisterAminoMsg(cdc, &MsgConfirm{}, "tss/Confirm")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitDEs{}, "tss/SubmitDEs")
	legacy.RegisterAminoMsg(cdc, &MsgRequestSignature{}, "tss/RequestSignature")
	legacy.RegisterAminoMsg(cdc, &MsgSubmitSignature{}, "tss/SubmitSignature")
	legacy.RegisterAminoMsg(cdc, &MsgActivate{}, "tss/Activate")
	legacy.RegisterAminoMsg(cdc, &MsgHealthCheck{}, "tss/HealthCheck")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "tss/UpdateParams")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateGroup{},
		&MsgReplaceGroup{},
		&MsgUpdateGroupFee{},
		&MsgSubmitDKGRound1{},
		&MsgSubmitDKGRound2{},
		&MsgComplain{},
		&MsgConfirm{},
		&MsgSubmitDEs{},
		&MsgRequestSignature{},
		&MsgSubmitSignature{},
		&MsgActivate{},
		&MsgHealthCheck{},
		&MsgUpdateParams{},
	)
	registry.RegisterInterface(
		"tss.v1beta1.Content",
		(*Content)(nil),
		&TextRequestingSignature{},
	)
}

// RegisterRequestSignatureTypeCodec registers an external request signature content type defined
// in another module for the internal ModuleCdc. This allows the MsgRequestSignature
// to be correctly Amino encoded and decoded.
//
// NOTE: This should only be used for applications that are still using a concrete
// Amino codec for serialization.
func RegisterRequestSignatureTypeCodec(o interface{}, name string) {
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
