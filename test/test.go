package main

import (
	"fmt"
	"os"
	"peer-store/services"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error while reading the file path. Please check your file path")
	}
	baseDir := currentDir + "/test/"
	file_path_string := baseDir + "test.pdf"
	l, err := services.FileSpliterService(file_path_string, 10, baseDir)

	if err != nil {
		fmt.Println("Error while chunking files : " + err.Error())
	}

	fmt.Println(l)

}
