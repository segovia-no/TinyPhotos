package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const COMPRESSED_FILENAME_TEXT = "_compressed"
const COMPRESSED_FOLDER_NAME = "compressed"

func GenerateCompressedFilename(filePath string) (string, error) {
	filepathSlice := strings.Split(filePath, ".")
	if len(filepathSlice) < 2 {
		return "", errors.New("missing extension")
	}

	extensionPos := len(filepathSlice) - 1
	if strings.ToLower(filepathSlice[extensionPos]) != "jpg" && strings.ToLower(filepathSlice[extensionPos]) != "jpeg" {
		return "", errors.New("this file is not an JPG or JPEG")
	}

	filepathNoExt := strings.Join(filepathSlice[:extensionPos], "")
	compressedFilename := filepathNoExt + COMPRESSED_FILENAME_TEXT + "." + filepathSlice[extensionPos]

	return compressedFilename, nil
}

func GetAllJPGFilePathsInFolder(folderPath string) ([]string, error) {
	if folderPath == "" {
		return []string{}, errors.New("folder path is empty")
	}

	if folderPath[len(folderPath)-1:] != string(filepath.Separator) {
		folderPath = folderPath + string(filepath.Separator)
	}

	filenames, err := os.ReadDir(folderPath)
	if err != nil {
		return []string{}, errors.New("Couldnt read folder path: " + err.Error())
	}

	jpegFilePaths := []string{}
	for _, entry := range filenames {
		if entry.IsDir() {
			continue
		}

		filepathSlice := strings.Split(entry.Name(), ".")
		if len(filepathSlice) < 2 {
			continue
		}

		extensionPos := len(filepathSlice) - 1
		if strings.ToLower(filepathSlice[extensionPos]) != "jpg" && strings.ToLower(filepathSlice[extensionPos]) != "jpeg" {
			continue
		}

		jpegFilePaths = append(jpegFilePaths, folderPath+entry.Name())
	}

	return jpegFilePaths, nil
}

func CreateCompressedForFolder(folderPath string) (string, error) {
	compressedFolderPath := filepath.Join(folderPath, COMPRESSED_FOLDER_NAME) + string(filepath.Separator)

	err := os.MkdirAll(compressedFolderPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return compressedFolderPath, nil
}
