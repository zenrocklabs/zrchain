// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zrchain/zentp/params.proto

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

// Params defines the parameters for the module.
type Params struct {
	SolanaRelayerKeyId  uint64 `protobuf:"varint,1,opt,name=solana_relayer_key_id,json=solanaRelayerKeyId,proto3" json:"solana_relayer_key_id,omitempty"`
	ZrchainRelayerKeyId uint64 `protobuf:"varint,2,opt,name=zrchain_relayer_key_id,json=zrchainRelayerKeyId,proto3" json:"zrchain_relayer_key_id,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_32fd83cb82e1f0ea, []int{0}
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

func (m *Params) GetSolanaRelayerKeyId() uint64 {
	if m != nil {
		return m.SolanaRelayerKeyId
	}
	return 0
}

func (m *Params) GetZrchainRelayerKeyId() uint64 {
	if m != nil {
		return m.ZrchainRelayerKeyId
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "zrchain.zentp.Params")
}

func init() { proto.RegisterFile("zrchain/zentp/params.proto", fileDescriptor_32fd83cb82e1f0ea) }

var fileDescriptor_32fd83cb82e1f0ea = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xaa, 0x2a, 0x4a, 0xce,
	0x48, 0xcc, 0xcc, 0xd3, 0xaf, 0x4a, 0xcd, 0x2b, 0x29, 0xd0, 0x2f, 0x48, 0x2c, 0x4a, 0xcc, 0x2d,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x85, 0xca, 0xe9, 0x81, 0xe5, 0xa4, 0x04, 0x13,
	0x73, 0x33, 0xf3, 0xf2, 0xf5, 0xc1, 0x24, 0x44, 0x85, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x98,
	0xa9, 0x0f, 0x62, 0x41, 0x44, 0x95, 0x26, 0x32, 0x72, 0xb1, 0x05, 0x80, 0x0d, 0x12, 0x32, 0xe4,
	0x12, 0x2d, 0xce, 0xcf, 0x49, 0xcc, 0x4b, 0x8c, 0x2f, 0x4a, 0xcd, 0x49, 0xac, 0x4c, 0x2d, 0x8a,
	0xcf, 0x4e, 0xad, 0x8c, 0xcf, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x09, 0x12, 0x82, 0x48,
	0x06, 0x41, 0xe4, 0xbc, 0x53, 0x2b, 0x3d, 0x53, 0x84, 0x8c, 0xb9, 0xc4, 0xa0, 0xf6, 0xa2, 0xeb,
	0x61, 0x02, 0xeb, 0x11, 0x86, 0xca, 0x22, 0x6b, 0xb2, 0x92, 0x7f, 0xb1, 0x40, 0x9e, 0xb1, 0xeb,
	0xf9, 0x06, 0x2d, 0x98, 0x5e, 0xfd, 0x0a, 0xa8, 0x8f, 0x20, 0x0e, 0x71, 0x0a, 0x38, 0xf1, 0x48,
	0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0,
	0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xb3, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd,
	0xe4, 0xfc, 0x5c, 0xfd, 0xa8, 0xd4, 0xbc, 0xa2, 0xfc, 0xe4, 0x6c, 0x5d, 0xb7, 0xfc, 0xd2, 0xbc,
	0x94, 0xc4, 0x92, 0xcc, 0xfc, 0x3c, 0x7d, 0x98, 0x79, 0x65, 0xa6, 0x70, 0x23, 0x4b, 0x2a, 0x0b,
	0x52, 0x8b, 0x93, 0xd8, 0xc0, 0x9e, 0x35, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x41, 0xbc, 0x56,
	0xbc, 0x42, 0x01, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.SolanaRelayerKeyId != that1.SolanaRelayerKeyId {
		return false
	}
	if this.ZrchainRelayerKeyId != that1.ZrchainRelayerKeyId {
		return false
	}
	return true
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
	if m.ZrchainRelayerKeyId != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.ZrchainRelayerKeyId))
		i--
		dAtA[i] = 0x10
	}
	if m.SolanaRelayerKeyId != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SolanaRelayerKeyId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SolanaRelayerKeyId != 0 {
		n += 1 + sovParams(uint64(m.SolanaRelayerKeyId))
	}
	if m.ZrchainRelayerKeyId != 0 {
		n += 1 + sovParams(uint64(m.ZrchainRelayerKeyId))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SolanaRelayerKeyId", wireType)
			}
			m.SolanaRelayerKeyId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SolanaRelayerKeyId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZrchainRelayerKeyId", wireType)
			}
			m.ZrchainRelayerKeyId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ZrchainRelayerKeyId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
