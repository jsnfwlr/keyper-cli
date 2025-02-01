package style_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jsnfwlr/keyper-cli/internal/feedback"
	"github.com/jsnfwlr/keyper-cli/internal/feedback/style"
)

func TestStyle(t *testing.T) {
	testCases := []struct {
		styles   []style.Attribute
		expected string
	}{
		{
			styles:   []style.Attribute{style.Bold},
			expected: "\033[1mStyle\033[0m Test",
		},
		{
			styles:   []style.Attribute{style.Faint},
			expected: "\033[2mStyle\033[0m Test",
		},
		{
			styles:   []style.Attribute{style.Italic},
			expected: "\033[3mStyle\033[0m Test",
		},
		{
			styles:   []style.Attribute{style.Underline},
			expected: "\033[4mStyle\033[0m Test",
		},
		{
			styles:   []style.Attribute{style.Faint, style.Bold},
			expected: "\033[1;2mStyle\033[0m Test",
		},
		{
			styles:   []style.Attribute{style.RedFG, style.Bold},
			expected: "\033[31;1mStyle\033[0m Test",
		},
		{
			styles:   []style.Attribute{style.RedFG, style.Bold, style.GreenFG},
			expected: "\033[31;1mStyle\033[0m Test",
		},
		{
			styles:   []style.Attribute{style.RedFG, style.Bold, style.GreenBG},
			expected: "\033[31;42;1mStyle\033[0m Test",
		},

		{
			styles:   []style.Attribute{style.GreenFG, style.Faint},
			expected: "\033[32;2mStyle\033[0m Test",
		},
	}

	results := make(map[string]string)
	maxLen := 0
	for _, tc := range testCases {
		labels := []string{}
		for _, style := range tc.styles {
			labels = append(labels, style.Name())
		}

		name := strings.Join(labels, ", ")
		if len(name) > maxLen {
			maxLen = len(name)
		}

		t.Run(strings.ReplaceAll(name, ",", ""), func(t *testing.T) {
			got := style.Apply("Style", tc.styles...) + " Test"
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
			results[name] = got
		})
	}

	pad := fmt.Sprintf("%%-%ds %%s\n", maxLen+2)
	for name, result := range results {
		t.Logf("\n"+pad, name, result)
		fmt.Printf("\t"+pad, name, result)
		feedback.Print(feedback.Required, false, "\t"+pad, name, result)
	}
}
