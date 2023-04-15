package domain

type EventName string

const (
	MessageEvent    EventName = "message"
	JoinEvent       EventName = "join"
	DisconnectEvent EventName = "disconnect"
)
