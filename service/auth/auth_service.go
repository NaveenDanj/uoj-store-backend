package auth

import (
	"errors"
	"peer-store/core/pki"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"
	"time"

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

	if user, err := GetUserByRegistrationNumber(userDTO.RegistrationNumber); err == nil {
		return user, errors.New("user already registered with this Registration number")
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

	otp := generateOTP(6)

	newUser := models.User{
		Username:           userDTO.Name,
		Email:              userDTO.Email,
		PassPhrase:         string(hashedPassPhrase),
		Password:           string(hashedPassword),
		PubKey:             pubKey,
		PrivateKeyPath:     privateKeyPath,
		RegistrationNumber: userDTO.RegistrationNumber,
		OTP:                otp,
	}

	if err := db.GetDB().Create(&newUser).Error; err != nil {
		return newUser, nil
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

func generateOTP(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}
