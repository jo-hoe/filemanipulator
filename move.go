package filemanipulator

import (
	"errors"
	"io"
	"os"
)

func doesFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func MoveFile(sourcePath, targetPath string) (err error) {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	if doesFileExist(targetPath) {
		// check if this is the same file
		filesEqual, err := areFileEqual(sourcePath, targetPath)
		if err != nil {
			return err
		}
		if filesEqual {
			// same file already exists and can be removed from source
			err = inputFile.Close()
			if err != nil {
				return err
			}
			err = os.Remove(sourcePath)
			if err != nil {
				return err
			}
			// stop process
			return nil
		} else {
			// remove destination file and continue
			err = os.Remove(targetPath)
			if err != nil {
				return err
			}
		}
	}

	outputFile, err := os.Create(targetPath)
	if err != nil {
		inputFile.Close()
		return err
	}
	defer func() {
		fileClosingError := outputFile.Close()
		if fileClosingError != nil {
			return
		}

		// check if copying was successfull
		if err != nil {
			return
		}

		// The copy was successful, so now delete the original file
		err = os.Remove(sourcePath)
	}()

	// actual file copy
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return err
	}

	return nil
}
