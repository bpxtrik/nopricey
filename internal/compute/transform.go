package compute

import (
	"strconv"
	"strings"
)

// Returns the price of a kilo unit (kg, liter, etc.). Uses string formatting, converts grams and mililiters to kg and liter respectively. Then calculates price. No error return, if price is 0 it will not be displayed
func ConvertToPricePerKiloUnit(price float64, qty string) float64 {
	parts := strings.Fields(qty)
	unit := parts[len(parts)-1]
	value, err := strconv.ParseFloat(parts[0], 64)

	if err != nil {
		return 0
	}

	switch unit {
	case "kg", "l":
		return price / value
	case "g", "ml":
		return (price / value) * 1000
	default:
		  return 0
	}
}
