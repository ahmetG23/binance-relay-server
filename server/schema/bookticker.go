package schema

type BookTicker struct {
	Symbol       string  `json:"s"`
	BestBidPrice string `json:"b"`
	BestAskPrice string `json:"a"`
}
