package migrations

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Name string
	Stmt string
}

type Interface interface {
	Schema() []Migration
	CreateMigrationTable(sqlx.Execer) error
	AddMigration(sqlx.Execer, string) error
	Migrated(sqlx.Queryer, string) (bool, error)
}

type SQLMigration struct {
}

func (*SQLMigration) Schema() []Migration {
	return []Migration{
		{
			Name: "create-quotes-table",
			Stmt: `CREATE TABLE quotes (
				  "quote_id"            VARCHAR(100)
				, "maker_asset_address" VARCHAR(100)
				, "taker_asset_address" VARCHAR(100)
				, "maker_asset_size"    TEXT
				, "taker_asset_size"    TEXT
				, "expiration"          NUMERIC
				, "server_time"         NUMERIC
				, "zero_ex_transaction_hash" VARCHAR(100)
				, PRIMARY KEY (quote_id)
			)`,
		},
		{
			Name: "create-transaction-infos-table",
			Stmt: `CREATE TABLE transaction_infos (
				 quote_id VARCHAR(100)
			   , transaction_info_bytes TEXT
			   , PRIMARY KEY(quote_id)
			)`,
		},
		{
			Name: "create-markets-table",
			Stmt: `CREATE TABLE markets (
				  "maker_asset_address" VARCHAR(10)
				, "taker_asset_addresses" TEXT
				, "trade_info"          TEXT
				, "quote_info"          TEXT
				, "metadata"            TEXT
				, PRIMARY KEY (maker_asset_address)
			)
		`,
		},
		{
			Name: "create-policy-table",
			Stmt: `CREATE TABLE policies (
				  "entry" VARCHAR(100)
				, PRIMARY KEY(entry)
			)
		`,
		},
		{
			Name: "create-trades-table",
			Stmt: `CREATE TABLE trades (
					  "quote_id" VARCHAR(100) REFERENCES quotes ON DELETE CASCADE
					, "tx_timestamp" NUMERIC
					, "tx_hash" VARCHAR(66)
					, "status" SMALLINT
				)
			`,
		},
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
	// After a call to Commit or Rollback, all operations on the transaction will fail with ErrTxDone
	defer tx.Rollback() // nolint:errcheck

	for _, step := range m.Schema() {
		if ok, err := m.Migrated(tx, step.Name); err != nil {
			return err
		} else if ok {
			continue
		}

		if _, err := tx.Exec(step.Stmt); err != nil {
			return fmt.Errorf("%w. query:%s", err, step.Stmt)
		}

		if err := m.AddMigration(tx, step.Name); err != nil {
			return err
		}
	}

	// defered Rollback won"t have effect after commit
	return tx.Commit()
}
