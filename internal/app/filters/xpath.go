package filters

import (
	"bytes"
	"fmt"

	"github.com/antchfx/xmlquery"
)

func XPath(input []byte, pattern string, index int) (val string, fault error) {
	doc, err := xmlquery.Parse(bytes.NewReader(input))
	if err != nil {
		return "", fmt.Errorf("could not parse xml: %w", err)
	}

	nodes, err := xmlquery.QueryAll(doc, pattern)
	if err != nil {
		return "", fmt.Errorf("could not match pattern: %w", err)
	}

	if len(nodes) == 0 {
		return "", fmt.Errorf("no nodes found for pattern %s", pattern)
	}

	if index >= len(nodes) {
		return "", fmt.Errorf("index %d out of range for pattern %s", index, pattern)
	}

	return nodes[index].InnerText(), nil
}
