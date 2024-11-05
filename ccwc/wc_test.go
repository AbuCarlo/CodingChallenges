package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"anthonyabunassar.com/coding-challenges/ccwc/options"
	"github.com/jessevdk/go-flags"
)

func parseOptions(line []string) (options.WcOptions, []string, error) {
	var options options.WcOptions
	p := flags.NewParser(&options, 0)
	files, err := p.ParseArgs(line)
    return options, files, err
}

func TestConsoleOutput(t *testing.T) {
	testCases := [][]string{
		{"./test-files/hello.txt"},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			out, err := exec.Command("wc", testCase...).Output()
			if err != nil {
                // You made a mistake.
				stderr := err.(*exec.ExitError).Stderr
				t.Fatal(string(stderr))
			} else {
				fmt.Println(string(out))
			}

            options, files, err := parseOptions(testCase)
            if err != nil {
                t.Fatal(err)
            }

            builder := strings.Builder{}
            ReadFiles(files, options, &builder)
            actual := builder.String()

            if actual != string(out) {
                t.Errorf("%v failed.\nExpected:\n%s\nGot:\n%s\n", testCase, string(out), actual)
            }
		})
	}

	// TODO: Externalize command-line parsing.
}

// TODO Generate PBTs