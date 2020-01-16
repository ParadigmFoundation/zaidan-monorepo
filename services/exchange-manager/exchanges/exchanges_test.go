package exchanges

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/em"
)

type ExchangesSuite struct {
	suite.Suite
	exchange em.Exchange
}

func assertEqualOrders(t *testing.T, a, b *types.ExchangeOrder) {
	assert.True(t, a.Equal(b), spew.Sdump(a, b))
}

func protoDiff(a, b proto.Message) string {
	return spew.Sdump(a, b)
}

// Block will run a test block. Blocks are not isolated,meaning that if one block fails, the whole test is aborted.
func (suite *ExchangesSuite) Block(op string, fn func(t *testing.T)) {
	fmt.Printf("%-20s ", op+" ...")
	fn(suite.T())
	fmt.Printf("done\n")
}

// TestCRUD will create an order, get it and cancel it
func (suite *ExchangesSuite) TestCRUD() {
	ctx := context.Background()

	order := &types.ExchangeOrder{
		Symbol: "BTC/USD",
		Price:  "1",
		Side:   types.ExchangeOrder_BUY,
		Amount: "0.001",
	}

	// Create
	suite.Block("CreateOrder", func(t *testing.T) {
		require.NoError(t, suite.exchange.CreateOrder(ctx, order))
		require.NotEmpty(t, order.Id)
	})

	// READ
	var totalOrders int
	suite.Block("GetOrder", func(t *testing.T) {
		found, err := suite.exchange.GetOrder(ctx, order.Id)
		require.NoError(t, err)
		require.NotNil(t, found)

		assertEqualOrders(t, order, found.Order)
	})
	suite.Block("GetTotalOrders", func(t *testing.T) {
		orders, err := suite.exchange.GetOpenOrders(ctx)
		require.NoError(t, err)
		totalOrders = len(orders.Array)
	})

	// DELETE
	suite.Block("CancelOrder", func(t *testing.T) {
		_, err := suite.exchange.CancelOrder(ctx, order.Id)
		require.NoError(t, err)
	})
	suite.Block("GetOrder", func(t *testing.T) {
		_, err := suite.exchange.GetOrder(ctx, order.Id)
		require.Error(t, err)
	})
	suite.Block("GetTotalOrders", func(t *testing.T) {
		orders, err := suite.exchange.GetOpenOrders(ctx)
		require.NoError(t, err)
		assert.Len(t, orders.Array, totalOrders-1)
	})
}
