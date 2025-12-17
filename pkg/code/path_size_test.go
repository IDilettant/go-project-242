package code

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetSize(t *testing.T) {
	td := filepath.Join("testdata")

	tests := []struct {
		name      string
		path      string
		want      int64
		wantErr   bool
		errAssert func(error) bool
	}{
		{
			name: "regular file: 5 bytes",
			path: filepath.Join(td, "file_5b.txt"),
			want: 5,
		},
		{
			name: "empty regular file: 0 bytes",
			path: filepath.Join(td, "file_0b.txt"),
			want: 0,
		},
		{
			name: "empty dir: 0 bytes",
			path: filepath.Join(td, "emptydir"),
			want: 0,
		},
		{
			name: "dir sums only first-level files, ignores nested dir",
			path: filepath.Join(td, "dirA"),
			want: 10,
		},
		{
			name: "dir with only subdir on first level => 0",
			path: filepath.Join(td, "dir_with_dir_only"),
			want: 0,
		},
		{
			name:    "non-existent path => error",
			path:    filepath.Join(td, "no_such_path"),
			wantErr: true,
			errAssert: func(err error) bool {
				return os.IsNotExist(err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSize(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (size=%d)", got)
				}

				if tt.errAssert != nil && !tt.errAssert(err) {
					t.Fatalf("unexpected error: %v", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("size mismatch: got %d, want %d", got, tt.want)
			}
		})
	}
}
