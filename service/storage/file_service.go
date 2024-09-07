package storage

import (
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
	fmt.Println("--------------------------")
	fmt.Println(user.PassPhrase)
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

	if _, err := file.Write(encrypted_data); err != nil {
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

	return &newMetaData, nil

}
