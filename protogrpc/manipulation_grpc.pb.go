// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: protogrpc/manipulation.proto

package protogrpc

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

// UsersAdminClient is the client API for UsersAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersAdminClient interface {
	AddUser(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error)
	DelUser(ctx context.Context, in *DelRequest, opts ...grpc.CallOption) (*DelResponse, error)
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error)
}

type usersAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersAdminClient(cc grpc.ClientConnInterface) UsersAdminClient {
	return &usersAdminClient{cc}
}

func (c *usersAdminClient) AddUser(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error) {
	out := new(AddResponse)
	err := c.cc.Invoke(ctx, "/protogrpc.UsersAdmin/AddUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersAdminClient) DelUser(ctx context.Context, in *DelRequest, opts ...grpc.CallOption) (*DelResponse, error) {
	out := new(DelResponse)
	err := c.cc.Invoke(ctx, "/protogrpc.UsersAdmin/DelUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersAdminClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, "/protogrpc.UsersAdmin/ListUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersAdminServer is the server API for UsersAdmin service.
// All implementations must embed UnimplementedUsersAdminServer
// for forward compatibility
type UsersAdminServer interface {
	AddUser(context.Context, *AddRequest) (*AddResponse, error)
	DelUser(context.Context, *DelRequest) (*DelResponse, error)
	ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error)
	/*mustEmbedUnimplementedUsersAdminServer()*/
}

// UnimplementedUsersAdminServer must be embedded to have forward compatible implementations.
type UnimplementedUsersAdminServer struct {
}

func (UnimplementedUsersAdminServer) AddUser(context.Context, *AddRequest) (*AddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedUsersAdminServer) DelUser(context.Context, *DelRequest) (*DelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelUser not implemented")
}
func (UnimplementedUsersAdminServer) ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedUsersAdminServer) mustEmbedUnimplementedUsersAdminServer() {}

// UnsafeUsersAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersAdminServer will
// result in compilation errors.
type UnsafeUsersAdminServer interface {
	mustEmbedUnimplementedUsersAdminServer()
}

func RegisterUsersAdminServer(s grpc.ServiceRegistrar, srv UsersAdminServer) {
	s.RegisterService(&UsersAdmin_ServiceDesc, srv)
}

func _UsersAdmin_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersAdminServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protogrpc.UsersAdmin/AddUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersAdminServer).AddUser(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersAdmin_DelUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersAdminServer).DelUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protogrpc.UsersAdmin/DelUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersAdminServer).DelUser(ctx, req.(*DelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersAdmin_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersAdminServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protogrpc.UsersAdmin/ListUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersAdminServer).ListUsers(ctx, req.(*ListUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UsersAdmin_ServiceDesc is the grpc.ServiceDesc for UsersAdmin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UsersAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protogrpc.UsersAdmin",
	HandlerType: (*UsersAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddUser",
			Handler:    _UsersAdmin_AddUser_Handler,
		},
		{
			MethodName: "DelUser",
			Handler:    _UsersAdmin_DelUser_Handler,
		},
		{
			MethodName: "ListUsers",
			Handler:    _UsersAdmin_ListUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protogrpc/manipulation.proto",
}
