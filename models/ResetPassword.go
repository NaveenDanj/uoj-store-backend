package models

import (
	"time"

	"gorm.io/gorm"
)

type ResetToken struct {
	gorm.Model
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId     uint      `json:"user_id"`
	ResetToken string    `json:"reset_token"`
	ExpireDate time.Time `json:"expire_date"`
	Type       string    `json:"type"`
}
