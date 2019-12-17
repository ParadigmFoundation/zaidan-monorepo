package grpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange"
	_ "github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange/binance"
	_ "github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange/coinbase"
	_ "github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange/gemini"
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
	types.RegisterOrderBookManagerServer(srv.grpc, srv)
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

func (s *Server) OrderBook(ctx context.Context, req *types.OrderBookRequest) (*types.OrderBookResponse, error) {
	return s.store.OrderBook(req.Exchange, req.Symbol)
}

type channSubscriber struct {
	ch chan *obm.Update
}

func newChannSubscriber(ch chan *obm.Update) *channSubscriber { return &channSubscriber{ch: ch} }

func (s *channSubscriber) OnUpdate(name string, update *obm.Update) error {
	s.ch <- update
	return nil
}

func (s *channSubscriber) OnSnapshot(name string, update *obm.Update) error {
	s.ch <- update
	return nil
}

func (s *Server) Updates(req *types.OrderBookUpdatesRequest, stream types.OrderBookManager_UpdatesServer) error {
	// Get the exchange driver
	x, err := exchange.Get(req.Request.Exchange)
	if err != nil {
		return err
	}

	// Creates the update and error channels
	upCh := make(chan *obm.Update)
	errCh := make(chan error)

	// Subscribe to the exchange in the background
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sub := newChannSubscriber(upCh)
		errCh <- x.Subscribe(ctx, sub, []string{req.Request.Symbol}...)
	}()

	// Handle events
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case err := <-errCh:
			return err
		case update := <-upCh:
			resp := newResponseFromUpdate(update)
			resp.Exchange = req.Request.Exchange
			resp.Symbol = req.Request.Symbol
			stream.Send(resp)
		}
	}
}

// newResponseFromUpdate creates a proto representation of an obm.Update
func newResponseFromUpdate(update *obm.Update) *types.OrderBookResponse {
	resp := &types.OrderBookResponse{
		LastUpdate: time.Now().Unix(),
	}

	for _, ask := range update.Asks {
		resp.Asks = append(resp.Asks, &types.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		})
	}

	for _, bid := range update.Bids {
		resp.Bids = append(resp.Bids, &types.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		})
	}

	return resp
}
