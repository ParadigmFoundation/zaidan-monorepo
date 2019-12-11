package sql

import (
	"os"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer"
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
