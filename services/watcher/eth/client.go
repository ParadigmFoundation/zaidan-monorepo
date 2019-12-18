package eth

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumToolkit struct {
	ethUrl 			   string
	Client             *ethclient.Client
	BlockHeaders       chan *types.Header
	SubscriptionErrors <-chan error
}

func Init (ethUrl string) *EthereumToolkit {
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		log.Fatal("Unable to connect to ethereum:" + err.Error())
	}

	channel := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), channel)
	if err != nil {
		log.Fatal("failed to subscribe" + err.Error())
	}

	return &EthereumToolkit{ ethUrl: ethUrl, Client: client, BlockHeaders: channel, SubscriptionErrors: sub.Err() }
}

func (etk *EthereumToolkit) Reset() {
	client, err := ethclient.Dial(etk.ethUrl)
	if err != nil {
		log.Fatal("Unable to reconnect to ethereum:" + err.Error())
	}

	sub, err := client.SubscribeNewHead(context.Background(), etk.BlockHeaders)
	if err != nil {
		log.Fatal("failed to subscribe" + err.Error())
	}

	etk.SubscriptionErrors = sub.Err()
}