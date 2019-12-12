package storetest

import (
	"math"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
)

type Suite struct {
	suite.Suite
	store store.Store
}

func (suite *Suite) SetStore(store store.Store) {
	suite.store = store
}

func (suite *Suite) TestTrades() {
	obj := &types.Trade{
		QuoteId:          "quote-id",
		MarketId:         "mkt-id",
		OrderHash:        "order-hash",
		TransactionHash:  "transaction-hash",
		TakerAddress:     "taker-address",
		Timestamp:        time.Now().Unix(),
		MakerAssetTicker: "m/a/t",
		TakerAssetTicker: "t/a/t",
		MakerAssetAmount: math.MaxUint64 / 2,
		TakerAssetAmount: math.MaxUint64 / 2,
	}
	suite.Require().NoError(
		suite.store.CreateTrade(obj),
	)

	suite.Run("Found", func() {
		found, err := suite.store.GetTrade(obj.QuoteId)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)
		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found,
		)
	})

	suite.Run("NotFound", func() {
		found, err := suite.store.GetTrade(uuid.New().String())
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
}

func (suite *Suite) TestQuotes() {
	obj := &types.Quote{
		MakerAssetTicker: "maker-asset-ticker",
		TakerAssetTicker: "taker-asset-ticker",
		MakerAssetSize:   "maker-asset-size",
		QuoteAssetSize:   "quote-asset-size",
		Expiration:       10,
		OrderHash:        "order-hash",
		Order:            `["this", "is", {"json": "format"}]`,
		FillTx:           "fill-tx",
	}

	suite.Require().NoError(
		suite.store.CreateQuote(obj),
	)
	suite.Require().Len(obj.QuoteId, 36,
		"CreateTrade should set a UUID",
	)

	suite.Run("Found", func() {
		found, err := suite.store.GetQuote(obj.QuoteId)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)
		suite.Assert().Zero(found.Expiration, "expiration should not be set by the store")

		// overwite Expiration so that proto.Equal does not fail
		found.Expiration = obj.Expiration
		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found,
		)
	})

	suite.Run("NotFound", func() {
		found, err := suite.store.GetTrade(uuid.New().String())
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
}
