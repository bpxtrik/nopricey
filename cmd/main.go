package main

import (
	"fmt"
	"net/http"

	"nopricey/internal/api"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Testing")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	h := api.Handler{DB: nil}
	mux := api.RegisterRoutes(&h)

	http.ListenAndServe(":10000", http.Handler(mux))
}

func main() {
	handleRequests()
}
