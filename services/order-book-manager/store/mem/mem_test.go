package mem

import (
	"math/rand"
	"runtime"
	"testing"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store/storetest"
	"github.com/stretchr/testify/require"
)

func TestMem(t *testing.T) {
	storetest.TestSuite(t, func(t *testing.T) (store.Store, func()) {
		return New(), func() {}
	})
}

func benchMem(b *testing.B, s store.Store, sym *obm.Symbol) {
	// generate 3k updates
	for i := 0; i < 3000; i++ {
		update := &obm.Update{Symbol: *sym,
			// Generate entries with prices ranging from 1 to 3000 and quantities ranging from 0 to 99
			Asks: obm.Entries{&obm.Entry{Price: float64(rand.Int31n(3000) + 1), Quantity: float64(rand.Int31n(300))}},
			Bids: obm.Entries{&obm.Entry{Price: float64(rand.Int31n(3000) + 1), Quantity: float64(rand.Int31n(300))}},
		}

		// update 5 exchanges
		s.OnUpdate("mem1", update)
		s.OnUpdate("mem2", update)
		s.OnUpdate("mem3", update)
		s.OnUpdate("mem4", update)
		s.OnUpdate("mem5", update)
	}
}

func BenchmarkMem(b *testing.B) {
	var s store.Store

	sym := obm.NewSymbol("BTC", "USD")
	for i := 0; i < b.N; i++ {
		s = New()
		benchMem(b, s, sym)
	}

	mkt, err := s.Market("mem1", sym.String())
	require.NoError(b, err)
	b.Logf("mem -> Asks: %d, Bids: %d", len(mkt.Asks), len(mkt.Bids))

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b.Logf("Total memory usage: %v/%v MiB (#%d)",
		m.Alloc/1024/1024, m.TotalAlloc/1024/1024, b.N,
	)
}
