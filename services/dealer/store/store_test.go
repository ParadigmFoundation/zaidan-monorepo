package store

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paradigmfoundation/zaidan-monorepo/services/dealer"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store, err := New("sqlite3", ":memory:")
	require.NoError(t, err)

	t.Run("Trade", Trade(store.Debug()))
}

func Trade(s *Store) func(t *testing.T) {
	fn := func(t *testing.T) {
		trade := &dealer.Trade{}
		s.CreateTrade(trade)
	}
	return fn
}
