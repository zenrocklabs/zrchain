// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zrchain/zentp/bridge.proto

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

// MBStatus represents the different possible states of a mint/burn operation.
type BridgeStatus int32

const (
	// Undefined: The status of the operation is not specified.
	BridgeStatus_BRIDGE_STATUS_UNSPECIFIED BridgeStatus = 0
	// Pending: The operation is currently being processed.
	BridgeStatus_BRIDGE_STATUS_PENDING BridgeStatus = 1
	// Completed: The operation has been successfully finalized.
	BridgeStatus_BRIDGE_STATUS_COMPLETED BridgeStatus = 2
	// Failed: The operation has failed.
	BridgeStatus_BRIDGE_STATUS_FAILED BridgeStatus = 4
)

var BridgeStatus_name = map[int32]string{
	0: "BRIDGE_STATUS_UNSPECIFIED",
	1: "BRIDGE_STATUS_PENDING",
	2: "BRIDGE_STATUS_COMPLETED",
	4: "BRIDGE_STATUS_FAILED",
}

var BridgeStatus_value = map[string]int32{
	"BRIDGE_STATUS_UNSPECIFIED": 0,
	"BRIDGE_STATUS_PENDING":     1,
	"BRIDGE_STATUS_COMPLETED":   2,
	"BRIDGE_STATUS_FAILED":      4,
}

func (x BridgeStatus) String() string {
	return proto.EnumName(BridgeStatus_name, int32(x))
}

func (BridgeStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e14233fdb9bf0972, []int{0}
}

// Bridge represents a mint and burn operation between two networks.
type Bridge struct {
	Id               uint64       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Denom            string       `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	Creator          string       `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	SourceAddress    string       `protobuf:"bytes,4,opt,name=source_address,json=sourceAddress,proto3" json:"source_address,omitempty"`
	SourceChain      string       `protobuf:"bytes,5,opt,name=source_chain,json=sourceChain,proto3" json:"source_chain,omitempty"`
	DestinationChain string       `protobuf:"bytes,6,opt,name=destination_chain,json=destinationChain,proto3" json:"destination_chain,omitempty"`
	Amount           uint64       `protobuf:"varint,7,opt,name=amount,proto3" json:"amount,omitempty"`
	RecipientAddress string       `protobuf:"bytes,8,opt,name=recipient_address,json=recipientAddress,proto3" json:"recipient_address,omitempty"`
	TxId             uint64       `protobuf:"varint,9,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	TxHash           string       `protobuf:"bytes,10,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	State            BridgeStatus `protobuf:"varint,11,opt,name=state,proto3,enum=zrchain.zentp.BridgeStatus" json:"state,omitempty"`
	BlockHeight      int64        `protobuf:"varint,12,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

func (m *Bridge) Reset()         { *m = Bridge{} }
func (m *Bridge) String() string { return proto.CompactTextString(m) }
func (*Bridge) ProtoMessage()    {}
func (*Bridge) Descriptor() ([]byte, []int) {
	return fileDescriptor_e14233fdb9bf0972, []int{0}
}
func (m *Bridge) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Bridge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Bridge.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Bridge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bridge.Merge(m, src)
}
func (m *Bridge) XXX_Size() int {
	return m.Size()
}
func (m *Bridge) XXX_DiscardUnknown() {
	xxx_messageInfo_Bridge.DiscardUnknown(m)
}

var xxx_messageInfo_Bridge proto.InternalMessageInfo

func (m *Bridge) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Bridge) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *Bridge) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *Bridge) GetSourceAddress() string {
	if m != nil {
		return m.SourceAddress
	}
	return ""
}

func (m *Bridge) GetSourceChain() string {
	if m != nil {
		return m.SourceChain
	}
	return ""
}

func (m *Bridge) GetDestinationChain() string {
	if m != nil {
		return m.DestinationChain
	}
	return ""
}

func (m *Bridge) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *Bridge) GetRecipientAddress() string {
	if m != nil {
		return m.RecipientAddress
	}
	return ""
}

func (m *Bridge) GetTxId() uint64 {
	if m != nil {
		return m.TxId
	}
	return 0
}

func (m *Bridge) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

func (m *Bridge) GetState() BridgeStatus {
	if m != nil {
		return m.State
	}
	return BridgeStatus_BRIDGE_STATUS_UNSPECIFIED
}

func (m *Bridge) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func init() {
	proto.RegisterEnum("zrchain.zentp.BridgeStatus", BridgeStatus_name, BridgeStatus_value)
	proto.RegisterType((*Bridge)(nil), "zrchain.zentp.Bridge")
}

func init() { proto.RegisterFile("zrchain/zentp/bridge.proto", fileDescriptor_e14233fdb9bf0972) }

var fileDescriptor_e14233fdb9bf0972 = []byte{
	// 488 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xcf, 0x6f, 0xda, 0x3e,
	0x18, 0xc6, 0x09, 0x3f, 0xbf, 0x35, 0x14, 0x51, 0x7f, 0xd9, 0xea, 0x52, 0x2d, 0x62, 0x93, 0x26,
	0xa1, 0x4d, 0x23, 0xda, 0x26, 0xf5, 0xce, 0x8f, 0xd0, 0x46, 0xea, 0x18, 0x02, 0x7a, 0xe9, 0x25,
	0x32, 0xb1, 0x45, 0xac, 0x0e, 0x1b, 0xd9, 0xce, 0x94, 0xf5, 0xb2, 0xeb, 0x8e, 0xfb, 0xb3, 0x76,
	0xec, 0x71, 0xc7, 0x09, 0xfe, 0x91, 0x09, 0x3b, 0x54, 0x65, 0x97, 0xc8, 0xef, 0xf3, 0x7c, 0xfc,
	0xe6, 0x49, 0xde, 0x17, 0xb4, 0xee, 0x65, 0x14, 0x63, 0xc6, 0xbd, 0x7b, 0xca, 0xf5, 0xda, 0x5b,
	0x48, 0x46, 0x96, 0xb4, 0xbb, 0x96, 0x42, 0x0b, 0x78, 0x9c, 0x79, 0x5d, 0xe3, 0xb5, 0x4e, 0xf0,
	0x8a, 0x71, 0xe1, 0x99, 0xa7, 0x25, 0x5a, 0xcd, 0xa5, 0x58, 0x0a, 0x73, 0xf4, 0x76, 0xa7, 0x4c,
	0xfd, 0xa7, 0xe7, 0x1a, 0x4b, 0xbc, 0x52, 0xd6, 0x7b, 0xf5, 0xa3, 0x00, 0xca, 0x7d, 0xf3, 0x12,
	0x58, 0x07, 0x79, 0x46, 0x90, 0xd3, 0x76, 0x3a, 0xc5, 0x69, 0x9e, 0x11, 0xd8, 0x04, 0x25, 0x42,
	0xb9, 0x58, 0xa1, 0x7c, 0xdb, 0xe9, 0x1c, 0x4d, 0x6d, 0x01, 0x11, 0xa8, 0x44, 0x92, 0x62, 0x2d,
	0x24, 0x2a, 0x18, 0x7d, 0x5f, 0xc2, 0xd7, 0xa0, 0xae, 0x44, 0x22, 0x23, 0x1a, 0x62, 0x42, 0x24,
	0x55, 0x0a, 0x15, 0x0d, 0x70, 0x6c, 0xd5, 0x9e, 0x15, 0xe1, 0x4b, 0x50, 0xcb, 0x30, 0x13, 0x0a,
	0x95, 0x0c, 0x54, 0xb5, 0xda, 0x60, 0x27, 0xc1, 0xb7, 0xe0, 0x84, 0x50, 0xa5, 0x19, 0xc7, 0x9a,
	0x09, 0x9e, 0x71, 0x65, 0xc3, 0x35, 0x9e, 0x18, 0x16, 0x7e, 0x0e, 0xca, 0x78, 0x25, 0x12, 0xae,
	0x51, 0xc5, 0x44, 0xcf, 0xaa, 0x5d, 0x13, 0x49, 0x23, 0xb6, 0x66, 0x94, 0xeb, 0xc7, 0x44, 0xff,
	0xd9, 0x26, 0x8f, 0xc6, 0x3e, 0xd4, 0xff, 0xa0, 0xa4, 0xd3, 0x90, 0x11, 0x74, 0x64, 0x7a, 0x14,
	0x75, 0x1a, 0x10, 0x78, 0x0a, 0x2a, 0x3a, 0x0d, 0x63, 0xac, 0x62, 0x04, 0xcc, 0xbd, 0xb2, 0x4e,
	0xaf, 0xb0, 0x8a, 0xe1, 0x7b, 0x50, 0x52, 0x1a, 0x6b, 0x8a, 0xaa, 0x6d, 0xa7, 0x53, 0xff, 0x70,
	0xde, 0x3d, 0x18, 0x4c, 0xd7, 0xfe, 0xcf, 0x99, 0xc6, 0x3a, 0x51, 0x53, 0x4b, 0xee, 0xbe, 0x7a,
	0xf1, 0x45, 0x44, 0x77, 0x61, 0x4c, 0xd9, 0x32, 0xd6, 0xa8, 0xd6, 0x76, 0x3a, 0x85, 0x69, 0xd5,
	0x68, 0x57, 0x46, 0x7a, 0xf3, 0x1d, 0xd4, 0x9e, 0xde, 0x84, 0x2f, 0xc0, 0x59, 0x7f, 0x1a, 0x0c,
	0x2f, 0xfd, 0x70, 0x36, 0xef, 0xcd, 0x6f, 0x66, 0xe1, 0xcd, 0x78, 0x36, 0xf1, 0x07, 0xc1, 0x28,
	0xf0, 0x87, 0x8d, 0x1c, 0x3c, 0x03, 0xcf, 0x0e, 0xed, 0x89, 0x3f, 0x1e, 0x06, 0xe3, 0xcb, 0x86,
	0x03, 0xcf, 0xc1, 0xe9, 0xa1, 0x35, 0xf8, 0xfc, 0x69, 0x72, 0xed, 0xcf, 0xfd, 0x61, 0x23, 0x0f,
	0x11, 0x68, 0x1e, 0x9a, 0xa3, 0x5e, 0x70, 0xed, 0x0f, 0x1b, 0xc5, 0xfe, 0xe4, 0xd7, 0xc6, 0x75,
	0x1e, 0x36, 0xae, 0xf3, 0x67, 0xe3, 0x3a, 0x3f, 0xb7, 0x6e, 0xee, 0x61, 0xeb, 0xe6, 0x7e, 0x6f,
	0xdd, 0xdc, 0xed, 0xc5, 0x92, 0xe9, 0x38, 0x59, 0x74, 0x23, 0xb1, 0xf2, 0x6e, 0x29, 0x97, 0x22,
	0xba, 0x7b, 0x37, 0x12, 0x09, 0x27, 0x66, 0x0e, 0xde, 0x7e, 0xbf, 0xbe, 0x5e, 0x78, 0x69, 0xb6,
	0x64, 0xfa, 0xdb, 0x9a, 0xaa, 0x45, 0xd9, 0x2c, 0xd9, 0xc7, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x83, 0xc4, 0x3e, 0x76, 0xd6, 0x02, 0x00, 0x00,
}

func (m *Bridge) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Bridge) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Bridge) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BlockHeight != 0 {
		i = encodeVarintBridge(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x60
	}
	if m.State != 0 {
		i = encodeVarintBridge(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x58
	}
	if len(m.TxHash) > 0 {
		i -= len(m.TxHash)
		copy(dAtA[i:], m.TxHash)
		i = encodeVarintBridge(dAtA, i, uint64(len(m.TxHash)))
		i--
		dAtA[i] = 0x52
	}
	if m.TxId != 0 {
		i = encodeVarintBridge(dAtA, i, uint64(m.TxId))
		i--
		dAtA[i] = 0x48
	}
	if len(m.RecipientAddress) > 0 {
		i -= len(m.RecipientAddress)
		copy(dAtA[i:], m.RecipientAddress)
		i = encodeVarintBridge(dAtA, i, uint64(len(m.RecipientAddress)))
		i--
		dAtA[i] = 0x42
	}
	if m.Amount != 0 {
		i = encodeVarintBridge(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x38
	}
	if len(m.DestinationChain) > 0 {
		i -= len(m.DestinationChain)
		copy(dAtA[i:], m.DestinationChain)
		i = encodeVarintBridge(dAtA, i, uint64(len(m.DestinationChain)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.SourceChain) > 0 {
		i -= len(m.SourceChain)
		copy(dAtA[i:], m.SourceChain)
		i = encodeVarintBridge(dAtA, i, uint64(len(m.SourceChain)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.SourceAddress) > 0 {
		i -= len(m.SourceAddress)
		copy(dAtA[i:], m.SourceAddress)
		i = encodeVarintBridge(dAtA, i, uint64(len(m.SourceAddress)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintBridge(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintBridge(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintBridge(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintBridge(dAtA []byte, offset int, v uint64) int {
	offset -= sovBridge(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Bridge) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovBridge(uint64(m.Id))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	l = len(m.SourceAddress)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	l = len(m.SourceChain)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	l = len(m.DestinationChain)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovBridge(uint64(m.Amount))
	}
	l = len(m.RecipientAddress)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	if m.TxId != 0 {
		n += 1 + sovBridge(uint64(m.TxId))
	}
	l = len(m.TxHash)
	if l > 0 {
		n += 1 + l + sovBridge(uint64(l))
	}
	if m.State != 0 {
		n += 1 + sovBridge(uint64(m.State))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovBridge(uint64(m.BlockHeight))
	}
	return n
}

func sovBridge(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBridge(x uint64) (n int) {
	return sovBridge(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Bridge) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
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
			return fmt.Errorf("proto: Bridge: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Bridge: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RecipientAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RecipientAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxId", wireType)
			}
			m.TxId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TxId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
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
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= BridgeStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBridge
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
func skipBridge(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBridge
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
					return 0, ErrIntOverflowBridge
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
					return 0, ErrIntOverflowBridge
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
				return 0, ErrInvalidLengthBridge
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBridge
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBridge
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBridge        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBridge          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBridge = fmt.Errorf("proto: unexpected end of group")
)
