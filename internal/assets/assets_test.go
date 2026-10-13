package assets

import (
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestGetWords(t *testing.T) {
	words, err := GetWords()
	if err != nil {
		t.Fatalf("the function returns an error: %v", err)
	}

	// check 50+ words
	if len(words) < 50 {
		t.Errorf("the number of words in dictionary is less than 50")
	}

	// check 5 letters in words
	for _, w := range words {
		if utf8.RuneCountInString(w) != 5 {
			t.Errorf("the number of letters in word is not 5: %q", w)
		}
	}

	// check russian words
	for _, w := range words {
		for _, s := range w {
			if !unicode.IsLetter(s) || !unicode.Is(unicode.Cyrillic, s) {
				t.Errorf("the word is not russian: %q", w)
			}
		}
	}
}
