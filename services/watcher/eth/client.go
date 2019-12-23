package eth

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumToolkit struct {
	ethUrl 			   string
	Client             *ethclient.Client
	BlockHeaders       chan *types.Header
	BlockHeadersSubscription ethereum.Subscription
}

func New(ethUrl string) *EthereumToolkit {
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		log.Fatal("Unable to connect to ethereum:" + err.Error())
	}

	channel := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), channel)
	if err != nil {
		log.Fatal("failed to subscribe" + err.Error())
	}

	return &EthereumToolkit{ ethUrl: ethUrl, Client: client, BlockHeaders: channel, BlockHeadersSubscription: sub }
}

func (etk *EthereumToolkit) Reset() {
	etk.BlockHeadersSubscription.Unsubscribe()

	client, err := ethclient.Dial(etk.ethUrl)
	if err != nil {
		log.Fatal("Unable to reconnect to ethereum:" + err.Error())
	}
	etk.Client = client

	sub, err := client.SubscribeNewHead(context.Background(), etk.BlockHeaders)
	if err != nil {
		log.Fatal("failed to subscribe" + err.Error())
	}

	etk.BlockHeadersSubscription = sub
}