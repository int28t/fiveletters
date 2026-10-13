// Package ui provides CLI for Five Letters
package ui

import (
	"bufio"
	"errors"
	"fiveletters/internal/game"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// ErrIncorrectInput indicates that input is not correct.
var ErrIncorrectInput = errors.New("incorrect input")

// PlayerStats represents the player statistics
type PlayerStats struct {
	Played int
	Won    int
	Lost   int
}

func printStats(out io.Writer, stats *PlayerStats) {
	winRate := 0
	if stats.Played > 0 {
		winRate = (stats.Won * 100) / stats.Played
	}
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "=== Статистика за этот запуск ===")
	_, _ = fmt.Fprintf(out, "Сыграно: %d | Побед: %d (%d%%) | Поражений: %d\n", stats.Played, stats.Won, winRate, stats.Lost)
	_, _ = fmt.Fprintln(out, "=================================")
}

func colorizeWord(states, guess string) string {
	stateRunes := []rune(states)
	guessRunes := []rune(guess)

	var sb strings.Builder
	for i := range guessRunes {
		switch stateRunes[i] {
		case '✅':
			sb.WriteString(colorGreen)
		case '🟡':
			sb.WriteString(colorYellow)
		case '❌':
			sb.WriteString(colorGray)
		}
		sb.WriteRune(guessRunes[i])
		sb.WriteString(colorReset)
	}
	return sb.String()
}

func printMenu(out io.Writer) {
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "======= 5 букв =======")
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "1. Новая игра")
	_, _ = fmt.Fprintln(out, "2. Выход")
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprint(out, "> Выберите пункт: ")
}

func readString(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}

	input := scanner.Text()
	data := strings.Fields(input)

	if len(data) != 1 {
		return "", ErrIncorrectInput
	}

	return data[0], nil
}

func readInt(scanner *bufio.Scanner) (int, error) {
	s, err := readString(scanner)
	if err != nil {
		return 0, err
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func getMenuChoice(scanner *bufio.Scanner, out io.Writer) int {
	for {
		choice, err := readInt(scanner)

		if errors.Is(err, io.EOF) {
			return 2
		}

		if err == nil && (choice == 1 || choice == 2) {
			return choice
		}
		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintln(out, "Некорректный ввод. Повторите попытку.")
		_, _ = fmt.Fprint(out, "> Выберите пункт: ")
	}
}

func getWord(scanner *bufio.Scanner, out io.Writer) (string, int) {
	for {
		input, err := readString(scanner)

		if errors.Is(err, io.EOF) {
			return "", 2
		}

		if err == nil {
			return input, 0
		}
		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintln(out, "Некорректный ввод. Повторите попытку.")
		_, _ = fmt.Fprint(out, "> Попробуйте угадать слово: ")
	}
}

func runGame(out io.Writer, scanner *bufio.Scanner, session *game.Session) int {
	_, _ = fmt.Fprintln(out)

	if history := session.History(); len(history) > 0 {
		_, _ = fmt.Fprintf(out, "История: %s\n", strings.Join(history, ", "))
	}

	_, _ = fmt.Fprint(out, "> Попробуйте угадать слово: ")
	guess, code := getWord(scanner, out)
	if code != 0 {
		return code
	}
	result, err := game.ApplyGuess(session, guess)
	if err != nil {
		switch {
		case errors.Is(err, game.ErrInvalidWordLen):
			_, _ = fmt.Fprintln(out, "Ошибка: Слово должно состоять из 5 букв.")
		case errors.Is(err, game.ErrWordIsNotInDictionary):
			_, _ = fmt.Fprintln(out, "Ошибка: Такого слова нет в словаре.")
		default:
			_, _ = fmt.Fprintln(out, "Ошибка: Некорректный ввод.")
		}
		return 0
	}
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintf(out, "%s %s\n", result.Result(), colorizeWord(result.Result(), guess))
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintf(out, "Осталось попыток: %d\n", session.MaxTries()-session.CurrentTries())
	return 0
}

func winMessage(tries int) string {
	tryWord := []string{"", "попытку", "попытки", "попытки", "попытки", "попыток", "попыток"}
	if tries > 6 {
		return fmt.Sprintf("Победа! Слово угадано за %d попыток\n", tries)
	}
	return fmt.Sprintf("Победа! Слово угадано за %d %s\n", tries, tryWord[tries])
}

func initGame(out io.Writer, scanner *bufio.Scanner, dictionary []string, seed int64, stats *PlayerStats) int {
	session, err := game.StartGame(dictionary, 6, seed)
	if err != nil {
		_, _ = fmt.Fprintln(out, "Ошибка: не получилось загрузить игру")
		return 1
	}
	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "Программа загадала слово из 5 русских букв")
	for session.Status() == game.StatusInProgress {
		code := runGame(out, scanner, &session)
		if code == 2 {
			return code
		}
	}
	if session.Status() == game.StatusLose {
		stats.Played++
		stats.Lost++
		_, _ = fmt.Fprintf(out, "Неудача. Загаданное слово: %s\n", session.Reveal())
		printStats(out, stats)
	}
	if session.Status() == game.StatusWin {
		stats.Played++
		stats.Won++
		_, _ = fmt.Fprint(out, winMessage(session.CurrentTries()))
		printStats(out, stats)
	}
	return 0
}

// Replay plays the game with prepared list of guesses
func Replay(out io.Writer, dictionary, guessList []string, seed int64) int {
	session, err := game.StartGame(dictionary, 6, seed)
	if err != nil {
		_, _ = fmt.Fprintln(out, "Ошибка: не получилось загрузить игру")
		return 1
	}

	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "Replay:")

	for _, guess := range guessList {
		guess = strings.TrimSpace(guess)

		result, err := game.ApplyGuess(&session, guess)
		if err != nil {
			_, _ = fmt.Fprintln(out, "Ошибка: Некорректный ввод.")
			return 1
		}

		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintf(out, "%s %s\n", result.Result(), colorizeWord(result.Result(), guess))

		if session.Status() != game.StatusInProgress {
			break
		}
	}

	switch session.Status() {
	case game.StatusLose:
		_, _ = fmt.Fprintf(out, "Неудача. Загаданное слово: %s\n", session.Reveal())
	case game.StatusWin:
		_, _ = fmt.Fprint(out, winMessage(session.CurrentTries()))
	case game.StatusInProgress:
		_, _ = fmt.Fprintf(out, "Игра не завершена. Осталось попыток: %d\n", session.MaxTries()-session.CurrentTries())
	}

	return 0
}

// Main starts the main loop of the game.
func Main(in io.Reader, out io.Writer, dictionary []string, seed int64) int {
	scanner := bufio.NewScanner(in)
	stats := &PlayerStats{}

	printMenu(out)
	choice := getMenuChoice(scanner, out)
	if choice == 2 {
		return 0
	}
	for {
		code := initGame(out, scanner, dictionary, seed, stats)
		if code == 1 {
			return 1
		}
		seed++
		if code == 2 {
			return 0
		}
		printMenu(out)
		choice = getMenuChoice(scanner, out)
		if choice == 2 {
			return 0
		}
	}
}
