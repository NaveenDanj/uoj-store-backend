package storage

import (
	"errors"
	"os"
	"peer-store/core/pki"
	"peer-store/models"
)

func UtilDecryptAndUse(gotFile *models.File, filepath string, key []byte, user *models.User, shouldDecrypt bool) (string, string, string, error) {

	rawFileData, err := os.ReadFile(gotFile.StoragePath)
	if err != nil {
		return "", "", "", errors.New("could not find the file")
	}

	var fileData []byte
	if shouldDecrypt {
		fileData, err = pki.Decrypt(string(rawFileData), key)
		if err != nil {
			return "", "", "", err
		}
	} else {
		fileData = rawFileData
	}

	createFolder("./disk/public/" + gotFile.FileId)
	uploadPath := "./disk/public/" + gotFile.FileId + "/" + gotFile.OriginalName

	file, err := os.Create(uploadPath)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	if _, err := file.Write(fileData); err != nil {
		return "", "", "", errors.New("cannot create output file")
	}

	if shouldDecrypt {
		publicKey, err := pki.LoadPublicKey([]byte(user.PubKey))
		if err != nil {
			return "", "", "", errors.New("cannot load public key")
		}

		if err := pki.VerifySign(uploadPath, gotFile.FileSignature, publicKey); err != nil {
			return "", "", "", errors.New("unauthorized file alteration detected")
		}
	}

	return uploadPath, gotFile.MimeType, "./disk/public/" + gotFile.FileId, nil
}
