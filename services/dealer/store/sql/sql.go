package sql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql/migrations"
)

type Store struct {
	db *sqlx.DB
}

func New(driver, dsn string) (*Store, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func migrate(db *sqlx.DB) error {
	var runner interface {
		Run(*sqlx.DB) error
	}

	switch db.DriverName() {
	case "postgres":
		runner = &migrations.Postgres{}
	default:
		runner = &migrations.Migrations{}
	}

	return runner.Run(db)
}

func (s *Store) CreateTrade(t *dealer.Trade) error {
	_, err := s.db.Exec(
		`INSERT INTO trades VALUES(?, ?)`,
		t.QuoteId, t.MarketId,
	)
	return err
}

func (s *Store) GetTrade(quoteId string) (*dealer.Trade, error) {
	t := dealer.Trade{}
	err := s.db.QueryRow(
		`SELECT quote_id, market_id FROM trades WHERE quote_id = ? LIMIT 1`,
		quoteId,
	).Scan(&t.QuoteId, &t.MarketId)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
