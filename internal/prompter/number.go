package prompter

import (
	"errors"
	"strconv"

	"github.com/manifoldco/promptui"
)

// Number - prompt the user to enter a number that answers your question and returns their answer
//
// Params:
//   - question: the question to ask the user
//   - prefill: a pre-filled value the user can edit - useful when overwriting existing values or when you want to provide an example value
func Number(question string, prefill int, allowBlank bool) (input int, fault error) {
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

	prompter.Default = strconv.Itoa(prefill)
	prompter.AllowEdit = true

	answer, err := prompter.Run()
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(answer)
	if err != nil {
		return 0, err
	}

	return number, nil
}
