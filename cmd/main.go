package main

import (
	"log"
	"net/http"

	"nopricey/internal"
	"nopricey/internal/api"
	"nopricey/internal/db"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Testing")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	database, err := db.New()
	if err != nil {
		log.Println(err.Error())
		return
	}

	hasOffers, err := db.HasOffers(database)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if !hasOffers {
		log.Println("current_offers is empty, running initial fetch")
		internal.FetchAndStoreAllOffers(database)
	}

	internal.StartDailyCron(func() {
		internal.FetchAndStoreAllOffers(database)
	})

	h := api.Handler{DB: database}
	mux := api.RegisterRoutes(&h)

	http.ListenAndServe(":10000", api.LoggingMiddleware(mux))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	handleRequests()
}
