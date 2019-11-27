package binance

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
	"github.com/adshao/go-binance"
)

type Exchange struct {
	sMutex  sync.RWMutex
	symbols map[string]string
}

func New() *Exchange {
	return &Exchange{
		symbols: make(map[string]string),
	}
}

func (x *Exchange) depthHandler(s store.Store) binance.WsDepthHandler {
	fn := func(event *binance.WsDepthEvent) {
		update, err := x.newUpdates(event)
		if err != nil {
			log.Printf("depthHandler: ERROR: %v", err)
			return
		}
		s.OnUpdate("binance", update)
	}

	return fn
}

func (x *Exchange) newSymbol(s string) string {
	x.sMutex.Lock()
	defer x.sMutex.Unlock()

	newSymbol := strings.Replace(s, "/", "", 1)
	x.symbols[newSymbol] = s
	return newSymbol
}

func (x *Exchange) symbol(s string) string {
	x.sMutex.RLock()
	defer x.sMutex.RUnlock()

	if found := x.symbols[s]; found != "" {
		return found
	}
	return s
}

func (x *Exchange) Subscribe(ctx context.Context, s store.Store, syms ...string) error {
	errHandler := func(err error) {
		log.Printf("ERROR: %+v", err)
	}

	var (
		doneCh = make(chan struct{})
		stopCh = make(chan struct{})
	)

	for i, _ := range syms {
		syms[i] = x.newSymbol(syms[i])
	}
	log.Printf("Binance querying: %q", syms)
	for _, sym := range syms {
		done, stop, err := binance.WsDepthServe(sym, x.depthHandler(s), errHandler)

		if err != nil {
			return fmt.Errorf("Subscribe(): %w", err)
		}

		go func() {
			select {
			case <-stop:
				stopCh <- struct{}{}
			case <-done:
				doneCh <- struct{}{}
			}
		}()
	}

	select {
	case <-ctx.Done():
		stopCh <- struct{}{}
		return ctx.Err()
	case <-doneCh:
		return nil
	}
}

func (x *Exchange) newUpdates(event *binance.WsDepthEvent) (*obm.Update, error) {
	var updates = obm.Update{
		Symbol: x.symbol(event.Symbol),
	}

	for _, bid := range event.Bids {
		entry, err := obm.NewEntryFromStrings(bid.Price, bid.Quantity)
		if err != nil {
			return nil, err
		}
		updates.Bids = append(updates.Bids, entry)
	}

	for _, ask := range event.Asks {
		entry, err := obm.NewEntryFromStrings(ask.Price, ask.Quantity)
		if err != nil {
			return nil, err
		}
		updates.Asks = append(updates.Asks, entry)
	}

	return &updates, nil
}
