package watching

import (
	"context"
	"log"
	"sync"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ethereum/go-ethereum/common"
)

type WatchedTransaction struct {
	TxHash common.Hash
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
}

var bg = context.Background()

func New(makerClient pb.MakerClient ) *TxWatching {
	watching := TxWatching{
		MakerClient:             makerClient,
		safeWatchedTransactions: SafeWatchedTransactions{ watchedTransactions: make(map[common.Hash]WatchedTransaction) },
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
	_, isPending, _:= eth.Client.TransactionByHash(bg, txHash)

	if !isPending && isWatched {
		txW.delete(txHash)
	}
}

func (txW *TxWatching) startWatchingBlocks() {
	for {

		select {
			case errors := <- eth.BlockHeadersSubscription.Err(): {
				log.Println("Subscription error! ", errors)
				log.Println("Attempting to reconnect")
				eth.Reset()
			}
			case headers, ok := <- eth.BlockHeaders: {
				txW.Lock()

				if !ok {
					log.Fatal("Headers channel failure.")
				}

				block, err := eth.Client.BlockByNumber(bg, headers.Number)
				if err != nil {
					log.Println("Error getting block number:", headers.Number.String(), err)
					log.Println("Attempting to reconnect")
					eth.Reset()
					eth.BlockHeaders <- headers
					return
				}

				for _, blockTx := range block.Transactions() {
					txHash := blockTx.Hash()

					if watchedTransaction, present := txW.IsWatched(txHash); present {
						log.Println("Found", txHash.String(), "in Block #", block.Number().String())
						txW.delete(txHash)

						receipt, err := eth.Client.TransactionReceipt(bg, txHash)
						if err != nil {
							log.Println("Failure getting receipt for watched transaction", txHash.String(), ":", err)
						}

						_, err = txW.MakerClient.OrderStatusUpdate(bg, &pb.OrderStatusUpdateRequest{
							QuoteId: watchedTransaction.QuoteId,
							Status:  uint32(receipt.Status),
						})
						if err != nil {
							log.Println("Failure calling maker for transaction ", txHash.String(), ":", err)
						}
						//TODO: Can we resolve/escalate the previous two errors for some intervention
					}
				}

				txW.Unlock()
			}
		}

	}
}



