#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Tool-specific configuration
TOOL_NAME="substack-reader"
BINARY_NAME="substack"

# Binary and session file directory: <tool_name>/scripts/
BINARY_DIR="$SCRIPT_DIR/scripts"
BIN_FILE="$BINARY_DIR/$BINARY_NAME"

echo "$TOOL_NAME CLI Tools Setup"
echo "================================"
echo ""

# Create directories
mkdir -p "$BINARY_DIR"

# Determine target OS and binary name
if [[ "$OSTYPE" == "darwin"* ]]; then
    TARGET_OS="darwin"
    BIN_NAME="$BINARY_NAME"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    TARGET_OS="windows"
    BIN_NAME="$BINARY_NAME.exe"
else
    TARGET_OS="linux"
    BIN_NAME="$BINARY_NAME"
fi

# Check for pre-built binary (platform-specific)
if [ -f "$PROJECT_ROOT/bin/$TARGET_OS/$BIN_NAME" ]; then
    echo "✓ Found pre-built binary for $TARGET_OS"
    cp "$PROJECT_ROOT/bin/$TARGET_OS/$BIN_NAME" "$BIN_FILE"
    chmod +x "$BIN_FILE"
    echo "✓ Installed: $BIN_FILE"
elif command -v go &> /dev/null; then
    echo "Building from source..."
    echo "✓ Go: $(go version)"
    echo "✓ Target OS: $TARGET_OS"
    GOOS="$TARGET_OS" go build -o "$BIN_FILE" "$SCRIPT_DIR/src/"
    chmod +x "$BIN_FILE"
    echo "✓ Built: $BIN_FILE"
else
    echo "Error: No pre-built binary found and Go is not installed"
    echo ""
    echo "Options:"
    echo "  1. Download pre-built binary from releases"
    echo "  2. Install Go from https://go.dev/dl/"
    exit 1
fi

# Done
echo ""
echo "================================"
echo "Setup complete!"
echo ""
echo "Binary and session file location:"
echo "  $BINARY_DIR"
echo ""
echo "Next step: Authenticate"
echo "  $BIN_FILE auth"
echo ""
echo "Usage:"
echo "  $BIN_FILE <command> -h"
