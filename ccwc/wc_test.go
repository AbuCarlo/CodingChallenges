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

var validTotalArgValues = []string{
	"always", "auto", "only", "never",
}

func generateValidArgs(t *rapid.T) []string {
	totalsGenerator := rapid.Custom(func(t *rapid.T) string {
		totalValue := rapid.SliceOfN(rapid.SampledFrom(validTotalArgValues), 1, 1).Draw(t, "tv")
		return "--total=" + totalValue[0]
	})
	totalArg := rapid.SliceOfN(totalsGenerator, 0, 1).Draw(t, "t")
	args := rapid.SliceOfDistinct(rapid.SampledFrom(validArguments), func(s string) string { return s }).Draw(t, "o")
	args = append(args, totalArg...)
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
		expected := string(out)

		wcOptions, files, _ := parseOptions(args)
		builder := strings.Builder{}
		ReadFiles(files, wcOptions, &builder)
		actual := builder.String()

		if actual != expected {
			t.Errorf("%v failed.\nExpected:\n%s\nGot:\n%s\n", args, string(out), actual)
		}
	})
}
