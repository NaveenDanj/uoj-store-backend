package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func Generate_pki_key_pair() error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Errorf("error generating private key:", err)
	}

	publicKey := &privateKey.PublicKey

	err = SavePEMKey("private_key.pem", privateKey)

	if err != nil {
		fmt.Errorf("error generating private key:", err.Error())
	}

	err = SavePublicPEMKey("public_key.pem", publicKey)

	if err != nil {
		fmt.Errorf("error generating private key:", err.Error())
	}

	return nil

}

func LoadPrivateKey(filepath string) (*rsa.PrivateKey, error) {

	keyData, err := os.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

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

func LoadPublicKey(filename string) (*rsa.PublicKey, error) {

	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

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

func SavePublicPEMKey(filename string, key *rsa.PublicKey) error {

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err.Error())
	}
	defer file.Close()

	pubkeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("error marshalling public key:", err)

	}

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	_, err = file.Write(publicKeyPEM)

	if err != nil {
		return fmt.Errorf("error writing public key to file:", err)

	}

	return nil

}
