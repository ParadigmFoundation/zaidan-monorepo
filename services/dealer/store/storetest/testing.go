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

	suite.Require().NoError(
		suite.Store.CreateQuote(obj),
	)
	suite.Require().Len(obj.QuoteId, 36,
		"UUID from quote should be of length 36",
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

func (suite *Suite) TestPolicies() {
	policy := "xxx"

	found, err := suite.Store.HasPolicy(policy)
	suite.Require().NoError(err)
	suite.Require().False(found)

	suite.Require().NoError(
		suite.Store.CreatePolicy(policy),
	)

	found, err = suite.Store.HasPolicy(policy)
	suite.Require().NoError(err)
	suite.Require().True(found)

	suite.Run("Idempotence", func() {
		p := "foo"
		suite.Require().NoError(suite.Store.CreatePolicy(p))
		suite.Require().NoError(suite.Store.CreatePolicy(p))
	})
}
