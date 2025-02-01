package filters

import (
	"os"
	"testing"

	"github.com/antchfx/xmlquery"
)

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	b, _ := os.ReadFile(path)
	return b
}

func TestXpath(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		pattern string
		index   int
		want    string
		fail    error
	}{
		{
			name:    "GitLab tags",
			input:   readFile(t, "testdata/gitlab_tags.atom"),
			pattern: "//entry//title",
			index:   0,
			want:    "4.6.1",
			fail:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := XPath(tc.input, tc.pattern, tc.index)
			if err != nil {
				if tc.fail == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if tc.fail != nil {
				t.Fatalf("expected error: %v", tc.fail)
			}

			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestXpathDirect(t *testing.T) {
	url := "https://gitlab.com/chrony/chrony/-/tags?format=atom"

	doc, err := xmlquery.LoadURL(url)
	if err != nil {
		t.Fatalf("could not load url: %v", err)
	}

	nodes, err := xmlquery.QueryAll(doc, "//entry//title")
	if err != nil {
		t.Fatalf("could not match pattern: %v", err)
	}

	if len(nodes) == 0 {
		t.Fatalf("no nodes found for pattern %s", "//entry//title")
	}

	t.Logf("%v", nodes[0].InnerText())
}
