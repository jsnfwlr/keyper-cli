package feedback

import (
	"fmt"
	"io"
	"os"
)

var (
	timestampFormat         = "2006-01-02 15:04:05.000"
	outputs         Outputs = Outputs{
		std:   os.Stdout,
		err:   os.Stderr,
		color: true,
	}
)

type Level int

const (
	Required Level = iota // 0
	Error                 // 1
	Warning               // 2
	Info                  // 3
	Extra                 // 4
	Debug                 // 5
	LevelCap Level = 5    // Always at the end and should be equal to the highest value above
)

type Outputs struct {
	std   io.Writer
	err   io.Writer
	color bool
}

func SetTimestampFormat(format string) {
	timestampFormat = format
}

func (l Level) String() string {
	switch l {
	case Required:
		return "Required"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Info:
		return "Info"
	case Extra:
		return "Extra"
	case Debug:
		return "Debug"
	default:
		return "Unknown"
	}
}

func SetOutputs(colorize bool, args ...io.Writer) {
	if len(args) == 0 {
		outputs = Outputs{
			std:   os.Stdout,
			err:   os.Stderr,
			color: colorize,
		}
		return
	}
	if len(args) == 1 {
		outputs = Outputs{
			std:   args[0],
			err:   args[0],
			color: colorize,
		}
		return
	}

	outputs = Outputs{
		std:   args[0],
		err:   args[1],
		color: colorize,
	}
}

func ClearLines(lines int) {
	for i := 0; i < lines; i++ {
		fmt.Fprint(outputs.std, "\033[2K\r")
	}
}

var clearBeforeNextLine = false

func Print(level Level, clearBeforeNext bool, format string, args ...interface{}) {
	if !ExceedsLimit(level) {
		if clearBeforeNextLine {
			fmt.Fprint(outputs.std, "\033[2K\r")
			clearBeforeNextLine = false
		}

		if !clearBeforeNext && format[len(format)-1:] != "\n" {
			format += "\n"
		}

		if clearBeforeNext {
			clearBeforeNextLine = true
		}

		fmt.Fprintf(outputs.std, format, args...)
	}
}

func Printf(level Level, format string, args ...interface{}) {
	if !ExceedsLimit(level) {
		fmt.Fprintf(outputs.std, format, args...)
	}
}

func Println(level Level, args ...interface{}) {
	if !ExceedsLimit(level) {
		fmt.Fprintln(outputs.std, args...)
	}
}
