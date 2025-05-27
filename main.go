package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SearchResult struct {
	FilePath string
	Line     string
}

func main() {
	pdfDir := "pdfs"
	targetName := "John Doe"
	numWorkers := 16 // Can be adjusted based on CPU

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

	// Progress counter with mutex
	var processedCount int
	var countMutex sync.Mutex

	// Step 2: Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				found, line := SearchPdfForName(file, targetName)

				// Update progress safely
				countMutex.Lock()
				processedCount++
				currentCount := processedCount
				countMutex.Unlock()

				// Print progress every 10 files
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

	// Step 4: Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Step 5: Collect and print results
	matchCount := 0
	for result := range resultChan {
		matchCount++
		fmt.Printf("\n✅ Match #%d in %s:\n%s\n", matchCount, result.FilePath, result.Line)
	}

	fmt.Printf("\nDone! Processed %d files, found %d matches.\n", totalFiles, matchCount)
}
