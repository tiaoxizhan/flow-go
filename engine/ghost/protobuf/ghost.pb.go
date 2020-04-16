// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ghost.proto

package ghost

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
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

type SubscribeRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeRequest) Reset()         { *m = SubscribeRequest{} }
func (m *SubscribeRequest) String() string { return proto.CompactTextString(m) }
func (*SubscribeRequest) ProtoMessage()    {}
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77edc9f77fb63d46, []int{0}
}

func (m *SubscribeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeRequest.Unmarshal(m, b)
}
func (m *SubscribeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeRequest.Marshal(b, m, deterministic)
}
func (m *SubscribeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeRequest.Merge(m, src)
}
func (m *SubscribeRequest) XXX_Size() int {
	return xxx_messageInfo_SubscribeRequest.Size(m)
}
func (m *SubscribeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeRequest proto.InternalMessageInfo

type SendEventRequest struct {
	ChannelId            uint32   `protobuf:"varint,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Message              []byte   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	TargetID             [][]byte `protobuf:"bytes,3,rep,name=targetID,proto3" json:"targetID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendEventRequest) Reset()         { *m = SendEventRequest{} }
func (m *SendEventRequest) String() string { return proto.CompactTextString(m) }
func (*SendEventRequest) ProtoMessage()    {}
func (*SendEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77edc9f77fb63d46, []int{1}
}

func (m *SendEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendEventRequest.Unmarshal(m, b)
}
func (m *SendEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendEventRequest.Marshal(b, m, deterministic)
}
func (m *SendEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendEventRequest.Merge(m, src)
}
func (m *SendEventRequest) XXX_Size() int {
	return xxx_messageInfo_SendEventRequest.Size(m)
}
func (m *SendEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendEventRequest proto.InternalMessageInfo

func (m *SendEventRequest) GetChannelId() uint32 {
	if m != nil {
		return m.ChannelId
	}
	return 0
}

func (m *SendEventRequest) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *SendEventRequest) GetTargetID() [][]byte {
	if m != nil {
		return m.TargetID
	}
	return nil
}

type FlowMessage struct {
	SenderID             []byte   `protobuf:"bytes,1,opt,name=senderID,proto3" json:"senderID,omitempty"`
	Message              []byte   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FlowMessage) Reset()         { *m = FlowMessage{} }
func (m *FlowMessage) String() string { return proto.CompactTextString(m) }
func (*FlowMessage) ProtoMessage()    {}
func (*FlowMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_77edc9f77fb63d46, []int{2}
}

func (m *FlowMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FlowMessage.Unmarshal(m, b)
}
func (m *FlowMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FlowMessage.Marshal(b, m, deterministic)
}
func (m *FlowMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FlowMessage.Merge(m, src)
}
func (m *FlowMessage) XXX_Size() int {
	return xxx_messageInfo_FlowMessage.Size(m)
}
func (m *FlowMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_FlowMessage.DiscardUnknown(m)
}

var xxx_messageInfo_FlowMessage proto.InternalMessageInfo

func (m *FlowMessage) GetSenderID() []byte {
	if m != nil {
		return m.SenderID
	}
	return nil
}

func (m *FlowMessage) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func init() {
	proto.RegisterType((*SubscribeRequest)(nil), "ghost.SubscribeRequest")
	proto.RegisterType((*SendEventRequest)(nil), "ghost.SendEventRequest")
	proto.RegisterType((*FlowMessage)(nil), "ghost.FlowMessage")
}

func init() { proto.RegisterFile("ghost.proto", fileDescriptor_77edc9f77fb63d46) }

var fileDescriptor_77edc9f77fb63d46 = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xdd, 0x4a, 0xc3, 0x30,
	0x14, 0xc7, 0x89, 0xc3, 0x8f, 0x9e, 0x55, 0x18, 0xe7, 0x42, 0x4b, 0x45, 0x28, 0xbd, 0xea, 0x55,
	0x26, 0x7a, 0x27, 0xde, 0x88, 0x9b, 0xd2, 0x0b, 0x45, 0xb2, 0x07, 0x90, 0x76, 0x39, 0x66, 0x83,
	0x2e, 0x99, 0x4d, 0xaa, 0xf8, 0x06, 0x3e, 0xb6, 0xac, 0xed, 0x32, 0x19, 0x78, 0xf9, 0xfb, 0xe7,
	0xe4, 0x7c, 0xfc, 0x60, 0xa8, 0x16, 0xc6, 0x3a, 0xbe, 0xae, 0x8d, 0x33, 0x78, 0xd8, 0x42, 0x7c,
	0xa1, 0x8c, 0x51, 0x15, 0x8d, 0xdb, 0xb0, 0x6c, 0xde, 0xc7, 0xb4, 0x5a, 0xbb, 0xef, 0xae, 0x26,
	0x45, 0x18, 0xcd, 0x9a, 0xd2, 0xce, 0xeb, 0x65, 0x49, 0x82, 0x3e, 0x1a, 0xb2, 0x2e, 0x55, 0x30,
	0x9a, 0x91, 0x96, 0xd3, 0x4f, 0xd2, 0xae, 0xcf, 0xf0, 0x12, 0x60, 0xbe, 0x28, 0xb4, 0xa6, 0xea,
	0x6d, 0x29, 0x23, 0x96, 0xb0, 0xec, 0x54, 0x04, 0x7d, 0x92, 0x4b, 0x8c, 0xe0, 0x78, 0x45, 0xd6,
	0x16, 0x8a, 0xa2, 0x83, 0x84, 0x65, 0xa1, 0xd8, 0x22, 0xc6, 0x70, 0xe2, 0x8a, 0x5a, 0x91, 0xcb,
	0x27, 0xd1, 0x20, 0x19, 0x64, 0xa1, 0xf0, 0x9c, 0x3e, 0xc0, 0xf0, 0xb1, 0x32, 0x5f, 0xcf, 0xbb,
	0x52, 0x4b, 0x5a, 0x52, 0x9d, 0x4f, 0xda, 0x09, 0xa1, 0xf0, 0xfc, 0xff, 0x80, 0xeb, 0x1f, 0x06,
	0xe1, 0xd3, 0xe6, 0xd0, 0x17, 0x23, 0xe9, 0xfe, 0x35, 0xc7, 0x3b, 0x08, 0xfc, 0xfa, 0x78, 0xce,
	0x3b, 0x23, 0xfb, 0x07, 0xc5, 0x67, 0xbc, 0xd3, 0xc2, 0xb7, 0x5a, 0xf8, 0x74, 0xa3, 0x05, 0x6f,
	0x21, 0xf0, 0x42, 0x76, 0xbf, 0xf7, 0x14, 0xc5, 0xd8, 0x3f, 0xfc, 0x59, 0xff, 0x8a, 0x95, 0x47,
	0x6d, 0xaf, 0x9b, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9b, 0xda, 0x04, 0x8d, 0x86, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GhostNodeAPIClient is the client API for GhostNodeAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GhostNodeAPIClient interface {
	// SendEvent submits and event to the internal Flow Libp2p network
	SendEvent(ctx context.Context, in *SendEventRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Subscribe returns all network messages
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (GhostNodeAPI_SubscribeClient, error)
}

type ghostNodeAPIClient struct {
	cc *grpc.ClientConn
}

func NewGhostNodeAPIClient(cc *grpc.ClientConn) GhostNodeAPIClient {
	return &ghostNodeAPIClient{cc}
}

func (c *ghostNodeAPIClient) SendEvent(ctx context.Context, in *SendEventRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/ghost.GhostNodeAPI/SendEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ghostNodeAPIClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (GhostNodeAPI_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GhostNodeAPI_serviceDesc.Streams[0], "/ghost.GhostNodeAPI/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &ghostNodeAPISubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GhostNodeAPI_SubscribeClient interface {
	Recv() (*FlowMessage, error)
	grpc.ClientStream
}

type ghostNodeAPISubscribeClient struct {
	grpc.ClientStream
}

func (x *ghostNodeAPISubscribeClient) Recv() (*FlowMessage, error) {
	m := new(FlowMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GhostNodeAPIServer is the server API for GhostNodeAPI service.
type GhostNodeAPIServer interface {
	// SendEvent submits and event to the internal Flow Libp2p network
	SendEvent(context.Context, *SendEventRequest) (*empty.Empty, error)
	// Subscribe returns all network messages
	Subscribe(*SubscribeRequest, GhostNodeAPI_SubscribeServer) error
}

// UnimplementedGhostNodeAPIServer can be embedded to have forward compatible implementations.
type UnimplementedGhostNodeAPIServer struct {
}

func (*UnimplementedGhostNodeAPIServer) SendEvent(ctx context.Context, req *SendEventRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEvent not implemented")
}
func (*UnimplementedGhostNodeAPIServer) Subscribe(req *SubscribeRequest, srv GhostNodeAPI_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}

func RegisterGhostNodeAPIServer(s *grpc.Server, srv GhostNodeAPIServer) {
	s.RegisterService(&_GhostNodeAPI_serviceDesc, srv)
}

func _GhostNodeAPI_SendEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GhostNodeAPIServer).SendEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ghost.GhostNodeAPI/SendEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GhostNodeAPIServer).SendEvent(ctx, req.(*SendEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GhostNodeAPI_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GhostNodeAPIServer).Subscribe(m, &ghostNodeAPISubscribeServer{stream})
}

type GhostNodeAPI_SubscribeServer interface {
	Send(*FlowMessage) error
	grpc.ServerStream
}

type ghostNodeAPISubscribeServer struct {
	grpc.ServerStream
}

func (x *ghostNodeAPISubscribeServer) Send(m *FlowMessage) error {
	return x.ServerStream.SendMsg(m)
}

var _GhostNodeAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ghost.GhostNodeAPI",
	HandlerType: (*GhostNodeAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEvent",
			Handler:    _GhostNodeAPI_SendEvent_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _GhostNodeAPI_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ghost.proto",
}
