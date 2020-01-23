package main

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	ggrpc "google.golang.org/grpc"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
)

var (
	ethAddress string
	makerUrl   string
	port       int
	bugsnagKey string

	cmd = &cobra.Command{
		Use:   "watcher",
		Short: "Zaidan Transaction Watcher",
		Run:   startup,
	}
)

func main() {
	configureFlags()
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func configureFlags() {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&ethAddress, "eth", "e", "wss://ropsten.infura.io/ws", "Ethereum RPC url")
	flags.IntVarP(&port, "port", "p", 5001, "gRPC listen port")
	flags.StringVarP(&makerUrl, "maker", "m", "localhost:5002", "Maker gRPC url")
}

func startup(_ /*cmd*/ *cobra.Command, _ /*args*/ []string) {
	log := logger.New("app")
	log.Info("Starting")

	if err := eth.Configure(ethAddress); err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to ethereum at", ethAddress)

	conn, err := ggrpc.Dial(makerUrl, ggrpc.WithInsecure())
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect maker client: %w", err))
	}
	makerClient := pb.NewMakerClient(conn)
	log.Info("Maker client configured for ", makerUrl)

	txWatching := watching.New(makerClient)

	watcherServer := grpc.NewServer(
		txWatching,
	)
	if err := watcherServer.Listen(strconv.Itoa(port)); err != nil {
		log.Fatal(fmt.Errorf("failed to listen: %w", err))
	}
}
