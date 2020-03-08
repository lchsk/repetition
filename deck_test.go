package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestDeck() *Deck {
	return &Deck{
		definitions: []Definition{
			Definition{
				from: "bugs",
				to:   "bunny",
			},
			Definition{
				from: "donald",
				to:   "duck",
			},
			Definition{
				from: "red",
				to:   "sox",
			},
		},
	}
}

func TestLoadDeck_simple_case(t *testing.T) {
	data := `
        [
            (red)
            (sox)
        ]

[(
                bugs
            )

            ( bunny )
        ]
    `

	deck := loadDeck(data)

	expected := []Definition{
		Definition{
			from: "red",
			to:   "sox",
		},
		Definition{
			from: "bugs",
			to:   "bunny",
		},
	}
	actual := deck.definitions

	assert.Equal(t, expected, actual)
}

func TestShuffleDeck(t *testing.T) {
	rand.Seed(10)
	deck := getTestDeck()

	deck.shuffle()

	expected := []Definition{
		Definition{
			from: "red",
			to:   "sox",
		},
		Definition{
			from: "bugs",
			to:   "bunny",
		},
		Definition{
			from: "donald",
			to:   "duck",
		},
	}
	actual := deck.definitions

	assert.Equal(t, expected, actual)
}

func TestGetRandomDefinition(t *testing.T) {
	rand.Seed(16)
	deck := getTestDeck()

	expected := &Definition{from: "bugs", to: "bunny"}
	actual := deck.getRandomDefinition()

	assert.Equal(t, expected, actual)
}
