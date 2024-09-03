package auth

import (
	"errors"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateNewUser(userDTO *dto.UserRequestDTO) (models.User, error) {

	if user, err := GetUserByEmail(userDTO.Email); err == nil {
		return user, errors.New("A User already registered with this email!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return models.User{}, errors.New("Internal server error!")

	}

	hashedPassPhrase, err := bcrypt.GenerateFromPassword([]byte(userDTO.PassPhrase), bcrypt.DefaultCost)

	if err != nil {
		return models.User{}, errors.New("Internal server error!")

	}

	newUser := models.User{
		Name:       userDTO.Name,
		Email:      userDTO.Email,
		PassPhrase: string(hashedPassPhrase),
		Password:   string(hashedPassword),
	}

	if err := db.GetDB().Create(&newUser).Error; err != nil {
		return newUser, nil
	}

	return newUser, nil

}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := db.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		return user, errors.New("User not found!")
	}

	return user, nil
}
