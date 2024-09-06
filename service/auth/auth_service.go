package auth

import (
	"errors"
	"peer-store/core/pki"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateNewUser(userDTO *dto.UserRequestDTO) (models.User, error) {

	if user, err := GetUserByEmail(userDTO.Email); err == nil {
		return user, errors.New("user already registered with this email")
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

	newUser := models.User{
		Name:           userDTO.Name,
		Email:          userDTO.Email,
		PassPhrase:     string(hashedPassPhrase),
		Password:       string(hashedPassword),
		PubKey:         pubKey,
		PrivateKeyPath: privateKeyPath,
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
