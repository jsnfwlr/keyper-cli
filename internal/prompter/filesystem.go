package prompter

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jsnfwlr/keyper-cli/internal/files"
	"github.com/manifoldco/promptui"
)

// GetFilesystem - prompt the user to enter a path to either a file or directory according to the question, then return their answer
//
// Params:
//   - question: the question to ask the user
//   - prefill: a pre-filled value the user can edit - useful when overwriting existing values or when you want to provide an example value
//   - allowBlank: whether the user can submit a blank value
//   - mode: the mode to use when checking the response from the user - See filesystem_v.Mode.go for more information
func GetFilesystem(question, prefill string, allowBlank bool, mode FSMode) (string, error) {
	prompter := promptui.Prompt{
		Label: question,
	}

	validate := FilesystemValidation{
		AllowBlank: allowBlank,
		Mode:       mode,
	}

	prompter.Validate = validate.Process

	if prefill != "" {
		prompter.Default = prefill
		prompter.AllowEdit = true
	}

	result, err := prompter.Run()
	if err != nil {
		return "", err
	}

	return filepath.Clean(result), nil
}

type FilesystemValidation struct {
	AllowBlank bool
	Mode       FSMode
}

var (
	ErrBlankNotAllowed    = errors.New("the path cannot be blank")
	ErrPathDoesNotExist   = errors.New("the path provided does not exist")
	ErrPathIsNotDirectory = errors.New("the path provided is not a directory")
	ErrPathIsNotFile      = errors.New("the path provided is not a file")
	ErrPathExists         = errors.New("the path provided already exists")
)

func (v FilesystemValidation) Process(i string) error {
	fmt.Printf("Answer:    \t'%s'\nAllowBlank:\t%t\n", i, v.AllowBlank)

	switch {
	case (!v.AllowBlank && i == ""):
		return ErrBlankNotAllowed
	case i == "":
		return nil
	case v.Mode.MustExist() && !files.Exists(filepath.Clean(i), v.Mode.FollowSymlinks()):
		return ErrPathDoesNotExist
	case v.Mode.MustExist() && files.Exists(filepath.Clean(i), v.Mode.FollowSymlinks()) && v.Mode.DirectoryOnly() && !files.IsDir(filepath.Clean(i), v.Mode.FollowSymlinks()):
		return ErrPathIsNotDirectory
	case v.Mode.MustExist() && files.Exists(filepath.Clean(i), v.Mode.FollowSymlinks()) && v.Mode.FileOnly() && !files.IsFile(filepath.Clean(i), v.Mode.FollowSymlinks()):
		return ErrPathIsNotFile
	case v.Mode.MustNotExist() && files.Exists(filepath.Clean(i), v.Mode.FollowSymlinks()):
		return ErrPathExists
	}

	return nil
}
