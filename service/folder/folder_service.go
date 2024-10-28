package folder

import (
	"errors"
	"fmt"
	"peer-store/db"
	"peer-store/models"
	"peer-store/service/storage"

	"gorm.io/gorm"
)

func GetFolderById(folderId uint, userId uint) (models.Folder, error) {

	var folder models.Folder
	if err := db.GetDB().Where("id = ?", folderId).Where("user_id = ?", userId).First(&folder).Error; err != nil {
		return models.Folder{}, errors.New("folder not found")
	}

	return folder, nil

}

func GetSubFolders(folderId string, userId uint) ([]*models.Folder, error) {
	var folders []*models.Folder

	if err := db.GetDB().Model(&models.Folder{}).Where("user_id  = ?", userId).Where("parent_id = ?", folderId).Where("id <> ?", folderId).Find(&folders).Error; err != nil {
		return folders, err
	}

	return folders, nil

}

func GetParentFolder(folderId string, userId uint) (models.Folder, error) {
	var folder models.Folder
	if err := db.GetDB().Where("id = ?", folderId).Where("user_id = ?", userId).First(&folder).Error; err != nil {
		return models.Folder{}, errors.New("folder not found")
	}

	var parentFolder models.Folder
	if err := db.GetDB().Where("id =?", folder.ParentID).Where("user_id = ?", userId).First(&parentFolder).Error; err != nil {
		return models.Folder{}, errors.New("parent folder not found")
	}

	return parentFolder, nil
}

func CheckFolderNameExist(folderName string, parentId uint, userId uint) (bool, error) {

	var folder *models.Folder

	err := db.GetDB().Model(&models.Folder{}).
		Where("user_id = ?", userId).
		Where("parent_id = ?", parentId).
		Where("name = ?", folderName).
		First(&folder).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

func DeleteFilesAndFoldersInsideFolder(folderId uint, user models.User) error {
	
	var subfolders []models.Folder
	if err := db.GetDB().Where("user_id = ?", user.ID).Where("parent_id = ?", folderId).Find(&subfolders).Error; err != nil {
		return fmt.Errorf("error finding subfolders: %w", err)
	}

	for _, subfolder := range subfolders {
		if err := DeleteFilesAndFoldersInsideFolder(subfolder.ID, user); err != nil {
			return fmt.Errorf("error deleting subfolder ID %d: %w", subfolder.ID, err)
		}
	}

	var files []models.File
	if err := db.GetDB().Where("user_id = ?", user.ID).Where("folder_id = ?", folderId).Find(&files).Error; err != nil {
		return fmt.Errorf("error finding files: %w", err)
	}

	for _, file := range files {
		
		if err := storage.FileDeleteService(file.FileId, &user); err != nil {
			return fmt.Errorf("error deleting file ID %s: %w", file.FileId, err)
		}

		if err := db.GetDB().Unscoped().Where("id = ?", file.ID).Delete(&models.File{}).Error; err != nil {
			return fmt.Errorf("error removing file record for ID %d: %w", file.ID, err)
		}
	}

	if err := db.GetDB().Unscoped().Where("user_id = ?", user.ID).Where("id = ?", folderId).Delete(&models.Folder{}).Error; err != nil {
		return fmt.Errorf("error deleting folder ID %d: %w", folderId, err)
	}

	return nil
}

func MoveFolder(folderId uint, destination_folder_id uint, userId uint) error {

	fs, err := GetFolderById(folderId, userId)

	if err != nil {
		return err
	}

	fd, err := GetFolderById(destination_folder_id, userId)

	if err != nil {
		return err
	}

	fs.ParentID = &fd.ID
	if err := db.GetDB().Save(&fs).Error; err != nil {
		return err
	}

	return nil

}

func GetFolderItems(folderId string, userId uint) ([]*models.Folder, []*models.File, error) {

	var folders []*models.Folder
	var files []*models.File

	if err := db.GetDB().Model(&models.Folder{}).Where("user_id  = ?", userId).Where("is_deleted = ?", false).Where("parent_id = ?", folderId).Where("id <> ?", folderId).Find(&folders).Error; err != nil {
		return folders, files, err
	}

	if err := db.GetDB().Model(&models.File{}).Where("user_id  = ?", userId).Where("folder_id = ?", folderId).Where("is_deleted = ?", false).Find(&files).Error; err != nil {
		return folders, files, err
	}

	return folders, files, nil

}
