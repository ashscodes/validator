package validator

// Compose combines a parser with validators, returning a parser that validates
// the parsed value in order. It stops at the first parsing or validation error
// and returns that error alongside the value produced by the parser. A value
// returned with an error is not guaranteed to have passed validation.
func Compose[T any](parser Parser[T], validators ...Validator[T]) Parser[T] {
	return func(input string) (T, error) {
		value, err := parser(input)
		if err != nil {
			return value, err
		}

		for _, validator := range validators {
			err := validator(value)
			if err != nil {
				return value, err
			}
		}

		return value, nil
	}
}
