package main

import (
	"fmt"
	"net/http"

	"nopricey/internal"
	"nopricey/internal/api"
	"nopricey/internal/db"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Testing")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	database, err := db.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	internal.StartDailyCron(func() {
		internal.FetchAndStoreAllOffers(database)
	})

	h := api.Handler{DB: database}
	mux := api.RegisterRoutes(&h)

	http.ListenAndServe(":10000", http.Handler(mux))
}

func main() {
	handleRequests()
}
