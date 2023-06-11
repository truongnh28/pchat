package domain

import "time"

type ChatMessage struct {
	SenderID    string    `json:"senderID,omitempty" json:"sender_id"`
	RecipientID string    `json:"recipientID,omitempty" json:"recipient_id"`
	Message     string    `json:"message,omitempty"`
	Time        time.Time `json:"time,omitempty"`
}
