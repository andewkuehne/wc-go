package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type counts struct {
	byteCount, wordCount, lineCount int
}

func countInput(input []byte, countBytes, countWords, countLines bool) counts {
	var (
		bytesCount, wordsCount, linesCount int
		inWord                             bool
	)
	for _, b := range input {
		bytesCount++
		switch b {
		case '\n':
			linesCount++
			inWord = false
		case ' ', '\t', '\r', '\f', '\v':
			inWord = false
		default:
			if !inWord {
				inWord = true
				wordsCount++
			}
		}
	}
	return counts{
		byteCount: bytesCount,
		wordCount: wordsCount,
		lineCount: linesCount,
	}
}

func countArgs(args []string) (counts, error) {
	// Set up flags
	countBytesPtr := flag.Bool("c", false, "print byte counts")
	countWordsPtr := flag.Bool("w", false, "print word counts")
	countLinesPtr := flag.Bool("l", false, "print line counts")
	flag.Parse()

	// Get file name from arguments
	fileName := flag.Arg(0)

	// Read file
	file, err := os.Open(fileName)
	if err != nil {
		return counts{}, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return counts{}, err
	}

	// Count data
	counts := countInput(data, *countBytesPtr, *countWordsPtr, *countLinesPtr)

	return counts, nil
}

func wc(args []string) {
	if len(args) == 0 {
		// No arguments, read from standard input
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
			os.Exit(1)
		}
		counts := countInput(data, true, true, true)
		fmt.Printf("%d %d %d\n", counts.lineCount, counts.wordCount, counts.byteCount)
	} else {
		// Parse arguments and count data from file
		counts, err := countArgs(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error counting data: %v\n", err)
			os.Exit(1)
		}

		// Print requested counts
		var countStrings []string
		if counts.lineCount > 0 {
			countStrings = append(countStrings, fmt.Sprintf("%d", counts.lineCount))
		}
		if counts.wordCount > 0 {
			countStrings = append(countStrings, fmt.Sprintf("%d", counts.wordCount))
		}
		if counts.byteCount > 0 {
			countStrings = append(countStrings, fmt.Sprintf("%d", counts.byteCount))
		}
		fmt.Println(strings.Join(countStrings, " "))
	}
}

func main() {
	wc(os.Args[1:])
}
