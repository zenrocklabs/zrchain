// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// SidecarServiceClient is the client API for SidecarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SidecarServiceClient interface {
	GetSidecarState(ctx context.Context, in *SidecarStateRequest, opts ...grpc.CallOption) (*SidecarStateResponse, error)
	GetSidecarStateByEthHeight(ctx context.Context, in *SidecarStateByEthHeightRequest, opts ...grpc.CallOption) (*SidecarStateResponse, error)
	GetBitcoinBlockHeaderByHeight(ctx context.Context, in *BitcoinBlockHeaderByHeightRequest, opts ...grpc.CallOption) (*BitcoinBlockHeaderResponse, error)
	GetLatestBitcoinBlockHeader(ctx context.Context, in *LatestBitcoinBlockHeaderRequest, opts ...grpc.CallOption) (*BitcoinBlockHeaderResponse, error)
	GetLatestEthereumNonceForAccount(ctx context.Context, in *LatestEthereumNonceForAccountRequest, opts ...grpc.CallOption) (*LatestEthereumNonceForAccountResponse, error)
	GetSolanaRecentBlockhash(ctx context.Context, in *SolanaRecentBlockhashRequest, opts ...grpc.CallOption) (*SolanaRecentBlockhashResponse, error)
	GetSolanaAccountInfo(ctx context.Context, in *SolanaAccountInfoRequest, opts ...grpc.CallOption) (*SolanaAccountInfoResponse, error)
}

type sidecarServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSidecarServiceClient(cc grpc.ClientConnInterface) SidecarServiceClient {
	return &sidecarServiceClient{cc}
}

func (c *sidecarServiceClient) GetSidecarState(ctx context.Context, in *SidecarStateRequest, opts ...grpc.CallOption) (*SidecarStateResponse, error) {
	out := new(SidecarStateResponse)
	err := c.cc.Invoke(ctx, "/api.SidecarService/GetSidecarState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sidecarServiceClient) GetSidecarStateByEthHeight(ctx context.Context, in *SidecarStateByEthHeightRequest, opts ...grpc.CallOption) (*SidecarStateResponse, error) {
	out := new(SidecarStateResponse)
	err := c.cc.Invoke(ctx, "/api.SidecarService/GetSidecarStateByEthHeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sidecarServiceClient) GetBitcoinBlockHeaderByHeight(ctx context.Context, in *BitcoinBlockHeaderByHeightRequest, opts ...grpc.CallOption) (*BitcoinBlockHeaderResponse, error) {
	out := new(BitcoinBlockHeaderResponse)
	err := c.cc.Invoke(ctx, "/api.SidecarService/GetBitcoinBlockHeaderByHeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sidecarServiceClient) GetLatestBitcoinBlockHeader(ctx context.Context, in *LatestBitcoinBlockHeaderRequest, opts ...grpc.CallOption) (*BitcoinBlockHeaderResponse, error) {
	out := new(BitcoinBlockHeaderResponse)
	err := c.cc.Invoke(ctx, "/api.SidecarService/GetLatestBitcoinBlockHeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sidecarServiceClient) GetLatestEthereumNonceForAccount(ctx context.Context, in *LatestEthereumNonceForAccountRequest, opts ...grpc.CallOption) (*LatestEthereumNonceForAccountResponse, error) {
	out := new(LatestEthereumNonceForAccountResponse)
	err := c.cc.Invoke(ctx, "/api.SidecarService/GetLatestEthereumNonceForAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sidecarServiceClient) GetSolanaRecentBlockhash(ctx context.Context, in *SolanaRecentBlockhashRequest, opts ...grpc.CallOption) (*SolanaRecentBlockhashResponse, error) {
	out := new(SolanaRecentBlockhashResponse)
	err := c.cc.Invoke(ctx, "/api.SidecarService/GetSolanaRecentBlockhash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sidecarServiceClient) GetSolanaAccountInfo(ctx context.Context, in *SolanaAccountInfoRequest, opts ...grpc.CallOption) (*SolanaAccountInfoResponse, error) {
	out := new(SolanaAccountInfoResponse)
	err := c.cc.Invoke(ctx, "/api.SidecarService/GetSolanaAccountInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SidecarServiceServer is the server API for SidecarService service.
// All implementations must embed UnimplementedSidecarServiceServer
// for forward compatibility
type SidecarServiceServer interface {
	GetSidecarState(context.Context, *SidecarStateRequest) (*SidecarStateResponse, error)
	GetSidecarStateByEthHeight(context.Context, *SidecarStateByEthHeightRequest) (*SidecarStateResponse, error)
	GetBitcoinBlockHeaderByHeight(context.Context, *BitcoinBlockHeaderByHeightRequest) (*BitcoinBlockHeaderResponse, error)
	GetLatestBitcoinBlockHeader(context.Context, *LatestBitcoinBlockHeaderRequest) (*BitcoinBlockHeaderResponse, error)
	GetLatestEthereumNonceForAccount(context.Context, *LatestEthereumNonceForAccountRequest) (*LatestEthereumNonceForAccountResponse, error)
	GetSolanaRecentBlockhash(context.Context, *SolanaRecentBlockhashRequest) (*SolanaRecentBlockhashResponse, error)
	GetSolanaAccountInfo(context.Context, *SolanaAccountInfoRequest) (*SolanaAccountInfoResponse, error)
	mustEmbedUnimplementedSidecarServiceServer()
}

// UnimplementedSidecarServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSidecarServiceServer struct {
}

func (UnimplementedSidecarServiceServer) GetSidecarState(context.Context, *SidecarStateRequest) (*SidecarStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSidecarState not implemented")
}
func (UnimplementedSidecarServiceServer) GetSidecarStateByEthHeight(context.Context, *SidecarStateByEthHeightRequest) (*SidecarStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSidecarStateByEthHeight not implemented")
}
func (UnimplementedSidecarServiceServer) GetBitcoinBlockHeaderByHeight(context.Context, *BitcoinBlockHeaderByHeightRequest) (*BitcoinBlockHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBitcoinBlockHeaderByHeight not implemented")
}
func (UnimplementedSidecarServiceServer) GetLatestBitcoinBlockHeader(context.Context, *LatestBitcoinBlockHeaderRequest) (*BitcoinBlockHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLatestBitcoinBlockHeader not implemented")
}
func (UnimplementedSidecarServiceServer) GetLatestEthereumNonceForAccount(context.Context, *LatestEthereumNonceForAccountRequest) (*LatestEthereumNonceForAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLatestEthereumNonceForAccount not implemented")
}
func (UnimplementedSidecarServiceServer) GetSolanaRecentBlockhash(context.Context, *SolanaRecentBlockhashRequest) (*SolanaRecentBlockhashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSolanaRecentBlockhash not implemented")
}
func (UnimplementedSidecarServiceServer) GetSolanaAccountInfo(context.Context, *SolanaAccountInfoRequest) (*SolanaAccountInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSolanaAccountInfo not implemented")
}
func (UnimplementedSidecarServiceServer) mustEmbedUnimplementedSidecarServiceServer() {}

// UnsafeSidecarServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SidecarServiceServer will
// result in compilation errors.
type UnsafeSidecarServiceServer interface {
	mustEmbedUnimplementedSidecarServiceServer()
}

func RegisterSidecarServiceServer(s grpc.ServiceRegistrar, srv SidecarServiceServer) {
	s.RegisterService(&SidecarService_ServiceDesc, srv)
}

func _SidecarService_GetSidecarState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SidecarStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServiceServer).GetSidecarState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SidecarService/GetSidecarState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServiceServer).GetSidecarState(ctx, req.(*SidecarStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SidecarService_GetSidecarStateByEthHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SidecarStateByEthHeightRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServiceServer).GetSidecarStateByEthHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SidecarService/GetSidecarStateByEthHeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServiceServer).GetSidecarStateByEthHeight(ctx, req.(*SidecarStateByEthHeightRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SidecarService_GetBitcoinBlockHeaderByHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BitcoinBlockHeaderByHeightRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServiceServer).GetBitcoinBlockHeaderByHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SidecarService/GetBitcoinBlockHeaderByHeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServiceServer).GetBitcoinBlockHeaderByHeight(ctx, req.(*BitcoinBlockHeaderByHeightRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SidecarService_GetLatestBitcoinBlockHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LatestBitcoinBlockHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServiceServer).GetLatestBitcoinBlockHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SidecarService/GetLatestBitcoinBlockHeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServiceServer).GetLatestBitcoinBlockHeader(ctx, req.(*LatestBitcoinBlockHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SidecarService_GetLatestEthereumNonceForAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LatestEthereumNonceForAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServiceServer).GetLatestEthereumNonceForAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SidecarService/GetLatestEthereumNonceForAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServiceServer).GetLatestEthereumNonceForAccount(ctx, req.(*LatestEthereumNonceForAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SidecarService_GetSolanaRecentBlockhash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SolanaRecentBlockhashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServiceServer).GetSolanaRecentBlockhash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SidecarService/GetSolanaRecentBlockhash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServiceServer).GetSolanaRecentBlockhash(ctx, req.(*SolanaRecentBlockhashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SidecarService_GetSolanaAccountInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SolanaAccountInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServiceServer).GetSolanaAccountInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.SidecarService/GetSolanaAccountInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServiceServer).GetSolanaAccountInfo(ctx, req.(*SolanaAccountInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SidecarService_ServiceDesc is the grpc.ServiceDesc for SidecarService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SidecarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.SidecarService",
	HandlerType: (*SidecarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSidecarState",
			Handler:    _SidecarService_GetSidecarState_Handler,
		},
		{
			MethodName: "GetSidecarStateByEthHeight",
			Handler:    _SidecarService_GetSidecarStateByEthHeight_Handler,
		},
		{
			MethodName: "GetBitcoinBlockHeaderByHeight",
			Handler:    _SidecarService_GetBitcoinBlockHeaderByHeight_Handler,
		},
		{
			MethodName: "GetLatestBitcoinBlockHeader",
			Handler:    _SidecarService_GetLatestBitcoinBlockHeader_Handler,
		},
		{
			MethodName: "GetLatestEthereumNonceForAccount",
			Handler:    _SidecarService_GetLatestEthereumNonceForAccount_Handler,
		},
		{
			MethodName: "GetSolanaRecentBlockhash",
			Handler:    _SidecarService_GetSolanaRecentBlockhash_Handler,
		},
		{
			MethodName: "GetSolanaAccountInfo",
			Handler:    _SidecarService_GetSolanaAccountInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/sidecar_service.proto",
}
