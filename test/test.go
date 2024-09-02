package main

import (
	"fmt"
	"os"
	core "peer-store/core/file_manager"
	"peer-store/core/pki"
	"peer-store/service"
)

func main() {

	Pki_test()

	err := service.UploadFileAsChunk("./test/test.pdf")

	if err != nil {
		fmt.Println(err.Error())
	}

	// FileEncryptTest()
	// file_chunk_test()

	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	panic("Error while reading the file path. Please check your file path")
	// }
	// baseDir := currentDir + "/test/"

	// core.ChunkerRollbackService(baseDir)
}

func Pki_test() {
	pki.Generate_pki_key_pair()
}

func File_chunk_test() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error while reading the file path. Please check your file path")
	}
	baseDir := currentDir + "./test/"
	file_path_string := baseDir + "Excel excecises 2021_csc_019.xlsx"
	l, err := core.FileSpliterService(file_path_string, 10, baseDir)

	if err != nil {
		fmt.Println("Error while chunking files : " + err.Error())
	}

	err = core.FileMerger(l.Sequence, l)

	if err != nil {
		fmt.Errorf(err.Error())
	}

}

func SignFileTest() {

}

// func FileEncryptTest() {
// 	// get public key
// 	pub_key, err := pki.LoadPublicKey("./public_key.pem")

// 	if err != nil {
// 		fmt.Errorf("Error while reading public key pem file")
// 	}

// 	out, err := pki.PublicKeyEncryption("hello world", pub_key)

// 	if err != nil {
// 		fmt.Errorf(err.Error())
// 	}

// 	fmt.Println(string(out))

// 	fmt.Println("---------------------------")

// 	pri_key, err := pki.LoadPrivateKey("./private_key.pem")

// 	if err != nil {
// 		fmt.Errorf(err.Error())
// 	}

// 	decrypted, err := pki.DecryptWithPrivateKey(out, pri_key)

// 	if err != nil {
// 		fmt.Errorf(err.Error())
// 	}

// 	fmt.Println("output is => ", decrypted)

// }
