#!/bin/bash
# Build script for entext-research-tool tools
# Usage: ./tool_build.sh [tool-name] [skills-dir]
#
# Available tools:
#   substack-reader  - Substack CLI tools
#   koyfin           - Koyfin CLI tools
#
# Arguments:
#   tool-name   - Name of the tool to build (required)
#   skills-dir  - AI skills directory for deployment (required)
#                 If not provided, will prompt for it

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TOOLS_DIR="$SCRIPT_DIR/tools"

# Global platform variables (set by detect_platform)
TARGET_OS=""
TARGET_ARCH=""

detect_platform() {
    case "$OSTYPE" in
        darwin*)      TARGET_OS="darwin" ;;
        msys*|win32*) TARGET_OS="windows" ;;
        *)            TARGET_OS="linux" ;;
    esac

    case "$(uname -m)" in
        x86_64)        TARGET_ARCH="amd64" ;;
        aarch64|arm64) TARGET_ARCH="arm64" ;;
        armv7l)        TARGET_ARCH="arm" ;;
        i386|i686)     TARGET_ARCH="386" ;;
        *)             TARGET_ARCH="amd64" ;;
    esac
}

get_binary_name() {
    local tool_name="$1"
    local name

    if [ "$tool_name" = "substack-reader" ]; then
        name="substack"
    else
        name="$tool_name"
    fi

    if [ "$TARGET_OS" = "windows" ]; then
        name="$name.exe"
    fi

    echo "$name"
}

build_tool() {
    local tool_name="$1"
    local tool_dir="$2"
    local skills_dir="$3"

    local binary_name bin_file
    binary_name=$(get_binary_name "$tool_name")
    bin_file="$skills_dir/scripts/$binary_name"

    if ! command -v go &> /dev/null; then
        cat <<EOF
Error: Go is not installed

Please install Go from https://go.dev/dl/
EOF
        exit 1
    fi

    printf "Building from source...\n"
    printf "✓ Go: %s\n" "$(go version)"
    printf "✓ Target OS: %s\n" "$TARGET_OS"
    printf "✓ Target Arch: %s\n" "$TARGET_ARCH"
    printf "✓ Deploying to: %s/scripts/\n" "$skills_dir"
    
    GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" go build -ldflags='-w -s' -o "$bin_file" "$tool_dir/src/"
    chmod +x "$bin_file"
    printf "✓ Built: %s\n" "$bin_file"

    cp "$tool_dir/SKILL.md" "$skills_dir"

    if [ -d "$tool_dir/utils" ]; then
        printf "Copying utilities...\n"
        cp "$tool_dir/utils/"* "$skills_dir/scripts"
        printf "✓ Utilities: %s/scripts/\n" "$skills_dir"
    fi

    cat <<EOF

================================
Setup complete!

Tool deployed to: $skills_dir/scripts/
Session file location: $skills_dir/scripts/session.json

Next step: Authenticate
  $skills_dir/scripts/$binary_name auth

Usage:
  $skills_dir/scripts/$binary_name <command> -h
EOF
}

show_usage() {
    cat <<EOF
Usage: $0 <tool-name> [skills-dir]

Arguments:
  tool-name   - Name of the tool to build
  skills-dir  - AI skills directory for deployment
                If not provided, will prompt for it

Available tools:
  substack-reader
  koyfin

Examples:
  $0 substack /<path-to-ai-tool-config>/skills
  $0 koyfin /<path-to-ai-tool-config>/skills
EOF
}

main() {
    if [ -z "${1:-}" ]; then
        show_usage
        exit 0
    fi

    local tool_name="$1"
    local skills_dir="${2:-}"
    local tool_dir="$TOOLS_DIR/$tool_name"

    if [ ! -d "$tool_dir" ]; then
        printf "Error: Tool '%s' not found\n\n" "$tool_name"
        show_usage
        exit 1
    fi

    if [ ! -f "$tool_dir/setup.sh" ]; then
        printf "Error: No setup.sh found for '%s'\n" "$tool_name"
        exit 1
    fi

    if [ -z "$skills_dir" ]; then
        printf "\nBuilding: %s\n\n" "$tool_name"
        read -p "Enter AI skills directory for deployment: " skills_dir
    fi

    if [ -z "$skills_dir" ]; then
        printf "Error: Skills directory is required\n"
        exit 1
    fi

    skills_dir="$skills_dir/$tool_name"
    mkdir -p "$skills_dir/scripts"

    # Detect platform once at startup
    detect_platform

    build_tool "$tool_name" "$tool_dir" "$skills_dir"
}

main "$@"
