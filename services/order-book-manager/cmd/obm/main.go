package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange/binance"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/exchange/coinbase"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm/store/mem"
)

type Server struct {
	store store.Store
}

func NewServer(store store.Store) *Server {
	return &Server{store: store}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	if len(split) != 4 {
		http.Error(w, "Invalid Path, try '/coinbase/BTC/USD'", 400)
		return
	}

	sym := strings.Join(split[2:4], "/")
	ob, err := srv.store.OrderBook(split[1], sym)
	if err != nil {
		http.Error(w, err.Error(), 404)
	}

	json.NewEncoder(w).Encode(ob)
}

type xchange struct {
	name    string
	symbols []string
}

type config struct {
	bind      string
	exchanges []xchange
}

func parseArgs() (*config, error) {
	bind := flag.String("bind", "localhost:8000", "Webserver binding address")
	xz := flag.StringArray("exchange", nil, "Specify exchanges and symbols using <exchange>:<symbol[,symbol|...], you can use this flag many times to specify more than one exchange.\nExample:\n\t--exchange=coinbase:BTC-USD,ETH-USD --echange=binance:BTCUSDT")
	flag.Parse()

	cfg := &config{
		bind: *bind,
	}

	for _, x := range *xz {
		split := strings.Split(x, ":")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid format %s, use <echange>:<symbol[,symbol|...]", split)
		}

		xch := xchange{
			name:    split[0],
			symbols: strings.Split(split[1], ","),
		}
		if len(xch.symbols) == 0 || xch.symbols[0] == "" {
			return nil, fmt.Errorf("invalid format %s, no symbols defined, use <echange>:<symbol[,symbol|...]", split)
		}

		cfg.exchanges = append(cfg.exchanges, xch)
	}
	if len(cfg.exchanges) == 0 {
		return nil, fmt.Errorf("please specify at least one exchange")
	}

	return cfg, nil
}

func main() {
	cfg, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	store := mem.New()
	srv := NewServer(store)

	log.Printf("API Listening on %s", cfg.bind)
	go func() { http.ListenAndServe(cfg.bind, srv) }()

	errCh := make(chan error)
	ctx := context.Background()
	for _, i := range cfg.exchanges {
		var xch exchange.Exchange
		switch i.name {
		case "coinbase":
			xch = coinbase.New()
		case "binance":
			xch = binance.New()
		default:
			log.Fatalf("unknown exchange %s. Available exchanges: [coinbase,binance]", i.name)
		}

		syms := i.symbols
		go func() {
			err := xch.Subscribe(ctx, store, syms...)
			if err != nil {
				errCh <- err
			}
		}()
	}

	log.Fatal(<-errCh)
}
