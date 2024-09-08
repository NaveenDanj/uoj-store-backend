package dto

type FileUploadDTO struct {
	PassPhrase string `form:"passPhrase" binding:"required"`
}

type FileDownloadRequestDTO struct {
	PassPhrase string `json:"passPhrase" binding:"required"`
	FileId     string `json:"fileId" binding:"required"`
}
