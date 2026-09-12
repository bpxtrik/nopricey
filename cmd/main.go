package main

import (
	"fmt"
	"net/http"
	"nopricey/internal/api"
	"nopricey/internal/compute"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Testing")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	h := api.Handler{DB: nil}
	mux := api.RegisterRoutes(&h)

	resp, err := compute.BunnprisWeekly()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(len(resp))

	http.ListenAndServe(":10000", http.Handler(mux))
}

func main() {
	handleRequests()
}
