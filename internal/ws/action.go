package ws

type Event string

// Subscribed message
const (
	JoinUser  Event = "JoinUser"
	LeaveUser Event = "LeaveUser"

	JoinGroup  Event = "JoinGroup"
	LeaveGroup Event = "LeaveGroup"

	StartTyping Event = "StartTyping"
	StopTyping  Event = "StopTyping"
)

// Emit message
const (
	NewMessage    Event = "NewMessage"
	EditMessage   Event = "EditMessage"
	DeleteMessage Event = "DeleteMessage"
)
