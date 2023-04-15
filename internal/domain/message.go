package domain

import "time"

type ChatMessage struct {
	SenderID    string    `json:"senderID,omitempty"`
	RecipientID string    `json:"recipientID,omitempty"`
	Message     string    `json:"message,omitempty"`
	Time        time.Time `json:"time,omitempty"`
}
