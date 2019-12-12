package main

import (
	"log"
	"net"
	"strconv"

	pb "../proto"
	"../server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	gethAddress string
	makerEndpoint string
	port int

	cmd = &cobra.Command{
		Use:   "watcher",
		Short: "Zaidan Transaction Watcher",
		Run: startup,
	}
)


func main() {
	configureFlags()
	cmd.Execute()
}

func configureFlags() {
	flags := cmd.PersistentFlags()
	flags.StringVar(&gethAddress, "geth", "wss://eth-ropsten.ws.alchemyapi.io/ws/nUUajaRKoZM-645b32rSRMwNVhW2EP3w", "Geth endpoint")
	flags.IntVarP(&port, "port", "p", 5001, "gRPC listen port")
	flags.StringVar(&makerEndpoint, "maker", "localhost:5002", "Maker gRPC endpoint")
}

func startup (cmd *cobra.Command, args []string) {
	log.Println("Starting")
	lis, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	watcherServer := &server.WatcherServer{ GethAddress: gethAddress, MakerEndpoint: makerEndpoint }
	if err := watcherServer.Init(); err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}
	pb.RegisterWatcherServer(s, watcherServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
