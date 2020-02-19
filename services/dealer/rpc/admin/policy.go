package admin

import (
	"errors"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc/policy"
)

// AddTaker adds a taker to the registry to whitelist/blacklist them (depending on mode)
func (srv *Service) AddTaker(addr string) error {
	if srv.policy == nil {
		return errors.New("policy not enabled")
	}

	return srv.policy.Store().CreatePolicy(addr)
}

// RemoveTaker removes a taker from the registry (see _addTaker)
func (srv *Service) RemoveTaker(addr string) error {
	if srv.policy == nil {
		return errors.New("policy not enabled")
	}

	return srv.policy.Store().RemovePolicy(addr)
}

// GetTakers fetches a table of all takers in the registry
func (srv *Service) GetTakers() ([]string, error) {
	return srv.policy.Store().ListPolicies()
}

// GetTakerAuthorization fetches authorization status for a taker by their address
func (srv *Service) GetTakerAuthorization(addr string) (interface{}, error) {
	if srv.policy == nil {
		return nil, errors.New("policy not enabled")
	}

	r, ok, err := srv.policy.AuthStatus(addr)
	if err != nil {
		return nil, err
	}

	return struct {
		Reason     policy.Reason
		Authorized bool
	}{r, ok}, nil
}
