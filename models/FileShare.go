package models

import (
	"time"

	"gorm.io/gorm"
)

type FileShare struct {
	gorm.Model
	Id            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	FileId        string    `json:"file_id"`
	OwnerId       uint      `json:"owner_id"`
	Token         string    `json:"token" gorm:"uniqueIndex"`
	IsPublic      bool      `json:"is_public"`
	ExpireDate    time.Time `json:"expire_date"`
	Status        string    `json:"status"`
	Note          string    `json:"note"`
	SharedAt      time.Time `json:"shared_at"`
	DownloadLimit int       `json:"download_limit"`
	EncryptionKey string    `json:"encryption_key"`
}

type FileShareUser struct {
	gorm.Model
	Id            int  `json:"id" gorm:"primaryKey;autoIncrement"`
	FileShareId   uint `json:"file_share_id"`
	UserId        uint `json:"user_id"`
	DownloadCount int  `json:"download_count"`
}
