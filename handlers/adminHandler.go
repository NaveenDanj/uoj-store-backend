package handlers

import (
	"fmt"
	"net/http"
	"peer-store/core/pki"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"
	"peer-store/service"
	"peer-store/service/auth"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ActivateAccount(c *gin.Context) {
	var requestDto dto.ActivateAccountRequestDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := db.GetDB().Model(&models.User{}).Where("id = ?", requestDto.UserId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.IsActive = requestDto.Status
	if err := db.GetDB().Save(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error while changing account status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account status changed",
		"user":    user,
	})

}

func CreateAdminUser(c *gin.Context) {

	var requestDto dto.CreateAdminRequesDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := auth.GetUserByEmail(requestDto.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already registered with this email"})
		return
	}

	if _, err := auth.GetUserByUsername(requestDto.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already registered with this username"})
		return
	}

	otp := auth.GenerateOTP(6)

	newAdmin := models.User{
		Username:       requestDto.Username,
		Email:          requestDto.Email,
		Role:           "Admin",
		PassPhrase:     "",
		Password:       "",
		PubKey:         "",
		PrivateKeyPath: "",
		IsVerified:     true,
		IsActive:       false,
		OTP:            otp,
	}

	if err := db.GetDB().Create(&newAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating new admin user"})
		return
	}

	token := uuid.New().String()
	currentTime := time.Now().UTC()
	newTime := currentTime.Add(24 * time.Hour).UTC()

	resetPasssword := models.ResetToken{
		UserId:     newAdmin.ID,
		ResetToken: token,
		ExpireDate: newTime,
		Type:       "AdminAccountCreation",
	}

	if err := db.GetDB().Create(&resetPasssword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating setup account token"})
		return
	}

	html := service.ProcessSetupAdminAccountEmail(requestDto.Username, token)

	if err := service.SendEmail(requestDto.Email, "UOJ-Store Admin account setup invite", html); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while sending the otp email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": newAdmin})

}

func AdminAccountSetup(c *gin.Context) {

	var requestDto dto.AdminAccountSetupDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resetToken models.ResetToken
	if err := db.GetDB().Where("reset_token = ?", requestDto.Token).Where("type = ?", "AdminAccountCreation").First(&resetToken).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid link"})
		return
	}

	today := time.Now().UTC()

	if resetToken.ExpireDate.UTC().Compare(today) < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link expired"})
		return
	}

	var user models.User

	if err := db.GetDB().Where("id = ?", resetToken.UserId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestDto.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hasing password"})
		return
	}

	hashedPassPhrase, err := bcrypt.GenerateFromPassword([]byte(requestDto.Passphrase), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hasing Passphrase"})
		return
	}

	// generate keys
	privateKeyPath, pubKey, err := pki.GeneratePkiKeyPair(requestDto.Passphrase)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating public key and private key"})
		return
	}

	fmt.Println("Everythin created ---> " + string(hashedPassword))
	// update user model
	err = db.GetDB().Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"password":         string(hashedPassword),
		"pass_phrase":      string(hashedPassPhrase),
		"pub_key":          pubKey,
		"private_key_path": privateKeyPath,
		"is_active":        true,
	}).Error

	if err != nil {
		fmt.Println("Error updating admin account:", err)
	}

	db.GetDB().Unscoped().Where("reset_token = ?", requestDto.Token).Where("type = ?", "AdminAccountCreation").Delete(&models.ResetToken{})

	c.JSON(http.StatusOK, gin.H{"user": user})
}
