package prompter

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

// Bool - prompt the user to answer a yes/no question and return their answer
//
// Params:
//   - question: the question to ask the user
func Bool(question string, prefill bool) bool {
	question = strings.TrimSuffix(question, "?")

	def := "n"
	if prefill {
		def = "y"
	}

	chooser := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
		Validate: func(i string) error {
			if strings.EqualFold(i, "y") || strings.EqualFold(i, "n") {
				return nil
			}
			return errors.New("please entire either 'y' or 'n'")
		},
		Default: def,
	}

	i, _ := chooser.Run()

	return strings.EqualFold(i, "y")
}
