package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MockClient struct {

}

func (m MockClient) TransactionByHash(c context.Context, h common.Hash) (*types.Transaction, bool, error) {
	fmt.Println("mocked tbh",h)
	emptyTx := types.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)
	return emptyTx, false, nil
}

func (m MockClient) TransactionReceipt(c context.Context, h common.Hash) (*types.Receipt, error) {
	passed := types.NewReceipt(nil, false, 0)
	return passed, nil
}

func (m MockClient) BlockByNumber(context.Context, *big.Int) (*types.Block, error) {
	return nil, nil
}


type Sub struct {}

func (s Sub) Unsubscribe() {

}

func (s Sub) Err() <- chan error {
	return nil
}

func Mock() {
	once.Do(
		func () {
			channel := make(chan *types.Header)
			sub := Sub{}
			client := &MockClient{}

			Client = client
			BlockHeaders = channel
			BlockHeadersSubscription = sub
		},
	)
}
