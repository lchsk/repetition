package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStringsBetween_simple_case(t *testing.T) {
	matches := getStringsBetween("  (test)  ", '(', ')')

	assert.Equal(t, matches, []string{"test"})
}

func TestGetStringsBetween_simple_case_with_spaces(t *testing.T) {
	matches := getStringsBetween("  ( test )  ", '(', ')')

	assert.Equal(t, matches, []string{"test"})
}

func TestGetStringsBetween_nested_delimiters(t *testing.T) {
	matches := getStringsBetween("  ( (test) )  ", '(', ')')

	assert.Equal(t, matches, []string{"(test)"})
}

func TestGetStringsBetween_empty(t *testing.T) {
	matches := getStringsBetween("", '(', ')')

	assert.Equal(t, matches, []string{})
}

func TestGetStringsBetween_broken_input(t *testing.T) {
	matches := getStringsBetween("(test", '(', ')')

	assert.Equal(t, matches, []string{})
}
