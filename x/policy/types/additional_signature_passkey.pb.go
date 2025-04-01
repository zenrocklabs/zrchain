// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zrchain/policy/additional_signature_passkey.proto

package types

import (
	fmt "fmt"
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

// AdditionalSignaturePasskey is a message that contains passkey signature data
type AdditionalSignaturePasskey struct {
	RawId             []byte `protobuf:"bytes,1,opt,name=raw_id,json=rawId,proto3" json:"raw_id,omitempty"`
	AuthenticatorData []byte `protobuf:"bytes,2,opt,name=authenticator_data,json=authenticatorData,proto3" json:"authenticator_data,omitempty"`
	ClientDataJson    []byte `protobuf:"bytes,3,opt,name=client_data_json,json=clientDataJson,proto3" json:"client_data_json,omitempty"`
	Signature         []byte `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *AdditionalSignaturePasskey) Reset()         { *m = AdditionalSignaturePasskey{} }
func (m *AdditionalSignaturePasskey) String() string { return proto.CompactTextString(m) }
func (*AdditionalSignaturePasskey) ProtoMessage()    {}
func (*AdditionalSignaturePasskey) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed2e7f8c91479d0c, []int{0}
}
func (m *AdditionalSignaturePasskey) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AdditionalSignaturePasskey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AdditionalSignaturePasskey.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AdditionalSignaturePasskey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdditionalSignaturePasskey.Merge(m, src)
}
func (m *AdditionalSignaturePasskey) XXX_Size() int {
	return m.Size()
}
func (m *AdditionalSignaturePasskey) XXX_DiscardUnknown() {
	xxx_messageInfo_AdditionalSignaturePasskey.DiscardUnknown(m)
}

var xxx_messageInfo_AdditionalSignaturePasskey proto.InternalMessageInfo

func (m *AdditionalSignaturePasskey) GetRawId() []byte {
	if m != nil {
		return m.RawId
	}
	return nil
}

func (m *AdditionalSignaturePasskey) GetAuthenticatorData() []byte {
	if m != nil {
		return m.AuthenticatorData
	}
	return nil
}

func (m *AdditionalSignaturePasskey) GetClientDataJson() []byte {
	if m != nil {
		return m.ClientDataJson
	}
	return nil
}

func (m *AdditionalSignaturePasskey) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*AdditionalSignaturePasskey)(nil), "zrchain.policy.AdditionalSignaturePasskey")
}

func init() {
	proto.RegisterFile("zrchain/policy/additional_signature_passkey.proto", fileDescriptor_ed2e7f8c91479d0c)
}

var fileDescriptor_ed2e7f8c91479d0c = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xcf, 0x4a, 0xf4, 0x30,
	0x14, 0x47, 0x9b, 0xef, 0xd3, 0x01, 0x83, 0x0c, 0x1a, 0x10, 0x8a, 0x48, 0x10, 0x57, 0xb3, 0x99,
	0x16, 0x11, 0x74, 0xad, 0x88, 0xa0, 0x2b, 0xff, 0xec, 0x66, 0x53, 0xee, 0x24, 0x61, 0x1a, 0xa7,
	0x26, 0x25, 0xb9, 0x75, 0xac, 0x4f, 0xe1, 0x73, 0xf8, 0x24, 0x2e, 0x67, 0xe9, 0x52, 0xda, 0x17,
	0x91, 0x49, 0xed, 0x88, 0xdb, 0x5f, 0xce, 0x81, 0xdc, 0x43, 0x8f, 0x5f, 0x9d, 0xc8, 0x41, 0x9b,
	0xb4, 0xb4, 0x85, 0x16, 0x75, 0x0a, 0x52, 0x6a, 0xd4, 0xd6, 0x40, 0x91, 0x79, 0x3d, 0x33, 0x80,
	0x95, 0x53, 0x59, 0x09, 0xde, 0xcf, 0x55, 0x9d, 0x94, 0xce, 0xa2, 0x65, 0xc3, 0x1f, 0x25, 0xe9,
	0x94, 0xa3, 0x77, 0x42, 0xf7, 0xcf, 0xd7, 0xda, 0x43, 0x6f, 0xdd, 0x76, 0x12, 0xdb, 0xa3, 0x03,
	0x07, 0x8b, 0x4c, 0xcb, 0x98, 0x1c, 0x92, 0xd1, 0xf6, 0xfd, 0xa6, 0x83, 0xc5, 0xb5, 0x64, 0x63,
	0xca, 0xa0, 0xc2, 0x5c, 0x19, 0xd4, 0x02, 0xd0, 0xba, 0x4c, 0x02, 0x42, 0xfc, 0x2f, 0x20, 0xbb,
	0x7f, 0x5e, 0x2e, 0x01, 0x81, 0x8d, 0xe8, 0x8e, 0x28, 0xb4, 0x32, 0x18, 0xb8, 0xec, 0xd1, 0x5b,
	0x13, 0xff, 0x0f, 0xf0, 0xb0, 0xdb, 0x57, 0xd4, 0x8d, 0xb7, 0x86, 0x1d, 0xd0, 0xad, 0xf5, 0xcf,
	0xe3, 0x8d, 0x80, 0xfc, 0x0e, 0x17, 0x77, 0x1f, 0x0d, 0x27, 0xcb, 0x86, 0x93, 0xaf, 0x86, 0x93,
	0xb7, 0x96, 0x47, 0xcb, 0x96, 0x47, 0x9f, 0x2d, 0x8f, 0x26, 0x67, 0x33, 0x8d, 0x79, 0x35, 0x4d,
	0x84, 0x7d, 0x4a, 0x27, 0xca, 0x38, 0x2b, 0xe6, 0xe3, 0x2b, 0x5b, 0x19, 0x09, 0xab, 0xc3, 0xd2,
	0xbe, 0xd3, 0xf3, 0x69, 0xfa, 0xd2, 0xc7, 0xc2, 0xba, 0x54, 0x7e, 0x3a, 0x08, 0x59, 0x4e, 0xbe,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x82, 0x17, 0x17, 0xe1, 0x4b, 0x01, 0x00, 0x00,
}

func (m *AdditionalSignaturePasskey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AdditionalSignaturePasskey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AdditionalSignaturePasskey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintAdditionalSignaturePasskey(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ClientDataJson) > 0 {
		i -= len(m.ClientDataJson)
		copy(dAtA[i:], m.ClientDataJson)
		i = encodeVarintAdditionalSignaturePasskey(dAtA, i, uint64(len(m.ClientDataJson)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.AuthenticatorData) > 0 {
		i -= len(m.AuthenticatorData)
		copy(dAtA[i:], m.AuthenticatorData)
		i = encodeVarintAdditionalSignaturePasskey(dAtA, i, uint64(len(m.AuthenticatorData)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.RawId) > 0 {
		i -= len(m.RawId)
		copy(dAtA[i:], m.RawId)
		i = encodeVarintAdditionalSignaturePasskey(dAtA, i, uint64(len(m.RawId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAdditionalSignaturePasskey(dAtA []byte, offset int, v uint64) int {
	offset -= sovAdditionalSignaturePasskey(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AdditionalSignaturePasskey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RawId)
	if l > 0 {
		n += 1 + l + sovAdditionalSignaturePasskey(uint64(l))
	}
	l = len(m.AuthenticatorData)
	if l > 0 {
		n += 1 + l + sovAdditionalSignaturePasskey(uint64(l))
	}
	l = len(m.ClientDataJson)
	if l > 0 {
		n += 1 + l + sovAdditionalSignaturePasskey(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovAdditionalSignaturePasskey(uint64(l))
	}
	return n
}

func sovAdditionalSignaturePasskey(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAdditionalSignaturePasskey(x uint64) (n int) {
	return sovAdditionalSignaturePasskey(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AdditionalSignaturePasskey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAdditionalSignaturePasskey
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
			return fmt.Errorf("proto: AdditionalSignaturePasskey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AdditionalSignaturePasskey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RawId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdditionalSignaturePasskey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RawId = append(m.RawId[:0], dAtA[iNdEx:postIndex]...)
			if m.RawId == nil {
				m.RawId = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuthenticatorData", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdditionalSignaturePasskey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AuthenticatorData = append(m.AuthenticatorData[:0], dAtA[iNdEx:postIndex]...)
			if m.AuthenticatorData == nil {
				m.AuthenticatorData = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientDataJson", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdditionalSignaturePasskey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientDataJson = append(m.ClientDataJson[:0], dAtA[iNdEx:postIndex]...)
			if m.ClientDataJson == nil {
				m.ClientDataJson = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAdditionalSignaturePasskey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAdditionalSignaturePasskey(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAdditionalSignaturePasskey
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
func skipAdditionalSignaturePasskey(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAdditionalSignaturePasskey
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
					return 0, ErrIntOverflowAdditionalSignaturePasskey
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
					return 0, ErrIntOverflowAdditionalSignaturePasskey
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
				return 0, ErrInvalidLengthAdditionalSignaturePasskey
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAdditionalSignaturePasskey
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAdditionalSignaturePasskey
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAdditionalSignaturePasskey        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAdditionalSignaturePasskey          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAdditionalSignaturePasskey = fmt.Errorf("proto: unexpected end of group")
)
