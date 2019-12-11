package rpc

type Service struct{}

// NewService creates a new Dealer JSONRPC service
func NewService() (*Service, error) {
	return &Service{}, nil
}
