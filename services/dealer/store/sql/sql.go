package sql

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql/migrations"
)

var _ store.Store = &Store{}

type Store struct {
	db *sqlx.DB
}

func New(driver, dsn string) (*Store, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
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

func (s *Store) CreateTrade(t *types.Trade) error {
	_, err := s.db.Exec(
		`INSERT INTO trades VALUES($1, $2)`,
		t.QuoteId, t.MarketId,
	)
	return err
}

func (s *Store) GetTrade(quoteId string) (*types.Trade, error) {
	t := types.Trade{}
	err := s.db.QueryRow(
		`SELECT quote_id, market_id FROM trades WHERE quote_id = $1 LIMIT 1`,
		quoteId,
	).Scan(&t.QuoteId, &t.MarketId)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) CreateQuote(q *types.Quote) error {
	q.QuoteId = uuid.New().String()
	_, err := s.db.Exec(
		`INSERT INTO quotes VALUES($1, $2, $3, $4, $5)`,
		q.QuoteId,
		q.MakerAssetTicker,
		q.TakerAssetTicker,
		q.MakerAssetSize,
		q.QuoteAssetSize,
	)
	return err
}

func (s *Store) GetQuote(quoteId string) (*types.Quote, error) {
	stmt := `
		SELECT
			quote_id,
			maker_asset_ticker,
			taker_asset_ticker,
			maker_asset_size,
			quote_asset_size
		FROM quotes
		WHERE quote_id = $1 LIMIT 1
	`
	q := types.Quote{}
	err := s.db.
		QueryRow(stmt, quoteId).
		Scan(
			&q.QuoteId,
			&q.MakerAssetTicker,
			&q.TakerAssetTicker,
			&q.MakerAssetSize,
			&q.QuoteAssetSize,
		)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
