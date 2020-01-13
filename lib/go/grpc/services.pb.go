// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services.proto

package grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	// 346 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xdd, 0x4e, 0xf2, 0x40,
	0x10, 0x86, 0xc3, 0x97, 0x4f, 0x7e, 0x06, 0xa3, 0x65, 0xc0, 0x04, 0xb8, 0x88, 0xc5, 0xd4, 0x98,
	0xf8, 0x73, 0x24, 0xc4, 0xd4, 0x98, 0x10, 0x23, 0x62, 0x88, 0x9e, 0xad, 0x65, 0x04, 0xd2, 0xa6,
	0x5b, 0x77, 0x07, 0x8d, 0x87, 0x5e, 0x8f, 0x37, 0x69, 0x6c, 0xbb, 0xb4, 0x20, 0x31, 0xf1, 0x6c,
	0xde, 0x67, 0xfa, 0x4c, 0x77, 0xda, 0x85, 0x3d, 0x43, 0xfa, 0x75, 0xe1, 0x93, 0x11, 0xb1, 0x56,
	0xac, 0xba, 0x75, 0x7e, 0x8f, 0x6d, 0x70, 0x3f, 0x4a, 0xe0, 0xdc, 0xe8, 0x29, 0xe9, 0xbe, 0x52,
	0xc1, 0x50, 0x46, 0x72, 0x46, 0x1a, 0x5d, 0xa8, 0xad, 0x18, 0x36, 0xc4, 0xaa, 0x1e, 0xd1, 0xcb,
	0x92, 0x0c, 0x77, 0xb1, 0x88, 0x4c, 0xac, 0x22, 0x43, 0x78, 0x0a, 0x95, 0xfb, 0x78, 0x2a, 0x99,
	0x0c, 0xb6, 0xf3, 0x76, 0x86, 0x7e, 0x11, 0x0f, 0x4b, 0xae, 0x82, 0x9d, 0xa1, 0x0c, 0x48, 0x63,
	0x0f, 0xaa, 0x1e, 0xf1, 0xed, 0x52, 0x31, 0xa1, 0x23, 0x6c, 0x69, 0xe5, 0x46, 0x81, 0x64, 0x2f,
	0x3d, 0x06, 0x18, 0xcc, 0xc9, 0x0f, 0x52, 0x05, 0x45, 0x1e, 0xac, 0xd4, 0x5c, 0x63, 0xa9, 0xe6,
	0x3e, 0x40, 0x3d, 0x39, 0xc7, 0x1d, 0x4b, 0x5e, 0x1a, 0xbc, 0x86, 0x46, 0x21, 0xa6, 0x47, 0xc6,
	0x8e, 0xf8, 0xc1, 0xec, 0xcc, 0xee, 0xb6, 0x56, 0x36, 0xfa, 0xf3, 0x1f, 0xd4, 0xae, 0x14, 0x4f,
	0x64, 0x18, 0x12, 0xe3, 0x09, 0xd4, 0x07, 0x9a, 0x24, 0x53, 0x22, 0x60, 0x53, 0x14, 0x92, 0x9d,
	0xd6, 0x5a, 0x87, 0xd9, 0x66, 0x67, 0xb0, 0xef, 0x11, 0x8f, 0x55, 0x40, 0x51, 0x5f, 0x86, 0x32,
	0xf2, 0xbf, 0xd7, 0xf3, 0x88, 0xb3, 0x90, 0xaf, 0x57, 0x64, 0x6b, 0xee, 0x25, 0xcf, 0x49, 0xff,
	0xd9, 0x3d, 0x87, 0x5d, 0x8f, 0xf8, 0x22, 0x0c, 0xd5, 0x5b, 0x22, 0xb6, 0x44, 0x31, 0x5a, 0xf5,
	0x60, 0x83, 0x66, 0x72, 0x0f, 0xaa, 0x63, 0x2d, 0x23, 0xf3, 0x4c, 0x1a, 0x1d, 0x61, 0xcb, 0xfc,
	0xff, 0xe5, 0x24, 0xfb, 0x5a, 0x23, 0xa8, 0x4c, 0x24, 0xfb, 0x73, 0xd2, 0xe8, 0x81, 0x93, 0x94,
	0xc9, 0x33, 0xd2, 0xe7, 0x85, 0x8a, 0xb0, 0x2d, 0x36, 0x91, 0x9d, 0xd5, 0xd9, 0xd2, 0x49, 0x67,
	0xf6, 0xcb, 0x8f, 0xff, 0x67, 0x3a, 0xf6, 0x9f, 0xca, 0xc9, 0x05, 0x3f, 0xfa, 0x0a, 0x00, 0x00,
	0xff, 0xff, 0xf6, 0xa9, 0xbd, 0xba, 0xff, 0x02, 0x00, 0x00,
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

// UnimplementedOrderBookManagerServer can be embedded to have forward compatible implementations.
type UnimplementedOrderBookManagerServer struct {
}

func (*UnimplementedOrderBookManagerServer) OrderBook(ctx context.Context, req *OrderBookRequest) (*OrderBookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderBook not implemented")
}
func (*UnimplementedOrderBookManagerServer) Updates(req *OrderBookUpdatesRequest, srv OrderBookManager_UpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method Updates not implemented")
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

// MakerServer is the server API for Maker service.
type MakerServer interface {
	GetQuote(context.Context, *GetQuoteRequest) (*GetQuoteResponse, error)
	CheckQuote(context.Context, *CheckQuoteRequest) (*CheckQuoteResponse, error)
}

// UnimplementedMakerServer can be embedded to have forward compatible implementations.
type UnimplementedMakerServer struct {
}

func (*UnimplementedMakerServer) GetQuote(ctx context.Context, req *GetQuoteRequest) (*GetQuoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuote not implemented")
}
func (*UnimplementedMakerServer) CheckQuote(ctx context.Context, req *CheckQuoteRequest) (*CheckQuoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckQuote not implemented")
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services.proto",
}

// OrderStatusClient is the client API for OrderStatus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrderStatusClient interface {
	OrderStatusUpdate(ctx context.Context, in *OrderStatusUpdateRequest, opts ...grpc.CallOption) (*OrderStatusUpdateResponse, error)
}

type orderStatusClient struct {
	cc *grpc.ClientConn
}

func NewOrderStatusClient(cc *grpc.ClientConn) OrderStatusClient {
	return &orderStatusClient{cc}
}

func (c *orderStatusClient) OrderStatusUpdate(ctx context.Context, in *OrderStatusUpdateRequest, opts ...grpc.CallOption) (*OrderStatusUpdateResponse, error) {
	out := new(OrderStatusUpdateResponse)
	err := c.cc.Invoke(ctx, "/OrderStatus/OrderStatusUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderStatusServer is the server API for OrderStatus service.
type OrderStatusServer interface {
	OrderStatusUpdate(context.Context, *OrderStatusUpdateRequest) (*OrderStatusUpdateResponse, error)
}

// UnimplementedOrderStatusServer can be embedded to have forward compatible implementations.
type UnimplementedOrderStatusServer struct {
}

func (*UnimplementedOrderStatusServer) OrderStatusUpdate(ctx context.Context, req *OrderStatusUpdateRequest) (*OrderStatusUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderStatusUpdate not implemented")
}

func RegisterOrderStatusServer(s *grpc.Server, srv OrderStatusServer) {
	s.RegisterService(&_OrderStatus_serviceDesc, srv)
}

func _OrderStatus_OrderStatusUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderStatusUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderStatusServer).OrderStatusUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderStatus/OrderStatusUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderStatusServer).OrderStatusUpdate(ctx, req.(*OrderStatusUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrderStatus_serviceDesc = grpc.ServiceDesc{
	ServiceName: "OrderStatus",
	HandlerType: (*OrderStatusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OrderStatusUpdate",
			Handler:    _OrderStatus_OrderStatusUpdate_Handler,
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
	GetTokenBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	GetEtherBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
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

func (c *hotWalletClient) GetTokenBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/GetTokenBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) GetEtherBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/GetEtherBalance", in, out, opts...)
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
	GetTokenBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	GetEtherBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	GetAllowance(context.Context, *GetAllowanceRequest) (*GetAllowanceResponse, error)
	Transfer(context.Context, *TransferRequest) (*TransferResponse, error)
}

// UnimplementedHotWalletServer can be embedded to have forward compatible implementations.
type UnimplementedHotWalletServer struct {
}

func (*UnimplementedHotWalletServer) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (*UnimplementedHotWalletServer) GetTokenBalance(ctx context.Context, req *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTokenBalance not implemented")
}
func (*UnimplementedHotWalletServer) GetEtherBalance(ctx context.Context, req *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEtherBalance not implemented")
}
func (*UnimplementedHotWalletServer) GetAllowance(ctx context.Context, req *GetAllowanceRequest) (*GetAllowanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllowance not implemented")
}
func (*UnimplementedHotWalletServer) Transfer(ctx context.Context, req *TransferRequest) (*TransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
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

func _HotWallet_GetTokenBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).GetTokenBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/GetTokenBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).GetTokenBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HotWallet_GetEtherBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).GetEtherBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/GetEtherBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).GetEtherBalance(ctx, req.(*GetBalanceRequest))
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
			MethodName: "GetTokenBalance",
			Handler:    _HotWallet_GetTokenBalance_Handler,
		},
		{
			MethodName: "GetEtherBalance",
			Handler:    _HotWallet_GetEtherBalance_Handler,
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

// UnimplementedWatcherServer can be embedded to have forward compatible implementations.
type UnimplementedWatcherServer struct {
}

func (*UnimplementedWatcherServer) WatchTransaction(ctx context.Context, req *WatchTransactionRequest) (*WatchTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WatchTransaction not implemented")
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
