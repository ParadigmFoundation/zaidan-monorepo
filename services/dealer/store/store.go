package store

import (
	"github.com/paradigmfoundation/zaidan-monorepo/services/dealer"
)

type Store interface {
	CreateTrade(*dealer.Trade) error
}
