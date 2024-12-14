package availabletokens

import (
	"binance-server/connection"
	"binance-server/fapi"
	"sync"
)

type AvailableTokenList struct {
	tokens     []fapi.FetchedToken
	mp         map[string]bool
	Websockets *connection.WebSockets
	mu         sync.Mutex
}

func NewAvailableTokenList(_tokens []fapi.FetchedToken) *AvailableTokenList {
	list := AvailableTokenList{tokens: _tokens, mp: make(map[string]bool), mu: sync.Mutex{}, Websockets: connection.NewWebSockets()}
	for _, token := range list.tokens {
		list.mp[token.Symbol] = true
		list.Websockets.AddWebsocket(token.Symbol, connection.NewWebsocket(token.Symbol))
	}
	return &list
}

func (list *AvailableTokenList) IsAvailable(symbol string) bool {
	_, exists := list.mp[symbol]
	return exists
}
