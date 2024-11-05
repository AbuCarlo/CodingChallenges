package main

// -ldflags="-X 'main.Version=v1.0.0'"

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"golang.org/x/exp/constraints"

	_ "embed"

	"github.com/jessevdk/go-flags"

	"anthonyabunassar.com/cc/wc/input"
)

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

func printSingleFile(options WcOptions, w io.Writer, result input.WcResult) {
	if options.IsDefault() {
		// The long-standing default for wc: "newline, word, and byte counts"
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

func printMultipleFiles(options WcOptions, results []input.WcResult, w io.Writer) {
	if len(results) == 1 {
		printSingleFile(options, os.Stdout, results[0])
		return
	}
	// According to https://man7.org/linux/man-pages/man1/wc.1p.html
	// "When any option is specified, wc shall report only the
	// information requested by the specified options."
	var totalWords, totalBytes, totalChars, totalLines int64
	var longestLine int
	// The minimum width for a number is 7 if we're going to print totals.
	// See https://github.com/coreutils/coreutils/blob/0c9d372c96f2c7ce8c259c5563a48d1816fe611d/src/wc.c#L702
	for _, result := range results {
		totalBytes += result.Bytes
		totalChars += result.Chars
		totalLines += int64(result.Lines)
		totalWords += int64(result.Words)
		longestLine = max(longestLine, result.Width)
	}

	maxLinesLength := adjustPrintWidth(countDecimalChars(totalLines))
	maxCharLength := adjustPrintWidth(countDecimalChars(totalChars))
	maxByteLength := adjustPrintWidth(countDecimalChars(totalBytes))
	maxWordsLength := adjustPrintWidth(countDecimalChars(totalWords))
	maxWidthLength := adjustPrintWidth(countDecimalChars(longestLine))

	// There will be more bytes than anything else.
	for _, result := range results {

		if options.IsDefault() {
			// The long-standing default for wc: "newline, word, and byte counts"
			fmt.Fprintf(w, "%*d %*d %*d", maxLinesLength, result.Lines, maxWordsLength, result.Words, maxByteLength, result.Bytes)
			if result.FileName == "-" {
				fmt.Fprintln(w)
			} else {
				fmt.Fprintln(w, " ", result.FileName)
			}
		} else {
			{
				var data []string
				// From the wc help: "newline, word, character, byte, maximum line length"
				if options.Lines {
					data = append(data, fmt.Sprintf("%*d", maxLinesLength, result.Lines))
				}
				if options.Words {
					data = append(data, fmt.Sprintf("%*d", maxWordsLength, result.Words))
				}
				if options.Chars {
					data = append(data, fmt.Sprintf("%*d", maxCharLength, result.Chars))
				}
				if options.Bytes {
					data = append(data, fmt.Sprintf("%*d", maxByteLength, result.Bytes))
				}
				if options.Width {
					data = append(data, fmt.Sprintf("%*d", maxWidthLength, result.Width))
				}
				fmt.Fprintln(w, strings.Join(data, " "))
			}
		}
	}

	// Print totals.

	if options.IsDefault() {
		fmt.Fprintf(w, "%*d %*d %*d total\n", maxLinesLength, totalLines, maxWordsLength, totalWords, maxByteLength, totalBytes)
	} else {
		var data []string
		if options.Lines {
			data = append(data, fmt.Sprintf("%*d", maxLinesLength, totalLines))
		}
		if options.Words {
			data = append(data, fmt.Sprintf("%*d", maxWordsLength, totalWords))
		}
		if options.Chars {
			data = append(data, fmt.Sprintf("%*d", maxCharLength, totalChars))
		}
		if options.Bytes {
			data = append(data, fmt.Sprintf("%*d", maxByteLength, totalBytes))
		}
		if options.Width {
			data = append(data, fmt.Sprintf("%*d", maxWidthLength, longestLine))
		}
		fmt.Fprint(w, strings.Join(data, " "))
		fmt.Fprintln(w, " total")
	}

	/*
		By default, the standard output shall contain an entry for each
		input file of the form:

			"%d %d %d %s\n", <newlines>, <words>, <bytes>, <file>

		If the -m option is specified, the number of characters shall
		replace the <bytes> field in this format.
	*/
}

//go:embed version-preface.txt
var versionPreface string

//go:embed gnu-help.txt
var gnuUsage string

//go:embed gnu-version.txt
var gnuVersion string

// WARN Unfortunately, the Windows shell will create this as a UTF-16 file.
// Regenerate it thus:  git describe --always --long --dirty --tags > version.txt

//go:embed version.txt
var ccVersion string

func main() {
	var options WcOptions
	p := flags.NewParser(&options, 0)
	p.Usage = gnuUsage
	files, err := p.Parse()

	if err != nil {
		// TODO: What does wc do?
		fmt.Fprint(os.Stderr, gnuUsage)
		os.Exit(1)
	}

	if options.Help {
		// TODO What does wc do?
		flag.CommandLine.Output().Write([]byte(gnuUsage))
		os.Exit(0)
	}

	if options.Version {
		// See warning above.
		// TODO What does wc do?
		printVersion()
		os.Exit(0)
	}

	var channels []chan input.WcResult
	if len(files) == 0 {
		result := input.ReadSingleFileInternal("-")
		printSingleFile(options, os.Stdout, result)
		os.Exit(0)
	}

	for _, f := range files {
		c := make(chan input.WcResult)
		channels = append(channels, c)

		go func(f string) {
			c <- input.ReadSingleFileInternal(f)
			close(c)
		}(f)
	}

	var results []input.WcResult
	for _, c := range channels {
		result := <-c
		results = append(results, result)
	}

	printMultipleFiles(options, results, os.Stdout)
}

func printVersion() {
	divider := "================================================================================"

	flag.CommandLine.Output().Write([]byte(versionPreface))
	flag.CommandLine.Output().Write([]byte(divider + "\n"))
	v, _ := debug.ReadBuildInfo()
	flag.CommandLine.Output().Write([]byte(gnuVersion))
	flag.CommandLine.Output().Write([]byte(divider + "\n"))

	flag.CommandLine.Output().Write([]byte("Golang version by Anthony A. Nassar: " + ccVersion))
	flag.CommandLine.Output().Write([]byte("Built with Go version: " + v.GoVersion + "\n"))
}
