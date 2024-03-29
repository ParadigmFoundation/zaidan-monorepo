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

func buildQuote() *types.Quote {
	return &types.Quote{
		QuoteId:               uuid.New().String(),
		MakerAssetAddress:     "maker-asset-address",
		TakerAssetAddress:     "taker-asset-address",
		MakerAssetSize:        "maker-asset-size",
		TakerAssetSize:        "taker-asset-size",
		Expiration:            time.Now().Add(1 * time.Second).Unix(),
		ServerTime:            time.Now().Unix(),
		ZeroExTransactionHash: "tx-hash",
		ZeroExTransactionInfo: &types.ZeroExTransactionInfo{
			Order: &types.SignedOrder{
				ChainId:         1,
				ExchangeAddress: "exchange-address",
			},
			Transaction: &types.ZeroExTransaction{
				Salt:                  "salty",
				Data:                  "data",
				ExpirationTimeSeconds: time.Now().Unix(),
			},
		},
	}
}

func (suite *Suite) TestQuotes() {
	obj := buildQuote()

	suite.Require().NoError(
		suite.Store.CreateQuote(obj),
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
		found, err := suite.Store.GetQuote(uuid.New().String())
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
}

func (suite *Suite) TestMarkets() {
	obj := &types.Market{
		MakerAssetAddress:   "0xfoo",
		TakerAssetAddresses: []string{"0xbar", "0xbuzz"},
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
		found, err := suite.Store.GetMarket(obj.MakerAssetAddress)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)

		suite.Assert().True(proto.Equal(obj, found),
			"\nexpected: %s\ngot:      %s", obj, found)
	})
}

func buildTrade() *types.Trade {
	return &types.Trade{
		Quote:       &types.Quote{},
		TxTimestamp: time.Now().UnixNano(),
	}
}

func (suite *Suite) TestTrades() {
	quote := buildQuote()
	suite.Require().NoError(
		suite.Store.CreateQuote(quote),
	)

	// Create the Trade
	trade := buildTrade()
	trade.Quote = quote
	suite.Require().NoError(
		suite.Store.CreateTrade(trade),
	)

	// Get the trade
	found, err := suite.Store.GetTrade(quote.QuoteId)
	suite.Require().NoError(err)
	suite.Assert().True(
		proto.Equal(trade, found),
	)

	suite.Run("StatusUpdate", func() {
		suite.Require().NoError(
			suite.Store.UpdateTradeStatus(quote.QuoteId, types.Trade_ERROR),
		)
		found, err := suite.Store.GetTrade(quote.QuoteId)
		suite.Require().NoError(err)
		suite.Assert().Equal(types.Trade_ERROR, found.Status)

		err2 := suite.Store.UpdateTradeStatus("this-id-does-not-exist", types.Trade_ERROR)
		suite.Require().Error(err2)
		suite.Assert().Equal(store.ErrQuoteDoesNotExist, err2)
	})

	suite.Run("NonExistent Quote", func() {
		t := buildTrade()
		t.Quote.QuoteId = "this-id-does-not-exist"
		err := suite.Store.CreateTrade(t)
		suite.Require().Error(err)
		suite.Assert().Equal(store.ErrQuoteDoesNotExist, err)
	})
}
