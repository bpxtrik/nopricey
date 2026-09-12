package api

import (
	"encoding/json"
	"io"
	"net/http"

	"nopricey/internal"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *pgxpool.Pool
}

func writeJSON(w *http.ResponseWriter, status int, body any) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(status)
	json.NewEncoder(*w).Encode(body)
}

func (h Handler) BunnprisWeekly(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("https://api.etilbudsavis.dk/v2/offers?dealer_ids=5b11sm&order_by=-expires")
	if err != nil {
		writeJSON(&w, http.StatusBadGateway, internal.Response{Detail: err.Error()})
		return
	}
	// resp.Body is a long list of articles. Based on this it is possible to search and filter thru it.
	// NOTE: PROBLEM: It is in norwegian
	// NOTE: POTENTIAL SOLUTION: Use norwegian as searching language
	// NOTE: POTENTIAL ALGO: 
	// 1. Potentially save the response to a localdb
	// 2. All data combined into a list
	// 3. Take inputs 1 by 1 from the user provided list
	// 4. With some function check if the input is in any of the json
	// 5. If only 1 found, return that if multiple decide based on price OR price/kg/L
	// 6. If none no return
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
