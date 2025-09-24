# Makefile for HTTPRunner

# Variables
BINARY_NAME = http-runner
SOURCE_DIR = ./
BUILD_DIR = ./bin
GO_FILES = $(wildcard $(SOURCE_DIR)/*.go)

# Default target
all: build

# Build the binary
build: $(GO_FILES)
	@echo "Building the binary..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

# Run the binary with arguments
run: build
	@echo "Running the binary..."
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Help
help:
	@echo "Makefile for HTTPRunner"
	@echo "Usage:"
	@echo "  make build    - Build the binary"
	@echo "  make run ARGS='--method GET --url http://example.com --count 10 --verbose' - Build and run the binary with arguments"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make help     - Show this help message"

.PHONY: all build run clean help