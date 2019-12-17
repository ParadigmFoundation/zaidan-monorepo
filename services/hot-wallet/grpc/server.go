package grpc

import (
	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	grpc *grpc.Server
}

func NewServer(service types.HotWalletServer) *Server {
	srv := &Server{
		grpc: grpc.NewServer(),
	}
	types.RegisterHotWalletServer(srv.grpc, service)
	return srv
}
