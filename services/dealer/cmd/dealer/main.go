package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/peterbourgon/ff"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
)

func main() {
	fs := flag.NewFlagSet("dealer", flag.ExitOnError)
	var (
		db        = fs.String("db", "sqlite3", "Database driver [sqlite3|postgres]")
		dsn       = fs.String("dsn", ":memory:", "Database's Data Source Name (see driver's doc for details)")
		makerBind = fs.String("maker", "0.0.0.0:50051", "The port to connect to a maker server over gRPC on")
		hwBind    = fs.String("hw", "0.0.0.0:42001", "The port to connect to a hot-wallet server over gRPC on")
	)
	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("DEALER"),
	)

	store, err := sql.New(*db, *dsn)
	if err != nil {
		log.Fatal(err)
	}

	cfg := core.DealerConfig{
		MakerBindAddress:     *makerBind,
		HotWalletBindAddress: *hwBind,
	}
	dealer, err := core.NewDealer(context.Background(), store, cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := gethrpc.NewServer()
	service, err := rpc.NewService(dealer)
	if err != nil {
		log.Fatal(err)
	}

	// setup root logger to log to stdout
	gethlog.Root().SetHandler(gethlog.StdoutHandler)

	if err := server.RegisterName("dealer", service); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe("0.0.0.0:8000", server.WebsocketHandler([]string{"*"})); err != nil {
		log.Fatal(err)
	}
}
