package sql

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
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

	// Map json tags to map to tables
	db.Mapper = reflectx.NewMapper("json")

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
	stmt := `
		INSERT INTO trades VALUES (
			:quote_id,
			:market_id,
			:order_hash,
			:transaction_hash,
			:taker_address,
			:timestamp,
			:maker_asset_ticker,
			:taker_asset_ticker,
			:maker_asset_amount,
			:taker_asset_amount
		);
	`
	_, err := s.db.NamedExec(stmt, t)
	return err
}

func (s *Store) GetTrade(quoteId string) (*types.Trade, error) {
	stmt := `
		SELECT
		  "quote_id"
		, "market_id"
		, "order_hash"
		, "transaction_hash"
		, "taker_address"
		, "timestamp"
		, "maker_asset_ticker"
		, "taker_asset_ticker"
		, "maker_asset_amount"
		, "taker_asset_amount"
		FROM trades WHERE quote_id = $1 LIMIT 1
	`
	t := types.Trade{}
	if err := s.db.Get(&t, stmt, quoteId); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) CreateQuote(q *types.Quote) error {
	q.QuoteId = uuid.New().String()
	orderBytes, err := proto.Marshal(q.Order)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`INSERT INTO quotes VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		q.QuoteId,
		q.MakerAssetTicker,
		q.TakerAssetTicker,
		q.MakerAssetSize,
		q.QuoteAssetSize,
		q.Expiration,
		q.ServerTime,
		q.OrderHash,
		orderBytes,
		q.FillTx,
	)
	return err
}

func (s *Store) GetQuote(quoteId string) (*types.Quote, error) {
	var orderBytes []byte

	stmt := `
		SELECT
		  "quote_id"
		, "maker_asset_ticker"
		, "taker_asset_ticker"
		, "maker_asset_size"
		, "quote_asset_size"
		, "expiration"
		, "server_time"
		, "order_hash"
		, "order"
		, "fill_tx"
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
			&q.Expiration,
			&q.ServerTime,
			&q.OrderHash,
			&orderBytes,
			&q.FillTx,
		)
	if err != nil {
		return nil, err
	}

	var order types.SignedOrder
	if err := proto.Unmarshal(orderBytes, &order); err != nil {
		return nil, err
	}
	q.Order = &order

	return &q, nil
}
