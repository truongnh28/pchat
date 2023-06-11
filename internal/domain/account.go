package domain

type (
	UpdatePassword struct {
		UserName string
		Password string
	}
	Account struct {
		UserId      uint64 `json:"userId,omitempty"`
		Username    string `json:"username,omitempty"`
		Email       string `json:"email,omitempty"`
		PhoneNumber string `json:"phoneNumber,omitempty"`
		Status      string `json:"status,omitempty"`
	}
)
