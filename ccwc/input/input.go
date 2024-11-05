package input

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

type WcResult struct {
	FileName string
	Bytes    int64
	Chars    int64
	Words    int
	Lines    int
	Width    int
	err      error
}

func CountWords(s string) int {
	/*
		From the man page: "A word is a non-zero-length sequence of characters delimited by white space.

		Let's treat the first non-whitespace character of a line as being preceded by whitespace.
	*/
	wasWhitespace := true
	words := 0
	for _, r := range s {
		if !unicode.IsSpace(r) {
			if wasWhitespace {
				words += 1
			}
			wasWhitespace = false
		} else {
			wasWhitespace = true
		}
	}
	// We've already counted the last word on this line.
	return words
}

func ReadSingleReader(r *bufio.Reader) WcResult {
	result := WcResult{}
	// The file begins on a "word boundary.""
	wasWhitespace := true
	currentWidth := 0

	for {

		r, size, err := r.ReadRune()
		if err != nil {
			if err != io.EOF {
				// TODO What is our plan here?
				result.err = err
			} else {
				// wc does not count EOF as the end of a word or a line, 
				// but it might affect the maximum line width.
				result.Width = max(result.Width, currentWidth)
			}
			break
		}
		currentWidth += 1
		result.Bytes += int64(size)
		result.Chars += 1
		if r == '\n' {
			result.Width = max(result.Width, currentWidth)
			currentWidth = 0
			result.Lines += 1
		}
		if !unicode.IsSpace(r) {
			if wasWhitespace {
				result.Words += 1
			}
			wasWhitespace = false
		} else {
			wasWhitespace = true
		}
	}
	return result
}

func ReadSingleFileInternal(f string) WcResult {
	var result WcResult
	result.FileName = f
	inputFile, err := os.Open(f)
	if err != nil {
		result.err = err
		return result
	}

	var r *bufio.Reader
	if f == "-" || f == "" {
		r = bufio.NewReader(os.Stdin)
	} else {
		// The default buffer size is 4K. Performance test?
		// TODO https://www.reddit.com/r/golang/comments/i1cro6/on_choosing_a_buffer_size/
		// TODO What it a line is longer than the buffer?
		r = bufio.NewReaderSize(inputFile, 65536)
	}

	result = ReadSingleReader(r)
	if f != "" {
		result.FileName = f
	}

	if err := inputFile.Close(); err != nil {
		// TODO Do we care?
		result.err = err
	}

	return result
}