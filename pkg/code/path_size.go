package code

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetSize returns size of a file or, for a directory, sums sizes of files in the first level
func GetSize(path string) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if info.Mode().IsRegular() {
		return info.Size(), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, e := range entries {
		childPath := filepath.Join(path, e.Name())

		childInfo, err := os.Lstat(childPath)
		if err != nil {
			return 0, err
		}

		if childInfo.IsDir() {
			continue
		}

		total += childInfo.Size()
	}

	return total, nil
}

func FormatSizeOutput(size int64, path string) string {
	return fmt.Sprintf("%d\t%s", size, path)
}

