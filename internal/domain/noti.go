package domain

type Notification struct {
	UserId  string              `json:"user_id,omitempty"`
	Message NotificationMessage `json:"message"`
}

type NotificationMessage struct {
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	ImageURL string `json:"image,omitempty"`
}
