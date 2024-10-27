package dto

type ActivateAccountRequestDTO struct {
	UserId uint   `json:"userId" binding:"required"`
	Status bool   `json:"status"`
	Role   string `json:"role" binding:"required"`
}

type CreateAdminRequesDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Status   bool   `json:"status"`
}

type AdminAccountSetupDTO struct {
	Token      string `json:"token" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
	Passphrase string `json:"passphrase" binding:"required,len=32"`
}
