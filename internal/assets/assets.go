// Package assets provides dictionary JSON parser.
package assets

import (
	_ "embed"
	"encoding/json"
)

//go:embed words.json
var b []byte

// GetWords returns a slice of words from JSON dictionary.
func GetWords() ([]string, error) {
	var words []string
	if err := json.Unmarshal(b, &words); err != nil {
		return nil, err
	}
	return words, nil
}
