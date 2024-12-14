package subscription

import (
	"binance-server/availabletokens"
)

type Subscription struct {
	Pairs []string `json:"pairs"`
	Mp    map[string]bool
	Quits map[string]chan struct{}
	List  *availabletokens.AvailableTokenList
}
