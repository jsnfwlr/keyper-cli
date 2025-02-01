package feedback

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jsnfwlr/keyper-cli/internal/app"
)

func HandleFatalErr(err error) {
	if err == nil {
		return
	}

	HandleErr(err)

	os.Exit(1)
}

func HandleErr(err error) {
	if err == nil {
		return
	}

	if errors.As(err, &ExtendedError{}) {
		fmt.Fprintf(outputs.err, "%v\n", err)
		return
	}

	file, line, ok := getCaller(1)
	if ok {
		fmt.Fprintf(outputs.err, "%s:%d - Error: %+v\n", file, line, err)

		return
	}

	fmt.Fprintf(outputs.err, "Error: %+v\n", err)
}

func HandleWErr(format string, err error, args ...any) {
	if err == nil {
		return
	}

	comboArgs := append([]any{err}, args...)

	switch {
	case format == "" && len(args) == 0:
		HandleErr(err)
	case format != "" && len(args) == 0:
		HandleErr(fmt.Errorf(format, err))
	case format == "" && len(args) > 0:
		format = fmt.Sprintf("%s%s", "%w", strings.Repeat(" %v", len(comboArgs)-1))
		HandleErr(fmt.Errorf(format, comboArgs...))
	default:
		HandleErr(fmt.Errorf(format, comboArgs...))
	}
}

func HandleFatalWErr(format string, err error, args ...any) {
	if err == nil {
		return
	}

	comboArgs := append([]any{err}, args...)

	switch {
	case format == "" && len(args) == 0:
		HandleFatalErr(err)
	case format != "" && len(args) == 0:
		HandleFatalErr(fmt.Errorf(format, err))
	case format == "" && len(args) > 0:
		format = fmt.Sprintf("%s%s", "%w", strings.Repeat(" %v", len(comboArgs)-1))
		HandleFatalErr(fmt.Errorf(format, comboArgs...))
	default:
		HandleFatalErr(fmt.Errorf(format, comboArgs...))
	}
}

type ExtendedError struct {
	Err        error
	OriginFile string
	OriginLine int
}

func isDevBuild() bool {
	return strings.Contains(app.GetCurrentVersion(), "dev") || strings.Contains(app.GetCurrentVersion(), "alpha") || strings.Contains(app.GetCurrentVersion(), "beta")
}

func AddOrigin(err error, depth, offset int) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, &ExtendedError{}) {
		return err
	}

	file, line, ok := getCaller(depth)
	if !ok {
		return ExtendedError{Err: err}
	}

	return ExtendedError{Err: err, OriginFile: file, OriginLine: line + offset}
}

func (e ExtendedError) Error() string {
	if !isDevBuild() || e.OriginFile == "" {
		return e.Err.Error()
	}

	return fmt.Sprintf("%s:%d - Error: %s", e.OriginFile, e.OriginLine, e.Err)
}

// Make sure we the file and line number of the caller don't point to this file
func getCaller(skip int) (string, int, bool) {
	if !isDevBuild() {
		return "", 0, false
	}

	var (
		file string
		line int
		ok   bool
	)

	_, curFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", 0, false
	}

	for {
		_, file, line, ok = runtime.Caller(skip)
		if file != curFile {
			break
		}

		skip++
	}

	return file, line, ok
}
