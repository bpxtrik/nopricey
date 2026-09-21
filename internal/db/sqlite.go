package db

import (
	"database/sql"
	"fmt"
	"slices"
	"sort"

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

// offerKey makes rows comparable
func offerKey(heading, description string, price, prePrice float64, currency, unitSymbol string, sizeFrom, sizeTo float64, runFrom, runTill string) string {
	return fmt.Sprintf("%s|%s|%.4f|%.4f|%s|%s|%.4f|%.4f|%s|%s",
		heading, description, price, prePrice, currency, unitSymbol, sizeFrom, sizeTo, runFrom, runTill)
}

// offersChanged reports whether offers differs from whats currently stored
// in current_offers for store. Catalogs are published weekly, so most daily
// calls into InsertNewOffers should find nothing changed.
func offersChanged(db *sql.DB, store string, offers []fetch.TjekOffer) (bool, error) {
	rows, err := db.Query(`
		SELECT heading, description, price, pre_price, currency, unit_symbol, size_from, size_to, run_from, run_till
		FROM current_offers WHERE store = ?
	`, store)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var current []string
	for rows.Next() {
		var heading, description, currency, unitSymbol, runFrom, runTill string
		var price, prePrice, sizeFrom, sizeTo float64

		if err := rows.Scan(&heading, &description, &price, &prePrice, &currency, &unitSymbol, &sizeFrom, &sizeTo, &runFrom, &runTill); err != nil {
			return false, err
		}

		current = append(current, offerKey(heading, description, price, prePrice, currency, unitSymbol, sizeFrom, sizeTo, runFrom, runTill))
	}
	if err := rows.Err(); err != nil {
		return false, err
	}

	fresh := make([]string, 0, len(offers))
	for _, o := range offers {
		fresh = append(fresh, offerKey(o.Heading, o.Description, o.Pricing.Price, o.Pricing.PrePrice,
			o.Pricing.Currency, o.Quantity.Unit.Symbol, o.Quantity.Size.From, o.Quantity.Size.To, o.RunFrom, o.RunTill))
	}

	sort.Strings(current)
	sort.Strings(fresh)

	return !slices.Equal(current, fresh), nil
}

// InsertNewOffers replaces the contents of `current_offers` with offers for
// store, but only if offers actually differs from whats already stored 
// so calling this daily is a no-op on the days the catalog hasnt changed.
// When it does write, whatever was in `current_offers` is archived into
// `old_offers` first, so briefly both tables hold the same data, before
// its cleared and the new offers are inserted. Runs as a single
// transaction so a failure midway leaves current_offers untouched rather
// than half-updated.
func InsertNewOffers(db *sql.DB, store string, offers []fetch.TjekOffer) error {
	changed, err := offersChanged(db, store, offers)
	if err != nil {
		return err
	}
	if !changed {
		fmt.Println(store, "offers unchanged, skipping update")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO old_offers (store, heading, description, price, pre_price, currency, unit_symbol, size_from, size_to, run_from, run_till, fetched_at)
		SELECT store, heading, description, price, pre_price, currency, unit_symbol, size_from, size_to, run_from, run_till, fetched_at
		FROM current_offers WHERE store = ?
	`, store)
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM current_offers WHERE store = ?", store); err != nil {
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
