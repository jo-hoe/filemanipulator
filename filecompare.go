package filemanipulator

import (
	"bytes"
	"io"
	"os"
)

// compareBufferSize is the chunk size used when comparing files byte-by-byte.
const compareBufferSize = 64 * 1024

// AreFileEqual reports whether the two files have identical contents.
// It compares the files chunk-by-chunk and returns as soon as a difference
// is found, so differing files are detected without reading either file in
// full.
func AreFileEqual(leftFilePath string, rightFilePath string) (equal bool, err error) {
	leftFile, err := os.Open(leftFilePath)
	if err != nil {
		return false, err
	}
	defer func() {
		if closeErr := leftFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	rightFile, err := os.Open(rightFilePath)
	if err != nil {
		return false, err
	}
	defer func() {
		if closeErr := rightFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	leftBuffer := make([]byte, compareBufferSize)
	rightBuffer := make([]byte, compareBufferSize)
	for {
		leftCount, leftErr := io.ReadFull(leftFile, leftBuffer)
		rightCount, rightErr := io.ReadFull(rightFile, rightBuffer)

		if !bytes.Equal(leftBuffer[:leftCount], rightBuffer[:rightCount]) {
			return false, nil
		}

		leftDone := leftErr == io.EOF || leftErr == io.ErrUnexpectedEOF
		rightDone := rightErr == io.EOF || rightErr == io.ErrUnexpectedEOF
		if leftDone || rightDone {
			// both files ended together -> equal; otherwise different lengths
			return leftDone && rightDone, nil
		}
		if leftErr != nil {
			return false, leftErr
		}
		if rightErr != nil {
			return false, rightErr
		}
	}
}
