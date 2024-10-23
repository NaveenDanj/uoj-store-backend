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

	expireDate, err := time.Parse(time.RFC3339, fileShareDetails.ExpireDate)

	if err != nil {
		return "", err
	}

	today := time.Now()

	if expireDate.Compare(today) < 0 {
		return "", errors.New("invalid expire date")
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

func RevokeLink(token string, user *models.User) error {

	var fileShare *models.FileShare

	if err := db.GetDB().Model(&models.FileShare{}).Where("Token = ?", token).First(&fileShare).Error; err != nil {
		return err
	}

	if fileShare.Status == "Revoked" {
		return errors.New("link already revoked")
	}

	var requestFile *models.File

	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", fileShare.FileId).First(&requestFile).Error; err != nil {
		return err
	}

	if user.ID != requestFile.UserId {
		return errors.New("permission denied")
	}

	if err := db.GetDB().Model(&models.FileShare{}).Where("id = ?", fileShare.Id).Update("status", "Revoked").Error; err != nil {
		return err
	}

	if err := DeleteFile("./disk/shared/" + requestFile.FileId); err != nil {
		return err
	}

	return nil

}

func DownloadSharedFile(token string, user *models.User) (string, string, string, error) {

	var sharedFile *models.FileShare

	if err := db.GetDB().Model(&models.FileShare{}).Where("token = ?", token).Find(&sharedFile).Error; err != nil {
		return "", "", "", err
	}

	var file *models.File

	if err := db.GetDB().Model(&models.File{}).Where("file_id = ?", sharedFile.FileId).Find(&file).Error; err != nil {
		return "", "", "", err
	}

	// check expiry date
	if sharedFile.Status == "Revoked" || sharedFile.ExpireDate.Before(time.Now().UTC()) {
		return "", "", "", errors.New("link is expired or revoked")
	}

	// check file permission
	var userShare *models.FileShareUser

	if err := db.GetDB().Model(&models.FileShareUser{}).Where("user_id =?", user.ID).Where("file_share_id =?", sharedFile.Id).First(&userShare).Error; err != nil {
		return "", "", "", err
	}

	// check download count
	if userShare.DownloadCount >= sharedFile.DownloadLimit {
		return "", "", "", errors.New("download limit reached")
	}

	userShare.DownloadCount++

	if err := db.GetDB().Save(userShare).Error; err != nil {
		return "", "", "", err
	}

	// decrypt the file
	rawFileData, err := os.ReadFile("./disk/shared/" + file.FileId)

	if err != nil {
		return "", "", "", errors.New("could not find the file")
	}

	decryptedFileData, err := pki.Decrypt(string(rawFileData), []byte(sharedFile.EncryptionKey))

	if err != nil {
		return "", "", "", err
	}

	createFolder("./disk/public/" + file.FileId)

	uploadPath := "./disk/public/" + file.FileId + "/" + file.OriginalName

	newFile, err := os.Create(uploadPath)

	if err != nil {
		return "", "", "", err
	}
	defer newFile.Close()

	if _, err := newFile.Write(decryptedFileData); err != nil {
		return "", "", "", errors.New("cannot create decrypted file")
	}

	// publicKey, err := pki.LoadPublicKey([]byte(user.PubKey))
	_, err = pki.LoadPublicKey([]byte(user.PubKey))

	if err != nil {
		return "", "", "", errors.New("cannot load public key")
	}

	// check the file checksum
	// if err := pki.VerifySign(uploadPath, file.FileSignature, publicKey); err != nil {
	// 	return "", "", "", errors.New("unauthorized file alteration detected")
	// }

	return uploadPath, file.MimeType, file.FileId, nil

}
