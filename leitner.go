package main

import "fmt"

type Box struct {
	definitions []Definition
}

type Leitner struct {
	boxCount  int
	sessionNo int
	boxes     []Box
}

func (leitner *Leitner) getBox() *Box {
	var box *Box = nil

	for true {
		box = &leitner.boxes[leitner.sessionNo%leitner.boxCount]

		if len(box.definitions) == 0 {
			leitner.sessionNo++
		} else {
			break
		}
	}

	return box
}

func (leitner *Leitner) getNextBox() *Box {
	return &leitner.boxes[(leitner.sessionNo+1)%leitner.boxCount]
}

func (leitner *Leitner) getPreviousBox() *Box {
	return &leitner.boxes[(leitner.sessionNo-1)%leitner.boxCount]
}

func (leitner *Leitner) moveDefinitionToNextBox(definition *Definition) {
	nextBox := leitner.getNextBox()

	nextBox.definitions = append(nextBox.definitions, *definition)
}

func (leitner *Leitner) getDefinition() *Definition {
	box := leitner.getBox()

	definition := box.definitions[0]
	box.definitions = box.definitions[1:]

	return &definition
}

func initLeitner(boxCount int, allDefinitions []Definition) *Leitner {
	boxes := make([]Box, boxCount)
	var firstBox *Box = &boxes[0]

	for _, def := range allDefinitions {
		fmt.Println(def)
		firstBox.definitions = append(firstBox.definitions, def)
	}

	return &Leitner{
		boxCount:  boxCount,
		sessionNo: 0,
		boxes:     boxes,
	}
}
