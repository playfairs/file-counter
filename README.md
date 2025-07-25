# File Counter - Advanced File System Scanner

A high-performance Go application that scans the entire file system starting from the root directory `/` and counts all files while displaying real-time progress.

## Features

- **Real-time Progress Display**: Shows live updates of scanned files, directories, errors, and more
- **Concurrent Processing**: Uses multiple goroutines for efficient scanning
- **Graceful Shutdown**: Handle Ctrl+C interrupts cleanly
- **Comprehensive Statistics**: File count, directory count, total size, scan speed, and error tracking
- **Smart Error Handling**: Continues scanning even when encountering permission errors
- **System Directory Skipping**: Automatically skips problematic system directories like `/proc`, `/sys`, `/dev`

## Requirements

- Go 1.21 or later
- Unix-like operating system (Linux, macOS)
- Root/sudo privileges recommended for complete system access

## Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd "File Counter"
   ```
3. Build the applications:
   ```bash
   go build -o file-counter .
   go build -o file-counter-demo ./cmd/demo
   ```

## Usage

### Full System Scan
```bash
./file-counter              # Basic scan
sudo ./file-counter         # With root privileges (recommended)
```

### Demo Version (for testing)
```bash
./file-counter-demo                    # Scan current directory
./file-counter-demo /path/to/scan     # Scan specific directory
```

Root privileges provide access to all system files and directories that would otherwise be restricted.

## Output Example

```
=== File Counter - Advanced File System Scanner ===
Scanning entire file system from root /
Note: This may take a very long time and require elevated permissions
Use 'sudo' for full system access if needed

Starting file system scan from: /
Using 8 worker goroutines
Press Ctrl+C to stop at any time

Scanned Files: 1,245,678 | Dirs: 156,789 | Errors: 23 | Skipped: 45 | Size: 2.3 TB | Time: 5m32s
Current: /Users/username/Documents/projects/large-file.zip
Last Error: Error accessing /private/var/db/ConfigurationProfiles: permission denied

=== FINAL RESULTS ===
Total Files Scanned: 1,245,678
Total Directories: 156,789
Total Errors: 23
Total Skipped: 45
Total Data Size: 2.3 TB
Total Time: 8m45s
Average Speed: 2,375.32 files/second
Average File Size: 1.9 MB
Items per Second: 2,673.21

Scan completed with 23 errors (permission denied, etc.)
```

## Understanding the Output

- **Scanned Files**: Total number of regular files found
- **Dirs**: Total number of directories processed
- **Errors**: Files/directories that couldn't be accessed (usually permission issues)
- **Skipped**: System directories automatically skipped for safety
- **Size**: Total size of all scanned files
- **Current**: The file/directory currently being processed
- **Last Error**: Most recent error encountered

## Performance Considerations

- **Memory Usage**: The application uses minimal memory as it doesn't store file lists
- **CPU Usage**: Uses multiple goroutines (default: 2x CPU cores) for parallel processing
- **I/O Performance**: Optimized for fast directory traversal
- **Large File Systems**: Can handle millions of files efficiently

## Safety Features

- **System Directory Protection**: Automatically skips dangerous system directories
- **Graceful Interruption**: Ctrl+C stops the scan cleanly and shows partial results
- **Error Resilience**: Continues scanning even when individual files cause errors
- **Permission Handling**: Gracefully handles permission denied errors

## Troubleshooting

### Permission Errors
If you see many permission errors, run with sudo:
```bash
sudo ./file-counter
```

### Testing Before Full Scan
Use the demo version to test on smaller directories first:
```bash
./file-counter-demo ~/Documents
```

### Slow Performance
- Close unnecessary applications to free up I/O resources
- Consider running on an SSD for better performance
- The scan speed depends on your storage device and file system structure

### High Memory Usage
The application is designed to use minimal memory, but scanning very large directories with millions of files might use more system resources.

## Project Structure

```
File Counter/
├── main.go              # Main application (full system scan)
├── scanner.go           # Core scanning logic
├── scanner_test.go      # Unit tests
├── cmd/
│   └── demo/
│       └── main.go      # Demo application
├── build.sh             # Build script
├── Makefile             # Build automation
├── test_build.go        # Build verification
└── README.md            # Documentation
```

## Technical Details

- **Language**: Go 1.21+
- **Concurrency**: Worker pool pattern with configurable goroutines
- **File System API**: Uses Go's `filepath.Walk` and `os.Lstat`
- **Progress Updates**: Real-time updates every 50ms
- **Architecture**: Concurrent producer-consumer pattern

## Contributing

Feel free to submit issues or pull requests to improve the application.

## License

This project is open source. Use it freely for personal or commercial purposes.

## Warning

This application scans the ENTIRE file system starting from root (`/`). On systems with large amounts of data, this can:
- Take several hours to complete
- Generate significant disk I/O
- Require substantial system permissions

Always ensure you have adequate time and system resources before starting a full system scan.# file-counter
