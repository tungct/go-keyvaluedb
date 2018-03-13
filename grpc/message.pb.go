// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

/*
Package grpc is a generated protocol buffer package.

It is generated from these files:
	message.proto

It has these top-level messages:
	Message
	MessageResponse
*/
package grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc1 "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Id      int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content" json:"content,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Message) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type MessageResponse struct {
	Content string `protobuf:"bytes,1,opt,name=content" json:"content,omitempty"`
}

func (m *MessageResponse) Reset()                    { *m = MessageResponse{} }
func (m *MessageResponse) String() string            { return proto.CompactTextString(m) }
func (*MessageResponse) ProtoMessage()               {}
func (*MessageResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MessageResponse) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "grpc.Message")
	proto.RegisterType((*MessageResponse)(nil), "grpc.MessageResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc1.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc1.SupportPackageIsVersion4

// Client API for SendMessage service

type SendMessageClient interface {
	SendMessage(ctx context.Context, in *Message, opts ...grpc1.CallOption) (*MessageResponse, error)
}

type sendMessageClient struct {
	cc *grpc1.ClientConn
}

func NewSendMessageClient(cc *grpc1.ClientConn) SendMessageClient {
	return &sendMessageClient{cc}
}

func (c *sendMessageClient) SendMessage(ctx context.Context, in *Message, opts ...grpc1.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := grpc1.Invoke(ctx, "/grpc.SendMessage/SendMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SendMessage service

type SendMessageServer interface {
	SendMessage(context.Context, *Message) (*MessageResponse, error)
}

func RegisterSendMessageServer(s *grpc1.Server, srv SendMessageServer) {
	s.RegisterService(&_SendMessage_serviceDesc, srv)
}

func _SendMessage_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc1.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendMessageServer).SendMessage(ctx, in)
	}
	info := &grpc1.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.SendMessage/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendMessageServer).SendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _SendMessage_serviceDesc = grpc1.ServiceDesc{
	ServiceName: "grpc.SendMessage",
	HandlerType: (*SendMessageServer)(nil),
	Methods: []grpc1.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _SendMessage_SendMessage_Handler,
		},
	},
	Streams:  []grpc1.StreamDesc{},
	Metadata: "message.proto",
}

func init() { proto.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2f, 0x2a, 0x48, 0x56,
	0x32, 0xe6, 0x62, 0xf7, 0x85, 0x08, 0x0b, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30,
	0x6a, 0xb0, 0x06, 0x31, 0x65, 0xa6, 0x08, 0x49, 0x70, 0xb1, 0x27, 0xe7, 0xe7, 0x95, 0xa4, 0xe6,
	0x95, 0x48, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x4a, 0xda, 0x5c, 0xfc, 0x50, 0x4d,
	0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0xc8, 0x8a, 0x19, 0x51, 0x14, 0x1b, 0xb9, 0x70,
	0x71, 0x07, 0xa7, 0xe6, 0xa5, 0xc0, 0x6c, 0x31, 0x45, 0xe5, 0xf2, 0xea, 0x81, 0x9c, 0xa1, 0x07,
	0xe5, 0x4a, 0x89, 0xa2, 0x70, 0x61, 0xa6, 0x2b, 0x31, 0x24, 0xb1, 0x81, 0x1d, 0x6d, 0x0c, 0x08,
	0x00, 0x00, 0xff, 0xff, 0x2a, 0x6a, 0xd8, 0x16, 0xc5, 0x00, 0x00, 0x00,
}