package grpc

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"sort"
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

func TestWatchTransactionAppendEndpoint(t *testing.T) {
	eth.Mock()

	testTxString := "0x71b044c65962a23ed50a6081177b2ec2711b32d9fb1c9b2c7a4b6d711bf98FF0"
	testTxHash := common.HexToHash(testTxString)
	ws := WatcherServer{
		TxWatching: watching.New(),
		log: testLogger,
	}
	_, watched := ws.TxWatching.IsWatched(testTxHash)
	assert.False(t, watched)
	r, err := ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: testTxString, QuoteId: "test", StatusUrls: []string{"one", "two"} })
	assert.NoErrorf(t, err, "Unexpected error")
	assert.True(t, r.IsWatched)
	watcher, _ := ws.TxWatching.IsWatched(testTxHash)
	assert.Equal(t, 2, len(watcher.UpdateUrls), "Not enough urls")

	r, err = ws.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: testTxString, QuoteId: "test", StatusUrls: []string{"five", "two"} })
	assert.NoErrorf(t, err, "Unexpected error")
	assert.True(t, r.IsWatched)
	watcher, _ = ws.TxWatching.IsWatched(testTxHash)
	assert.Equal(t, 3, len(watcher.UpdateUrls), "Not enough urls")
	sort.Strings(watcher.UpdateUrls)
	assert.Equal(t, watcher.UpdateUrls[sort.SearchStrings(watcher.UpdateUrls, "one")], "one")
	assert.Equal(t, watcher.UpdateUrls[sort.SearchStrings(watcher.UpdateUrls, "two")], "two")
	assert.Equal(t, watcher.UpdateUrls[sort.SearchStrings(watcher.UpdateUrls, "five")], "five")
}
