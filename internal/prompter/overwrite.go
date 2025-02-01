package prompter

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// ANP - All None Prompt - the overwrite setting for a field
// All - do not ask the user if they want to overwrite subsequent values, just ask the user for new values and overwrite the existing values
// None - do not ask the user if they want to overwrite subsequent values, just keep the original value, unless it is invalid
// Prompt - ask the user if they want to overwrite subsequent valid values, if they say yes, ask them for a new value, if they say no, keep the original value. Invalid existing values should skip the prompt and just ask the user for a new value.
type ANP int

const (
	OverwriteAll ANP = iota
	OverwriteNone
	PromptBeforeEach
)

type ANPOption struct {
	Label       string
	Description string
	Icon        string
	Value       ANP
}

type Overwriter struct {
	Mode ANP
}

// OverwriteMode - prompt the user to select an overwrite mode and return their answer
// The value returned should be used in subsequent calls to OverwriteBool, OverwriteSelect
// and OverwriteText to determine how to handle the replacement of the original values.
//
// Params:
//   - question: the question to ask the user about overwriting values
//   - hideHelp: hide the help text that explains how to answer the prompt
func NewOverwriter(question string, hideHelp bool) (Overwriter, error) {
	options := []ANPOption{
		{
			Label:       "All",
			Description: "Ignore all existing values and ask for new values for all fields",
			Icon:        "▸ ",
			Value:       OverwriteAll,
		},
		{
			Label:       "None",
			Description: "Only ask for values for new fields",
			Icon:        "▸ ",
			Value:       OverwriteNone,
		},
		{
			Label:       "Prompt",
			Description: "Ask if fields should be overwritten on a field by field basis",
			Icon:        "▸ ",
			Value:       PromptBeforeEach,
		},
	}
	check := promptui.Select{
		Label: question,
		Items: options,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Label }}?",
			Active:   "{{ .Icon | bold | cyan }}{{ .Label }}",
			Inactive: "  {{ .Label }}",
			Details:  "{{ .Description | faint }}",
			Selected: fmt.Sprintf("%s {{ .Label | cyan }}", question),
		},
		HideHelp: hideHelp,
	}

	choice, _, err := check.Run()
	if err != nil {
		return Overwriter{}, err
	}

	ow := Overwriter{
		Mode: options[choice].Value,
	}

	return ow, nil
}

// OverwriteBool - prompt the user to answer a yes/no question and return their answer,
// with special consideration for how existing values should be overwritten
//
// Params:
//   - question:  the question to ask the user
//   - hideHelp:  hide the help text that explains how to answer the overwrite prompt
//   - original:  the original value of the field
//   - prompt:    the function to call to prompt the user for a new value for the field - should be prompt.Bool()
func (ow Overwriter) Bool(question string, hideHelp, original bool, prompt func() bool) (bool, error) {
	switch ow.Mode {
	case OverwriteNone:
		return original, nil
	case OverwriteAll:
		return prompt(), nil

	default:
		check := promptui.Select{
			Label: question,
			Items: []stringOption{
				{
					Label:       "Yes",
					Description: "Enter a new value, overwriting the existing value",
					Icon:        "▸ ",
					Value:       "Yes",
				},
				{
					Label:       "No",
					Description: "Keep the existing value",
					Icon:        "▸ ",
					Value:       "No",
				},
			},
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .Label }}?",
				Active:   "{{ .Icon | bold | cyan }}{{ .Label }}",
				Inactive: "  {{ .Label }}",
				Details:  "{{ .Description | faint }}",
				Selected: fmt.Sprintf("%s {{ .Label | cyan }}", question),
			},
			HideHelp: hideHelp,
		}

		result, _, err := check.Run()
		if err != nil {
			return false, err
		}

		if result == 0 {
			return prompt(), nil
		}

		return original, nil
	}
}

// OverwriteSelect - prompt the user to select an option from a list that answers your
// question and return their answer, with special consideration for how existing values
// should be overwritten
//
// Params:
//   - question:  the question to ask the user
//   - hideHelp:  hide the help text that explains how to answer the overwrite prompt
//   - original:  the original value of the field
//   - prompt:    the function to call to prompt the user for a new value for the field - should be prompt.Select()
func (ow Overwriter) Select(question string, hideHelp bool, original string, prompt func() string) (string, error) {
	switch ow.Mode {
	case OverwriteNone:
		return original, nil
	case OverwriteAll:
		return prompt(), nil

	default:
		check := promptui.Select{
			Label: question,
			Items: []stringOption{
				{
					Label:       "Yes",
					Description: "Enter a new value, overwriting the existing value",
					Icon:        "▸ ",
					Value:       "Yes",
				},
				{
					Label:       "No",
					Description: "Keep the existing value",
					Icon:        "▸ ",
					Value:       "No",
				},
			},
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .Label }}?",
				Active:   "{{ .Icon | bold | cyan }}{{ .Label }}",
				Inactive: "  {{ .Label }}",
				Details:  "{{ .Description | faint }}",
				Selected: fmt.Sprintf("%s {{ .Label | cyan }}", question),
			},
			HideHelp: hideHelp,
		}

		result, _, err := check.Run()
		if err != nil {
			return "", err
		}

		if result == 0 {
			return prompt(), nil
		}

		return original, nil
	}
}

// OverwriteText -prompt the user to enter some text that answers your question and
// return their answer, with special consideration for how existing values should be
// overwritten
//
// Params:
//   - question:  the question to ask the user
//   - hideHelp:  hide the help text that explains how to answer the overwrite prompt
//   - original:  the original value of the field
//   - prompt:    the function to call to prompt the user for a new value for the field - should be prompt.Text()
func (ow Overwriter) Text(question string, hideHelp bool, original string, prompt func() string) (string, error) {
	switch ow.Mode {
	case OverwriteNone:
		return original, nil
	case OverwriteAll:
		return prompt(), nil
	default:
		check := promptui.Select{
			Label: question,
			Items: []stringOption{
				{
					Label:       "Yes",
					Description: "Enter a new value, overwriting the existing value",
					Icon:        "▸ ",
					Value:       "Yes",
				},
				{
					Label:       "No",
					Description: "Keep the existing value",
					Icon:        "▸ ",
					Value:       "No",
				},
			},
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .Label }}?",
				Active:   "{{ .Icon | bold | cyan }}{{ .Label }}",
				Inactive: "  {{ .Label }}",
				Details:  "{{ .Description | faint }}",
				Selected: fmt.Sprintf("%s {{ .Label | cyan }}", question),
			},
			HideHelp: hideHelp,
		}

		result, _, err := check.Run()
		if err != nil {
			return "", err
		}

		if result == 0 {
			return prompt(), nil
		}

		return original, nil
	}
}

// GetOverwritePassword - prompt the user to enter a secret that answers your question and return their answer,
// with special consideration for how existing values should be overwritten. The value entered will be masked
// with asterisks.
//
// Params:
//   - question:  the question to ask the user
//   - hideHelp:  hide the help text that explains how to answer the overwrite prompt
//   - original:  the original value of the field
//   - prompt:    the function to call to prompt the user for a new value for the field - should be prompt.Password()
func (ow Overwriter) Password(question string, hideHelp bool, original string, prompt func() string) (string, error) {
	switch ow.Mode {
	case OverwriteNone:
		return original, nil
	case OverwriteAll:
		return prompt(), nil

	default:
		check := promptui.Select{
			Label: question,
			Items: []stringOption{
				{
					Label:       "Yes",
					Description: "Enter a new value, overwriting the existing value",
					Icon:        "▸ ",
					Value:       "Yes",
				},
				{
					Label:       "No",
					Description: "Keep the existing value",
					Icon:        "▸ ",
					Value:       "No",
				},
			},
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .Label }}?",
				Active:   "{{ .Icon | bold | cyan }}{{ .Label }}",
				Inactive: "  {{ .Label }}",
				Details:  "{{ .Description | faint }}",
				Selected: fmt.Sprintf("%s {{ .Label | cyan }}", question),
			},
			HideHelp: hideHelp,
		}

		result, _, err := check.Run()
		if err != nil {
			return "", err
		}

		if result == 0 {
			return prompt(), nil
		}

		return original, nil
	}
}
