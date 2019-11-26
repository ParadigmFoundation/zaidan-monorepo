package store

import "github.com/ParadigmFoundation/zaidan-monorepo/services/obm"

type Store interface {
	OnSnapshot(string, *obm.Update) error
	OnUpdate(string, *obm.Update) error
	OrderBook(exchange, symbol string) (*obm.OrderBook, error)
}
