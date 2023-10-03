package webrtc

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/whatvn/denny"
)

type Event string

const (
	Candidate Event = "candidate"
	Answer    Event = "answer"
	Offer     Event = "offer"
)

type websocketMessage struct {
	Event Event  `json:"event"`
	Data  string `json:"data"`
}

func (s *websocketMessage) Encode() []byte {
	encoding, err := sonic.Marshal(s)
	if err != nil {
		denny.GetLogger(context.Background()).Error("Encode webrtc err: ", err)
	}
	return encoding
}
