package fetch

// NOTE: PROBLEM: It is in norwegian
// NOTE: POTENTIAL SOLUTION: Use norwegian as searching language
// NOTE: POTENTIAL ALGO:
// 1. Potentially save the response to a localdb
// 2. All data combined into a list
// 3. Take inputs 1 by 1 from the user provided list
// 4. With some function check if the input is in any of the json
// 5. If only 1 found, return that if multiple decide based on price OR price/kg/L
// 6. If none no return

// Bunnpris's Tjek dealer ID, found via /v2/dealers/search?query=Bunnpris.
const bunnprisDealerID = "5b11sm"

func BunnprisWeekly() ([]TjekOffer, error) {
	return fetchTjekOffers(bunnprisDealerID)
}
