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
	Words    int64
	Lines    int64
	Width    int64
	Error    error
}

/*
	From the documentation:

	"Unless the environment variable POSIXLY_CORRECT is set, GNU wc treats the following Unicode characters
	as white space even if the current locale does not: U+00A0 NO-BREAK SPACE, U+2007 FIGURE SPACE,
	U+202F NARROW NO-BREAK SPACE, and U+2060 WORD JOINER."
*/

var posixlyCorrect bool

func init() {
	_, posixlyCorrect = os.LookupEnv("POSIXLY_CORRECT")
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
	precedingWhitespace := true
	currentWidth := int64(0)

	for {

		r, size, err := r.ReadRune()
		if err != nil {
			if err != io.EOF {
				// TODO What is our plan here?
				result.Error = err
			} else {
				// wc does not count EOF as the end of a word or a line,
				// but it might affect the maximum line width.
				result.Width = max(result.Width, currentWidth)
			}
			break
		}
		currentWidth += 1
		result.Bytes += int64(size)
		/*
			GNU wc does not count encoding errors as characters.

			See also the documentation for bufio.Reader.ReadRune().
		*/
		if r != '\ufffd' {
			result.Chars += 1
		}

		if r == '\n' {
			/*
				GNU wc counts \n characters. A file with no \n will have 0 "lines," per the documentation.
			*/
			result.Width = max(result.Width, currentWidth)
			currentWidth = 0
			result.Lines += 1
		}
		/*
			See the initialization of posixlyCorrect, above.
		*/
		if unicode.IsSpace(r) || posixlyCorrect && (r == '\u2007' || r == '\u202f' || r == '\u2060') {
			precedingWhitespace = true
		} else {
			if precedingWhitespace {
				result.Words += 1
			}
			precedingWhitespace = false
		}
	}
	return result
}

func ReadSingleFileInternal(f string) WcResult {
	var result WcResult
	result.FileName = f
	inputFile, err := os.Open(f)
	if err != nil {
		// wc does not stop on an error.
		result.Error = err
		return result
	}

	var r *bufio.Reader
	if f == "-" || f == "" {
		r = bufio.NewReader(os.Stdin)
	} else {
		// The default buffer size is 4K. Nonetheless, wc uses the much larger
		// buffer size frequently recommended for most architectures. See
		// https://github.com/coreutils/coreutils/blob/master/src/ioblksize.h
		r = bufio.NewReaderSize(inputFile, 262144)
	}

	result = ReadSingleReader(r)
	if f != "" {
		result.FileName = f
	}

	if err := inputFile.Close(); err != nil {
		// At this point the error doesn't matter.
		// Let's just output it later.
		result.Error = err
	}

	return result
}
