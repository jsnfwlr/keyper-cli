package filters

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

func JQ(input []byte, pattern string) (string, error) {
	jq, err := gojq.Parse(pattern)
	if err != nil {
		return "", fmt.Errorf("could not compile json filter: %w", err)
	}

	var bodyMap any
	if err := json.Unmarshal(input, &bodyMap); err != nil {
		return "", fmt.Errorf("could not unmarshal json source: %w", err)
	}

	iter := jq.Run(bodyMap)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if fmt.Sprintf("%v", v) != "" && fmt.Sprintf("%v", v) != string(input) {
			return fmt.Sprintf("%v", v), nil
		}
	}
	return "", fmt.Errorf("could not find value with pattern %s in json source", pattern)
}
