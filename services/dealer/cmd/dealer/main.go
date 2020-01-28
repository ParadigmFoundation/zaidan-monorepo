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

func unQuote(s string) string {
	isQuote := func(c byte) bool {
		return c == '\'' || c == '"'
	}

	n := len(s)
	if isQuote(s[0]) && s[0] == s[n-1] {
		return s[1 : n-1]
	}
	return s
}

func main() {
	fs := flag.NewFlagSet("dealer", flag.ExitOnError)
	var (
		db            = fs.String("db", "sqlite3", "Database driver [sqlite3|postgres]")
		dsn           = fs.String("dsn", ":memory:", "Database's Data Source Name (see driver's doc for details)")
		bind          = fs.String("bind", ":8000", "Server binding address")
		makerBind     = fs.String("maker", "localhost:50051", "The address to connect to a maker server over gRPC on")
		hwBind        = fs.String("hw", "localhost:42001", "The address to connect to a hot-wallet server over gRPC on")
		watcherBind   = fs.String("watcher", "localhost:43001", "The address to connect to a watcher server over gRPC on")
		orderDuration = fs.Int64("order-duration", 600, "The duration in seconds that signed order should be valid for")
		policyBlack   = fs.Bool("policy.blacklist", false, "Enable BlackList policy mode")
		policyWhite   = fs.Bool("policy.whitelist", false, "Enable WhiteList policy mode")
	)
	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("DEALER"),
	)
	*dsn = unQuote(*dsn)

	log := logger.New("app")
	log.WithFields(logrus.Fields{"db": *db, "dsn": *dsn}).Info("Initializing database")
	store, err := sql.New(*db, *dsn)
	if err != nil {
		log.Fatal(err)
	}

	cfg := core.DealerConfig{
		MakerBindAddress:     *makerBind,
		HotWalletBindAddress: *hwBind,
		WatcherBindAddress:   *watcherBind,
		OrderDuration:        *orderDuration,
	}
	log.WithField("cfg", cfg).Info("Dealer")
	dealer, err := core.NewDealer(context.Background(), store, cfg)
	if err != nil {
		log.Fatal(err)
	}

	service, err := rpc.NewService(dealer)
	if err != nil {
		log.Fatal(err)
	}

	if *policyBlack && *policyWhite {
		log.Fatal("Can't use both -policy.blacklist and -policy.whitelist")
	}

	if *policyBlack || *policyWhite {
		var mode rpc.PolicyMode
		if *policyBlack {
			log.Print("Using Blacklist mode")
			mode = rpc.PolicyBlackList
		}
		if *policyWhite {
			log.Print("Using Whitelist mode")
			mode = rpc.PolicyWhiteList
		}

		service.WithPolicy(mode, store)
	}

	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("bind", *bind).Info("Server started")
	server := &http.Server{Handler: service}
	log.Fatal(server.Serve(ln))
}
