package main

import (
	"os/exec"
)

func CopyExifMetadata(sourceFile string, targetFile string) error {
	cmd := exec.Command("exiftool", "-overwrite_original", "-TagsFromFile", sourceFile, "-all:all>all:all", targetFile)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
