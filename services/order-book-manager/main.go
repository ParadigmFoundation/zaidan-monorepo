package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/order-book-manager/exchange/coinbase"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/order-book-manager/exchange"
)

type Market struct {
	Bids map[float64]float64
	Asks map[float64]float64
}

func NewMarket() *Market {
	return &Market{
		Bids: make(map[float64]float64),
		Asks: make(map[float64]float64),
	}
}

type Server struct {
	markets map[string]*Market
	m       sync.RWMutex
}

func NewServer() *Server {
	return &Server{markets: make(map[string]*Market)}
}

func (srv *Server) OnSnapshot(up *exchange.Update) error {
	return srv.onUpdate(up)
}

func (srv *Server) OnChange(up *exchange.Update) error {
	return srv.onUpdate(up)
}

func (srv *Server) onUpdate(up *exchange.Update) error {
	srv.m.Lock()
	defer srv.m.Unlock()

	sym := up.Symbol.String()

	var mkt *Market
	if mkt = srv.markets[sym]; mkt == nil {
		mkt = NewMarket()
		srv.markets[sym] = mkt
	}

	var newBids int
	for _, bid := range up.Bids {
		p := bid.Price
		q := bid.Quantity
		if q == 0 {
			newBids--
			delete(mkt.Bids, p)
		} else {
			newBids++
			mkt.Bids[p] = q
		}
	}

	var newAsks int
	for _, ask := range up.Asks {
		p := ask.Price
		q := ask.Quantity
		if q == 0 {
			newAsks--
			delete(mkt.Asks, p)
		} else {
			newAsks++
			mkt.Asks[p] = q
		}
	}

	/*
		fmt.Printf("[%s] Bids: %d (%d), Asks: %d (%d)\n", sym,
			len(mkt.Bids), newBids,
			len(mkt.Asks), newAsks,
		)
	*/
	return nil
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	if len(split) != 3 {
		http.Error(w, "Invalid Path, try '/BTC/USD'", 400)
		return
	}

	sym := exchange.NewSymbol(split[1], split[2]).String()
	mkt := srv.markets[sym]
	if mkt == nil {
		http.Error(w, "Market "+sym+" does not exist", 404)
		return
	}

	srv.m.RLock()
	srv.m.RUnlock()

	var bids exchange.EntriesByPriceDesc
	for p, q := range mkt.Bids {
		entry := exchange.Entry{Price: p, Quantity: q}
		bids = append(bids, &entry)
	}
	sort.Sort(bids)

	var asks exchange.EntriesByPriceAsc
	for p, q := range mkt.Asks {
		entry := exchange.Entry{Price: p, Quantity: q}
		asks = append(asks, &entry)
	}
	sort.Sort(asks)

	json.NewEncoder(w).Encode(struct {
		Symbol string                      `json:"symbol"`
		Asks   exchange.EntriesByPriceAsc  `json:"asks"`
		Bids   exchange.EntriesByPriceDesc `json:"bids"`
	}{
		Symbol: sym, Asks: asks, Bids: bids,
	})
}

func main() {
	xch := coinbase.New()
	//xch := binance.New()

	srv := NewServer()

	go func() {
		http.ListenAndServe(":8000", srv)
	}()

	ctx := context.Background()
	err := xch.Subscribe(ctx,
		exchange.Callbacks{OnSnapshot: srv.OnSnapshot, OnChange: srv.OnChange},
		exchange.NewSymbol("BTC", "USD"),
		exchange.NewSymbol("ETH", "USD"),
		exchange.NewSymbol("LTC", "USD"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
