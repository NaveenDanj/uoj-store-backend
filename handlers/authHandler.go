package handlers

import (
	"net/http"
	"peer-store/dto"
	"peer-store/service/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateNewUser(c *gin.Context) {

	var requestJSON dto.UserRequestDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate key pairs and store public key on the user table and private key on the backend

	user, err := auth.CreateNewUser(&requestJSON)

	if err != nil && user.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	authToken, err := auth.GenerateJWT(user.ID, user.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "User account created successfully",
		"user":      requestJSON,
		"authToken": authToken,
	})

}

func UserSignIn(c *gin.Context) {
	var requestJSON dto.UserSignInDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.GetUserByUsername(requestJSON.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestJSON.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username or password is in-correct, Please try again.",
		})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user account is not verified",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user account is not activated",
		})
		return
	}

	authToken, err := auth.GenerateJWT(user.ID, user.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login success",
		"user":      user,
		"authToken": authToken,
	})

}

func GetCurrentUser(c *gin.Context) {

	user, _ := c.Get("currentUser")

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

func Logout(c *gin.Context) {

	authToken := c.GetHeader("Authorization")
	err := auth.BlockToken(authToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout success!",
	})

}

func GetUserPublicKey(c *gin.Context) {

	userId := c.Param("userId")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	pubKey, err := auth.GetPublicKey(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while fetching user public key ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Public key retrieved successfully",
		"pubKey":  pubKey,
	})

}

func VerifyAccount(c *gin.Context) {
	var requestJSON dto.VerfyAccountDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	res := auth.VerifyAccount(requestJSON.OTP, requestJSON.Email)

	if !res {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Account verification failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
	return

}

func ResetPasswordSendMail(c *gin.Context) {
	var requestJSON dto.ResetPasswordSendMailDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := auth.GetUserByEmail(requestJSON.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to send reset password link"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to send reset password link"})
		return
	}

	if _, err := auth.ResetPasswordGenerateLink(user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to send reset password link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func ResetPasswordSetPassword(c *gin.Context) {
	var requestJSON dto.ResetPasswordNewPasswordDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := auth.HandleResetPassword(requestJSON.Token, requestJSON.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reseted successfully!"})

}
