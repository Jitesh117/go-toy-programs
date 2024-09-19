package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	reverseFlag := flag.Bool("r", false, "reverse the result of comparisons")
	numericFlag := flag.Bool("n", false, "compare according to string numerical value")
	uniqueFlag := flag.Bool("u", false, "output only the first of an equal run")
	outputFlag := flag.String("o", "", "write result to output file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE]...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
		return
	}

	if *numericFlag {
		sort.Slice(lines, func(i, j int) bool {
			ni, err1 := strconv.ParseFloat(lines[i], 64)
			nj, err2 := strconv.ParseFloat(lines[j], 64)
			if err1 != nil || err2 != nil {
				return lines[i] < lines[j]
			}
			if *reverseFlag {
				return ni > nj
			}
			return ni < nj
		})
	} else {
		sort.SliceStable(lines, func(i, j int) bool {
			if *reverseFlag {
				return lines[i] > lines[j]
			}
			return lines[i] < lines[j]
		})
	}

	if *uniqueFlag {
		lines = unique(lines)
	}

	if *outputFlag != "" {
		err := writeToFile(*outputFlag, lines)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error writing to file:", err)
		}
	} else {
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}

func unique(lines []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, line := range lines {
		if _, ok := seen[line]; !ok {
			seen[line] = struct{}{}
			result = append(result, line)
		}
	}
	return result
}

func writeToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
