package store

import (
	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type Store interface {
	CreateTrade(*types.Trade) error
	GetTrade(string) (*types.Trade, error)
}
