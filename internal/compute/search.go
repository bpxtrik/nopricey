package compute

import (
	"fmt"
	"log"
	"nopricey/internal/fetch"
	"strings"
)

// Stores maps a store name to its weekly-offer fetcher, covering every
// store currently wired up.
var Stores = map[string]func() ([]fetch.TjekOffer, error){
	"Bunnpris": fetch.BunnprisWeekly,
	"Rema1000": fetch.Rema1000Weekly,
	"Kiwi":     fetch.KiwiWeekly,
	"Extra": 	fetch.ExtraWeekly,
}

// FindOffers looks up query (matched case-insensitively against each
// offer's heading/description) across every store in Stores, returning the
// matching offers keyed by store name.
func FindOffers(query string) map[string][]fetch.TjekOffer {
	query = strings.ToLower(query)
	results := make(map[string][]fetch.TjekOffer)

	for store, fetch := range Stores {
		offers, err := fetch()
		if err != nil {
			log.Println(store, "error:", err)
			continue
		}
		for _, o := range offers {
			if strings.Contains(strings.ToLower(o.Heading), query) || strings.Contains(strings.ToLower(o.Description), query) {
				results[store] = append(results[store], o)
			}
		}
	}

	return results
}

// FormatQuantity renders an offer's size/unit, e.g. "400 g" or "0.5 l", from
// the API's structured quantity field rather than parsing the freeform
// description text.
func FormatQuantity(o fetch.TjekOffer) string {
	if o.Quantity.Unit.Symbol == "" {
		return ""
	}
	if o.Quantity.Size.To > o.Quantity.Size.From {
		return fmt.Sprintf("%g-%g %s", o.Quantity.Size.From, o.Quantity.Size.To, o.Quantity.Unit.Symbol)
	}
	return fmt.Sprintf("%g %s", o.Quantity.Size.From, o.Quantity.Unit.Symbol)
}
