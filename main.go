package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type SearchResult struct {
	FilePath string
	Line     string
}

func main() {
	start := time.Now() // Start timer

	pdfDir := "pdfs"
	targetName := "Jhon Doe"
	numWorkers := 20

	var foundResults []SearchResult

	files := []string{}

	// Step 1: Collect all PDF files
	err := filepath.Walk(pdfDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".pdf") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking PDF directory:", err)
		return
	}

	totalFiles := len(files)
	if totalFiles == 0 {
		fmt.Println("No PDF files found.")
		return
	}
	fmt.Printf("Starting search for \"%s\" in %d PDF files...\n", targetName, totalFiles)

	// Channels
	fileChan := make(chan string)
	resultChan := make(chan SearchResult)
	var wg sync.WaitGroup

	// Progress counter
	var processedCount int
	var countMutex sync.Mutex

	// Step 2: Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				found, line := SearchPdfForName(file, targetName)

				// Update progress
				countMutex.Lock()
				processedCount++
				currentCount := processedCount
				countMutex.Unlock()

				if currentCount%10 == 0 || currentCount == totalFiles {
					fmt.Printf("Processed %d / %d files...\n", currentCount, totalFiles)
				}

				if found {
					resultChan <- SearchResult{FilePath: file, Line: line}
				}
			}
		}()
	}

	// Step 3: Feed files to workers
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// Step 4: Close results when all done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Step 5: Collect results
	matchCount := 0
	for result := range resultChan {
		matchCount++
		fmt.Printf("\n✅ Match #%d in %s:\n%s\n", matchCount, result.FilePath, result.Line)
		foundResults = append(foundResults, result)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nDone! Processed %d files, found %d matches.\n", totalFiles, matchCount)
	fmt.Printf("Results found: %s \n", foundResults)
	fmt.Printf("⏱️  Total time: %s\n", elapsed.Round(time.Millisecond))
}
