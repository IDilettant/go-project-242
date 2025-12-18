package code

import (
	"code/pkg/size"
)

// GetPathSize returns the formatted size of a file or directory with optional flags.
func GetPathSize(path string, recursive, human, all bool) (string, error) {
	opts := size.Options{
		All:       all,
		Recursive: recursive,
	}

	pathSize, err := size.GetSize(path, opts)
	if err != nil {
		return "", err
	}

	formatedSize := size.FormatSize(pathSize, human)

	return formatedSize, nil
}
