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

func StoreUploadedFile(mimeData string, fileData *FileUploadMetaData, user *models.User, passPhrase string) (*models.File, error) {

	fileId := uuid.New().String()

	privateKeyRaw, err := pki.DecryptPemFile(user.PrivateKeyPath, passPhrase)

	if err != nil {
		return &models.File{}, err
	}

	privateKey, err := pki.LoadPrivateKey(privateKeyRaw)

	if err != nil {
		return &models.File{}, err
	}

	signature, err := pki.SignFile(fileData.FilePath, privateKey)

	if err != nil {
		return &models.File{}, err
	}

	rawFileData, err := os.ReadFile(fileData.FilePath)

	if err != nil {
		return &models.File{}, err
	}

	encrypted_data, err := pki.Encrypt(rawFileData, []byte(passPhrase))

	if err != nil {
		return &models.File{}, err
	}

	uploadPath := "./disk/storage/" + fileId

	file, err := os.Create(uploadPath)
	if err != nil {
		return &models.File{}, err
	}
	defer file.Close()

	if _, err := file.Write([]byte(encrypted_data)); err != nil {
		return &models.File{}, err
	}

	newMetaData := models.File{
		FileId:        fileId,
		UserId:        user.ID,
		OriginalName:  fileData.Filename,
		FileSize:      fileData.Size,
		MimeType:      mimeData,
		UploadDate:    time.Now().UTC(),
		StoragePath:   uploadPath,
		FileSignature: signature,
		DownloadCount: 0,
		LastAccessed:  time.Now().UTC(),
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

func HandleDownloadProcess(fileId string, user *models.User, passPhrase string) (string, string, error) {

	// check the owner
	var gotFile models.File
	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", fileId).Where("user_id = ?", user.ID).First(&gotFile).Error; err != nil {
		return "", "", errors.New("permission denied")
	}

	// check the passPhrase
	if !ValidatePassPhrase(passPhrase, user) {
		return "", "", errors.New("invalid pass phrase")
	}

	// decrypt the file
	rawFileData, err := os.ReadFile(gotFile.StoragePath)

	if err != nil {
		return "", "", errors.New("could not find the file")
	}

	decryptedFileData, err := pki.Decrypt(string(rawFileData), []byte(passPhrase))

	if err != nil {
		return "", "", err
	}

	createFolder("./disk/public/" + gotFile.FileId)

	uploadPath := "./disk/public/" + gotFile.FileId + "/" + gotFile.OriginalName

	file, err := os.Create(uploadPath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	if _, err := file.Write(decryptedFileData); err != nil {
		return "", "", errors.New("cannot create decrypted file")
	}

	publicKey, err := pki.LoadPublicKey([]byte(user.PubKey))

	if err != nil {
		return "", "", errors.New("cannot load public key")
	}

	// check the file checksum
	if err := pki.VerifySign(uploadPath, gotFile.FileSignature, publicKey); err != nil {
		return "", "", errors.New("unauthorized filed alteration detected")
	}

	// increment download count
	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", gotFile.FileId).Update("download_count", gotFile.DownloadCount+1).Error; err != nil {
		return "", "", err
	}

	return uploadPath, gotFile.MimeType, nil

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
