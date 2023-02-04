package filemanipulator

import (
	"errors"
	"io"
	"os"
)

type MountedFileHandler struct {
}

func (m *MountedFileHandler) DoesFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func (m *MountedFileHandler) Open(filePath string) (reader io.ReadWriteCloser, err error) {
	return os.Open(filePath)
}

func (m *MountedFileHandler) Remove(filePath string) (err error) {
	return os.Remove(filePath)
}
func (m *MountedFileHandler) Create(filePath string) (reader io.ReadWriteCloser, err error) {
	return os.Create(filePath)
}
