package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"peer-store/core/pki"
	"peer-store/db"
	"peer-store/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type FileUploadMetaData struct {
	Filename     string
	Size         int64
	TempFolder   string
	UploadedFile *os.File
	FilePath     string
}

func ValidatePassPhrase(passPhrase string, user *models.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PassPhrase), []byte(passPhrase))
	return err == nil
}

func FileUploader(file multipart.File, header *multipart.FileHeader) (FileUploadMetaData, error) {

	tempFolder := "./disk/chunker/" + uuid.New().String()

	err := os.MkdirAll(tempFolder, os.ModePerm)

	if err != nil {
		return FileUploadMetaData{}, err
	}

	dst, err := os.Create(filepath.Join(tempFolder, header.Filename))
	if err != nil {
		return FileUploadMetaData{}, err
	}

	defer dst.Close()

	buf := make([]byte, 1*1024*1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return FileUploadMetaData{}, err
		}
		if n == 0 {
			break
		}

		if _, err := dst.Write(buf[:n]); err != nil {
			return FileUploadMetaData{}, err
		}

	}

	fileUploadMetaData := FileUploadMetaData{
		Filename:     header.Filename,
		Size:         header.Size,
		TempFolder:   tempFolder,
		UploadedFile: dst,
		FilePath:     tempFolder + "/" + header.Filename,
	}

	return fileUploadMetaData, nil

}

func StoreUploadedFile(mimeData string, fileData *FileUploadMetaData, user *models.User, passPhrase string, folder_id uint, shouldEncrypt bool) (*models.File, error) {
	fileId := uuid.New().String()

	// Initialize variables for conditional signature and encryption
	var signature string
	var encrypted_data string
	var err error

	// Encrypt and sign only if shouldEncrypt is true
	if shouldEncrypt {
		// Load private key and sign file
		privateKeyRaw, err := pki.DecryptPemFile(user.PrivateKeyPath, passPhrase)
		if err != nil {
			return &models.File{}, err
		}

		privateKey, err := pki.LoadPrivateKey(privateKeyRaw)
		if err != nil {
			return &models.File{}, err
		}

		dat, err := pki.SignFile(fileData.FilePath, privateKey)
		signature = string(dat)
		if err != nil {
			return &models.File{}, err
		}

		// Read and encrypt file data
		rawFileData, err := os.ReadFile(fileData.FilePath)
		if err != nil {
			return &models.File{}, err
		}

		encrypted_data, err = pki.Encrypt(rawFileData, []byte(passPhrase))
		if err != nil {
			return &models.File{}, err
		}
	} else {
		// Only read file data if encryption is not required
		data, err := os.ReadFile(fileData.FilePath)
		encrypted_data = string(data)
		if err != nil {
			return &models.File{}, err
		}
	}

	// Create and save the file
	uploadPath := "./disk/storage/" + fileId
	file, err := os.Create(uploadPath)
	if err != nil {
		return &models.File{}, err
	}
	defer file.Close()

	if _, err := file.Write([]byte(encrypted_data)); err != nil {
		return &models.File{}, err
	}

	// Set up metadata and save to DB
	newMetaData := models.File{
		FolderID:      folder_id,
		FileId:        fileId,
		UserId:        user.ID,
		OriginalName:  fileData.Filename,
		FileSize:      fileData.Size,
		MimeType:      mimeData,
		UploadDate:    time.Now().UTC(),
		StoragePath:   uploadPath,
		FileSignature: []byte(signature),
		DownloadCount: 0,
		LastAccessed:  time.Now().UTC(),
	}

	if !shouldEncrypt {
		newMetaData.IsPublic = true
	}

	if err := db.GetDB().Model(models.File{}).Create(&newMetaData).Error; err != nil {
		return &models.File{}, err
	}

	DeleteFolder(fileData.TempFolder)

	return &newMetaData, nil
}

func createFolder(folderPath string) error {
	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %w", err)
		}
		fmt.Println("Folder created:", folderPath)
	} else {
		fmt.Println("Folder already exists:", folderPath)
	}
	return nil
}

func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	return nil

}

func DeleteFolder(folderPath string) error {
	err := os.RemoveAll(folderPath)
	if err != nil {
		return err
	}
	return nil

}

func HandleDownloadProcess(fileId string, user *models.User, passPhrase string, shouldDecrypt bool) (string, string, error) {

	// check the owner
	var gotFile models.File
	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", fileId).Where("user_id = ?", user.ID).First(&gotFile).Error; err != nil {
		return "", "", errors.New("permission denied")
	}

	// check the passPhrase
	if shouldDecrypt {
		if !ValidatePassPhrase(passPhrase, user) {
			return "", "", errors.New("invalid pass phrase")
		}
	}

	path, mime, _, err := UtilDecryptAndUse(&gotFile, gotFile.StoragePath, []byte(passPhrase), user, shouldDecrypt)

	if err != nil {
		return "", "", errors.New("unauthorized filed alteration detected")
	}
	// increment download count
	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", gotFile.FileId).Update("download_count", gotFile.DownloadCount+1).Error; err != nil {
		return "", "", err
	}

	return path, mime, nil

}

func FileDeleteService(fileId string, user *models.User) error {

	// check the owner of the file
	var file *models.File

	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", fileId).First(&file).Error; err != nil {
		return err
	}

	if file.UserId != user.ID {
		return errors.New("permission denied")
	}

	// delete the file from the database
	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", fileId).Delete(&file).Error; err != nil {
		return err
	}

	// delete the file from the disk
	if err := DeleteFile(file.StoragePath); err != nil {
		return err
	}

	return nil

}

func GetUserFiles(userId uint) ([]*models.File, error) {

	var files []*models.File

	if err := db.GetDB().Model(&models.File{}).Where("user_id  = ?", userId).Find(&files).Error; err != nil {
		return files, err
	}

	return files, nil

}
