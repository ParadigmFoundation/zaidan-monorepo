package storetest

import (
	"testing"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Factory func(t *testing.T) (store.Store, func())

func TestSuite(t *testing.T, f Factory) {
	cases := []struct {
		name string
		run  func(t *testing.T, store store.Store)
	}{
		{"ValidMarket", ValidMarket},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			s, closer := f(t)
			defer closer()

			test.run(t, s)
		})
	}
}

func ValidMarket(t *testing.T, store store.Store) {
	sym := "BTC-USD"
	store.OnUpdate("test", &obm.Update{
		Symbol: sym,
		Bids: obm.Entries{
			{Price: 2, Quantity: 1},
			{Price: 3, Quantity: 1},
			{Price: 1, Quantity: 1},
		},
		Asks: obm.Entries{
			{Price: 2, Quantity: 1},
			{Price: 3, Quantity: 1},
			{Price: 1, Quantity: 1},
		},
	})

	ob, err := store.OrderBook("test", sym)
	require.NoError(t, err)
	require.NotNil(t, ob)

	assert.Equal(t, "test", ob.Exchange)
	assert.Equal(t, sym, ob.Symbol)

	t.Run("Bids", func(t *testing.T) {
		bids := ob.Bids
		require.Len(t, bids, 3)
		assert.Equal(t, &obm.Entry{Price: 3, Quantity: 1}, bids[0])
		assert.Equal(t, &obm.Entry{Price: 2, Quantity: 1}, bids[1])
		assert.Equal(t, &obm.Entry{Price: 1, Quantity: 1}, bids[2])
	})

	t.Run("Asks", func(t *testing.T) {
		asks := ob.Asks
		require.Len(t, ob.Bids, 3)
		assert.Equal(t, &obm.Entry{Price: 1, Quantity: 1}, asks[0])
		assert.Equal(t, &obm.Entry{Price: 2, Quantity: 1}, asks[1])
		assert.Equal(t, &obm.Entry{Price: 3, Quantity: 1}, asks[2])
	})
}
