package code

import (
	"fmt"
)

const unitBase = 1024

var units = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// FormatSize formats a size in bytes.
// If human is false, it returns the size in bytes with "B" suffix.
// If human is true, it returns a human-readable representation using base 1024.
func FormatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", size)
	}

	value := float64(size)
	unit := 0

	for value >= unitBase && unit < len(units)-1 {
		value /= unitBase
		unit++
	}

	if unit == 0 {
		return fmt.Sprintf("%d%s", int64(value), units[unit])
	}

	return fmt.Sprintf("%.1f%s", value, units[unit])
}

// FormatOutput formats size and path as a single output line.
func FormatOutput(sizeStr, path string) string {
	return fmt.Sprintf("%s\t%s", sizeStr, path)
}
