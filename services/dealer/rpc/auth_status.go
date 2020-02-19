package rpc

import (
	"encoding/json"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc/policy"
)

type authStatusResponse struct {
	Authorized bool
	Reason     policy.Reason
}

func (status *authStatusResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		status.Authorized, status.Reason,
	})
}

func (svc *Service) AuthStatus(addr string) (*authStatusResponse, error) {
	if svc.policy == nil {
		return &authStatusResponse{Authorized: true, Reason: policy.WhiteListed}, nil
	}

	r, ok, err := svc.policy.AuthStatus(addr)
	if err != nil {
		return nil, err
	}

	return &authStatusResponse{
		Authorized: ok,
		Reason:     r,
	}, nil
}
