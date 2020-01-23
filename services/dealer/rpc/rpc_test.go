package rpc

import (
	"context"
	"testing"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPC(t *testing.T) {
	store, err := sql.New("sqlite3", ":memory:")
	require.NoError(t, err)

	dealer, err := core.NewDealer(context.Background(), store, core.DealerConfig{})
	require.NoError(t, err)

	service, err := NewService(dealer)
	require.NoError(t, err)

	t.Run("AuthStatus", func(t *testing.T) {
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
}
