package size

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	All       bool
	Recursive bool
}

// GetSize returns size of a file or, for a directory, sums sizes of files in the first level
func GetSize(path string, opts Options) (int64, error) {
	if !opts.All && isHidden(path) {
		return 0, nil
	}

	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	switch {
	case info.Mode().IsRegular():
		return getFileSize(info), nil

	case info.IsDir():
		return getDirSize(path, opts)

	default:
		return 0, fmt.Errorf("%w: %s", ErrUnsupportedFileType, path)
	}
}

func getFileSize(info os.FileInfo) int64 {
	return info.Size()
}

func getDirSize(path string, opts Options) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var total int64

	for _, e := range entries {
		name := e.Name()

		if !opts.All && isHidden(name) {
			continue
		}

		childPath := filepath.Join(path, name)

		childInfo, err := os.Lstat(childPath)
		if err != nil {
			return 0, err
		}

		switch {
		case childInfo.Mode().IsRegular():
			total += getFileSize(childInfo)

		case childInfo.IsDir() && opts.Recursive:
			dirSize, err := getDirSize(childPath, opts)
			if err != nil {
				return 0, err
			}

			total += dirSize
		}
	}

	return total, nil
}

func isHidden(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}
