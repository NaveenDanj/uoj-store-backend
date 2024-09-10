package storage

import (
	"errors"
	"os"
	"peer-store/core/pki"
	"peer-store/db"
	"peer-store/dto"
	"peer-store/models"
	"time"

	"github.com/google/uuid"
)

func GenerateShareLink(fileShareDetails *dto.FileShareRequestDTO, user *models.User) (string, error) {

	token := uuid.New().String()
	encryption_key := uuid.New().String()
	encryption_key = encryption_key[:32]

	var requestFile models.File

	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", fileShareDetails.FileId).First(&requestFile).Error; err != nil {
		return "", err
	}

	if user.ID != requestFile.UserId {
		return "", errors.New("permission denied")
	}

	if !ValidatePassPhrase(fileShareDetails.PassPhrase, user) {
		return "", errors.New("invalid pass phrase")
	}

	rawFileData, err := os.ReadFile(requestFile.StoragePath)

	if err != nil {
		return "", errors.New("could not find the file")
	}

	decryptedFileData, err := pki.Decrypt(string(rawFileData), []byte(fileShareDetails.PassPhrase))

	if err != nil {
		return "", err
	}

	out, err := pki.Encrypt(decryptedFileData, []byte(encryption_key))

	if err != nil {
		return "", err
	}

	uploadPath := "./disk/shared/" + requestFile.FileId

	file, err := os.Create(uploadPath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err := file.Write([]byte(out)); err != nil {
		return "", errors.New("cannot create encrypted file")
	}

	expireDate, err := time.Parse(time.RFC3339, fileShareDetails.ExpireDate)

	if err != nil {
		return "", err
	}

	fileShare := models.FileShare{
		Token:         token,
		FileId:        requestFile.FileId,
		OwnerId:       user.ID,
		IsPublic:      false,
		ExpireDate:    expireDate,
		Status:        "Shared",
		Note:          fileShareDetails.Note,
		SharedAt:      time.Now().UTC(),
		DownloadLimit: fileShareDetails.DownloadLimit,
		EncryptionKey: encryption_key,
	}

	if err := db.GetDB().Create(&fileShare).Error; err != nil {
		return "", err
	}

	for _, user := range fileShareDetails.Users {
		userShare := models.FileShareUser{
			UserId:        user.UserId,
			FileShareId:   uint(fileShare.Id),
			DownloadCount: 0,
		}

		if err := db.GetDB().Create(&userShare).Error; err != nil {
			return "", err
		}

	}

	return token, nil

}

func RevokeLink(shareId string, fileId string, user *models.User) error {

	var requestFile *models.File

	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", fileId).First(&requestFile).Error; err != nil {
		return err
	}

	if user.ID != requestFile.UserId {
		return errors.New("permission denied")
	}

	if err := db.GetDB().Model(&models.FileShare{}).Where("id = ?", shareId).Update("status", "Revoked").Error; err != nil {
		return err
	}

	return nil
}
