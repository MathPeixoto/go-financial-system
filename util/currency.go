package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	BRL = "BRL"
)

// Currency is a type for currencies
type Currency string

// IsSupportedCurrency checks if the currency is valid
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, BRL:
		return true
	}
	return false
}
