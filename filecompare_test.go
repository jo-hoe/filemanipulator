package filemanipulator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAreFileEqual(t *testing.T) {
	_, leftDirectory, rightDirectory, fileName := setupTestMoveEnvironment(t)

	type args struct {
		leftFilePath  string
		rightFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "right file does not exist",
			args: args{
				leftFilePath:  filepath.Join(leftDirectory, fileName),
				rightFilePath: "",
			},
			want:    false,
			wantErr: true,
		}, {
			name: "left file does not exist",
			args: args{
				leftFilePath:  "",
				rightFilePath: filepath.Join(rightDirectory, fileName),
			},
			want:    false,
			wantErr: true,
		}, {
			name: "positive test",
			args: args{
				leftFilePath:  filepath.Join(leftDirectory, fileName),
				rightFilePath: filepath.Join(leftDirectory, fileName),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AreFileEqual(tt.args.leftFilePath, tt.args.rightFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreFileEqual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AreFileEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAreFileEqual_Content exercises the byte-by-byte comparison over content
// that spans multiple read chunks and differs only in the final byte, which a
// short-circuiting comparison must still detect.
func TestAreFileEqual_Content(t *testing.T) {
	dir := t.TempDir()

	// data large enough to span more than one internal comparison chunk
	identical := make([]byte, compareBufferSize*2+123)
	for i := range identical {
		identical[i] = byte(i % 251)
	}
	differsLastByte := make([]byte, len(identical))
	copy(differsLastByte, identical)
	differsLastByte[len(differsLastByte)-1] ^= 0xFF

	tests := []struct {
		name string
		left []byte
		righ []byte
		want bool
	}{
		{name: "both empty", left: []byte{}, righ: []byte{}, want: true},
		{name: "empty vs non-empty", left: []byte{}, righ: []byte("x"), want: false},
		{name: "same length differ", left: []byte("hello"), righ: []byte("hellO"), want: false},
		{name: "different length", left: []byte("hell"), righ: []byte("hello"), want: false},
		{name: "multi-chunk identical", left: identical, righ: identical, want: true},
		{name: "multi-chunk differ in last byte", left: identical, righ: differsLastByte, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leftPath := filepath.Join(dir, "left_"+tt.name)
			rightPath := filepath.Join(dir, "right_"+tt.name)
			if err := os.WriteFile(leftPath, tt.left, 0o600); err != nil {
				t.Fatalf("could not write left file: %v", err)
			}
			if err := os.WriteFile(rightPath, tt.righ, 0o600); err != nil {
				t.Fatalf("could not write right file: %v", err)
			}

			got, err := AreFileEqual(leftPath, rightPath)
			if err != nil {
				t.Fatalf("AreFileEqual() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("AreFileEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
