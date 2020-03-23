package main

import (
	"bufio"
	"flag"
	"fmt"
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
	debug    *bool
	deckPath *string
	order    *string
}

func printDebug(deck *Deck) {
	leitner := deck.leitner

	fmt.Printf("boxes in stage %d:\n", leitner.stage)

	for _, box := range leitner.boxesInCurrentStage {
		fmt.Printf("\t- box %d\tdefinitions: %d\n", box.boxNumber, len(box.definitions))
	}

	for _, box := range leitner.boxes {
		fmt.Println(box.boxNumber)
		fmt.Println(box.definitions)
	}

	fmt.Println("Movements")

	for def, boxNumber := range deck.leitner.movements {
		fmt.Printf("\t%s -> %d\n", def, boxNumber)
	}
}

func readCommandLine() *CommandLine {
	command := CommandLine{}
	command.debug = flag.Bool("debug", false, "Debug mode")
	command.deckPath = flag.String("deck-path", "", "Path to deck file")
	command.order = flag.String("order", "standard", "Question or answer first (standard, reversed, random")

	flag.Parse()

	return &command
}

func getQuestionAnswer(cmd *CommandLine, def *Definition) (string, string) {
	if *cmd.order == "reversed" {
		return def.to, def.from
	}

	if *cmd.order == "random" {
		if rand.Float32() < 0.5 {
			return def.to, def.from
		} else {
			return def.from, def.to
		}
	}

	// order == 'standard'
	return def.from, def.to
}

func setupEndOfSessionHandler(session *Session) {
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

		os.Exit(0)
	}()
}

func prepareQuestion(command *CommandLine, deck *Deck) (bool, string, string) {
	leitner := deck.leitner

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

	question, answer := getQuestionAnswer(command, leitner.currentDefinition)

	return false, question, answer
}

func recordAnswer(userAnswer string, correctAnswer string, session *Session, leitner *Leitner) {
	def := leitner.currentDefinition

	if userAnswer == correctAnswer {
		nextBox := leitner.currentBox + 1
		if nextBox >= leitner.boxCount {
			nextBox = leitner.boxCount - 1
		}
		leitner.movements[def] = nextBox

		session.correctAnswers++
		fmt.Printf("\n%s\n\n", aurora.Green("============ CORRECT ============"))
	} else {
		prevBox := leitner.currentBox - 1
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
	setupEndOfSessionHandler(session)

	command := readCommandLine()

	data, err := loadFile(*command.deckPath)

	if err != nil {
		fmt.Printf("File '%s' does not exist\n", *command.deckPath)
		os.Exit(1)
	}

	deck := loadDeck(data)
	deck.shuffle()

	input := bufio.NewScanner(os.Stdin)

	leitner := deck.leitner

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
