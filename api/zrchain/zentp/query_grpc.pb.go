// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: zrchain/zentp/query.proto

package zentp

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
	Query_Params_FullMethodName                = "/zrchain.zentp.Query/Params"
	Query_Mints_FullMethodName                 = "/zrchain.zentp.Query/Mints"
	Query_Burns_FullMethodName                 = "/zrchain.zentp.Query/Burns"
	Query_Stats_FullMethodName                 = "/zrchain.zentp.Query/Stats"
	Query_QuerySolanaROCKSupply_FullMethodName = "/zrchain.zentp.Query/QuerySolanaROCKSupply"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Query defines the gRPC querier service.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a list of Mints.
	Mints(ctx context.Context, in *QueryMintsRequest, opts ...grpc.CallOption) (*QueryMintsResponse, error)
	// Queries a list of Burns items.
	Burns(ctx context.Context, in *QueryBurnsRequest, opts ...grpc.CallOption) (*QueryBurnsResponse, error)
	// Stats queries the total amounts of mints and burns for an address
	Stats(ctx context.Context, in *QueryStatsRequest, opts ...grpc.CallOption) (*QueryStatsResponse, error)
	// QuerySolanaROCKSupply queries the amount of ROCK on Solana.
	QuerySolanaROCKSupply(ctx context.Context, in *QuerySolanaROCKSupplyRequest, opts ...grpc.CallOption) (*QuerySolanaROCKSupplyResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Mints(ctx context.Context, in *QueryMintsRequest, opts ...grpc.CallOption) (*QueryMintsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryMintsResponse)
	err := c.cc.Invoke(ctx, Query_Mints_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Burns(ctx context.Context, in *QueryBurnsRequest, opts ...grpc.CallOption) (*QueryBurnsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryBurnsResponse)
	err := c.cc.Invoke(ctx, Query_Burns_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Stats(ctx context.Context, in *QueryStatsRequest, opts ...grpc.CallOption) (*QueryStatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryStatsResponse)
	err := c.cc.Invoke(ctx, Query_Stats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QuerySolanaROCKSupply(ctx context.Context, in *QuerySolanaROCKSupplyRequest, opts ...grpc.CallOption) (*QuerySolanaROCKSupplyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuerySolanaROCKSupplyResponse)
	err := c.cc.Invoke(ctx, Query_QuerySolanaROCKSupply_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility.
//
// Query defines the gRPC querier service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of Mints.
	Mints(context.Context, *QueryMintsRequest) (*QueryMintsResponse, error)
	// Queries a list of Burns items.
	Burns(context.Context, *QueryBurnsRequest) (*QueryBurnsResponse, error)
	// Stats queries the total amounts of mints and burns for an address
	Stats(context.Context, *QueryStatsRequest) (*QueryStatsResponse, error)
	// QuerySolanaROCKSupply queries the amount of ROCK on Solana.
	QuerySolanaROCKSupply(context.Context, *QuerySolanaROCKSupplyRequest) (*QuerySolanaROCKSupplyResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) Mints(context.Context, *QueryMintsRequest) (*QueryMintsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mints not implemented")
}
func (UnimplementedQueryServer) Burns(context.Context, *QueryBurnsRequest) (*QueryBurnsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Burns not implemented")
}
func (UnimplementedQueryServer) Stats(context.Context, *QueryStatsRequest) (*QueryStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedQueryServer) QuerySolanaROCKSupply(context.Context, *QuerySolanaROCKSupplyRequest) (*QuerySolanaROCKSupplyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuerySolanaROCKSupply not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}
func (UnimplementedQueryServer) testEmbeddedByValue()               {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	// If the following call pancis, it indicates UnimplementedQueryServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Query_ServiceDesc, srv)
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

func _Query_Mints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Mints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Mints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Mints(ctx, req.(*QueryMintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Burns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBurnsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Burns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Burns_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Burns(ctx, req.(*QueryBurnsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Stats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Stats(ctx, req.(*QueryStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QuerySolanaROCKSupply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySolanaROCKSupplyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QuerySolanaROCKSupply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QuerySolanaROCKSupply_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QuerySolanaROCKSupply(ctx, req.(*QuerySolanaROCKSupplyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "zrchain.zentp.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Mints",
			Handler:    _Query_Mints_Handler,
		},
		{
			MethodName: "Burns",
			Handler:    _Query_Burns_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _Query_Stats_Handler,
		},
		{
			MethodName: "QuerySolanaROCKSupply",
			Handler:    _Query_QuerySolanaROCKSupply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "zrchain/zentp/query.proto",
}
