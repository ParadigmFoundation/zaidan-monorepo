package sql

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer"
)

type SQLSuite struct {
	suite.Suite
	store *Store
}

func (suite *SQLSuite) SetupTest() {
	s, err := New("sqlite3", ":memory:")
	suite.Require().NoError(err)
	suite.store = s
}

func (suite *SQLSuite) TestTrades() {
	t := &dealer.Trade{
		QuoteId:  "test-id",
		MarketId: "mkt-id",
	}
	suite.Require().NoError(
		suite.store.CreateTrade(t),
	)

	found, err := suite.store.GetTrade(t.QuoteId)
	suite.Require().NoError(err)
	suite.Require().NotNil(found)

	suite.Assert().True(proto.Equal(found, t))
}

func TestSQL(t *testing.T) {
	suite.Run(t, &SQLSuite{})
}
