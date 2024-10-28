package handlers

import (
	"net/http"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"
	"peer-store/service"
	"peer-store/service/folder"
	"peer-store/service/storage"

	"github.com/gin-gonic/gin"
)

func CreateFolder(c *gin.Context) {

	var requestDto dto.CreateFolderDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, e := c.Get("currentUser")

	if !e {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthenticated!"})
		return
	}

	f, err := folder.GetFolderById(requestDto.ParentId, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Folder not found or you don't have permission to access it"})
		return
	}

	if exist, _ := folder.CheckFolderNameExist(requestDto.Name, f.ID, user.(models.User).ID); exist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Folder name already exists"})
		return
	}

	newFolder := models.Folder{
		Name:          requestDto.Name,
		UserId:        user.(models.User).ID,
		ParentID:      &f.ID,
		SpecialFolder: "",
	}

	if err := db.GetDB().Create(&newFolder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while creating folder DB record"})
		return
	}

	service.CreateNewNotification(user.(models.User).ID, "New folder '"+newFolder.Name+"' created successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "New folder created successfully!",
		"folder":  newFolder,
	})

}

func GetSubFolders(c *gin.Context) {
	parentId := c.Param("parentId")

	if parentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Parent Id is required",
		})
		return
	}

	user, _ := c.Get("currentUser")

	sub, err := folder.GetSubFolders(parentId, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching subfolders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folders": sub,
	})

}

func RenameFolder(c *gin.Context) {
	var requestDto dto.RenameFolderDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")

	f, err := folder.GetFolderById(requestDto.FolderId, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Folder not found or you don't have permission to access it"})
		return
	}

	if f.SpecialFolder != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You can't rename this folder"})
		return
	}

	f.Name = requestDto.Name
	if err := db.GetDB().Save(&f).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while renaming folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Folder renamed successfully!",
		"folder":  f,
	})

}

func DeleteFolder(c *gin.Context) {
	var requestDto dto.DeleteFolderDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")

	_, err := folder.GetFolderById(requestDto.FolderId, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Folder not found or you don't have permission to access it"})
		return
	}

	if err := folder.DeleteFilesAndFoldersInsideFolder(requestDto.FolderId, user.(models.User)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Folder deleted successfully!",
	})

}

func MoveFolder(c *gin.Context) {
	var requestDto dto.MoveFolderDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")

	_, err := folder.GetFolderById(requestDto.FolderId, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Folder not found or you don't have permission to access it"})
		return
	}

	_, err = folder.GetFolderById(requestDto.DestinationFolderID, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Folder not found or you don't have permission to access it"})
		return
	}

	if err := folder.MoveFolder(requestDto.FolderId, requestDto.DestinationFolderID, user.(models.User).ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Eror while moving folders!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Folder moved successfully!",
	})

}

func GetFolderItems(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id is required",
		})
		return
	}

	user, _ := c.Get("currentUser")

	folders, files, err := folder.GetFolderItems(id, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Eror while fetching files and folders!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files, "folders": folders})

}

func MoveToTrash(c *gin.Context) {
	var requestDto dto.MoveToTrashDTO

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")

	fs, err := folder.GetFolderById(requestDto.FolderId, user.(models.User).ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Folder not found or you don't have permission to access it"})
		return
	}

	fs.IsDeleted = true
	if err := db.GetDB().Save(&fs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot move file to the destination"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Folder moved to trash successfully!",
	})

}

func EmptyTrash(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)
	userID := user.(models.User).ID

	var deletedFolders []models.Folder
	if err := db.GetDB().Where("user_id = ? AND is_deleted = ?", userID, true).Find(&deletedFolders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving folders"})
		return
	}

	for _, _folder := range deletedFolders {
		if err := folder.DeleteFilesAndFoldersInsideFolder(_folder.ID, user.(models.User)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error emptying trash for folders"})
			return
		}
	}

	var deletedFiles []models.File
	if err := db.GetDB().Where("user_id = ? AND is_deleted = ?", userID, true).Find(&deletedFiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving files"})
		return
	}

	for _, file := range deletedFiles {
		if err := storage.FileDeleteService(file.FileId, &currentUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error emptying trash for files"})
			return
		}

		if err := db.GetDB().Unscoped().Where("id = ?", file.ID).Delete(&models.File{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting file record"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trash emptied successfully"})
}
