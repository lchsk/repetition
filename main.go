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
	deckPath *string
}

func readCommandLine() *CommandLine {
	command := CommandLine{}
	command.deckPath = flag.String("deck-path", "", "Path to deck file")

	flag.Parse()

	return &command
}

func main() {
	rand.Seed(time.Now().UnixNano())
	session := Session{}

	command := readCommandLine()

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

	data, err := loadFile(*command.deckPath)

	if err != nil {
		fmt.Printf("File '%s' does not exist\n", *command.deckPath)
		os.Exit(1)
	}

	deck := loadDeck(data)
	deck.shuffle()

	input := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Println(deck.leitner.boxes)
		def := deck.leitner.getDefinition()

		// def := deck.getRandomDefinition()
		fmt.Printf("%s: \n%s\n\n%s:\n", aurora.Yellow("Question"), def.from, aurora.Yellow("Answer"))

		input.Scan()

		if input.Text() == def.to {
			deck.leitner.moveDefinitionToNextBox(def)
			session.correctAnswers++
			fmt.Printf("\n%s\n\n", aurora.Green("============ CORRECT ============"))
		} else {
			session.wrongAnswers++

			fmt.Printf("\n%s\n\n", aurora.Red("============ WRONG ============"))
			fmt.Printf("%s:\n%s\n\n", aurora.Blue("Correct answer"), def.to)
		}
	}
}
