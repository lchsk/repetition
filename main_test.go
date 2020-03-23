package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDeck() *Deck {
	return &Deck{
		definitions: definitions,
		leitner:     initLeitner(3, definitions),
	}
}

func getCommand(order string) *CommandLine {
	debug := false
	deckPath := ""

	return &CommandLine{
		debug:    &debug,
		deckPath: &deckPath,
		order:    &order,
	}
}

func getSession() *Session {
	return &Session{
		correctAnswers: 0,
		wrongAnswers:   0,
	}
}

func checkBoxes(t *testing.T, leitner *Leitner, definitions1 []Definition, definitions2 []Definition, definitions3 []Definition) {
	box1 := leitner.boxes[0]
	box2 := leitner.boxes[1]
	box3 := leitner.boxes[2]

	assert.Equal(t, definitions1, box1.definitions)
	assert.Equal(t, definitions2, box2.definitions)
	assert.Equal(t, definitions3, box3.definitions)
}

// Go through the question-answer flow.
// The goal is the check whether we're asking some questions (the ones we're getting wrong) more often.
func TestMain__question_answer_flow(t *testing.T) {
	deck := getDeck()
	leitner := deck.leitner

	command := getCommand("standard")
	session := getSession()

	stats := make(map[string]int)

	assert.Equal(t, 2, leitner.stage)
	checkBoxes(t, leitner, []Definition{defToGo, defToBe, defToSee, defToSleep}, []Definition(nil), []Definition(nil))

	// Stage = 0

	_, q, answer := prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 0, leitner.stage)
	assert.Equal(t, &defToBe, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{defToGo, defToSee, defToSleep}, []Definition(nil), []Definition(nil))

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 0, leitner.stage)
	assert.Equal(t, &defToGo, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{defToSee, defToSleep}, []Definition(nil), []Definition(nil))

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 0, leitner.stage)
	assert.Equal(t, &defToSee, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{defToSleep}, []Definition(nil), []Definition(nil))

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 0, leitner.stage)
	assert.Equal(t, &defToSleep, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition(nil), []Definition(nil))

	// Stage = 1

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 1, leitner.stage)
	assert.Equal(t, &defToBe, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{defToGo, defToSee, defToSleep}, []Definition(nil))

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 1, leitner.stage)
	assert.Equal(t, &defToGo, leitner.currentDefinition)
	recordAnswer("wrong", answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{defToSee, defToSleep}, []Definition(nil))

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 1, leitner.stage)
	assert.Equal(t, &defToSee, leitner.currentDefinition)
	recordAnswer("wrong", answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{defToSleep}, []Definition(nil))

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 1, leitner.stage)
	assert.Equal(t, &defToSleep, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{}, []Definition(nil))

	// Stage = 2

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 2, leitner.stage)
	assert.Equal(t, &defToGo, leitner.currentDefinition)
	recordAnswer("wrong", answer, session, leitner)
	checkBoxes(t, leitner, []Definition{defToSee}, []Definition{}, []Definition{defToBe, defToSleep})

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 2, leitner.stage)
	assert.Equal(t, &defToSee, leitner.currentDefinition)
	recordAnswer("wrong", answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{}, []Definition{defToBe, defToSleep})

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 2, leitner.stage)
	assert.Equal(t, &defToBe, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{}, []Definition{defToSleep})

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 2, leitner.stage)
	assert.Equal(t, &defToSleep, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{}, []Definition{})

	// Stage = 0

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 0, leitner.stage)
	assert.Equal(t, &defToGo, leitner.currentDefinition)
	recordAnswer("wrong", answer, session, leitner)
	checkBoxes(t, leitner, []Definition{defToSee}, []Definition{}, []Definition{defToBe, defToSleep})

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 0, leitner.stage)
	assert.Equal(t, &defToSee, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{}, []Definition{defToBe, defToSleep})

	// Stage = 1

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 1, leitner.stage)
	assert.Equal(t, &defToGo, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{defToSee}, []Definition{defToBe, defToSleep})

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 1, leitner.stage)
	assert.Equal(t, &defToSee, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{}, []Definition{defToBe, defToSleep})

	// Stage = 2

	_, q, answer = prepareQuestion(command, deck)
	stats[q]++
	assert.Equal(t, 2, leitner.stage)
	assert.Equal(t, &defToGo, leitner.currentDefinition)
	recordAnswer(answer, answer, session, leitner)
	checkBoxes(t, leitner, []Definition{}, []Definition{}, []Definition{defToBe, defToSee, defToSleep})

	assert.Equal(t, map[string]int{"andare": 6, "dormire": 3, "essere": 3, "vedere": 5}, stats)
}
