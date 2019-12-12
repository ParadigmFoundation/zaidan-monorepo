package server

import (
	"context"
	"testing"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/stretchr/testify/assert"
)

func TestWatchTransaction(t *testing.T) {
	ws := WatcherServer{ GethAddress: "wss://eth-ropsten.ws.alchemyapi.io/ws/nUUajaRKoZM-645b32rSRMwNVhW2EP3w" } //TODO geth
	ws.Init() //TODO stub?
	transaction, err := ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: "0x71b044c65962a23ed50a6081177b2ec2711b32d9fb1c9b2c7a4b6d711bf98210"})
	assert.Equal(t, false, transaction.IsPending)
	assert.Equal(t, nil, err)
}

func TestGethInitFailure(t *testing.T) {
	ws := WatcherServer{ GethAddress: "fork" }
	err := ws.Init()
	assert.NotNil(t, err)
}