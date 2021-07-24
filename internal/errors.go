package internal

import (
	"fmt"
)

type (
	// GXSError are internal gxs errors.
	GXSError struct {
		category string
		message  string
	}
)

// NewGXSError creates a new gxs error.
func NewGXSError(category, message string) error {
	return &GXSError{
		category: category,
		message: message,
	}
}

func (err *GXSError) Error() string {
	return fmt.Sprintf("%s: %v", err.category, err.message)
}
