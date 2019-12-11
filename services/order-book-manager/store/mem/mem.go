package mem

import (
	"sort"
	"sync"
	"time"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
)

type market struct {
	lastUpdate time.Time
	bids       map[float64]float64
	asks       map[float64]float64
}

// these operates as db indexes
type markets map[string]*market
type exchanges map[string]markets

type Store struct {
	xch exchanges
	m   sync.RWMutex
}

func New() *Store {
	return &Store{
		xch: make(map[string]markets),
	}
}

func (s *Store) OnSnapshot(name string, update *obm.Update) error {
	return s.doUpdate(name, update)
}

func (s *Store) OnUpdate(name string, update *obm.Update) error {
	return s.doUpdate(name, update)
}

func (s *Store) findOrCreateMarket(name, symbol string) *market {
	if s.xch[name] == nil {
		s.xch[name] = make(map[string]*market)
	}

	mkt := s.xch[name][symbol]
	if mkt == nil {
		mkt = &market{
			bids: make(map[float64]float64),
			asks: make(map[float64]float64),
		}
		s.xch[name][symbol] = mkt
	}

	return mkt
}

func (s *Store) doUpdate(name string, update *obm.Update) error {
	s.m.Lock()
	defer s.m.Unlock()

	mkt := s.findOrCreateMarket(name, update.Symbol)

	for _, bid := range update.Bids {
		p, q := bid.Price, bid.Quantity
		if q == 0 {
			delete(mkt.bids, p)
		} else {
			mkt.bids[p] = q
		}
	}

	for _, ask := range update.Asks {
		p, q := ask.Price, ask.Quantity
		if q == 0 {
			delete(mkt.asks, p)
		} else {
			mkt.asks[p] = q
		}
	}
	mkt.lastUpdate = time.Now()

	return nil
}

func (s *Store) OrderBook(exchange, symbol string) (*grpc.OrderBookResponse, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	mkt := s.findOrCreateMarket(exchange, symbol)

	var asks grpc.OrderBookEntriesByPriceAsc
	for p, q := range mkt.asks {
		asks = append(asks, &grpc.OrderBookEntry{Price: p, Quantity: q})
	}
	sort.Sort(asks)

	var bids grpc.OrderBookEntriesByPriceDesc
	for p, q := range mkt.bids {
		bids = append(bids, &grpc.OrderBookEntry{Price: p, Quantity: q})
	}
	sort.Sort(bids)

	ob := &grpc.OrderBookResponse{
		Exchange: exchange,
		Symbol:   symbol,
		Asks:     asks,
		Bids:     bids,
	}
	if !mkt.lastUpdate.IsZero() {
		ob.LastUpdate = mkt.lastUpdate.Unix()
	}

	return ob, nil
}
