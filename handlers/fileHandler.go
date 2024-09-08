package handlers

import (
	"net/http"
	"peer-store/dto"
	"peer-store/models"
	"peer-store/service/storage"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {

	var requestForm dto.FileUploadDTO

	if err := c.ShouldBind(&requestForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	userObj, _ := user.(models.User)

	file, header, err := c.Request.FormFile("upfile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file : " + err.Error()})
		return
	}

	if !storage.ValidatePassPhrase(requestForm.PassPhrase, &userObj) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file : Invalid pass phrase"})
		return
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}

	// Detect the MIME type
	mimeType := http.DetectContentType(buffer)

	metaData, err := storage.FileUploader(file, header)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file : " + err.Error()})
		return
	}

	UploadedFileData, err := storage.StoreUploadedFile(mimeType, &metaData, &userObj, requestForm.PassPhrase)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"FileData": UploadedFileData,
	})

}

func DownloadFile(c *gin.Context) {

}
