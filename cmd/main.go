package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Testing")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":10000", nil)
}

func main() {
	handleRequests()
}
