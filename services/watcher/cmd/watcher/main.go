package main

import (
	"fmt"
	"strconv"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logging"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
	"github.com/spf13/cobra"
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
		logging.Fatal(err)
	}
}

func configureFlags() {
	flags := cmd.PersistentFlags()
	flags.StringVarP(&ethAddress, "eth", "e", "wss://ropsten.infura.io/ws", "Ethereum RPC url")
	flags.IntVarP(&port, "port", "p", 5001, "gRPC listen port")
	flags.StringVarP(&makerUrl, "maker", "m", "localhost:5002", "Maker gRPC url")
	flags.StringVarP(&bugsnagKey, "bugsnag", "b", "", "Bugsnag project key")
}

func startup(_ /*cmd*/ *cobra.Command, _ /*args*/ []string) {
	logging.ConfigureBugsnag(bugsnagKey)
	logging.Info("Starting")

	if err := eth.Configure(ethAddress); err != nil {
		logging.Fatal(err)
	}
	logging.Info("Connected to ethereum at", ethAddress)



	txWatching := watching.New(makerUrl)

	watcherServer := grpc.NewServer(
		txWatching,
	)
	if err := watcherServer.Listen(strconv.Itoa(port)); err != nil {
		logging.Fatal(fmt.Errorf("failed to listen: %v", err))
	}
}
