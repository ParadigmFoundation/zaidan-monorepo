package rpc

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
)

type Service struct {
	dealer     *core.Dealer
	server     *rpc.Server
	policyMode PolicyMode
	policy     store.Policy
	log        *logger.Entry
}

// NewService creates a new Dealer JSONRPC service
func NewService(dealer *core.Dealer) (*Service, error) {
	srv := &Service{
		dealer: dealer,
		server: rpc.NewServer(),
		log:    logger.New("rpc", logger.HandleEthLog()),
	}

	if err := srv.server.RegisterName("dealer", srv); err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/ws" {
		srv.log.Infof("New WebSocket connection")
		srv.server.WebsocketHandler([]string{"*"}).ServeHTTP(w, r)
	} else {
		srv.log.Infof("%s %s", r.Method, path)
		srv.server.ServeHTTP(w, r)
	}
}

func (srv *Service) WithPolicy(mode PolicyMode, policy store.Policy) *Service {
	srv.policyMode = mode
	srv.policy = policy
	return srv
}
