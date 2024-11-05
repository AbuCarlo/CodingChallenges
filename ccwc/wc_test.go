package main

import (
	"os/exec"
	"strings"
	"testing"

	"anthonyabunassar.com/coding-challenges/ccwc/options"
	"github.com/jessevdk/go-flags"
	"pgregory.net/rapid"
)

func parseOptions(line []string) (options.WcOptions, []string, error) {
	var options options.WcOptions
	p := flags.NewParser(&options, 0)
	files, err := p.ParseArgs(line)
	return options, files, err
}

var validArguments = []string{
	"-c", "-m", "-l", "-L", "-w",
}

func generateValidArgs(t *rapid.T) []string {
	args := rapid.SliceOfDistinct(rapid.SampledFrom(validArguments), func(s string) string { return s }).Draw(t, "o")
	return args
}

func TestConsoleOutput(t *testing.T) {
	inputs := []string{
		"./test-files/hello.txt",
	}

	rapid.Check(t, func(t *rapid.T) {
		args := generateValidArgs(t)
		args = append(args, inputs...)
		out, err := exec.Command("wc", args...).Output()
		if err != nil {
			// Something is wrong with the test, or wc is not on $PATH.
			stderr := err.(*exec.ExitError).Stderr
			t.Fatal(string(stderr))
		}
        // GnuUtils on Windows output \r\n, but Go does not!
		expected := strings.ReplaceAll(string(out), "\r\n", "\n")

		wcOptions, files, _ := parseOptions(args)
		builder := strings.Builder{}
		ReadFiles(files, wcOptions, &builder)
		actual := builder.String()
		if actual != expected {
			t.Errorf("%v failed.\nExpected:\n%s\nGot:\n%s\n", args, string(out), actual)
		}
	})
}
