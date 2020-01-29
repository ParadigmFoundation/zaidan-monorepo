package core

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type MakerMock struct {
	GetMarketsFn func(*types.GetMarketsRequest) (*types.GetMarketsResponse, error)
}

func (m *MakerMock) GetQuote(context.Context, *types.GetQuoteRequest, ...grpc.CallOption) (*types.GetQuoteResponse, error) {
	panic("not implemented")
}

func (m *MakerMock) CheckQuote(context.Context, *types.CheckQuoteRequest, ...grpc.CallOption) (*types.CheckQuoteResponse, error) {
	panic("not implemented")
}

func (m *MakerMock) GetMarkets(_ context.Context, in *types.GetMarketsRequest, _ ...grpc.CallOption) (*types.GetMarketsResponse, error) {
	if fn := m.GetMarketsFn; fn != nil {
		return fn(in)
	}
	return &types.GetMarketsResponse{}, nil
}

type HWMock struct {
	GetTradeInfoFn func() (*types.TradeInfo, error)
}

func (d *HWMock) CreateOrder(context.Context, *types.CreateOrderRequest, ...grpc.CallOption) (*types.CreateOrderResponse, error) {
	panic("not implemented")
}

func (d *HWMock) ValidateOrder(context.Context, *types.ValidateOrderRequest, ...grpc.CallOption) (*types.ValidateOrderResponse, error) {
	panic("not implemented")
}

func (d *HWMock) GetAllowance(context.Context, *types.GetAllowanceRequest, ...grpc.CallOption) (*types.GetAllowanceResponse, error) {
	panic("not implemented")
}

func (d *HWMock) SetAllowance(context.Context, *types.SetAllowanceRequest, ...grpc.CallOption) (*types.SetAllowanceResponse, error) {
	panic("not implemented")
}

func (d *HWMock) GetTokenBalance(context.Context, *types.GetBalanceRequest, ...grpc.CallOption) (*types.GetBalanceResponse, error) {
	panic("not implemented")
}

func (d *HWMock) GetEtherBalance(context.Context, *types.GetBalanceRequest, ...grpc.CallOption) (*types.GetBalanceResponse, error) {
	panic("not implemented")
}

func (d *HWMock) TransferEther(context.Context, *types.TransferRequest, ...grpc.CallOption) (*types.TransferResponse, error) {
	panic("not implemented")
}

func (d *HWMock) TransferToken(context.Context, *types.TransferRequest, ...grpc.CallOption) (*types.TransferResponse, error) {
	panic("not implemented")
}

func (d *HWMock) SendTransaction(context.Context, *types.SendTransactionRequest, ...grpc.CallOption) (*types.SendTransactionResponse, error) {
	panic("not implemented")
}

func (d *HWMock) ExecuteZeroExTransaction(context.Context, *types.ExecuteZeroExTransactionRequest, ...grpc.CallOption) (*types.ExecuteZeroExTransactionResponse, error) {
	panic("not implemented")
}

func (d *HWMock) GetTradeInfo(context.Context, *empty.Empty, ...grpc.CallOption) (*types.TradeInfo, error) {
	if fn := d.GetTradeInfoFn; fn != nil {
		return fn()
	}
	return &types.TradeInfo{}, nil
}
