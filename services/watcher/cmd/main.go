package main

import (
	"fmt"
	"strconv"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logging"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
	"github.com/spf13/cobra"
	ggrpc "google.golang.org/grpc"
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
	flags.StringVarP(&ethAddress, "eth", "e", "wss://eth-ropsten.ws.alchemyapi.io/ws/nUUajaRKoZM-645b32rSRMwNVhW2EP3w", "Ethereum RPC url")
	flags.IntVarP(&port, "port", "p", 5001, "gRPC listen port")
	flags.StringVarP(&makerUrl, "maker", "m", "localhost:5002", "Maker gRPC url")
	flags.StringVarP(&bugsnagKey, "bugsnag", "b", "", "Bugsnag project key")
}

func startup(_ /*cmd*/ *cobra.Command, _ /*args*/ []string) {
	logging.ConfigureBugsnag(bugsnagKey)
	logging.Info("Starting")

	ethToolkit := eth.New(ethAddress)
	logging.Info("Connected to ethereum at", ethAddress)

	conn, err := ggrpc.Dial(makerUrl, ggrpc.WithInsecure())
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to connect maker client: %v", err))
	}
	makerClient :=  pb.NewMakerClient(conn)
	logging.Info("Maker client configured for", makerUrl)

	txWatching := watching.New(ethToolkit, makerClient)


	watcherServer := grpc.NewServer(
		ethToolkit,
		txWatching,
	)
	if err := watcherServer.Listen(strconv.Itoa(port)); err != nil {
		logging.Fatal(fmt.Errorf("failed to listen: %v", err))
	}
}
