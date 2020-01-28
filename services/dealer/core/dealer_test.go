package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type makerMock struct {
	t          *testing.T
	getMarkets func(*types.GetMarketsRequest) (*types.GetMarketsResponse, error)
}

func (m *makerMock) GetQuote(context.Context, *types.GetQuoteRequest, ...grpc.CallOption) (*types.GetQuoteResponse, error) {
	panic("not implemented")
}

func (m *makerMock) CheckQuote(context.Context, *types.CheckQuoteRequest, ...grpc.CallOption) (*types.CheckQuoteResponse, error) {
	panic("not implemented")
}

func (m *makerMock) GetMarkets(_ context.Context, in *types.GetMarketsRequest, _ ...grpc.CallOption) (*types.GetMarketsResponse, error) {
	require.NotNil(m.t, m.getMarkets, "getMarkets not defined in mock")
	return m.getMarkets(in)
}

func TestDealerGetMarkets(t *testing.T) {
	mkts := []*types.Market{
		&types.Market{MakerAssetAddress: "1"},
		&types.Market{MakerAssetAddress: "2"},
		&types.Market{MakerAssetAddress: "3"},
		&types.Market{MakerAssetAddress: "4"},
		&types.Market{MakerAssetAddress: "5"},
		&types.Market{MakerAssetAddress: "6"},
		&types.Market{MakerAssetAddress: "7"},
		&types.Market{MakerAssetAddress: "8"},
		&types.Market{MakerAssetAddress: "9"},
	}
	mock := &makerMock{
		t: t,
		getMarkets: func(req *types.GetMarketsRequest) (*types.GetMarketsResponse, error) {
			return &types.GetMarketsResponse{Markets: mkts}, nil

		},
	}

	d := &Dealer{makerClient: mock}
	t.Run("GetAll", func(t *testing.T) {
		resp, _ := d.GetMarkets("mAddr", "tAddr", 0, 0)
		assert.Len(t, resp, len(mkts))
	})

	t.Run("First", func(t *testing.T) {
		resp, _ := d.GetMarkets("mAddr", "tAddr", 0, 1)
		require.Len(t, resp, 1)
		assert.Equal(t, mkts[0], resp[0])
	})

	t.Run("Last", func(t *testing.T) {
		resp, _ := d.GetMarkets("mAddr", "tAddr", 8, 1)
		require.Len(t, resp, 1)
		assert.Equal(t, mkts[8], resp[0])
	})

	t.Run("Middle", func(t *testing.T) {
		// third page, two per page
		resp, _ := d.GetMarkets("mAddr", "tAddr", 3, 2)
		require.Len(t, resp, 2)
		assert.Equal(t, mkts[6], resp[0])
		assert.Equal(t, mkts[7], resp[1])
	})

}
