package watching

import (
	"context"
	"fmt"
	"log"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
)

type WatchedTransaction struct {
	TxHash common.Hash
	QuoteId string
}

type TxWatching struct {
	EthClient           *ethclient.Client
	WatchedTransactions map[common.Hash]WatchedTransaction
	MakerEndpoint       string
	MakerClient         pb.MakerClient
}

var bg = context.Background()

func (txW *TxWatching) Init() error {
	txW.WatchedTransactions = make(map[common.Hash]WatchedTransaction)

	conn, err := grpc.Dial(txW.MakerEndpoint, grpc.WithInsecure())

	if err != nil {
		return fmt.Errorf("failed to connect maker client " + err.Error())
	}

	txW.MakerClient = pb.NewMakerClient(conn)

	headerChan := make(chan *types.Header)
	sub, err := txW.EthClient.SubscribeNewHead(bg, headerChan)
	if err != nil {
		return fmt.Errorf("failed to subscribe " + err.Error())
	}

	// TODO use error channel to attempt to reset connect a few times before crashing
	go txW.watchBlock(headerChan, sub.Err())

	return nil
}

func (txW *TxWatching) watchBlock(headerChannel <-chan *types.Header, errorChannel <-chan error) {
	for {
		select {
			case errors := <- errorChannel:
				//TODO reset connection
				log.Fatal("Subscription error!", errors, len(headerChannel))
			case headers, ok := <- headerChannel:
				if !ok {
					fmt.Println("Headers died: ", len(headerChannel), ok)
				}

				block, err := txW.EthClient.BlockByNumber(bg, headers.Number)
				if err != nil {
					fmt.Println("Error getting block number: ", headers.Number.String(), err) //TODO Error handling
					return
				}

				for _, blockTx := range block.Transactions() {
					txHash := blockTx.Hash()

					if watchedTransaction, present := txW.WatchedTransactions[txHash]; present {
						log.Println("Found ", txHash.String(), " in Block #", block.Number().String())
						delete(txW.WatchedTransactions, txHash)

						receipt, err := txW.EthClient.TransactionReceipt(bg, txHash)
						if err != nil {
							fmt.Println(err) //TODO Error handling
						}

						//TODO CALL TO CONFIRM  //TODO Error handling
						_, err = txW.MakerClient.OrderStatusUpdate(context.Background(), &pb.OrderStatusUpdateRequest{
							QuoteId: watchedTransaction.QuoteId,
							Status:  uint32(receipt.Status),
						})

						if err != nil {
							log.Println("Failure calling maker: ", err)
						}
					}
				}
		}
	}
}



