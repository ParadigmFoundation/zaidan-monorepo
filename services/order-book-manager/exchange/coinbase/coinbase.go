package coinbase

import (
	"context"
	"fmt"
	"log"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/order-book-manager/exchange"
	"github.com/gorilla/websocket"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

const (
	FEED_URL = "wss://ws-feed.pro.coinbase.com"
)

type Exchange struct {
}

func New() *Exchange {
	return &Exchange{}
}

func (x *Exchange) dial(ctx context.Context) (*websocket.Conn, error) {
	var ws websocket.Dialer
	conn, _, err := ws.DialContext(ctx, FEED_URL, nil)
	if err != nil {
		return nil, fmt.Errorf("Coinbase: DialContext(): %w", err)
	}
	return conn, nil
}

func (x *Exchange) Subscribe(ctx context.Context, cb exchange.Callbacks, sym *exchange.Symbol, syms ...*exchange.Symbol) error {
	ws, err := x.dial(ctx)
	if err != nil {
		return err
	}

	// Build the ProductIDs based on the symbols
	ids := []string{sym.String()}
	for _, s := range syms {
		ids = append(ids, s.String())
	}

	req := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			coinbasepro.MessageChannel{Name: "level2", ProductIds: ids},
		},
	}

	log.Printf("-> %v", req)
	if err := ws.WriteJSON(req); err != nil {
		return err
	}

	for {
		var msg coinbasepro.Message
		if err := ws.ReadJSON(&msg); err != nil {
			return err
		}

		var fn func(*exchange.Update) error
		switch msg.Type {
		case "snapshot":
			fn = cb.OnSnapshot
		case "l2update":
			fn = cb.OnChange
		}

		if fn != nil {
			updates, err := newUpdates(&msg)
			if err != nil {
				return err
			}

			if err := fn(updates); err != nil {
				return err
			}
		}
	}
}

// NewUpdate returns a new exchange.Update given a coinbasepro.Message
func newUpdates(msg *coinbasepro.Message) (*exchange.Update, error) {
	sym, err := exchange.NewSymbolFromString(msg.ProductID, "-")
	if err != nil {
		return nil, err
	}

	var updates = exchange.Update{
		Symbol: *sym,
	}

	for _, bid := range msg.Bids {
		entry, err := exchange.NewEntryFromStrings(bid.Price, bid.Size)
		if err != nil {
			return nil, err
		}
		updates.Bids = append(updates.Bids, entry)
	}

	for _, ask := range msg.Asks {
		entry, err := exchange.NewEntryFromStrings(ask.Price, ask.Size)
		if err != nil {
			return nil, err
		}

		updates.Asks = append(updates.Asks, entry)
	}

	for _, change := range msg.Changes {
		entry, err := exchange.NewEntryFromStrings(change.Price, change.Size)
		if err != nil {
			return nil, err
		}

		switch change.Side {
		case "buy":
			updates.Bids = append(updates.Bids, entry)
		case "sell":
			updates.Asks = append(updates.Asks, entry)
		}
	}

	return &updates, nil
}
