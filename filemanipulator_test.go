package filemanipulator

import (
	"errors"
	"io"
	"testing"
)

type mockReadWriteCloser struct {
	io.Reader
	io.Writer
	io.Closer
	closed bool
}

func (m *mockReadWriteCloser) Close() error {
	m.closed = true
	return nil
}

type MockFileHandler struct {
	doesFileExist      bool
	openreaderwriter   io.ReadWriteCloser
	openerror          error
	createreaderwriter io.ReadWriteCloser
	createerror        error
	removeerror        error
}

func (m *MockFileHandler) DoesFileExist(filePath string) bool {
	return m.doesFileExist
}

func (m *MockFileHandler) Open(filePath string) (readerWriter io.ReadWriteCloser, err error) {
	return m.openreaderwriter, m.openerror
}

func (m *MockFileHandler) Remove(filePath string) (err error) {
	return m.removeerror
}
func (m *MockFileHandler) Create(filePath string) (readerWriter io.ReadWriteCloser, err error) {
	return m.createreaderwriter, m.createerror
}

func TestFileManipulator_MoveFile(t *testing.T) {
	type args struct {
		sourcePath string
		targetPath string
	}
	tests := []struct {
		name    string
		m       *FileManipulator
		args    args
		wantErr bool
	}{
		{
			name: "non existing file",
			m: &FileManipulator{&MockFileHandler{
				openerror:          errors.New(""),
				openreaderwriter:   &mockReadWriteCloser{},
				createreaderwriter: &mockReadWriteCloser{},
			}},
			args: args{
				sourcePath: "",
				targetPath: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.MoveFile(tt.args.sourcePath, tt.args.targetPath); (err != nil) != tt.wantErr {
				t.Errorf("FileManipulator.MoveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
