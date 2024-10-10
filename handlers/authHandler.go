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

	user, err := auth.GetUserByEmail(requestJSON.Email)

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

	authToken, err := auth.GenerateJWT(user.ID, user.Email)

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
