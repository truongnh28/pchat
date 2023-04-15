package wc

import (
	"chat-app/internal/domain"
	"chat-app/pkg/repositories"
	"github.com/gorilla/websocket"
)

// SocketEventStruct is a type for socket events which
type SocketEventStruct struct {
	EventName    domain.EventName `json:"eventName"`
	EventPayload interface{}      `json:"eventPayload"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub         *Hub
	wsConn      *websocket.Conn
	sendChannel chan SocketEventStruct
	userId      string
}

// Hub maintains the set of active ConnectedClients and broadcasts messages to the ConnectedClients.
type Hub struct {
	ConnectedClients  map[*Client]bool
	BroadcastChannel  chan []byte
	RegisterChannel   chan *Client
	UnregisterChannel chan *Client
	UserRepo          repositories.UserRepository
	MessageRepo       repositories.MessageRepository
}
