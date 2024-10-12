package server

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)
/**
 * call Init function with pairs, print and timeoutSec
 * It connects to my WebSocket server and subscribes to the given pairs
 * pairs: []string - list of pairs to subscribe
 * print: bool - print the received messages
 * timeoutSec: int - timeout in seconds
*/
func Init(pairs []string, print bool, timeoutSec int) {
	firstPrinted := false
	startTime := time.Now()

	wsURL := "ws://localhost:8080/ws"

	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatal("Failed to parse the WebSocket URL:", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket:", err)
	}
	defer conn.Close()

	json := NewSubscription(pairs).GetJson()

	// send subscription message
	err = conn.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	// uncomment this function to simulate a subscription overwriting

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	if len(pairs) > 0 {
	// 		pairs = pairs[:len(pairs)-1]
	// 	}
	// 	pairs = append(pairs, "ADAUSDT")
	// 	json = NewSubscription(pairs).GetJson()
	// 	err = conn.WriteMessage(websocket.TextMessage, json)
	// 	if err != nil {
	// 		log.Println("Write error:", err)
	// 		return
	// 	}
	// }()

	done := make(chan struct{})
	// Channel to signal when the connection is closed
	closeConn := make(chan struct{})
	var wg sync.WaitGroup

	// Goroutine to read messages from WebSocket
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(done)
		for {
			select {
			case <-closeConn:
				// Exit the goroutine when the connection is closed
				return
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					if !websocket.IsCloseError(err) {
						log.Println("Read error:", err)
					}
					return
				}
				if !firstPrinted {
					log.Println("SERVER first receive time elapse:", time.Since(startTime))
					firstPrinted = true
				}
				if print {
					log.Printf("SERVER Received: %s\n", message)
				}
			}
		}
	}()

	// Channel to handle interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	timeout := time.After(time.Duration(timeoutSec) * time.Second)

	// Main loop to handle interrupt signals and keep the connection alive
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")
			// Signal the goroutine to stop reading messages
			close(closeConn)
			// Wait for the goroutine to finish
			wg.Wait()
			// Close the WebSocket connection here
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			conn.Close()
			return
		case <-timeout:
			log.Println("Timeout reached, closing connection...")
			// Signal the goroutine to stop reading messages
			close(closeConn)
			// Wait for the goroutine to finish
			wg.Wait()
			// Close the WebSocket connection here
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			conn.Close()
			return
		}
	}
}
