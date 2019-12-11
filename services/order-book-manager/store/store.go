package store

import (
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange"
)

type Store interface {
	exchange.Subscriber
	OrderBook(exchange, symbol string) (*grpc.OrderBookResponse, error)
}
