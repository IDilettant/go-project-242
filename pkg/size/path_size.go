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

	if info.Mode().IsRegular() {
		return info.Size(), nil
	}

	if !info.IsDir() {
		return 0, fmt.Errorf("%w: %s", ErrUnsupportedFileType, path)
	}

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

		if childInfo.IsDir() {
			if !opts.Recursive {
				continue
			}

			dirSize, err := GetSize(childPath, opts)
			if err != nil {
				return 0, err
			}

			total += dirSize
		}

		if childInfo.Mode().IsRegular() {
			total += childInfo.Size()
		}
	}

	return total, nil
}

func isHidden(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}
