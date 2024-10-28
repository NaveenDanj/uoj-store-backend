package dto

type UserRequestDTO struct {
	Name       string `json:"name" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
	PassPhrase string `json:"passphrase" binding:"required,len=32"`
	Password   string `json:"password" binding:"required,min=8"`
}

type UserSignInDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PrivateSessionSignInDTO struct {
	SessionId uint `json:"session_id" binding:"required"`
}

type VerfyAccountDTO struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type ResetPasswordSendMailDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordNewPasswordDTO struct {
	Password string `json:"password" binding:"required,min=8"`
	Token    string `json:"token" binding:"required"`
}

type PassPhraseDTO struct {
	PassPhrase string `json:"pass_phrase" binding:"required"`
}

type UpdateUserProfileRequestDTO struct {
	TimoutTime uint `json:"timeout_time" binding:"required"`
}
