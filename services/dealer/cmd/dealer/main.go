package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/peterbourgon/ff"
	"github.com/sirupsen/logrus"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
)

func main() {
	fs := flag.NewFlagSet("dealer", flag.ExitOnError)
	var (
		db        = fs.String("db", "sqlite3", "Database driver [sqlite3|postgres]")
		dsn       = fs.String("dsn", ":memory:", "Database's Data Source Name (see driver's doc for details)")
		bind      = fs.String("bind", ":8000", "Server binding address")
		makerBind = fs.String("maker", "0.0.0.0:50051", "The port to connect to a maker server over gRPC on")
		hwBind    = fs.String("hw", "0.0.0.0:42001", "The port to connect to a hot-wallet server over gRPC on")
	)
	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("DEALER"),
	)

	log := logger.New("app")
	log.WithFields(logrus.Fields{"db": *db, "dsn": *dsn}).Info("Initializing database")
	store, err := sql.New(*db, *dsn)
	if err != nil {
		log.Fatal(err)
	}
	_ = store

	cfg := core.DealerConfig{
		MakerBindAddress:     *makerBind,
		HotWalletBindAddress: *hwBind,
	}
	log.WithField("cfg", cfg).Info("Dealer")
	dealer, err := core.NewDealer(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	service, err := rpc.NewService(dealer)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("bind", *bind).Info("Server started")
	server := &http.Server{Handler: service}
	log.Fatal(server.Serve(ln))
}
