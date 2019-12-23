package main

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"

	"github.com/gogo/protobuf/jsonpb"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type server struct {
	client types.HotWalletClient
}

func (s *server) getBalance(w http.ResponseWriter, req *http.Request) {
	var balreq types.GetBalanceRequest
	if err := jsonpb.Unmarshal(req.Body, &balreq); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	balres, err := s.client.GetBalance(context.Background(), &balreq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := new(jsonpb.Marshaler).Marshal(w, balres); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *server) getAllowance(w http.ResponseWriter, req *http.Request) {
	var alreq types.GetAllowanceRequest
	if err := jsonpb.Unmarshal(req.Body, &alreq); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	alres, err := s.client.GetAllowance(context.Background(), &alreq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := new(jsonpb.Marshaler).Marshal(w, alres); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	conn, err := grpc.Dial("0.0.0.0:42001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	svr := &server{types.NewHotWalletClient(conn)}
	mux := http.NewServeMux()
	mux.HandleFunc("/balance", svr.getBalance)
	mux.HandleFunc("/allowance", svr.getAllowance)

	log.Fatal(http.ListenAndServe("0.0.0.0:7999", mux))
}
