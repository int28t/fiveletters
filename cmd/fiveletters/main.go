// Package main implements the Five Letters console game.
package main

import (
	"fiveletters/internal/assets"
	"fiveletters/internal/ui"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func run() int {
	var seed int64
	var replay string

	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "seed value for random word pick")
	flag.StringVar(&replay, "replay", "", "game replay: words separated by commas")

	flag.Parse()

	if seed < 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: Seed не может быть отрицательным.")
		return 2
	}

	dictionary, err := assets.GetWords()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка: Не получилось прочитать словарь.")
		return 1
	}

	if replay != "" {
		guessList := strings.Split(replay, ",")
		return ui.Replay(os.Stdout, dictionary, guessList, seed)
	}

	return ui.Main(os.Stdin, os.Stdout, dictionary, seed)
}

func main() {
	code := run()
	os.Exit(code)
}
