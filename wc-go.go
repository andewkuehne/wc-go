package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"unicode"
)

type counts struct {
	byteCount, wordCount, lineCount int64
}

func countInput(input io.Reader, countBytes, countWords, countLines bool) counts {
	var (
		bytesCount, wordsCount, linesCount int64
	)
	scanner := bufio.NewScanner(input)
	if countLines {
		scanner.Split(bufio.ScanLines)
	}
	if countWords {
		scanner.Split(bufio.ScanWords)
	}
	inWord := false
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		bytesCount += int64(len(lineBytes))
		if countWords {
			for _, b := range lineBytes {
				isSpace := unicode.IsSpace(rune(b))
				if isSpace {
					inWord = false
				} else if !inWord {
					inWord = true
					wordsCount++
				}
			}
		}
		if countLines {
			linesCount++
		}
	}
	return counts{
		byteCount: bytesCount,
		wordCount: wordsCount,
		lineCount: linesCount,
	}
}

func wcFile(filePath string, countBytes, countWords, countLines bool, totalCounts *counts, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()
	fileCounts := countInput(file, countBytes, countWords, countLines)

	if countLines {
		fmt.Printf("%8d", fileCounts.lineCount)
	}
	if countWords {
		fmt.Printf("%8d", fileCounts.wordCount)
	}
	if countBytes {
		fmt.Printf("%8d", fileCounts.byteCount)
	}
	fmt.Printf(" %s\n", filePath)

	if totalCounts != nil {
		totalCounts.lineCount += fileCounts.lineCount
		totalCounts.wordCount += fileCounts.wordCount
		totalCounts.byteCount += fileCounts.byteCount
	}
}

func wc(args []string) {
	// Set up flags
	countBytesPtr := flag.Bool("c", false, "print byte counts")
	countWordsPtr := flag.Bool("w", false, "print word counts")
	countLinesPtr := flag.Bool("l", false, "print line counts")
	flag.Parse()

	// Check which flags are set
	printBytes := *countBytesPtr
	printWords := *countWordsPtr
	printLines := *countLinesPtr

	// Get file paths
	filePaths := flag.Args()

	// Determine total counts
	var totalCounts *counts
	if len(filePaths) > 1 {
		totalCounts = &counts{}
	}

	// Count data for each file in separate goroutine
	var wg sync.WaitGroup
	for _, filePath := range filePaths {
		if filePath == "-" {
			// Read from standard input
			wg.Add(1)
			go func() {
				defer wg.Done()
				counts := countInput(os.Stdin, printBytes, printWords, printLines)
				fmt.Printf("%8d %8d %8d\n", counts.lineCount, counts.wordCount, counts.byteCount)
				if totalCounts != nil {
					totalCounts.lineCount += counts.lineCount
					totalCounts.wordCount += counts.wordCount
					totalCounts.byteCount += counts.byteCount
				}
			}()
		} else {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				fileContents, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
				counts := countInput(bytes.NewReader(fileContents), printBytes, printWords, printLines)
				fmt.Printf("%8d %8d %8d %s\n", counts.lineCount, counts.wordCount, counts.byteCount, filePath)
				if totalCounts != nil {
					totalCounts.lineCount += counts.lineCount
					totalCounts.wordCount += counts.wordCount
					totalCounts.byteCount += counts.byteCount
				}
			}(filePath)
		}
	}
	wg.Wait()

	// Print total counts if there were multiple files or standard input was used
	if totalCounts != nil {
		if printLines {
			fmt.Printf("%8d", totalCounts.lineCount)
		}
		if printWords {
			fmt.Printf("%8d", totalCounts.wordCount)
		}
		if printBytes {
			fmt.Printf("%8d", totalCounts.byteCount)
		}
		fmt.Println(" total")
	} else if len(filePaths) == 1 {
		// If only one file path was provided, print its counts
		counts := countInput(os.Stdin, printBytes, printWords, printLines)
		fmt.Printf("%8d %8d %8d %s\n", counts.lineCount, counts.wordCount, counts.byteCount, filePaths[0])
	}
}

func main() {
	wc(os.Args[1:])
}
