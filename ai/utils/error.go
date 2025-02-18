package utils

import (
	"fmt"
)

const (
	colorRed   = "\033[31m"
	colorBlue  = "\033[34m"
	colorReset = "\033[0m"
)

// NewErrorWithSource returns a colorful error message
func NewErrorWithSource(err error, source string) error {
	return fmt.Errorf("\n%sSource:%s %s\n%sError:%s %v",
		colorBlue, colorReset, source,
		colorRed, colorReset, err,
	)
}
