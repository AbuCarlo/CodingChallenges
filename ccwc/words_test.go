package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestWordCount(t *testing.T) {
	tests := []struct {
		word     string
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
	testCases := []struct {
		input    string
		expected int
	}{
		{"Hello", 0},
		{"Hello\n", 1},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(testCase.input))
			result := readSingleReader(reader)
			if result.Lines != testCase.expected {
				t.Errorf("Test Case %d expected %d; got %d", i, testCase.expected, result.Lines)
			}
		})
	}
}

func TestLongLine(t *testing.T) {
	// TODO Extract this to separate function.
	f, err := os.CreateTemp("", "*.txt");
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	size := 1000
	for range size {
		w.WriteString("*")
	}
	w.Flush()
	f.Close()

	expected := WcResult{FileName: f.Name(), Bytes: int64(size), Lines: 0, Chars: int64(size), Width: size}
	actual := readSingleFileInternal(f.Name())
	if actual != expected {
		t.Errorf("Got %+v from generated file of %d characters.", actual, size)
	}
}
