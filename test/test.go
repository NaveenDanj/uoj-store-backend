package main

import (
	"fmt"
	"os"
	core "peer-store/core/file_manager"
	"peer-store/core/pki"
)

func main() {
	// pki_test()
	// file_chunk_test()

	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error while reading the file path. Please check your file path")
	}
	baseDir := currentDir + "/test/"

	core.ChunkerRollbackService(baseDir)
}

func pki_test() {
	pki.Generate_pki_key_pair()
}

func file_chunk_test() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error while reading the file path. Please check your file path")
	}
	baseDir := currentDir + "/test/"
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
