//go:generate protoc -I ../../.. --go_out=plugins=grpc:.. ../../../proto/watcher.proto

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
	flags.StringVar(&gethAddress, "geth", "https://ropsten.infura.io", "Geth endpoint")
	flags.IntVarP(&port, "port", "p", 5001, "gRPC listen port")
}

func startup (cmd *cobra.Command, args []string) {
	log.Println("Starting")
	lis, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	watcherServer := &server.WatcherServer{ GethAddress: gethAddress }
	if watcherServer.Init() != nil {
		log.Fatalf("failed connect to geth: %v", err)
	}
	pb.RegisterWatcherServer(s, watcherServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
