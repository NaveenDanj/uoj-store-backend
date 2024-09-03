package dto

type UserRequestDTO struct {
	Name       string `json:"name" binding:"required,min=3,max=20"`
	Email      string `json:"email" binding:"required,email"`
	PassPhrase string `json:"age" binding:"required,len=32"`
	Password   string `json:"password" binding:"required,min=8"`
}
