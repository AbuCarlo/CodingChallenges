package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type WcResult struct {
	f string
	bites int64
	chars int64
	// words int
	lines int
	longest int
	err error
}

func readSingleFile(f string, c chan<- WcResult) {
	// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-a-file-using-go
	var result WcResult
	result.f = f

	inputFile, err := os.Open(f)
    if err != nil {
        result.err = err
		c <- result
		return
    }

	// The default buffer size is 4K. Why not 64K?
	// TODO: https://www.reddit.com/r/golang/comments/i1cro6/on_choosing_a_buffer_size/
	r := bufio.NewReaderSize(inputFile, 65536)

	// buf := make([]byte, 65536)
    for {
		s, err := r.ReadString('\n')
        // n, err := r.Read(buf)
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
        result.err = err
	}

	c <- result
}

func main() {
	characters := flag.Bool("c", false, "------")

	flag.Parse()

	if !*characters {
		fmt.Fprintln(os.Stderr, "Only the -c option is accepted.")
		// TODO: What does 
		os.Exit(0)
	}

	if len(flag.Args()) > 1 {
		os.Stderr.WriteString("Only one file is currently accepted.")
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		os.Stderr.WriteString("At least one file must be specified.")
		os.Exit(0)
	}

	var channels []chan WcResult
	for _, f := range flag.Args() {
		c := make(chan WcResult)
		channels = append(channels, c)
		go readSingleFile(f, c)
	}

	for _, c := range channels {
		result := <-c
		fmt.Printf("%6d bytes, %6d chars, %6d lines in %s\n", result.bites, result.chars, result.lines, result.f)
	}
}
