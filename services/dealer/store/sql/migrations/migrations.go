package migrations

import (
	"database/sql"
	"reflect"

	"github.com/jmoiron/sqlx"
)

type Migration func() string

type Migrations struct {
}

func (m *Migrations) CreateTradesTable() string {
	return `
	CREATE TABLE trades (
		  quote_id  VARCHAR(100)
		, market_id VARCHAR(100)
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
	INSERT INTO migrations (name) VALUES (?)
	`
	_, err := tx.Exec(q, name)
	return err
}

func (m *Migrations) Migrated(tx *sql.Tx, name string) (bool, error) {
	q := `
	SELECT COUNT(*) FROM migrations WHERE name = ?
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

	mz := []Migration{
		m.CreateTradesTable,
	}

	for _, fn := range mz {
		name := reflect.TypeOf(fn).String()
		stmt := fn()

		if ok, err := m.Migrated(tx, name); err != nil {
			return err
		} else if ok {
			continue
		}

		if _, err := tx.Exec(stmt); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		m.AddMigration(tx, name)
	}

	return tx.Commit()
}
