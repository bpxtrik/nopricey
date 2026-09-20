package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./my.db")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer db.Close()
	fmt.Println("Connected to DB")

	_, err = CreateTable(db)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Table created")


	return db, nil
}

func CreateTable(db *sql.DB) (sql.Result, error) {
	schema := `CREATE TABLE IF NOT EXISTS offers (
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
	)`

	return db.Exec(schema)
}


