// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/email-sender.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	EmailSenderApi_SendEmail_FullMethodName = "/proto.EmailSenderApi/SendEmail"
)

// EmailSenderApiClient is the client API for EmailSenderApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmailSenderApiClient interface {
	// SendEmail - RPC method for sending email.
	SendEmail(ctx context.Context, in *EmailSenderRequest, opts ...grpc.CallOption) (*EmailSenderResponse, error)
}

type emailSenderApiClient struct {
	cc grpc.ClientConnInterface
}

func NewEmailSenderApiClient(cc grpc.ClientConnInterface) EmailSenderApiClient {
	return &emailSenderApiClient{cc}
}

func (c *emailSenderApiClient) SendEmail(ctx context.Context, in *EmailSenderRequest, opts ...grpc.CallOption) (*EmailSenderResponse, error) {
	out := new(EmailSenderResponse)
	err := c.cc.Invoke(ctx, EmailSenderApi_SendEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailSenderApiServer is the server API for EmailSenderApi service.
// All implementations must embed UnimplementedEmailSenderApiServer
// for forward compatibility
type EmailSenderApiServer interface {
	// SendEmail - RPC method for sending email.
	SendEmail(context.Context, *EmailSenderRequest) (*EmailSenderResponse, error)
	mustEmbedUnimplementedEmailSenderApiServer()
}

// UnimplementedEmailSenderApiServer must be embedded to have forward compatible implementations.
type UnimplementedEmailSenderApiServer struct {
}

func (UnimplementedEmailSenderApiServer) SendEmail(context.Context, *EmailSenderRequest) (*EmailSenderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmail not implemented")
}
func (UnimplementedEmailSenderApiServer) mustEmbedUnimplementedEmailSenderApiServer() {}

// UnsafeEmailSenderApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmailSenderApiServer will
// result in compilation errors.
type UnsafeEmailSenderApiServer interface {
	mustEmbedUnimplementedEmailSenderApiServer()
}

func RegisterEmailSenderApiServer(s grpc.ServiceRegistrar, srv EmailSenderApiServer) {
	s.RegisterService(&EmailSenderApi_ServiceDesc, srv)
}

func _EmailSenderApi_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailSenderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailSenderApiServer).SendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailSenderApi_SendEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailSenderApiServer).SendEmail(ctx, req.(*EmailSenderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EmailSenderApi_ServiceDesc is the grpc.ServiceDesc for EmailSenderApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmailSenderApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.EmailSenderApi",
	HandlerType: (*EmailSenderApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmail",
			Handler:    _EmailSenderApi_SendEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/email-sender.proto",
}
