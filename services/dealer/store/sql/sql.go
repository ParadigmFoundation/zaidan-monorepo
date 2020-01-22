package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
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

func (s *Store) CreateTrade(t *types.Trade) error {
	stmt := `
		INSERT INTO trades VALUES (
			:quote_id,
			:market_id,
			:order_hash,
			:transaction_hash,
			:taker_address,
			:timestamp,
			:maker_asset_address,
			:taker_asset_address,
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

func (s *Store) CreateAsset(a *types.Asset) error {
	_, err := s.db.Exec(
		`INSERT INTO assets VALUES($1, $2, $3, $4, $5)`,
		a.Ticker,
		a.Name,
		a.Decimals,
		a.NetworkId,
		a.Address,
	)
	return err
}

func (s *Store) GetAssetBy(key, value string) (*types.Asset, error) {
	stmt := "SELECT * FROM assets WHERE " + key + " = $1 LIMIT 1"
	asset := types.Asset{}
	if err := s.db.Get(&asset, stmt, value); err != nil {
		return nil, err
	}
	return &asset, nil

}

func (s *Store) GetAsset(ticker string) (*types.Asset, error) {
	return s.GetAssetBy("ticker", ticker)
}

func (s *Store) GetAssetByAddress(addr string) (*types.Asset, error) {
	return s.GetAssetBy("address", addr)
}

func (s *Store) CreateMarket(mkt *types.Market) error {
	mkt.Id = uuid.New().String()
	_, err := s.db.Exec(
		`INSERT INTO markets VALUES($1, $2, $3, $4, $5, $6)`,
		mkt.Id,
		mkt.MakerAssetTicker,
		StringSlice(mkt.TakerAssetTickers),
		mkt.TradeInfo,
		mkt.QuoteInfo,
		MapStringString(mkt.Metadata),
	)
	return err
}

func (s *Store) GetMarket(id string) (*types.Market, error) {
	stmt := `SELECT * FROM markets where id = $1 LIMIT 1`
	mkt := types.Market{}
	var tickers StringSlice
	var metadata MapStringString

	err := s.db.QueryRow(stmt, id).Scan(
		&mkt.Id,
		&mkt.MakerAssetTicker,
		&tickers,
		&mkt.TradeInfo,
		&mkt.QuoteInfo,
		&metadata,
	)
	if err != nil {
		return nil, err
	}
	mkt.TakerAssetTickers = tickers
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

	stmt := `SELECT COUNT(*) FROM policies WHERE entry = ?`
	if err := s.db.QueryRow(stmt, t).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

type StringSlice []string

func (ss StringSlice) Value() (driver.Value, error) { return json.Marshal(ss) }
func (ss *StringSlice) Scan(v interface{}) error    { return json.Unmarshal(utils.AnyToBytes(v), ss) }

type MapStringString map[string]string

func (mss MapStringString) Value() (driver.Value, error) { return json.Marshal(mss) }
func (mss *MapStringString) Scan(v interface{}) error    { return json.Unmarshal(utils.AnyToBytes(v), mss) }
