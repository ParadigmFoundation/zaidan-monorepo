package sql

import (
	"os"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type SQLSuite struct {
	suite.Suite

	driver string
	dsn    string

	store *Store
}

func (suite *SQLSuite) SetupTest() {
	s, err := New(suite.driver, suite.dsn)
	suite.Require().NoError(err)
	suite.store = s
}

func (suite *SQLSuite) TestTrades() {
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

func (suite *SQLSuite) TestQuotes() {
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

func TestSQL(t *testing.T) {
	t.Run("SQLite3", func(t *testing.T) {
		env := "DEALER_TEST_SQLITE"
		dsn := os.Getenv(env)
		if dsn == "" {
			dsn = ":memory:"
		}
		suite.Run(t, &SQLSuite{driver: "sqlite3", dsn: dsn})
	})

	t.Run("Postgres", func(t *testing.T) {
		env := "DEALER_TEST_PSQL"
		dsn := os.Getenv(env)
		if dsn == "" {
			t.Skipf("%s not defined", env)
		}

		suite.Run(t, &SQLSuite{driver: "postgres", dsn: dsn})
	})
}
