package prompter

import (
	"errors"

	"github.com/manifoldco/promptui"
)

// Text - prompt the user to enter some text that answers your question and return their answer
//
// Params:
//   - question: the question to ask the user
//   - prefill: a pre-filled value the user can edit - useful when overwriting existing values or when you want to provide an example value
//   - allowBlank: whether the user can submit a blank value
func Text(question, prefill string, allowBlank bool) (string, error) {
	prompter := promptui.Prompt{
		Label: question,
	}

	if !allowBlank {
		prompter.Validate = func(i string) error {
			if i == "" {
				return errors.New("a non-blank value is required")
			}
			return nil
		}
	}
	if prefill != "" {
		prompter.Default = prefill
		prompter.AllowEdit = true
	}

	result, err := prompter.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}
