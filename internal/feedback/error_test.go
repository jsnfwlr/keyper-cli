package feedback_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jsnfwlr/keyper-cli/internal/app"
	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/stretchr/testify/assert"
)

const (
	alphaVer int = 0
	betaVer  int = 1
	rcVer    int = 2
	gaVer    int = 3
)

var buf *bytes.Buffer

func setOutput(t *testing.T) *bytes.Buffer {
	t.Helper()

	if buf != nil {
		return buf
	}

	buf = &bytes.Buffer{}

	feedback.SetOutputs(false, buf, buf)

	return buf
}

func TestErrorsWithOrigin(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("can't run the test if we can't find the current working directory: %s", err)
	}

	buf := setOutput(t)

	baseErr := errors.New("error msg")

	_, _, line, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("can't get the caller")
	}

	testCases := []struct {
		name      string
		err       error
		expectErr string
		expectBuf string
		ver       int
	}{
		{
			name:      "alpha single origin",
			err:       feedback.AddOrigin(baseErr, 1, getOffset(t, line-1)),
			expectErr: fmt.Sprintf("%s:%d - Error: error msg", filepath.Join(cwd, "error_test.go"), line-1),
			expectBuf: fmt.Sprintf("%s:%d - Error: error msg\n", filepath.Join(cwd, "error_test.go"), line-1),
			ver:       alphaVer,
		},
		{
			name:      "beta single origin",
			err:       feedback.AddOrigin(baseErr, 1, getOffset(t, line-1)),
			expectErr: fmt.Sprintf("%s:%d - Error: error msg", filepath.Join(cwd, "error_test.go"), line-1),
			expectBuf: fmt.Sprintf("%s:%d - Error: error msg\n", filepath.Join(cwd, "error_test.go"), line-1),
			ver:       betaVer,
		},
		{
			name:      "rc single origin",
			err:       feedback.AddOrigin(baseErr, 1, getOffset(t, line-1)),
			expectErr: "error msg",
			expectBuf: "error msg\n",
			ver:       rcVer,
		},
		{
			name:      "ga single origin",
			err:       feedback.AddOrigin(baseErr, 1, getOffset(t, line-1)),
			expectErr: "error msg",
			expectBuf: "error msg\n",
			ver:       gaVer,
		},
	}

	for _, tc := range testCases {
		setVer(t, tc.ver)
		t.Log(app.GetCurrentVersion())
		t.Run(fmt.Sprintf("%s stringer", tc.name), func(t *testing.T) {
			if tc.err.Error() != tc.expectErr {
				t.Errorf("expected\n\t`%s`\n, got\n\t`%s`", tc.expectErr, tc.err.Error())
			}
		})
		t.Run(fmt.Sprintf("%s handler", tc.name), func(t *testing.T) {
			feedback.HandleErr(tc.err)

			if buf.String() != tc.expectBuf {
				t.Errorf("expected\n\t`%s`\n, got\n\t`%s`", tc.expectBuf, buf.String())
			}
		})
		buf.Reset()

	}
}

func TestErrorsWithWrap(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("can't run the test if we can't find the current working directory: %s", err)
	}

	buf := setOutput(t)

	testCases := []struct {
		name      string
		err       error
		format    string
		args      []any
		expectErr bool
		expectStr string
		ver       int
	}{
		{
			ver:       alphaVer,
			name:      "alpha err fmt args",
			err:       errors.New("error msg"),
			format:    "%[2]s fmt: %[1]w",
			args:      []any{"arg"},
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: arg fmt: error msg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       alphaVer,
			name:      "alpha err args",
			err:       errors.New("error msg"),
			args:      []any{"arg"},
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: error msg arg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       alphaVer,
			name:      "alpha err",
			err:       errors.New("error msg"),
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: error msg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       alphaVer,
			name:      "alpha err fmt",
			err:       errors.New("error msg"),
			format:    "fmt: %w",
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: fmt: error msg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       alphaVer,
			name:      "alpha nil fmt",
			err:       nil,
			format:    "fmt: %w",
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       alphaVer,
			name:      "alpha nil fmt args",
			err:       nil,
			format:    "fmt: %w",
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       alphaVer,
			name:      "alpha nil args",
			err:       nil,
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       betaVer,
			name:      "beta err fmt args",
			err:       errors.New("error msg"),
			format:    "%[2]s fmt: %[1]w",
			args:      []any{"arg"},
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: arg fmt: error msg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       betaVer,
			name:      "beta err args",
			err:       errors.New("error msg"),
			args:      []any{"arg"},
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: error msg arg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       betaVer,
			name:      "beta err",
			err:       errors.New("error msg"),
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: error msg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       betaVer,
			name:      "beta err fmt",
			err:       errors.New("error msg"),
			format:    "fmt: %w",
			expectErr: true,
			expectStr: fmt.Sprintf("%s:%%d - Error: fmt: error msg\n", filepath.Join(cwd, "error_test.go")),
		},
		{
			ver:       betaVer,
			name:      "beta nil fmt",
			err:       nil,
			format:    "fmt: %w",
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       betaVer,
			name:      "beta nil fmt args",
			err:       nil,
			format:    "fmt: %w",
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       betaVer,
			name:      "beta nil args",
			err:       nil,
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       rcVer,
			name:      "rc err fmt args",
			err:       errors.New("error msg"),
			format:    "%[2]s fmt: %[1]w",
			args:      []any{"arg"},
			expectErr: true,
			expectStr: "Error: arg fmt: error msg\n",
		},
		{
			ver:       rcVer,
			name:      "rc err args",
			err:       errors.New("error msg"),
			args:      []any{"arg"},
			expectErr: true,
			expectStr: "Error: error msg arg\n",
		},
		{
			ver:       rcVer,
			name:      "rc err",
			err:       errors.New("error msg"),
			expectErr: true,
			expectStr: "Error: error msg\n",
		},
		{
			ver:       rcVer,
			name:      "rc err fmt",
			err:       errors.New("error msg"),
			format:    "fmt: %w",
			expectErr: true,
			expectStr: "Error: fmt: error msg\n",
		},
		{
			ver:       rcVer,
			name:      "rc nil fmt",
			err:       nil,
			format:    "fmt: %w",
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       rcVer,
			name:      "rc nil fmt args",
			err:       nil,
			format:    "fmt: %w",
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       rcVer,
			name:      "rc nil args",
			err:       nil,
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       gaVer,
			name:      "ga err fmt args",
			err:       errors.New("error msg"),
			format:    "%[2]s fmt: %[1]w",
			args:      []any{"arg"},
			expectErr: true,
			expectStr: "Error: arg fmt: error msg\n",
		},
		{
			ver:       gaVer,
			name:      "ga err args",
			err:       errors.New("error msg"),
			args:      []any{"arg"},
			expectErr: true,
			expectStr: "Error: error msg arg\n",
		},
		{
			ver:       gaVer,
			name:      "ga err",
			err:       errors.New("error msg"),
			expectErr: true,
			expectStr: "Error: error msg\n",
		},
		{
			ver:       gaVer,
			name:      "ga err fmt",
			err:       errors.New("error msg"),
			format:    "fmt: %w",
			expectErr: true,
			expectStr: "Error: fmt: error msg\n",
		},
		{
			ver:       gaVer,
			name:      "ga nil fmt",
			err:       nil,
			format:    "fmt: %w",
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       gaVer,
			name:      "ga nil fmt args",
			err:       nil,
			format:    "fmt: %w",
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
		{
			ver:       gaVer,
			name:      "ga nil args",
			err:       nil,
			args:      []any{"arg"},
			expectErr: false,
			expectStr: "",
		},
	}

	for _, tc := range testCases {
		setVer(t, tc.ver)

		t.Run(tc.name, func(t *testing.T) {
			feedback.HandleWErr(tc.format, tc.err, tc.args...)

			if !assert.Equal(t, tc.expectErr, tc.err != nil) {
				return
			}

			if tc.expectErr {
				tc.expectStr = updateLine(t, tc.expectStr, 7)
			}

			if !assert.Equal(t, len(tc.expectStr), len(buf.String())) {
				t.Logf("\nwant: %s", tc.expectStr)
				t.Logf("\ngot:  %s", buf.String())
				return
			}

			assert.Regexp(t, tc.expectStr, buf.String())
		})
		buf.Reset()
	}
}

func updateLine(t *testing.T, origin string, offset int) string {
	t.Helper()

	_, _, line, ok := runtime.Caller(1)
	if !ok {
		return origin
	}

	if strings.Contains(origin, "%d") {
		return fmt.Sprintf(origin, line-offset)
	}

	return origin
}

func getOffset(t *testing.T, errLine int) int {
	t.Helper()

	_, _, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("can't get the caller")
	}

	return errLine - line
}

func setVer(t *testing.T, ver int) {
	t.Helper()

	switch ver {
	case alphaVer:
		app.Alpha(t)
	case betaVer:
		app.Beta(t)
	case rcVer:
		app.ReleaseCandidate(t)
	case gaVer:
		app.GeneralAvailability(t)
	}
}
