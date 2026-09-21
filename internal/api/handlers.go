package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"nopricey/internal"
	"nopricey/internal/compute"
	"nopricey/internal/db"
)

type Handler struct {
	DB *sql.DB
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func (h *Handler) SingleItem(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	if query == "" {
		writeJSON(w, 400, internal.Response{Detail: "Query param missing"})
		return
	}
	all, err := db.GetCurrentOffers(h.DB)
	if err != nil {
		log.Println("GetCurrentOffers error:", err.Error())
		writeJSON(w, 500, internal.Response{Detail: "Internal error"})
		return
	}

	qr := compute.FilterOffers(all, query)

	if len(qr) == 0 {
		writeJSON(w, 200, internal.Response{Detail: "No product found!"})
		return
	}
	var results []internal.Offer
	for store, offers := range qr {
		for _, o := range offers {
			qty := compute.FormatQuantity(o)
			results = append(results, internal.Offer{
				Store:       store,
				ProductName: o.Heading,
				ProductDesc: o.Description,
				Price:       o.Pricing.Price,
				Quantity:    qty,
				PricePKU:    compute.ConvertToPricePerKiloUnit(o.Pricing.Price, qty),
			})

		}
	}

	writeJSON(w, 200, results)
}
