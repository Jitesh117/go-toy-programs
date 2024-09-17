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
	printLineNos      bool
	ignoreBlanks      bool
	squeezeBlanks     bool
	displayLineBreaks bool
	displayTabs       bool
)

func init() {
	flag.BoolVar(&printLineNos, "n", false, "Number all output lines")
	flag.BoolVar(&ignoreBlanks, "b", false, "Number non-blank lines")
	flag.BoolVar(&squeezeBlanks, "s", false, "Squeeze multiple blank lines into one")
	flag.BoolVar(&displayLineBreaks, "E", false, "Display $ at the end of each line")
	flag.BoolVar(&displayTabs, "T", false, "Display tab characters as ^I")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		printContents(os.Stdin)
	} else {
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cat: %v\n", err)
				continue
			}
			defer file.Close()
			printContents(file)
		}
	}
}

func printContents(r io.Reader) {
	scanner := bufio.NewScanner(r)
	lineNo := 1
	previousWasBlank := false

	for scanner.Scan() {
		line := scanner.Text()

		if squeezeBlanks && line == "" {
			if previousWasBlank {
				continue
			}
			previousWasBlank = true
		} else {
			previousWasBlank = false
		}

		if displayTabs {
			line = strings.ReplaceAll(line, "\t", "^I")
		}

		if displayLineBreaks {
			line += "$"
		}

		if printLineNos {
			fmt.Printf("%6d  %s\n", lineNo, line)
			lineNo++
		} else if ignoreBlanks {
			if line != "" {
				fmt.Printf("%6d  %s\n", lineNo, line)
				lineNo++
			} else {
				fmt.Println(line)
			}
		} else {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "cat: error reading input: %v\n", err)
	}
}
