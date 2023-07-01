package ws

import (
	redis2 "chat-app/pkg/client/redis"
)

// Hub maintains the set of active ConnectedClients and broadcasts messages to the ConnectedClients.
type Hub struct {
	ConnectedClients  map[*Client]bool
	RegisterChannel   chan *Client
	UnregisterChannel chan *Client
	Rooms             map[*Room]bool
	redisClient       redis2.IRedisClient
}

func NewHub(redisCli redis2.IRedisClient) *Hub {
	return &Hub{
		ConnectedClients:  make(map[*Client]bool),
		RegisterChannel:   make(chan *Client),
		UnregisterChannel: make(chan *Client),
		Rooms:             make(map[*Room]bool),
		redisClient:       redisCli,
	}
}

// Run will execute Go Routines to check incoming Socket events
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.RegisterChannel:
			h.registerClient(client)
		case client := <-h.UnregisterChannel:
			h.unregisterClient(client)
		}
	}
}

// registerClient will handle the join event for new socket users
func (h *Hub) registerClient(client *Client) {
	h.ConnectedClients[client] = true
}

// unregisterClient will handle the disconnect event for new socket users
func (h *Hub) unregisterClient(client *Client) {
	delete(h.ConnectedClients, client)
}

// BroadcastToRoom sends the given message to all clients connected to the given room
func (h *Hub) BroadcastToRoom(message []byte, roomId string) {
	room := h.findRoomById(roomId)
	if room != nil {
		room.publishRoomMessage(message)
	}
}

func (h *Hub) findRoomById(id string) *Room {
	var foundRoom *Room
	for room := range h.Rooms {
		if room.GetId() == id {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (h *Hub) createRoom(roomId string) *Room {
	room := NewRoom(roomId, h.redisClient)
	h.Rooms[room] = true

	go room.RunRoom()

	return room
}
