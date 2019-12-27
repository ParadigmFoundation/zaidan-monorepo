package eth

import (
	"context"
	"fmt"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logging"
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
		logging.Fatal(fmt.Errorf("unable to connect to ethereum: %v", err))
	}

	channel := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), channel)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to subscribe: %v", err))
	}

	return &EthereumToolkit{ ethUrl: ethUrl, Client: client, BlockHeaders: channel, BlockHeadersSubscription: sub }
}

func (etk *EthereumToolkit) Reset() {
	etk.BlockHeadersSubscription.Unsubscribe()

	client, err := ethclient.Dial(etk.ethUrl)
	if err != nil {
		logging.Fatal(fmt.Errorf("unable to reconnect to ethereum: %v", err))
	}
	etk.Client = client

	sub, err := client.SubscribeNewHead(context.Background(), etk.BlockHeaders)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to subscribe: %v", err))
	}

	etk.BlockHeadersSubscription = sub
}