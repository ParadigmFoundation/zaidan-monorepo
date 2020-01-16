package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/em/exchanges"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/em/grpc"
	"github.com/peterbourgon/ff"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

func atLeastOne(ss ...string) bool {
	for _, s := range ss {
		if s != "" {
			return true
		}
	}
	return false
}

func main() {
	fs := flag.NewFlagSet("em", flag.ExitOnError)
	var (
		bind                = fs.String("bind", "localhost:8080", "Server listen address")
		coinbase_sandbox    = fs.Bool("coinbase-sandbox", true, "Use Coinbase sandbox API")
		coinbase_key        = fs.String("coinbase-key", "", "Coinbase API key")
		coinbase_passphrase = fs.String("coinbase-passphrase", "", "Coinbase Passphrase")
		coinbase_secret     = fs.String("coinbase-secret", "", "Coinbase Secret")
	)
	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("EM"),
	)

	srv := grpc.NewServer()
	if atLeastOne(*coinbase_key, *coinbase_passphrase, *coinbase_secret) {
		log.Printf("Enabling Coinbase accounts:")
		var url = "https://api-public.sandbox.pro.coinbase.com"
		if *coinbase_sandbox == false {
			url = ""
		}
		x := exchanges.NewCoinbase(&coinbasepro.ClientConfig{
			BaseURL:    url,
			Key:        *coinbase_key,
			Passphrase: *coinbase_passphrase,
			Secret:     *coinbase_secret,
		})
		accz, err := x.Client().GetAccounts()
		if err != nil {
			log.Fatal(fmt.Errorf("Coinbase: %w", err))
		}
		for _, acc := range accz {
			log.Printf("-> %-6s balance: %s", acc.Currency, acc.Balance)
		}
		srv.RegisterExchange("coinbase", x)
	}

	log.Printf("Listening on %s", *bind)
	log.Fatal(srv.Listen(*bind))
}
