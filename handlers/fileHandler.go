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

	_, err = file.Seek(0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to reset file pointer: " + err.Error()})
		return
	}

	metaData, err := storage.FileUploader(file, header)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file : " + err.Error()})
		return
	}

	UploadedFileData, err := storage.StoreUploadedFile(mimeType, &metaData, &userObj, requestForm.PassPhrase, requestForm.FolderId)

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

	var requestForm dto.FileDownloadRequestDTO

	if err := c.ShouldBindJSON(&requestForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	userObj, _ := user.(models.User)

	path, mimeType, err := storage.HandleDownloadProcess(requestForm.FileId, &userObj, requestForm.PassPhrase)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", mimeType)
	c.File(path)

	// delete the file if it exists
	storage.DeleteFile(path)
	storage.DeleteFolder("./disk/public/" + requestForm.FileId)

}

func DeleteFile(c *gin.Context) {
	var requestForm dto.FileDeleteRequestDTO

	if err := c.ShouldBindJSON(&requestForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")
	userObj, _ := user.(models.User)

	if err := storage.FileDeleteService(requestForm.FileId, &userObj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})

}

func GetUserFiles(c *gin.Context) {

	user, _ := c.Get("currentUser")
	userObj, _ := user.(models.User)

	files, err := storage.GetUserFiles(userObj.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})

	return
}
