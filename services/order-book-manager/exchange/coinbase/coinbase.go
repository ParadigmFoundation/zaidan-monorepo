package coinbase

import (
	"context"
	"fmt"
	"log"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
	"github.com/gorilla/websocket"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

const (
	FEED_URL = "wss://ws-feed.pro.coinbase.com"
)

var _ exchange.Exchange = &Exchange{}

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

func (x *Exchange) Subscribe(ctx context.Context, store store.Store, syms ...string) error {
	ws, err := x.dial(ctx)
	if err != nil {
		return err
	}

	req := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			coinbasepro.MessageChannel{Name: "level2", ProductIds: syms},
		},
	}

	log.Printf("Coinbase querying: %q", syms)

	if err := ws.WriteJSON(req); err != nil {
		return err
	}

	for {
		var msg coinbasepro.Message
		if err := ws.ReadJSON(&msg); err != nil {
			return err
		}

		var fn func(string, *obm.Update) error
		switch msg.Type {
		case "snapshot":
			fn = store.OnSnapshot
		case "l2update":
			fn = store.OnUpdate
		}

		if fn != nil {
			updates, err := newUpdates(&msg)
			if err != nil {
				return err
			}

			if err := fn("coinbase", updates); err != nil {
				return err
			}
		}
	}
}

// NewUpdate returns a new obm.Update given a coinbasepro.Message
func newUpdates(msg *coinbasepro.Message) (*obm.Update, error) {
	var updates = obm.Update{
		Symbol: msg.ProductID,
	}

	for _, bid := range msg.Bids {
		entry, err := obm.NewEntryFromStrings(bid.Price, bid.Size)
		if err != nil {
			return nil, err
		}
		updates.Bids = append(updates.Bids, entry)
	}

	for _, ask := range msg.Asks {
		entry, err := obm.NewEntryFromStrings(ask.Price, ask.Size)
		if err != nil {
			return nil, err
		}

		updates.Asks = append(updates.Asks, entry)
	}

	for _, change := range msg.Changes {
		entry, err := obm.NewEntryFromStrings(change.Price, change.Size)
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
