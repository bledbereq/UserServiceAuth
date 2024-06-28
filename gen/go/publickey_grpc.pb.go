// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: publickey.proto

package ssov1

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

// GetPublicKeyClient is the client API for GetPublicKey service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GetPublicKeyClient interface {
	PublicKey(ctx context.Context, in *PublicKeyRequest, opts ...grpc.CallOption) (*PublicKeyResponse, error)
}

type getPublicKeyClient struct {
	cc grpc.ClientConnInterface
}

func NewGetPublicKeyClient(cc grpc.ClientConnInterface) GetPublicKeyClient {
	return &getPublicKeyClient{cc}
}

func (c *getPublicKeyClient) PublicKey(ctx context.Context, in *PublicKeyRequest, opts ...grpc.CallOption) (*PublicKeyResponse, error) {
	out := new(PublicKeyResponse)
	err := c.cc.Invoke(ctx, "/publickey.GetPublicKey/PublicKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetPublicKeyServer is the server API for GetPublicKey service.
// All implementations must embed UnimplementedGetPublicKeyServer
// for forward compatibility
type GetPublicKeyServer interface {
	PublicKey(context.Context, *PublicKeyRequest) (*PublicKeyResponse, error)
	mustEmbedUnimplementedGetPublicKeyServer()
}

// UnimplementedGetPublicKeyServer must be embedded to have forward compatible implementations.
type UnimplementedGetPublicKeyServer struct {
}

func (UnimplementedGetPublicKeyServer) PublicKey(context.Context, *PublicKeyRequest) (*PublicKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublicKey not implemented")
}
func (UnimplementedGetPublicKeyServer) mustEmbedUnimplementedGetPublicKeyServer() {}

// UnsafeGetPublicKeyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GetPublicKeyServer will
// result in compilation errors.
type UnsafeGetPublicKeyServer interface {
	mustEmbedUnimplementedGetPublicKeyServer()
}

func RegisterGetPublicKeyServer(s grpc.ServiceRegistrar, srv GetPublicKeyServer) {
	s.RegisterService(&GetPublicKey_ServiceDesc, srv)
}

func _GetPublicKey_PublicKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GetPublicKeyServer).PublicKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/publickey.GetPublicKey/PublicKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GetPublicKeyServer).PublicKey(ctx, req.(*PublicKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GetPublicKey_ServiceDesc is the grpc.ServiceDesc for GetPublicKey service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GetPublicKey_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "publickey.GetPublicKey",
	HandlerType: (*GetPublicKeyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublicKey",
			Handler:    _GetPublicKey_PublicKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "publickey.proto",
}
