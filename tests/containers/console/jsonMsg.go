package console

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-git/go-git/v5/plumbing/color"
	"github.com/morikuni/aec"
)

// JSONMessage defines a message struct. It describes
// the created time, where it from, status, ID of the
// message. It's used for docker events.
type JSONMessage struct {
	Stream          string        `json:"stream,omitempty"`
	Status          string        `json:"status,omitempty"`
	Progress        *JSONProgress `json:"progressDetail,omitempty"`
	ProgressMessage string        `json:"progress,omitempty"` // deprecated
	ID              string        `json:"id,omitempty"`
	From            string        `json:"from,omitempty"`
	Time            int64         `json:"time,omitempty"`
	TimeNano        int64         `json:"timeNano,omitempty"`
	Error           *JSONError    `json:"errorDetail,omitempty"`
	ErrorMessage    string        `json:"error,omitempty"` // deprecated
	// Aux contains out-of-band data, such as digests for push signing and image id after building.
	Aux *json.RawMessage `json:"aux,omitempty"`
}

func clearLine(out io.Writer) {
	eraseMode := aec.EraseModes.All
	cl := aec.EraseLine(eraseMode)
	_, _ = fmt.Fprint(out, cl)
}

func cursorUp(out io.Writer, l uint) {
	_, _ = fmt.Fprint(out, aec.Up(l))
}

func cursorDown(out io.Writer, l uint) {
	_, _ = fmt.Fprint(out, aec.Down(l))
}

// Display prints the JSONMessage to out. If isTerminal is true, it erases
// the entire current line when displaying the progressbar. It returns an
// error if the [JSONMessage.Error] field is non-nil.
func (jm *JSONMessage) Display(out io.Writer, isTerminal bool) error {
	if jm.Error != nil {
		return jm.Error
	}
	var endLine string
	if isTerminal && jm.Stream == "" && jm.Progress != nil {
		clearLine(out)
		endLine = "\r"
		_, _ = fmt.Fprint(out, endLine)
	} else if jm.Progress != nil && jm.Progress.String() != "" { // disable progressbar in non-terminal
		return nil
	}
	if jm.TimeNano != 0 {
		_, _ = fmt.Fprintf(out, "%s ", time.Unix(0, jm.TimeNano).Format(RFC3339NanoFixed))
	} else if jm.Time != 0 {
		_, _ = fmt.Fprintf(out, "%s ", time.Unix(jm.Time, 0).Format(RFC3339NanoFixed))
	}
	if jm.ID != "" {
		_, _ = fmt.Fprintf(out, "%s: ", jm.ID)
	}
	if jm.From != "" {
		_, _ = fmt.Fprintf(out, "(from %s) ", jm.From)
	}

	switch {
	case jm.Progress != nil && isTerminal:
		_, _ = fmt.Fprintf(out, "%s%s%s %s%s", color.Green, jm.Status, color.Reset, jm.Progress.String(), endLine)
	case jm.ProgressMessage != "": // deprecated
		_, _ = fmt.Fprintf(out, "%s%s%s %s%s", color.Blue, jm.Status, color.Reset, jm.ProgressMessage, endLine)
	case jm.Stream != "":
		_, _ = fmt.Fprintf(out, "%s%s%s%s", color.Red, jm.Stream, color.Reset, endLine)
	default:
		_, _ = fmt.Fprintf(out, "%s%s%s%s\n", color.Yellow, jm.Status, color.Reset, endLine)
	}
	return nil
}

// DisplayJSONMessagesStream reads a JSON message stream from in, and writes
// each [JSONMessage] to out. It returns an error if an invalid JSONMessage
// is received, or if a JSONMessage containers a non-zero [JSONMessage.Error].
//
// Presentation of the JSONMessage depends on whether a terminal is attached,
// and on the terminal width. Progress bars ([JSONProgress]) are suppressed
// on narrower terminals (< 110 characters).
//
//   - isTerminal describes if out is a terminal, in which case it prints
//     a newline ("\n") at the end of each line and moves the cursor while
//     displaying.
//   - terminalFd is the fd of the current terminal (if any), and used
//     to get the terminal width.
//   - auxCallback allows handling the [JSONMessage.Aux] field. It is
//     called if a JSONMessage contains an Aux field, in which case
//     DisplayJSONMessagesStream does not present the JSONMessage.
func DisplayJSONMessagesStream(in io.Reader, out io.Writer, terminalFd uintptr, isTerminal bool, auxCallback func(JSONMessage)) error {
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
