package models

import "gorm.io/gorm"

type AccessToken struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
	Blockd bool   `json:"blocked"`
}
