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

show_usage() {
    echo "Usage: $0 <tool-name> [skills-dir]"
    echo ""
    echo "Arguments:"
    echo "  tool-name   - Name of the tool to build"
    echo "  skills-dir  - AI skills directory for deployment"
    echo "                If not provided, will prompt for it"
    echo ""
    echo "Available tools:"
    for dir in "$TOOLS_DIR"/*/; do
        if [ -d "$dir" ] && [ -f "$dir/setup.sh" ]; then
            echo "  $(basename "$dir")"
        fi
    done
    echo ""
    echo "Examples:"
    echo "  $0 substack /<path-to-ai-tool-config>/skills"
    echo "  $0 koyfin /<path-to-ai-tool-config>/skills"
}

list_available_tools() {
    echo "Available tools:"
    for dir in "$SCRIPT_DIR"/*/; do
        if [ -d "$dir" ] && [ -f "$dir/setup.sh" ]; then
            echo "  $(basename "$dir")"
        fi
    done
}

detect_platform() {
    local os arch
    case "$OSTYPE" in
        darwin*)      os="darwin" ;;
        msys*|win32*) os="windows" ;;
        *)            os="linux" ;;
    esac

    arch=$(uname -m)
    case "$arch" in
        x86_64)        arch="amd64" ;;
        aarch64|arm64) arch="arm64" ;;
        armv7l)        arch="arm" ;;
        i386|i686)     arch="386" ;;
    esac

    echo "$os $arch"
}

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
    local skills_dir="$3"
    local target_os="$4"
    local target_arch="$5"

    local project_root binary_name bin_name bin_file
    project_root="$(dirname "$tool_dir")"
    binary_name=$(get_binary_name "$tool_name" "$target_os")
    bin_name="$binary_name"
    bin_file="$skills_dir/scripts/$binary_name"

    if [ -f "$project_root/bin/$target_os/$bin_name" ]; then
        echo "✓ Found pre-built binary for $target_os"
        cp "$project_root/bin/$target_os/$bin_name" "$bin_file"
        chmod +x "$bin_file"
        echo "✓ Installed: $bin_file"
    elif command -v go &> /dev/null; then
        echo "Building from source..."
        echo "✓ Go: $(go version)"
        echo "✓ Target OS: $target_os"
        echo "✓ Target Arch: $target_arch"
        echo "✓ Deploying to: $skills_dir/scripts/"
        GOOS="$target_os" GOARCH="$target_arch" go build -ldflags='-w -s' -o "$bin_file" "$tool_dir/src/"
        chmod +x "$bin_file"
        echo "✓ Built: $bin_file"
    else
        echo "Error: No pre-built binary found and Go is not installed"
        echo ""
        echo "Options:"
        echo "  1. Download pre-built binary from releases"
        echo "  2. Install Go from https://go.dev/dl/"
        exit 1
    fi

    cp "$tool_dir/SKILL.md" "$skills_dir"

    if [ -d "$tool_dir/utils" ]; then
        echo "Copying utilities..."
        cp "$tool_dir/utils/"* "$skills_dir/scripts"
        echo "✓ Utilities: $skills_dir/scripts/"
    fi

    local skill_file="$skills_dir/$tool_name.json"
    cat > "$skill_file" << EOF
{
    "name": "$tool_name",
    "path": "$skills_dir/scripts",
    "binary": "$binary_name",
    "registered": "$(date -Iseconds)"
}
EOF

    echo ""
    echo "================================"
    echo "Setup complete!"
    echo ""
    echo "Tool deployed to: $skills_dir/scripts/"
    echo "Session file location: $skills_dir/scripts/session.json"
    echo ""
    echo "Next step: Authenticate"
    echo "  $skills_dir/scripts/$binary_name auth"
    echo ""
    echo "Usage:"
    echo "  $skills_dir/scripts/$binary_name <command> -h"

    if [ "$tool_name" = "koyfin" ]; then
        echo ""
        echo "Python utilities:"
        echo "  python3 $skills_dir/scripts/excel_export.py -h"
        echo "  $skills_dir/scripts/koyfin snapshot -kids <list_of_koyfin_ids> | python3 $skills_dir/scripts/excel_export.py -o snapshot.xlsx"
    fi

    echo ""
    echo "Tool registered: $skill_file"
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
        echo "Error: Tool '$tool_name' not found"
        echo ""
        list_available_tools
        exit 1
    fi

    if [ ! -f "$tool_dir/setup.sh" ]; then
        echo "Error: No setup.sh found for '$tool_name'"
        exit 1
    fi

    if [ -z "$skills_dir" ]; then
        echo ""
        echo "Building: $tool_name"
        echo ""
        read -p "Enter AI skills directory for deployment: " skills_dir
    fi

    if [ -z "$skills_dir" ]; then
        echo "Error: Skills directory is required"
        exit 1
    fi

    skills_dir="$skills_dir/$tool_name"
    mkdir -p "$skills_dir/scripts"

    read -r target_os target_arch <<< "$(detect_platform)"

    build_tool "$tool_name" "$tool_dir" "$skills_dir" "$target_os" "$target_arch"
}

main "$@"
