// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services.proto

package grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("services.proto", fileDescriptor_8e16ccb8c5307b32) }

var fileDescriptor_8e16ccb8c5307b32 = []byte{
	// 254 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xd1, 0x4a, 0xc3, 0x30,
	0x14, 0x86, 0x19, 0xe8, 0xb0, 0x67, 0x20, 0xeb, 0x99, 0x37, 0xf6, 0x21, 0x22, 0x54, 0x04, 0xaf,
	0x37, 0xd0, 0x21, 0x0c, 0xd1, 0x21, 0x82, 0x77, 0xb1, 0x3b, 0x74, 0xa3, 0xa3, 0x89, 0x27, 0xa7,
	0x82, 0x0f, 0xe5, 0x13, 0xf8, 0x72, 0x62, 0xdb, 0xa4, 0x19, 0xee, 0xee, 0x3f, 0x5f, 0xf2, 0x85,
	0xff, 0x10, 0x38, 0x77, 0xc4, 0x9f, 0xbb, 0x82, 0x9c, 0xb2, 0x6c, 0xc4, 0x64, 0x13, 0xf9, 0xb2,
	0x7e, 0xc8, 0xef, 0x60, 0xfa, 0xc8, 0x1b, 0xe2, 0xb9, 0x31, 0xd5, 0x4a, 0xd7, 0xba, 0x24, 0xc6,
	0x1c, 0x92, 0xc0, 0x30, 0x55, 0x21, 0x3f, 0xd3, 0x47, 0x43, 0x4e, 0x32, 0x8c, 0x91, 0xb3, 0xa6,
	0x76, 0x94, 0xff, 0x8c, 0xe0, 0x74, 0xa5, 0x2b, 0x62, 0xbc, 0x82, 0xb3, 0x7b, 0x92, 0xa7, 0xc6,
	0x08, 0xe1, 0x54, 0xf9, 0xe8, 0xdd, 0x34, 0x22, 0x9d, 0x8a, 0x37, 0x00, 0x8b, 0x2d, 0x15, 0x55,
	0xa7, 0xa0, 0x1a, 0x06, 0x2f, 0xcd, 0x0e, 0x58, 0xaf, 0x3d, 0x40, 0xda, 0xd6, 0x58, 0x8b, 0x96,
	0xc6, 0xbd, 0xd8, 0x8d, 0x16, 0xc2, 0x4b, 0xf5, 0x8f, 0xf9, 0x47, 0xb2, 0x63, 0x47, 0x7d, 0xfb,
	0xef, 0x11, 0x24, 0x4b, 0x23, 0xaf, 0x7a, 0xbf, 0x27, 0xf9, 0xdb, 0x7f, 0xbd, 0x2b, 0xeb, 0xf6,
	0x3a, 0xa6, 0x2a, 0xe4, 0x61, 0xff, 0x08, 0xf5, 0x6d, 0x72, 0x48, 0x96, 0xda, 0x6d, 0xbd, 0x13,
	0xf2, 0xe0, 0x44, 0xa8, 0x77, 0x6e, 0x61, 0xb2, 0x60, 0xd2, 0x42, 0x9d, 0x35, 0x53, 0xd1, 0xe4,
	0xbd, 0x8b, 0x43, 0xd8, 0x99, 0xf3, 0xf1, 0xdb, 0x49, 0xc9, 0xb6, 0x78, 0x1f, 0xb7, 0x9f, 0x78,
	0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xef, 0xbe, 0x73, 0x65, 0xe3, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrderBookManagerClient is the client API for OrderBookManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrderBookManagerClient interface {
	OrderBook(ctx context.Context, in *OrderBookRequest, opts ...grpc.CallOption) (*OrderBookResponse, error)
}

type orderBookManagerClient struct {
	cc *grpc.ClientConn
}

func NewOrderBookManagerClient(cc *grpc.ClientConn) OrderBookManagerClient {
	return &orderBookManagerClient{cc}
}

func (c *orderBookManagerClient) OrderBook(ctx context.Context, in *OrderBookRequest, opts ...grpc.CallOption) (*OrderBookResponse, error) {
	out := new(OrderBookResponse)
	err := c.cc.Invoke(ctx, "/OrderBookManager/OrderBook", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderBookManagerServer is the server API for OrderBookManager service.
type OrderBookManagerServer interface {
	OrderBook(context.Context, *OrderBookRequest) (*OrderBookResponse, error)
}

func RegisterOrderBookManagerServer(s *grpc.Server, srv OrderBookManagerServer) {
	s.RegisterService(&_OrderBookManager_serviceDesc, srv)
}

func _OrderBookManager_OrderBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderBookManagerServer).OrderBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderBookManager/OrderBook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderBookManagerServer).OrderBook(ctx, req.(*OrderBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrderBookManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "OrderBookManager",
	HandlerType: (*OrderBookManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OrderBook",
			Handler:    _OrderBookManager_OrderBook_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}

// MakerClient is the client API for Maker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MakerClient interface {
	GetQuote(ctx context.Context, in *GetQuoteRequest, opts ...grpc.CallOption) (*GetQuoteResponse, error)
	CheckQuote(ctx context.Context, in *CheckQuoteRequest, opts ...grpc.CallOption) (*CheckQuoteResponse, error)
	OrderStatusUpdate(ctx context.Context, in *OrderStatusUpdateRequest, opts ...grpc.CallOption) (*OrderStatusUpdateResponse, error)
}

type makerClient struct {
	cc *grpc.ClientConn
}

func NewMakerClient(cc *grpc.ClientConn) MakerClient {
	return &makerClient{cc}
}

func (c *makerClient) GetQuote(ctx context.Context, in *GetQuoteRequest, opts ...grpc.CallOption) (*GetQuoteResponse, error) {
	out := new(GetQuoteResponse)
	err := c.cc.Invoke(ctx, "/Maker/GetQuote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *makerClient) CheckQuote(ctx context.Context, in *CheckQuoteRequest, opts ...grpc.CallOption) (*CheckQuoteResponse, error) {
	out := new(CheckQuoteResponse)
	err := c.cc.Invoke(ctx, "/Maker/CheckQuote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *makerClient) OrderStatusUpdate(ctx context.Context, in *OrderStatusUpdateRequest, opts ...grpc.CallOption) (*OrderStatusUpdateResponse, error) {
	out := new(OrderStatusUpdateResponse)
	err := c.cc.Invoke(ctx, "/Maker/OrderStatusUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MakerServer is the server API for Maker service.
type MakerServer interface {
	GetQuote(context.Context, *GetQuoteRequest) (*GetQuoteResponse, error)
	CheckQuote(context.Context, *CheckQuoteRequest) (*CheckQuoteResponse, error)
	OrderStatusUpdate(context.Context, *OrderStatusUpdateRequest) (*OrderStatusUpdateResponse, error)
}

func RegisterMakerServer(s *grpc.Server, srv MakerServer) {
	s.RegisterService(&_Maker_serviceDesc, srv)
}

func _Maker_GetQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetQuoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MakerServer).GetQuote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Maker/GetQuote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MakerServer).GetQuote(ctx, req.(*GetQuoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Maker_CheckQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckQuoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MakerServer).CheckQuote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Maker/CheckQuote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MakerServer).CheckQuote(ctx, req.(*CheckQuoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Maker_OrderStatusUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderStatusUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MakerServer).OrderStatusUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Maker/OrderStatusUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MakerServer).OrderStatusUpdate(ctx, req.(*OrderStatusUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Maker_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Maker",
	HandlerType: (*MakerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetQuote",
			Handler:    _Maker_GetQuote_Handler,
		},
		{
			MethodName: "CheckQuote",
			Handler:    _Maker_CheckQuote_Handler,
		},
		{
			MethodName: "OrderStatusUpdate",
			Handler:    _Maker_OrderStatusUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}

// HotWalletClient is the client API for HotWallet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HotWalletClient interface {
	SignOrder(ctx context.Context, in *SignOrderRequest, opts ...grpc.CallOption) (*SignOrderResponse, error)
	HashOrder(ctx context.Context, in *HashOrderRequest, opts ...grpc.CallOption) (*HashOrderResponse, error)
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
}

type hotWalletClient struct {
	cc *grpc.ClientConn
}

func NewHotWalletClient(cc *grpc.ClientConn) HotWalletClient {
	return &hotWalletClient{cc}
}

func (c *hotWalletClient) SignOrder(ctx context.Context, in *SignOrderRequest, opts ...grpc.CallOption) (*SignOrderResponse, error) {
	out := new(SignOrderResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/SignOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) HashOrder(ctx context.Context, in *HashOrderRequest, opts ...grpc.CallOption) (*HashOrderResponse, error) {
	out := new(HashOrderResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/HashOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HotWalletServer is the server API for HotWallet service.
type HotWalletServer interface {
	SignOrder(context.Context, *SignOrderRequest) (*SignOrderResponse, error)
	HashOrder(context.Context, *HashOrderRequest) (*HashOrderResponse, error)
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
}

func RegisterHotWalletServer(s *grpc.Server, srv HotWalletServer) {
	s.RegisterService(&_HotWallet_serviceDesc, srv)
}

func _HotWallet_SignOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).SignOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/SignOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).SignOrder(ctx, req.(*SignOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HotWallet_HashOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HashOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).HashOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/HashOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).HashOrder(ctx, req.(*HashOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HotWallet_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HotWallet_serviceDesc = grpc.ServiceDesc{
	ServiceName: "HotWallet",
	HandlerType: (*HotWalletServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignOrder",
			Handler:    _HotWallet_SignOrder_Handler,
		},
		{
			MethodName: "HashOrder",
			Handler:    _HotWallet_HashOrder_Handler,
		},
		{
			MethodName: "CreateOrder",
			Handler:    _HotWallet_CreateOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}
