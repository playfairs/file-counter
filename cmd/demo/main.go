package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"file-counter/pkg/scanner"
)

// Demo version that scans a specific directory instead of root
func main() {
	var scanPath string

	// Check if a path was provided as command line argument
	if len(os.Args) > 1 {
		scanPath = os.Args[1]
	} else {
		// Default to current directory for demo
		var err error
		scanPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Verify the path exists
	if _, err := os.Stat(scanPath); os.IsNotExist(err) {
		fmt.Printf("Error: Path '%s' does not exist\n", scanPath)
		os.Exit(1)
	}

	fmt.Println("=== File Counter Demo ===")
	fmt.Printf("Scanning directory: %s\n", scanPath)
	fmt.Println("This is a demo version for testing purposes")
	fmt.Println()

	// Create scanner
	fileScanner := scanner.NewScanner()

	// Start the scan
	startTime := time.Now()
	result := fileScanner.Start(scanPath)
	duration := time.Since(startTime)

	// Display results
	fmt.Printf("\n=== DEMO SCAN RESULTS ===\n")
	fmt.Printf("Scanned Path: %s\n", scanPath)
	fmt.Printf("Total Files: %d\n", result.TotalFiles)
	fmt.Printf("Total Directories: %d\n", result.TotalDirs)
	fmt.Printf("Total Errors: %d\n", result.TotalErrors)
	fmt.Printf("Total Size: %s\n", scanner.FormatBytes(result.TotalBytes))
	fmt.Printf("Scan Time: %v\n", duration.Truncate(time.Millisecond))

	if result.TotalFiles > 0 {
		fmt.Printf("Average File Size: %s\n", scanner.FormatBytes(result.TotalBytes/result.TotalFiles))
		fmt.Printf("Files per Second: %.2f\n", result.FilesPerSecond)
	}

	fmt.Println("\nDemo completed successfully!")
	fmt.Println("To scan the entire system, use the main program: ./file-counter")
}

// Simple demo function to show directory structure
func showDirectoryTree(root string, maxDepth int) {
	fmt.Printf("\nDirectory structure (max depth %d):\n", maxDepth)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors in demo
		}

		// Calculate depth
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		depth := len(filepath.SplitList(rel)) - 1
		if rel == "." {
			depth = 0
		}

		if depth > maxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create indentation
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}

		// Print file/directory with appropriate symbol
		symbol := "📄"
		if info.IsDir() {
			symbol = "📁"
		}

		name := filepath.Base(path)
		if depth == 0 {
			name = path
		}

		fmt.Printf("%s%s %s", indent, symbol, name)

		if !info.IsDir() {
			fmt.Printf(" (%s)", scanner.FormatBytes(info.Size()))
		}

		fmt.Println()

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}
