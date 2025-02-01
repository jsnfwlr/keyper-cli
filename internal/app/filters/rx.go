package filters

import (
	"fmt"
	"regexp"
)

type ErrorCode string

const (
	ErrBadPattern ErrorCode = "bad pattern"
	ErrNoMatch    ErrorCode = "no match"
	ErrOutOfRange ErrorCode = "out of range"
)

type RXError struct {
	Code    ErrorCode
	Pattern string
	Index   int
	Matches int
	Raw     string
}

func (e RXError) Error() string {
	switch e.Code {
	case ErrBadPattern:
		return fmt.Sprintf("too many sub-matches (%d) for pattern '%s'", e.Matches, e.Pattern)
	case ErrNoMatch:
		return fmt.Sprintf("could not find match for pattern '%s' in input", e.Pattern)
	case ErrOutOfRange:
		return fmt.Sprintf("index (%d) exceeds number of matches (%d) for pattern '%s'", e.Index, e.Matches, e.Pattern)
	default:
		return e.Raw
	}
}

func RX(input []byte, pattern string, index int) (val string, fault error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return "", RXError{Raw: err.Error()}
	}

	if !r.MatchString(string(input)) {
		return "", RXError{Pattern: pattern, Code: ErrNoMatch}
	}

	matches := r.FindAllStringSubmatch(string(input), -1)

	if index >= len(matches) {
		return "", RXError{Index: index, Matches: len(matches), Pattern: pattern, Code: ErrOutOfRange}
	}

	if len(matches[index]) > 2 {
		return "", RXError{Matches: len(matches[index]) - 1, Pattern: pattern, Code: ErrBadPattern}
	}

	return matches[index][1], nil
}
