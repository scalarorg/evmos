// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:LGPL-3.0-only

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: evmos/revenue/v1/query.proto

package revenuev1

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
	Query_Revenues_FullMethodName           = "/evmos.revenue.v1.Query/Revenues"
	Query_Revenue_FullMethodName            = "/evmos.revenue.v1.Query/Revenue"
	Query_Params_FullMethodName             = "/evmos.revenue.v1.Query/Params"
	Query_DeployerRevenues_FullMethodName   = "/evmos.revenue.v1.Query/DeployerRevenues"
	Query_WithdrawerRevenues_FullMethodName = "/evmos.revenue.v1.Query/WithdrawerRevenues"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Revenues retrieves all registered revenues
	Revenues(ctx context.Context, in *QueryRevenuesRequest, opts ...grpc.CallOption) (*QueryRevenuesResponse, error)
	// Revenue retrieves a registered revenue for a given contract address
	Revenue(ctx context.Context, in *QueryRevenueRequest, opts ...grpc.CallOption) (*QueryRevenueResponse, error)
	// Params retrieves the revenue module params
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// DeployerRevenues retrieves all revenues that a given deployer has
	// registered
	DeployerRevenues(ctx context.Context, in *QueryDeployerRevenuesRequest, opts ...grpc.CallOption) (*QueryDeployerRevenuesResponse, error)
	// WithdrawerRevenues retrieves all revenues with a given withdrawer
	// address
	WithdrawerRevenues(ctx context.Context, in *QueryWithdrawerRevenuesRequest, opts ...grpc.CallOption) (*QueryWithdrawerRevenuesResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Revenues(ctx context.Context, in *QueryRevenuesRequest, opts ...grpc.CallOption) (*QueryRevenuesResponse, error) {
	out := new(QueryRevenuesResponse)
	err := c.cc.Invoke(ctx, Query_Revenues_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Revenue(ctx context.Context, in *QueryRevenueRequest, opts ...grpc.CallOption) (*QueryRevenueResponse, error) {
	out := new(QueryRevenueResponse)
	err := c.cc.Invoke(ctx, Query_Revenue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DeployerRevenues(ctx context.Context, in *QueryDeployerRevenuesRequest, opts ...grpc.CallOption) (*QueryDeployerRevenuesResponse, error) {
	out := new(QueryDeployerRevenuesResponse)
	err := c.cc.Invoke(ctx, Query_DeployerRevenues_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) WithdrawerRevenues(ctx context.Context, in *QueryWithdrawerRevenuesRequest, opts ...grpc.CallOption) (*QueryWithdrawerRevenuesResponse, error) {
	out := new(QueryWithdrawerRevenuesResponse)
	err := c.cc.Invoke(ctx, Query_WithdrawerRevenues_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Revenues retrieves all registered revenues
	Revenues(context.Context, *QueryRevenuesRequest) (*QueryRevenuesResponse, error)
	// Revenue retrieves a registered revenue for a given contract address
	Revenue(context.Context, *QueryRevenueRequest) (*QueryRevenueResponse, error)
	// Params retrieves the revenue module params
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// DeployerRevenues retrieves all revenues that a given deployer has
	// registered
	DeployerRevenues(context.Context, *QueryDeployerRevenuesRequest) (*QueryDeployerRevenuesResponse, error)
	// WithdrawerRevenues retrieves all revenues with a given withdrawer
	// address
	WithdrawerRevenues(context.Context, *QueryWithdrawerRevenuesRequest) (*QueryWithdrawerRevenuesResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Revenues(context.Context, *QueryRevenuesRequest) (*QueryRevenuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revenues not implemented")
}
func (UnimplementedQueryServer) Revenue(context.Context, *QueryRevenueRequest) (*QueryRevenueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revenue not implemented")
}
func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) DeployerRevenues(context.Context, *QueryDeployerRevenuesRequest) (*QueryDeployerRevenuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployerRevenues not implemented")
}
func (UnimplementedQueryServer) WithdrawerRevenues(context.Context, *QueryWithdrawerRevenuesRequest) (*QueryWithdrawerRevenuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawerRevenues not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Revenues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRevenuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Revenues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Revenues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Revenues(ctx, req.(*QueryRevenuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Revenue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRevenueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Revenue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Revenue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Revenue(ctx, req.(*QueryRevenueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DeployerRevenues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDeployerRevenuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DeployerRevenues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_DeployerRevenues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DeployerRevenues(ctx, req.(*QueryDeployerRevenuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_WithdrawerRevenues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryWithdrawerRevenuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).WithdrawerRevenues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_WithdrawerRevenues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).WithdrawerRevenues(ctx, req.(*QueryWithdrawerRevenuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "evmos.revenue.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Revenues",
			Handler:    _Query_Revenues_Handler,
		},
		{
			MethodName: "Revenue",
			Handler:    _Query_Revenue_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "DeployerRevenues",
			Handler:    _Query_DeployerRevenues_Handler,
		},
		{
			MethodName: "WithdrawerRevenues",
			Handler:    _Query_WithdrawerRevenues_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "evmos/revenue/v1/query.proto",
}