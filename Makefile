.PHONY: all build clean test lint run run-plugin install-hooks

# Binary names
BINARY_NAME=clustermind
PLUGIN_BINARY_NAME=clustermind-k9s-plugin

# Build output directory
BUILD_DIR=bin

all: lint test build

build:
	@echo "Building..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/clustermind
	go build -o $(BUILD_DIR)/$(PLUGIN_BINARY_NAME) ./cmd/clustermind-k9s-plugin

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

run-plugin: build
	@echo "Running $(PLUGIN_BINARY_NAME)..."
	./$(BUILD_DIR)/$(PLUGIN_BINARY_NAME)

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

test:
	@echo "Testing..."
	go test -v ./...

lint:
	@echo "Linting..."
	golangci-lint run

install-hooks:
	@echo "Installing git hooks..."
	chmod +x scripts/pre-commit.sh
	ln -sf ../../scripts/pre-commit.sh .git/hooks/pre-commit
	@echo "Hooks installed!"
