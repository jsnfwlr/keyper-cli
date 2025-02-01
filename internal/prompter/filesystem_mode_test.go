package prompter_test

import (
	"testing"

	"github.com/jsnfwlr/keyper-cli/internal/prompter"
)

func TestFileSystemModes(t *testing.T) {
	testCases := []struct {
		name         string
		actionSets   []map[string]prompter.FSFlag
		expectations map[prompter.FSFlag]bool
	}{
		{
			name: "+FO+SL",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.FileOnly},
				{"Set": prompter.FollowSymlinks},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       true,
				prompter.DirectoryOnly:  false,
				prompter.FollowSymlinks: true,
				prompter.MustExist:      true,
				prompter.MustNotExist:   false,
			},
		},
		{
			name: "+DO+FO+SL",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.DirectoryOnly},
				{"Set": prompter.FileOnly},
				{"Set": prompter.FollowSymlinks},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       true,
				prompter.DirectoryOnly:  false,
				prompter.FollowSymlinks: true,
				prompter.MustExist:      true,
				prompter.MustNotExist:   false,
			},
		},
		{
			name: "+DO+FO+NE+SL",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.DirectoryOnly},
				{"Set": prompter.FileOnly},
				{"Set": prompter.MustNotExist},
				{"Set": prompter.FollowSymlinks},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       false,
				prompter.DirectoryOnly:  false,
				prompter.FollowSymlinks: true,
				prompter.MustExist:      false,
				prompter.MustNotExist:   true,
			},
		},
		{
			name: "+DO+FO+NE+ME",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.DirectoryOnly},
				{"Set": prompter.FileOnly},
				{"Set": prompter.MustNotExist},
				{"Set": prompter.MustExist},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       false,
				prompter.DirectoryOnly:  false,
				prompter.FollowSymlinks: false,
				prompter.MustExist:      true,
				prompter.MustNotExist:   false,
			},
		},
		{
			name: "+DO-DO+NE+ME",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.DirectoryOnly},
				{"Clear": prompter.DirectoryOnly},
				{"Set": prompter.MustNotExist},
				{"Set": prompter.MustExist},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       false,
				prompter.DirectoryOnly:  false,
				prompter.FollowSymlinks: false,
				prompter.MustExist:      true,
				prompter.MustNotExist:   false,
			},
		},
		{
			name: "+NE~DO",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.MustNotExist},
				{"Toggle": prompter.DirectoryOnly},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       false,
				prompter.DirectoryOnly:  true,
				prompter.FollowSymlinks: false,
				prompter.MustExist:      true,
				prompter.MustNotExist:   false,
			},
		},
		{
			name: "+DO~ME",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.DirectoryOnly},
				{"Toggle": prompter.MustExist},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       false,
				prompter.DirectoryOnly:  true,
				prompter.FollowSymlinks: false,
				prompter.MustExist:      false,
				prompter.MustNotExist:   false,
			},
		},
		{
			name: "+DO+NE",
			actionSets: []map[string]prompter.FSFlag{
				{"Set": prompter.DirectoryOnly},
				{"Toggle": prompter.MustNotExist},
			},
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       false,
				prompter.DirectoryOnly:  false,
				prompter.FollowSymlinks: false,
				prompter.MustExist:      false,
				prompter.MustNotExist:   true,
			},
		},
		{
			name: "blank",
			expectations: map[prompter.FSFlag]bool{
				prompter.FileOnly:       false,
				prompter.DirectoryOnly:  false,
				prompter.FollowSymlinks: false,
				prompter.MustExist:      false,
				prompter.MustNotExist:   false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := prompter.FSMode{}
			for _, actionSet := range tc.actionSets {
				for action, mode := range actionSet {
					switch action {
					case "Set":
						m.Set(mode)
					case "Clear":
						m.Clear(mode)
					case "Toggle":
						m.Toggle(mode)
					}
				}
			}

			for flag, expected := range tc.expectations {
				if m.Has(flag) != expected {
					t.Fatalf("expected %t for flag %s, got %t", expected, flag, m.Has(flag))
				}
			}
		})
	}
}
