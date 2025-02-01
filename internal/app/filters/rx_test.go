package filters

import (
	"errors"
	"regexp/syntax"
	"testing"
)

func TestRX(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		input   []byte
		index   int
		want    string
		fail    error
		msg     string
	}{
		{
			name:    "first complete match of one",
			pattern: `(v\d+\.\d+\.\d+)`,
			input:   []byte("Application v1.0.0-beta.1 is available"),
			index:   0,
			want:    "v1.0.0",
			fail:    nil,
			msg:     "",
		},

		{
			name:    "first complete match of two",
			pattern: `(v\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 0,
			want:  "v1.0.0",
			fail:  nil,
			msg:   "",
		},
		{
			name:    "second complete match of two",
			pattern: `(v\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 1,
			want:  "v0.9.0",
			fail:  nil,
			msg:   "",
		},
		{
			name:    "second sub-match of two",
			pattern: `v(\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 1,
			want:  "0.9.0",
			fail:  nil,
			msg:   "",
		},
		{
			name:    "no match",
			pattern: `y(\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 1,
			fail:  RXError{Code: ErrNoMatch, Pattern: `y(\d+\.\d+\.\d+)`},
			msg:   "could not find match for pattern 'y(\\d+\\.\\d+\\.\\d+)' in input",
		},
		{
			name:    "too many sub-matches",
			pattern: `(v)(\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 1,
			fail:  RXError{Code: ErrBadPattern, Pattern: `(v)(\d+\.\d+\.\d+)`, Matches: 2},
			msg:   "too many sub-matches (2) for pattern '(v)(\\d+\\.\\d+\\.\\d+)'",
		},
		{
			name:    "even empty sub-matches count",
			pattern: `(y)?(\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 1,
			fail:  RXError{Code: ErrBadPattern, Pattern: `(y)?(\d+\.\d+\.\d+)`, Matches: 2},
			msg:   "too many sub-matches (2) for pattern '(y)?(\\d+\\.\\d+\\.\\d+)'",
		},
		{
			name:    "unexpected paren",
			pattern: `\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 1,
			fail:  RXError{Raw: RegexSyntaxError(t, syntax.ErrUnexpectedParen, "\\d+\\.\\d+\\.\\d+)")},
			msg:   "error parsing regexp: unexpected ): `\\d+\\.\\d+\\.\\d+)`",
		},
		{
			name:    "missing closing paren",
			pattern: `((\d+\.\d+\.\d+)`,
			input: []byte(`Application v1.0.0-beta.1 released
		Application v0.9.0 released`),
			index: 1,
			fail:  RXError{Raw: RegexSyntaxError(t, syntax.ErrMissingParen, "((\\d+\\.\\d+\\.\\d+)")},
			msg:   "error parsing regexp: missing closing ): `((\\d+\\.\\d+\\.\\d+)`",
		},
		{
			name:    "no matches",
			pattern: `(\d+\.\d+\.\d+)`,
			input:   []byte("Application v1.0.0-beta.1 is available"),
			index:   2,
			fail:    RXError{Code: ErrOutOfRange, Index: 2, Matches: 1, Pattern: `(\d+\.\d+\.\d+)`},
			msg:     "index (2) exceeds number of matches (1) for pattern '(\\d+\\.\\d+\\.\\d+)'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RX(tc.input, tc.pattern, tc.index)

			switch {
			case err == nil && tc.fail != nil:
				t.Errorf("expected error\n%v", tc.fail)
			case err != nil && tc.fail == nil:
				t.Errorf("unexpected error\n%v", err)
			case err != nil && tc.fail != nil && !errors.Is(err, tc.fail):
				t.Errorf("incorrect error:\n\twant %[1]T\n%[1]v\n\tgot %[2]T\n%[2]v", tc.fail, err)
			case err != nil && tc.fail != nil && errors.Is(err, tc.fail) && tc.msg != err.Error():
				t.Errorf("incorrect error message:\n\twant\n%[1]v\n\tgot\n%[2]v", tc.msg, err.Error())
			case err == nil && tc.fail == nil && got != tc.want:
				t.Errorf("incorrect result:\n\twant\n%s\n\tgot\n%s", tc.want, got)

			}
		})
	}
}

func RegexSyntaxError(t *testing.T, code syntax.ErrorCode, expr string) string {
	t.Helper()
	err := &syntax.Error{Code: code, Expr: expr}
	return err.Error()
}
