package sql

import (
	"os"
	"testing"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/storetest"
	"github.com/stretchr/testify/suite"
)

type SQLSuite struct {
	storetest.Suite
	driver string
	dsn    string
}

func (suite *SQLSuite) SetupTest() {
	s, err := New(suite.driver, suite.dsn)
	suite.Require().NoError(err)
	suite.Store = s
}

func (suite *SQLSuite) TearDownTest() {
	db := suite.Suite.Store.(*Store).db
	defer db.Close()

	// wipe
	for _, t := range []string{"migrations", "transaction_info", "quotes", "markets"} {
		_, err := db.Exec("DROP TABLE IF EXISTS " + t)
		suite.Require().NoError(err)
	}
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
