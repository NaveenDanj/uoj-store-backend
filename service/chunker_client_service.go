package service

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"
	"peer-store/config"
	core "peer-store/core/file_manager"
	"peer-store/core/pki"
	"peer-store/core/types"
)

func LoadRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {

	private_pem_file_enc, err := pki.DecryptPemFile(config.CONFIG.PrivatePEMFilePath+".enc", config.CONFIG.PassPhrase)

	if err != nil {
		return nil, nil, fmt.Errorf("Error while decrypting PEM files. Please check your pass-prase is correct : " + err.Error())
	}

	public_pem_file_enc, err := pki.DecryptPemFile(config.CONFIG.PublicPEMFilePath+".enc", config.CONFIG.PassPhrase)

	if err != nil {
		return nil, nil, fmt.Errorf("Error while decrypting PEM files. Please check your pass-prase is correct : " + err.Error())
	}

	// decrypt them
	private_pem_file, err := pki.LoadPrivateKey(private_pem_file_enc)

	if err != nil {
		return nil, nil, fmt.Errorf("Error while decrypting PEM files 1 : " + err.Error())
	}

	public_pem_file, err := pki.LoadPublicKey(public_pem_file_enc)

	if err != nil {
		return nil, nil, fmt.Errorf("Error while decrypting PEM files 2 : " + err.Error())
	}

	return private_pem_file, public_pem_file, nil

}

func UploadFileAsChunk(filepath string) error {
	// load pem files
	private_key, _, err := LoadRSAKeys()

	if err != nil {
		return fmt.Errorf("Error while decrypting PEM files. Please check your pass-prase is correct : " + err.Error())
	}

	// chunk the files
	file_chunk_data, err := core.FileSpliterService(filepath, 10, config.CONFIG.FileCWD)

	if err != nil {
		return fmt.Errorf("Error while file chunking : " + err.Error())
	}

	for _, path := range file_chunk_data.Sequence {

		// sign them
		data, err := pki.SignFile(path, private_key)

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		file_chunk_data.DigitalSignature = append(file_chunk_data.DigitalSignature, string(data))

		// encrypt them
		chunk_file_data, err := os.ReadFile(path)

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		encrypted_data, err := pki.Encrypt(chunk_file_data, []byte(config.CONFIG.PassPhrase))

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		partFile, err := os.Create(path)

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		partFile.Write(encrypted_data)

	}

	fileUploadDTO := types.FileUploadDTO{
		FileChunkMetaData: &file_chunk_data,
		Sequence:          file_chunk_data.Sequence,
	}

	outData, err := json.Marshal(fileUploadDTO)

	if err != nil {
		fmt.Errorf("Error converting struct to JSON: %v", err)
	}

	out, err := pki.Encrypt(outData, []byte(config.CONFIG.PassPhrase))

	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Println("Out string for server is -> ", string(out))

	// call server to store file chunks in the main-net and return file-meta data and main-net router map

	return nil

}
