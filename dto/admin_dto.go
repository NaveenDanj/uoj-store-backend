package dto

type ActivateAccountRequestDTO struct {
	UserId uint `json:"userId" binding:"required"`
	Status bool `json:"status" binding:"required"`
}
