.PHONY: help build clean test vet run-auth run-inbox run-search run-article release-all release-linux release-darwin release-windows install uninstall uninstall-all setup clean-session deps tidy

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=substack
BINARY_DIR=bin

# Source directories
SRC_DIR=tools/substack-reader/src
INTERNAL_DIR=internal
SETUP_SCRIPT=tools/substack-reader/setup.sh

# Build output paths
BUILD_LINUX=$(BINARY_DIR)/linux/$(BINARY_NAME)
BUILD_DARWIN=$(BINARY_DIR)/darwin/$(BINARY_NAME)
BUILD_WINDOWS=$(BINARY_DIR)/windows/$(BINARY_NAME).exe
BUILD_LOCAL=$(BINARY_DIR)/$(BINARY_NAME)

# Default target
help:
	@echo "Substack CLI - Makefile Commands"
	@echo "================================"
	@echo ""
	@echo "Build:"
	@echo "  make build           - Build for current platform"
	@echo "  make release-all     - Build for all platforms (linux, darwin, windows)"
	@echo "  make release-linux   - Build for Linux"
	@echo "  make release-darwin  - Build for macOS (Intel + Apple Silicon)"
	@echo "  make release-windows - Build for Windows"
	@echo ""
	@echo "Install/Setup:"
	@echo "  make setup           - Run full setup (build + install + auth session)"
	@echo "  make install         - Build and install binary only"
	@echo "  make uninstall       - Remove binary from platform-specific path"
	@echo "  make uninstall-all   - Remove binary + session file + config"
	@echo "  make clean-session   - Remove only the session file"
	@echo ""
	@echo "Test & Lint:"
	@echo "  make test            - Run tests"
	@echo "  make vet             - Run go vet"
	@echo "  make clean           - Clean build artifacts"
	@echo ""
	@echo "Run (requires session):"
	@echo "  make run-auth        - Run auth command"
	@echo "  make run-inbox       - Run inbox command"
	@echo "  make run-search      - Run search command"
	@echo "  make run-article     - Run article command"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps            - Download dependencies"
	@echo "  make tidy            - Tidy go.mod"

# Build for current platform
build: $(BINARY_DIR)
	$(GOBUILD) -o $(BUILD_LOCAL) ./$(SRC_DIR)/
	@echo "✓ Built: $(BUILD_LOCAL)"

# Create bin directory
$(BINARY_DIR):
	mkdir -p $(BINARY_DIR)
	mkdir -p $(BINARY_DIR)/linux
	mkdir -p $(BINARY_DIR)/darwin
	mkdir -p $(BINARY_DIR)/windows

# Build for all platforms
release-all: release-linux release-darwin release-windows
	@echo "✓ Built all platforms"

# Build for Linux
release-linux: $(BINARY_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_LINUX) ./$(SRC_DIR)/
	@echo "✓ Built: $(BUILD_LINUX)"

# Build for macOS (Intel & Apple Silicon)
release-darwin: $(BINARY_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DARWIN)_amd64 ./$(SRC_DIR)/
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DARWIN)_arm64 ./$(SRC_DIR)/
	@echo "✓ Built: $(BUILD_DARWIN)_amd64 and $(BUILD_DARWIN)_arm64"

# Build for Windows
release-windows: $(BINARY_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_WINDOWS) ./$(SRC_DIR)/
	@echo "✓ Built: $(BUILD_WINDOWS)"

# Run tests
test:
	$(GOTEST) -v ./$(INTERNAL_DIR)/...

# Run go vet
vet:
	$(GOVET) ./...

# Clean build artifacts
clean:
	rm -rf $(BINARY_DIR)
	@echo "✓ Cleaned build artifacts"

# Full setup (calls setup.sh)
setup:
	@echo "Running full setup..."
	@bash $(SETUP_SCRIPT)

# Install to platform-specific location (matches setup.sh paths)
install: build
	@echo "Installing to platform-specific location..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		mkdir -p $$HOME/bin && cp $(BUILD_LOCAL) $$HOME/bin/$(BINARY_NAME) && chmod +x $$HOME/bin/$(BINARY_NAME); \
		echo "✓ Installed to $$HOME/bin/$(BINARY_NAME)"; \
	elif [ "$$(uname)" = "Linux" ]; then \
		mkdir -p $$HOME/.local/bin && cp $(BUILD_LOCAL) $$HOME/.local/bin/$(BINARY_NAME) && chmod +x $$HOME/.local/bin/$(BINARY_NAME); \
		echo "✓ Installed to $$HOME/.local/bin/$(BINARY_NAME)"; \
	else \
		echo "Manual installation required for this platform"; \
		echo "Copy $(BUILD_LOCAL) to your preferred location"; \
	fi

# Uninstall binary from platform-specific location
uninstall:
	@echo "Removing binary from platform-specific location..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		if [ -f $$HOME/bin/$(BINARY_NAME) ]; then \
			rm -f $$HOME/bin/$(BINARY_NAME) && echo "✓ Removed from $$HOME/bin/$(BINARY_NAME)"; \
		else \
			echo "Binary not found at $$HOME/bin/$(BINARY_NAME)"; \
		fi; \
	elif [ "$$(uname)" = "Linux" ]; then \
		if [ -f $$HOME/.local/bin/$(BINARY_NAME) ]; then \
			rm -f $$HOME/.local/bin/$(BINARY_NAME) && echo "✓ Removed from $$HOME/.local/bin/$(BINARY_NAME)"; \
		else \
			echo "Binary not found at $$HOME/.local/bin/$(BINARY_NAME)"; \
		fi; \
	else \
		echo "Manual uninstallation required for this platform"; \
		echo "Remove the binary from your installation directory"; \
	fi

# Uninstall everything (binary + session + config)
uninstall-all: uninstall clean-session
	@echo ""
	@echo "Full uninstall complete."
	@echo "To remove SKILL.md from config directory, run:"
	@echo "  rm -rf ~/.config/substack-reader  # Linux"
	@echo "  rm -rf ~/Library/Application\\ Support/substack-reader  # macOS"

# Clean only session file (for security)
clean-session:
	@echo "Removing session file..."
	@if [ -n "$$SUBSTACK_SESSION_FILE" ]; then \
		rm -f "$$SUBSTACK_SESSION_FILE" && echo "✓ Removed: $$SUBSTACK_SESSION_FILE"; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		rm -f $$HOME/Library/Application\\ Support/substack-reader/session.json && echo "✓ Removed session file (macOS)"; \
	elif [ "$$(uname)" = "Linux" ]; then \
		rm -f $$HOME/.config/substack-reader/session.json && echo "✓ Removed session file (Linux)"; \
	else \
		echo "Session file location unknown - remove manually"; \
	fi

# Run auth command
run-auth: build
	./$(BUILD_LOCAL) auth

# Run inbox command
run-inbox: build
	./$(BUILD_LOCAL) inbox

# Run search command
run-search: build
	./$(BUILD_LOCAL) search -query "test"

# Run article command
run-article: build
	./$(BUILD_LOCAL) article -post-id 123456

# Download dependencies
deps:
	$(GOGET) ./...

# Tidy go.mod
tidy:
	$(GOMOD) tidy
