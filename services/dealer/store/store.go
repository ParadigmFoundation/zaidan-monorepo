package store

import (
	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type Store interface {
	CreateQuote(*types.Quote) error
	GetQuote(string) (*types.Quote, error)
	CreateMarket(*types.Market) error
	GetMarket(string) (*types.Market, error)
}
