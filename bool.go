package validator

import (
	"fmt"
	"strings"
)

// ParseBool parses a string as a boolean value.
//
// It accepts true, 1, t, yes, and y as true, and false, 0, f, no, and n as
// false, without regard to case.
func ParseBool(input string) (bool, error) {
	switch strings.ToLower(input) {
	case "true", "1", "t", "yes", "y":
		return true, nil
	case "false", "0", "f", "no", "n":
		return false, nil
	default:
		return false, &ParseError{
			Cause: fmt.Errorf("invalid boolean value: %s", input),
			Input: input,
			Type:  "bool",
		}
	}
}

// False reports an error when value is true.
func False(value bool) error {
	if value {
		return &ValidationError[bool]{
			Code:    FalseValue,
			Message: "value must be false",
			Value:   value,
		}
	}
	return nil
}

// True reports an error when value is false.
func True(value bool) error {
	if !value {
		return &ValidationError[bool]{
			Code:    TrueValue,
			Message: "value must be true",
			Value:   value,
		}
	}
	return nil
}
