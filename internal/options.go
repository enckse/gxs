package internal

import (
	"strings"
)

type (
	// Option are CLI argument options.
	Option struct {
		asciiNoDelimiter bool
	}
)

// NoDelimiterASCII indicates if the ascii delimiter setting is on.
func (o Option) NoDelimiterASCII() bool {
	return o.asciiNoDelimiter
}

func toBool(s string) (bool, error) {
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	return false, NewOptionsError("invalid boolean value")
}

// NewOptionsError will create a new options-based error.
func NewOptionsError(message string) error {
	return NewGXSError("options", message)
}

// Set will set a CLI key=value property.
func (o *Option) Set(value string) error {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return NewOptionsError("invalid key=value pair")
	}
	switch parts[0] {
	case "ascii-no-delimiter":
		b, err := toBool(parts[1])
		if err != nil {
			return err
		}
		o.asciiNoDelimiter = b
	default:
		return NewOptionsError("unknown option")
	}
	return nil
}
