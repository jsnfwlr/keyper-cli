package feedback

import (
	"fmt"
	"time"
)

func Log(level Level, format string, args ...interface{}) {
	if !ExceedsLimit(level) {
		format, args := addTimestampLevel(level, time.Now().Format(timestampFormat), format, args...)

		if format[len(format)-1:] != "\n" {
			format += "\n"
		}

		fmt.Fprintf(outputs.std, format, args...)
	}
}

func Logf(level Level, format string, args ...interface{}) {
	if !ExceedsLimit(level) {
		format, args := addTimestampLevel(level, time.Now().Format(timestampFormat), format, args...)

		fmt.Fprintf(outputs.std, format, args...)
	}
}

func Logln(level Level, args ...interface{}) {
	if !ExceedsLimit(level) {
		fmt.Fprintf(outputs.std, "%s - %s - %s\n", time.Now().Format(timestampFormat), level, fmt.Sprint(args...))
	}
}

func addTimestampLevel(level Level, timestamp, format string, args ...interface{}) (string, []interface{}) {
	xArgs := append([]interface{}{timestamp, level}, args...)

	return "%s - %s - " + format, xArgs
}
