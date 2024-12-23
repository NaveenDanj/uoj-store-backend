package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	FileId        string    `json:"file_id" gorm:"uniqueIndex"`
	FolderID      uint      `json:"folder_id"`
	UserId        uint      `json:"user_id"`
	OriginalName  string    `json:"original_name"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	UploadDate    time.Time `json:"upload_date"`
	StoragePath   string    `json:"storage_path"`
	FileSignature []byte    `json:"file_signature"`
	DownloadCount int       `json:"download_count"`
	LastAccessed  time.Time `json:"last_accessed"`
	IsFavourite   bool      `json:"is_favourite"`
	IsDeleted     bool      `json:"is_deleted"`
	IsPublic      bool      `json:"is_public"`
	IsInSession   bool      `json:"is_in_session"`
}

type Folder struct {
	gorm.Model
	UserId        uint   `json:"user_id"`
	Name          string `json:"name" gorm:"not null"`
	ParentID      *uint  `json:"parent_id" gorm:"index"`
	SpecialFolder string `json:"special_folder"`
	IsDeleted     bool   `json:"is_deleted"`
}
