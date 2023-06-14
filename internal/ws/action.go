package ws

type Event string

// Subscribed message
const (
	JoinRoom  Event = "JoinRoom"
	LeaveRoom Event = "LeaveRoom"

	StartTyping Event = "StartTyping"
	StopTyping  Event = "StopTyping"
)

// Emit message
const (
	NewMessage    Event = "NewMessage"
	EditMessage   Event = "EditMessage"
	DeleteMessage Event = "DeleteMessage"
)
