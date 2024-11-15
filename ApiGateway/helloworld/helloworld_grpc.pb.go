// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0--rc2
// source: helloworld.proto

package helloworld

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
	HelloWorldService_SayHello_FullMethodName = "/HelloWorldService/SayHello"
)

// HelloWorldServiceClient is the client API for HelloWorldService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloWorldServiceClient interface {
	SayHello(ctx context.Context, in *HelloWorld, opts ...grpc.CallOption) (*HelloWorld, error)
}

type helloWorldServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloWorldServiceClient(cc grpc.ClientConnInterface) HelloWorldServiceClient {
	return &helloWorldServiceClient{cc}
}

func (c *helloWorldServiceClient) SayHello(ctx context.Context, in *HelloWorld, opts ...grpc.CallOption) (*HelloWorld, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HelloWorld)
	err := c.cc.Invoke(ctx, HelloWorldService_SayHello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HelloWorldServiceServer is the server API for HelloWorldService service.
// All implementations must embed UnimplementedHelloWorldServiceServer
// for forward compatibility.
type HelloWorldServiceServer interface {
	SayHello(context.Context, *HelloWorld) (*HelloWorld, error)
	mustEmbedUnimplementedHelloWorldServiceServer()
}

// UnimplementedHelloWorldServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedHelloWorldServiceServer struct{}

func (UnimplementedHelloWorldServiceServer) SayHello(context.Context, *HelloWorld) (*HelloWorld, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedHelloWorldServiceServer) mustEmbedUnimplementedHelloWorldServiceServer() {}
func (UnimplementedHelloWorldServiceServer) testEmbeddedByValue()                           {}

// UnsafeHelloWorldServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloWorldServiceServer will
// result in compilation errors.
type UnsafeHelloWorldServiceServer interface {
	mustEmbedUnimplementedHelloWorldServiceServer()
}

func RegisterHelloWorldServiceServer(s grpc.ServiceRegistrar, srv HelloWorldServiceServer) {
	// If the following call pancis, it indicates UnimplementedHelloWorldServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&HelloWorldService_ServiceDesc, srv)
}

func _HelloWorldService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloWorld)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloWorldServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HelloWorldService_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloWorldServiceServer).SayHello(ctx, req.(*HelloWorld))
	}
	return interceptor(ctx, in, info, handler)
}

// HelloWorldService_ServiceDesc is the grpc.ServiceDesc for HelloWorldService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HelloWorldService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "HelloWorldService",
	HandlerType: (*HelloWorldServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _HelloWorldService_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "helloworld.proto",
}
