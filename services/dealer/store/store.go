package store

import (
	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type Policy interface {
	CreatePolicy(string) error
	HasPolicy(string) (bool, error)
}

type Store interface {
	Policy
	CreateQuote(*types.Quote) error
	GetQuote(string) (*types.Quote, error)
	CreateMarket(*types.Market) error
	GetMarket(string) (*types.Market, error)
}
