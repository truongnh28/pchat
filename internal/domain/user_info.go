package domain

type UserDetail struct {
	ID          string
	Username    string `json:"name,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Password    string `json:"password,omitempty"`
	Online      bool
	SocketId    string
}

// MessageConversation is a universal struct for mapping the conversations
type MessageConversation struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Message     string `json:"message"`
	SenderID    string `json:"SenderID"`
	RecipientID string `json:"recipientID"`
}

// UserDetailsRequest represents payload for Login and Registration request
type UserDetailsRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// UserDetailsResponse represents payload for Login and Registration response
type UserDetailsResponse struct {
	Username string `json:"username"`
	UserID   string `json:"userID"`
	Online   string `json:"online"`
}
