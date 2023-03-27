package utils

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "qwertyuiopasdfghjklzxcvbnm"
)

var (
	currencyList = []string{"USD", "EUR", "INR"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(26)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomMoney() int64 {
	return randomInt(0, 1000)
}
func RandomOwner() string {
	return randomString(6)
}
func RandomCurrencies() string {
	rg := rand.Intn(len(currencyList))
	return currencyList[rg]
}
