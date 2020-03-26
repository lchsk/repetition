package main

import (
	"sort"
)

type Box struct {
	BoxNumber   int          `json:"box_number"`
	Definitions []Definition `json:"definitions"`
}

type Leitner struct {
	BoxCount  int   `json:"box_count"`
	SessionNo int   `json:"session_no"`
	Boxes     []Box `json:"boxes"`

	// 0 - 1st box
	// 1 - 1st box, 2nd box
	// n - 1st box, 2nd box, ..., n - 1 box, n box
	Stage int `json:"stage"`

	BoxesInCurrentStage []*Box `json:"-"`

	movements map[*Definition]int

	CurrentDefinition *Definition `json:"-"`
	CurrentBox        int         `json:"-"`
}

func sortDefinitions(leitner *Leitner) {
	for _, box := range leitner.Boxes {
		sort.Slice(box.Definitions, func(i, j int) bool {
			return box.Definitions[i].To < box.Definitions[j].To
		})
	}
}

func (leitner *Leitner) move() {
	for def, boxNumber := range leitner.movements {
		box := &leitner.Boxes[boxNumber]

		box.Definitions = append(box.Definitions, *def)
	}

	sortDefinitions(leitner)

	leitner.movements = make(map[*Definition]int)
}

func (leitner *Leitner) isCurrentStageEmpty() bool {
	if len(leitner.BoxesInCurrentStage) == 0 {
		return true
	}

	for _, box := range leitner.BoxesInCurrentStage {
		if len(box.Definitions) > 0 {
			return false
		}
	}

	return true
}

func (leitner *Leitner) setupStage() {
	leitner.BoxesInCurrentStage = make([]*Box, leitner.Stage+1)

	for i := 0; i <= leitner.Stage; i++ {
		leitner.BoxesInCurrentStage[i] = &leitner.Boxes[i]
	}
}

func (leitner *Leitner) maybeChangeStage() {
	leitner.Stage++
	if leitner.Stage >= leitner.BoxCount {
		leitner.Stage = 0
	}
}

func (leitner *Leitner) isFirstBox(box *Box) bool {
	return box.BoxNumber == 0
}

func (leitner *Leitner) isLastBox(box *Box) bool {
	return box.BoxNumber == leitner.BoxCount-1
}

func (leitner *Leitner) getBox(change int) *Box {
	var box *Box = nil

	tries := 0

	for true {
		if tries > leitner.BoxCount {
			return nil
		}

		box = &leitner.Boxes[leitner.SessionNo%leitner.BoxCount]

		if len(box.Definitions) == 0 {
			leitner.SessionNo += change
		} else {
			break
		}

		tries++
	}

	return box
}

// Calculate next box in line.
// Selected box doesn't need to have any definitions.
func (leitner *Leitner) getAnyNextBox() *Box {
	box := &leitner.Boxes[(leitner.SessionNo)%leitner.BoxCount]

	if leitner.isLastBox(box) {
		return box
	}

	return &leitner.Boxes[(leitner.SessionNo+1)%leitner.BoxCount]
}

// Calculate previous box in line.
// Selected box doesn't need to have any definitions.
func (leitner *Leitner) getAnyPreviousBox() *Box {
	box := &leitner.Boxes[(leitner.SessionNo)%leitner.BoxCount]

	if leitner.isFirstBox(box) {
		return box
	}

	return &leitner.Boxes[(leitner.SessionNo-1)%leitner.BoxCount]
}

func (leitner *Leitner) moveDefinitionToNextBox(definition *Definition) {
	nextBox := leitner.getAnyNextBox()

	nextBox.Definitions = append(nextBox.Definitions, *definition)
}

func (leitner *Leitner) moveDefinitionToPrevBox(definition *Definition) {
	prevBox := leitner.getAnyPreviousBox()

	prevBox.Definitions = append(prevBox.Definitions, *definition)
}

func (leitner *Leitner) getDefinition() {
	leitner.CurrentBox = -1
	leitner.CurrentDefinition = nil

	for _, box := range leitner.BoxesInCurrentStage {
		if len(box.Definitions) > 0 {
			leitner.CurrentBox = box.BoxNumber
			leitner.CurrentDefinition = &box.Definitions[0]
			box.Definitions = box.Definitions[1:]

			break
		}
	}
}

func initLeitner(boxCount int, allDefinitions []Definition) *Leitner {
	boxes := make([]Box, boxCount)

	for i := 0; i < boxCount; i++ {
		boxes[i] = Box{
			BoxNumber:   i,
			Definitions: []Definition{},
		}
	}

	var firstBox *Box = &boxes[0]

	for _, def := range allDefinitions {
		firstBox.Definitions = append(firstBox.Definitions, def)
	}

	return &Leitner{
		BoxCount:  boxCount,
		SessionNo: 0,
		Boxes:     boxes,
		// Stage will get set to 0 automatically
		Stage:               boxCount - 1,
		BoxesInCurrentStage: make([]*Box, 0),
		movements:           make(map[*Definition]int),
		CurrentDefinition:   nil,
		CurrentBox:          0,
	}
}
