package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	countWords bool
	countLines bool
	countChars bool
	countBytes bool
)

func init() {
	flag.BoolVar(&countWords, "w", false, "count words")
	flag.BoolVar(&countLines, "l", false, "count lines")
	flag.BoolVar(&countChars, "c", false, "count chars")
	flag.BoolVar(&countBytes, "b", false, "count bytes")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if !countWords && !countLines && !countChars && !countBytes {
		countWords, countLines, countChars = true, true, true
	}
	if len(args) == 0 {
		count(os.Stdin, "")
	} else {
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "wc: %v\n", err)
				continue
			}
			count(file, filename)
			file.Close()
		}
	}
}

func count(r io.Reader, filename string) {
	scanner := bufio.NewScanner(r)
	var lines, words, chars, bytes int

	for scanner.Scan() {
		line := scanner.Text()
		lines++
		words += len(strings.Fields(line))
		chars += len(line)
		bytes += len([]byte(line)) + 1
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "wc: %v\n", err)
	}

	output := ""
	if countLines {
		output += fmt.Sprintf("%d\t", lines)
	}
	if countWords {
		output += fmt.Sprintf("%d\t", words)
	}

	if countChars {
		output += fmt.Sprintf("%d\t", chars)
	}
	if countBytes {
		output += fmt.Sprintf("%d\t", bytes)
	}

	if filename != "" {
		output += filename
	}
	fmt.Println(output)
}
