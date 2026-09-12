package compute

import (
	"encoding/json"
	"net/http"
)

// NOTE: PROBLEM: It is in norwegian
// NOTE: POTENTIAL SOLUTION: Use norwegian as searching language
// NOTE: POTENTIAL ALGO:
// 1. Potentially save the response to a localdb
// 2. All data combined into a list
// 3. Take inputs 1 by 1 from the user provided list
// 4. With some function check if the input is in any of the json
// 5. If only 1 found, return that if multiple decide based on price OR price/kg/L
// 6. If none no return

type BunnprisOffer struct {
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

func BunnprisWeekly() ([]BunnprisOffer, error) {
	resp, err := http.Get("https://api.etilbudsavis.dk/v2/offers?dealer_ids=5b11sm&order_by=-expires")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var offers []BunnprisOffer
	err = json.NewDecoder(resp.Body).Decode(&offers)
	if err != nil {
		return nil, err
	}

	return offers, nil
}
