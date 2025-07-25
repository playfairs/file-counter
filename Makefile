# File Counter Makefile

BINARY_NAME=file-counter
GO_FILES=$(wildcard *.go)
BUILD_DIR=build

all: build

build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .
	@echo "Build complete!"

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete!"

run: build
	@echo "Starting file system scan..."
	./$(BINARY_NAME)
run-sudo: build
	@echo "Starting file system scan with sudo..."
	sudo ./$(BINARY_NAME)

run-demo: build-demo
	@echo "Running demo version..."
	./$(DEMO_BINARY)
build-demo:
	@echo "Building demo application..."
	go build -o $(DEMO_BINARY) ./cmd/demo

build-linux: $(BUILD_DIR)
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(DEMO_BINARY)-linux-amd64 ./cmd/demo

build-darwin: $(BUILD_DIR)
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(DEMO_BINARY)-darwin-amd64 ./cmd/demo
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(DEMO_BINARY)-darwin-arm64 ./cmd/demo

build-windows: $(BUILD_DIR)
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(DEMO_BINARY)-windows-amd64.exe ./cmd/demo

build-all: $(BUILD_DIR) build-linux build-darwin build-windows
	@echo "Cross-compilation complete!"

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

fmt:
	@echo "Formatting code..."
	go fmt ./...
test:
	@echo "Running tests..."
	go test -v ./...

vet:
	@echo "Vetting code..."
	go vet ./...
lint:
	@echo "Running linter..."
	golangci-lint run

dev: fmt vet build
release: clean build-all
	@echo "Release builds created in $(BUILD_DIR)/"

help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  clean       - Remove build artifacts"
	@echo "  run         - Build and run the application"
	@echo "  run-sudo    - Build and run with sudo privileges"
	@echo "  build-linux - Cross-compile for Linux"
	@echo "  build-darwin- Cross-compile for macOS"
	@echo "  build-windows- Cross-compile for Windows"
	@echo "  build-all   - Cross-compile for all platforms"
	@echo "  build-demo  - Build demo application"
	@echo "  run-demo    - Build and run demo application"
	@echo "  deps        - Install Go dependencies"
	@echo "  fmt         - Format Go code"
	@echo "  test        - Run tests"
	@echo "  vet         - Run go vet"
	@echo "  lint        - Run golangci-lint"
	@echo "  dev         - Run development workflow (fmt, vet, build)"
	@echo "  release     - Create release builds for all platforms"
	@echo "  help        - Show this help message"

.PHONY: all build clean run run-sudo run-demo build-demo build-linux build-darwin build-windows build-all deps fmt test vet lint dev release help
