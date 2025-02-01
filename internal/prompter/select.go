package prompter

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type stringOption struct {
	Label       string
	Description string
	Icon        string
	Value       string
}

// Select - prompt the user to select an option from a list that answers your question and return their answer
//
// Params:
//   - question: the question to ask the user
//   - hideHelp: hide the help text that explains how to answer the prompt
//   - options:  the options to present to the user
func Select(question string, hideHelp bool, options ...string) (string, error) {
	var items []stringOption
	for _, option := range options {
		items = append(items, stringOption{
			Label:       option,
			Description: "",
			Icon:        "▸ ",
			Value:       option,
		})
	}

	chooser := promptui.Select{
		Label:    question,
		Items:    items,
		HideHelp: hideHelp,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Label }}?",
			Active:   "{{ .Icon | bold | cyan }}{{ .Label }}",
			Inactive: "  {{ .Label }}",
			Details:  "{{ .Description | faint }}",
			Selected: fmt.Sprintf("%s {{ .Label | cyan }}", question),
		},
	}

	i, _, err := chooser.Run()
	if err != nil {
		return "", err
	}

	return options[i], nil
}
