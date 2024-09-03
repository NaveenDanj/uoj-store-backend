package handlers

import (
	"net/http"
	"peer-store/dto"
	"peer-store/service/auth"

	"github.com/gin-gonic/gin"
)

func CreateNewUser(c *gin.Context) {

	var requestJSON dto.UserRequestDTO

	if err := c.ShouldBindJSON(&requestJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.CreateNewUser(&requestJSON)

	if err != nil && user.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is already used in another account!",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	authToken, err := auth.GenerateJWT(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"message":   "User account created successfully",
		"user":      requestJSON,
		"authToken": authToken,
	})

}
