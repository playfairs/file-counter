package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"file-counter/pkg/scanner"
)

func main() {
	fmt.Println("=== File Counter - Advanced File System Scanner ===")
	fmt.Println("Scanning entire file system from root /")
	fmt.Println("Note: This may take a very long time and require elevated permissions")
	fmt.Println("Use 'sudo' for full system access if needed")
	fmt.Println()

	fileScanner := scanner.NewScanner()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	resultChan := make(chan *scanner.ScanResult, 1)
	go func() {
		result := fileScanner.Start("/")
		resultChan <- result
	}()

	var result *scanner.ScanResult
	select {
	case <-sigChan:
		fmt.Println("\n\nReceived interrupt signal. Stopping scan...")
		fileScanner.Stop()
		select {
		case result = <-resultChan:
		case <-make(chan struct{}):
		}
		fmt.Println("Scan interrupted by user.")
	case result = <-resultChan:
		fmt.Println("\n\nScan completed!")
	}
	if result != nil {
		fmt.Printf("\n=== FINAL RESULTS ===\n")
		fmt.Printf("Total Files Scanned: %d\n", result.TotalFiles)
		fmt.Printf("Total Directories: %d\n", result.TotalDirs)
		fmt.Printf("Total Errors: %d\n", result.TotalErrors)
		fmt.Printf("Total Skipped: %d\n", result.TotalSkipped)
		fmt.Printf("Total Data Size: %s\n", scanner.FormatBytes(result.TotalBytes))
		fmt.Printf("Total Time: %v\n", result.Duration.Truncate(1))
		fmt.Printf("Average Speed: %.2f files/second\n", result.FilesPerSecond)

		if result.TotalFiles > 0 {
			avgFileSize := float64(result.TotalBytes) / float64(result.TotalFiles)
			fmt.Printf("Average File Size: %s\n", scanner.FormatBytes(int64(avgFileSize)))
		}

		totalItems := result.TotalFiles + result.TotalDirs
		if totalItems > 0 {
			itemsPerSecond := float64(totalItems) / result.Duration.Seconds()
			fmt.Printf("Items per Second: %.2f\n", itemsPerSecond)
		}

		if result.TotalErrors > 0 {
			fmt.Printf("\nScan completed with %d errors (permission denied, etc.)\n", result.TotalErrors)
		} else {
			fmt.Printf("\nScan completed successfully with no errors!\n")
		}
	}

	fmt.Println("\nThank you for using File Counter.")
}
