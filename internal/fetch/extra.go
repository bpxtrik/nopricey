package fetch

// id from: /v2/dealers/search?query=Extra&r_lat=59.9139&r_lng=10.7522&r_radius=50000.
const extraDealerId = "80742m"

func ExtraWeekly() ([]TjekOffer, error) {
	return fetchTjekOffers(extraDealerId)
}
