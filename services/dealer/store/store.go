package store

import (
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer"
)

type Store interface {
	CreateTrade(*dealer.Trade) error
}
