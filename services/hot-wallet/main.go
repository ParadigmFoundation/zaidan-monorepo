package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/caarlos0/env/v6"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/hw/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/hw/zeroex"
)

// a test server
type server struct {
	pvr *eth.Provider
	mux *http.ServeMux
}

// test server config
type config struct {
	Ethurl   string `env:"ETHEREUM_JSONRPC_URL" envDefault:"http://localhost:8545"`
	Mnemonic string `env:"MNEMONIC" envDefault:"concert load couple harbor equip island argue ramp clarify fence smart topic"`
}

func newServer(provider *eth.Provider) *server {
	mux := http.NewServeMux()
	svr := &server{pvr: provider, mux: mux}

	svr.mux.HandleFunc("/order/hash", svr.hashOrder)
	svr.mux.HandleFunc("/order/sign", svr.signOrder)

	return svr
}

func (s *server) start(path string, errChan chan<- error) {
	errChan <- http.ListenAndServe(path, s.mux)
}

func (s *server) hashOrder(w http.ResponseWriter, r *http.Request) {
	var order zeroex.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := order.ComputeOrderHash()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to compute hash, likely an invalid order: %v", err), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, `{"orderHash":"%s"}`, hash.Hex())
}

func (s *server) signOrder(w http.ResponseWriter, r *http.Request) {
	var order zeroex.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signedOrder, err := zeroex.SignOrder(s.pvr, &order)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sign order, likely an invalid order: %v", err), http.StatusInternalServerError)
	}

	signedOrderJson, err := json.Marshal(signedOrder)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal signed order: %v", err), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s", string(signedOrderJson))
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	path := accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}
	provider, err := eth.NewProvider(cfg.Ethurl, cfg.Mnemonic, path)
	if err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error)
	server := newServer(provider)
	go func() {
		server.start("0.0.0.0:8000", errChan)
	}()

	log.Fatal(<-errChan)
}
