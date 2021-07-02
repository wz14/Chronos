// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Core.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Message struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Sender               uint32   `protobuf:"varint,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver             uint32   `protobuf:"varint,3,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Data                 []byte   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c8d6e0ab9231f09, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message) GetSender() uint32 {
	if m != nil {
		return m.Sender
	}
	return 0
}

func (m *Message) GetReceiver() uint32 {
	if m != nil {
		return m.Receiver
	}
	return 0
}

func (m *Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Zero struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Zero) Reset()         { *m = Zero{} }
func (m *Zero) String() string { return proto.CompactTextString(m) }
func (*Zero) ProtoMessage()    {}
func (*Zero) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c8d6e0ab9231f09, []int{1}
}

func (m *Zero) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Zero.Unmarshal(m, b)
}
func (m *Zero) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Zero.Marshal(b, m, deterministic)
}
func (m *Zero) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Zero.Merge(m, src)
}
func (m *Zero) XXX_Size() int {
	return xxx_messageInfo_Zero.Size(m)
}
func (m *Zero) XXX_DiscardUnknown() {
	xxx_messageInfo_Zero.DiscardUnknown(m)
}

var xxx_messageInfo_Zero proto.InternalMessageInfo

type ECMsg struct {
	Index                uint64   `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Msglength            uint64   `protobuf:"varint,2,opt,name=msglength,proto3" json:"msglength,omitempty"`
	Merkleroot           []byte   `protobuf:"bytes,3,opt,name=merkleroot,proto3" json:"merkleroot,omitempty"`
	Merklepath           [][]byte `protobuf:"bytes,4,rep,name=merklepath,proto3" json:"merklepath,omitempty"`
	Merkleindex          []int64  `protobuf:"varint,5,rep,packed,name=merkleindex,proto3" json:"merkleindex,omitempty"`
	ErasureCode          []byte   `protobuf:"bytes,6,opt,name=erasureCode,proto3" json:"erasureCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ECMsg) Reset()         { *m = ECMsg{} }
func (m *ECMsg) String() string { return proto.CompactTextString(m) }
func (*ECMsg) ProtoMessage()    {}
func (*ECMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c8d6e0ab9231f09, []int{2}
}

func (m *ECMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ECMsg.Unmarshal(m, b)
}
func (m *ECMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ECMsg.Marshal(b, m, deterministic)
}
func (m *ECMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ECMsg.Merge(m, src)
}
func (m *ECMsg) XXX_Size() int {
	return xxx_messageInfo_ECMsg.Size(m)
}
func (m *ECMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ECMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ECMsg proto.InternalMessageInfo

func (m *ECMsg) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *ECMsg) GetMsglength() uint64 {
	if m != nil {
		return m.Msglength
	}
	return 0
}

func (m *ECMsg) GetMerkleroot() []byte {
	if m != nil {
		return m.Merkleroot
	}
	return nil
}

func (m *ECMsg) GetMerklepath() [][]byte {
	if m != nil {
		return m.Merklepath
	}
	return nil
}

func (m *ECMsg) GetMerkleindex() []int64 {
	if m != nil {
		return m.Merkleindex
	}
	return nil
}

func (m *ECMsg) GetErasureCode() []byte {
	if m != nil {
		return m.ErasureCode
	}
	return nil
}

type TS struct {
	Dummy                bool     `protobuf:"varint,1,opt,name=dummy,proto3" json:"dummy,omitempty"`
	Hash                 []byte   `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Sender               uint32   `protobuf:"varint,3,opt,name=sender,proto3" json:"sender,omitempty"`
	Num                  uint64   `protobuf:"varint,4,opt,name=num,proto3" json:"num,omitempty"`
	TS                   []uint64 `protobuf:"varint,5,rep,packed,name=TS,proto3" json:"TS,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TS) Reset()         { *m = TS{} }
func (m *TS) String() string { return proto.CompactTextString(m) }
func (*TS) ProtoMessage()    {}
func (*TS) Descriptor() ([]byte, []int) {
	return fileDescriptor_6c8d6e0ab9231f09, []int{3}
}

func (m *TS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TS.Unmarshal(m, b)
}
func (m *TS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TS.Marshal(b, m, deterministic)
}
func (m *TS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TS.Merge(m, src)
}
func (m *TS) XXX_Size() int {
	return xxx_messageInfo_TS.Size(m)
}
func (m *TS) XXX_DiscardUnknown() {
	xxx_messageInfo_TS.DiscardUnknown(m)
}

var xxx_messageInfo_TS proto.InternalMessageInfo

func (m *TS) GetDummy() bool {
	if m != nil {
		return m.Dummy
	}
	return false
}

func (m *TS) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *TS) GetSender() uint32 {
	if m != nil {
		return m.Sender
	}
	return 0
}

func (m *TS) GetNum() uint64 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *TS) GetTS() []uint64 {
	if m != nil {
		return m.TS
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "Message")
	proto.RegisterType((*Zero)(nil), "Zero")
	proto.RegisterType((*ECMsg)(nil), "ECMsg")
	proto.RegisterType((*TS)(nil), "TS")
}

func init() { proto.RegisterFile("Core.proto", fileDescriptor_6c8d6e0ab9231f09) }

var fileDescriptor_6c8d6e0ab9231f09 = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x51, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0xfd, 0x92, 0x4c, 0xd2, 0xf4, 0x7e, 0x55, 0x64, 0x10, 0x09, 0x45, 0x24, 0x64, 0x15, 0x10,
	0x2a, 0xe8, 0x1b, 0x18, 0x5c, 0xd6, 0xc5, 0xa4, 0xab, 0xee, 0xa6, 0x9d, 0x4b, 0x12, 0x6c, 0x32,
	0x65, 0x26, 0x15, 0x7d, 0x3a, 0x5f, 0x4d, 0xe6, 0xa6, 0x9a, 0xb8, 0x3b, 0x3f, 0x03, 0xe7, 0xcc,
	0xb9, 0x00, 0x85, 0x36, 0xb8, 0x3a, 0x1a, 0xdd, 0xeb, 0x4c, 0xc2, 0x6c, 0x8d, 0xd6, 0xca, 0x0a,
	0xf9, 0x25, 0xf8, 0x8d, 0x4a, 0xbc, 0xd4, 0xcb, 0xe7, 0xc2, 0x6f, 0x14, 0xbf, 0x81, 0xc8, 0x62,
	0xa7, 0xd0, 0x24, 0x7e, 0xea, 0xe5, 0x17, 0xe2, 0xcc, 0xf8, 0x12, 0x62, 0x83, 0x7b, 0x6c, 0xde,
	0xd1, 0x24, 0x01, 0x39, 0xbf, 0x9c, 0x73, 0x60, 0x4a, 0xf6, 0x32, 0x61, 0xa9, 0x97, 0x2f, 0x04,
	0xe1, 0x2c, 0x02, 0xb6, 0x45, 0xa3, 0xb3, 0x2f, 0x0f, 0xc2, 0x97, 0x62, 0x6d, 0x2b, 0x7e, 0x0d,
	0x61, 0xd3, 0x29, 0xfc, 0xa0, 0x30, 0x26, 0x06, 0xc2, 0x6f, 0x61, 0xde, 0xda, 0xea, 0x80, 0x5d,
	0xd5, 0xd7, 0x14, 0xc9, 0xc4, 0x28, 0xf0, 0x3b, 0x80, 0x16, 0xcd, 0xdb, 0x01, 0x8d, 0xd6, 0x3d,
	0xe5, 0x2e, 0xc4, 0x44, 0x19, 0xfd, 0xa3, 0xec, 0xeb, 0x84, 0xa5, 0xc1, 0xe8, 0x3b, 0x85, 0xa7,
	0xf0, 0x7f, 0x60, 0x43, 0x72, 0x98, 0x06, 0x79, 0x20, 0xa6, 0x92, 0x7b, 0x81, 0x46, 0xda, 0x93,
	0xc1, 0x42, 0x2b, 0x4c, 0x22, 0x8a, 0x98, 0x4a, 0x59, 0x0d, 0xfe, 0xa6, 0x74, 0xed, 0xd5, 0xa9,
	0x6d, 0x3f, 0xa9, 0x7d, 0x2c, 0x06, 0xe2, 0x7e, 0x5e, 0x4b, 0x3b, 0x14, 0x5f, 0x08, 0xc2, 0x93,
	0x05, 0x83, 0x3f, 0x0b, 0x5e, 0x41, 0xd0, 0x9d, 0x5a, 0x1a, 0x89, 0x09, 0x07, 0xdd, 0xf6, 0x9b,
	0x92, 0x4a, 0x31, 0xe1, 0x6f, 0xca, 0xc7, 0x7b, 0x98, 0xbd, 0x6a, 0x85, 0x85, 0xee, 0x5c, 0xad,
	0x12, 0x3b, 0xf5, 0x73, 0xa5, 0x78, 0x75, 0x46, 0xcb, 0x70, 0x45, 0xb3, 0xfe, 0x7b, 0x8e, 0xb7,
	0x91, 0x94, 0xfb, 0x87, 0xe3, 0x6e, 0x17, 0xd1, 0x51, 0x9f, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x15, 0xfb, 0x21, 0x8e, 0xe2, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NodeConClient is the client API for NodeCon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeConClient interface {
	SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Zero, error)
}

type nodeConClient struct {
	cc *grpc.ClientConn
}

func NewNodeConClient(cc *grpc.ClientConn) NodeConClient {
	return &nodeConClient{cc}
}

func (c *nodeConClient) SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Zero, error) {
	out := new(Zero)
	err := c.cc.Invoke(ctx, "/NodeCon/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeConServer is the server API for NodeCon service.
type NodeConServer interface {
	SendMessage(context.Context, *Message) (*Zero, error)
}

// UnimplementedNodeConServer can be embedded to have forward compatible implementations.
type UnimplementedNodeConServer struct {
}

func (*UnimplementedNodeConServer) SendMessage(ctx context.Context, req *Message) (*Zero, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func RegisterNodeConServer(s *grpc.Server, srv NodeConServer) {
	s.RegisterService(&_NodeCon_serviceDesc, srv)
}

func _NodeCon_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeConServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NodeCon/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeConServer).SendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _NodeCon_serviceDesc = grpc.ServiceDesc{
	ServiceName: "NodeCon",
	HandlerType: (*NodeConServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _NodeCon_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Core.proto",
}
