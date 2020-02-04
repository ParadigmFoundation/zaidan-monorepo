package grpc

import (
	"net"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"google.golang.org/grpc"
)

// Server represents a hot-wallet gRPC server
type Server struct {
	grpc *grpc.Server
}

// NewServer creates a new gRPC hot-wallet server provided a service
func NewServer(service types.HotWalletServer) *Server {
	log := logger.New("grpc")
	opt := grpc.UnaryInterceptor(logger.UnaryServerInterceptor(log))
	srv := &Server{
		grpc: grpc.NewServer(opt),
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

// Stop triggers a graceful shutdown of the gRPC server
func (s *Server) Stop() { s.grpc.GracefulStop() }

// ForceStop forcefully terminates all open connections (only use if Stop fails)
func (s *Server) ForceStop() { s.grpc.Stop() }
