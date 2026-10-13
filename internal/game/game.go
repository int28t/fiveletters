// Package game provides game engine for Five Letters.
package game

import (
	"errors"
	"math/rand"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Status is a current game state.
type Status int

const (
	// StatusInProgress shows that game session is still running.
	StatusInProgress Status = iota
	// StatusWin shows that game ends with user's win.
	StatusWin
	// StatusLose shows that game ends with user's lose.
	StatusLose
)

// Session represents the state of one game.
type Session struct {
	answer       string
	maxTries     int
	history      []string
	currentTries int
	status       Status
	dictionary   map[string]struct{}
}

//nolint:revive
type GameSession = Session

// Status is a getter for game session's status.
func (g *Session) Status() Status {
	return g.status
}

// MaxTries is a getter for game session's max tries number.
func (g *Session) MaxTries() int {
	return g.maxTries
}

// CurrentTries is a getter for game session's current tries number.
func (g *Session) CurrentTries() int {
	return g.currentTries
}

// History is a getter for game session's history.
func (g *Session) History() []string {
	if g.history == nil {
		return nil
	}
	cp := slices.Clone(g.history)
	return cp
}

// Reveal is a getter for game session's answer
func (g *Session) Reveal() string {
	return g.answer
}

// GuessResult represents the result of one accepted guess.
type GuessResult struct {
	result []rune
}

// Result is a getter for GuessResult result
func (g GuessResult) Result() string {
	return string(g.result)
}

var (
	// ErrInvalidWordLen indicates that the word length is not exactly 5 letters.
	ErrInvalidWordLen = errors.New("the number of letters in word is not 5")

	// ErrWordIsNotCyrillic indicates that the word contains non-Cyrillic characters.
	ErrWordIsNotCyrillic = errors.New("the word is not сyrillic")

	// ErrInvalidDictionaryLen indicates that the dictionary contains fewer than 50 words.
	ErrInvalidDictionaryLen = errors.New("dictionary should have at least 50 words")

	// ErrRepeatedDictionaryValues indicates that the dictionary contains duplicate words.
	ErrRepeatedDictionaryValues = errors.New("dictionary should not have repeated values")

	// ErrInvalidMaxAttempts indicates that the maximum number of attempts is zero or negative.
	ErrInvalidMaxAttempts = errors.New("maxAttempts should be > 0")

	// ErrGameAlreadyFinished indicates that the game is over and no more attempts are allowed.
	ErrGameAlreadyFinished = errors.New("game is already finished")

	// ErrWordIsNotInDictionary indicates that the word is not in the dictionary.
	ErrWordIsNotInDictionary = errors.New("the word is not in the dictionary")
)

func validateWord(w string) error {
	if utf8.RuneCountInString(w) != 5 {
		return ErrInvalidWordLen
	}
	for _, s := range w {
		if !unicode.IsLetter(s) || !unicode.Is(unicode.Cyrillic, s) {
			return ErrWordIsNotCyrillic
		}
	}
	return nil
}

func validateDictionary(d []string) error {
	if len(d) < 50 {
		return ErrInvalidDictionaryLen
	}
	cp := slices.Clone(d)
	slices.Sort(cp)
	unique := slices.Compact(cp)
	if len(unique) != len(d) {
		return ErrRepeatedDictionaryValues
	}
	for _, w := range d {
		if err := validateWord(w); err != nil {
			return err
		}
	}
	return nil
}

// StartGame creates a new game session.
func StartGame(dictionary []string, maxAttempts int, seed int64) (Session, error) {
	if maxAttempts <= 0 {
		return Session{}, ErrInvalidMaxAttempts
	}

	if err := validateDictionary(dictionary); err != nil {
		return Session{}, err
	}

	source := rand.NewSource(seed)
	r := rand.New(source)
	randomIndex := r.Intn(len(dictionary))
	word := strings.ToLower(dictionary[randomIndex])

	mp := make(map[string]struct{})
	for _, w := range dictionary {
		mp[strings.ToLower(w)] = struct{}{}
	}

	session := Session{
		answer:       word,
		maxTries:     maxAttempts,
		history:      make([]string, 0, maxAttempts),
		currentTries: 0,
		status:       StatusInProgress,
		dictionary:   mp,
	}

	return session, nil
}

func makeWordState(answer, guess string) ([]rune, error) {
	answerRunes := []rune(answer)
	guessRunes := []rune(guess)
	if err := validateWord(answer); err != nil {
		return nil, err
	}
	if err := validateWord(guess); err != nil {
		return nil, err
	}

	result := make([]rune, 5)

	var used [5]bool

	for i := range 5 {
		if guessRunes[i] == answerRunes[i] {
			result[i] = '✅'
			used[i] = true
		}
	}

	for i := range 5 {
		if guessRunes[i] == answerRunes[i] {
			continue
		}
		for j := range 5 {
			if (guessRunes[i] == answerRunes[j]) && !used[j] {
				used[j] = true
				result[i] = '🟡'
				break
			}
		}
		if result[i] == 0 {
			result[i] = '❌'
		}
	}

	return result, nil
}

// ApplyGuess validates and applies a guess to the session.
func ApplyGuess(session *Session, guess string) (GuessResult, error) {
	if session.status != StatusInProgress || session.currentTries == session.maxTries {
		return GuessResult{}, ErrGameAlreadyFinished
	}

	if utf8.RuneCountInString(guess) != 5 {
		return GuessResult{}, ErrInvalidWordLen
	}

	guess = strings.ToLower(guess)

	_, ok := session.dictionary[guess]

	if !ok {
		return GuessResult{}, ErrWordIsNotInDictionary
	}

	session.currentTries++
	session.history = append(session.history, guess)

	if guess == session.answer {
		session.status = StatusWin
	} else if session.currentTries == session.maxTries {
		session.status = StatusLose
	}

	result, err := makeWordState(session.answer, guess)
	if err != nil {
		return GuessResult{}, err
	}

	return GuessResult{
		result: result,
	}, nil
}
