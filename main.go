package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/logrusorgru/aurora"
)

type Session struct {
	correctAnswers int
	wrongAnswers   int
}

type CommandLine struct {
	debug         *bool
	deckPath      *string
	order         *string
	convertFromKV *string
}

func printDebug(deck *Deck) {
	leitner := deck.Leitner

	fmt.Printf("boxes in stage %d:\n", leitner.Stage)

	for _, box := range leitner.BoxesInCurrentStage {
		fmt.Printf("\t- box %d\tdefinitions: %d\n", box.BoxNumber, len(box.Definitions))
	}

	for _, box := range leitner.Boxes {
		fmt.Println(box.BoxNumber)
		fmt.Println(box.Definitions)
	}

	fmt.Println("Movements")

	for def, boxNumber := range deck.Leitner.movements {
		fmt.Printf("\t%s -> %d\n", def, boxNumber)
	}
}

func readCommandLine() *CommandLine {
	command := CommandLine{}
	command.debug = flag.Bool("debug", false, "Debug mode")
	command.deckPath = flag.String("deck-path", "", "Path to deck file")
	command.order = flag.String("order", "standard", "Question or answer first (standard, reversed, random")
	command.convertFromKV = flag.String("convert-from-kv", "", "Convert file from key-value pairs to deck")

	flag.Parse()

	return &command
}

func getQuestionAnswer(cmd *CommandLine, def *Definition) (string, string) {
	if *cmd.order == "reversed" {
		return def.To, def.From
	}

	if *cmd.order == "random" {
		if rand.Float32() < 0.5 {
			return def.To, def.From
		} else {
			return def.From, def.To
		}
	}

	// order == 'standard'
	return def.From, def.To
}

func saveDeck(deck *Deck, deckPath string) {
	leitner := deck.Leitner

	leitner.movements[leitner.CurrentDefinition] = 0
	leitner.move()

	file, err := json.MarshalIndent(deck, "", " ")

	if err != nil {
		fmt.Printf("Cannot save the deck history file %s\n", err)
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.history.json", deckPath), file, 0644)

	if err != nil {
		fmt.Printf("Cannot save the deck history file %s\n", err)
	}
}

func setupEndOfSessionHandler(session *Session, deck *Deck, deckPath string) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println(aurora.Blue("\nSession summary"))

		fmt.Printf("\tCorrect: %d\n", session.correctAnswers)
		fmt.Printf("\tWrong: %d\n", session.wrongAnswers)

		total := session.correctAnswers + session.wrongAnswers

		if total != 0 {
			fmt.Printf("\tPct: %0.2f%%\n", float64(session.correctAnswers)/float64(total)*100)
		}

		saveDeck(deck, deckPath)

		os.Exit(0)
	}()
}

func prepareQuestion(command *CommandLine, deck *Deck) (bool, string, string) {
	leitner := deck.Leitner

	if leitner.isCurrentStageEmpty() {
		leitner.move()
		leitner.maybeChangeStage()
		leitner.setupStage()

		if leitner.isCurrentStageEmpty() {
			return true, "", ""
		}
	}

	if *command.debug {
		printDebug(deck)
	}

	leitner.getDefinition()

	question, answer := getQuestionAnswer(command, leitner.CurrentDefinition)

	return false, question, answer
}

func recordAnswer(userAnswer string, correctAnswer string, session *Session, leitner *Leitner) {
	def := leitner.CurrentDefinition

	if userAnswer == correctAnswer {
		nextBox := leitner.CurrentBox + 1
		if nextBox >= leitner.BoxCount {
			nextBox = leitner.BoxCount - 1
		}
		leitner.movements[def] = nextBox

		session.correctAnswers++
		fmt.Printf("\n%s\n\n", aurora.Green("============ CORRECT ============"))
	} else {
		prevBox := leitner.CurrentBox - 1
		if prevBox < 0 {
			prevBox = 0
		}
		leitner.movements[def] = prevBox

		session.wrongAnswers++

		fmt.Printf("\n%s\n\n", aurora.Red("============ WRONG ============"))
		fmt.Printf("%s:\n%s\n\n", aurora.Blue("Correct answer"), correctAnswer)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	session := &Session{}

	command := readCommandLine()

	if *command.convertFromKV != "" {
		err := convertKeyValueToDeckFile(*command.convertFromKV)

		if err == nil {
			os.Exit(0)
		}

		fmt.Println(fmt.Errorf("error: %s", err))
		os.Exit(1)
	}

	data, err := loadFile(*command.deckPath)

	if err != nil {
		fmt.Printf("File '%s' does not exist\n", *command.deckPath)
		os.Exit(1)
	}

	// TODO: Load only missing definitions, if the history file is present

	deck := loadDeck(data)
	deck.shuffle()

	jsonFile, _ := os.Open(fmt.Sprintf("%s.history.json", *command.deckPath))
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &deck)

	deck.Definitions = []Definition{}

	setupEndOfSessionHandler(session, deck, *command.deckPath)

	input := bufio.NewScanner(os.Stdin)

	leitner := deck.Leitner

	for true {
		cont, question, answer := prepareQuestion(command, deck)

		if cont {
			continue
		}

		fmt.Printf("%s: \n%s\n\n%s:\n", aurora.Yellow("Question"), question, aurora.Yellow("Answer"))

		input.Scan()

		recordAnswer(input.Text(), answer, session, leitner)
	}
}
