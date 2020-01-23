package migrations

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Migration func() string

type Interface interface {
	Schema() map[string]string
	CreateMigrationTable(sqlx.Execer) error
	AddMigration(sqlx.Execer, string) error
	Migrated(sqlx.Queryer, string) (bool, error)
}

type SQLMigration struct {
}

func (*SQLMigration) Schema() map[string]string {
	return map[string]string{
		"create-quotes-table": `
			CREATE TABLE quotes (
				  "quote_id"           VARCHAR(100)
				, "maker_asset_address" VARCHAR(100)
				, "taker_asset_address" VARCHAR(100)
				, "maker_asset_size"   TEXT
				, "taker_asset_size"   TEXT
				, "expiration"         INTEGER
				, "server_time"        INTEGER
				, "order_hash"         VARCHAR(100)
				, "zero_ex_transaction_hash" VARCHAR(100)
				, PRIMARY KEY (quote_id)
			)`,

		"create-orders-table": `
			CREATE TABLE signed_orders (
				 quote_id VARCHAR(100)
			   , order_bytes TEXT
			   , PRIMARY KEY(quote_id)
			)`,

		"create-markets-table": `
			CREATE TABLE markets (
				  "maker_asset_address" VARCHAR(10)
				, "taker_asset_addresses" TEXT
				, "trade_info"          TEXT
				, "quote_info"          TEXT
				, "metadata"            TEXT
				, PRIMARY KEY (maker_asset_address)
			)
		`,
	}
}

func (*SQLMigration) CreateMigrationTable(db sqlx.Execer) error {
	q := `
	CREATE TABLE IF NOT EXISTS migrations(
		  name VARCHAR(255)
		, UNIQUE(name)
	)
	`
	_, err := db.Exec(q)
	return err
}

func (*SQLMigration) AddMigration(db sqlx.Execer, name string) error {
	q := `
	INSERT INTO migrations (name) VALUES ($1);
	`
	_, err := db.Exec(q, name)
	return err
}

func (*SQLMigration) Migrated(db sqlx.Queryer, name string) (bool, error) {
	q := `
	SELECT COUNT(*) FROM migrations WHERE name = $1
	`
	row := db.QueryRowx(q, name)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func Run(db *sqlx.DB, m Interface) error {
	if err := m.CreateMigrationTable(db); err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	// After a call to Commit or Rollback, all operations on the transaction fail with ErrTxDone.
	defer tx.Rollback() // nolint:errcheck

	for name, stmt := range m.Schema() {
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
