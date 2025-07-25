//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	fmt.Println("=== File Counter Build Test ===")

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current directory: %s\n", pwd)
	fmt.Printf("Go version: %s\n", runtime.Version())

	requiredFiles := []string{"main.go", "pkg/scanner/scanner.go", "go.mod", "cmd/demo/main.go"}
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("Error: Required file %s not found\n", file)
			os.Exit(1)
		}
		fmt.Printf("✓ Found: %s\n", file)
	}

	fmt.Println("\nRunning go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error running go mod tidy: %v\n%s\n", err, output)
		os.Exit(1)
	}
	fmt.Println("✓ go mod tidy completed")

	fmt.Println("\nBuilding main application...")
	cmd = exec.Command("go", "build", "-o", "file-counter", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error building main application: %v\n%s\n", err, output)
		os.Exit(1)
	}
	fmt.Println("✓ Main application built successfully")

	fmt.Println("\nBuilding demo application...")
	cmd = exec.Command("go", "build", "-o", "file-counter-demo", "./cmd/demo")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error building demo application: %v\n%s\n", err, output)
		os.Exit(1)
	}
	fmt.Println("✓ Demo application built successfully")

	binaries := []string{"file-counter", "file-counter-demo"}
	for _, binary := range binaries {
		if info, err := os.Stat(binary); err == nil {
			fmt.Printf("✓ Binary created: %s (size: %d bytes)\n", binary, info.Size())
		} else {
			fmt.Printf("✗ Binary not found: %s\n", binary)
			os.Exit(1)
		}
	}
	fmt.Println("\nTesting demo application with current directory...")
	cmd = exec.Command("./file-counter-demo", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Demo test failed: %v\n", err)
	}

	fmt.Println("\n Build Test Completed Successfully")
	fmt.Println("\nUsage:")
	fmt.Println("  Full system scan:     ./file-counter")
	fmt.Println("  Demo (current dir):   ./file-counter-demo")
	fmt.Println("  Demo (custom path):   ./file-counter-demo /path/to/directory")
	fmt.Println("\nNote: Use 'sudo ./file-counter' for full system access")
}
