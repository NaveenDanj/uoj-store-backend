package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func GeneratePkiKeyPair(passPhrase string) (string, string, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("error generating private key: %s", err)
	}

	publicKey := &privateKey.PublicKey

	privatePemFilePath := "./disk/keys/" + uuid.New().String() + ".pem"

	if err := SavePEMKey(privatePemFilePath, privateKey); err != nil {
		return "", "", fmt.Errorf("error generating private key: %s", err.Error())
	}

	if err := EncryptPemFile(privatePemFilePath, privatePemFilePath+".enc", passPhrase); err != nil {
		return "", "", fmt.Errorf("error encrypting private key: %s", err.Error())
	}

	pubKey, err := PublicKeyToBytes(publicKey)

	if err != nil {
		return "", "", fmt.Errorf("error encrypting private key: %s", err.Error())
	}

	if err := os.Remove(privatePemFilePath); err != nil {
		return "", "", fmt.Errorf("error deleting pem file: %s", err.Error())
	}

	return privatePemFilePath + ".enc", string(pubKey), nil

}

func LoadPrivateKey(keyData []byte) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil

}

func LoadPublicKey(keyData []byte) (*rsa.PublicKey, error) {

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPubKey, nil

}

func SavePEMKey(filename string, key *rsa.PrivateKey) error {

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file:", err)
	}
	defer file.Close()

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	_, err = file.Write(privateKeyPEM)

	if err != nil {
		return fmt.Errorf("error writing private key to file:", err)
	}

	return nil

}

func PublicKeyToBytes(key *rsa.PublicKey) ([]byte, error) {

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, fmt.Errorf("error marshalling public key:", err)
	}

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	return publicKeyPEM, nil

}
