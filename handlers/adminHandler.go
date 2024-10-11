package handlers

import (
	"net/http"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, gin.H{"message": "Account activated successfully"})

}
