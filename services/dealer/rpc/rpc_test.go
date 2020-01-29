package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/utils/ptr"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
)

func TestRPC(t *testing.T) {
	store, err := sql.New("sqlite3", ":memory:")
	require.NoError(t, err)

	dealer, err := core.NewDealer(context.Background(), store, core.DealerConfig{})
	require.NoError(t, err)

	t.Run("AuthStatus", func(t *testing.T) {
		service, err := NewService(dealer)
		require.NoError(t, err)

		require.NoError(t, store.CreatePolicy("xxx"))

		defer service.WithPolicy(0, nil)

		t.Run("Blacklist", func(t *testing.T) {
			service.WithPolicy(PolicyBlackList, store)
			resp, err := service.AuthStatus("xxx")
			require.NoError(t, err)
			assert.False(t, resp.Authorized)
			assert.Equal(t, "BLACKLISTED", resp.Reason)
		})

		t.Run("Whitelist", func(t *testing.T) {
			service.WithPolicy(PolicyWhiteList, store)
			resp, err := service.AuthStatus("xxx")
			require.NoError(t, err)
			assert.True(t, resp.Authorized)
			assert.Equal(t, "WHITELISTED", resp.Reason)
		})
	})

	t.Run("GetMarkets", func(t *testing.T) {
		mkts := make([]*types.Market, 99)
		m := &core.MakerMock{
			T: t,
			GetMarketsFn: func(req *types.GetMarketsRequest) (*types.GetMarketsResponse, error) {
				return &types.GetMarketsResponse{Markets: mkts}, nil
			},
		}
		service, err := NewService(dealer.WithMakerClient(m))
		require.NoError(t, err)

		t.Run("Defaults", func(t *testing.T) {
			resp, err := service.GetMarkets(nil, nil, nil, nil)
			require.NoError(t, err)

			assert.Len(t, resp.record, DEFAULT_PER_PAGE)
			assert.Equal(t, resp.total, DEFAULT_PER_PAGE)
			assert.Equal(t, resp.page, 0)
		})

		t.Run("1PerPage", func(t *testing.T) {
			resp, err := service.GetMarkets(nil, nil, ptr.FromInt(1), ptr.FromInt(1))
			require.NoError(t, err)

			assert.Len(t, resp.record, 1)
			assert.Equal(t, resp.total, 1)
			assert.Equal(t, resp.page, 1)
		})
	})
}
