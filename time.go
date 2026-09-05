package validator

import (
	"fmt"
	"strings"
	"time"
)

// TimeRange represents a time interval with a start and end time.
type TimeRange struct {
	End   time.Time
	Start time.Time
}

// Contains returns true if the value time is within the TimeRange (inclusive).
func (tr TimeRange) Contains(value time.Time) bool {
	return !tr.Start.IsZero() &&
		!tr.End.IsZero() &&
		!tr.End.Before(tr.Start) &&
		!value.Before(tr.Start) &&
		!value.After(tr.End)
}

// Duration returns the duration between the start and end time of the
// TimeRange. It returns 0 if either start or end is zero, or if end is before
// start.
func (tr TimeRange) Duration() time.Duration {
	if tr.Start.IsZero() ||
		tr.End.IsZero() ||
		tr.End.Before(tr.Start) {
		return 0
	}
	return tr.End.Sub(tr.Start)
}

// Overlaps returns true if the ranges share any instant, including a touching
// endpoint. It returns false if either range has a zero endpoint or an end
// before its start.
func (tr TimeRange) Overlaps(other TimeRange) bool {
	return !tr.Start.IsZero() &&
		!tr.End.IsZero() &&
		!tr.End.Before(tr.Start) &&
		!other.Start.IsZero() &&
		!other.End.IsZero() &&
		!other.End.Before(other.Start) &&
		!tr.Start.After(other.End) &&
		!tr.End.Before(other.Start)
}

// ParseDuration parses a duration string and returns a time.Duration. It
// returns a ParseError if the input cannot be parsed as a duration.
func ParseDuration(input string) (time.Duration, error) {
	value, err := time.ParseDuration(input)
	if err != nil {
		return 0, &ParseError{
			Cause: err,
			Input: input,
			Type:  "duration",
		}
	}
	return value, nil
}

// parseTimeWithFormat parses a time string using the specified format. It
// returns a ParseError if the input cannot be parsed with the given format.
func parseTimeWithFormat(input, format string) (time.Time, error) {
	value, err := time.Parse(format, input)
	if err != nil {
		return time.Time{}, &ParseError{
			Cause: err,
			Input: input,
			Type:  "time",
		}
	}
	return value, nil
}

// splitTimeRange splits a time range string into start and end components. It
// supports '|', ';', and '~' as separators. It returns a ParseError if the
// input is empty or does not contain a valid separator.
func splitTimeRange(input string) (string, string, error) {
	if input == "" {
		return "", "", &ParseError{
			Cause: ErrInputEmpty,
			Input: input,
			Type:  "time range",
		}
	}

	for _, splitChar := range []rune{'|', ';', '~'} {
		parts := strings.Split(input, string(splitChar))
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	}

	return "", "", &ParseError{
		Cause: fmt.Errorf("invalid separator, must be '|', ';', or '~'"),
		Input: input,
		Type:  "time range",
	}
}

// ParseTime parses a time string in RFC3339 format. It returns a ParseError if
// the input cannot be parsed as a time.
func ParseTime(input string) (time.Time, error) {
	return parseTimeWithFormat(input, time.RFC3339)
}

// ParseTimeRange parses a time range string into a TimeRange. The input string
// must contain two RFC3339-formatted times separated by '|', ';', or '~'.
// It returns a ParseError if the input cannot be split or either time cannot
// be parsed. It does not validate endpoint order or reject zero times;
// use ValidTimeRange to validate the parsed range.
func ParseTimeRange(input string) (TimeRange, error) {
	startStr, endStr, err := splitTimeRange(input)
	if err != nil {
		return TimeRange{}, err
	}

	start, err := ParseTime(startStr)
	if err != nil {
		return TimeRange{}, err
	}

	end, err := ParseTime(endStr)
	if err != nil {
		return TimeRange{}, err
	}

	return TimeRange{
		Start: start,
		End:   end,
	}, nil
}

// ParseTimeRangeWithFormats returns a parser function that parses time ranges
// separated by '|', ';', or '~' using the specified formats. It tries each
// format in order and accepts the first that parses both endpoints using the
// same layout. If no formats are provided, it defaults to ParseTimeRange with
// RFC3339 format. Use ValidTimeRange to validate the parsed range.
func ParseTimeRangeWithFormats(formats ...string) Parser[TimeRange] {
	if len(formats) == 0 {
		return func(input string) (TimeRange, error) {
			return ParseTimeRange(input)
		}
	}

	return func(input string) (TimeRange, error) {
		startStr, endStr, err := splitTimeRange(input)
		if err != nil {
			return TimeRange{}, err
		}

		var lastErr error
		for _, format := range formats {
			start, err := parseTimeWithFormat(startStr, format)
			if err != nil {
				lastErr = err
				continue
			}

			end, err := parseTimeWithFormat(endStr, format)
			if err != nil {
				lastErr = err
				continue
			}

			return TimeRange{
				Start: start,
				End:   end,
			}, nil
		}
		return TimeRange{}, lastErr
	}
}

// ParseTimeWithFormats returns a parser function that parses times using the
// specified formats. If no formats are provided, it defaults to using ParseTime
// with RFC3339 format. It tries each format in order and returns the first
// successful parse.
func ParseTimeWithFormats(formats ...string) Parser[time.Time] {
	if len(formats) == 0 {
		return func(input string) (time.Time, error) {
			return ParseTime(input)
		}
	}

	return func(input string) (time.Time, error) {
		var lastErr error
		for _, format := range formats {
			value, err := parseTimeWithFormat(input, format)
			if err != nil {
				lastErr = err
				continue
			}
			return value, nil
		}

		return time.Time{}, &ParseError{
			Cause: lastErr,
			Input: input,
			Type:  "time",
		}
	}
}

// NewTimeRange creates and returns a new TimeRange with the specified start and
// end times.
func NewTimeRange(start, end time.Time) TimeRange {
	return TimeRange{
		Start: start,
		End:   end,
	}
}

// ValidTimeRange validates a TimeRange and returns an error if it is invalid. A
// TimeRange is invalid if the start or end time is zero, or if the end time is
// before the start time.
func ValidTimeRange(tr TimeRange) error {
	if tr.Start.IsZero() {
		return &ValidationError[TimeRange]{
			Code:    InvalidRange,
			Message: "start time is zero",
			Value:   tr,
		}
	}

	if tr.End.IsZero() {
		return &ValidationError[TimeRange]{
			Code:    InvalidRange,
			Message: "end time is zero",
			Value:   tr,
		}
	}

	if tr.End.Before(tr.Start) {
		return &ValidationError[TimeRange]{
			Code:    InvalidRange,
			Message: "end time is before start time",
			Value:   tr,
		}
	}
	return nil
}
