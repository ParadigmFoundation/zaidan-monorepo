package watching

import (
	"context"
	"sync"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ethereum/go-ethereum/common"
)

type WatchedTransaction struct {
	TxHash  common.Hash
	QuoteId string
}

type SafeWatchedTransactions struct {
	watchedTransactions map[common.Hash]WatchedTransaction
	sync.Mutex
}

type TxWatching struct {
	MakerEndpoint           string
	MakerClient             pb.MakerClient
	safeWatchedTransactions SafeWatchedTransactions
	log                     *logger.Entry
}

var bg = context.Background()

func New(makerClient pb.MakerClient) *TxWatching {
	watching := TxWatching{
		MakerClient:             makerClient,
		safeWatchedTransactions: SafeWatchedTransactions{watchedTransactions: make(map[common.Hash]WatchedTransaction)},
		log:                     logger.New("watching"),
	}

	go watching.startWatchingBlocks()

	return &watching
}

func (txW *TxWatching) Lock() {
	txW.safeWatchedTransactions.Lock()
}

func (txW *TxWatching) Unlock() {
	txW.safeWatchedTransactions.Unlock()
}

func (txW *TxWatching) IsWatched(txHash common.Hash) (WatchedTransaction, bool) {
	value, present := txW.safeWatchedTransactions.watchedTransactions[txHash]
	return value, present
}

func (txW *TxWatching) Watch(txHash common.Hash, quoteId string) {
	txW.safeWatchedTransactions.watchedTransactions[txHash] = WatchedTransaction{
		TxHash:  txHash,
		QuoteId: quoteId,
	}
}

func (txW *TxWatching) delete(txHash common.Hash) {
	delete(txW.safeWatchedTransactions.watchedTransactions, txHash)
}

func (txW *TxWatching) RequestRemoval(txHash common.Hash) {
	_, isWatched := txW.IsWatched(txHash)
	_, isPending, _ := eth.Client.TransactionByHash(bg, txHash)

	if !isPending && isWatched {
		txW.delete(txHash)
	}
}

func (txW *TxWatching) startWatchingBlocks() {
	for {

		select {
		case err := <-eth.BlockHeadersSubscription.Err():
			{
				txW.log.WithError(err).Error("subscription error")
				txW.log.Info("Attempting to reconnect")
				eth.Reset()
			}
		case headers, ok := <-eth.BlockHeaders:
			{
				txW.Lock()

				if !ok {
					txW.log.Fatal("headers channel failure")
				}

				block, err := eth.Client.BlockByNumber(bg, headers.Number)
				if err != nil {
					txW.log.WithError(err).Errorf("error getting block number %s", headers.Number.String())
					txW.log.Info("Attempting to reonnect")
					eth.Reset()
					eth.BlockHeaders <- headers
					return
				}

				for _, blockTx := range block.Transactions() {
					txHash := blockTx.Hash()

					if watchedTransaction, present := txW.IsWatched(txHash); present {
						txW.log.Info("Found", txHash.String(), "in Block #", block.Number().String())
						txW.delete(txHash)

						receipt, err := eth.Client.TransactionReceipt(bg, txHash)
						if err != nil {
							txW.log.WithError(err).Errorf("failure getting receipt for watched transaction %s", txHash.String())
						}

						_, err = txW.MakerClient.OrderStatusUpdate(bg, &pb.OrderStatusUpdateRequest{
							QuoteId: watchedTransaction.QuoteId,
							Status:  uint32(receipt.Status),
						})
						if err != nil {
							txW.log.WithError(err).Errorf("failure calling maker for transaction %s", txHash.String())
						}
					}
				}

				txW.Unlock()
			}
		}
	}
}
