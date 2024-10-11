package dto

type UserRequestDTO struct {
	Name               string `json:"name" binding:"required,min=3,max=100"`
	Email              string `json:"email" binding:"required,email"`
	RegistrationNumber string `json:"registration_number" binding:"required,len=10"`
	PassPhrase         string `json:"passphrase" binding:"required,len=32"`
	Password           string `json:"password" binding:"required,min=8"`
}

type UserSignInDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerfyAccountDTO struct {
	UserId string `json:"user_id" binding:"required"`
	OTP    string `json:"otp" binding:"required"`
}

type ResetPasswordSendMailDTO struct {
	Email string `json:"email" binding:"required,email"`
}
