package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

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

	// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-a-file-using-go
	count := int64(0)

	inputFile, err := os.Open(flag.Arg(0))
    if err != nil {
        panic(err)
    }

	defer func() {
        if err := inputFile.Close(); err != nil {
            panic(err)
        }
    }()

	r := bufio.NewReader(inputFile)

	// TODO: https://www.reddit.com/r/golang/comments/i1cro6/on_choosing_a_buffer_size/
	buf := make([]byte, 65536)
    for {
        n, err := r.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }

        if n == 0 {
            break
        }

        count += int64(n)
    }

	fmt.Printf("%6d %s\n", count, flag.Arg(0))
}
