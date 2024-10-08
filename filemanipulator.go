package filemanipulator

import (
	"fmt"
	"io"
	"os"
)

type FileProtocolHandler interface {
	DoesFileExist(filePath string) bool
	Open(filePath string) (readerWriter io.ReadWriteCloser, err error)
	Create(filePath string) (readerWriter io.ReadWriteCloser, err error)
	Remove(filePath string) (err error)
}

type FileManipulator struct {
	handler FileProtocolHandler
}

func NewFileManipulator(handler FileProtocolHandler) *FileManipulator {
	return &FileManipulator{
		handler: handler,
	}
}

func (m *FileManipulator) MoveFile(sourcePath, targetPath string) (err error) {
	inputFile, err := m.handler.Open(sourcePath)
	if err != nil {
		return err
	}

	if m.handler.DoesFileExist(targetPath) {
		// check if this is the same file
		filesEqual, err := AreFileEqual(sourcePath, targetPath)
		if err != nil {
			return err
		}
		if filesEqual {
			// same file already exists and can be removed from source
			err = inputFile.Close()
			if err != nil {
				return err
			}
			err = m.handler.Remove(sourcePath)
			if err != nil {
				return err
			}
			// stop process
			return nil
		} else {
			// remove destination file and continue
			err = m.handler.Remove(targetPath)
			if err != nil {
				return err
			}
		}
	}

	// create an alias for filename while copying is ongoing
	tempFileName := fmt.Sprintf("%s.part", targetPath)
	outputFile, err := m.handler.Create(tempFileName)
	if err != nil {
		inputFile.Close()
		return err
	}
	defer func() {
		fileClosingError := outputFile.Close()
		if fileClosingError != nil {
			return
		}

		// check if copying was successful
		if err != nil {
			return
		}

		// rename to intended name
		err = os.Rename(tempFileName, targetPath)
		if err != nil {
			return
		}

		// The copy was successful, now delete the original file
		err = m.handler.Remove(sourcePath)
	}()

	// actual file copy
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return err
	}

	return err
}
