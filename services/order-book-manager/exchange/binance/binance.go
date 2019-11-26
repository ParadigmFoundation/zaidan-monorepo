package binance

import (
	"context"
	"fmt"
	"log"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
	"github.com/adshao/go-binance"
)

type Exchange struct {
}

func New() *Exchange {
	return &Exchange{}
}

func (x *Exchange) depthHandler(s store.Store) binance.WsDepthHandler {
	fn := func(event *binance.WsDepthEvent) {
		update, err := newUpdates(event)
		if err != nil {
			log.Printf("depthHandler: ERROR: %v", err)
			return
		}
		s.OnUpdate("binance", update)
	}

	return fn
}

func (x *Exchange) Subscribe(ctx context.Context, s store.Store, syms ...string) error {
	errHandler := func(err error) {
		log.Printf("ERROR: %+v", err)
	}

	var (
		doneCh = make(chan struct{})
		stopCh = make(chan struct{})
	)

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

func newUpdates(event *binance.WsDepthEvent) (*obm.Update, error) {
	var updates = obm.Update{
		Symbol: event.Symbol,
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
