package domain

import (
	"path/filepath"
	"time"
)

type MessageType string

const (
	MessageText  MessageType = "text"
	MessageFile  MessageType = "file"
	MessageImage MessageType = "image"
)

var fileExtensions = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".bmp",
	".tiff",
	".tif",
	".webp",
	".svg",
}

func StringToMessageType(s, url string) MessageType {
	if s == "image" {
		ext := filepath.Ext(url)
		for _, i := range fileExtensions {
			if ext == i {
				return MessageImage
			}
		}
		return MessageFile
	}
	if s == "raw" || s == "video" || s == "pdf" || s == "sprite" {
		return MessageFile
	}
	panic("not support resource type")
}

type ChatMessage struct {
	SenderID    string      `json:"sender_id,omitempty"`
	RecipientID string      `json:"recipient_id,omitempty"`
	Message     string      `json:"message,omitempty"`
	Time        time.Time   `json:"time,omitempty"`
	FileName    string      `json:"file_name,omitempty"`
	Height      uint32      `json:"height,omitempty"`
	Width       uint32      `json:"width,omitempty"`
	FileSize    uint32      `json:"file_size,omitempty"`
	URL         string      `json:"url,omitempty"`
	Type        MessageType `json:"type,omitempty"`
}

type RoomChatShortDetail struct {
	RoomName  string      `json:"room_name,omitempty"`
	RoomImage string      `json:"room_image,omitempty"`
	RoomId    string      `json:"room_id,omitempty"`
	IsGroup   bool        `json:"is_group"`
	Message   ChatMessage `json:"message"`
}
