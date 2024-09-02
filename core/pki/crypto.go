package pki

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func EncryptPemFile(filepath string, outputFilePath string, passphrase string) error {

	file, err := os.ReadFile(filepath)

	if err != nil {
		return fmt.Errorf("error in file opening pem file")
	}

	key := genKeyFromPassPhrase(passphrase)

	encryptedData, err := encryptAESGCM(key, file)
	if err != nil {
		return fmt.Errorf("failed to encrypt PEM file: %w", err)
	}

	err = os.WriteFile(outputFilePath, encryptedData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write encrypted PEM file: %w", err)
	}

	return nil

}

func DecryptPemFile(filepath string, passPhrase string) ([]byte, error) {
	encryptedData, err := os.ReadFile(filepath)

	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted PEM file: %w", err)
	}

	key := genKeyFromPassPhrase(passPhrase)

	decryptedData, err := decryptAESGCM(key, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt PEM file: %w", err)
	}

	return decryptedData, nil

}

func genKeyFromPassPhrase(passPhrase string) []byte {
	key := make([]byte, 32)
	copy(key, []byte(passPhrase))
	return key
}

func encryptAESGCM(key []byte, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plainText, nil)
	return ciphertext, nil
}

func decryptAESGCM(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}