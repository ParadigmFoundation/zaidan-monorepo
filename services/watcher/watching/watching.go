package watching

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

type TxWatching struct {
	Context   context.Context
	EthClient *ethclient.Client
	TxHashes  []common.Hash
}

func (txW *TxWatching) Init() error {

	headerChan := make(chan *types.Header)
	sub, err := txW.EthClient.SubscribeNewHead(context.Background(), headerChan)
	if err != nil {
		return fmt.Errorf("failed to subscribe" + err.Error())

	}
	
	go txW.watchBlock(headerChan)
	sub.Err() //TODO use error channel to attempt to reset connect a few times before crashing

	return nil
}

func (txW *TxWatching) watchBlock(headers chan *types.Header) {
	for {
		select {
			case headers := <-headers:
				block, err:= txW.EthClient.BlockByNumber(context.Background(), headers.Number)
				if err != nil {
					fmt.Println(err)
				}
				var foundTxs []common.Hash
				for _, blockTx := range block.Transactions() {
					txHash := blockTx.Hash()
					for wI, watchTx := range txW.TxHashes {
						if txHash == watchTx {
							log.Info("Found ", txHash, " in Block #", block.Number().String())
							foundTxs = append(foundTxs, txHash)
							txW.TxHashes[wI] = txW.TxHashes[len(txW.TxHashes)-1]
							txW.TxHashes = txW.TxHashes[:len(txW.TxHashes)-1]

							//TODO CALL TO CONFIRM
						}
					}
				}
		}
	}
}



