package ws

import (
	"context"
	"encoding/json"
	"github.com/whatvn/denny"
)

// SocketMessage is a type for socket events which
type SocketMessage struct {
	Event   Event `json:"event"`
	Payload any   `json:"payload"`
}

// ReceivedMessage represents a received websocket message
type ReceivedMessage struct {
	Event   Event `json:"event"`
	Payload any   `json:"payload"`
}

func (s *SocketMessage) Encode() []byte {
	encoding, err := json.Marshal(s)
	if err != nil {
		denny.GetLogger(context.Background()).Error("Encode Socket err: ", err)
	}
	return encoding
}
