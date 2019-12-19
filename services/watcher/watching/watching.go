package watching

import (
	"context"
	"log"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ethereum/go-ethereum/common"
)

type WatchedTransaction struct {
	TxHash common.Hash
	QuoteId string
}

type TxWatching struct {
	EthToolkit          *eth.EthereumToolkit
	watchedTransactions map[common.Hash]WatchedTransaction
	MakerEndpoint       string
	MakerClient         pb.MakerClient
}

var bg = context.Background()

func Init(ethToolkit *eth.EthereumToolkit, makerClient pb.MakerClient ) *TxWatching {
	watching := TxWatching{
		EthToolkit:    ethToolkit,
		MakerClient: makerClient,
	}
	watching.watchedTransactions = make(map[common.Hash]WatchedTransaction)

	go watching.watchBlocks()

	return &watching
}

func (txW *TxWatching) IsWatched(txHash common.Hash) (WatchedTransaction, bool) {
	value, present := txW.watchedTransactions[txHash]
	return value, present
}

func (txW *TxWatching) Watch(txHash common.Hash, quoteId string) {
	txW.watchedTransactions[txHash] = WatchedTransaction{
		TxHash:  txHash,
		QuoteId: quoteId,
	}
}

func (txW *TxWatching) watchBlocks() {
	for {
		select {
			case errors := <- txW.EthToolkit.BlockHeadersSubscription.Err(): {
				log.Println("Subscription error! ", errors)
				log.Println("Attempting to reconnect")
				txW.EthToolkit.Reset()
			}
			case headers, ok := <- txW.EthToolkit.BlockHeaders: {
				log.Println(headers.Number.String())// TODO remove this
				if !ok {
					log.Fatal("Headers channel failure.")
				}

				block, err := txW.EthToolkit.Client.BlockByNumber(bg, headers.Number)
				if err != nil {
					log.Println("Error getting block number:", headers.Number.String(), err)
					log.Println("Attempting to reconnect")
					txW.EthToolkit.Reset()
					txW.EthToolkit.BlockHeaders <- headers
					return
				}

				for _, blockTx := range block.Transactions() {
					txHash := blockTx.Hash()

					if watchedTransaction, present := txW.watchedTransactions[txHash]; present {
						log.Println("Found", txHash.String(), "in Block #", block.Number().String())
						delete(txW.watchedTransactions, txHash)

						receipt, err := txW.EthToolkit.Client.TransactionReceipt(bg, txHash)
						if err != nil {
							log.Println(err) //TODO Error handling
						}

						//TODO CALL TO CONFIRM  //TODO Error handling
						_, err = txW.MakerClient.OrderStatusUpdate(context.Background(), &pb.OrderStatusUpdateRequest{
							QuoteId: watchedTransaction.QuoteId,
							Status:  uint32(receipt.Status),
						})

						if err != nil {
							log.Println("Failure calling maker:", err)
						}
					}
				}
			}
		}
	}
}



