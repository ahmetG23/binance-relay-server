package server

import (
	"encoding/json"
	"log"
)

type Subscription struct {
	Method string `json:"method"`
	Pairs []string `json:"pairs"`
}


func NewSubscription(pairs []string) *Subscription {
	return &Subscription{
		Method: "SUBSCRIBE",
		Pairs: pairs,
	}
}

func (s *Subscription) GetJson() []byte {
	json, err := json.Marshal(s)
	if err != nil {
		log.Fatal("Failed to marshal subscription message:", err)
	}
	return json
}