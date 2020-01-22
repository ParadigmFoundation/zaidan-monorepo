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
		MakerAssetAddress: "maker-asset-address",
		TakerAssetAddress: "taker-asset-address",
		MakerAssetSize:    "maker-asset-size",
		TakerAssetSize:    "taker-asset-size",
		Expiration:        time.Now().Add(1 * time.Second).Unix(),
		ServerTime:        time.Now().Unix(),
		OrderHash:         "order-hash",
		Order: &types.SignedOrder{
			ChainId:         1,
			ExchangeAddress: "exchange-address",
		},
		ZeroExTransactionHash: "tx-hash",
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
		QuoteId:           "quote-id",
		MarketId:          "mkt-id",
		OrderHash:         "order-hash",
		TransactionHash:   "transaction-hash",
		TakerAddress:      "taker-address",
		Timestamp:         time.Now().Unix(),
		MakerAssetAddress: "m/a/t",
		TakerAssetAddress: "t/a/t",
		MakerAssetAmount:  "10000000000000000",
		TakerAssetAmount:  "99999999999999999",
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

func (suite *Suite) TestMarkets() {
	obj := &types.Market{
		MakerAssetTicker:  "FOO/BAR",
		TakerAssetTickers: []string{"FOO/BAR", "XXX/YYY"},
		TradeInfo: &types.TradeInfo{
			ChainId:  123,
			GasPrice: "210000",
			GasLimit: "12000000000",
		},
		QuoteInfo: &types.QuoteInfo{
			MinSize: "100000000000000",
			MaxSize: "10000000000000000000000000",
		},
		Metadata: map[string]string{
			"this": "is",
			"a":    "test",
		},
	}
	suite.Require().NoError(
		suite.Store.CreateMarket(obj),
	)

	suite.Run("found", func() {
		found, err := suite.Store.GetMarket(obj.Id)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)

		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found)
	})
}
