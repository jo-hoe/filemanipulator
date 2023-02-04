package filemanipulator

import (
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
