package original

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/**
 * call Init function with pairs, print and timeoutSec
 * It connects directly to the binance's server
 * pairs: []string - list of pairs to subscribe
 * print: bool - print the received messages
 * timeoutSec: int - timeout in seconds
*/
func Init(pairs []string, print bool, timeoutSec int) {
	firstPrinted := false
	startTime := time.Now()

	wsURL := "wss://fstream.binance.com/stream?streams="

	for i, pair := range pairs {
		if i > 0 {
			wsURL += "/"
		}
		wsURL += strings.ToLower(pair) + "@bookTicker"
	}
	log.Println(wsURL)

	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatal("Failed to parse the WebSocket URL:", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	// Goroutine to read messages from WebSocket
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Read error:", err)
				}
				return 
			}
			if !firstPrinted {
				log.Println("ORIGINAL first receive time elapse:", time.Since(startTime))
				firstPrinted = true
			}
			if print {
				log.Printf("ORIGINAL Received: %s\n", message)
			}
		}
	}()

	// Channel to handle interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Timeout after 20 seconds
	timeout := time.After(time.Duration(timeoutSec) * time.Second)

	// Main loop to handle interrupt signals and keep the connection alive
	for {
		select {
		case <-done:
			return 
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")
			// Cleanly close the connection by sending a close message and then waiting for the server to close the connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close message error:", err)
				return 
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return 
		case <-timeout:
			log.Println("Timeout reached, closing connection...")
			// Cleanly close the connection by sending a close message and then waiting for the server to close the connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close message error:", err)
				return 
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return 
		}
	}
}
