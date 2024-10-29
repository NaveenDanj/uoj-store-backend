package models

import "gorm.io/gorm"

type AccessToken struct {
	gorm.Model
	UserId    uint   `json:"user_id"`
	Token     string `gorm:"uniqueIndex" json:"token"`
	IsSession bool   `json:"is_session"`
	Blocked   bool   `json:"blocked"`
}
