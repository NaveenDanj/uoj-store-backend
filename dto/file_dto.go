package dto

type FileUploadDTO struct {
	PassPhrase string `form:"passPhrase" binding:"required"`
}

type FileDownloadRequestDTO struct {
	PassPhrase string `form:"passPhrase" binding:"required"`
	FileId     string `form:"fileId" binding:"required"`
}
