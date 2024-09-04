package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `json:"name"`
	Email          string `gorm:"uniqueIndex" json:"email"`
	PassPhrase     string `json:"-"`
	Password       string `json:"-"`
	PubKey         string `json:"-"`
	PrivateKeyPath string `json:"-"`
}
