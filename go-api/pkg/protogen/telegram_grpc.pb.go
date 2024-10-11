// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.26.1
// source: proto/telegram.proto

package protogen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Telega_GetPortfolioList_FullMethodName           = "/Telega/GetPortfolioList"
	Telega_GetPortfolioSummaryMessage_FullMethodName = "/Telega/GetPortfolioSummaryMessage"
)

// TelegaClient is the client API for Telega service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TelegaClient interface {
	GetPortfolioList(ctx context.Context, in *PortfolioListRequest, opts ...grpc.CallOption) (*PortfolioListResponse, error)
	GetPortfolioSummaryMessage(ctx context.Context, in *PortfolioRequest, opts ...grpc.CallOption) (*PortfolioSummaryResponse, error)
}

type telegaClient struct {
	cc grpc.ClientConnInterface
}

func NewTelegaClient(cc grpc.ClientConnInterface) TelegaClient {
	return &telegaClient{cc}
}

func (c *telegaClient) GetPortfolioList(ctx context.Context, in *PortfolioListRequest, opts ...grpc.CallOption) (*PortfolioListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PortfolioListResponse)
	err := c.cc.Invoke(ctx, Telega_GetPortfolioList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telegaClient) GetPortfolioSummaryMessage(ctx context.Context, in *PortfolioRequest, opts ...grpc.CallOption) (*PortfolioSummaryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PortfolioSummaryResponse)
	err := c.cc.Invoke(ctx, Telega_GetPortfolioSummaryMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TelegaServer is the server API for Telega service.
// All implementations must embed UnimplementedTelegaServer
// for forward compatibility
type TelegaServer interface {
	GetPortfolioList(context.Context, *PortfolioListRequest) (*PortfolioListResponse, error)
	GetPortfolioSummaryMessage(context.Context, *PortfolioRequest) (*PortfolioSummaryResponse, error)
	mustEmbedUnimplementedTelegaServer()
}

// UnimplementedTelegaServer must be embedded to have forward compatible implementations.
type UnimplementedTelegaServer struct {
}

func (UnimplementedTelegaServer) GetPortfolioList(context.Context, *PortfolioListRequest) (*PortfolioListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortfolioList not implemented")
}
func (UnimplementedTelegaServer) GetPortfolioSummaryMessage(context.Context, *PortfolioRequest) (*PortfolioSummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortfolioSummaryMessage not implemented")
}
func (UnimplementedTelegaServer) mustEmbedUnimplementedTelegaServer() {}

// UnsafeTelegaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TelegaServer will
// result in compilation errors.
type UnsafeTelegaServer interface {
	mustEmbedUnimplementedTelegaServer()
}

func RegisterTelegaServer(s grpc.ServiceRegistrar, srv TelegaServer) {
	s.RegisterService(&Telega_ServiceDesc, srv)
}

func _Telega_GetPortfolioList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PortfolioListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelegaServer).GetPortfolioList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Telega_GetPortfolioList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelegaServer).GetPortfolioList(ctx, req.(*PortfolioListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Telega_GetPortfolioSummaryMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PortfolioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelegaServer).GetPortfolioSummaryMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Telega_GetPortfolioSummaryMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelegaServer).GetPortfolioSummaryMessage(ctx, req.(*PortfolioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Telega_ServiceDesc is the grpc.ServiceDesc for Telega service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Telega_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Telega",
	HandlerType: (*TelegaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPortfolioList",
			Handler:    _Telega_GetPortfolioList_Handler,
		},
		{
			MethodName: "GetPortfolioSummaryMessage",
			Handler:    _Telega_GetPortfolioSummaryMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/telegram.proto",
}