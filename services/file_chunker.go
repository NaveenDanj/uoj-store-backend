package services

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileChunkMetaData struct {
	Sequence  []string
	FileInfo  fs.FileInfo
	Extension string
}

func FileSpliterService(path string, numberOfChunks int64, outDir string) (FileChunkMetaData, error) {

	file, err := os.Open(path)
	chunckFileList := make([]string, 0)

	if err != nil {
		return FileChunkMetaData{}, fmt.Errorf("error reading splitting file data")
	}

	fileInfo, err := file.Stat()

	defer file.Close()

	if err != nil {
		return FileChunkMetaData{}, fmt.Errorf("failed to get file info: %w", err)
	}

	partSize := fileInfo.Size() / numberOfChunks
	baseName := uuid.New().String()
	buffer := make([]byte, partSize)
	fileExtension := filepath.Ext(file.Name())

	for i := int64(0); i < int64(numberOfChunks); i++ {

		bytesRead, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			return FileChunkMetaData{}, fmt.Errorf("failed to read file: %w", err)
		}

		partFileName := fmt.Sprintf("%s%d.part", baseName, i+1)
		partFilePath := filepath.Join(outDir, partFileName)
		partFile, err := os.Create(partFilePath)
		chunckFileList = append(chunckFileList, partFilePath)

		if err != nil {
			return FileChunkMetaData{}, fmt.Errorf("failed to create part file: %w", err)
		}

		_, err = partFile.Write(buffer[:bytesRead])
		if err != nil {
			partFile.Close()
			return FileChunkMetaData{}, fmt.Errorf("failed to write to part file: %w", err)
		}

	}

	chunkData := FileChunkMetaData{
		Sequence:  chunckFileList,
		FileInfo:  fileInfo,
		Extension: fileExtension,
	}

	return chunkData, nil

}

func ChunkerRollbackService() {
	return
}
