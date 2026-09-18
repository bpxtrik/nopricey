package fetch

// id via: /v2/dealers/search?query=Rema&r_lat=59.9139&r_lng=10.7522&r_radius=50000 —
// the plain query search returns the Danish REMA 1000 dealer instead, so
// the geolocation params are required to resolve the Norwegian one.
const rema1000DealerID = "faa0Ym"

func Rema1000Weekly() ([]TjekOffer, error) {
	return fetchTjekOffers(rema1000DealerID)
}
