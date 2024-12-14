package connection

type WebSockets struct {
	mp map[string]*Websocket
}

func NewWebSockets() *WebSockets {
	return &WebSockets{mp: make(map[string]*Websocket)}
}

func (ws *WebSockets) AddWebsocket(symbol string, websocket *Websocket) {
	ws.mp[symbol] = websocket
}

func (ws *WebSockets) GetWebsocket(symbol string) *Websocket {
	websocket, ok := ws.mp[symbol]
	if !ok {
		return nil
	}
	return websocket
}

func (ws *WebSockets) GetAllWebsockets() map[string]*Websocket {
	return ws.mp
}