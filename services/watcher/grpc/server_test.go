package grpc

import (
	"context"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
	"testing"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/stretchr/testify/assert"
)

func TestWatchTransaction(t *testing.T) {
	if err := eth.Configure("wss://ropsten.infura.io/ws"); err != nil {
		assert.Error(t, err, "Test connection failed.")
	}
	ws := WatcherServer{
		TxWatching: watching.New(nil),
	}
	ws.TxWatching.MakerClient = MockMakerClient{}
	transaction, err := ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: "0x71b044c65962a23ed50a6081177b2ec2711b32d9fb1c9b2c7a4b6d711bf98210", QuoteId: "test"})
	assert.Equal(t, false, transaction.IsPending)
	assert.Equal(t, nil, err)
}