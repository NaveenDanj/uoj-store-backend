package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func SignFile(filepath string, privateKey *rsa.PrivateKey) ([]byte, error) {

	file_data, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	defer file_data.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file_data); err != nil {
		return nil, fmt.Errorf("failed to hash file content: %w", err)
	}

	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, fmt.Errorf("failed to sign file: %w", err)
	}

	return signature, nil

}

func VerifySign(filepath string, signature []byte, publicKey *rsa.PublicKey) error {

	file, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	hash := sha256.New()
	// Read the file content and write it to the hash
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to hash file content: %w", err)
	}

	hashed := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, signature)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil

}
