package code

import "errors"

var ErrUnsupportedFileType = errors.New("error: path is not a regular file or directory")
