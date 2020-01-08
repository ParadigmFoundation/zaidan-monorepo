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
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xdd, 0x4a, 0xf3, 0x40,
	0x10, 0x86, 0x09, 0x7c, 0x5f, 0x7f, 0xa6, 0x22, 0xe9, 0xb4, 0x42, 0xdb, 0x8b, 0xd8, 0x4a, 0x44,
	0x50, 0x3c, 0xb2, 0x3d, 0x88, 0x08, 0x45, 0xac, 0x4a, 0xc1, 0xb3, 0x35, 0x1d, 0xdb, 0xd2, 0x90,
	0x8d, 0xbb, 0x13, 0xc5, 0x43, 0xaf, 0xcb, 0x3b, 0xf2, 0x2a, 0xc4, 0x24, 0xdb, 0x4d, 0x6b, 0xf1,
	0x6c, 0xde, 0x67, 0xf3, 0xbc, 0xd9, 0x09, 0x81, 0x43, 0x43, 0xfa, 0x75, 0x15, 0x91, 0x11, 0xa9,
	0x56, 0xac, 0x06, 0x2d, 0x7e, 0x4f, 0x6d, 0x08, 0x3e, 0x3c, 0xf0, 0x6f, 0xf4, 0x9c, 0xf4, 0x48,
	0xa9, 0xf5, 0x44, 0x26, 0x72, 0x41, 0x1a, 0x03, 0x68, 0x6e, 0x18, 0xb6, 0xc5, 0x66, 0x9e, 0xd2,
	0x4b, 0x46, 0x86, 0x07, 0x58, 0x45, 0x26, 0x55, 0x89, 0x21, 0x3c, 0x87, 0xfa, 0x43, 0x3a, 0x97,
	0x4c, 0x06, 0x7b, 0xee, 0xb8, 0x44, 0x7f, 0x88, 0xc7, 0x5e, 0xf0, 0xe9, 0xc1, 0xff, 0x89, 0x5c,
	0x93, 0xc6, 0x21, 0x34, 0x42, 0xe2, 0xdb, 0x4c, 0x31, 0xa1, 0x2f, 0xec, 0x68, 0xed, 0x76, 0x85,
	0x94, 0x6f, 0x3d, 0x05, 0x18, 0x2f, 0x29, 0x5a, 0x17, 0x0a, 0x0a, 0x17, 0xac, 0xd4, 0xd9, 0x62,
	0xa5, 0x76, 0x0d, 0xed, 0xfc, 0x22, 0x77, 0x2c, 0x39, 0x33, 0xc5, 0x25, 0xb1, 0x2f, 0x7e, 0x31,
	0x5b, 0x32, 0xd8, 0x77, 0x54, 0x74, 0x05, 0x5f, 0x1e, 0x34, 0xaf, 0x14, 0xcf, 0x64, 0x1c, 0x13,
	0xe3, 0x19, 0xb4, 0xc6, 0x9a, 0x24, 0x53, 0x2e, 0x60, 0x47, 0x54, 0x92, 0x6d, 0xeb, 0x6e, 0x43,
	0xb7, 0x4a, 0x48, 0x3c, 0x92, 0xb1, 0x4c, 0xa2, 0x9f, 0x55, 0x5c, 0x70, 0xab, 0x54, 0x59, 0xa9,
	0x5d, 0xc0, 0x41, 0x48, 0x7c, 0x19, 0xc7, 0xea, 0x2d, 0x17, 0xbb, 0xa2, 0x1a, 0xad, 0x7a, 0xb4,
	0x43, 0x4b, 0x79, 0x08, 0x8d, 0x7b, 0x2d, 0x13, 0xf3, 0x4c, 0x1a, 0x7d, 0x61, 0x47, 0xf7, 0xbd,
	0x1d, 0x29, 0x97, 0x9d, 0x42, 0x7d, 0x26, 0x39, 0x5a, 0x92, 0xc6, 0x10, 0xfc, 0x7c, 0xcc, 0x9f,
	0x91, 0x11, 0xaf, 0x54, 0x82, 0x3d, 0xb1, 0x8b, 0x6c, 0x57, 0x7f, 0xcf, 0x49, 0xd1, 0x39, 0xaa,
	0x3d, 0xfe, 0x5b, 0xe8, 0x34, 0x7a, 0xaa, 0xe5, 0x7f, 0xe4, 0xc9, 0x77, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x1d, 0xde, 0x17, 0xd1, 0xb0, 0x02, 0x00, 0x00,
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
	Updates(ctx context.Context, in *OrderBookUpdatesRequest, opts ...grpc.CallOption) (OrderBookManager_UpdatesClient, error)
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

func (c *orderBookManagerClient) Updates(ctx context.Context, in *OrderBookUpdatesRequest, opts ...grpc.CallOption) (OrderBookManager_UpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_OrderBookManager_serviceDesc.Streams[0], "/OrderBookManager/Updates", opts...)
	if err != nil {
		return nil, err
	}
	x := &orderBookManagerUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OrderBookManager_UpdatesClient interface {
	Recv() (*OrderBookResponse, error)
	grpc.ClientStream
}

type orderBookManagerUpdatesClient struct {
	grpc.ClientStream
}

func (x *orderBookManagerUpdatesClient) Recv() (*OrderBookResponse, error) {
	m := new(OrderBookResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OrderBookManagerServer is the server API for OrderBookManager service.
type OrderBookManagerServer interface {
	OrderBook(context.Context, *OrderBookRequest) (*OrderBookResponse, error)
	Updates(*OrderBookUpdatesRequest, OrderBookManager_UpdatesServer) error
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

func _OrderBookManager_Updates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OrderBookUpdatesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OrderBookManagerServer).Updates(m, &orderBookManagerUpdatesServer{stream})
}

type OrderBookManager_UpdatesServer interface {
	Send(*OrderBookResponse) error
	grpc.ServerStream
}

type orderBookManagerUpdatesServer struct {
	grpc.ServerStream
}

func (x *orderBookManagerUpdatesServer) Send(m *OrderBookResponse) error {
	return x.ServerStream.SendMsg(m)
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
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Updates",
			Handler:       _OrderBookManager_Updates_Handler,
			ServerStreams: true,
		},
	},
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
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	GetAllowance(ctx context.Context, in *GetAllowanceRequest, opts ...grpc.CallOption) (*GetAllowanceResponse, error)
	Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error)
}

type hotWalletClient struct {
	cc *grpc.ClientConn
}

func NewHotWalletClient(cc *grpc.ClientConn) HotWalletClient {
	return &hotWalletClient{cc}
}

func (c *hotWalletClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/GetBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) GetAllowance(ctx context.Context, in *GetAllowanceRequest, opts ...grpc.CallOption) (*GetAllowanceResponse, error) {
	out := new(GetAllowanceResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/GetAllowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error) {
	out := new(TransferResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/Transfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HotWalletServer is the server API for HotWallet service.
type HotWalletServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	GetAllowance(context.Context, *GetAllowanceRequest) (*GetAllowanceResponse, error)
	Transfer(context.Context, *TransferRequest) (*TransferResponse, error)
}

func RegisterHotWalletServer(s *grpc.Server, srv HotWalletServer) {
	s.RegisterService(&_HotWallet_serviceDesc, srv)
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

func _HotWallet_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/GetBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).GetBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HotWallet_GetAllowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllowanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).GetAllowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/GetAllowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).GetAllowance(ctx, req.(*GetAllowanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HotWallet_Transfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/Transfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).Transfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HotWallet_serviceDesc = grpc.ServiceDesc{
	ServiceName: "HotWallet",
	HandlerType: (*HotWalletServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _HotWallet_CreateOrder_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _HotWallet_GetBalance_Handler,
		},
		{
			MethodName: "GetAllowance",
			Handler:    _HotWallet_GetAllowance_Handler,
		},
		{
			MethodName: "Transfer",
			Handler:    _HotWallet_Transfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}

// WatcherClient is the client API for Watcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WatcherClient interface {
	WatchTransaction(ctx context.Context, in *WatchTransactionRequest, opts ...grpc.CallOption) (*WatchTransactionResponse, error)
}

type watcherClient struct {
	cc *grpc.ClientConn
}

func NewWatcherClient(cc *grpc.ClientConn) WatcherClient {
	return &watcherClient{cc}
}

func (c *watcherClient) WatchTransaction(ctx context.Context, in *WatchTransactionRequest, opts ...grpc.CallOption) (*WatchTransactionResponse, error) {
	out := new(WatchTransactionResponse)
	err := c.cc.Invoke(ctx, "/Watcher/WatchTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WatcherServer is the server API for Watcher service.
type WatcherServer interface {
	WatchTransaction(context.Context, *WatchTransactionRequest) (*WatchTransactionResponse, error)
}

func RegisterWatcherServer(s *grpc.Server, srv WatcherServer) {
	s.RegisterService(&_Watcher_serviceDesc, srv)
}

func _Watcher_WatchTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WatchTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WatcherServer).WatchTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Watcher/WatchTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WatcherServer).WatchTransaction(ctx, req.(*WatchTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Watcher_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Watcher",
	HandlerType: (*WatcherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WatchTransaction",
			Handler:    _Watcher_WatchTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}
