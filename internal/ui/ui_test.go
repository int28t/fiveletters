package ui

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestWinMessage(t *testing.T) {
	tests := []struct {
		name     string
		tries    int
		expected string
	}{
		{"1 try", 1, "Победа! Слово угадано за 1 попытку\n"},
		{"2 tries", 2, "Победа! Слово угадано за 2 попытки\n"},
		{"5 tries", 5, "Победа! Слово угадано за 5 попыток\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := winMessage(tt.tries); result != tt.expected {
				t.Errorf("winMessage(%d) = %q, want %q", tt.tries, result, tt.expected)
			}
		})
	}
}

func TestPrintMenu(t *testing.T) {
	out := &bytes.Buffer{}
	printMenu(out)
	result := out.String()

	if !strings.Contains(result, "1. Новая игра") {
		t.Errorf("printMenu() missing 'Новая игра'")
	}
	if !strings.Contains(result, "2. Выход") {
		t.Errorf("printMenu() missing 'Выход'")
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedStr string
		expectedErr error
	}{
		{"valid", "слово\n", "слово", nil},
		{"multiple words", "два слова\n", "", ErrIncorrectInput},
		{"EOF", "", "", io.EOF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			str, err := readString(scanner)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("readString() error = %v, want %v", err, tt.expectedErr)
			}
			if str != tt.expectedStr {
				t.Errorf("readString() = %v, want %v", str, tt.expectedStr)
			}
		})
	}
}

func TestReadInt(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedInt int
		expectedErr error
	}{
		{"valid", "1\n", 1, nil},
		{"invalid", "abc\n", 0, ErrIncorrectInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			val, err := readInt(scanner)
			if tt.expectedErr != nil && err == nil {
				t.Errorf("readInt() expected error, got nil")
			}
			if val != tt.expectedInt {
				t.Errorf("readInt() = %v, want %v", val, tt.expectedInt)
			}
		})
	}
}

func TestGetMenuChoice(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedChoice int
	}{
		{"valid choice 1", "1\n", 1},
		{"valid choice 2", "2\n", 2},
		{"invalid then valid", "3\n1\n", 1},
		{"EOF", "", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			out := &bytes.Buffer{}
			if choice := getMenuChoice(scanner, out); choice != tt.expectedChoice {
				t.Errorf("getMenuChoice() = %d, want %d", choice, tt.expectedChoice)
			}
		})
	}
}

func TestGetWord(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedWord string
		expectedCode int
	}{
		{"valid word", "кошка\n", "кошка", 0},
		{"invalid then valid", "два слова\nсобака\n", "собака", 0},
		{"EOF", "", "", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			out := &bytes.Buffer{}
			word, code := getWord(scanner, out)

			if word != tt.expectedWord || code != tt.expectedCode {
				t.Errorf("getWord() = (%q, %d), want (%q, %d)", word, code, tt.expectedWord, tt.expectedCode)
			}
		})
	}
}
