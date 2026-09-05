package validator

// Parser parses a string and returns a value of type T or an error.
type Parser[T any] func(string) (T, error)

// Transform transforms a value of type A into a value of type B and returns the
// transformed value or an error.
type Transform[A, B any] func(A) (B, error)

// Validator validates a value of type T and returns an error if validation
// fails.
type Validator[T any] func(T) error
