package util

import (
	"golang.org/x/text/currency"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvyxwz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer betwwen min and mx
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random currency
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{
		currency.EUR.String(),
		currency.USD.String(),
		currency.BRL.String(),
	}

	n := len(currencies)
	return currencies[rand.Intn(n)]
}