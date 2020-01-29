package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type MakerMock struct {
	T            *testing.T
	GetMarketsFn func(*types.GetMarketsRequest) (*types.GetMarketsResponse, error)
}

func (m *MakerMock) GetQuote(context.Context, *types.GetQuoteRequest, ...grpc.CallOption) (*types.GetQuoteResponse, error) {
	panic("not implemented")
}

func (m *MakerMock) CheckQuote(context.Context, *types.CheckQuoteRequest, ...grpc.CallOption) (*types.CheckQuoteResponse, error) {
	panic("not implemented")
}

func (m *MakerMock) GetMarkets(_ context.Context, in *types.GetMarketsRequest, _ ...grpc.CallOption) (*types.GetMarketsResponse, error) {
	fn := m.GetMarketsFn
	require.NotNil(m.T, fn, "getMarkets not defined in mock")
	return fn(in)
}
