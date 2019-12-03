package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
)

type Server struct {
	store store.Store
	grpc  *grpc.Server
}

func NewServer(store store.Store) *Server {
	srv := &Server{
		store: store,
		grpc:  grpc.NewServer(),
	}
	obm.RegisterOrderBookManagerServer(srv.grpc, srv)
	return srv
}

func (s *Server) Listen(bind string) error {
	l, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}

	return s.grpc.Serve(l)
}

func (s *Server) Stop() { s.grpc.GracefulStop() }

func (s *Server) OrderBook(ctx context.Context, req *obm.OrderBookRequest) (*obm.OrderBookResponse, error) {
	return s.store.OrderBook(req.Exchange, req.Symbol)
}
