package main

import (
	"net/http"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/paradigmfoundation/zaidan-monorepo/services/dealer/rpc"
)

func main() {
	server := gethrpc.NewServer()
	service, err := rpc.NewService()
	if err != nil {
		panic(err)
	}

	if err := server.RegisterName("dealer", service); err != nil {
		panic(err)
	}

	http.ListenAndServe("0.0.0.0:8000", server.WebsocketHandler([]string{"*"}))
}
