package gemini

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
	"github.com/gorilla/websocket"
)

const (
	FEED_URL = "wss://api.sandbox.gemini.com/v1/marketdata/"
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

func (x *Exchange) Subscribe(ctx context.Context, store store.Store, syms ...string) error {
	for i, _ := range syms {
		syms[i] = x.newSymbol(syms[i])
	}
	log.Printf("Gemini querying: %q", syms)

	for _, sym := range syms {
		x.subscribe(ctx, store, sym)
	}

	return nil
}

func (x *Exchange) subscribe(ctx context.Context, store store.Store, sym string) error {
	url := FEED_URL + sym //+ "?trades=false&auctions=false"
	c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil
	}

	go func() {
		if err := x.handleWs(c, store, sym); err != nil {
			// TODO(gchaincl): please
			panic(err)
		}
	}()

	select {}
}

func (x *Exchange) handleWs(c *websocket.Conn, s store.Store, sym string) error {
	for {
		var msg Message
		err := c.ReadJSON(&msg)
		if err != nil {
			return err
		}

		if msg.Type != "update" {
			continue
		}

		update := &obm.Update{
			Symbol: x.symbol(sym),
		}

		for _, event := range msg.Events {
			if event.Type != "change" {
				continue
			}

			entry, err := obm.NewEntryFromStrings(event.Price, event.Remaining)
			if err != nil {
				return err
			}

			switch event.Side {
			case "bid":
				update.Bids = append(update.Bids, entry)
			case "ask":
				update.Asks = append(update.Asks, entry)
			}
		}

		if len(update.Bids) != 0 || len(update.Asks) != 0 {
			s.OnUpdate("gemini", update)
		}
	}
}
