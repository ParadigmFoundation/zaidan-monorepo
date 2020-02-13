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

var testLogger = logger.New("server_test")

func TestWatchTransaction(t *testing.T) {
	eth.Mock()

	ws := WatcherServer{
		TxWatching: watching.New(),
		log: testLogger,
	}
	transaction, err := ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: "0x71b044c65962a23ed50a6081177b2ec2711b32d9fb1c9b2c7a4b6d711bf9821C", QuoteId: "test"})
	assert.Equal(t, false, transaction.IsPending)
	assert.Equal(t, nil, err)
}

func TestWatchTransactionFailure(t *testing.T) {
	eth.Mock()

	ws := WatcherServer{
		TxWatching: watching.New(),
		log: testLogger,
	}
	_, err := ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: "0x71b044c65962a23ed50a6081177b2ec2711b32d9fb1c9b2c7a4b6d711bf98219", QuoteId: "test"})
	assert.Errorf(t, err, "testing")
}
