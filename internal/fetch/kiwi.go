package fetch

// id via: /v2/dealers/search?query=Kiwi&r_lat=59.9139&r_lng=10.7522&r_radius=50000.
const kiwiDealerID = "257bxm"

func KiwiWeekly() ([]TjekOffer, error) {
	return fetchTjekOffers(kiwiDealerID)
}
