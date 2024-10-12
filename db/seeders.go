package db

import (
	"fmt"
	"peer-store/core/pki"
	"peer-store/models"

	"golang.org/x/crypto/bcrypt"
)

// SeedAdminAccount seeds the admin account if it doesn't exist
func SeedAdminAccount() {
	var admin models.User
	err := GetDB().Where("role = ?", "Admin").First(&admin).Error
	passphrase := "naveen is my name. can you see m"
	password := "admin123"
	if err != nil && err.Error() == "record not found" {

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		hashedPassPhrase, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)

		if err != nil {
			fmt.Println("Error checking admin account:", err)
		}

		privateKeyPath, pubKey, err := pki.GeneratePkiKeyPair(passphrase)

		if err != nil {
			fmt.Println("Error checking admin account:", err)
		}

		admin = models.User{
			Username:       "admin",
			Email:          "admin@example.com",
			Password:       string(hashedPassword),
			Role:           "Admin",
			PassPhrase:     string(hashedPassPhrase),
			PubKey:         pubKey,
			PrivateKeyPath: privateKeyPath,
			IsVerified:     true,
			IsActive:       true,
			OTP:            "123456",
		}

		err = GetDB().Create(&admin).Error

		if err != nil {
			fmt.Println("Failed to create admin account:", err)
		} else {
			fmt.Println("Admin account created successfully.")
		}
	} else if err == nil {
		fmt.Println("Admin account already exists.")
	} else {
		fmt.Println("Error checking admin account:", err)
	}
}
