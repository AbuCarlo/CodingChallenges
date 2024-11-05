package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	_ "embed"

	"github.com/jessevdk/go-flags"

	"anthonyabunassar.com/coding-challenges/ccwc/format/format.go"
	"anthonyabunassar.com/coding-challenges/ccwc/input"
	"anthonyabunassar.com/coding-challenges/ccwc/options"
)

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
	var options options.WcOptions
	p := flags.NewParser(&options, 0)
	p.Usage = gnuUsage
	files, err := p.Parse()

	if err != nil {
		fmt.Fprint(os.Stderr, gnuUsage)
		os.Exit(1)
	}

	if options.Help {
		// wc prints usage to standard output unless there's an error:
		// https://github.com/coreutils/coreutils/blob/fbfda4df1a1f73b6e9e5e8bd00f2ddb06c6f219b/src/wc.c#L155
		fmt.Print(gnuUsage)
		os.Exit(0)
	}

	if options.Version {
		printVersion()
		os.Exit(0)
	}

	ReadFiles(files, options, os.Stdout)
}

func ReadFiles(files []string, options options.WcOptions, w io.Writer) {
	var channels []chan input.WcResult
	if len(files) == 0 {
		result := input.ReadSingleFileInternal("")
		format.ReportSingleFile(result, options, w)
		return
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

	format.ReportMultipleFiles(results, options, w)
}

func printVersion() {
	divider := "================================================================================"
	/*
		wc writes version informaton to standard output.

		wc --version > version.txt 2> version-error.txt
	*/

	fmt.Print(versionPreface)
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	v, _ := debug.ReadBuildInfo()
	fmt.Print(gnuVersion)
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
	// This string will already have a line ending if created by redirecting Git's output.
	fmt.Printf("Golang version by Anthony A. Nassar: %s", ccVersion)
	fmt.Printf("Built with Go version: %s\n", v.GoVersion)
	fmt.Println()
}
