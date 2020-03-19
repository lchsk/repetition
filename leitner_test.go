package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defToGo Definition = Definition{
	from: "andare",
	to:   "to go",
}

var defToBe Definition = Definition{
	from: "essere",
	to:   "to be",
}

var defToSee Definition = Definition{
	from: "vedere",
	to:   "to see",
}

var defToSleep Definition = Definition{
	from: "dormire",
	to:   "to sleep",
}

var definitions []Definition = []Definition{
	{
		from: "andare",
		to:   "to go",
	},
	{
		from: "essere",
		to:   "to be",
	},
	{
		from: "vedere",
		to:   "to see",
	},
	{
		from: "dormire",
		to:   "to sleep",
	},
}

func getTestLeitner() *Leitner {

	return &Leitner{
		boxCount:  3,
		sessionNo: 0,
		boxes: []Box{
			{
				definitions: []Definition{
					{
						from: "andare",
						to:   "to go",
					},
					{
						from: "essere",
						to:   "to be",
					},
					{
						from: "vedere",
						to:   "to see",
					},
					{
						from: "dormire",
						to:   "to sleep",
					},
				},
			},
		},
	}
}

func TestInitLeitner(t *testing.T) {
	leitner := initLeitner(3, definitions)

	assert.Equal(t, len(leitner.boxes), 3)

	box1 := &leitner.boxes[0]
	box2 := &leitner.boxes[1]
	box3 := &leitner.boxes[2]

	assert.Equal(t, 0, box1.boxNumber)
	assert.Equal(t, 1, box2.boxNumber)
	assert.Equal(t, 2, box3.boxNumber)

	assert.Equal(t, definitions, box1.definitions)
	assert.Equal(t, []Definition(nil), box2.definitions)
	assert.Equal(t, []Definition(nil), box3.definitions)
	assert.Equal(t, leitner.boxCount-1, leitner.stage)
	assert.Equal(t, make([]*Box, 0), leitner.boxesInCurrentStage)
	assert.Equal(t, make(map[*Definition]int), leitner.movements)
	assert.Equal(t, (*Definition)(nil), leitner.currentDefinition)
	assert.Equal(t, 0, leitner.currentBox)
}

func TestIsCurrentStageEmpty_initial_state(t *testing.T) {
	leitner := initLeitner(3, definitions)

	assert.True(t, leitner.isCurrentStageEmpty())
}

func TestIsCurrentStageEmpty_no_definitions_in_boxes(t *testing.T) {
	leitner := initLeitner(3, definitions)

	leitner.boxesInCurrentStage = append(leitner.boxesInCurrentStage, &Box{
		definitions: []Definition{},
	})

	assert.True(t, leitner.isCurrentStageEmpty())
}

func TestIsCurrentStageEmpty_definition_in_box(t *testing.T) {
	leitner := initLeitner(3, definitions)

	leitner.boxesInCurrentStage = append(leitner.boxesInCurrentStage, &Box{
		definitions: []Definition{
			defToSleep,
		},
	})

	assert.False(t, leitner.isCurrentStageEmpty())
}

func TestMove(t *testing.T) {
	leitner := initLeitner(3, []Definition{})

	leitner.movements[&defToBe] = 0
	leitner.movements[&defToSleep] = 1
	leitner.movements[&defToGo] = 2
	leitner.movements[&defToSee] = 2

	leitner.move()

	assert.Equal(t, []Definition{defToBe}, leitner.boxes[0].definitions)
	assert.Equal(t, []Definition{defToSleep}, leitner.boxes[1].definitions)

	actual := make(map[Definition]struct{})

	for _, def := range leitner.boxes[2].definitions {
		actual[def] = struct{}{}
	}

	expected := map[Definition]struct{}{
		defToGo:  struct{}{},
		defToSee: struct{}{},
	}

	assert.Equal(t, expected, actual)
}

func TestNextBox(t *testing.T) {
	leitner := initLeitner(3, definitions)

	box1 := &leitner.boxes[0]
	box2 := &leitner.boxes[1]
	box3 := &leitner.boxes[2]

	assert.Equal(t, box1, leitner.getBox(1))
	leitner.sessionNo++

	// Still box1, as other boxes are empty
	assert.Equal(t, box1, leitner.getBox(1))

	box2.definitions = append(box2.definitions, defToBe)
	box3.definitions = append(box2.definitions, defToBe)

	leitner.sessionNo++
	assert.Equal(t, box2, leitner.getBox(1))

	leitner.sessionNo++
	assert.Equal(t, box3, leitner.getBox(1))

	leitner.sessionNo++
	assert.Equal(t, box1, leitner.getBox(1))
}

func TestNextBox__not_getting_into_endless_loop(t *testing.T) {
	// Init with no definitions
	leitner := initLeitner(3, nil)

	assert.Equal(t, (*Box)(nil), leitner.getBox(1))
}

// func TestGetNextBox(t *testing.T) {
// leitner := initLeitner(3, definitions)

// box1 := &leitner.boxes[0]
// box2 := &leitner.boxes[1]
// box3 := &leitner.boxes[2]

// No other boxes have definitions, so it's still box1
// assert.Equal(t, box1, leitner.getNextBox())
// assert.Equal(t, box1, leitner.getPreviousBox())

// box2.definitions = append(box2.definitions, defToBe)

// assert.Equal(t, box2, leitner.getNextBox())
// assert.Equal(t, box1, leitner.getNextBox())
// assert.Equal(t, box2, leitner.getPreviousBox())
// assert.Equal(t, box1, leitner.getPreviousBox())

// box3.definitions = append(box3.definitions, defToBe)

// assert.Equal(t, box2, leitner.getNextBox())
// assert.Equal(t, box3, leitner.getNextBox())
// assert.Equal(t, box1, leitner.getNextBox())
// assert.Equal(t, box3, leitner.getPreviousBox())
// assert.Equal(t, box2, leitner.getPreviousBox())
// assert.Equal(t, box1, leitner.getPreviousBox())
// assert.Equal(t, box3, leitner.getPreviousBox())
// }
