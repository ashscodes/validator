package validator

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"
)

func TestBooleanValidators(t *testing.T) {
	err := False(false)
	if err != nil {
		t.Errorf("False(false) returned unexpected error: %v", err)
	}

	err = True(true)
	if err != nil {
		t.Errorf("True(true) returned unexpected error: %v", err)
	}

	for name, validator := range map[string]func(bool) error{
		"False": False,
		"True":  True,
	} {
		err = validator(name == "False")
		if err == nil {
			t.Errorf("%s(%v) returned nil, want validation error", name, name == "False")
		}
	}
}

func TestCompose(t *testing.T) {
	check := Compose(ParseInt, Min(1), Max(10))
	got, err := check("5")
	if err != nil || got != 5 {
		t.Errorf("Compose valid input = %d, %v; want 5, nil", got, err)
	}

	_, err = check("0")
	if err == nil {
		t.Error("Compose should reject values below the minimum")
	}

	_, err = check("invalid")
	if err == nil || !strings.Contains(err.Error(), "failed to parse input") {
		t.Errorf("Compose invalid input error = %v", err)
	}
}


func TestNumberParsers(t *testing.T) {
	tests := []struct {
		name  string
		parse func(string) (any, error)
		input string
		want  any
	}{
		{
			"ParseFloat32",
			func(input string) (any, error) { return ParseFloat32(input) },
			"1.25",
			float32(1.25),
		},
		{
			"ParseFloat64",
			func(input string) (any, error) { return ParseFloat64(input) },
			"1.25",
			float64(1.25),
		},
		{
			"ParseInt",
			func(input string) (any, error) { return ParseInt(input) },
			"-12",
			-12,
		},
		{
			"ParseInt8",
			func(input string) (any, error) { return ParseInt8(input) },
			"-12",
			int8(-12),
		},
		{
			"ParseInt16",
			func(input string) (any, error) { return ParseInt16(input) },
			"-12",
			int16(-12),
		},
		{
			"ParseInt32",
			func(input string) (any, error) { return ParseInt32(input) },
			"-12",
			int32(-12),
		},
		{
			"ParseInt64",
			func(input string) (any, error) { return ParseInt64(input) },
			"-12",
			int64(-12),
		},
		{
			"ParseUint",
			func(input string) (any, error) { return ParseUint(input) },
			"12",
			uint(12),
		},
		{
			"ParseUint8",
			func(input string) (any, error) { return ParseUint8(input) },
			"12",
			uint8(12),
		},
		{
			"ParseUint16",
			func(input string) (any, error) { return ParseUint16(input) },
			"12",
			uint16(12),
		},
		{
			"ParseUint32",
			func(input string) (any, error) { return ParseUint32(input) },
			"12",
			uint32(12),
		},
		{
			"ParseUint64",
			func(input string) (any, error) { return ParseUint64(input) },
			"12",
			uint64(12),
		},
	}

	for _, test := range tests {
		got, err := test.parse(test.input)
		if err != nil {
			t.Errorf("%s(%q) returned unexpected error: %v", test.name, test.input, err)
			continue
		}

		if got != test.want {
			t.Errorf("%s(%q) = %v, want %v", test.name, test.input, got, test.want)
		}
	}

	_, err := ParseInt("not-a-number")
	if err == nil {
		t.Error("ParseInt(invalid input) returned nil error")
	}

	_, err = ParseUint("-1")
	if err == nil {
		t.Error("ParseUint(negative input) returned nil error")
	}
}

func TestNumberValidators(t *testing.T) {
	err := InRange(1, 10)(5)
	if err != nil {
		t.Errorf("InRange accepted value returned unexpected error: %v", err)
	}

	check := NotInRange(1, 10)
	err = check(11)
	if err != nil {
		t.Errorf("NotInRange outside value returned unexpected error: %v", err)
	}

	check = Max(10)
	err = check(10)
	if err != nil {
		t.Errorf("Max boundary returned unexpected error: %v", err)
	}

	check = Min(1)
	err = check(1)
	if err != nil {
		t.Errorf("Min boundary returned unexpected error: %v", err)
	}


	tests := []struct{
		name string
		check Validator[int]
		input int
	}{
		{ "InRange", InRange(1, 10), 5 },
		{ "NotInRange", NotInRange(1, 10), 11 },
		{ "Max", Max(10), 10 },
		{ "Min", Min(1), 1 },
	}

	for _, test := range tests {
		err := test.check(test.input)
		if err != nil {
			t.Errorf("%s(%v) returned unexpected error: %v", test.name, test.input, err)
		}
	}
}

func TestNumberValidatorsRejectNaN(t *testing.T) {
	tests := []struct {
		name  string
		check Validator[float64]
	}{
		{"InRange", InRange(1.0, 10.0)},
		{"NotInRange", NotInRange(1.0, 10.0)},
		{"Max", Max(10.0)},
		{"Min", Min(1.0)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.check(math.NaN())
			if err == nil {
				t.Fatal("accepted NaN, want validation error")
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"YES", true},
		{"0", false},
		{"No", false},
	}

	for _, test := range tests {
		got, err := ParseBool(test.input)
		if err != nil {
			t.Fatalf("ParseBool(%q) returned unexpected error: %v", test.input, err)
		}

		if got != test.want {
			t.Errorf("ParseBool(%q) = %v, want %v", test.input, got, test.want)
		}
	}

	got, err := ParseBool("maybe")
	if got {
		t.Errorf("ParseBool(%q) = true, want false", "maybe")
	}

	_, ok := errors.AsType[*ParseError](err)
	if !ok {
		t.Fatalf("ParseBool(%q) error = %T, want *ParseError", "maybe", err)
	}
}

func TestStringFunctions(t *testing.T) {
	value, err := ParseString("hello")
	if err != nil || value != "hello" {
		t.Fatalf("ParseString(%q) = %q, %v; want %q, nil", "hello", value, err, "hello")
	}

	validators := []struct {
		name  string
		value string
		check Validator[string]
	}{
		{"AlphaNumeric", "abc123", AlphaNumeric},
		{"Latin", "cafe", Latin},
		{"Latin with diacritic", "café", Latin},
		{"Numeric", "012345", Numeric},
		{"HasPrefix", "prefix-value", HasPrefix("prefix-")},
		{"HasSuffix", "value-suffix", HasSuffix("-suffix")},
	}

	for _, test := range validators {
		err = test.check(test.value)
		if err != nil {
			t.Errorf("%s(%q) returned unexpected error: %v", test.name, test.value, err)
		}
	}

	invalid := []struct {
		name  string
		value string
		check Validator[string]
	}{
		{"AlphaNumeric", "abc-123", AlphaNumeric},
		{"Latin", "cafe2", Latin},
		{"Numeric", "12a", Numeric},
		{"HasPrefix", "value", HasPrefix("prefix")},
		{"HasSuffix", "value", HasSuffix("suffix")},
	}

	for _, test := range invalid {
		err := test.check(test.value)
		if err == nil {
			t.Errorf("%s(%q) returned nil, want validation error", test.name, test.value)
		}
	}
}

func TestTimeFunctions(t *testing.T) {
	start := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	rangeValue := NewTimeRange(start, end)

	if !rangeValue.Contains(start) || !rangeValue.Contains(end) {
		t.Error("TimeRange.Contains should include both boundaries")
	}

	if rangeValue.Contains(end.Add(time.Nanosecond)) {
		t.Error("TimeRange.Contains accepted a value after the end")
	}

	got := rangeValue.Duration()
	if got != 2*time.Hour {
		t.Errorf("TimeRange.Duration() = %v, want %v", got, 2*time.Hour)
	}

	if !rangeValue.Overlaps(NewTimeRange(end, end.Add(time.Hour))) {
		t.Error("TimeRange.Overlaps should include touching ranges")
	}

	if rangeValue.Overlaps(NewTimeRange(end.Add(time.Hour), end.Add(2*time.Hour))) {
		t.Error("TimeRange.Overlaps reported non-overlapping ranges")
	}

	err := ValidTimeRange(rangeValue)
	if err != nil {
		t.Errorf("ValidTimeRange returned unexpected error: %v", err)
	}

	input := start.Format(time.RFC3339) + "|" + end.Format(time.RFC3339)
	parsed, err := ParseTimeRange(input)
	if err != nil {
		t.Fatalf("ParseTimeRange(%q) returned unexpected error: %v", input, err)
	}

	if !parsed.Start.Equal(start) || !parsed.End.Equal(end) {
		t.Errorf("ParseTimeRange(%q) = %#v, want start %v and end %v", input, parsed, start, end)
	}

	_, err = ParseTime("invalid")
	if err == nil {
		t.Error("ParseTime(invalid input) returned nil error")
	}

	_, err = ParseTimeRange("invalid")
	if err == nil {
		t.Error("ParseTimeRange(invalid input) returned nil error")
	}

	custom := ParseTimeWithFormats("02/01/2006")
	tm, err := custom("01/02/2026")
	if err != nil || !tm.Equal(time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("ParseTimeWithFormats returned %v, %v", tm, err)
	}
}
