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
	CreateTrade(*types.Trade) error
	GetTrade(string) (*types.Trade, error)
	CreateQuote(*types.Quote) error
	GetQuote(string) (*types.Quote, error)
	CreateAsset(*types.Asset) error
	GetAsset(string) (*types.Asset, error)
	CreateMarket(*types.Market) error
	GetMarket(string) (*types.Market, error)
	GetAssetByAddress(string) (*types.Asset, error)
}
