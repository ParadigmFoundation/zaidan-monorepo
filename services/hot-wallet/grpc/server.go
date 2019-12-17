package grpc

import (
	"net"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"google.golang.org/grpc"
)

// Server represents a hot-wallet gRPC server
type Server struct {
	grpc *grpc.Server
}

// NewServer creates a new gRPC hot-wallet server provided a service
func NewServer(service types.HotWalletServer) *Server {
	srv := &Server{
		grpc: grpc.NewServer(),
	}
	types.RegisterHotWalletServer(srv.grpc, service)
	return srv
}

// Listen binds the gRPC server to the provided bind address (TCP)
func (s *Server) Listen(bind string) error {
	l, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}

	return s.grpc.Serve(l)
}
