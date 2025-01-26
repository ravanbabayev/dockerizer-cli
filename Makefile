BINARY_NAME=dockerizer
VERSION=1.0.0
BUILD_DIR=build

.PHONY: all build clean test

all: clean build

build:
	@echo "Building for multiple platforms..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	go build -o "$(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe" cmd/main.go
	set GOOS=linux&& go build -o "$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64" cmd/main.go
	set GOOS=darwin&& go build -o "$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64" cmd/main.go
	@echo "Done!"

install: build
	@echo "Installing..."
	go install ./cmd
	@echo "Done!"

clean:
	@echo "Cleaning..."
	@if exist "$(BUILD_DIR)" rd /s /q "$(BUILD_DIR)"
	@echo "Done!"

test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Done!"

# Create release archives
release: build
	@echo "Creating release archives..."
	cd "$(BUILD_DIR)" && tar -czf "$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz" "$(BINARY_NAME)-linux-amd64"
	cd "$(BUILD_DIR)" && tar -czf "$(BINARY_NAME)-darwin-amd64-$(VERSION).tar.gz" "$(BINARY_NAME)-darwin-amd64"
	cd "$(BUILD_DIR)" && powershell Compress-Archive -Path "$(BINARY_NAME)-windows-amd64.exe" -DestinationPath "$(BINARY_NAME)-windows-amd64-$(VERSION).zip"
	@echo "Done!" 