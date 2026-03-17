#!/bin/bash
set -e

echo "Koyfin CLI Tools Setup"
echo "======================"
echo ""

# Platform-specific paths
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    CONFIG_DIR="$HOME/Library/Application Support/koyfin"
    BINARY_DIR="$HOME/bin"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    # Windows (Git Bash/WSL)
    CONFIG_DIR="$APPDATA/koyfin"
    BINARY_DIR="$LOCALAPPDATA/Programs/koyfin"
else
    # Linux and others (XDG Base Directory)
    if [ -n "$XDG_CONFIG_HOME" ]; then
        CONFIG_DIR="$XDG_CONFIG_HOME/koyfin"
    else
        CONFIG_DIR="$HOME/.config/koyfin"
    fi
    BINARY_DIR="$HOME/.local/bin"
fi

SESSION_FILE="$CONFIG_DIR/session.json"
BIN_FILE="$BINARY_DIR/koyfin"
UTILS_DIR="$BINARY_DIR/koyfin-utils"

# Create directories
mkdir -p "$CONFIG_DIR"
mkdir -p "$BINARY_DIR"

# Check for pre-built binary first
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

if [ -f "$PROJECT_ROOT/bin/koyfin" ]; then
    echo "✓ Found pre-built binary"
    cp "$PROJECT_ROOT/bin/koyfin" "$BIN_FILE"
    chmod +x "$BIN_FILE"
    echo "✓ Installed: $BIN_FILE"
elif command -v go &> /dev/null; then
    echo "Building from source..."
    echo "✓ Go: $(go version)"
    go build -o "$BIN_FILE" "$SCRIPT_DIR/src/"
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

# Create session file
if [ -f "$SESSION_FILE" ]; then
    echo "✓ Session exists: $SESSION_FILE"
else
    echo "Enter Koyfin credentials:"
    echo "(Session stored at: $SESSION_FILE)"
    echo ""
    read -p "Email: " EMAIL
    read -s -p "Password: " PASS
    echo ""
    
    if [ -z "$EMAIL" ] || [ -z "$PASS" ]; then
        echo "Error: Email and password are required"
        exit 1
    fi
    
    echo "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" > "$SESSION_FILE"
    chmod 600 "$SESSION_FILE"
    echo "✓ Session created"
fi

# Copy Python utilities
if [ -d "$SCRIPT_DIR/utils" ]; then
    echo ""
    echo "Copying Python utilities..."
    mkdir -p "$UTILS_DIR"
    cp "$SCRIPT_DIR/utils/"*.py "$UTILS_DIR/"
    cp "$SCRIPT_DIR/utils/requirements.txt" "$UTILS_DIR/" 2>/dev/null || true
    echo "✓ Python utilities: $UTILS_DIR"
fi

# Detect Python command
if command -v python3 &> /dev/null; then
    PYTHON_CMD="python3"
elif command -v python &> /dev/null; then
    PYTHON_CMD="python"
elif command -v py &> /dev/null; then
    PYTHON_CMD="py"
elif command -v py3 &> /dev/null; then
    PYTHON_CMD="py3"
else
    PYTHON_CMD=""
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
        echo "Run 'source $SHELL_RC' or restart terminal to use 'koyfin' command"
    else
        echo "Add to PATH manually:"
        echo "  export PATH=\"\$PATH:$BINARY_DIR\""
    fi
fi

# Done
echo ""
echo "======================"
echo "Setup complete!"
echo ""
echo "Usage:"
if [[ ":$PATH:" == *":$BINARY_DIR:"* ]]; then
    echo "  koyfin <command> -h"
    echo "  koyfin search -q Apple"
else
    echo "  $BIN_FILE <command> -h"
    echo "  $BIN_FILE search -q Apple"
fi

if [ -n "$PYTHON_CMD" ] && [ -d "$UTILS_DIR" ]; then
    echo ""
    echo "Python utilities:"
    echo "  $PYTHON_CMD $UTILS_DIR/excel_export.py -h"
    echo "  koyfin snapshot -kids <list_of_koyfin_ids> | $PYTHON_CMD $UTILS_DIR/excel_export.py -o snapshot.xlsx"
fi
