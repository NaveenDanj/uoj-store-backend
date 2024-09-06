package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("upfile")

	if err != nil {
		// Respond with a 400 Bad Request if there's an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file : " + err.Error()})
		return
	}

	tempFolder := "./disk/chunker/" + uuid.New().String()

	err = os.MkdirAll(tempFolder, os.ModePerm)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file : " + err.Error()})
		return
	}

	dst, err := os.Create(filepath.Join(tempFolder, header.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file on server"})
		return
	}

	defer dst.Close() // Ensure the destination file is closed after writing

	// Stream data from the incoming file to the destination file
	buf := make([]byte, 1024*1024) // 1MB buffer size
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
			return
		}
		if n == 0 {
			break
		}

		// Write the buffer to the destination file
		if _, err := dst.Write(buf[:n]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "File uploaded successfully",
		"file_path":   filepath.Join(tempFolder, header.Filename),
		"temp_folder": tempFolder,
		"filename":    header.Filename,
		"size":        header.Size,
	})

}
