package main

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestWordCount(t *testing.T) {
	tests := []struct{
		word string
		expected int
	}{
		{"", 0},
		{"               ", 0},
		{"       a       ", 1},
		{"Hello", 1},
		{" Hello", 1},
		{"Hello, world!", 2},
		{"Hello!\nGoodbye!", 2},
		{"Hello!\nGoodbye!\n", 2},
		{"\n\nHello!\n\nGoodbye!\n\n", 2},
	}
	for i, testCase := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
            actual := countWords(testCase.word)
            if actual != testCase.expected {
                t.Errorf("Test Case %d expected %d; got %d", i, testCase.expected, actual)
            }
        })
	}
}

func TestLineCounts(t *testing.T) {
	testCases := []struct{
		input string
		expected int
	}{
		{"Hello", 0},
		{"Hello\n", 1},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(testCase.input))
            result := readSingleReader(reader)
            if result.lines != testCase.expected {
                t.Errorf("Test Case %d expected %d; got %d", i, testCase.expected, result.lines)
            }
        })
	}
}
