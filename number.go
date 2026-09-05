package validator

import (
	"fmt"
	"strconv"
)

// Number represents the numeric types supported by the numeric validators.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// ParseFloat32 parses input as a 32-bit floating-point number.
func ParseFloat32(input string) (float32, error) {
	value, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "float32",
		}
	}
	return float32(value), nil
}

// ParseFloat64 parses input as a 64-bit floating-point number.
func ParseFloat64(input string) (float64, error) {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "float64",
		}
	}
	return value, nil
}

// ParseInt parses input as an int.
func ParseInt(input string) (int, error) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "int",
		}
	}
	return value, nil
}

// ParseInt8 parses input as an 8-bit signed integer.
func ParseInt8(input string) (int8, error) {
	value, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "int8",
		}
	}
	return int8(value), nil
}

// ParseInt16 parses input as a 16-bit signed integer.
func ParseInt16(input string) (int16, error) {
	value, err := strconv.ParseInt(input, 10, 16)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "int16",
		}
	}
	return int16(value), nil
}

// ParseInt32 parses input as a 32-bit signed integer.
func ParseInt32(input string) (int32, error) {
	value, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "int32",
		}
	}
	return int32(value), nil
}

// ParseInt64 parses input as a 64-bit signed integer.
func ParseInt64(input string) (int64, error) {
	value, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "int64",
		}
	}
	return value, nil
}

// ParseUint parses input as a uint.
func ParseUint(input string) (uint, error) {
	value, err := strconv.ParseUint(input, 10, 0)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "uint",
		}
	}
	return uint(value), nil
}

// ParseUint8 parses input as an 8-bit unsigned integer.
func ParseUint8(input string) (uint8, error) {
	value, err := strconv.ParseUint(input, 10, 8)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "uint8",
		}
	}
	return uint8(value), nil
}

// ParseUint16 parses input as a 16-bit unsigned integer.
func ParseUint16(input string) (uint16, error) {
	value, err := strconv.ParseUint(input, 10, 16)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "uint16",
		}
	}
	return uint16(value), nil
}

// ParseUint32 parses input as a 32-bit unsigned integer.
func ParseUint32(input string) (uint32, error) {
	value, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "uint32",
		}
	}
	return uint32(value), nil
}

// ParseUint64 parses input as a 64-bit unsigned integer.
func ParseUint64(input string) (uint64, error) {
	value, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "uint64",
		}
	}
	return value, nil
}

const (
	valueLimitMessage = "value must %s than or equal to %v"
	valueRangeMessage = "value must %s between %v and %v"
)

// InRange returns a validator that accepts values between minimum and maximum,
// inclusive. NaN values are rejected.
func InRange[T Number](minimum, maximum T) Validator[T] {
	return func(value T) error {
		if !(value >= minimum && value <= maximum) {
			return &ValidationError[T]{
				Code:    OutOfRange,
				Message: fmt.Sprintf(valueRangeMessage, "be", minimum, maximum),
				Value:   value,
			}
		}
		return nil
	}
}

// Max returns a validator that accepts values less than or equal to maximum.
// NaN values are rejected.
func Max[T Number](maximum T) Validator[T] {
	return func(value T) error {
		if !(value <= maximum) {
			return &ValidationError[T]{
				Code:    AboveMax,
				Message: fmt.Sprintf(valueLimitMessage, "be less", maximum),
				Value:   value,
			}
		}
		return nil
	}
}

// Min returns a validator that accepts values greater than or equal to minimum.
// NaN values are rejected.
func Min[T Number](minimum T) Validator[T] {
	return func(value T) error {
		if !(value >= minimum) {
			return &ValidationError[T]{
				Code:    BelowMin,
				Message: fmt.Sprintf(valueLimitMessage, "be greater", minimum),
				Value:   value,
			}
		}
		return nil
	}
}

// NotInRange returns a validator that accepts values strictly below minimum
// or strictly above maximum. NaN values are rejected.
func NotInRange[T Number](minimum, maximum T) Validator[T] {
	return func(value T) error {
		if !(value < minimum || value > maximum) {
			return &ValidationError[T]{
				Code:    WithinRange,
				Message: fmt.Sprintf(valueRangeMessage, "not be", minimum, maximum),
				Value:   value,
			}
		}
		return nil
	}
}
