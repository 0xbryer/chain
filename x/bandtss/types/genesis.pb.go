// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bandtss/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	github_com_bandprotocol_chain_v2_pkg_tss "github.com/bandprotocol/chain/v2/pkg/tss"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the bandtss module's genesis state.
type GenesisState struct {
	// params defines all the paramiters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// members is an array containing member information.
	Members []Member `protobuf:"bytes,2,rep,name=members,proto3" json:"members"`
	// current_group_id is the current group id of the module.
	CurrentGroupID github_com_bandprotocol_chain_v2_pkg_tss.GroupID `protobuf:"varint,3,opt,name=current_group_id,json=currentGroupId,proto3,casttype=github.com/bandprotocol/chain/v2/pkg/tss.GroupID" json:"current_group_id,omitempty"`
	// signing_count is the number of signings.
	SigningCount uint64 `protobuf:"varint,4,opt,name=signing_count,json=signingCount,proto3" json:"signing_count,omitempty"`
	// signings is the bandtss signing info.
	Signings []Signing `protobuf:"bytes,5,rep,name=signings,proto3" json:"signings"`
	// signing_request_mappings is the list of mapping between tss signing id and bandtss signing id.
	SigningIDMappings []SigningIDMappingGenesis `protobuf:"bytes,6,rep,name=signing_id_mappings,json=signingIdMappings,proto3" json:"signing_id_mappings"`
	// replacement is the replacement information of the current group and new group.
	Replacement Replacement `protobuf:"bytes,7,opt,name=replacement,proto3" json:"replacement"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_521cec35d8e53d42, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetMembers() []Member {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *GenesisState) GetCurrentGroupID() github_com_bandprotocol_chain_v2_pkg_tss.GroupID {
	if m != nil {
		return m.CurrentGroupID
	}
	return 0
}

func (m *GenesisState) GetSigningCount() uint64 {
	if m != nil {
		return m.SigningCount
	}
	return 0
}

func (m *GenesisState) GetSignings() []Signing {
	if m != nil {
		return m.Signings
	}
	return nil
}

func (m *GenesisState) GetSigningIDMappings() []SigningIDMappingGenesis {
	if m != nil {
		return m.SigningIDMappings
	}
	return nil
}

func (m *GenesisState) GetReplacement() Replacement {
	if m != nil {
		return m.Replacement
	}
	return Replacement{}
}

// Params defines the set of module parameters.
type Params struct {
	// active_duration is the duration where a member can be active without interaction.
	ActiveDuration time.Duration `protobuf:"bytes,1,opt,name=active_duration,json=activeDuration,proto3,stdduration" json:"active_duration"`
	// reward_percentage is the percentage of block rewards allocated to active TSS validators after being allocated to
	// oracle rewards.
	RewardPercentage uint64 `protobuf:"varint,2,opt,name=reward_percentage,json=rewardPercentage,proto3" json:"reward_percentage,omitempty"`
	// inactive_penalty_duration is the duration where a member cannot activate back after inactive.
	InactivePenaltyDuration time.Duration `protobuf:"bytes,3,opt,name=inactive_penalty_duration,json=inactivePenaltyDuration,proto3,stdduration" json:"inactive_penalty_duration"`
	// jail_penalty_duration is the duration where a member cannot activate back after jail.
	JailPenaltyDuration time.Duration `protobuf:"bytes,4,opt,name=jail_penalty_duration,json=jailPenaltyDuration,proto3,stdduration" json:"jail_penalty_duration"`
	// fee is the tokens that will be paid per signing.
	Fee github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,5,rep,name=fee,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"fee"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_521cec35d8e53d42, []int{1}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetActiveDuration() time.Duration {
	if m != nil {
		return m.ActiveDuration
	}
	return 0
}

func (m *Params) GetRewardPercentage() uint64 {
	if m != nil {
		return m.RewardPercentage
	}
	return 0
}

func (m *Params) GetInactivePenaltyDuration() time.Duration {
	if m != nil {
		return m.InactivePenaltyDuration
	}
	return 0
}

func (m *Params) GetJailPenaltyDuration() time.Duration {
	if m != nil {
		return m.JailPenaltyDuration
	}
	return 0
}

func (m *Params) GetFee() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Fee
	}
	return nil
}

// SigningIDMappingGenesis defines the genesis state for the signingIDMapping.
type SigningIDMappingGenesis struct {
	// signing_id is the tss signing id.
	SigningID github_com_bandprotocol_chain_v2_pkg_tss.SigningID `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3,casttype=github.com/bandprotocol/chain/v2/pkg/tss.SigningID" json:"signing_id,omitempty"`
	// bandtss_signing_id is the signing id being referred in bandtss module.
	BandtssSigningID SigningID `protobuf:"varint,2,opt,name=bandtss_signing_id,json=bandtssSigningId,proto3,casttype=SigningID" json:"bandtss_signing_id,omitempty"`
}

func (m *SigningIDMappingGenesis) Reset()         { *m = SigningIDMappingGenesis{} }
func (m *SigningIDMappingGenesis) String() string { return proto.CompactTextString(m) }
func (*SigningIDMappingGenesis) ProtoMessage()    {}
func (*SigningIDMappingGenesis) Descriptor() ([]byte, []int) {
	return fileDescriptor_521cec35d8e53d42, []int{2}
}
func (m *SigningIDMappingGenesis) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SigningIDMappingGenesis) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SigningIDMappingGenesis.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SigningIDMappingGenesis) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigningIDMappingGenesis.Merge(m, src)
}
func (m *SigningIDMappingGenesis) XXX_Size() int {
	return m.Size()
}
func (m *SigningIDMappingGenesis) XXX_DiscardUnknown() {
	xxx_messageInfo_SigningIDMappingGenesis.DiscardUnknown(m)
}

var xxx_messageInfo_SigningIDMappingGenesis proto.InternalMessageInfo

func (m *SigningIDMappingGenesis) GetSigningID() github_com_bandprotocol_chain_v2_pkg_tss.SigningID {
	if m != nil {
		return m.SigningID
	}
	return 0
}

func (m *SigningIDMappingGenesis) GetBandtssSigningID() SigningID {
	if m != nil {
		return m.BandtssSigningID
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "bandtss.v1beta1.GenesisState")
	proto.RegisterType((*Params)(nil), "bandtss.v1beta1.Params")
	proto.RegisterType((*SigningIDMappingGenesis)(nil), "bandtss.v1beta1.SigningIDMappingGenesis")
}

func init() { proto.RegisterFile("bandtss/v1beta1/genesis.proto", fileDescriptor_521cec35d8e53d42) }

var fileDescriptor_521cec35d8e53d42 = []byte{
	// 673 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x4f, 0x4f, 0xd4, 0x4c,
	0x18, 0xdf, 0xb2, 0xfb, 0x2e, 0x30, 0xf0, 0xc2, 0x32, 0x60, 0xe8, 0x12, 0xdd, 0x22, 0x5e, 0xf6,
	0x62, 0x0b, 0x18, 0x63, 0xe2, 0xcd, 0xee, 0x46, 0x82, 0x91, 0x48, 0xca, 0xc1, 0xc4, 0xc4, 0x34,
	0xd3, 0x76, 0x28, 0x95, 0xed, 0x4c, 0xd3, 0x99, 0x45, 0xf9, 0x0e, 0x1e, 0x3c, 0xfa, 0x19, 0x3c,
	0xfa, 0x29, 0x38, 0x19, 0x8e, 0x9e, 0x8a, 0x29, 0xdf, 0x82, 0x93, 0xe9, 0xcc, 0xb4, 0xbb, 0xee,
	0x4a, 0xe4, 0x04, 0x33, 0xbf, 0x3f, 0xcf, 0x33, 0x7d, 0x7e, 0xcf, 0x82, 0x07, 0x1e, 0x22, 0x01,
	0x67, 0xcc, 0x3a, 0xdb, 0xf1, 0x30, 0x47, 0x3b, 0x56, 0x88, 0x09, 0x66, 0x11, 0x33, 0x93, 0x94,
	0x72, 0x0a, 0x97, 0x15, 0x6c, 0x2a, 0x78, 0x63, 0x2d, 0xa4, 0x21, 0x15, 0x98, 0x55, 0xfc, 0x27,
	0x69, 0x1b, 0x9d, 0x90, 0xd2, 0x70, 0x80, 0x2d, 0x71, 0xf2, 0x86, 0xc7, 0x56, 0x30, 0x4c, 0x11,
	0x8f, 0x28, 0x51, 0xf8, 0x54, 0x95, 0xd2, 0x56, 0xc9, 0x7d, 0xca, 0x62, 0xca, 0x2c, 0x0f, 0x31,
	0x5c, 0x51, 0x7c, 0x1a, 0x95, 0xf2, 0xb6, 0xc4, 0x5d, 0x59, 0x57, 0x1e, 0x24, 0xb4, 0xf5, 0xb9,
	0x01, 0x16, 0xf7, 0x64, 0xcb, 0x47, 0x1c, 0x71, 0x0c, 0x9f, 0x82, 0x66, 0x82, 0x52, 0x14, 0x33,
	0x5d, 0xdb, 0xd4, 0xba, 0x0b, 0xbb, 0xeb, 0xe6, 0xc4, 0x13, 0xcc, 0x43, 0x01, 0xdb, 0x8d, 0x8b,
	0xcc, 0xa8, 0x39, 0x8a, 0x0c, 0x9f, 0x81, 0xd9, 0x18, 0xc7, 0x1e, 0x4e, 0x99, 0x3e, 0xb3, 0x59,
	0xff, 0xab, 0xee, 0x40, 0xe0, 0x4a, 0x57, 0xb2, 0x61, 0x02, 0x5a, 0xfe, 0x30, 0x4d, 0x31, 0xe1,
	0x6e, 0x98, 0xd2, 0x61, 0xe2, 0x46, 0x81, 0x5e, 0xdf, 0xd4, 0xba, 0x0d, 0xfb, 0x65, 0x9e, 0x19,
	0x4b, 0x3d, 0x89, 0xed, 0x15, 0xd0, 0x7e, 0xff, 0x26, 0x33, 0xb6, 0xc3, 0x88, 0x9f, 0x0c, 0x3d,
	0xd3, 0xa7, 0xb1, 0xf8, 0x0a, 0xe2, 0x19, 0x3e, 0x1d, 0x58, 0xfe, 0x09, 0x8a, 0x88, 0x75, 0xb6,
	0x6b, 0x25, 0xa7, 0xa1, 0x55, 0xd4, 0x55, 0x1a, 0x67, 0xc9, 0x1f, 0xf7, 0x08, 0xe0, 0x23, 0xf0,
	0x3f, 0x8b, 0x42, 0x12, 0x91, 0xd0, 0xf5, 0xe9, 0x90, 0x70, 0xbd, 0x51, 0x94, 0x73, 0x16, 0xd5,
	0x65, 0xaf, 0xb8, 0x83, 0xcf, 0xc1, 0x9c, 0x3a, 0x33, 0xfd, 0x3f, 0xf1, 0x20, 0x7d, 0xea, 0x41,
	0x47, 0x92, 0xa0, 0x5e, 0x54, 0xf1, 0x21, 0x03, 0xab, 0x65, 0x81, 0x28, 0x70, 0x63, 0x94, 0x24,
	0xc2, 0xa6, 0x29, 0x6c, 0xba, 0xb7, 0xd9, 0xec, 0xf7, 0x0f, 0x24, 0x53, 0x8d, 0xc3, 0x6e, 0x17,
	0xb6, 0x79, 0x66, 0xac, 0x4c, 0x12, 0x98, 0xb3, 0xa2, 0xfc, 0xf7, 0x83, 0xf2, 0x0a, 0xf6, 0xc1,
	0x42, 0x8a, 0x93, 0x01, 0xf2, 0x71, 0x8c, 0x09, 0xd7, 0x67, 0xc5, 0xf0, 0xee, 0x4f, 0x15, 0x73,
	0x46, 0x1c, 0xd5, 0xf7, 0xb8, 0x6c, 0xeb, 0x7b, 0x1d, 0x34, 0xe5, 0x7c, 0xe1, 0x6b, 0xb0, 0x8c,
	0x7c, 0x1e, 0x9d, 0x61, 0xb7, 0x0c, 0xa3, 0x4a, 0x44, 0xdb, 0x94, 0x69, 0x35, 0xcb, 0xb4, 0x9a,
	0x7d, 0x45, 0xb0, 0xe7, 0x0a, 0xc7, 0xaf, 0x57, 0x86, 0xe6, 0x2c, 0x49, 0x6d, 0x89, 0xc0, 0x17,
	0x60, 0x25, 0xc5, 0x1f, 0x51, 0x1a, 0xb8, 0x09, 0x4e, 0x7d, 0x4c, 0x38, 0x0a, 0xb1, 0x3e, 0x23,
	0xe6, 0xbc, 0x96, 0x67, 0x46, 0xcb, 0x11, 0xe0, 0x61, 0x85, 0x39, 0xad, 0x74, 0xe2, 0x06, 0xba,
	0xa0, 0x1d, 0x11, 0xd5, 0x52, 0x82, 0x09, 0x1a, 0xf0, 0xf3, 0x51, 0x6b, 0xf5, 0xbb, 0xb7, 0xb6,
	0x5e, 0xba, 0x1c, 0x4a, 0x93, 0xaa, 0xc7, 0xb7, 0xe0, 0xde, 0x07, 0x14, 0x0d, 0xa6, 0xcd, 0x1b,
	0x77, 0x37, 0x5f, 0x2d, 0x1c, 0x26, 0x8d, 0xdf, 0x83, 0xfa, 0x31, 0xc6, 0x2a, 0x47, 0x6d, 0x53,
	0x2d, 0x60, 0xb1, 0xad, 0xd5, 0x5c, 0x7a, 0x34, 0x22, 0xf6, 0x76, 0x61, 0xf3, 0xed, 0xca, 0xe8,
	0x8e, 0x65, 0x5c, 0xad, 0xb6, 0xfc, 0xf3, 0x98, 0x05, 0xa7, 0x16, 0x3f, 0x4f, 0x30, 0x13, 0x02,
	0xe6, 0x14, 0xbe, 0x5b, 0x3f, 0x34, 0xb0, 0x7e, 0x4b, 0x88, 0xa0, 0x07, 0xc0, 0x28, 0x8b, 0x62,
	0x80, 0x0d, 0xbb, 0x97, 0x67, 0xc6, 0x7c, 0x25, 0xb8, 0xc9, 0x8c, 0xdd, 0x3b, 0xef, 0x54, 0xa5,
	0x72, 0xe6, 0xab, 0x08, 0xc2, 0x37, 0x00, 0xaa, 0x98, 0xb9, 0x63, 0xb5, 0xe4, 0x70, 0x1f, 0x16,
	0xc3, 0xb5, 0x25, 0x3a, 0x5e, 0x72, 0x54, 0xdf, 0x69, 0x79, 0x7f, 0xc2, 0x81, 0xfd, 0xea, 0x22,
	0xef, 0x68, 0x97, 0x79, 0x47, 0xfb, 0x95, 0x77, 0xb4, 0x2f, 0xd7, 0x9d, 0xda, 0xe5, 0x75, 0xa7,
	0xf6, 0xf3, 0xba, 0x53, 0x7b, 0xf7, 0xef, 0xed, 0xff, 0x54, 0xfe, 0x36, 0xca, 0xef, 0xe4, 0x35,
	0x05, 0xe5, 0xc9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3e, 0xdf, 0xa7, 0x7b, 0xa9, 0x05, 0x00,
	0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Replacement.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if len(m.SigningIDMappings) > 0 {
		for iNdEx := len(m.SigningIDMappings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SigningIDMappings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Signings) > 0 {
		for iNdEx := len(m.Signings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Signings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if m.SigningCount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SigningCount))
		i--
		dAtA[i] = 0x20
	}
	if m.CurrentGroupID != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.CurrentGroupID))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Members) > 0 {
		for iNdEx := len(m.Members) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Members[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Fee) > 0 {
		for iNdEx := len(m.Fee) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Fee[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	n3, err3 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.JailPenaltyDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.JailPenaltyDuration):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintGenesis(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x22
	n4, err4 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.InactivePenaltyDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.InactivePenaltyDuration):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintGenesis(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x1a
	if m.RewardPercentage != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.RewardPercentage))
		i--
		dAtA[i] = 0x10
	}
	n5, err5 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.ActiveDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ActiveDuration):])
	if err5 != nil {
		return 0, err5
	}
	i -= n5
	i = encodeVarintGenesis(dAtA, i, uint64(n5))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *SigningIDMappingGenesis) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SigningIDMappingGenesis) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SigningIDMappingGenesis) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BandtssSigningID != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.BandtssSigningID))
		i--
		dAtA[i] = 0x10
	}
	if m.SigningID != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SigningID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Members) > 0 {
		for _, e := range m.Members {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.CurrentGroupID != 0 {
		n += 1 + sovGenesis(uint64(m.CurrentGroupID))
	}
	if m.SigningCount != 0 {
		n += 1 + sovGenesis(uint64(m.SigningCount))
	}
	if len(m.Signings) > 0 {
		for _, e := range m.Signings {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.SigningIDMappings) > 0 {
		for _, e := range m.SigningIDMappings {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.Replacement.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ActiveDuration)
	n += 1 + l + sovGenesis(uint64(l))
	if m.RewardPercentage != 0 {
		n += 1 + sovGenesis(uint64(m.RewardPercentage))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.InactivePenaltyDuration)
	n += 1 + l + sovGenesis(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.JailPenaltyDuration)
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Fee) > 0 {
		for _, e := range m.Fee {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *SigningIDMappingGenesis) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningID != 0 {
		n += 1 + sovGenesis(uint64(m.SigningID))
	}
	if m.BandtssSigningID != 0 {
		n += 1 + sovGenesis(uint64(m.BandtssSigningID))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Members", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Members = append(m.Members, Member{})
			if err := m.Members[len(m.Members)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentGroupID", wireType)
			}
			m.CurrentGroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentGroupID |= github_com_bandprotocol_chain_v2_pkg_tss.GroupID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningCount", wireType)
			}
			m.SigningCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signings", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signings = append(m.Signings, Signing{})
			if err := m.Signings[len(m.Signings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningIDMappings", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SigningIDMappings = append(m.SigningIDMappings, SigningIDMappingGenesis{})
			if err := m.SigningIDMappings[len(m.SigningIDMappings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Replacement", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Replacement.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActiveDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.ActiveDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardPercentage", wireType)
			}
			m.RewardPercentage = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardPercentage |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InactivePenaltyDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.InactivePenaltyDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JailPenaltyDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.JailPenaltyDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Fee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Fee = append(m.Fee, types.Coin{})
			if err := m.Fee[len(m.Fee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SigningIDMappingGenesis) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SigningIDMappingGenesis: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SigningIDMappingGenesis: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningID", wireType)
			}
			m.SigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningID |= github_com_bandprotocol_chain_v2_pkg_tss.SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BandtssSigningID", wireType)
			}
			m.BandtssSigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BandtssSigningID |= SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
