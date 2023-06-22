package ws

type Event string

// Subscribed message
const (
	JoinRoom  Event = "JoinRoom"
	LeaveRoom Event = "LeaveRoom"

	JoinUser  Event = "JoinUser"
	LeaveUser Event = "LeaveUser"

	StartTyping Event = "StartTyping"
	StopTyping  Event = "StopTyping"
)

// Emit message
const (
	NewMessage    Event = "NewMessage"
	EditMessage   Event = "EditMessage"
	DeleteMessage Event = "DeleteMessage"
)
