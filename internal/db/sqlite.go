package db

import (
	"database/sql"
	"fmt"

	"nopricey/internal/fetch"

	_ "modernc.org/sqlite"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./my.db")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Connected to DB")

	_, err = CreateTable(db)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return db, nil
}

func CreateTable(db *sql.DB) (sql.Result, error) {
	schema := `CREATE TABLE IF NOT EXISTS current_offers (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		store          TEXT NOT NULL,
		heading        TEXT NOT NULL,
		description    TEXT NOT NULL,
		price          REAL NOT NULL,
		pre_price      REAL NOT NULL,
		currency       TEXT NOT NULL,
		unit_symbol    TEXT NOT NULL,
		size_from      REAL NOT NULL,
		size_to        REAL NOT NULL,
		run_from       TEXT NOT NULL,
		run_till       TEXT NOT NULL,
		fetched_at     TEXT NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS old_offers (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		store          TEXT NOT NULL,
		heading        TEXT NOT NULL,
		description    TEXT NOT NULL,
		price          REAL NOT NULL,
		pre_price      REAL NOT NULL,
		currency       TEXT NOT NULL,
		unit_symbol    TEXT NOT NULL,
		size_from      REAL NOT NULL,
		size_to        REAL NOT NULL,
		run_from       TEXT NOT NULL,
		run_till       TEXT NOT NULL,
		fetched_at     TEXT NOT NULL DEFAULT (datetime('now'))
	)
	`

	return db.Exec(schema)
}

// InsertNewOffers replaces the contents of `current_offers` with offers for
// store. Whatever was in `current_offers` is archived into `old_offers`
// first, so briefly both tables hold the same data, before it's cleared and
// the new offers are inserted. Runs as a single transaction so a failure
// midway leaves current_offers untouched rather than half-updated.
func InsertNewOffers(db *sql.DB, store string, offers []fetch.TjekOffer) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO old_offers (store, heading, description, price, pre_price, currency, unit_symbol, size_from, size_to, run_from, run_till, fetched_at)
		SELECT store, heading, description, price, pre_price, currency, unit_symbol, size_from, size_to, run_from, run_till, fetched_at
		FROM current_offers
	`)
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM current_offers"); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO current_offers (store, heading, description, price, pre_price, currency, unit_symbol, size_from, size_to, run_from, run_till)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, o := range offers {
		_, err := stmt.Exec(store, o.Heading, o.Description, o.Pricing.Price, o.Pricing.PrePrice,
			o.Pricing.Currency, o.Quantity.Unit.Symbol, o.Quantity.Size.From, o.Quantity.Size.To,
			o.RunFrom, o.RunTill)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
