package models

import (
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	UserId      string        `gorm:"column:id;primaryKey;"`
	UserName    string        `gorm:"index:username_idx_uni,unique"`
	Email       string        `gorm:"column:email;unique"`
	PhoneNumber string        `gorm:"column:phone_number;unique"`
	Password    string        `gorm:"column:password"`
	Status      AccountStatus `gorm:"column:status"`
}

type AccountStatus string

const (
	Active  AccountStatus = "Active"
	Blocked AccountStatus = "Blocked"
)

type Group struct {
	*gorm.Model
	GroupId   string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name;unique"`
	AvatarUrl string `gorm:"column:avatar_url"`
}

type Room struct {
	*gorm.Model
	GroupId string `gorm:"column:group_id;primaryKey"`
	UserId  string `gorm:"column:user_id;primaryKey"`
}
