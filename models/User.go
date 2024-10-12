package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string `gorm:"uniqueIndex" json:"username"`
	Email          string `gorm:"uniqueIndex" json:"email"`
	Role           string `json:"role"`
	PassPhrase     string `json:"-"`
	Password       string `json:"-"`
	PubKey         string `json:"-"`
	PrivateKeyPath string `json:"-"`
	IsVerified     bool   `json:"is_verified"`
	IsActive       bool   `json:"is_active"`
	OTP            string `json:"-"`
}
