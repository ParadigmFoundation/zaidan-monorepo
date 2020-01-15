package rpc

import (
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
)

type Service struct {
	dealer *core.Dealer
}

// NewService creates a new Dealer JSONRPC service
func NewService(dealer *core.Dealer) (*Service, error) {
	return &Service{
		dealer: dealer,
	}, nil
}
