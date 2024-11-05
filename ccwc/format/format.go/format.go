package format

import (
	"fmt"
	"io"
	"os"
	"strings"

	"anthonyabunassar.com/coding-challenges/ccwc/input"
	"anthonyabunassar.com/coding-challenges/ccwc/options"
	"golang.org/x/exp/constraints"
)

func ReportSingleFile(result input.WcResult, options options.WcOptions, w io.Writer) {
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "ccwc: %s: %e\n", result.FileName, result.Error)
		// TODO Print totals nonetheless?
	}
	// TODO Handle --total==always and --total==only
	if options.IsDefault() {
		width := max(countDecimalChars(result.Lines), countDecimalChars(result.Words), countDecimalChars(result.Bytes))
		// The long-standing default for wc: "newline, word, and byte counts"
		fmt.Fprintf(w, "%*d %*d %*d", width, result.Lines, width, result.Words, width, result.Bytes)
		// TODO This is not really right: "-" can appear on the command line.
		if result.FileName == "-" {
			fmt.Fprintln(w)
		} else {
			fmt.Fprintln(w, "", result.FileName)
		}
	} else {
		var values []int64

		// From the wc help: "newline, word, character, byte, maximum line length"
		if options.Lines {
			values = append(values, result.Lines)
		}
		if options.Words {
			values = append(values, result.Words)
		}
		if options.Chars {
			values = append(values, result.Chars)
		}
		if options.Bytes {
			values = append(values, result.Bytes)
		}
		if options.Width {
			values = append(values, result.Width)
		}
		var largestValue int64
		for _, v := range values {
			largestValue = max(largestValue, v)
		}
		width := countDecimalChars(largestValue)
		for _, v := range values {
			fmt.Fprintf(w, "%*d ", width, v)
		}
		fmt.Fprintln(w, result.FileName)
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

func adjustIntegerOutputWidth(w int) int {
	if w < 6 {
		// This is the default in wc.
		return 7
	}
	return w + 1
}

func ReportMultipleFiles(results []input.WcResult, options options.WcOptions, w io.Writer) {
	if len(results) == 1 {
		ReportSingleFile(results[0], options, w)
		return
	}
	// According to https://man7.org/linux/man-pages/man1/wc.1p.html
	// "When any option is specified, wc shall report only the
	// information requested by the specified options."
	var totalWords, totalBytes, totalChars, totalLines int64
	var longestLine int64
	// The minimum width for a number is 7 if we're going to print totals.
	// See https://github.com/coreutils/coreutils/blob/0c9d372c96f2c7ce8c259c5563a48d1816fe611d/src/wc.c#L702
	for _, result := range results {
		totalBytes += result.Bytes
		totalChars += result.Chars
		totalLines += result.Lines
		totalWords += result.Words
		longestLine = max(longestLine, result.Width)
	}

	maxLinesLength := adjustIntegerOutputWidth(countDecimalChars(totalLines))
	maxCharLength := adjustIntegerOutputWidth(countDecimalChars(totalChars))
	maxByteLength := adjustIntegerOutputWidth(countDecimalChars(totalBytes))
	maxWordsLength := adjustIntegerOutputWidth(countDecimalChars(totalWords))
	maxWidthLength := adjustIntegerOutputWidth(countDecimalChars(longestLine))

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

		fmt.Fprintln(w, strings.Join(data, " "), "total")
	}
}