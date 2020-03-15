package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

type Definition struct {
	from string
	to   string
}

type Deck struct {
	definitions []Definition
	leitner     *Leitner
}

func (deck *Deck) getRandomDefinition() *Definition {
	return &deck.definitions[rand.Intn(len(deck.definitions))]
}

func (deck *Deck) shuffle() {
	rand.Shuffle(len(deck.definitions), func(i, j int) {
		deck.definitions[i], deck.definitions[j] = deck.definitions[j], deck.definitions[i]
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
			from: words[0],
			to:   words[1],
		}

		definitions = append(definitions, definition)
	}

	return &Deck{
		definitions: definitions,
		leitner:     initLeitner(3, definitions),
	}
}
