package dto

type UserRequestDTO struct {
	Name       string `json:"name" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
	PassPhrase string `json:"passphrase" binding:"required,len=32"`
	Password   string `json:"password" binding:"required,min=8"`
}

type UserSignInDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
