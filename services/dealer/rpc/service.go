package rpc

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
)

type Service struct {
	dealer *core.Dealer
	server *rpc.Server
}

// NewService creates a new Dealer JSONRPC service
func NewService(dealer *core.Dealer) (*Service, error) {
	srv := &Service{
		dealer: dealer,
		server: rpc.NewServer(),
	}

	if err := srv.server.RegisterName("dealer", srv); err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/ws" {
		srv.server.WebsocketHandler([]string{"*"}).ServeHTTP(w, r)
	} else {
		srv.server.ServeHTTP(w, r)
	}
}
