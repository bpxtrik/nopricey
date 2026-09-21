package internal

import (
	"database/sql"
	"log"

	"nopricey/internal/compute"
	"nopricey/internal/db"
)

// FetchAndStoreAllOffers fetches every store's weekly offers and hands them
// to db.InsertNewOffers. That call only touches current_offers/old_offers
// when the fetched offers actually differ from what's already stored, so
// it's safe to run this daily even though the underlying catalogs only
// change weekly.
func FetchAndStoreAllOffers(database *sql.DB) {
	for store, fetchOffers := range compute.Stores {
		offers, err := fetchOffers()
		if err != nil {
			log.Println(store, "fetch error:", err.Error())
			continue
		}
		log.Printf("%s: fetched %d offers", store, len(offers))

		if err := db.InsertNewOffers(database, store, offers); err != nil {
			log.Println(store, "insert error:", err.Error())
		}
	}
}
