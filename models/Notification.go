package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserId  uint   `json:"user_id"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}
