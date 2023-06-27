package domain

import (
	"chat-app/helper"
	"time"
)

type (
	UpdatePassword struct {
		UserName string
		Password string
	}
	User struct {
		UserId      string            `json:"userId,omitempty"`
		Username    string            `json:"username,omitempty"`
		Email       string            `json:"email,omitempty"`
		PhoneNumber string            `json:"phoneNumber,omitempty"`
		Status      string            `json:"status,omitempty"`
		Code        string            `json:"code,omitempty"`
		DateOfBirth time.Time         `json:"date_of_birth,omitempty"`
		Gender      helper.GenderType `json:"gender,omitempty"`
		Url         string            `json:"url,omitempty"`
	}
)
