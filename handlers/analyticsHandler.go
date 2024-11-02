package handlers

import (
	"database/sql"
	"net/http"
	"peer-store/db"
	"peer-store/models"

	"github.com/gin-gonic/gin"
)

func GetTotalStorageByMimeType(c *gin.Context) {
	user, _ := c.Get("currentUser")
	userID := user.(models.User).ID

	mimeCategories := map[string][]string{

		// Image types
		"image": {
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"image/bmp", "image/tiff", "image/svg+xml", "image/heic",
		},

		// Video types
		"video": {
			"video/mp4", "video/avi", "video/mkv", "video/webm",
			"video/x-msvideo", "video/quicktime", "video/mpeg", "video/x-matroska",
		},

		// Audio types
		"audio": {
			"audio/mpeg", "audio/wav", "audio/ogg", "audio/mp4",
			"audio/aac", "audio/flac", "audio/x-ms-wma", "audio/x-wav",
		},

		// Document types
		"document": {
			"application/pdf", "application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.ms-powerpoint",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"application/rtf", "text/plain", "text/csv", "text/html", "application/epub+zip",
		},
	}

	storageUsage := map[string]float64{
		"image":    0,
		"video":    0,
		"audio":    0,
		"document": 0,
	}

	for category, mimeTypes := range mimeCategories {
		var totalBytes sql.NullInt64
		if err := db.GetDB().
			Model(&models.File{}).
			Where("user_id = ? AND mime_type IN ?", userID, mimeTypes).
			Select("SUM(file_size)").
			Scan(&totalBytes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error calculating storage usage"})
			return
		}

		if totalBytes.Valid {
			storageUsage[category] = float64(totalBytes.Int64) / (1024 * 1024)
		} else {
			storageUsage[category] = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"storage_usage": storageUsage,
	})
}

func GetTotalStorageUsage(c *gin.Context) {
	user, _ := c.Get("currentUser")
	userID := user.(models.User).ID

	var totalBytes sql.NullInt64
	if err := db.GetDB().
		Model(&models.File{}).
		Where("user_id = ?", userID).
		Select("SUM(file_size)").
		Scan(&totalBytes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error calculating total storage usage"})
		return
	}

	totalUsage := 0.0
	if totalBytes.Valid {
		totalUsage = float64(totalBytes.Int64) / (1024 * 1024)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_usage": totalUsage,
	})
}

func GetTopStorageUsers(c *gin.Context) {
	var users []struct {
		UserID     uint    `json:"user_id"`
		TotalUsage float64 `json:"total_usage"`
	}

	if err := db.GetDB().
		Model(&models.File{}).
		Select("user_id, SUM(file_size) AS total_usage").
		Group("user_id").
		Order("total_usage DESC").
		Limit(5).
		Scan(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving top storage users"})
		return
	}

	for i := range users {
		users[i].TotalUsage /= (1024 * 1024)
	}

	c.JSON(http.StatusOK, gin.H{
		"top_storage_users": users,
	})
}

func GetTopUploadDates(c *gin.Context) {
	var results []struct {
		UploadDate string  `json:"upload_date"`
		TotalSize  float64 `json:"total_size"`
	}

	if err := db.GetDB().
		Model(&models.File{}).
		Select("DATE(created_at) AS upload_date, SUM(file_size) AS total_size").
		Group("upload_date").
		Order("total_size DESC").
		Limit(5).
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving top upload dates"})
		return
	}

	type formattedResult struct {
		UploadDate string  `json:"upload_date"`
		TotalSize  float64 `json:"total_size"`
	}

	var formattedResults []formattedResult
	for _, result := range results {
		formattedResults = append(formattedResults, formattedResult{
			UploadDate: result.UploadDate,
			TotalSize:  result.TotalSize / (1024 * 1024),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"top_upload_dates": formattedResults,
	})
}

func GetTopFolders(c *gin.Context) {
	results := []struct {
		FolderID  uint    `json:"folder_id"`
		TotalSize float64 `json:"total_size"`
	}{}

	user, _ := c.Get("currentUser")
	userID := user.(models.User).ID

	if err := db.GetDB().
		Model(&models.File{}).
		Select("folder_id, SUM(file_size) AS total_size").
		Where("user_id = ?", userID).
		Group("folder_id").
		Order("total_size DESC").
		Limit(5).
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving top folders"})
		return
	}

	for i := range results {
		results[i].TotalSize /= (1024 * 1024)
	}

	c.JSON(http.StatusOK, gin.H{
		"top_folders": results,
	})
}

func GetTopFoldersDetailed(c *gin.Context) {
	results := []struct {
		FolderID   uint    `json:"folder_id"`
		FolderName string  `json:"folder_name"`
		TotalSize  float64 `json:"total_size"`
	}{}

	user, _ := c.Get("currentUser")
	userID := user.(models.User).ID

	if err := db.GetDB().
		Model(&models.File{}).
		Select("files.folder_id, folders.name AS folder_name, SUM(files.file_size) AS total_size").
		Joins("JOIN folders ON folders.id = files.folder_id").
		Where("files.user_id = ?", userID).
		Group("files.folder_id, folders.name").
		Order("total_size DESC").
		Limit(5).
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving top folders"})
		return
	}

	for i := range results {
		results[i].TotalSize /= (1024 * 1024)
	}

	c.JSON(http.StatusOK, gin.H{
		"top_folders": results,
	})
}

func GetTopFiles(c *gin.Context) {
	results := []struct {
		FileID       string  `json:"file_id"`
		OriginalName string  `json:"original_name"`
		FileSize     float64 `json:"file_size"`
	}{}

	user, _ := c.Get("currentUser")
	userID := user.(models.User).ID

	if err := db.GetDB().
		Model(&models.File{}).
		Select("file_id, original_name, file_size").
		Where("user_id = ?", userID).
		Order("file_size DESC").
		Limit(8).
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving top files"})
		return
	}

	for i := range results {
		results[i].FileSize /= (1024 * 1024)
	}

	c.JSON(http.StatusOK, gin.H{
		"top_files": results,
	})
}
