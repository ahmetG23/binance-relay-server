/**
 * @api {get} /ws Start WebSocket server
 * @apiName StartWebSocket
 * @apiGroup WebSocket
 * @apiDescription Start a WebSocket server that listens for incoming connections and reads subscriptions from clients.
 */

package main

import (
	"binance-server/availabletokens"
	"binance-server/fapi"
	"binance-server/schema"
	"binance-server/subscription"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/websocket"
)

var tokens []fapi.FetchedToken

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func readSubscriptions(conn *websocket.Conn, list *availabletokens.AvailableTokenList, ch chan *subscription.Subscription) {
	defer close(ch)

	// Set a close handler
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Client connection is closed with code: %d, text: %s", code, text)
		// Close all websockets in the list
		for _, ws := range list.Websockets.GetAllWebsockets() {
			close(ws.Quit)
		}
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			} else {
				log.Println("Read error:", err)
			}
			return
		}
		log.Println("Received subscription:", string(message))

		sub := &subscription.Subscription{
			Mp:    make(map[string]bool),
			Quits: make(map[string]chan struct{}),
			List:  list,
		}

		err = json.Unmarshal(message, &sub)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// make all pairs lower case
		for i, pair := range sub.Pairs {
			sub.Pairs[i] = strings.ToLower(pair)
		}

		var onlyAvailable []string
		for _, pair := range sub.Pairs {
			if list.IsAvailable(pair) {
				onlyAvailable = append(onlyAvailable, pair)
			} else {
				log.Println("Pair", pair, "is not available")
			}
		}
		sub.Pairs = onlyAvailable
		ch <- sub
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	list := availabletokens.NewAvailableTokenList(tokens)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic in handleWebSocket: %v", r)
        }
        conn.Close()
    }()

	subCh := make(chan *subscription.Subscription)
	tickerCh := make(chan schema.BookTicker, 100) // Buffered channel to improve performance
	quit := make(chan struct{})

	sub := &subscription.Subscription{
		Mp:    make(map[string]bool),
		Quits: make(map[string]chan struct{}),
		List:  list,
	}

	go readSubscriptions(conn, list, subCh)

	for {
		select {
		case newSub, ok := <-subCh:
			if !ok {
				return
			}
			oldPairs := sub.Pairs
			sub.Pairs = newSub.Pairs
			sub.Read(tickerCh, oldPairs)

		case bookTicker := <-tickerCh:
			jsonData, err := json.Marshal(bookTicker)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				close(quit)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Write error:", err)
				close(quit)
				return
			}
		case <-quit:
			return
		}
	}
}

func main() {
	tokens = fapi.Fetch()
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}