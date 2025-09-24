package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/moby/term"
)

// JSONProgress describes a progress message in a JSON stream.
type JSONProgress struct {
	// Current is the current status and value of the progress made towards Total.
	Current int64 `json:"current,omitempty"`
	// Total is the end value describing when we made 100% progress for an operation.
	Total int64 `json:"total,omitempty"`
	// Start is the initial value for the operation.
	Start int64 `json:"start,omitempty"`
	// HideCounts. if true, hides the progress count indicator (xB/yB).
	HideCounts bool `json:"hidecounts,omitempty"`
	// Units is the unit to print for progress. It defaults to "bytes" if empty.
	Units string `json:"units,omitempty"`

	// terminalFd is the fd of the current terminal, if any. It is used
	// to get the terminal width.
	terminalFd uintptr

	// nowFunc is used to override the current time in tests.
	nowFunc func() time.Time

	// winSize is used to override the terminal width in tests.
	winSize int
}

func (p *JSONProgress) String() string {
	var (
		width       = p.width()
		pbBox       string
		numbersBox  string
		timeLeftBox string
	)
	if p.Current <= 0 && p.Total <= 0 {
		return ""
	}
	if p.Total <= 0 {
		switch p.Units {
		case "":
			return fmt.Sprintf("%8v", units.HumanSize(float64(p.Current)))
		default:
			return fmt.Sprintf("%d %s", p.Current, p.Units)
		}
	}

	percentage := int(float64(p.Current)/float64(p.Total)*100) / 2
	if percentage > 50 {
		percentage = 50
	}
	if width > 110 {
		// this number can't be negative gh#7136
		numSpaces := 0
		if 50-percentage > 0 {
			numSpaces = 50 - percentage
		}
		pbBox = fmt.Sprintf("[%s>%s] ", strings.Repeat("=", percentage), strings.Repeat(" ", numSpaces))
	}

	switch {
	case p.HideCounts:
	case p.Units == "": // no units, use bytes
		current := units.HumanSize(float64(p.Current))
		total := units.HumanSize(float64(p.Total))

		numbersBox = fmt.Sprintf("%8v/%v", current, total)

		if p.Current > p.Total {
			// remove total display if the reported current is wonky.
			numbersBox = fmt.Sprintf("%8v", current)
		}
	default:
		numbersBox = fmt.Sprintf("%d/%d %s", p.Current, p.Total, p.Units)

		if p.Current > p.Total {
			// remove total display if the reported current is wonky.
			numbersBox = fmt.Sprintf("%d %s", p.Current, p.Units)
		}
	}

	if p.Current > 0 && p.Start > 0 && percentage < 50 {
		fromStart := p.now().Sub(time.Unix(p.Start, 0))
		perEntry := fromStart / time.Duration(p.Current)
		left := time.Duration(p.Total-p.Current) * perEntry
		left = (left / time.Second) * time.Second

		if width > 50 {
			timeLeftBox = " " + left.String()
		}
	}
	return pbBox + numbersBox + timeLeftBox
}

// now returns the current time in UTC, but can be overridden in tests
// by setting JSONProgress.nowFunc to a custom function.
func (p *JSONProgress) now() time.Time {
	if p.nowFunc != nil {
		return p.nowFunc()
	}
	return time.Now().UTC()
}

// width returns the current terminal's width, but can be overridden
// in tests by setting JSONProgress.winSize to a non-zero value.
func (p *JSONProgress) width() int {
	if p.winSize != 0 {
		return p.winSize
	}
	ws, err := term.GetWinsize(p.terminalFd)
	if err == nil {
		return int(ws.Width)
	}
	return 200
}
