BINARY_NAME=dockerizer
VERSION=1.0.0
BUILD_DIR=build

.PHONY: all build clean test

all: clean build

build:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 cmd/main.go
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 cmd/main.go
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe cmd/main.go
	@echo "Done!"

install: build
	@echo "Installing..."
	@go install ./cmd
	@echo "Done!"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Done!"

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Done!"

# Create release archives
release: build
	@echo "Creating release archives..."
	@cd $(BUILD_DIR) && \
		tar czf $(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz $(BINARY_NAME)-linux-amd64 && \
		tar czf $(BINARY_NAME)-darwin-amd64-$(VERSION).tar.gz $(BINARY_NAME)-darwin-amd64 && \
		zip $(BINARY_NAME)-windows-amd64-$(VERSION).zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Done!" 