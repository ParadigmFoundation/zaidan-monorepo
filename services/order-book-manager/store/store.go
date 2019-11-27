package store

import (
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange"
)

type Store interface {
	exchange.Subscriber
	OrderBook(exchange, symbol string) (*obm.OrderBook, error)
}
