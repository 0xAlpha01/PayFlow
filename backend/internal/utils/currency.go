package utils

import (
	"fmt"
	"strings"
)

// FormatNaira formats a monetary amount (stored as a plain numeric major-unit
// value, e.g. 50000) as a Nigerian Naira display string, e.g. "₦50,000.00".
func FormatNaira(amount float64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}
	whole := int64(amount)
	cents := int64((amount-float64(whole))*100 + 0.5)
	if cents >= 100 {
		whole++
		cents -= 100
	}

	digits := fmt.Sprintf("%d", whole)
	var grouped strings.Builder
	for i, d := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			grouped.WriteByte(',')
		}
		grouped.WriteRune(d)
	}

	return fmt.Sprintf("%s₦%s.%02d", sign, grouped.String(), cents)
}
