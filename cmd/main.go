package main

import (
	"fmt"
	"net/http"

	"nopricey/internal/api"
	"nopricey/internal/db"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Testing")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	db, err := db.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	h := api.Handler{DB: db}
	mux := api.RegisterRoutes(&h)

	http.ListenAndServe(":10000", http.Handler(mux))
}

func main() {
	handleRequests()
}
