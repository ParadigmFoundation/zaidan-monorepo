package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/utils"
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
	var m migrations.Interface

	switch db.DriverName() {
	case "postgres":
		m = &migrations.Postgres{}
	default:
		m = &migrations.SQLMigration{}
	}

	return migrations.Run(db, m)
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

func (s *Store) CreateQuote(q *types.Quote) error {
	return s.Atomic(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO quotes VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			q.QuoteId,
			q.MakerAssetAddress,
			q.TakerAssetAddress,
			q.MakerAssetSize,
			q.TakerAssetSize,
			q.Expiration,
			q.ServerTime,
			q.OrderHash,
			q.FillTx,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`INSERT INTO signed_orders(quote_id, order_bytes) VALUES ($1, $2)`, q.QuoteId, q.Order)
		return err
	})

}

func (s *Store) GetQuote(quoteId string) (*types.Quote, error) {
	stmt := `
		SELECT
		  q.quote_id
		, q.maker_asset_address
		, q.taker_asset_address
		, q.maker_asset_size
		, q.taker_asset_size
		, q.expiration
		, q.server_time
		, q.order_hash
		, q.fill_tx
		, s.order_bytes as "order"
		FROM quotes q
		LEFT JOIN signed_orders s
		ON q.quote_id = s.quote_id
		WHERE q.quote_id = $1 LIMIT 1
	`
	q := types.Quote{}
	err := s.db.Get(&q, stmt, quoteId)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (s *Store) CreateMarket(mkt *types.Market) error {
	_, err := s.db.Exec(
		`INSERT INTO markets VALUES($1, $2, $3, $4, $5)`,
		mkt.MakerAssetAddress,
		StringSlice(mkt.TakerAssetAddresses),
		mkt.TradeInfo,
		mkt.QuoteInfo,
		MapStringString(mkt.Metadata),
	)
	return err
}

func (s *Store) GetMarket(makerAssetAddress string) (*types.Market, error) {
	stmt := ` SELECT * FROM markets where maker_asset_address = $1 LIMIT 1`
	mkt := types.Market{}
	var addresses StringSlice
	var metadata MapStringString

	err := s.db.QueryRow(stmt, makerAssetAddress).Scan(
		&mkt.MakerAssetAddress,
		&addresses,
		&mkt.TradeInfo,
		&mkt.QuoteInfo,
		&metadata,
	)
	if err != nil {
		return nil, err
	}
	mkt.TakerAssetAddresses = addresses
	mkt.Metadata = metadata

	return &mkt, nil
}

type StringSlice []string

func (ss StringSlice) Value() (driver.Value, error) { return json.Marshal(ss) }
func (ss *StringSlice) Scan(v interface{}) error    { return json.Unmarshal(utils.AnyToBytes(v), ss) }

type MapStringString map[string]string

func (mss MapStringString) Value() (driver.Value, error) { return json.Marshal(mss) }
func (mss *MapStringString) Scan(v interface{}) error    { return json.Unmarshal(utils.AnyToBytes(v), mss) }
