package main

import (
	"context"
	"log"
	"net/http"

	"github.com/spf13/pflag"

	"google.golang.org/grpc"

	"github.com/gogo/protobuf/jsonpb"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type config struct {
	hwUrl string
	bind  string
}

type server struct {
	client types.HotWalletClient
}

func (s *server) getBalance(w http.ResponseWriter, req *http.Request) {
	var balreq types.GetBalanceRequest
	if err := jsonpb.Unmarshal(req.Body, &balreq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balres, err := s.client.GetBalance(context.Background(), &balreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := new(jsonpb.Marshaler).Marshal(w, balres); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *server) getAllowance(w http.ResponseWriter, req *http.Request) {
	var alReq types.GetAllowanceRequest
	if err := jsonpb.Unmarshal(req.Body, &alReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	alRes, err := s.client.GetAllowance(context.Background(), &alReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := new(jsonpb.Marshaler).Marshal(w, alRes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	var cfg config
	pflag.StringVar(&cfg.hwUrl, "server", "0.0.0.0:42001", "host and port for the hot-wallet server")
	pflag.StringVar(&cfg.bind, "bind", "0.0.0.0:7999", "host and port to bind HTTP server to")

	conn, err := grpc.Dial(cfg.hwUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	svr := &server{types.NewHotWalletClient(conn)}
	mux := http.NewServeMux()
	mux.HandleFunc("/balance", svr.getBalance)
	mux.HandleFunc("/allowance", svr.getAllowance)

	log.Fatal(http.ListenAndServe(cfg.bind, mux))
}
