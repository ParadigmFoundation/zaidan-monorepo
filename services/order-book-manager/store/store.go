package store

import "github.com/ParadigmFoundation/zaidan-monorepo/services/obm"

type Store interface {
	OnSnapshot(string, *obm.Update) error
	OnUpdate(string, *obm.Update) error
	Market(exchange, symbol string) (*obm.Market, error)
}
