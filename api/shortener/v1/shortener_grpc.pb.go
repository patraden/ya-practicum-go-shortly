// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: shortener/v1/shortener.proto

package api

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
	URLShortenerService_ShortenURL_FullMethodName     = "/shortener.v1.URLShortenerService/ShortenURL"
	URLShortenerService_GetOriginalURL_FullMethodName = "/shortener.v1.URLShortenerService/GetOriginalURL"
)

// URLShortenerServiceClient is the client API for URLShortenerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerServiceClient interface {
	ShortenURL(ctx context.Context, in *ShortenURLRequest, opts ...grpc.CallOption) (*ShortenURLResponse, error)
	GetOriginalURL(ctx context.Context, in *GetOriginalURLRequest, opts ...grpc.CallOption) (*GetOriginalURLResponse, error)
}

type uRLShortenerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerServiceClient(cc grpc.ClientConnInterface) URLShortenerServiceClient {
	return &uRLShortenerServiceClient{cc}
}

func (c *uRLShortenerServiceClient) ShortenURL(ctx context.Context, in *ShortenURLRequest, opts ...grpc.CallOption) (*ShortenURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenURLResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_ShortenURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) GetOriginalURL(ctx context.Context, in *GetOriginalURLRequest, opts ...grpc.CallOption) (*GetOriginalURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetOriginalURLResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_GetOriginalURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServiceServer is the server API for URLShortenerService service.
// All implementations must embed UnimplementedURLShortenerServiceServer
// for forward compatibility.
type URLShortenerServiceServer interface {
	ShortenURL(context.Context, *ShortenURLRequest) (*ShortenURLResponse, error)
	GetOriginalURL(context.Context, *GetOriginalURLRequest) (*GetOriginalURLResponse, error)
	mustEmbedUnimplementedURLShortenerServiceServer()
}

// UnimplementedURLShortenerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLShortenerServiceServer struct{}

func (UnimplementedURLShortenerServiceServer) ShortenURL(context.Context, *ShortenURLRequest) (*ShortenURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenURL not implemented")
}
func (UnimplementedURLShortenerServiceServer) GetOriginalURL(context.Context, *GetOriginalURLRequest) (*GetOriginalURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedURLShortenerServiceServer) mustEmbedUnimplementedURLShortenerServiceServer() {}
func (UnimplementedURLShortenerServiceServer) testEmbeddedByValue()                             {}

// UnsafeURLShortenerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServiceServer will
// result in compilation errors.
type UnsafeURLShortenerServiceServer interface {
	mustEmbedUnimplementedURLShortenerServiceServer()
}

func RegisterURLShortenerServiceServer(s grpc.ServiceRegistrar, srv URLShortenerServiceServer) {
	// If the following call pancis, it indicates UnimplementedURLShortenerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLShortenerService_ServiceDesc, srv)
}

func _URLShortenerService_ShortenURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).ShortenURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_ShortenURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).ShortenURL(ctx, req.(*ShortenURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOriginalURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_GetOriginalURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).GetOriginalURL(ctx, req.(*GetOriginalURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortenerService_ServiceDesc is the grpc.ServiceDesc for URLShortenerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortenerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.v1.URLShortenerService",
	HandlerType: (*URLShortenerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenURL",
			Handler:    _URLShortenerService_ShortenURL_Handler,
		},
		{
			MethodName: "GetOriginalURL",
			Handler:    _URLShortenerService_GetOriginalURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shortener/v1/shortener.proto",
}
