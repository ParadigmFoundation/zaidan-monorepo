package storetest

import (
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
		QuoteId:  "quote-id",
		MarketId: "mkt-id",
	}
	suite.Require().NoError(
		suite.store.CreateTrade(obj),
	)

	suite.Run("Found", func() {
		found, err := suite.store.GetTrade(obj.QuoteId)
		suite.Require().NoError(err)
		suite.Require().NotNil(found)
		suite.Assert().True(proto.Equal(found, obj))
	})

	suite.Run("NotFound", func() {
		found, err := suite.store.GetTrade(uuid.New().String())
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
}

func (suite *Suite) TestQuotes() {
	obj := &types.Quote{}

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
		suite.Assert().True(proto.Equal(found, obj))
	})

	suite.Run("NotFound", func() {
		found, err := suite.store.GetTrade(uuid.New().String())
		suite.Assert().Error(err)
		suite.Assert().Nil(found)
	})
}
