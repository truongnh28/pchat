package domain

type Group struct {
	Id       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	ImageUrl string   `json:"image_url,omitempty"`
	ImageId  uint     `json:"image_id,omitempty"`
	UserId   []string `json:"user_id,omitempty"`
}
