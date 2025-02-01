package prompter_test

import (
	"testing"

	"github.com/jsnfwlr/keyper-cli/internal/prompter"
)

func newMode(t *testing.T, flags ...prompter.FSFlag) prompter.FSMode {
	var mode prompter.FSMode
	for _, flag := range flags {
		mode.Set(flag)
	}

	return mode
}

func TestFilesystemValidation(t *testing.T) {
	testCases := []struct {
		name        string
		answer      string
		allowBlank  bool
		mode        prompter.FSMode
		expectError error
	}{
		{
			name:        "blank answer when not allowed",
			answer:      "",
			allowBlank:  false,
			expectError: prompter.ErrBlankNotAllowed,
		},
		{
			name:        "blank answer when allowed",
			answer:      "",
			allowBlank:  true,
			expectError: nil,
		},
		{
			name:        "file only does not exist",
			answer:      "file/does/not/exist",
			allowBlank:  false,
			mode:        newMode(t, prompter.FileOnly),
			expectError: prompter.ErrPathDoesNotExist,
		},
		{
			name:        "dir only does not exist",
			answer:      "dir/does/not/exist",
			allowBlank:  false,
			mode:        newMode(t, prompter.DirectoryOnly),
			expectError: prompter.ErrPathDoesNotExist,
		},
		{
			name:        "blank answer, with allow blank, but must exist",
			answer:      "",
			allowBlank:  true,
			mode:        newMode(t, prompter.MustExist),
			expectError: nil,
		},
		{
			name:        "blank answer, with allow blank, but file only",
			answer:      "",
			allowBlank:  true,
			mode:        newMode(t, prompter.FileOnly),
			expectError: nil,
		},
		{
			name:        "blank answer, with allow blank, but dir only",
			answer:      "",
			allowBlank:  true,
			mode:        newMode(t, prompter.DirectoryOnly),
			expectError: nil,
		},
		{
			name:        "blank answer, with allow blank, but must not exist",
			answer:      "",
			allowBlank:  true,
			mode:        newMode(t, prompter.MustNotExist),
			expectError: nil,
		},
		{
			name:        "file when dir only",
			answer:      "./filesystem_test.go",
			allowBlank:  false,
			mode:        newMode(t, prompter.DirectoryOnly),
			expectError: prompter.ErrPathIsNotDirectory,
		},
		{
			name:        "dir when file only",
			answer:      "../prompt",
			allowBlank:  false,
			mode:        newMode(t, prompter.FileOnly),
			expectError: prompter.ErrPathIsNotFile,
		},

		{
			name:        "file when file only",
			answer:      "./filesystem_test.go",
			allowBlank:  false,
			mode:        newMode(t, prompter.FileOnly),
			expectError: nil,
		},
		{
			name:        "dir when dir only",
			answer:      "../prompt",
			allowBlank:  false,
			mode:        newMode(t, prompter.DirectoryOnly),
			expectError: nil,
		},
		{
			name:        "file when must not exist",
			answer:      "./filesystem_test.go",
			allowBlank:  false,
			mode:        newMode(t, prompter.MustNotExist),
			expectError: prompter.ErrPathExists,
		},
		{
			name:        "dir when must not exist",
			answer:      "../prompt",
			allowBlank:  false,
			mode:        newMode(t, prompter.MustNotExist),
			expectError: prompter.ErrPathExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := prompter.FilesystemValidation{
				AllowBlank: tc.allowBlank,
				Mode:       tc.mode,
			}

			err := v.Process(tc.answer)
			if err != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}
