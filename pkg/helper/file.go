package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindFile(filename string) (string, error) {
	// start searching in the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// recursively search for the file by going upwards one level
	for {
		// check if the file exists in the current directory
		filePath := filepath.Join(currentDir, filename)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}

		// go up one level
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// reached the root directory, file not found
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("file not found: '%s'", filename)
}
