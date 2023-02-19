package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type counts struct {
	byteCount, wordCount, lineCount int
}

func countInput(input io.Reader, countBytes, countWords, countLines bool) counts {
	var (
		bytesCount, wordsCount, linesCount int
	)
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)
	inWord := false
	for scanner.Scan() {
		bytesCount += len(scanner.Bytes())
		if countWords {
			if !inWord {
				inWord = true
				wordsCount++
			}
		}
	}
	if countLines {
		linesCount = 1
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

	// Count data
	counts := countInput(file, *countBytesPtr, *countWordsPtr, *countLinesPtr)

	return counts, nil
}

func wc(args []string) {
	if len(args) == 0 {
		// No arguments, read from standard input
		counts := countInput(os.Stdin, true, true, true)
		fmt.Fprint(os.Stdout, counts.lineCount, " ", counts.wordCount, " ", counts.byteCount, "\n")
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
		fmt.Fprint(os.Stdout, strings.Join(countStrings, " "), " ", flag.Arg(0), "\n")
	}
}

func main() {
	wc(os.Args[1:])
}
