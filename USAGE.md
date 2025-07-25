# Quick Start Guide

This guide will help you get the File Counter up and running quickly.

## Quick Build & Run

### 1. Build Both Applications
```bash
cd "File Counter"
go run test_build.go

go build -o file-counter .
go build -o file-counter-demo ./cmd/demo
```

### 2. Test with Demo First
```bash
./file-counter-demo
./file-counter-demo ~/Documents
./file-counter-demo /usr/local
```

### 3. Run Full System Scan
```bash
./file-counter
sudo ./file-counter
```

## What You'll See

### Demo Output
```
=== File Counter Demo ===
Scanning directory: /Users/username/Documents
This is a demo version for testing purposes

Starting file system scan from: /Users/username/Documents
Using 8 worker goroutines
Press Ctrl+C to stop at any time

=== DEMO SCAN RESULTS ===
Scanned Path: /Users/username/Documents
Total Files: 15,432
Total Directories: 2,156
Total Errors: 3
Total Size: 4.2 GB
Scan Time: 1.2s
Average File Size: 285.3 KB
Files per Second: 12,860.00
```

### Full System Scan Output
```
=== File Counter - Advanced File System Scanner ===
Scanning entire file system from root /

Starting file system scan from: /
Using 8 worker goroutines
Press Ctrl+C to stop at any time

Scanned Files: 1,245,678 | Dirs: 156,789 | Errors: 23 | Skipped: 45 | Size: 2.3 TB | Time: 5m32s
Current: /System/Library/Frameworks/WebKit.framework/Resources/file.bin

=== FINAL RESULTS ===
Total Files Scanned: 1,245,678
Total Directories: 156,789
Total Errors: 23
Total Skipped: 45
Total Data Size: 2.3 TB
Total Time: 8m45s
Average Speed: 2,375.32 files/second
```

## Build Options

### Using Make (if available)
```bash
make build          
make build-demo     
make run           
make run-demo      
make run-sudo      
```

### Using Build Script
```bash
chmod +x build.sh
./build.sh
```

### Manual Build Commands
```bash
go build -o file-counter .
go build -o file-counter-demo ./cmd/demo

go build -ldflags="-s -w" -o file-counter .
```

## Troubleshooting

### Permission Errors
```bash
sudo ./file-counter
```

### Build Errors
```bash
rm -f file-counter file-counter-demo
go mod tidy
go build -o file-counter .
go build -o file-counter-demo ./cmd/demo
```

### Testing Build
```bash
go run test_build.go
```

## Recommended Workflow

1. **First Time Setup**:
   ```bash
   cd "File Counter"
   go run test_build.go
   ```

2. **Test on Small Directory**:
   ```bash
   ./file-counter-demo ~/Desktop
   ```

3. **Full System Scan**:
   ```bash
   sudo ./file-counter
   ```

4. **Stop Anytime**: Press `Ctrl+C` to stop gracefully

## Project Structure
```
File Counter/
├── main.go                  # Main application (scans from /)
├── pkg/scanner/
│   ├── scanner.go          # Core scanning logic
│   └── scanner_test.go     # Tests
├── cmd/demo/
│   └── main.go             # Demo application
├── build.sh                # Build script
├── test_build.go           # Build verification
├── Makefile               # Build automation
├── go.mod                 # Go module
├── README.md              # Full documentation
└── USAGE.md               # This file
```

## Important Notes

- **Full system scan** can take **hours** on large systems
- **Root privileges** (`sudo`) required for complete access
- **Demo version** is perfect for testing before full scan
- **Press Ctrl+C** anytime to stop and see partial results
- **System directories** like `/proc`, `/sys` are automatically skipped for safety

## Need Help?

- Check `README.md` for detailed documentation
- Run `go run test_build.go` to verify your build
- Use demo version first: `./file-counter-demo .`
- For full system scan: `sudo ./file-counter`

Happy scanning