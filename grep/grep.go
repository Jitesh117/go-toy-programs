package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type Options struct {
	invertMatch       bool
	recursive         bool
	ignoreCase        bool
	lineNumber        bool
	count             bool
	onlyMatching      bool
	quiet             bool
	wordRegexp        bool
	beforeContext     int
	afterContext      int
	contextSeparator  string
	maxCount          int
	filesWithMatches  bool
	filesWithoutMatch bool
	binaryFilesTreat  string
}

func grep(pattern string, reader io.Reader, opts Options) error {
	var re *regexp.Regexp
	var err error

	if opts.ignoreCase {
		re, err = regexp.Compile("(?i)" + pattern)
	} else {
		re, err = regexp.Compile(pattern)
	}
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	lineNum := 0
	matchCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		match := re.MatchString(line)

		if opts.wordRegexp {
			match = re.MatchString("\\b" + line + "\\b")
		}

		if (opts.invertMatch && !match) || (!opts.invertMatch && match) {
			if !opts.quiet && !opts.count && !opts.filesWithMatches && !opts.filesWithoutMatch {
				if opts.lineNumber {
					fmt.Printf("%d:", lineNum)
				}
				if opts.onlyMatching {
					matches := re.FindAllString(line, -1)
					for _, m := range matches {
						fmt.Println(m)
					}
				} else {
					fmt.Println(line)
				}
			}
			matchCount++
			if opts.maxCount > 0 && matchCount >= opts.maxCount {
				break
			}
		}
	}

	if opts.count {
		fmt.Println(matchCount)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processFile(pattern string, filename string, opts Options) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return false
	}
	defer file.Close()

	if opts.binaryFilesTreat != "" {
		// Implement binary file detection and handling here
	}

	hasMatch := false
	err = grep(pattern, file, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	} else {
		hasMatch = true
	}

	if opts.filesWithMatches && hasMatch {
		fmt.Println(filename)
	} else if opts.filesWithoutMatch && !hasMatch {
		fmt.Println(filename)
	}

	return hasMatch
}

func processDirectory(pattern string, dirname string, opts Options) {
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			processFile(pattern, path, opts)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func main() {
	opts := Options{}

	flag.BoolVar(&opts.invertMatch, "v", false, "Invert match")
	flag.BoolVar(&opts.recursive, "r", false, "Recursive search")
	flag.BoolVar(&opts.ignoreCase, "i", false, "Ignore case")
	flag.BoolVar(&opts.lineNumber, "n", false, "Print line number with output lines")
	flag.BoolVar(&opts.count, "c", false, "Print only a count of selected lines")
	flag.BoolVar(&opts.onlyMatching, "o", false, "Print only the matched parts of a matching line")
	flag.BoolVar(&opts.quiet, "q", false, "Quiet mode: suppress normal output")
	flag.BoolVar(&opts.wordRegexp, "w", false, "Match only whole words")
	flag.IntVar(&opts.beforeContext, "B", 0, "Print NUM lines of leading context")
	flag.IntVar(&opts.afterContext, "A", 0, "Print NUM lines of trailing context")
	flag.StringVar(
		&opts.contextSeparator,
		"context-separator",
		"--",
		"Set the context separator string",
	)
	flag.IntVar(&opts.maxCount, "m", 0, "Stop after NUM matches")
	flag.BoolVar(
		&opts.filesWithMatches,
		"l",
		false,
		"Print only names of FILEs with selected lines",
	)
	flag.BoolVar(
		&opts.filesWithoutMatch,
		"L",
		false,
		"Print only names of FILEs with no selected lines",
	)
	flag.StringVar(
		&opts.binaryFilesTreat,
		"binary-files",
		"",
		"Type of binary file contents (binary, text, without-match)",
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] PATTERN [FILE...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	pattern := flag.Arg(0)

	if flag.NArg() == 1 {
		err := grep(pattern, os.Stdin, opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	} else {
		for _, filename := range flag.Args()[1:] {
			fileInfo, err := os.Stat(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				continue
			}
			if fileInfo.IsDir() && opts.recursive {
				processDirectory(pattern, filename, opts)
			} else if fileInfo.IsDir() && !opts.recursive {
				fmt.Fprintf(os.Stderr, "Error: %s is a directory, use -r to search recursively\n", filename)
			} else {
				processFile(pattern, filename, opts)
			}
		}
	}
}
