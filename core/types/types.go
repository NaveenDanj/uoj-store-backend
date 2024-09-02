package types

import "io/fs"

type FileChunkMetaData struct {
	Sequence         []string    `json:"sequence"`
	FileInfo         fs.FileInfo `json:"fileInfo"`
	Extension        string      `json:"extension"`
	OriginalName     string      `json:"originalName"`
	DigitalSignature []string    `json:"digitalSignature"`
}

type FileUploadDTO struct {
	FileChunkMetaData *FileChunkMetaData
	Sequence          []string
}
