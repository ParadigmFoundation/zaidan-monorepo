package grpc

import (
	"context"
	"testing"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/stretchr/testify/assert"
)

func TestWatchTransaction(t *testing.T) {
	if err := eth.Configure("wss://eth-ropsten.ws.alchemyapi.io/ws/AAv0PpPC5GE3nqbj99bLqVhIsQKg7C-7"); err != nil {
		assert.NoError(t, err, "Test connection failed.")
		t.Fatal()
	}
	ws := WatcherServer{
		TxWatching: watching.New(),
		log: logger.New("test"),
	}
	transaction, err := ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: "0x71b044c65962a23ed50a6081177b2ec2711b32d9fb1c9b2c7a4b6d711bf98210", QuoteId: "test"})
	assert.Equal(t, false, transaction.IsPending)
	assert.Equal(t, nil, err)
}
