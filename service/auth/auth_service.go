package auth

import (
	"errors"
	"peer-store/core/pki"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"
	"peer-store/service"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func CreateNewUser(userDTO *dto.UserRequestDTO) (models.User, error) {

	if user, err := GetUserByEmail(userDTO.Email); err == nil {
		return user, errors.New("user already registered with this email")
	}

	if user, err := GetUserByUsername(userDTO.Name); err == nil {
		return user, errors.New("user already registered with this username")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return models.User{}, errors.New("internal server error")
	}

	hashedPassPhrase, err := bcrypt.GenerateFromPassword([]byte(userDTO.PassPhrase), bcrypt.DefaultCost)

	if err != nil {
		return models.User{}, errors.New("internal server error")
	}

	// generate keys
	privateKeyPath, pubKey, err := pki.GeneratePkiKeyPair(userDTO.PassPhrase)

	if err != nil {
		return models.User{}, errors.New("internal server error")
	}

	otp := GenerateOTP(6)
	sessionId := GenerateOTP(8)

	newUser := models.User{
		Username:       userDTO.Name,
		Email:          userDTO.Email,
		PassPhrase:     string(hashedPassPhrase),
		Password:       string(hashedPassword),
		PubKey:         pubKey,
		PrivateKeyPath: privateKeyPath,
		OTP:            otp,
		SessionId:      sessionId,
	}

	html := service.ProcessOTPEmail(otp, userDTO.Name)

	if err := service.SendEmail(userDTO.Email, "UOJ-Store Account verification OTP", html); err != nil {
		return models.User{}, errors.New("error while sending OTP")
	}

	if err := db.GetDB().Create(&newUser).Error; err != nil {
		return newUser, err
	}

	return newUser, nil

}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := db.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}

	return user, nil

}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := db.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

func GetUserByRegistrationNumber(regNo string) (models.User, error) {
	var user models.User
	if err := db.GetDB().Where("registration_number = ?", regNo).First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}

	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	if err := db.GetDB().Find(&users).Error; err != nil {
		return users, errors.New("error while fetching user data")
	}

	return users, nil

}

func GetPublicKey(userId string) (string, error) {
	var pubKey string
	if err := db.GetDB().Model(models.User{}).Select("pub_key").Where("id = ?", userId).First(&pubKey).Error; err != nil {
		return "", err
	}
	return pubKey, nil
}

func BlockToken(token string) error {

	var tokenRecord models.AccessToken

	if err := db.GetDB().Where("token = ?", token).First(&tokenRecord).Error; err != nil {
		return err
	}

	if err := db.GetDB().Model(&models.AccessToken{}).Where("token = ?", token).Update("blocked", true).Error; err != nil {
		return err
	}

	return nil
}

func IsBlocked(token string) bool {
	var tokenRecord models.AccessToken
	if err := db.GetDB().Where("token = ?", token).First(&tokenRecord).Error; err != nil {
		return true
	}
	return tokenRecord.Blocked
}

func VerifyAccount(otp string, userEmail string) bool {
	var user models.User

	if err := db.GetDB().Where("email = ?", userEmail).First(&user).Error; err != nil {
		return false
	}

	if user.IsVerified {
		return false
	}

	if user.OTP != otp {
		return false
	}

	user.IsVerified = true
	if err := db.GetDB().Save(&user).Error; err != nil {
		return false
	}

	return true
}

func ResetPasswordGenerateLink(userid uint) (bool, error) {

	var user models.User
	if err := db.GetDB().Where("id = ?", userid).First(&user).Error; err != nil {
		return false, errors.New("user not found")
	}

	if !user.IsVerified {
		return false, errors.New("user is not verified")
	}

	if !user.IsActive {
		return false, errors.New("User is not active")
	}

	db.GetDB().Where("user_id = ?", user.ID).Delete(&models.ResetToken{})

	token := uuid.New().String()
	currentTime := time.Now().UTC()
	newTime := currentTime.Add(30 * time.Minute).UTC()

	resetPasssword := models.ResetToken{
		UserId:     user.ID,
		ResetToken: token,
		ExpireDate: newTime,
		Type:       "PasswordReset",
	}

	if err := db.GetDB().Create(&resetPasssword).Error; err != nil {
		return false, errors.New("error while generating reset password link")
	}

	link := "http://localhost:5173/auth/reset-password?token=" + token

	// send the email
	html := service.ProcessResetPasswordEmail(user.Username, link)

	if err := service.SendEmail(user.Email, "UOJ-Store Password reset link", html); err != nil {
		return false, errors.New("error while sending password reset link")
	}

	return true, nil

}

func HandleResetPassword(token string, newPassword string) error {
	var resetPassword models.ResetToken
	if err := db.GetDB().Where("reset_token = ?", token).Where("type = ?", "PasswordReset").First(&resetPassword).Error; err != nil {
		return errors.New("Invalid password reset link")
	}

	today := time.Now().UTC()

	if resetPassword.ExpireDate.UTC().Compare(today) < 0 {
		return errors.New("link expired")
	}

	var user models.User

	if err := db.GetDB().Where("id = ?", resetPassword.UserId).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("reset password failed")
	}

	db.GetDB().Model(&models.User{}).Where("id = ?", user.ID).Update("password", string(hashedPassword))
	// db.GetDB().Where("reset_token = ?", token).Delete(&models.ResetPassword{})
	db.GetDB().Unscoped().Where("reset_token = ?", token).Where("type = ?", "PasswordReset").Delete(&models.ResetToken{})

	return nil

}

func GenerateOTP(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}
