package models

import "gorm.io/gorm"

type Account struct {
	*gorm.Model
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
