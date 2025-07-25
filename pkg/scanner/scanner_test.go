package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewScanner(t *testing.T) {
	scanner := NewScanner()
	if scanner == nil {
		t.Fatal("NewScanner() returned nil")
	}

	if scanner.workerCount <= 0 {
		t.Error("Worker count should be positive")
	}

	if scanner.startTime.IsZero() {
		t.Error("Start time should be set")
	}
}

func TestScannerStop(t *testing.T) {
	scanner := NewScanner()

	// Should not panic
	scanner.Stop()

	// Context should be cancelled
	select {
	case <-scanner.ctx.Done():
		// Good, context was cancelled
	case <-time.After(100 * time.Millisecond):
		t.Error("Context was not cancelled after Stop()")
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, test := range tests {
		result := FormatBytes(test.input)
		if result != test.expected {
			t.Errorf("FormatBytes(%d) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestShouldSkipPath(t *testing.T) {
	s := NewScanner()

	tests := []struct {
		path     string
		expected bool
	}{
		{"/proc", true},
		{"/proc/cpuinfo", true},
		{"/sys", true},
		{"/sys/class", true},
		{"/dev", true},
		{"/dev/null", true},
		{"/run", true},
		{"/tmp", true},
		{"/var/run", true},
		{"/var/lock", true},
		{"/var/tmp", true},
		{"/home", false},
		{"/usr", false},
		{"/etc", false},
		{"/var/log", false},
	}

	for _, test := range tests {
		result := s.ShouldSkipPath(test.path)
		if result != test.expected {
			t.Errorf("ShouldSkipPath(%s) = %v, expected %v", test.path, result, test.expected)
		}
	}
}

func TestScannerWithTempDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "file-counter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := []string{
		"file1.txt",
		"file2.txt",
		"subdir/file3.txt",
		"subdir/file4.txt",
		"subdir/nested/file5.txt",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	s := NewScanner()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.ctx = ctx

	result := s.Start(tmpDir)

	if result == nil {
		t.Fatal("Scanner returned nil result")
	}

	if result.TotalFiles != 5 {
		t.Errorf("Expected 5 files, got %d", result.TotalFiles)
	}
	if result.TotalDirs < 3 {
		t.Errorf("Expected at least 3 directories, got %d", result.TotalDirs)
	}

	if result.TotalBytes == 0 {
		t.Error("Expected some bytes to be counted")
	}
	if result.Duration <= 0 {
		t.Error("Duration should be positive")
	}
}

func TestScanResultCalculations(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := "Hello, World!"
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	s := NewScanner()
	s.ProcessPath(tmpFile.Name())

	if s.fileCount != 1 {
		t.Errorf("Expected file count 1, got %d", s.fileCount)
	}
	if s.bytesScanned != int64(len(testContent)) {
		t.Errorf("Expected %d bytes, got %d", len(testContent), s.bytesScanned)
	}
}

func BenchmarkFormatBytes(b *testing.B) {
	sizes := []int64{1024, 1048576, 1073741824, 1099511627776}

	for i := 0; i < b.N; i++ {
		for _, size := range sizes {
			FormatBytes(size)
		}
	}
}
