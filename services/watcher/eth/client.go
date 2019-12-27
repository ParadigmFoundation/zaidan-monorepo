package eth

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethUrl 			   string
	Client             *ethclient.Client
	BlockHeaders       chan *types.Header
	BlockHeadersSubscription ethereum.Subscription
)

func Configure(ethreumUrl string) {
	client, err := ethclient.Dial(ethreumUrl)
	if err != nil {
		log.Fatal("Unable to connect to ethereum:" + err.Error())
	}

	channel := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), channel)
	if err != nil {
		log.Fatal("failed to subscribe" + err.Error())
	}

	ethUrl = ethreumUrl
	Client = client
	BlockHeaders = channel
	BlockHeadersSubscription = sub
}

func Reset() {
	BlockHeadersSubscription.Unsubscribe()

	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		log.Fatal("Unable to reconnect to ethereum:" + err.Error())
	}
	Client = client

	sub, err := client.SubscribeNewHead(context.Background(), BlockHeaders)
	if err != nil {
		log.Fatal("failed to subscribe" + err.Error())
	}

	BlockHeadersSubscription = sub
}