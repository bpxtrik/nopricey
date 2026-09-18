package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TjekOffer is the structured weekly-discount ("tilbud") shape returned by
// Tjek's public eTilbudsavis API (api.etilbudsavis.dk, no auth required).
// Several Norwegian grocery chains publish their weekly catalogs through
// this same platform (confirmed dealer IDs live in each store's own file),
// so one fetcher covers all of them 
type TjekOffer struct {
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Pricing     struct {
		Price    float64 `json:"price"`
		PrePrice float64 `json:"pre_price"`
		Currency string  `json:"currency"`
	} `json:"pricing"`
	Quantity struct {
		Unit struct {
			Symbol string `json:"symbol"`
		} `json:"unit"`
		Size struct {
			From float64 `json:"from"`
			To   float64 `json:"to"`
		} `json:"size"`
	} `json:"quantity"`
	RunFrom string `json:"run_from"`
	RunTill string `json:"run_till"`
}

// The offers endpoint paginates with offset/limit, not page/p.
// limit is capped server-side at 100
// (anything higher returns HTTP 400 PAGINATION_INVALID_LIMIT), so a single
// request can't get everything, keep paging with offset until a short page
// tells us we've reached the end.
const tjekPageSize = 100

func fetchTjekOffers(dealerID string) ([]TjekOffer, error) {
	var all []TjekOffer

	for offset := 0; ; offset += tjekPageSize {
		url := fmt.Sprintf(
			"https://api.etilbudsavis.dk/v2/offers?dealer_ids=%s&order_by=-expires&limit=%d&offset=%d",
			dealerID, tjekPageSize, offset,
		)

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("tjek offers request failed for dealer %s: status %d", dealerID, resp.StatusCode)
		}

		var page []TjekOffer
		err = json.NewDecoder(resp.Body).Decode(&page)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < tjekPageSize {
			break
		}
	}

	return all, nil
}
