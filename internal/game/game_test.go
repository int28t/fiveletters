package game

import (
	"errors"
	"slices"
	"testing"
)

func TestValidateWord_Positive(t *testing.T) {
	tests := []struct {
		name string
		word string
	}{
		{
			name: "word len is 5",
			word: "мышка",
		},
		{
			name: "uppercase letters",
			word: "МЫШКА",
		},
		{
			name: "word with cyrillic symbols",
			word: "океан",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWord(tt.word)
			if err != nil {
				t.Errorf("validateWord(%q) error = %v, want <nil>", tt.word, err)
			}
		})
	}
}

func TestValidateWord_Negative(t *testing.T) {
	tests := []struct {
		name        string
		word        string
		expectedErr error
	}{
		{
			name:        "word len is less than 5",
			word:        "кот",
			expectedErr: ErrInvalidWordLen,
		},
		{
			name:        "empty string",
			word:        "",
			expectedErr: ErrInvalidWordLen,
		},
		{
			name:        "word len is more than 5",
			word:        "тобразование",
			expectedErr: ErrInvalidWordLen,
		},
		{
			name:        "word with numbers",
			word:        "ко111",
			expectedErr: ErrWordIsNotCyrillic,
		},
		{
			name:        "word with spec symbols",
			word:        "?()!*",
			expectedErr: ErrWordIsNotCyrillic,
		},
		{
			name:        "word with non-alpha characters",
			word:        "岁月的吻痕",
			expectedErr: ErrWordIsNotCyrillic,
		},
		{
			name:        "word with latin symbols",
			word:        "green",
			expectedErr: ErrWordIsNotCyrillic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWord(tt.word)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("validateWord(%q) error = %v, want %v", tt.word, err, tt.expectedErr)
			}
		})
	}
}

var validDict = []string{"какао", "ангел", "атлас", "базар", "бокал", "веник", "ветка", "горох", "диван", "дрозд", "жизнь", "забор", "зебра", "икона", "искра", "камыш", "лента", "лимон", "мотор", "мышка", "налог", "номер", "олень", "осень", "пирог", "повар", "радио", "рубин", "стена", "сокол", "товар", "туфля", "укроп", "успех", "филин", "форма", "хвост", "хомяк", "цапля", "чехол", "число", "щетка", "щегол", "эклер", "банан", "весна", "гроза", "маска", "пилот", "штора"}

func TestValidateDictionary_Positive(t *testing.T) {
	tests := []struct {
		name       string
		dictionary []string
	}{
		{
			name:       "exactly 50 words",
			dictionary: validDict,
		},
		{
			name:       "> 50 words",
			dictionary: append(slices.Clone(validDict), "арбуз"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDictionary(tt.dictionary)
			if err != nil {
				t.Errorf("validateDictionary() error = %v, want <nil>", err)
			}
		})
	}
}

func TestValidateDictionary_Negative(t *testing.T) {
	tests := []struct {
		name        string
		dictionary  []string
		expectedErr error
	}{
		{
			name:        "dictionary with < 50 words",
			dictionary:  validDict[:49],
			expectedErr: ErrInvalidDictionaryLen,
		},
		{
			name:        "empty dictionary",
			dictionary:  []string{},
			expectedErr: ErrInvalidDictionaryLen,
		},
		{
			name:        "not unique words",
			dictionary:  append(slices.Clone(validDict[:49]), validDict[0]),
			expectedErr: ErrRepeatedDictionaryValues,
		},
		{
			name:        "invalid word inside (numbers)",
			dictionary:  append(slices.Clone(validDict[:49]), "ко111"),
			expectedErr: ErrWordIsNotCyrillic,
		},
		{
			name:        "invalid word inside (length)",
			dictionary:  append(slices.Clone(validDict[:49]), "кот"),
			expectedErr: ErrInvalidWordLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDictionary(tt.dictionary)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("validateDictionary() error = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

func TestStartGame_Positive(t *testing.T) {
	tests := []struct {
		name        string
		dictionary  []string
		maxAttempts int
		seed        int64
		check       func(t *testing.T, s Session)
	}{
		{
			name:        "valid session initialization",
			dictionary:  validDict,
			maxAttempts: 10,
			seed:        28,
			check: func(t *testing.T, s Session) {
				if s.Status() != StatusInProgress {
					t.Errorf("status = %v, want %v", s.Status(), StatusInProgress)
				}
				if s.MaxTries() != 10 {
					t.Errorf("maxTries = %d, want 10", s.MaxTries())
				}
				if s.CurrentTries() != 0 {
					t.Errorf("currentTries = %d, want 0", s.CurrentTries())
				}
				if s.answer == "" {
					t.Errorf("answer is empty")
				}
				if len(s.History()) != 0 {
					t.Errorf("history len = %d, want 0", len(s.History()))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := StartGame(tt.dictionary, tt.maxAttempts, tt.seed)
			if err != nil {
				t.Errorf("StartGame() error = %v, want <nil>", err)
			}

			if err == nil && tt.check != nil {
				tt.check(t, session)
			}
		})
	}
}

func TestStartGame_Negative(t *testing.T) {
	tests := []struct {
		name        string
		dictionary  []string
		maxAttempts int
		seed        int64
		expectedErr error
	}{
		{
			name:        "zero maxAttempts",
			dictionary:  validDict,
			maxAttempts: 0,
			seed:        28,
			expectedErr: ErrInvalidMaxAttempts,
		},
		{
			name:        "negative maxAttempts",
			dictionary:  validDict,
			maxAttempts: -100,
			seed:        28,
			expectedErr: ErrInvalidMaxAttempts,
		},
		{
			name:        "invalid dictionary",
			dictionary:  []string{},
			maxAttempts: 6,
			seed:        28,
			expectedErr: ErrInvalidDictionaryLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StartGame(tt.dictionary, tt.maxAttempts, tt.seed)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("StartGame() error = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

func TestStartGame_Determinism(t *testing.T) {
	seed := int64(28)

	session1, err1 := StartGame(validDict, 6, seed)
	if err1 != nil {
		t.Fatalf("unexpected error: %v", err1)
	}

	session2, err2 := StartGame(validDict, 6, seed)
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}

	if session1.answer != session2.answer {
		t.Errorf("expected answers to be equal, got %q and %q", session1.answer, session2.answer)
	}

	session3, _ := StartGame(validDict, 6, seed+1)
	if session1.answer == session3.answer {
		t.Errorf("expected answers to be different for different seeds")
	}
}

func TestMakeWordState_Positive(t *testing.T) {
	tests := []struct {
		name   string
		answer string
		guess  string
		want   []rune
	}{
		{
			name:   "task requirement: отбор / около",
			answer: "отбор",
			guess:  "около",
			want:   []rune{'✅', '❌', '🟡', '❌', '❌'},
		},
		{
			name:   "all correct",
			answer: "кошка",
			guess:  "кошка",
			want:   []rune{'✅', '✅', '✅', '✅', '✅'},
		},
		{
			name:   "no matches",
			answer: "весна",
			guess:  "поток",
			want:   []rune{'❌', '❌', '❌', '❌', '❌'},
		},
		{
			name:   "mixed matches without duplicates",
			answer: "пирог",
			guess:  "порог",
			want:   []rune{'✅', '❌', '✅', '✅', '✅'},
		},
		{
			name:   "yellow and green matches",
			answer: "сокол",
			guess:  "голос",
			want:   []rune{'❌', '✅', '🟡', '✅', '🟡'},
		},
		{
			name:   "duplicate letters in guess but not in answer",
			answer: "катер",
			guess:  "топот",
			want:   []rune{'🟡', '❌', '❌', '❌', '❌'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := makeWordState(tt.answer, tt.guess)
			if err != nil {
				t.Errorf("makeWordState() error = %v, want <nil>", err)
			}

			if !slices.Equal(tt.want, result) {
				t.Errorf("makeWordState() result = %s, want %s", string(result), string(tt.want))
			}
		})
	}
}

func TestMakeWordState_Negative(t *testing.T) {
	tests := []struct {
		name        string
		answer      string
		guess       string
		expectedErr error
	}{
		{
			name:        "invalid answer length",
			answer:      "кот",
			guess:       "мышка",
			expectedErr: ErrInvalidWordLen,
		},
		{
			name:        "invalid guess chars",
			answer:      "мышка",
			guess:       "ко111",
			expectedErr: ErrWordIsNotCyrillic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := makeWordState(tt.answer, tt.guess)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("makeWordState() error = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

func TestApplyGuess_Positive(t *testing.T) {
	mockDict := map[string]struct{}{
		"кошка": {},
		"мышка": {},
		"потоп": {},
		"какао": {},
	}

	tests := []struct {
		name       string
		setup      func() *Session
		guess      string
		wantResult []rune
		check      func(t *testing.T, s *Session)
	}{
		{
			name: "valid correct guess (win)",
			setup: func() *Session {
				return &Session{
					answer:       "кошка",
					status:       StatusInProgress,
					currentTries: 0,
					maxTries:     6,
					dictionary:   mockDict,
					history:      make([]string, 0, 6),
				}
			},
			guess:      "кошка",
			wantResult: []rune{'✅', '✅', '✅', '✅', '✅'},
			check: func(t *testing.T, s *Session) {
				if s.Status() != StatusWin {
					t.Errorf("status = %v, want %v", s.Status(), StatusWin)
				}
				if s.CurrentTries() != 1 {
					t.Errorf("currentTries = %d, want 1", s.CurrentTries())
				}
				if len(s.History()) != 1 || s.History()[0] != "кошка" {
					t.Errorf("history = %v, want [кошка]", s.History())
				}
			},
		},
		{
			name: "valid incorrect guess (game continues)",
			setup: func() *Session {
				return &Session{
					answer:       "кошка",
					status:       StatusInProgress,
					currentTries: 0,
					maxTries:     6,
					dictionary:   mockDict,
					history:      make([]string, 0, 6),
				}
			},
			guess:      "мышка",
			wantResult: []rune{'❌', '❌', '✅', '✅', '✅'},
			check: func(t *testing.T, s *Session) {
				if s.Status() != StatusInProgress {
					t.Errorf("status = %v, want %v", s.Status(), StatusInProgress)
				}
				if s.CurrentTries() != 1 {
					t.Errorf("currentTries = %d, want 1", s.CurrentTries())
				}
			},
		},
		{
			name: "valid incorrect guess (game over - lose)",
			setup: func() *Session {
				return &Session{
					answer:       "кошка",
					status:       StatusInProgress,
					currentTries: 5,
					maxTries:     6,
					dictionary:   mockDict,
					history:      make([]string, 0, 6),
				}
			},
			guess:      "мышка",
			wantResult: []rune{'❌', '❌', '✅', '✅', '✅'},
			check: func(t *testing.T, s *Session) {
				if s.Status() != StatusLose {
					t.Errorf("status = %v, want %v", s.Status(), StatusLose)
				}
				if s.CurrentTries() != 6 {
					t.Errorf("currentTries = %d, want 6", s.CurrentTries())
				}
			},
		},
		{
			name: "uppercase input is handled correctly",
			setup: func() *Session {
				return &Session{
					answer:       "кошка",
					status:       StatusInProgress,
					currentTries: 0,
					maxTries:     6,
					dictionary:   mockDict,
					history:      make([]string, 0, 6),
				}
			},
			guess:      "МЫШКА",
			wantResult: []rune{'❌', '❌', '✅', '✅', '✅'},
			check: func(t *testing.T, s *Session) {
				if len(s.History()) != 1 || s.History()[0] != "мышка" {
					t.Errorf("history = %v, expected lowered [мышка]", s.History())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := tt.setup()

			res, err := ApplyGuess(session, tt.guess)
			if err != nil {
				t.Errorf("ApplyGuess() error = %v, want <nil>", err)
			}

			if !slices.Equal(tt.wantResult, res.result) {
				t.Errorf("ApplyGuess() result = %s, want %s", string(res.result), string(tt.wantResult))
			}

			if err == nil && tt.check != nil {
				tt.check(t, session)
			}
		})
	}
}

func TestApplyGuess_Negative(t *testing.T) {
	mockDict := map[string]struct{}{
		"кошка": {},
		"мышка": {},
		"потоп": {},
		"какао": {},
	}

	tests := []struct {
		name        string
		setup       func() *Session
		guess       string
		expectedErr error
	}{
		{
			name: "game already finished (win)",
			setup: func() *Session {
				return &Session{status: StatusWin, currentTries: 1, maxTries: 6}
			},
			guess:       "кошка",
			expectedErr: ErrGameAlreadyFinished,
		},
		{
			name: "game already finished (max tries reached)",
			setup: func() *Session {
				return &Session{status: StatusInProgress, currentTries: 6, maxTries: 6}
			},
			guess:       "мышка",
			expectedErr: ErrGameAlreadyFinished,
		},
		{
			name: "invalid guess length",
			setup: func() *Session {
				return &Session{status: StatusInProgress, currentTries: 0, maxTries: 6}
			},
			guess:       "кот",
			expectedErr: ErrInvalidWordLen,
		},
		{
			name: "guess not in dictionary",
			setup: func() *Session {
				return &Session{
					status:       StatusInProgress,
					currentTries: 0,
					maxTries:     6,
					dictionary:   mockDict,
				}
			},
			guess:       "птица",
			expectedErr: ErrWordIsNotInDictionary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := tt.setup()

			_, err := ApplyGuess(session, tt.guess)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ApplyGuess() error = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}
