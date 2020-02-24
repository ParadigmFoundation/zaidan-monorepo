package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/em"
)

type Server struct {
	m         sync.RWMutex
	exchanges map[string]em.Exchange
	grpc      *grpc.Server
}

func NewServer() *Server {
	srv := &Server{
		exchanges: make(map[string]em.Exchange),
		grpc:      grpc.NewServer(),
	}
	types.RegisterExchangeManagerServer(srv.grpc, srv)
	return srv
}

func (srv *Server) Listen(bind string) error {
	l, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}

	return srv.grpc.Serve(l)
}

func (srv *Server) RegisterExchange(name string, exchange em.Exchange) {
	srv.m.Lock()
	srv.exchanges[name] = exchange
	srv.m.Unlock()
}

func (srv *Server) getExchange(name string) (em.Exchange, error) {
	srv.m.RLock()
	defer srv.m.RUnlock()

	exchange, ok := srv.exchanges[name]
	if !ok {
		return nil, fmt.Errorf("exchange %s not found", name)
	}
	return exchange, nil
}

func (srv *Server) CreateOrder(ctx context.Context, req *types.ExchangeCreateOrderRequest) (*types.ExchangeOrderResponse, error) {
	ex, err := srv.getExchange(req.Exchange)
	if err != nil {
		return nil, err
	}

	if err := ex.CreateOrder(ctx, req.Order); err != nil {
		return nil, err
	}

	return &types.ExchangeOrderResponse{
		Order: req.Order,
	}, nil
}

func (srv *Server) GetOrder(ctx context.Context, req *types.ExchangeOrderRequest) (*types.ExchangeOrderResponse, error) {
	ex, err := srv.getExchange(req.Exchange)
	if err != nil {
		return nil, err
	}

	order, err := ex.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (srv *Server) GetOpenOrders(ctx context.Context, req *types.GetOpenOrdersRequest) (*types.ExchangeOrderArrayResponse, error) {
	ex, err := srv.getExchange(req.Exchange)
	if err != nil {
		return nil, err
	}

	return ex.GetOpenOrders(ctx)
}

func (srv *Server) CancelOrder(ctx context.Context, req *types.ExchangeOrderRequest) (*empty.Empty, error) {
	ex, err := srv.getExchange(req.Exchange)
	if err != nil {
		return nil, err
	}

	return ex.CancelOrder(ctx, req.Id)
}
