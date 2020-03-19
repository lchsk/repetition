package main

type Box struct {
	boxNumber   int
	definitions []Definition
}

type Leitner struct {
	boxCount  int
	sessionNo int
	boxes     []Box

	// 0 - 1st box
	// 1 - 1st box, 2nd box
	// n - 1st box, 2nd box, ..., n - 1 box, n box
	stage int

	boxesInCurrentStage []*Box

	movements map[*Definition]int

	currentDefinition *Definition
	currentBox        int
}

func (leitner *Leitner) move() {
	for def, boxNumber := range leitner.movements {
		box := &leitner.boxes[boxNumber]

		box.definitions = append(box.definitions, *def)
	}

	leitner.movements = make(map[*Definition]int)
}

func (leitner *Leitner) isCurrentStageEmpty() bool {
	if len(leitner.boxesInCurrentStage) == 0 {
		return true
	}

	for _, box := range leitner.boxesInCurrentStage {
		if len(box.definitions) > 0 {
			return false
		}
	}

	return true
}

func (leitner *Leitner) setupStage() {
	leitner.boxesInCurrentStage = make([]*Box, leitner.stage+1)

	for i := 0; i <= leitner.stage; i++ {
		leitner.boxesInCurrentStage[i] = &leitner.boxes[i]
	}
}

func (leitner *Leitner) maybeChangeStage() {
	leitner.stage++
	if leitner.stage >= leitner.boxCount {
		leitner.stage = 0
	}
}

func (leitner *Leitner) isFirstBox(box *Box) bool {
	return box.boxNumber == 0
}

func (leitner *Leitner) isLastBox(box *Box) bool {
	return box.boxNumber == leitner.boxCount-1
}

func (leitner *Leitner) getBox(change int) *Box {
	var box *Box = nil

	tries := 0

	for true {
		if tries > leitner.boxCount {
			return nil
		}

		box = &leitner.boxes[leitner.sessionNo%leitner.boxCount]

		if len(box.definitions) == 0 {
			leitner.sessionNo += change
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
	box := &leitner.boxes[(leitner.sessionNo)%leitner.boxCount]

	if leitner.isLastBox(box) {
		return box
	}

	return &leitner.boxes[(leitner.sessionNo+1)%leitner.boxCount]
}

// Calculate previous box in line.
// Selected box doesn't need to have any definitions.
func (leitner *Leitner) getAnyPreviousBox() *Box {
	box := &leitner.boxes[(leitner.sessionNo)%leitner.boxCount]

	if leitner.isFirstBox(box) {
		return box
	}

	return &leitner.boxes[(leitner.sessionNo-1)%leitner.boxCount]
}

func (leitner *Leitner) moveDefinitionToNextBox(definition *Definition) {
	nextBox := leitner.getAnyNextBox()

	nextBox.definitions = append(nextBox.definitions, *definition)
}

func (leitner *Leitner) moveDefinitionToPrevBox(definition *Definition) {
	prevBox := leitner.getAnyPreviousBox()

	prevBox.definitions = append(prevBox.definitions, *definition)
}

func (leitner *Leitner) getDefinition() {
	leitner.currentBox = -1
	leitner.currentDefinition = nil

	for _, box := range leitner.boxesInCurrentStage {
		if len(box.definitions) > 0 {
			leitner.currentBox = box.boxNumber
			leitner.currentDefinition = &box.definitions[0]
			box.definitions = box.definitions[1:]

			break
		}
	}
}

func initLeitner(boxCount int, allDefinitions []Definition) *Leitner {
	boxes := make([]Box, boxCount)

	for i := 0; i < boxCount; i++ {
		boxes[i] = Box{
			boxNumber: i,
		}
	}

	var firstBox *Box = &boxes[0]

	for _, def := range allDefinitions {
		firstBox.definitions = append(firstBox.definitions, def)
	}

	return &Leitner{
		boxCount:  boxCount,
		sessionNo: 0,
		boxes:     boxes,
		// Stage will get set to 0 automatically
		stage:               boxCount - 1,
		boxesInCurrentStage: make([]*Box, 0),
		movements:           make(map[*Definition]int),
		currentDefinition:   nil,
		currentBox:          0,
	}
}
