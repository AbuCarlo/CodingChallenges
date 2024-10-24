package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"golang.org/x/exp/constraints"
)

import _ "embed"

type WcResult struct {
	f     string
	bites int64
	chars int64
	// words int
	lines   int
	longest int
	err     error
}

func readSingleFileInternal(f string) WcResult {
	var result WcResult
	result.f = f
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
		// TODO: https://www.reddit.com/r/golang/comments/i1cro6/on_choosing_a_buffer_size/
		 r = bufio.NewReaderSize(inputFile, 65536)
	}

	for {
		s, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			result.err = err
		}
		if err == io.EOF {
			break
		}
		result.bites += int64(r.Buffered())
		result.chars += int64(len(s))
		result.lines += 1
		result.longest = max(len(s), result.longest)
	}

	if err := inputFile.Close(); err != nil {
		// TODO: Do we care?
		result.err = err
	}

	return result
}

func readSingleFile(f string, c chan<- WcResult) {
	c <- readSingleFileInternal(f)
	// The *sender* must close a channel.
	close(c)
}

func printSingleFile(w io.Writer, result WcResult) {
	fmt.Fprintf(w, "%6d bytes, %6d chars, %6d lines in %s\n", result.bites, result.chars, result.lines, result.f)
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

func printSingleFiles(results []WcResult, w io.Writer) {
	if len(results) == 1 {
		printSingleFile(os.Stdout, results[0])
		return
	}
	// According to https://man7.org/linux/man-pages/man1/wc.1p.html
	// "When any option is specified, wc shall report only the
	// information requested by the specified options."
	var totalBytes, totalChars, totalLines int64
	// The minimum width for a number is 7 if we're going to print totals.
	// See https://github.com/coreutils/coreutils/blob/0c9d372c96f2c7ce8c259c5563a48d1816fe611d/src/wc.c#L702
	for _, result := range results {
		totalBytes += result.bites
		totalChars += result.chars
		totalLines += int64(result.lines)
	}

	maxByteLength := adjustPrintWidth(countDecimalChars(totalBytes))
	maxCharLength := adjustPrintWidth(countDecimalChars(totalChars))
	maxLinesLength := adjustPrintWidth(countDecimalChars(totalLines))
	// There will be more bytes than anything else.
	for _, result := range results {
		fmt.Fprintf(w, "%*d %*d %*d %s\n", maxLinesLength, result.lines, maxCharLength, result.chars, maxByteLength, result.bites, result.f)
	}

	/*

	   By default, the standard output shall contain an entry for each
       input file of the form:

           "%d %d %d %s\n", <newlines>, <words>, <bytes>, <file>

       If the -m option is specified, the number of characters shall
       replace the <bytes> field in this format.
	*/

	fmt.Fprintf(w, "%*d %*d %*d total\n", maxLinesLength, totalLines, maxCharLength, totalChars, maxByteLength, totalBytes)
}

//go:embed wc-help.md
var usage string

func main() {
	// wc prints only byte counts by default.
	help := flag.Bool("h", false, "display this help and exit")
	characters := flag.Bool("c", false, "print the character counts")
	// bites := flag.Bool("b", true, "print the byte counts")
	// newlines := flag.Bool("l", true, "print the newline counts")
	// maximumWidths = flag.Bool("L", false, "print the maximum display width")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	flag.Parse()

	if *help {
		// TODO: Explain the discrepancy from POSIX.
		flag.CommandLine.Output().Write([]byte(usage))
		os.Exit(0)
	}

	// On controlling the usage output better:
	// https://stackoverflow.com/a/23726033/476942

	if !*characters {
		fmt.Fprintln(os.Stderr, "Only the -c option is accepted.")
		// wc returns 1 for invalid options. It also returns 1
		// for missing files.
		os.Exit(1)
	}

	var channels []chan WcResult
	if len(flag.Args()) == 0 {
		c := make(chan WcResult)
		channels = append(channels, c)
		// This is how it appears in columnar output.
		go readSingleFile("-", c)
	}

	for _, f := range flag.Args() {
		c := make(chan WcResult)
		channels = append(channels, c)
		go readSingleFile(f, c)
	}

	var results []WcResult
	for _, c := range channels {
		result := <-c
		results = append(results, result)
	}
	printSingleFiles(results, os.Stdout)
}
