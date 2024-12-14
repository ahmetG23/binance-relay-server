package subscription

import (
	"binance-server/schema"
	"encoding/json"
	"log"
)

func (sub *Subscription) readBookTicker(symbol string, tickerch chan<- schema.BookTicker) {
    bookTicker := schema.BookTicker{Symbol: symbol}
    websocket := sub.List.Websockets.GetWebsocket(symbol)

    // Start reading from the websocket
    websocket.Connect()
    go websocket.Start()

    for {
        select {
        case msg := <-websocket.Ch:
            err := json.Unmarshal(msg, &bookTicker)
            if err != nil {
                log.Println("Error unmarshalling JSON:", err)
            } else {
                tickerch <- bookTicker
            }
        case <-sub.Quits[symbol]:
            websocket.Close()  
            return
        }
    }
}

func (sub *Subscription) Read(tickerch chan<- schema.BookTicker, oldPairs []string) {
    newPairs := make(map[string]bool)

    for _, symbol := range sub.Pairs {
        newPairs[symbol] = true

        if sub.Mp[symbol] {
            continue
        } else {
            sub.Mp[symbol] = true
            sub.Quits[symbol] = make(chan struct{})
            go sub.readBookTicker(symbol, tickerch)
        }
    }

    // Close channels for old pairs that are no longer in use
    for _, symbol := range oldPairs {
        if !newPairs[symbol] {
            if quitCh, ok := sub.Quits[symbol]; ok {
                close(quitCh)
                delete(sub.Mp, symbol)
                delete(sub.Quits, symbol)
            }
        }
    }
}