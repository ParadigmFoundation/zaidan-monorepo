package main

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
)

var (
	ethAddress string
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
	flags.StringVarP(&ethAddress, "eth", "e", "wss://eth-ropsten.ws.alchemyapi.io/ws/AAv0PpPC5GE3nqbj99bLqVhIsQKg7C-7", "Ethereum RPC url")
	flags.IntVarP(&port, "port", "p", 5001, "gRPC listen port")
}

func startup(_ /*cmd*/ *cobra.Command, _ /*args*/ []string) {
	log := logger.New("app")
	log.Info("Starting")

	if err := eth.Configure(ethAddress); err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to ethereum at ", ethAddress)

	txWatching := watching.New()

	watcherServer := grpc.NewServer(
		txWatching,
	)
	if err := watcherServer.Listen(strconv.Itoa(port)); err != nil {
		log.Fatal(fmt.Errorf("failed to listen: %w", err))
	}
}
