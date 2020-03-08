package main

import (
	"io/ioutil"
	"math/rand"
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
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(data), nil
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
