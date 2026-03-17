#!/bin/bash
set -e

echo "Substack Reader CLI Tools Setup"
echo "================================"
echo ""

# Platform-specific paths
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    CONFIG_DIR="$HOME/Library/Application Support/substack-reader"
    BINARY_DIR="$HOME/bin"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    # Windows (Git Bash/WSL)
    CONFIG_DIR="$APPDATA/substack-reader"
    BINARY_DIR="$LOCALAPPDATA/Programs/substack-reader"
else
    # Linux and others (XDG Base Directory)
    if [ -n "$XDG_CONFIG_HOME" ]; then
        CONFIG_DIR="$XDG_CONFIG_HOME/substack-reader"
    else
        CONFIG_DIR="$HOME/.config/substack-reader"
    fi
    BINARY_DIR="$HOME/.local/bin"
fi

BIN_FILE="$BINARY_DIR/substack"

# Create directories
mkdir -p "$CONFIG_DIR"
mkdir -p "$BINARY_DIR"

# Check for pre-built binary first
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Determine target OS and binary name
if [[ "$OSTYPE" == "darwin"* ]]; then
    TARGET_OS="darwin"
    BIN_NAME="substack"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    TARGET_OS="windows"
    BIN_NAME="substack.exe"
else
    TARGET_OS="linux"
    BIN_NAME="substack"
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

# Copy SKILL.md for AI assistants (Gemini, Qwen, Claude Code)
if [ -f "$SCRIPT_DIR/SKILL.md" ]; then
    echo ""
    echo "Copying SKILL.md for AI assistants..."
    mkdir -p "$CONFIG_DIR"
    cp "$SCRIPT_DIR/SKILL.md" "$CONFIG_DIR/"
    echo "✓ SKILL.md: $CONFIG_DIR/SKILL.md"
fi

# Setup PATH
echo ""
if [[ ":$PATH:" != *":$BINARY_DIR:"* ]]; then
    echo "Adding $BINARY_DIR to PATH..."

    # Detect shell and platform
    SHELL_RC=""
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS - common shells
        if [[ "$SHELL" == *"zsh"* ]]; then
            SHELL_RC="$HOME/.zprofile"
        elif [[ "$SHELL" == *"bash"* ]]; then
            SHELL_RC="$HOME/.bash_profile"
        fi
    elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        # Windows - skip shell profile modification
        echo ""
        echo "On Windows, add to PATH manually:"
        echo "  System Properties > Environment Variables > Path > New"
        echo "  Add: $BINARY_DIR"
        SHELL_RC=""
    else
        # Linux
        if [[ "$SHELL" == *"zsh"* ]]; then
            SHELL_RC="$HOME/.zshrc"
        elif [[ "$SHELL" == *"bash"* ]]; then
            SHELL_RC="$HOME/.bashrc"
        fi
    fi

    if [ -n "$SHELL_RC" ] && ! grep -q "$BINARY_DIR" "$SHELL_RC" 2>/dev/null; then
        echo "" >> "$SHELL_RC"
        echo 'export PATH="$PATH:'"$BINARY_DIR"'"' >> "$SHELL_RC"
        echo "✓ Added to $SHELL_RC"
        echo ""
        echo "Run 'source $SHELL_RC' or restart terminal to use 'substack' command"
    fi
fi

# Done
echo ""
echo "================================"
echo "Setup complete!"
echo ""
echo "Next step: Authenticate with Substack"
echo "  substack auth"
echo ""
echo "Usage:"
if [[ ":$PATH:" == *":$BINARY_DIR:"* ]]; then
    echo "  substack <command> -h"
    echo "  substack inbox"
    echo "  substack search -query \"AI\""
else
    echo "  $BIN_FILE <command> -h"
    echo "  $BIN_FILE inbox"
    echo "  $BIN_FILE search -query \"AI\""
fi
