package validator

import (
	"fmt"
	"strings"
	"unicode"
)

// ParseString parses the input string and returns it unchanged. It always
// succeeds and returns no error.
func ParseString(input string) (string, error) {
	return input, nil
}

const (
	messageInvalidCharIndex = "character %c at index %d is not %s"
)

// AlphaNumeric validates that the string contains only alphanumeric characters
// (0-9, A-Z, a-z). Returns a ValidationError if any non-alphanumeric character
// is found.
func AlphaNumeric(value string) error {
	for i, r := range value {
		if r < '0' ||
			(r > '9' && r < 'A') ||
			(r > 'Z' && r < 'a') ||
			r > 'z' {
			return &ValidationError[string]{
				Code:    InvalidCharacter,
				Message: fmt.Sprintf(messageInvalidCharIndex, r, i, "alphanumeric"),
				Value:   value,
			}
		}
	}
	return nil
}

// HasPrefix validates that the string starts with the specified prefix. Returns
// a ValidationError if the string does not have the prefix.
func HasPrefix(prefix string) Validator[string] {
	return func(value string) error {
		if !strings.HasPrefix(value, prefix) {
			return &ValidationError[string]{
				Code:    InvalidPrefix,
				Message: fmt.Sprintf("string %q does not have prefix %q", value, prefix),
				Value:   value,
			}
		}
		return nil
	}
}

// HasSuffix validates that the string ends with the specified suffix. Returns
// a ValidationError if the string does not have the suffix.
func HasSuffix(suffix string) Validator[string] {
	return func(value string) error {
		if !strings.HasSuffix(value, suffix) {
			return &ValidationError[string]{
				Code:    InvalidSuffix,
				Message: fmt.Sprintf("string %q does not have suffix %q", value, suffix),
				Value:   value,
			}
		}
		return nil
	}
}

// Latin validates that the string contains only alphabetic characters from the
// Latin script. Precomposed letters with diacritics, such as é, are accepted;
// separate combining marks are rejected. It returns a ValidationError for the
// first rune that is not a Latin letter.
func Latin(value string) error {
	for i, r := range value {
		if !unicode.IsLetter(r) || !unicode.Is(unicode.Latin, r) {
			return &ValidationError[string]{
				Code:    InvalidCharacter,
				Message: fmt.Sprintf(messageInvalidCharIndex, r, i, "latin"),
				Value:   value,
			}
		}
	}
	return nil
}

// Numeric validates that the string contains only numeric characters (0-9).
// Returns a ValidationError if any non-numeric character is found.
func Numeric(value string) error {
	for i, r := range value {
		if r < '0' || r > '9' {
			return &ValidationError[string]{
				Code:    InvalidCharacter,
				Message: fmt.Sprintf(messageInvalidCharIndex, r, i, "numeric"),
				Value:   value,
			}
		}
	}
	return nil
}
