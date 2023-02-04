package filemanipulator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFile(t *testing.T) {
	_, leftDirectory, rightDirectory, fileName := setupTestMoveEnvironment(t)
	target := filepath.Join(rightDirectory, fileName)
	origin := filepath.Join(leftDirectory, fileName)

	manipluator := NewFileManipulator(&MountedFileHandler{})
	err := manipluator.MoveFile(origin, target)

	if err != nil {
		t.Errorf("found error '%+v'", err)
	}
	if _, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		t.Errorf("file '%s' was not found", fileName)
	}
	if _, err := os.Stat(origin); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file '%s' was removed", fileName)
	}
}

func TestMoveFileOnToExistingFile(t *testing.T) {
	_, leftDirectory, rightDirectory, fileName := setupTestMoveEnvironment(t)

	target := filepath.Join(rightDirectory, fileName)
	origin := filepath.Join(leftDirectory, fileName)

	file, err := os.Create(target)
	if err != nil {
		t.Error("could not create file")
	}
	file.Close()

	manipluator := NewFileManipulator(&MountedFileHandler{})
	err = manipluator.MoveFile(origin, target)

	if err != nil {
		t.Errorf("found error '%+v'", err)
	}
	if _, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		t.Errorf("file '%s' was not found", fileName)
	}
	if _, err := os.Stat(origin); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file '%s' was removed", fileName)
	}
}

func TestMoveFileOnToCorruptExistingFile(t *testing.T) {
	_, leftDirectory, rightDirectory, fileName := setupTestMoveEnvironment(t)

	target := filepath.Join(rightDirectory, fileName)
	origin := filepath.Join(leftDirectory, fileName)

	file, err := os.Create(target)
	if err != nil {
		t.Error("could not create file")
	}
	_, err = file.WriteString("corrupt")
	if err != nil {
		t.Error("could not write to file")
	}
	file.Close()

	manipluator := NewFileManipulator(&MountedFileHandler{})
	err = manipluator.MoveFile(origin, target)

	if err != nil {
		t.Errorf("found error '%+v'", err)
	}
	if _, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		t.Errorf("file '%s' was not found", fileName)
	}
	if _, err := os.Stat(origin); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("file '%s' was removed", fileName)
	}
}
