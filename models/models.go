package models

import (
	"chat-app/helper"
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId      string            `gorm:"column:user_id;primary_key;"`
	UserName    string            `gorm:"index:username_idx_uni,unique"`
	Email       string            `gorm:"column:email;unique"`
	PhoneNumber string            `gorm:"column:phone_number;unique"`
	Password    string            `gorm:"column:password"`
	Status      AccountStatus     `gorm:"column:status"`
	DateOfBirth time.Time         `gorm:"column:date_of_birth"`
	Gender      helper.GenderType `gorm:"column:gender"`
	FileId      uint              `gorm:"column:file_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type AccountStatus string

const (
	Active  AccountStatus = "Active"
	Blocked AccountStatus = "Blocked"
)

type Group struct {
	GroupId   string `gorm:"column:group_id;primary_key"`
	Name      string `gorm:"column:name"`
	FileId    uint   `gorm:"column:file_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Room struct {
	*gorm.Model
	GroupId string `gorm:"column:group_id;primary_key"`
	UserId  string `gorm:"column:user_id;primary_key"`
}

type File struct {
	*gorm.Model
	AssetID          string `gorm:"column:asset_id"`
	PublicID         string `gorm:"column:public_id"`
	AssetFolder      string `gorm:"column:asset_folder"`
	DisplayName      string `gorm:"column:display_name"`
	URL              string `gorm:"column:url"`
	SecureURL        string `gorm:"column:secure_url"`
	OriginalFilename string `gorm:"column:original_filename"`
	Height           uint32 `gorm:"column:height"`
	Width            uint32 `gorm:"column:width"`
	FileSize         uint32 `gorm:"column:file_size"`
	ResourceType     string `gorm:"column:resource_type"`
}

type Friend struct {
	*gorm.Model
	UserId1 string `gorm:"column:user_id_1;primary_key"`
	UserId2 string `gorm:"column:user_id_2;primary_key"`
}
