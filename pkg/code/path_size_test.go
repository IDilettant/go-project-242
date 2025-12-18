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
		opts      Options
		want      int64
		wantErr   bool
		errAssert func(error) bool
	}{
		{
			name: "file/regular",
			path: filepath.Join(td, "file_5b.txt"),
			opts: Options{},
			want: 5,
		},
		{
			name: "file/empty",
			path: filepath.Join(td, "file_0b.txt"),
			opts: Options{},
			want: 0,
		},
		{
			name: "dir/empty",
			path: filepath.Join(td, "emptydir"),
			opts: Options{},
			want: 0,
		},

		{
			name: "dir/dirA non-recursive ignores hidden and nested",
			path: filepath.Join(td, "dirA"),
			opts: Options{All: false, Recursive: false},
			want: 10,
		},
		{
			name: "dir/dirA non-recursive includes hidden files",
			path: filepath.Join(td, "dirA"),
			opts: Options{All: true, Recursive: false},
			want: 19,
		},
		{
			name: "dir/dirA recursive ignores hidden",
			path: filepath.Join(td, "dirA"),
			opts: Options{All: false, Recursive: true},
			want: 14,
		},
		{
			name: "dir/dirA recursive includes hidden files and directories",
			path: filepath.Join(td, "dirA"),
			opts: Options{All: true, Recursive: true},
			want: 28,
		},

		{
			name: "dir/dir_with_dir_only non-recursive",
			path: filepath.Join(td, "dir_with_dir_only"),
			opts: Options{Recursive: false},
			want: 0,
		},
		{
			name: "dir/dir_with_dir_only recursive",
			path: filepath.Join(td, "dir_with_dir_only"),
			opts: Options{Recursive: true},
			want: 2,
		},

		{
			name: "hidden/file ignored when all=false",
			path: filepath.Join(td, "dirA", ".hidden_9b.txt"),
			opts: Options{All: false},
			want: 0,
		},
		{
			name: "hidden/file included when all=true",
			path: filepath.Join(td, "dirA", ".hidden_9b.txt"),
			opts: Options{All: true},
			want: 9,
		},

		{
			name:    "error/path does not exist",
			path:    filepath.Join(td, "no_such_path"),
			opts:    Options{},
			wantErr: true,
			errAssert: func(err error) bool {
				return os.IsNotExist(err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSize(tt.path, tt.opts)

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
