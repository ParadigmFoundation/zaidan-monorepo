package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/peterbourgon/ff"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/core"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
)

func main() {
	fs := flag.NewFlagSet("dealer", flag.ExitOnError)
	var (
		db          = fs.String("db", "sqlite3", "Database driver [sqlite3|postgres]")
		dsn         = fs.String("dsn", ":memory:", "Database's Data Source Name (see driver's doc for details)")
		bind        = fs.String("bind", ":8000", "Server binding address")
		makerBind   = fs.String("maker", "0.0.0.0:50051", "The port to connect to a maker server over gRPC on")
		hwBind      = fs.String("hw", "0.0.0.0:42001", "The port to connect to a hot-wallet server over gRPC on")
		policyBlack = fs.Bool("policy.blacklist", false, "Enable BlackList policy mode")
		policyWhite = fs.Bool("policy.whitelist", false, "Enable WhiteList policy mode")
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

	if err := http.ListenAndServe(*bind, service); err != nil {
		log.Fatal(err)
	}
}
