package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"


	"github.com/ethereum/go-ethereum/accounts"

	"github.com/caarlos0/env/v6"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet"
	"github.com/ParadigmFoundation/zaidan-monorepo/common/eth"
)

// a test server
type server struct {
	pvr *eth.Provider
	mux *http.ServeMux

	gsvr *grpc.Server
}

// test server config
type config struct {
	Ethurl   string `env:"ETHEREUM_JSONRPC_URL" envDefault:"http://localhost:8545"`
	Mnemonic string `env:"MNEMONIC" envDefault:"concert load couple harbor equip island argue ramp clarify fence smart topic"`
}

func newServer(provider *eth.Provider) *server {
	mux := http.NewServeMux()
	gsvr := grpc.NewServer(provider, 0, 0)
	svr := &server{pvr: provider, mux: mux, gsvr: gsvr}

	svr.mux.HandleFunc("/order/hash", svr.hashOrder)
	svr.mux.HandleFunc("/order/sign", svr.signOrder)
	svr.mux.HandleFunc("/order/marshal", svr.marshalOrder)
	svr.mux.HandleFunc("/order/create", svr.createOrder)

	return svr
}

func (s *server) start(path string, errChan chan<- error) {
	errChan <- http.ListenAndServe(path, s.mux)
}

func (s *server) createOrder(w http.ResponseWriter, r *http.Request) {
	var createOrderRequest *hw.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&createOrderRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signedOrderResponse, err := s.gsvr.CreateOrder(context.Background(), createOrderRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := json.Marshal(signedOrderResponse.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(output))
}

func (s *server) marshalOrder(w http.ResponseWriter, r *http.Request) {
	var order hw.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	zrxOrder, err := order.ToZeroExOrder()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := zrxOrder.ComputeOrderHash()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, hash.Hex())
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
		server.start("0.0.0.0:7999", errChan)
	}()

	log.Fatal(<-errChan)
}
