package dto

type FileUploadDTO struct {
	PassPhrase string `form:"passPhrase" binding:"required"`
}
