package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"golang.org/x/exp/constraints"

	_ "embed"

	"github.com/jessevdk/go-flags"
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

type WcOptions struct {
	Chars   bool `short:"m" long:"chars" description:"print the character counts"`
	Bytes   bool `short:"c" long:"bytes" description:"print the byte counts"`
	Lines   bool `short:"l" long:"lines" description:"print the newline counts"`
	Width   bool `short:"L" long:"max-line-length" description:"print the maximum display width"`
	Words   bool `short:"w" long:"words" description:"print the word count"`
	Help    bool `long:"help" description:"display this help and exit"`
	Version bool `long:"version" description:"output version information and exit"`
}

func (o WcOptions) IsDefault() bool {
	return !o.Bytes && !o.Chars && !o.Lines && !o.Width && !o.Words
}

func readSingleFileInternal(f string) WcResult {
	var result WcResult
	result.FileName = f
	inputFile, err := os.Open(f)
	if err != nil {
		result.err = err
		return result
	}

	var r *bufio.Reader
	if f == "-" {
		r = bufio.NewReader(os.Stdin)
	} else {
		// The default buffer size is 4K. Performance test?
		// TODO https://www.reddit.com/r/golang/comments/i1cro6/on_choosing_a_buffer_size/
		// TODO What it a line is longer than the buffer?
		r = bufio.NewReaderSize(inputFile, 65536)
	}

	result = readSingleReader(r)
	result.FileName = f

	if err := inputFile.Close(); err != nil {
		// TODO Do we care?
		result.err = err
	}

	return result
}

func readSingleReader(r *bufio.Reader) WcResult {
	result := WcResult{}
	for {
		// The delimiter is included in the return value. This allows us
		// to emulate wc's behavior, which is to count a line only if it
		// end with the delimiter, and not if it ends with EOF.
		s, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			result.err = err
		}

		if len(s) > 0 {
			result.Words += countWords(s)
			// WARN The byte count from a string.Reader is always 0!
			result.Bytes += int64(r.Buffered())
			result.Chars += int64(len(s))
			result.Width = max(len(s), result.Width)
			if s[len(s)-1] == '\n' {
				result.Lines += 1

			}
		}

		if err == io.EOF {
			break
		}
	}
	return result
}

func countWords(s string) int {
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

func printSingleFile(options WcOptions, w io.Writer, result WcResult) {
	if options.IsDefault() {
		// The long-standing default for wc: "newline, word, and byte counts"
		// TODO Check source for wc and confirm that this is the format.
		fmt.Fprintf(w, "%6d %6d %6d", result.Lines, result.Words, result.Bytes)
		if result.FileName == "-" {
			fmt.Fprintln(w)
		} else {
			fmt.Fprintln(w, " ", result.FileName)
		}
	} else {
		var data []string
		// From the wc help: "newline, word, character, byte, maximum line length"
		if options.Lines {
			data = append(data, fmt.Sprintf("%6d", result.Lines))
		}
		if options.Words {
			data = append(data, fmt.Sprintf("%6d", result.Words))
		}
		if options.Chars {
			data = append(data, fmt.Sprintf("%6d", result.Chars))
		}
		if options.Bytes {
			data = append(data, fmt.Sprintf("%6d", result.Bytes))
		}
		if options.Width {
			data = append(data, fmt.Sprintf("%6d", result.Width))
		}
		fmt.Fprintln(w, strings.Join(data, " "))
	}
}

func countDecimalChars[I constraints.Integer](in I) int {
	// Hat tip to https://stackoverflow.com/a/68124773/476942
	// I've removed the overflow check, since we're not going
	// to encounter numbers > 1e18.
	x, count := I(10), 1
	for x <= in {
		x *= 10
		count++
	}
	return count
}

func adjustPrintWidth(w int) int {
	if w < 6 {
		// This is the default in wc.
		return 7
	}
	return w + 1
}

func printSingleFiles(options WcOptions, results []WcResult, w io.Writer) {
	if len(results) == 1 {
		printSingleFile(options, os.Stdout, results[0])
		return
	}
	// According to https://man7.org/linux/man-pages/man1/wc.1p.html
	// "When any option is specified, wc shall report only the
	// information requested by the specified options."
	var totalWords, totalBytes, totalChars, totalLines int64
	// The minimum width for a number is 7 if we're going to print totals.
	// See https://github.com/coreutils/coreutils/blob/0c9d372c96f2c7ce8c259c5563a48d1816fe611d/src/wc.c#L702
	for _, result := range results {
		totalBytes += result.Bytes
		totalChars += result.Chars
		totalLines += int64(result.Lines)
		totalWords += int64(result.Words)
	}

	maxLinesLength := adjustPrintWidth(countDecimalChars(totalLines))
	maxCharLength := adjustPrintWidth(countDecimalChars(totalChars))
	maxByteLength := adjustPrintWidth(countDecimalChars(totalBytes))
	maxWordsLength := adjustPrintWidth(countDecimalChars(totalWords))
	// There will be more bytes than anything else.
	for _, result := range results {
		fmt.Fprintf(w, "%*d %*d %*d %*d %s\n", maxLinesLength, result.Lines, maxWordsLength, result.Words, maxCharLength, result.Chars, maxByteLength, result.Bytes, result.FileName)
	}

	/*
		By default, the standard output shall contain an entry for each
		input file of the form:

			"%d %d %d %s\n", <newlines>, <words>, <bytes>, <file>

		If the -m option is specified, the number of characters shall
		replace the <bytes> field in this format.
	*/

	fmt.Fprintf(w, "%*d %*d %*d %*d total\n", maxLinesLength, totalLines, maxWordsLength, totalWords, maxCharLength, totalChars, maxByteLength, totalBytes)
}

//go:embed wc-help.md
var usage string

func main() {
	var options WcOptions
	files, err := flags.Parse(&options)

	if err != nil {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	if options.Help {
		// TODO Rewrite the usage for Windows.
		flag.CommandLine.Output().Write([]byte(usage))
		os.Exit(0)
	}

	if options.Version {
		flag.CommandLine.Output().Write([]byte("Version 0.0.1"))
		os.Exit(0)
	}

	var channels []chan WcResult
	if len(files) == 0 {
		result := readSingleFileInternal("-")
		printSingleFile(options, os.Stdout, result)
		os.Exit(0)
	}

	for _, f := range files {
		c := make(chan WcResult)
		channels = append(channels, c)

		go func(f string) {
			c <- readSingleFileInternal(f)
			close(c)
		}(f)
	}

	var results []WcResult
	for _, c := range channels {
		result := <-c
		results = append(results, result)
	}

	printSingleFiles(options, results, os.Stdout)
}
