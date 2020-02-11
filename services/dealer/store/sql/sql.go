package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"

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

	if driver == "sqlite3" {
		if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
			return nil, err
		}
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
			`INSERT INTO quotes VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
			q.QuoteId,
			q.MakerAssetAddress,
			q.TakerAssetAddress,
			q.MakerAssetSize,
			q.TakerAssetSize,
			q.Expiration,
			q.ServerTime,
			q.ZeroExTransactionHash,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`INSERT INTO transaction_infos(quote_id, transaction_info_bytes) VALUES ($1, $2)`, q.QuoteId, q.ZeroExTransactionInfo)
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
		, q.zero_ex_transaction_hash
		, t.transaction_info_bytes as "zero_ex_transaction_info"
		FROM quotes q
		LEFT JOIN transaction_infos t
		ON q.quote_id = t.quote_id
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
	stmt := `SELECT * FROM markets where maker_asset_address = $1 LIMIT 1`
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

func (s *Store) CreatePolicy(t string) error {
	exists, err := s.HasPolicy(t)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = s.db.Exec("INSERT INTO policies VALUES($1)", t)
	return err
}

func (s *Store) HasPolicy(t string) (bool, error) {
	var count int

	stmt := `SELECT COUNT(*) FROM policies WHERE entry = $1`
	if err := s.db.QueryRow(stmt, t).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Store) CreateTrade(t *types.Trade) error {
	_, err := s.db.Exec(
		"INSERT INTO trades VALUES($1, $2, $3, $4)",
		t.Quote.QuoteId,
		t.TxTimestamp,
		t.TxHash,
		types.Trade_NOT_SET,
	)

	return WrapError(err)
}

func (s *Store) GetTrade(id string) (*types.Trade, error) {
	quote, err := s.GetQuote(id)
	if err != nil {
		return nil, WrapError(err)
	}

	var trade = &types.Trade{
		Quote: quote,
	}
	stmt := `SELECT tx_timestamp, tx_hash, status FROM trades WHERE quote_id = $1 LIMIT 1`
	if err := s.db.QueryRow(stmt, id).Scan(
		&trade.TxTimestamp,
		&trade.TxHash,
		&trade.Status,
	); err != nil {
		return nil, WrapError(err)
	}
	return trade, nil
}

func (s *Store) UpdateTradeStatus(id string, status types.Trade_Status) error {
	res, err := s.db.Exec("UPDATE trades SET status = $1 WHERE quote_id = $2", status, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return store.ErrQuoteDoesNotExist
	}

	return nil
}

// WrapError converts a DB specific error into a local error
// if the error can't be handled, the original error is returned
func WrapError(err error) error {
	switch t := err.(type) {
	case sqlite3.Error:
		if t.Code == sqlite3.ErrConstraint {
			return store.ErrQuoteDoesNotExist
		}
	case *pq.Error:
		if t.Code.Name() == "foreign_key_violation" {
			return store.ErrQuoteDoesNotExist
		}
	}

	return err
}

type StringSlice []string

func (ss StringSlice) Value() (driver.Value, error) { return json.Marshal(ss) }
func (ss *StringSlice) Scan(v interface{}) error    { return json.Unmarshal(utils.AnyToBytes(v), ss) }

type MapStringString map[string]string

func (mss MapStringString) Value() (driver.Value, error) { return json.Marshal(mss) }
func (mss *MapStringString) Scan(v interface{}) error    { return json.Unmarshal(utils.AnyToBytes(v), mss) }
