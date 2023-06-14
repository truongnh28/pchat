package domain

type Room struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}
