package size

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	filePerm = 0o644
	dirPerm  = 0o755
)

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()

	require.NoError(t, os.WriteFile(path, data, filePerm))
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()

	require.NoError(t, os.MkdirAll(path, dirPerm))
}

func makeSymlink(t *testing.T, target, link string) {
	t.Helper()

	require.NoError(t, os.Symlink(target, link))
}

func TestGetSize(t *testing.T) {
	td := filepath.Join("testdata")

	type testCase struct {
		name      string
		path      string
		opts      Options
		want      int64
		errAssert func(error) bool
	}

	tests := []testCase{
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
			name: "dir/dirA recursive excludes hidden files",
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
			name: "error/path does not exist",
			path: filepath.Join(td, "no_such_path"),
			opts: Options{},
			errAssert: func(err error) bool {
				return os.IsNotExist(err)
			},
		},
		{
			name: "file/unicode name",
			path: filepath.Join(td, "файл_6b.txt"),
			opts: Options{},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.path

			got, err := GetSize(path, tt.opts)

			if tt.errAssert != nil {
				require.Error(t, err)
				require.Truef(t, tt.errAssert(err), "unexpected error: %v", err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetSize_Symlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink tests are skipped on windows")
	}

	type testCase struct {
		name      string
		opts      Options
		want      int64
		errAssert func(error) bool
		setup     func(t *testing.T) string
	}

	tests := []testCase{
		{
			name: "in-dir ignored when summing",
			opts: Options{All: true, Recursive: false},
			want: 5,
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()

				target := filepath.Join(dir, "file_5b.txt")
				writeFile(t, target, []byte("12345"))

				link := filepath.Join(dir, "link_to_file")
				makeSymlink(t, target, link)

				return dir
			},
		},
		{
			name: "to-dir ignored even in recursive mode",
			opts: Options{All: true, Recursive: true},
			want: 4,
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()

				realSub := filepath.Join(dir, "subdir")
				mkdirAll(t, realSub)

				writeFile(t, filepath.Join(realSub, "a_4b.txt"), []byte("1234"))

				linkSub := filepath.Join(dir, "link_to_subdir")
				makeSymlink(t, realSub, linkSub)

				return dir
			},
		},
		{
			name: "root path unsupported",
			opts: Options{},
			errAssert: func(err error) bool {
				return err != nil && errors.Is(err, ErrUnsupportedFileType)
			},
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()

				target := filepath.Join(dir, "file_3b.txt")
				writeFile(t, target, []byte("123"))

				link := filepath.Join(dir, "link_root")
				makeSymlink(t, target, link)

				return link
			},
		},
	}

	for _, tt := range tests {
		t.Run("symlink/"+tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.setup(t)
			got, err := GetSize(path, tt.opts)

			if tt.errAssert != nil {
				require.Error(t, err)
				require.Truef(t, tt.errAssert(err), "unexpected error: %v", err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

