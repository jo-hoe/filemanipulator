package filemanipulator

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
)

func calculateFileHash(filePath string) (result []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return result, err
	}
	result = hasher.Sum(nil)
	return result, err
}

func AreFileEqual(leftFilePath string, rightFilePath string) (bool, error) {
	leftFileHash, err := calculateFileHash(leftFilePath)
	if err != nil {
		return false, err
	}
	rightFileHash, err := calculateFileHash(rightFilePath)
	if err != nil {
		return false, err
	}
	return bytes.Equal(leftFileHash, rightFileHash), nil
}
