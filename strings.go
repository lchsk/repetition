package main

import "strings"

// Get contents between characters, e.g. everything between ( and ).
// It works even if ( and ) and nested.
func getStringsBetween(data string, delimiter1 rune, delimiter2 rune) []string {
	matches := []string{}

	var braces int
	var begin int

	for index, char := range data {
		if char == delimiter1 {
			if braces == 0 {
				begin = index
			}

			braces++
		}

		if char == delimiter2 {
			if braces == 1 {
				matches = append(matches, strings.TrimSpace(data[begin+1:index]))
			}

			braces--
		}
	}

	return matches
}
