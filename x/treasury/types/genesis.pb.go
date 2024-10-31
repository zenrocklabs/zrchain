// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zrchain/treasury/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the treasury module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params         Params                   `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	PortId         string                   `protobuf:"bytes,2,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	Keys           []Key                    `protobuf:"bytes,3,rep,name=keys,proto3" json:"keys"`
	KeyRequests    []KeyRequest             `protobuf:"bytes,4,rep,name=key_requests,json=keyRequests,proto3" json:"key_requests"`
	SignRequests   []SignRequest            `protobuf:"bytes,5,rep,name=sign_requests,json=signRequests,proto3" json:"sign_requests"`
	SignTxRequests []SignTransactionRequest `protobuf:"bytes,6,rep,name=sign_tx_requests,json=signTxRequests,proto3" json:"sign_tx_requests"`
	IcaTxRequests  []ICATransactionRequest  `protobuf:"bytes,7,rep,name=ica_tx_requests,json=icaTxRequests,proto3" json:"ica_tx_requests"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1b80491bb414464, []int{0}
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

func (m *GenesisState) GetPortId() string {
	if m != nil {
		return m.PortId
	}
	return ""
}

func (m *GenesisState) GetKeys() []Key {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *GenesisState) GetKeyRequests() []KeyRequest {
	if m != nil {
		return m.KeyRequests
	}
	return nil
}

func (m *GenesisState) GetSignRequests() []SignRequest {
	if m != nil {
		return m.SignRequests
	}
	return nil
}

func (m *GenesisState) GetSignTxRequests() []SignTransactionRequest {
	if m != nil {
		return m.SignTxRequests
	}
	return nil
}

func (m *GenesisState) GetIcaTxRequests() []ICATransactionRequest {
	if m != nil {
		return m.IcaTxRequests
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "zrchain.treasury.GenesisState")
}

func init() { proto.RegisterFile("zrchain/treasury/genesis.proto", fileDescriptor_e1b80491bb414464) }

var fileDescriptor_e1b80491bb414464 = []byte{
	// 408 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0xd2, 0xcf, 0xea, 0xd3, 0x40,
	0x10, 0x07, 0xf0, 0xc4, 0xd6, 0xfc, 0xe8, 0xb6, 0xd5, 0x1a, 0x14, 0x43, 0xb0, 0xb1, 0x78, 0x31,
	0x08, 0x66, 0xa1, 0xe2, 0x41, 0x3c, 0x59, 0xf1, 0x4f, 0xf1, 0x22, 0x69, 0x05, 0xe9, 0xa5, 0x6c,
	0xd3, 0x25, 0x5d, 0x42, 0x76, 0xe3, 0xee, 0x46, 0x1a, 0x9f, 0xc2, 0xc7, 0xf0, 0xe8, 0x63, 0xf4,
	0xd8, 0x9b, 0x9e, 0x44, 0xda, 0x83, 0xaf, 0x21, 0xd9, 0xa4, 0x0d, 0x6d, 0xfa, 0xbb, 0x94, 0xed,
	0x7c, 0x67, 0x3e, 0x4c, 0x60, 0x80, 0xf3, 0x8d, 0x07, 0x2b, 0x44, 0x28, 0x94, 0x1c, 0x23, 0x91,
	0xf2, 0x0c, 0x86, 0x98, 0x62, 0x41, 0x84, 0x97, 0x70, 0x26, 0x99, 0xd9, 0x2b, 0x73, 0xef, 0x90,
	0xdb, 0x77, 0x50, 0x4c, 0x28, 0x83, 0xea, 0xb7, 0x68, 0xb2, 0xef, 0x86, 0x2c, 0x64, 0xea, 0x09,
	0xf3, 0x57, 0x59, 0xed, 0xd7, 0xe8, 0x04, 0x71, 0x14, 0x97, 0xb2, 0x6d, 0xd7, 0xe2, 0x08, 0x67,
	0x65, 0x56, 0xdf, 0x2a, 0x4e, 0x02, 0x41, 0x42, 0x5a, 0xe4, 0x8f, 0x7e, 0x35, 0x40, 0xe7, 0x5d,
	0xb1, 0xe7, 0x44, 0x22, 0x89, 0xcd, 0x97, 0xc0, 0x28, 0x70, 0x4b, 0x1f, 0xe8, 0x6e, 0x7b, 0x68,
	0x79, 0xe7, 0x7b, 0x7b, 0x1f, 0x55, 0x3e, 0x6a, 0x6d, 0xfe, 0x3c, 0xd4, 0x7e, 0xfc, 0xfb, 0xf9,
	0x44, 0xf7, 0xcb, 0x11, 0xf3, 0x3e, 0xb8, 0x4a, 0x18, 0x97, 0x73, 0xb2, 0xb4, 0x6e, 0x0c, 0x74,
	0xb7, 0xe5, 0x1b, 0xf9, 0xdf, 0xf1, 0xd2, 0x84, 0xa0, 0x19, 0xe1, 0x4c, 0x58, 0x8d, 0x41, 0xc3,
	0x6d, 0x0f, 0xef, 0xd5, 0xcd, 0x0f, 0x38, 0x1b, 0x35, 0x73, 0xd0, 0x57, 0x8d, 0xe6, 0x1b, 0xd0,
	0x89, 0x70, 0x36, 0xe7, 0xf8, 0x4b, 0x8a, 0x85, 0x14, 0x56, 0x53, 0x0d, 0x3e, 0xb8, 0x38, 0xe8,
	0x17, 0x4d, 0xe5, 0x7c, 0x3b, 0x3a, 0x56, 0x84, 0xf9, 0x1e, 0x74, 0xf3, 0x8f, 0xad, 0x9c, 0x9b,
	0xca, 0xe9, 0xd7, 0x9d, 0x09, 0x09, 0xe9, 0x29, 0xd4, 0x11, 0x55, 0x49, 0x98, 0x9f, 0x41, 0x4f,
	0x49, 0x72, 0x5d, 0x61, 0x86, 0xc2, 0xdc, 0xcb, 0xd8, 0x94, 0x23, 0x2a, 0x50, 0x20, 0x09, 0x3b,
	0x73, 0x6f, 0xe5, 0xce, 0x74, 0x7d, 0x94, 0x3f, 0x81, 0xdb, 0x24, 0x40, 0x27, 0xf0, 0x95, 0x82,
	0x1f, 0xd7, 0xe1, 0xf1, 0xeb, 0x57, 0xd7, 0xba, 0x5d, 0x12, 0xa0, 0x8a, 0x1d, 0x4d, 0x36, 0x3b,
	0x47, 0xdf, 0xee, 0x1c, 0xfd, 0xef, 0xce, 0xd1, 0xbf, 0xef, 0x1d, 0x6d, 0xbb, 0x77, 0xb4, 0xdf,
	0x7b, 0x47, 0x9b, 0xbd, 0x08, 0x89, 0x5c, 0xa5, 0x0b, 0x2f, 0x60, 0x31, 0x9c, 0x61, 0xca, 0x59,
	0x10, 0x3d, 0x7d, 0xcb, 0x52, 0xba, 0x44, 0x39, 0x0a, 0x0f, 0x17, 0xf3, 0xf5, 0x39, 0x5c, 0x57,
	0x67, 0x23, 0xb3, 0x04, 0x8b, 0x85, 0xa1, 0xae, 0xe6, 0xd9, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x46, 0x36, 0x0f, 0x50, 0xed, 0x02, 0x00, 0x00,
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
	if len(m.IcaTxRequests) > 0 {
		for iNdEx := len(m.IcaTxRequests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.IcaTxRequests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.SignTxRequests) > 0 {
		for iNdEx := len(m.SignTxRequests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SignTxRequests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.SignRequests) > 0 {
		for iNdEx := len(m.SignRequests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SignRequests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.KeyRequests) > 0 {
		for iNdEx := len(m.KeyRequests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.KeyRequests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Keys) > 0 {
		for iNdEx := len(m.Keys) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Keys[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.PortId) > 0 {
		i -= len(m.PortId)
		copy(dAtA[i:], m.PortId)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.PortId)))
		i--
		dAtA[i] = 0x12
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
	l = len(m.PortId)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Keys) > 0 {
		for _, e := range m.Keys {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.KeyRequests) > 0 {
		for _, e := range m.KeyRequests {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.SignRequests) > 0 {
		for _, e := range m.SignRequests {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.SignTxRequests) > 0 {
		for _, e := range m.SignTxRequests {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.IcaTxRequests) > 0 {
		for _, e := range m.IcaTxRequests {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
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
				return fmt.Errorf("proto: wrong wireType = %d for field PortId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PortId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Keys", wireType)
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
			m.Keys = append(m.Keys, Key{})
			if err := m.Keys[len(m.Keys)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyRequests", wireType)
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
			m.KeyRequests = append(m.KeyRequests, KeyRequest{})
			if err := m.KeyRequests[len(m.KeyRequests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignRequests", wireType)
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
			m.SignRequests = append(m.SignRequests, SignRequest{})
			if err := m.SignRequests[len(m.SignRequests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignTxRequests", wireType)
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
			m.SignTxRequests = append(m.SignTxRequests, SignTransactionRequest{})
			if err := m.SignTxRequests[len(m.SignTxRequests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IcaTxRequests", wireType)
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
			m.IcaTxRequests = append(m.IcaTxRequests, ICATransactionRequest{})
			if err := m.IcaTxRequests[len(m.IcaTxRequests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
