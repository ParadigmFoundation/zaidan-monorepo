package rpc

import (
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
)

type Service struct {
	dealer     *core.Dealer
	policyMode PolicyMode
	policy     store.Policy
}

// NewService creates a new Dealer JSONRPC service
func NewService(dealer *core.Dealer) (*Service, error) {
	return &Service{
		dealer: dealer,
	}, nil
}

func (srv *Service) WithPolicy(mode PolicyMode, policy store.Policy) *Service {
	srv.policyMode = mode
	srv.policy = policy
	return srv
}
