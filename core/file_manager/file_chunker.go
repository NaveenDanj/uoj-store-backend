package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"peer-store/core/types"

	"github.com/google/uuid"
)

func FileSpliterService(path string, numberOfChunks int64, outDir string) (types.FileChunkMetaData, error) {

	file, err := os.Open(path)
	chunckFileList := make([]string, 0)

	if err != nil {
		return types.FileChunkMetaData{}, fmt.Errorf("error reading splitting file data")
	}

	fileInfo, err := file.Stat()

	defer file.Close()

	if err != nil {
		return types.FileChunkMetaData{}, fmt.Errorf("failed to get file info: %w", err)
	}

	partSize := fileInfo.Size() / numberOfChunks
	original_name := fileInfo.Name()
	baseName := uuid.New().String()
	buffer := make([]byte, partSize)
	fileExtension := filepath.Ext(file.Name())

	for i := int64(0); i < int64(numberOfChunks); i++ {

		bytesRead, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			return types.FileChunkMetaData{}, fmt.Errorf("failed to read file: %w", err)
		}

		partFileName := fmt.Sprintf("%s%d.part", baseName, i+1)
		partFilePath := filepath.Join(outDir, partFileName)
		partFile, err := os.Create(partFilePath)
		chunckFileList = append(chunckFileList, partFilePath)

		if err != nil {
			return types.FileChunkMetaData{}, fmt.Errorf("failed to create part file: %w", err)
		}

		_, err = partFile.Write(buffer[:bytesRead])
		if err != nil {
			partFile.Close()
			return types.FileChunkMetaData{}, fmt.Errorf("failed to write to part file: %w", err)
		}

	}

	chunkData := types.FileChunkMetaData{
		Sequence:     chunckFileList,
		FileInfo:     fileInfo,
		Extension:    fileExtension,
		OriginalName: original_name,
	}

	return chunkData, nil

}

func FileMerger(file_chunk_list []string, metaData types.FileChunkMetaData) error {

	// create new file for the merged file
	newFile, err := os.Create(metaData.OriginalName)

	if err != nil {
		return err
	}

	defer newFile.Close()

	for _, path := range file_chunk_list {

		chunkFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open chunk file %s: %w", path, err)
		}

		_, err = io.Copy(newFile, chunkFile)
		if err != nil {
			chunkFile.Close() // Ensure chunk file is closed before returning
			return fmt.Errorf("failed to write chunk %s to output file: %w", path, err)
		}

		chunkFile.Close()

	}

	return nil

}

func ChunkerRollbackService(base_dir string) {

	cwd := base_dir

	err := filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has the desired extension
		if !info.IsDir() && filepath.Ext(path) == ".part" {
			fmt.Println(path)
			err := os.Remove(path)

			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

}
