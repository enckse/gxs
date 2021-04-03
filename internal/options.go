package internal

import (
	"fmt"
	"strings"
)

type (
	Option struct {
		asciiNoDelimiter bool
	}
)

func toBool(s string) (bool, error) {
	if s == "true" {
		return true, nil
	} else {
		if s == "false" {
			return false, nil
		}
	}
	return false, fmt.Errorf("invalid boolean value")
}

func (o *Option) Set(value string) error {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key=value pair")
	}
	switch parts[0] {
	case "ascii-no-delimiter":
		b, err := toBool(parts[1])
		if err != nil {
			return err
		}
		o.asciiNoDelimiter = b
	default:
		return fmt.Errorf("unknown option")
	}
	return nil
}
