#!/bin/bash
# Build pre-built binaries for all tools and platforms
# Usage: ./bin/build.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TOOLS_DIR="$PROJECT_ROOT/tools"

# Platforms to build for
PLATFORMS=(
    "darwin amd64"
    "darwin arm64"
    "linux amd64"
    "linux arm64"
    "linux arm"
    "windows amd64"
    "windows 386"
)

# Get binary name for tool
get_binary_name() {
    local tool_name="$1"
    local os="$2"
    local name

    if [ "$tool_name" = "substack-reader" ]; then
        name="substack"
    else
        name="$tool_name"
    fi

    if [ "$os" = "windows" ]; then
        name="$name.exe"
    fi

    echo "$name"
}

build_tool() {
    local tool_name="$1"
    local tool_dir="$2"
    local os="$3"
    local arch="$4"

    local binary_name output_dir output_file
    binary_name=$(get_binary_name "$tool_name" "$os")
    output_dir="$SCRIPT_DIR/$os"
    output_file="$output_dir/$binary_name"

    mkdir -p "$output_dir"

    echo "Building $tool_name for $os/$arch..."
    GOOS="$os" GOARCH="$arch" go build -ldflags='-w -s' -o "$output_file" "$tool_dir/src/"
    chmod +x "$output_file"
    echo "✓ $output_file"
}

main() {
    echo "Building pre-built binaries for all tools..."
    echo ""

    for tool_dir in "$TOOLS_DIR"/*/; do
        if [ ! -d "$tool_dir" ]; then
            continue
        fi

        tool_name=$(basename "$tool_dir")

        # Skip if no src directory
        if [ ! -d "$tool_dir/src" ]; then
            echo "Skipping $tool_name (no src directory)"
            continue
        fi

        echo "================================"
        echo "Tool: $tool_name"
        echo "================================"

        for platform in "${PLATFORMS[@]}"; do
            read -r os arch <<< "$platform"
            build_tool "$tool_name" "$tool_dir" "$os" "$arch"
        done

        echo ""
    done

    echo "================================"
    echo "Build complete!"
    echo "Binaries located in: $SCRIPT_DIR/"
    echo ""
    ls -la "$SCRIPT_DIR"/*/
}

main "$@"
