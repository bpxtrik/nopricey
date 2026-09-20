package api

import (
	"encoding/json"
	"net/http"

	"nopricey/internal"
	"nopricey/internal/compute"
	"nopricey/internal/db"

	"database/sql"
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
	qr := compute.FindOffers(query)

	if len(qr) == 0 {
		writeJSON(w, 200, internal.Response{Detail: "No product found!"})
		return
	}
	var results []internal.Offer
	for store, offers := range qr {
		// fmt.Printf("%s:\n", store)
		for _, o := range offers {
			qty := compute.FormatQuantity(o)
			// fmt.Printf("  %s (%s) - %.2f %s\n", o.Heading, qty, o.Pricing.Price, o.Pricing.Currency)
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

	db.InsertNewOffers(h.DB, "Bunnpris", qr["Bunnpris"])

	writeJSON(w, 200, results)
}
