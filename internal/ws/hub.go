package ws

import (
	"chat-app/internal/service"
	"github.com/go-redis/redis/v8"
)

// Hub maintains the set of active ConnectedClients and broadcasts messages to the ConnectedClients.
type Hub struct {
	ConnectedClients  map[*Client]bool
	RegisterChannel   chan *Client
	UnregisterChannel chan *Client
	Rooms             map[*Room]bool
	UserService       service.UserService
	MessageService    service.MessageService
	redisClient       *redis.Client
}

func NewHub(
	userService service.UserService,
	messageService service.MessageService,
) *Hub {
	return &Hub{
		ConnectedClients:  make(map[*Client]bool),
		RegisterChannel:   make(chan *Client),
		UnregisterChannel: make(chan *Client),
		UserService:       userService,
		MessageService:    messageService,
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

// EmitToSpecificClient will emit the socket event to specific socket user(sendChannel message to single user)
// Based on the userID provided, it sends socket Payload to that user, whos id matches with the id of
// users stored in socket server.
func (h *Hub) EmitToSpecificClient(payload []byte, userId string) {
	for client := range h.ConnectedClients {
		if client.id == userId {
			select {
			case client.sendChannel <- payload:
			default:
				close(client.sendChannel)
				delete(h.ConnectedClients, client)
			}
		}
	}
}

// BroadcastSocketEventToAllClient will emit the socket events to all socket users
// function will sendChannel the socket events to all the users connected to Socket Server by using a ←sendChannel channel.
func (h *Hub) BroadcastSocketEventToAllClient(payload []byte) {
	for client := range h.ConnectedClients {
		select {
		case client.sendChannel <- payload:
		default:
			close(client.sendChannel)
			delete(h.ConnectedClients, client)
		}
	}
}

func (h *Hub) BroadcastSocketEventToAllClientExceptMe(payload []byte, myUserId string) {
	for client := range h.ConnectedClients {
		if client.id != myUserId {
			select {
			case client.sendChannel <- payload:
			default:
				close(client.sendChannel)
				delete(h.ConnectedClients, client)
			}
		}
	}
}

// BroadcastToRoom sends the given message to all clients connected to the given room
func (h *Hub) BroadcastToRoom(message []byte, roomId string) {
	if room := h.findRoomById(roomId); room != nil {
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

func (h *Hub) createRoom(id string) *Room {
	room := NewRoom(id, h.redisClient)
	go room.RunRoom()
	h.Rooms[room] = true

	return room
}
