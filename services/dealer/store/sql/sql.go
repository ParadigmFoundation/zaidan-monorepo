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

type AtomicFn func(tx *sqlx.Tx) error

// Atomic runs an AtomicFn inside a transaction (Begin).
// If fn returns an error the Tx is rolledback and the error is returned,
// otherwise the Tx is committed.
func (s *Store) Atomic(fn AtomicFn) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
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
		SELECT * FROM trades WHERE quote_id = $1 LIMIT 1
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

	return s.Atomic(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO quotes VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			q.QuoteId,
			q.MakerAssetTicker,
			q.TakerAssetTicker,
			q.MakerAssetSize,
			q.QuoteAssetSize,
			q.Expiration,
			q.ServerTime,
			q.OrderHash,
			q.FillTx,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`INSERT INTO signed_orders(quote_id, order_bytes) VALUES ($1, $2)`, q.QuoteId, orderBytes)
		return err
	})

}

func (s *Store) GetQuote(quoteId string) (*types.Quote, error) {
	var orderBytes []byte
	stmt := `
		SELECT
		  q.quote_id
		, q.maker_asset_ticker
		, q.taker_asset_ticker
		, q.maker_asset_size
		, q.quote_asset_size
		, q.expiration
		, q.server_time
		, q.order_hash
		, q.fill_tx
		, s.order_bytes
		FROM quotes q
		LEFT JOIN signed_orders s
		ON q.quote_id = s.quote_id
		WHERE q.quote_id = $1 LIMIT 1
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
			&q.FillTx,
			&orderBytes,
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
