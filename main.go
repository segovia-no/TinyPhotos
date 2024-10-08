package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var flags Flags
var tinifyClient TinifyClient

func main() {
	flags.parseFlags()

	if flags.log {
		logFile := setupLogToFile()
		defer logFile.Close()
	}

	apikey := flags.apikey
	if apikey == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading Tinify API key from .env file")
		}
		apikey = os.Getenv("TINIFY_API_KEY")
	}

	err := tinifyClient.SetAPIKey(apikey)
	if err != nil {
		log.Fatal("Couldnt set Tinify API key: " + err.Error())
	}

	if flags.file == "" && flags.bulkFromFolder == "" {
		log.Println("Use the -help flag to see how to use this program")
		os.Exit(0)
	}

	if flags.file != "" {
		processSingleFile(flags.file)
	}

	if flags.bulkFromFolder != "" {
		processFolder(flags.bulkFromFolder)
	}
}

func processSingleFile(filePath string) {
	startTime := time.Now()
	log.Println("Starting conversion of file: " + filePath)

	compressedFilename, err := GenerateCompressedFilename(filePath)
	if err != nil {
		log.Fatal("There was an issue with the provided filename: " + err.Error())
	}

	log.Println("Compressing file: " + filePath)
	tinifyResponse, err := tinifyClient.MakeRequest("/shrink", filePath)
	if err != nil {
		log.Fatal("Couldnt convert file: " + err.Error())
	}

	log.Println("Downloading compressed image: " + tinifyResponse.Headers.Location)
	err = tinifyClient.DownloadWithMetadata(tinifyResponse.Headers.Location, compressedFilename)
	if err != nil {
		log.Fatal("Couldnt download compressed image: " + err.Error())
	}

	log.Println("Writing metadata back to compressed image: " + compressedFilename)
	err = CopyExifMetadata(filePath, compressedFilename)
	if err != nil {
		log.Fatal("Coudlnt write metadata to compressed file: " + err.Error())
	}

	duration := time.Since(startTime)
	log.Printf("Done!, processing finished in %v", duration)
}

func processFolder(folderPath string) {
	log.Println("Starting conversion of folder: " + folderPath)

	jpegFilePaths, err := GetAllJPGFilePathsInFolder(folderPath)
	if err != nil {
		log.Fatal("Couldn't list files inside the requested folder path: " + err.Error())
	}

	compressedFolderPath, err := CreateCompressedForFolder(folderPath)
	if err != nil {
		log.Fatal("Couldn't create the target folder for the compressed files: ", err.Error())
	}

	if flags.maxRoutines < 2 {
		processFolderSync(jpegFilePaths, compressedFolderPath)
	} else {
		processFolderConcurrently(jpegFilePaths, compressedFolderPath)
	}

	log.Println("Done!")
}

func processFolderSync(jpegFilePaths []string, compressedFolderPath string) {
	totalJpegs := len(jpegFilePaths)
	for idx, fpath := range jpegFilePaths {
		processFileForFolder(fpath, compressedFolderPath, idx, totalJpegs)
	}
}

func processFolderConcurrently(jpegFilePaths []string, compressedFolderPath string) {
	var wg = sync.WaitGroup{}
	totalJpegs := len(jpegFilePaths)

	maxGoroutines := flags.maxRoutines
	if totalJpegs < maxGoroutines {
		maxGoroutines = totalJpegs
	}
	guard := make(chan struct{}, maxGoroutines)

	for i := 0; i < totalJpegs; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(fpath string, compressedFolderPath string, idx int, totalJpegs int) {
			processFileForFolder(jpegFilePaths[i], compressedFolderPath, i, totalJpegs)
			<-guard
			wg.Done()
		}(jpegFilePaths[i], compressedFolderPath, i, totalJpegs)
	}
	wg.Wait()
}

func processFileForFolder(fpath string, compressedFolderPath string, idx int, totalJpegs int) {
	startTime := time.Now()
	_, fname := filepath.Split(fpath)
	compressedFilePath := compressedFolderPath + fname

	log.Printf("[%d/%d] Starting processing of %s\n", idx+1, totalJpegs, fname)

	var compAttempt uint8 = 0
	compSucceeded := false
	for !compSucceeded {
		compAttempt++
		if compAttempt > uint8(flags.maxRetries) + 1 {
			log.Printf("[Failed %d/%d] Compression max attemps exhausted for: %s \n", idx+1, totalJpegs, fname)
			return
		}

		log.Printf("[%d/%d] Compressing file: %s - attempt %d \n", idx+1, totalJpegs, fname, compAttempt)
		tinifyResponse, err := tinifyClient.MakeRequest("/shrink", fpath)
		if err != nil {
			log.Printf("[Warning %d/%d] Couldnt convert file: %s \n", idx+1, totalJpegs, err.Error())
			continue
		}

		log.Printf("[%d/%d] Downloading compressed image: %s - attempt %d \n", idx+1, totalJpegs, fname, compAttempt)
		err = tinifyClient.DownloadWithMetadata(tinifyResponse.Headers.Location, compressedFilePath)
		if err != nil {
			log.Printf("[Warning %d/%d] Couldnt download compressed image: %s \n", idx+1, totalJpegs, err.Error())
			continue
		}

		compSucceeded = true
	}

	log.Printf("[%d/%d] Writing metadata back to compressed image: %s \n", idx+1, totalJpegs, fname)
	err := CopyExifMetadata(fpath, compressedFilePath)
	if err != nil {
		log.Printf("[Warning %d/%d] Coudlnt write metadata to compressed file: %s , tagging with _nometadata filename \n", idx+1, totalJpegs, err.Error())
		NoMetadataRenaming(compressedFilePath)
	}

	duration := time.Since(startTime)
	log.Printf("[%d/%d] Finished processing for image %s (took %v) \n", idx+1, totalJpegs, fname, duration)
}
