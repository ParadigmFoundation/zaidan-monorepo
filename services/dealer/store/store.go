package store

import (
	"errors"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

var (
	ErrQuoteDoesNotExist = errors.New("quote does not exist")
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
	CreateTrade(*types.Trade) error
	GetTrade(string) (*types.Trade, error)
	UpdateTradeStatus(string, types.Trade_Status) error
}
