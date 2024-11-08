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
	RootFolder     uint   `json:"root_folder"`
	SessionFolder  uint   `json:"session_folder"`
	WorkFolder     uint   `json:"work_folder"`
	PersonalFolder uint   `json:"personal_folder"`
	AcademicFolder uint   `json:"academic_folder"`
	SessionId      string `gorm:"uniqueIndex" json:"session_id"`
	SessionTime    uint   `json:"session_time"`
	MaxUploadSize  uint   `json:"max_upload_size"`
}
