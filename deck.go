package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

type Definition struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Deck struct {
	Definitions []Definition `json:"-"`
	Leitner     *Leitner     `json:"leitner"`
}

func (deck *Deck) getRandomDefinition() *Definition {
	return &deck.Definitions[rand.Intn(len(deck.Definitions))]
}

func (deck *Deck) shuffle() {
	rand.Shuffle(len(deck.Definitions), func(i, j int) {
		deck.Definitions[i], deck.Definitions[j] = deck.Definitions[j], deck.Definitions[i]
	})
}

func loadFile(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		line := scanner.Text()

		// Ignore empty lines or comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), err
}

func loadDeck(data string) *Deck {
	groups := getStringsBetween(data, '[', ']')

	var definitions []Definition

	for _, group := range groups {
		words := getStringsBetween(group, '(', ')')

		definition := Definition{
			From: words[0],
			To:   words[1],
		}

		definitions = append(definitions, definition)
	}

	return &Deck{
		Definitions: definitions,
		Leitner:     initLeitner(3, definitions),
	}
}
