package main

import (
	"fmt"
	"testing"
)

func TestWordCount(t *testing.T) {
	tests := []struct{
		word string
		expected int
	}{
		{"", 0},
		{"               ", 0},
		{"Hello", 1},
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
