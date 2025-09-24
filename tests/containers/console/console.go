package console

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	jsonMessage "github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

type Console struct{}

func Print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		line := jsonMessage.JSONMessage{}
		err := json.Unmarshal([]byte(lastLine), &line)
		if err != nil {
			return err
		}

		if line.Error != nil {
			return line.Error
		}
		_, isTerm := term.GetFdInfo(os.Stdout)
		_ = line.Display(os.Stdout, isTerm)
	}

	errLine := &ErrorLine{}
	if err := json.Unmarshal([]byte(lastLine), errLine); err != nil {
		return err
	}
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func Stream(rd io.Reader) error {
	fd, isTerm := term.GetFdInfo(os.Stdout)
	err := StreamJSON(rd, os.Stdout, fd, isTerm, nil)
	if err != nil {
		return err
	}
	return nil
}

func StreamJSON(in io.Reader, out io.Writer, terminalFd uintptr, isTerminal bool, auxCallback func(JSONMessage)) error {
	var (
		dec = json.NewDecoder(in)
		ids = make(map[string]uint)
	)

	for {
		var diff uint
		var jm JSONMessage
		if err := dec.Decode(&jm); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if jm.Aux != nil {
			if auxCallback != nil {
				auxCallback(jm)
			}
			continue
		}

		if jm.Progress != nil {
			jm.Progress.terminalFd = terminalFd
		}
		if jm.ID != "" && (jm.Progress != nil || jm.ProgressMessage != "") {
			line, ok := ids[jm.ID]
			if !ok {
				// NOTE: This approach of using len(id) to
				// figure out the number of lines of history
				// only works as long as we clear the history
				// when we output something that's not
				// accounted for in the map, such as a line
				// with no ID.
				line = uint(len(ids))
				ids[jm.ID] = line
				if isTerminal {
					_, _ = fmt.Fprintf(out, "\n")
				}
			}
			diff = uint(len(ids)) - line
			if isTerminal {
				cursorUp(out, diff)
			}
		} else {
			// When outputting something that isn't progress
			// output, clear the history of previous lines. We
			// don't want progress entries from some previous
			// operation to be updated (for example, pull -a
			// with multiple tags).
			ids = make(map[string]uint)
		}
		err := jm.Display(out, isTerminal)
		if jm.ID != "" && isTerminal {
			cursorDown(out, diff)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
