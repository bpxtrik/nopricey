package api

import "net/http"

func RegisterRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /bunnpris/weekly", h.BunnprisWeekly)

	return mux
}
