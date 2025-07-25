package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)
type Scanner struct {
	fileCount      int64
	dirCount       int64
	errorCount     int64
	skippedCount   int64
	bytesScanned   int64
	startTime      time.Time
	ctx            context.Context
	cancel         context.CancelFunc
	workerCount    int
	progressTicker *time.Ticker
	mu             sync.Mutex
	lastError      string
	currentPath    string
}
type ScanResult struct {
	TotalFiles     int64
	TotalDirs      int64
	TotalErrors    int64
	TotalSkipped   int64
	TotalBytes     int64
	Duration       time.Duration
	FilesPerSecond float64
}
func NewScanner() *Scanner {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scanner{
		startTime:      time.Now(),
		ctx:            ctx,
		cancel:         cancel,
		workerCount:    runtime.GOMAXPROCS(0) * 2,
		progressTicker: time.NewTicker(50 * time.Millisecond),
	}
}
func (s *Scanner) Start(rootPath string) *ScanResult {
	fmt.Printf("Starting file system scan from: %s\n", rootPath)
	fmt.Printf("Using %d worker goroutines\n", s.workerCount)
	fmt.Println("Press Ctrl+C to stop at any time")

	go s.displayProgress()

	pathChan := make(chan string, 1000)
	var wg sync.WaitGroup
	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)
		go s.worker(pathChan, &wg)
	}

	go func() {
		defer close(pathChan)
		s.walkDirectory(rootPath, pathChan)
	}()

	wg.Wait()
	s.progressTicker.Stop()
	duration := time.Since(s.startTime)
	filesPerSecond := float64(atomic.LoadInt64(&s.fileCount)) / duration.Seconds()

	return &ScanResult{
		TotalFiles:     atomic.LoadInt64(&s.fileCount),
		TotalDirs:      atomic.LoadInt64(&s.dirCount),
		TotalErrors:    atomic.LoadInt64(&s.errorCount),
		TotalSkipped:   atomic.LoadInt64(&s.skippedCount),
		TotalBytes:     atomic.LoadInt64(&s.bytesScanned),
		Duration:       duration,
		FilesPerSecond: filesPerSecond,
	}
}
func (s *Scanner) Stop() {
	s.cancel()
}
func (s *Scanner) worker(pathChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case path, ok := <-pathChan:
			if !ok {
				return
			}
			s.ProcessPath(path)
		case <-s.ctx.Done():
			return
		}
	}
}
func (s *Scanner) walkDirectory(root string, pathChan chan<- string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		select {
		case <-s.ctx.Done():
			return filepath.SkipDir
		default:
		}

		if err != nil {
			atomic.AddInt64(&s.errorCount, 1)
			s.setLastError(fmt.Sprintf("Error accessing %s: %v", path, err))
			return nil
		}

		s.setCurrentPath(path)

		select {
		case pathChan <- path:
		case <-s.ctx.Done():
			return filepath.SkipDir
		}

		return nil
	})
}
func (s *Scanner) ProcessPath(path string) {
	info, err := os.Lstat(path)
	if err != nil {
		atomic.AddInt64(&s.errorCount, 1)
		s.setLastError(fmt.Sprintf("Error getting info for %s: %v", path, err))
		return
	}

	if info.IsDir() {
		atomic.AddInt64(&s.dirCount, 1)
	} else {
		atomic.AddInt64(&s.fileCount, 1)
		atomic.AddInt64(&s.bytesScanned, info.Size())
	}

	if s.ShouldSkipPath(path) {
		atomic.AddInt64(&s.skippedCount, 1)
	}
}
func (s *Scanner) ShouldSkipPath(path string) bool {
	skipPaths := []string{
		"/proc", "/sys", "/dev", "/run", "/tmp",
		"/var/run", "/var/lock", "/var/tmp",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath || filepath.HasPrefix(path, skipPath+"/") {
			return true
		}
	}
	return false
}
func (s *Scanner) displayProgress() {
	for {
		select {
		case <-s.progressTicker.C:
			files := atomic.LoadInt64(&s.fileCount)
			dirs := atomic.LoadInt64(&s.dirCount)
			errors := atomic.LoadInt64(&s.errorCount)
			skipped := atomic.LoadInt64(&s.skippedCount)
			bytes := atomic.LoadInt64(&s.bytesScanned)
			elapsed := time.Since(s.startTime)

			currentPath := s.getCurrentPath()
			lastError := s.getLastError()

			fmt.Printf("\r\033[K")
			fmt.Printf("Scanned Files: %d | Dirs: %d | Errors: %d | Skipped: %d | Size: %s | Time: %v",
				files, dirs, errors, skipped, FormatBytes(bytes), elapsed.Truncate(time.Second))

			if len(currentPath) > 0 {
				if len(currentPath) > 80 {
					currentPath = "..." + currentPath[len(currentPath)-77:]
				}
				fmt.Printf("\nCurrent: %s", currentPath)
			}

			if len(lastError) > 0 && errors > 0 {
				if len(lastError) > 80 {
					lastError = lastError[:77] + "..."
				}
				fmt.Printf("\nLast Error: %s", lastError)
			}

			if len(currentPath) > 0 || len(lastError) > 0 {
				lines := 1
				if len(currentPath) > 0 {
					lines++
				}
				if len(lastError) > 0 {
					lines++
				}
				fmt.Printf("\033[%dA", lines-1)
			}

		case <-s.ctx.Done():
			return
		}
	}
}
func (s *Scanner) setCurrentPath(path string) {
	s.mu.Lock()
	s.currentPath = path
	s.mu.Unlock()
}
func (s *Scanner) getCurrentPath() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.currentPath
}
func (s *Scanner) setLastError(err string) {
	s.mu.Lock()
	s.lastError = err
	s.mu.Unlock()
}
func (s *Scanner) getLastError() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastError
}
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
