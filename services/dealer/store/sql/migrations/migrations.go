package migrations

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Migration func() string

type Migrations struct {
}

func (m *Migrations) CreateTradesTable() string {
	return `
	CREATE TABLE trades (
		  "quote_id"           VARCHAR(100)
		, "market_id"          VARCHAR(100)
		, "order_hash"         VARCHAR(100)
		, "transaction_hash"   VARCHAR(100)
		, "taker_address"      VARCHAR(100)
		, "timestamp"          INTEGER
		, "maker_asset_ticker" VARCHAR(10)
		, "taker_asset_ticker" VARCHAR(10)
		, "maker_asset_amount" BIGINT
		, "taker_asset_amount" BIGINT
		, PRIMARY KEY (quote_id)
	)`
}

func (m *Migrations) CreateQuotesTable() string {
	return `
	CREATE TABLE quotes (
		  "quote_id"           VARCHAR(100)
		, "maker_asset_ticker" VARCHAR(100)
		, "taker_asset_ticker" VARCHAR(100)
		, "maker_asset_size"   VARCHAR(100)
		, "quote_asset_size"   VARCHAR(100)
		, "order_hash"         VARCHAR(100)
		, "order"              TEXT
		, "fill_tx"            VARCHAR(100)
		, PRIMARY KEY (quote_id)
	)`
}

func (m *Migrations) CreateMigrationTable() string {
	return `
	CREATE TABLE IF NOT EXISTS migrations(
		  name VARCHAR(255)
		, UNIQUE(name)
	)
	`
}

func (m *Migrations) AddMigration(tx *sql.Tx, name string) error {
	q := `
	INSERT INTO migrations (name) VALUES ($1);
	`
	_, err := tx.Exec(q, name)
	return err
}

func (m *Migrations) Migrated(tx *sql.Tx, name string) (bool, error) {
	q := `
	SELECT COUNT(*) FROM migrations WHERE name = $1
	`
	row := tx.QueryRow(q, name)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *Migrations) Run(db *sqlx.DB) error {
	if _, err := db.Exec(m.CreateMigrationTable()); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// After a call to Commit or Rollback, all operations on the transaction fail with ErrTxDone.
	defer tx.Rollback()

	mz := []struct {
		name string
		fn   Migration
	}{
		{"create_trades_table", m.CreateTradesTable},
		{"create_quotes_table", m.CreateQuotesTable},
	}

	for _, op := range mz {
		name := op.name
		stmt := op.fn()

		if ok, err := m.Migrated(tx, name); err != nil {
			return err
		} else if ok {
			continue
		}

		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("%w. query:%s", err, stmt)
		}

		if err := m.AddMigration(tx, name); err != nil {
			return err
		}
	}

	// deferd Rollback won"t have effect after commit
	return tx.Commit()
}
