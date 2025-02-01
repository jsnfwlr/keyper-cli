package prompter

import (
	"github.com/manifoldco/promptui"
)

// Password - prompt the user to enter a secret that answers your question and return their answer. The value entered will be masked with asterisks.
//
// Params:
//   - question: the question to ask the user
func Password(label string, allowBlank bool) (string, error) {
	prompter := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}

	result, err := prompter.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}
