package tools

import (
	"fmt"
	"math"
)

// FormatMilliunits converts YNAB milliunits to a dollar string (e.g., "$12.50").
func FormatMilliunits(m int64) string {
	dollars := float64(m) / 1000.0
	if dollars < 0 {
		return fmt.Sprintf("-$%.2f", math.Abs(dollars))
	}
	return fmt.Sprintf("$%.2f", dollars)
}

// MilliunitsFromDollars converts a dollar amount to YNAB milliunits.
func MilliunitsFromDollars(d float64) int64 {
	return int64(math.Round(d * 1000))
}
