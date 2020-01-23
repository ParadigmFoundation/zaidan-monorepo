package rpc

import "encoding/json"

type authStatusResponse struct {
	Authorized bool
	Reason     string
}

func (status *authStatusResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		status.Authorized, status.Reason,
	})
}

func (svc *Service) AuthStatus(addr string) (*authStatusResponse, error) {
	if svc.policy == nil {
		return &authStatusResponse{Authorized: true, Reason: "WHITELISTED"}, nil
	}

	found, err := svc.policy.HasPolicy(addr)
	if err != nil {
		return nil, err
	}

	resp := &authStatusResponse{}
	switch svc.policyMode {
	case PolicyBlackList:
		resp.Authorized = !found
	case PolicyWhiteList:
		resp.Authorized = found
	}

	if resp.Authorized {
		resp.Reason = "WHITELISTED"
	} else {
		resp.Reason = "BLACKLISTED"
	}

	return resp, nil
}
