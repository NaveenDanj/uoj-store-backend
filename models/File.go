package models

import "time"

type File struct {
	FileId        string    `json:"file_id" gorm:"uniqueIndex"`
	UserId        uint      `json:"user_id"`
	OriginalName  string    `json:"original_name"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	UploadDate    time.Time `json:"upload_date"`
	StoragePath   string    `json:"storage_path"`
	FileSignature []byte    `json:"file_signature"`
	DownloadCount int       `json:"download_count"`
	LastAccessed  time.Time `json:"last_accessed"`
}
