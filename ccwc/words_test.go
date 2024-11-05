package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"anthonyabunassar.com/cc/wc/input"
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
			reader := bufio.NewReader(strings.NewReader(testCase.word))
			actual := input.ReadSingleReader(reader)
			if actual.Words != testCase.expected {
				t.Errorf("Test Case %d expected %d; got %d", i, testCase.expected, actual.Words)
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
			result := input.ReadSingleReader(reader)
			if result.Lines != testCase.expected {
				t.Errorf("Test Case %d expected %d; got %d", i, testCase.expected, result.Lines)
			}
		})
	}
}

func TestLongLine(t *testing.T) {
	f, err := os.CreateTemp("", "*.txt");
	if err != nil {
		t.Fatal(err)
	}
	w := bufio.NewWriter(f)
	size := 1000000
	for range size {
		w.WriteString("*")
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	// I've verified that this is how wc counts a long file with no spaces and no \n.
	expected := input.WcResult{FileName: f.Name(), Bytes: int64(size), Words: 1, Lines: 0, Chars: int64(size), Width: size}
	actual := input.ReadSingleFileInternal(f.Name())
	if actual != expected {
		t.Errorf("Got %+v from generated file of %d characters.", actual, size)
	}
}
