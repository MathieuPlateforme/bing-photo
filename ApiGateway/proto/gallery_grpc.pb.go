// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: proto/gallery.proto

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
	AlbumService_CreateAlbum_FullMethodName     = "/proto.AlbumService/CreateAlbum"
	AlbumService_GetAlbumsByUser_FullMethodName = "/proto.AlbumService/GetAlbumsByUser"
	AlbumService_UpdateAlbum_FullMethodName     = "/proto.AlbumService/UpdateAlbum"
	AlbumService_DeleteAlbum_FullMethodName     = "/proto.AlbumService/DeleteAlbum"
	AlbumService_GetPrivateAlbum_FullMethodName = "/proto.AlbumService/GetPrivateAlbum"
)

// AlbumServiceClient is the client API for AlbumService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AlbumServiceClient interface {
	CreateAlbum(ctx context.Context, in *CreateAlbumRequest, opts ...grpc.CallOption) (*CreateAlbumResponse, error)
	GetAlbumsByUser(ctx context.Context, in *GetAlbumsByUserRequest, opts ...grpc.CallOption) (*GetAlbumsByUserResponse, error)
	UpdateAlbum(ctx context.Context, in *UpdateAlbumRequest, opts ...grpc.CallOption) (*UpdateAlbumResponse, error)
	DeleteAlbum(ctx context.Context, in *DeleteAlbumRequest, opts ...grpc.CallOption) (*DeleteAlbumResponse, error)
	GetPrivateAlbum(ctx context.Context, in *GetPrivateAlbumRequest, opts ...grpc.CallOption) (*GetPrivateAlbumResponse, error)
}

type albumServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAlbumServiceClient(cc grpc.ClientConnInterface) AlbumServiceClient {
	return &albumServiceClient{cc}
}

func (c *albumServiceClient) CreateAlbum(ctx context.Context, in *CreateAlbumRequest, opts ...grpc.CallOption) (*CreateAlbumResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAlbumResponse)
	err := c.cc.Invoke(ctx, AlbumService_CreateAlbum_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *albumServiceClient) GetAlbumsByUser(ctx context.Context, in *GetAlbumsByUserRequest, opts ...grpc.CallOption) (*GetAlbumsByUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAlbumsByUserResponse)
	err := c.cc.Invoke(ctx, AlbumService_GetAlbumsByUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *albumServiceClient) UpdateAlbum(ctx context.Context, in *UpdateAlbumRequest, opts ...grpc.CallOption) (*UpdateAlbumResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAlbumResponse)
	err := c.cc.Invoke(ctx, AlbumService_UpdateAlbum_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *albumServiceClient) DeleteAlbum(ctx context.Context, in *DeleteAlbumRequest, opts ...grpc.CallOption) (*DeleteAlbumResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAlbumResponse)
	err := c.cc.Invoke(ctx, AlbumService_DeleteAlbum_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *albumServiceClient) GetPrivateAlbum(ctx context.Context, in *GetPrivateAlbumRequest, opts ...grpc.CallOption) (*GetPrivateAlbumResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPrivateAlbumResponse)
	err := c.cc.Invoke(ctx, AlbumService_GetPrivateAlbum_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlbumServiceServer is the server API for AlbumService service.
// All implementations must embed UnimplementedAlbumServiceServer
// for forward compatibility.
type AlbumServiceServer interface {
	CreateAlbum(context.Context, *CreateAlbumRequest) (*CreateAlbumResponse, error)
	GetAlbumsByUser(context.Context, *GetAlbumsByUserRequest) (*GetAlbumsByUserResponse, error)
	UpdateAlbum(context.Context, *UpdateAlbumRequest) (*UpdateAlbumResponse, error)
	DeleteAlbum(context.Context, *DeleteAlbumRequest) (*DeleteAlbumResponse, error)
	GetPrivateAlbum(context.Context, *GetPrivateAlbumRequest) (*GetPrivateAlbumResponse, error)
	mustEmbedUnimplementedAlbumServiceServer()
}

// UnimplementedAlbumServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAlbumServiceServer struct{}

func (UnimplementedAlbumServiceServer) CreateAlbum(context.Context, *CreateAlbumRequest) (*CreateAlbumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAlbum not implemented")
}
func (UnimplementedAlbumServiceServer) GetAlbumsByUser(context.Context, *GetAlbumsByUserRequest) (*GetAlbumsByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlbumsByUser not implemented")
}
func (UnimplementedAlbumServiceServer) UpdateAlbum(context.Context, *UpdateAlbumRequest) (*UpdateAlbumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAlbum not implemented")
}
func (UnimplementedAlbumServiceServer) DeleteAlbum(context.Context, *DeleteAlbumRequest) (*DeleteAlbumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAlbum not implemented")
}
func (UnimplementedAlbumServiceServer) GetPrivateAlbum(context.Context, *GetPrivateAlbumRequest) (*GetPrivateAlbumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrivateAlbum not implemented")
}
func (UnimplementedAlbumServiceServer) mustEmbedUnimplementedAlbumServiceServer() {}
func (UnimplementedAlbumServiceServer) testEmbeddedByValue()                      {}

// UnsafeAlbumServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AlbumServiceServer will
// result in compilation errors.
type UnsafeAlbumServiceServer interface {
	mustEmbedUnimplementedAlbumServiceServer()
}

func RegisterAlbumServiceServer(s grpc.ServiceRegistrar, srv AlbumServiceServer) {
	// If the following call pancis, it indicates UnimplementedAlbumServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AlbumService_ServiceDesc, srv)
}

func _AlbumService_CreateAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlbumServiceServer).CreateAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlbumService_CreateAlbum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlbumServiceServer).CreateAlbum(ctx, req.(*CreateAlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlbumService_GetAlbumsByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlbumsByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlbumServiceServer).GetAlbumsByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlbumService_GetAlbumsByUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlbumServiceServer).GetAlbumsByUser(ctx, req.(*GetAlbumsByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlbumService_UpdateAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlbumServiceServer).UpdateAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlbumService_UpdateAlbum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlbumServiceServer).UpdateAlbum(ctx, req.(*UpdateAlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlbumService_DeleteAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlbumServiceServer).DeleteAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlbumService_DeleteAlbum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlbumServiceServer).DeleteAlbum(ctx, req.(*DeleteAlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlbumService_GetPrivateAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPrivateAlbumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlbumServiceServer).GetPrivateAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlbumService_GetPrivateAlbum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlbumServiceServer).GetPrivateAlbum(ctx, req.(*GetPrivateAlbumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AlbumService_ServiceDesc is the grpc.ServiceDesc for AlbumService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AlbumService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AlbumService",
	HandlerType: (*AlbumServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAlbum",
			Handler:    _AlbumService_CreateAlbum_Handler,
		},
		{
			MethodName: "GetAlbumsByUser",
			Handler:    _AlbumService_GetAlbumsByUser_Handler,
		},
		{
			MethodName: "UpdateAlbum",
			Handler:    _AlbumService_UpdateAlbum_Handler,
		},
		{
			MethodName: "DeleteAlbum",
			Handler:    _AlbumService_DeleteAlbum_Handler,
		},
		{
			MethodName: "GetPrivateAlbum",
			Handler:    _AlbumService_GetPrivateAlbum_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gallery.proto",
}

const (
	MediaService_AddMedia_FullMethodName           = "/proto.MediaService/AddMedia"
	MediaService_GetMediaByUser_FullMethodName     = "/proto.MediaService/GetMediaByUser"
	MediaService_MarkAsPrivate_FullMethodName      = "/proto.MediaService/MarkAsPrivate"
	MediaService_GetPrivateMedia_FullMethodName    = "/proto.MediaService/GetPrivateMedia"
	MediaService_DownloadMedia_FullMethodName      = "/proto.MediaService/DownloadMedia"
	MediaService_DeleteMedia_FullMethodName        = "/proto.MediaService/DeleteMedia"
	MediaService_DetectSimilarMedia_FullMethodName = "/proto.MediaService/DetectSimilarMedia"
)

// MediaServiceClient is the client API for MediaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MediaServiceClient interface {
	AddMedia(ctx context.Context, in *AddMediaRequest, opts ...grpc.CallOption) (*AddMediaResponse, error)
	GetMediaByUser(ctx context.Context, in *GetMediaByUserRequest, opts ...grpc.CallOption) (*GetMediaByUserResponse, error)
	MarkAsPrivate(ctx context.Context, in *MarkAsPrivateRequest, opts ...grpc.CallOption) (*MarkAsPrivateResponse, error)
	GetPrivateMedia(ctx context.Context, in *GetPrivateMediaRequest, opts ...grpc.CallOption) (*GetPrivateMediaResponse, error)
	DownloadMedia(ctx context.Context, in *DownloadMediaRequest, opts ...grpc.CallOption) (*DownloadMediaResponse, error)
	DeleteMedia(ctx context.Context, in *DeleteMediaRequest, opts ...grpc.CallOption) (*DeleteMediaResponse, error)
	DetectSimilarMedia(ctx context.Context, in *DetectSimilarMediaRequest, opts ...grpc.CallOption) (*DetectSimilarMediaResponse, error)
}

type mediaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMediaServiceClient(cc grpc.ClientConnInterface) MediaServiceClient {
	return &mediaServiceClient{cc}
}

func (c *mediaServiceClient) AddMedia(ctx context.Context, in *AddMediaRequest, opts ...grpc.CallOption) (*AddMediaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddMediaResponse)
	err := c.cc.Invoke(ctx, MediaService_AddMedia_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) GetMediaByUser(ctx context.Context, in *GetMediaByUserRequest, opts ...grpc.CallOption) (*GetMediaByUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMediaByUserResponse)
	err := c.cc.Invoke(ctx, MediaService_GetMediaByUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) MarkAsPrivate(ctx context.Context, in *MarkAsPrivateRequest, opts ...grpc.CallOption) (*MarkAsPrivateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MarkAsPrivateResponse)
	err := c.cc.Invoke(ctx, MediaService_MarkAsPrivate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) GetPrivateMedia(ctx context.Context, in *GetPrivateMediaRequest, opts ...grpc.CallOption) (*GetPrivateMediaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPrivateMediaResponse)
	err := c.cc.Invoke(ctx, MediaService_GetPrivateMedia_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) DownloadMedia(ctx context.Context, in *DownloadMediaRequest, opts ...grpc.CallOption) (*DownloadMediaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DownloadMediaResponse)
	err := c.cc.Invoke(ctx, MediaService_DownloadMedia_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) DeleteMedia(ctx context.Context, in *DeleteMediaRequest, opts ...grpc.CallOption) (*DeleteMediaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMediaResponse)
	err := c.cc.Invoke(ctx, MediaService_DeleteMedia_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaServiceClient) DetectSimilarMedia(ctx context.Context, in *DetectSimilarMediaRequest, opts ...grpc.CallOption) (*DetectSimilarMediaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DetectSimilarMediaResponse)
	err := c.cc.Invoke(ctx, MediaService_DetectSimilarMedia_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MediaServiceServer is the server API for MediaService service.
// All implementations must embed UnimplementedMediaServiceServer
// for forward compatibility.
type MediaServiceServer interface {
	AddMedia(context.Context, *AddMediaRequest) (*AddMediaResponse, error)
	GetMediaByUser(context.Context, *GetMediaByUserRequest) (*GetMediaByUserResponse, error)
	MarkAsPrivate(context.Context, *MarkAsPrivateRequest) (*MarkAsPrivateResponse, error)
	GetPrivateMedia(context.Context, *GetPrivateMediaRequest) (*GetPrivateMediaResponse, error)
	DownloadMedia(context.Context, *DownloadMediaRequest) (*DownloadMediaResponse, error)
	DeleteMedia(context.Context, *DeleteMediaRequest) (*DeleteMediaResponse, error)
	DetectSimilarMedia(context.Context, *DetectSimilarMediaRequest) (*DetectSimilarMediaResponse, error)
	mustEmbedUnimplementedMediaServiceServer()
}

// UnimplementedMediaServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMediaServiceServer struct{}

func (UnimplementedMediaServiceServer) AddMedia(context.Context, *AddMediaRequest) (*AddMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMedia not implemented")
}
func (UnimplementedMediaServiceServer) GetMediaByUser(context.Context, *GetMediaByUserRequest) (*GetMediaByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMediaByUser not implemented")
}
func (UnimplementedMediaServiceServer) MarkAsPrivate(context.Context, *MarkAsPrivateRequest) (*MarkAsPrivateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkAsPrivate not implemented")
}
func (UnimplementedMediaServiceServer) GetPrivateMedia(context.Context, *GetPrivateMediaRequest) (*GetPrivateMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrivateMedia not implemented")
}
func (UnimplementedMediaServiceServer) DownloadMedia(context.Context, *DownloadMediaRequest) (*DownloadMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadMedia not implemented")
}
func (UnimplementedMediaServiceServer) DeleteMedia(context.Context, *DeleteMediaRequest) (*DeleteMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMedia not implemented")
}
func (UnimplementedMediaServiceServer) DetectSimilarMedia(context.Context, *DetectSimilarMediaRequest) (*DetectSimilarMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetectSimilarMedia not implemented")
}
func (UnimplementedMediaServiceServer) mustEmbedUnimplementedMediaServiceServer() {}
func (UnimplementedMediaServiceServer) testEmbeddedByValue()                      {}

// UnsafeMediaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MediaServiceServer will
// result in compilation errors.
type UnsafeMediaServiceServer interface {
	mustEmbedUnimplementedMediaServiceServer()
}

func RegisterMediaServiceServer(s grpc.ServiceRegistrar, srv MediaServiceServer) {
	// If the following call pancis, it indicates UnimplementedMediaServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MediaService_ServiceDesc, srv)
}

func _MediaService_AddMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).AddMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MediaService_AddMedia_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).AddMedia(ctx, req.(*AddMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_GetMediaByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMediaByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).GetMediaByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MediaService_GetMediaByUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).GetMediaByUser(ctx, req.(*GetMediaByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_MarkAsPrivate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkAsPrivateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).MarkAsPrivate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MediaService_MarkAsPrivate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).MarkAsPrivate(ctx, req.(*MarkAsPrivateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_GetPrivateMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPrivateMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).GetPrivateMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MediaService_GetPrivateMedia_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).GetPrivateMedia(ctx, req.(*GetPrivateMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_DownloadMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).DownloadMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MediaService_DownloadMedia_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).DownloadMedia(ctx, req.(*DownloadMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_DeleteMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).DeleteMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MediaService_DeleteMedia_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).DeleteMedia(ctx, req.(*DeleteMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MediaService_DetectSimilarMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetectSimilarMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServiceServer).DetectSimilarMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MediaService_DetectSimilarMedia_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServiceServer).DetectSimilarMedia(ctx, req.(*DetectSimilarMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MediaService_ServiceDesc is the grpc.ServiceDesc for MediaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MediaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MediaService",
	HandlerType: (*MediaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddMedia",
			Handler:    _MediaService_AddMedia_Handler,
		},
		{
			MethodName: "GetMediaByUser",
			Handler:    _MediaService_GetMediaByUser_Handler,
		},
		{
			MethodName: "MarkAsPrivate",
			Handler:    _MediaService_MarkAsPrivate_Handler,
		},
		{
			MethodName: "GetPrivateMedia",
			Handler:    _MediaService_GetPrivateMedia_Handler,
		},
		{
			MethodName: "DownloadMedia",
			Handler:    _MediaService_DownloadMedia_Handler,
		},
		{
			MethodName: "DeleteMedia",
			Handler:    _MediaService_DeleteMedia_Handler,
		},
		{
			MethodName: "DetectSimilarMedia",
			Handler:    _MediaService_DetectSimilarMedia_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gallery.proto",
}

const (
	UserService_CreateUser_FullMethodName = "/proto.UserService/CreateUser"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, UserService_CreateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility.
type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}
func (UnimplementedUserServiceServer) testEmbeddedByValue()                     {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	// If the following call pancis, it indicates UnimplementedUserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gallery.proto",
}
