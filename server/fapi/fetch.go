package fapi

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func Fetch() ([]FetchedToken) {
	url := "https://fapi.binance.com/fapi/v1/ticker/24hr"

	// Fetch data from Binance
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Fetch error: ", err)
	}
	defer resp.Body.Close()

	// Decode JSON response
	var tokens []FetchedToken
	err = json.NewDecoder(resp.Body).Decode(&tokens)
	if err != nil {
		log.Fatal("Decode error: ", err)
	}

	// filter out tokens that are not USDT
	var usdtTokens []FetchedToken
	for _, token := range tokens {
		if strings.Contains(token.Symbol, "USDT") {
			token.Symbol = strings.ToLower(token.Symbol)
			usdtTokens = append(usdtTokens, token)
		}
	}

	sort.Slice(usdtTokens, func(i, j int) bool {
		qv1,err1 := strconv.ParseFloat(usdtTokens[i].QuoteVolume, 64)
		qv2,err2 := strconv.ParseFloat(usdtTokens[j].QuoteVolume, 64)
		if err1 != nil || err2 != nil {
			log.Fatal("Error parsing float")
			log.Panic(err1)
		}
		return qv1 > qv2
	})

	// take only largest 150 tokens
	usdtTokens = usdtTokens[:150]
	log.Println("Fetched", len(usdtTokens), "tokens")
	return usdtTokens
}