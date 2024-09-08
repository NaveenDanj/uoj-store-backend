package dto

type FileUploadDTO struct {
	PassPhrase string `form:"passPhrase" binding:"required"`
}

type FileDownloadRequestDTO struct {
	PassPhrase string `json:"passPhrase" binding:"required"`
	FileId     string `json:"fileId" binding:"required"`
}

type FileDeleteRequestDTO struct {
	FileId string `json:"fileId" binding:"required"`
}
