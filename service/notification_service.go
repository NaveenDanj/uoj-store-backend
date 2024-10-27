package service

import (
	"peer-store/db"
	"peer-store/models"
)

func CreateNewNotification(userId uint, message string) bool {

	n := models.Notification{
		UserId:  userId,
		Message: message,
	}

	if err := db.GetDB().Create(&n).Error; err != nil {
		return false
	}

	return true
}
