// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: mail.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Mail_SendMail_FullMethodName = "/mail.Mail/SendMail"
)

// MailClient is the client API for Mail service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MailClient interface {
	SendMail(ctx context.Context, in *MailRequest, opts ...grpc.CallOption) (*MailReply, error)
}

type mailClient struct {
	cc grpc.ClientConnInterface
}

func NewMailClient(cc grpc.ClientConnInterface) MailClient {
	return &mailClient{cc}
}

func (c *mailClient) SendMail(ctx context.Context, in *MailRequest, opts ...grpc.CallOption) (*MailReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MailReply)
	err := c.cc.Invoke(ctx, Mail_SendMail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MailServer is the server API for Mail service.
// All implementations must embed UnimplementedMailServer
// for forward compatibility.
type MailServer interface {
	SendMail(context.Context, *MailRequest) (*MailReply, error)
	mustEmbedUnimplementedMailServer()
}

// UnimplementedMailServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMailServer struct{}

func (UnimplementedMailServer) SendMail(context.Context, *MailRequest) (*MailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMail not implemented")
}
func (UnimplementedMailServer) mustEmbedUnimplementedMailServer() {}
func (UnimplementedMailServer) testEmbeddedByValue()              {}

// UnsafeMailServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MailServer will
// result in compilation errors.
type UnsafeMailServer interface {
	mustEmbedUnimplementedMailServer()
}

func RegisterMailServer(s grpc.ServiceRegistrar, srv MailServer) {
	// If the following call pancis, it indicates UnimplementedMailServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Mail_ServiceDesc, srv)
}

func _Mail_SendMail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServer).SendMail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Mail_SendMail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServer).SendMail(ctx, req.(*MailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Mail_ServiceDesc is the grpc.ServiceDesc for Mail service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mail_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mail.Mail",
	HandlerType: (*MailServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMail",
			Handler:    _Mail_SendMail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mail.proto",
}
