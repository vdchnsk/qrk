package fs

import (
	"os"
	"strings"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	isFile := !info.IsDir()

	return isFile
}

func CanRunFile(path string) bool {
	fileExtension := ".qrk"

	isSupportedExtension := strings.HasSuffix(path, fileExtension)

	if !isSupportedExtension {
		return false
	}

	fileExists := FileExists(path)

	if !fileExists {
		return fileExists
	}

	return true
}
