package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func convertKeyValueToDeckFile(path string) error {
	definitions, err := loadKeyValueFile(path)

	if err != nil {
		return err
	}

	if len(definitions) == 0 {
		return errors.New("No definitions found")
	}

	return createDeckFile(path, definitions)
}

func createDeckFile(inputPath string, definitions []Definition) error {
	newPath := fmt.Sprintf("%s.deck", inputPath)

	if _, err := os.Stat(newPath); err == nil {
		return errors.New(fmt.Sprintf("output '%s' file already exists", newPath))
	}

	f, err := os.Create(newPath)
	defer f.Close()

	format := `[
    (%s)
    (%s)
]`

	entries := []string{}

	for _, def := range definitions {
		entry := fmt.Sprintf(format, def.From, def.To)

		entries = append(entries, entry)
	}

	_, err = f.WriteString((strings.Join(entries, "\n")))

	return err
}

func loadKeyValueFile(path string) ([]Definition, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return []Definition{}, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var definitions []Definition

	for scanner.Scan() {
		line := scanner.Text()

		// Ignore empty lines or comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := strings.Split(line, "=")

		if len(parts) == 2 {
			definitions = append(definitions, Definition{
				From: parts[0],
				To:   parts[1],
			})
		}
	}

	return definitions, nil
}
