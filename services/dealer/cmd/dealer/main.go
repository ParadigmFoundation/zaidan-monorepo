package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/peterbourgon/ff"

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
	_ = store

	cfg := core.DealerConfig{
		MakerBindAddress:     *makerBind,
		HotWalletBindAddress: *hwBind,
	}
	dealer, err := core.NewDealer(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	_ = dealer

	server := gethrpc.NewServer()
	service, err := rpc.NewService()
	if err != nil {
		panic(err)
	}

	if err := server.RegisterName("dealer", service); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe("0.0.0.0:8000", server.WebsocketHandler([]string{"*"})); err != nil {
		log.Fatal(err)
	}
}
