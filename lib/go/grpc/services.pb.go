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
	// 379 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xdf, 0x6e, 0xda, 0x30,
	0x14, 0xc6, 0x85, 0xb6, 0xc1, 0x38, 0x6c, 0x23, 0x18, 0xa6, 0x01, 0x0f, 0xe1, 0x4d, 0xd9, 0x26,
	0xed, 0xcf, 0xd5, 0x60, 0x55, 0xaa, 0x4a, 0xa8, 0x2a, 0xa1, 0x42, 0xea, 0x9d, 0x1b, 0x4e, 0x01,
	0x25, 0x8a, 0x53, 0xfb, 0xd0, 0xaa, 0x97, 0x7d, 0xae, 0x3e, 0x51, 0xdf, 0xa2, 0xc2, 0x89, 0x49,
	0x02, 0x51, 0xd5, 0xde, 0xf9, 0xfb, 0x1d, 0xfd, 0x3e, 0x1f, 0xc5, 0x0a, 0x7c, 0xd2, 0xa8, 0x6e,
	0xd6, 0x01, 0x6a, 0x9e, 0x28, 0x49, 0x72, 0xd8, 0xa2, 0xbb, 0xc4, 0x06, 0xf7, 0xbe, 0x06, 0xce,
	0xa9, 0x5a, 0xa0, 0x1a, 0x49, 0x19, 0x4e, 0x44, 0x2c, 0x96, 0xa8, 0x98, 0x0b, 0xcd, 0x1d, 0x63,
	0x1d, 0xbe, 0x3b, 0x4f, 0xf1, 0x7a, 0x83, 0x9a, 0x86, 0xac, 0x88, 0x74, 0x22, 0x63, 0x8d, 0xec,
	0x37, 0x34, 0xce, 0x93, 0x85, 0x20, 0xd4, 0xac, 0x9f, 0x8f, 0x33, 0xf4, 0x8c, 0xf8, 0xad, 0xe6,
	0x3e, 0xd4, 0xe0, 0xdd, 0x44, 0x84, 0xa8, 0xd8, 0x57, 0x78, 0xef, 0x21, 0x9d, 0x6d, 0x24, 0x21,
	0x73, 0xb8, 0x3d, 0x5a, 0xbb, 0x53, 0x20, 0xd9, 0xad, 0x3f, 0x01, 0xc6, 0x2b, 0x0c, 0xc2, 0x54,
	0x61, 0x3c, 0x0f, 0x56, 0xea, 0x96, 0x58, 0xa6, 0x9d, 0x40, 0xc7, 0x2c, 0xe2, 0x93, 0xa0, 0x8d,
	0x4e, 0x97, 0x64, 0x03, 0x7e, 0xc0, 0x6c, 0xc9, 0xb0, 0x6a, 0x94, 0x76, 0xb9, 0x8f, 0x6f, 0xa0,
	0x79, 0x2c, 0x69, 0x2e, 0xa2, 0x08, 0x89, 0xfd, 0x82, 0xd6, 0x58, 0xa1, 0x20, 0x34, 0x02, 0xeb,
	0xf2, 0x42, 0xb2, 0x6d, 0xbd, 0x32, 0xcc, 0x76, 0xfa, 0x0b, 0x1f, 0x3c, 0xa4, 0x7f, 0x51, 0x24,
	0x6f, 0x45, 0x1c, 0x20, 0xeb, 0xf1, 0x62, 0xb4, 0xee, 0xe7, 0x3d, 0x9a, 0xcb, 0x7e, 0x59, 0xf6,
	0x2b, 0x65, 0xbf, 0x4a, 0xfe, 0x03, 0x6d, 0x0f, 0x69, 0x26, 0x43, 0x8c, 0x47, 0x22, 0x32, 0x3e,
	0xdb, 0x5e, 0x93, 0x85, 0xfc, 0x4b, 0x16, 0x59, 0xc9, 0x3d, 0xa2, 0x15, 0xaa, 0x57, 0xbb, 0x3f,
	0xe0, 0xe3, 0x4c, 0x89, 0x58, 0x5f, 0xa1, 0x32, 0x05, 0xcc, 0xe1, 0x36, 0xe7, 0x4f, 0x9e, 0x93,
	0x43, 0xcb, 0xac, 0xfc, 0x32, 0xeb, 0x3f, 0xb4, 0x7d, 0x8c, 0x17, 0x86, 0x8b, 0x80, 0xd6, 0x32,
	0x66, 0x5f, 0xf8, 0x1e, 0xb1, 0x7a, 0xff, 0x70, 0x90, 0xbd, 0xf5, 0x14, 0x1a, 0x73, 0x41, 0xc1,
	0x76, 0x57, 0x0f, 0x1c, 0x73, 0x2c, 0x36, 0xf6, 0xf9, 0x3e, 0xb2, 0x95, 0x83, 0x8a, 0x49, 0xda,
	0x39, 0xaa, 0x5f, 0xbc, 0x5d, 0xaa, 0x24, 0xb8, 0xac, 0x9b, 0x1f, 0xf2, 0xfb, 0x53, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x78, 0xf8, 0xaa, 0x1c, 0xaf, 0x03, 0x00, 0x00,
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
	GetAllowance(ctx context.Context, in *GetAllowanceRequest, opts ...grpc.CallOption) (*GetAllowanceResponse, error)
	SetAllowance(ctx context.Context, in *SetAllowanceRequest, opts ...grpc.CallOption) (*SetAllowanceResponse, error)
	GetTokenBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	GetEtherBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	TransferEther(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error)
	TransferToken(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error)
	SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*SendTransactionResponse, error)
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

func (c *hotWalletClient) GetAllowance(ctx context.Context, in *GetAllowanceRequest, opts ...grpc.CallOption) (*GetAllowanceResponse, error) {
	out := new(GetAllowanceResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/GetAllowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) SetAllowance(ctx context.Context, in *SetAllowanceRequest, opts ...grpc.CallOption) (*SetAllowanceResponse, error) {
	out := new(SetAllowanceResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/SetAllowance", in, out, opts...)
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

func (c *hotWalletClient) TransferEther(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error) {
	out := new(TransferResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/TransferEther", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) TransferToken(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error) {
	out := new(TransferResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/TransferToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotWalletClient) SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*SendTransactionResponse, error) {
	out := new(SendTransactionResponse)
	err := c.cc.Invoke(ctx, "/HotWallet/SendTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HotWalletServer is the server API for HotWallet service.
type HotWalletServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetAllowance(context.Context, *GetAllowanceRequest) (*GetAllowanceResponse, error)
	SetAllowance(context.Context, *SetAllowanceRequest) (*SetAllowanceResponse, error)
	GetTokenBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	GetEtherBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	TransferEther(context.Context, *TransferRequest) (*TransferResponse, error)
	TransferToken(context.Context, *TransferRequest) (*TransferResponse, error)
	SendTransaction(context.Context, *SendTransactionRequest) (*SendTransactionResponse, error)
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

func _HotWallet_SetAllowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAllowanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).SetAllowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/SetAllowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).SetAllowance(ctx, req.(*SetAllowanceRequest))
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

func _HotWallet_TransferEther_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).TransferEther(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/TransferEther",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).TransferEther(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HotWallet_TransferToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).TransferToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/TransferToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).TransferToken(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HotWallet_SendTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotWalletServer).SendTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HotWallet/SendTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotWalletServer).SendTransaction(ctx, req.(*SendTransactionRequest))
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
			MethodName: "GetAllowance",
			Handler:    _HotWallet_GetAllowance_Handler,
		},
		{
			MethodName: "SetAllowance",
			Handler:    _HotWallet_SetAllowance_Handler,
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
			MethodName: "TransferEther",
			Handler:    _HotWallet_TransferEther_Handler,
		},
		{
			MethodName: "TransferToken",
			Handler:    _HotWallet_TransferToken_Handler,
		},
		{
			MethodName: "SendTransaction",
			Handler:    _HotWallet_SendTransaction_Handler,
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
