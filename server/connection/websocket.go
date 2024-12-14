package connection

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	conn *websocket.Conn
	url  string
	Ch   chan []byte
	Quit chan struct{}
	mu   sync.Mutex
}

func NewWebsocket(symbol string) *Websocket {
	return &Websocket{Ch: make(chan []byte),
		url: "wss://fstream.binance.com/ws/" + symbol + "@bookTicker", Quit: make(chan struct{})}
}

func (w *Websocket) Connect() {
	w.Quit = make(chan struct{})

	_conn, _, err := websocket.DefaultDialer.Dial(w.url, nil)
	log.Println("Connecting to WebSocket ", w.url)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	w.conn = _conn
}

func (w *Websocket) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	log.Println("Closing WebSocket ", w.url)
	w.conn.Close()
	close(w.Quit)
}

func (w *Websocket) Start() {
	for {
		select {
		case <-w.Quit:
			return
		default:
		}
		w.mu.Lock()
		msgType, message, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
		}
		w.mu.Unlock()
		if msgType == websocket.PingMessage {
			// Respond to ping with pong
			err := w.conn.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				log.Println("Error writing pong to WebSocket:", err)
			}
		} else {
			w.Ch <- message
		}
	}
}

func (w *Websocket) Write(message string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	err := w.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Error writing to WebSocket:", err)
	}
}
