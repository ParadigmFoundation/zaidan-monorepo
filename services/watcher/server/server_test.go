package server

import (
	"context"
	"testing"

	pb "../proto"
	"github.com/stretchr/testify/assert"
)

func TestWatchTransaction(t *testing.T) {
	ws := WatcherServer{ GethAddress: "https://ropsten.infura.io" } //TODO geth
	ws.Init() //TODO stub?
	transaction, err := ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: "0x71b044c65962a23ed50a6081177b2ec2711b32d9fb1c9b2c7a4b6d711bf98210"})
	assert.Equal(t, "0", transaction.TxStatus)
	assert.Equal(t, nil, err)
}

func TestGethInitFailure(t *testing.T) {
	ws := WatcherServer{ GethAddress: "fork" }
	err := ws.Init()
	assert.NotNil(t, err)
}