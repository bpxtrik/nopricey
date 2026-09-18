package api

import "net/http"

func RegisterRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /item", h.SingleItem)

	return mux
}
