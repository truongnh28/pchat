package wc

func NewHub() *Hub {
	return &Hub{
		ConnectedClients:  make(map[*Client]bool),
		RegisterChannel:   make(chan *Client),
		UnregisterChannel: make(chan *Client),
	}
}

// Run will execute Go Routines to check incoming Socket events
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.RegisterChannel:
			HandleUserRegisterEvent(h, client)
		case client := <-h.UnregisterChannel:
			HandleUserDisconnectEvent(h, client)
		}
	}
}
