package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type EthClient interface {
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)
	TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
}

var (
	ethUrl                   string
	Client                   EthClient
	BlockHeaders             chan *types.Header
	BlockHeadersSubscription ethereum.Subscription
)

func Configure(ethreumUrl string) error {
	client, err := ethclient.Dial(ethreumUrl)
	if err != nil {
		return fmt.Errorf("unable to connect to ethereum: %v", err.Error())
	}

	channel := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), channel)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err.Error())
	}

	ethUrl = ethreumUrl
	Client = client
	BlockHeaders = channel
	BlockHeadersSubscription = sub
	return nil
}

func Reset() {
	BlockHeadersSubscription.Unsubscribe()

	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		logrus.Fatal(fmt.Errorf("unable to reconnect to ethereum: %v", err))
	}
	Client = client

	sub, err := client.SubscribeNewHead(context.Background(), BlockHeaders)
	if err != nil {
		logrus.Fatal(fmt.Errorf("failed to subscribe: %v", err))
	}

	BlockHeadersSubscription = sub
}
