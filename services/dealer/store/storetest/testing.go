package storetest

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
)

type Suite struct {
	suite.Suite
	Store store.Store
}

func (suite *Suite) TestQuotes() {
	obj := &types.Quote{
		MakerAssetTicker: "maker-asset-ticker",
		TakerAssetTicker: "taker-asset-ticker",
		MakerAssetSize:   "maker-asset-size",
		QuoteAssetSize:   "quote-asset-size",
		Expiration:       time.Now().Add(1 * time.Second).Unix(),
		ServerTime:       time.Now().Unix(),
		OrderHash:        "order-hash",
		Order: &types.SignedOrder{
			ChainId:         1,
			ExchangeAddress: "exchange-address",
		},
		FillTx: "fill-tx",
	}

	suite.Require().NoError(
		suite.Store.CreateQuote(obj),
	)
	suite.Require().Len(obj.QuoteId, 36,
		"CreateTrade should set a UUID",
	)

	suite.Run("Found", func() {
		found, err := suite.Store.GetQuote(obj.QuoteId)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)

		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found,
		)
	})

	suite.Run("NotFound", func() {
		found, err := suite.Store.GetTrade(uuid.New().String())
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
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
		MakerAssetAmount: "10000000000000000",
		TakerAssetAmount: "99999999999999999",
	}
	suite.Require().NoError(
		suite.Store.CreateTrade(obj),
	)

	suite.Run("Found", func() {
		found, err := suite.Store.GetTrade(obj.QuoteId)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)
		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found,
		)
	})

	suite.Run("NotFound", func() {
		found, err := suite.Store.GetTrade(uuid.New().String())
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
}

func (suite *Suite) TestAssets() {
	obj := &types.Asset{
		Ticker:    "FOO/BAR",
		Name:      "Foo to Bar",
		Decimals:  18,
		NetworkId: 1,
		Address:   "0xdeadbeef",
	}

	suite.Require().NoError(
		suite.Store.CreateAsset(obj),
	)

	suite.Run("Found", func() {
		found, err := suite.Store.GetAsset(obj.Ticker)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)

		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found)
	})

	suite.Run("FoundByAddress", func() {
		found, err := suite.Store.GetAssetByAddress(obj.Address)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)

		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found)
	})

	suite.Run("NotFound", func() {
		found, err := suite.Store.GetAsset("XXX/YYY")
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
}
