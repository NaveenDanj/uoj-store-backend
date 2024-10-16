package dto

type FileUploadDTO struct {
	PassPhrase string `form:"passPhrase" binding:"required"`
	FolderId   uint   `form:"folder_id" binding:"required"`
}

type FileDownloadRequestDTO struct {
	PassPhrase string `json:"passPhrase" binding:"required"`
	FileId     string `json:"fileId" binding:"required"`
}

type FileDeleteRequestDTO struct {
	FileId string `json:"fileId" binding:"required"`
}

type FileShareUserRequestDTO struct {
	UserId        uint `json:"userId" binding:"required"`
	DownloadLimit int  `json:"downloadLimit" binding:"required"`
}

type FileShareRequestDTO struct {
	FileId        string                    `json:"fileId" binding:"required"`
	Users         []FileShareUserRequestDTO `json:"users" binding:"required"`
	PassPhrase    string                    `json:"passPhrase" binding:"required"`
	ExpireDate    string                    `json:"expireDate" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Note          string                    `json:"note,omitempty"`
	DownloadLimit int                       `json:"downloadLimit" binding:"required,gte=1"`
}

type LinkRevokeRequestDTO struct {
	Link string `json:"link" binding:"required"`
}
