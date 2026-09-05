package validator

import (
	"errors"
	"strings"
)

var (
	// ErrInputEmpty indicates that the provided input value is empty.
	ErrInputEmpty = errors.New("input is empty")
)

// ParseError reports a failure to parse an input value as a requested type.
type ParseError struct {
	Cause error
	Input string
	Type  string
}

// Error returns a description of the parsing failure.
func (err *ParseError) Error() string {
	return "failed to parse input " + err.Input + " as " + err.Type
}

// Unwrap returns the underlying parsing error.
func (err *ParseError) Unwrap() error {
	return err.Cause
}

// ValidationCode identifies the reason a validation error occurred.
type ValidationCode string

const (
	AboveMax         ValidationCode = "above_maximum"
	AfterEnd         ValidationCode = "after_end"
	BeforeStart      ValidationCode = "before_start"
	BelowMin         ValidationCode = "below_minimum"
	FalseValue       ValidationCode = "false_value"
	InvalidCharacter ValidationCode = "invalid_character"
	InvalidPrefix    ValidationCode = "invalid_prefix"
	InvalidRange     ValidationCode = "invalid_range"
	InvalidSuffix    ValidationCode = "invalid_suffix"
	ItemRequired     ValidationCode = "required"
	LengthTooLong    ValidationCode = "length_too_long"
	LengthTooShort   ValidationCode = "length_too_short"
	OutOfRange       ValidationCode = "out_of_range"
	TrueValue        ValidationCode = "true_value"
	WithinRange      ValidationCode = "within_range"
)

// ValidationError describes a validation failure for a value.
type ValidationError[T any] struct {
	Code    ValidationCode
	Message string
	Value   T
}

// Error returns the validation error message.
func (err *ValidationError[T]) Error() string {
	return err.Message
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []error

// Error returns all non-nil error messages separated by newlines.
func (errs ValidationErrors) Error() string {
	var builder strings.Builder
	for _, err := range errs {
		if err == nil {
			continue
		}

		if builder.Len() > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(err.Error())
	}
	return builder.String()
}

// Unwrap returns a copy of the contained errors.
func (errs ValidationErrors) Unwrap() []error {
	out := make([]error, len(errs))
	copy(out, errs)
	return out
}
