package ws

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/whatvn/denny"
)

// ReceivedMessage represents a received websocket message
type ReceivedMessage struct {
	Event   Event  `json:"event"`
	Room    string `json:"room"`
	Payload *any   `json:"payload"`
}

// SocketMessage is a type for socket events which
type SocketMessage struct {
	Event   Event `json:"event"`
	Payload any   `json:"payload"`
}

func (s *SocketMessage) Encode() []byte {
	encoding, err := sonic.Marshal(s)
	if err != nil {
		denny.GetLogger(context.Background()).Error("Encode Socket err: ", err)
	}
	return encoding
}
