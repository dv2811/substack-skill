---
name: entext-research-tool
description: Research tool for accessing external content platforms via CLI. Use when you need to retrieve inbox posts, article content, or search publications.
---

# entext-research-tool

Research tool for accessing external content platforms via command line.

## Installation

```bash
# Clone to your projects directory
git clone https://github.com/dv2811/substack-skill.git
cd substack-skill

# Build a tool
./tool_build.sh <tool-name>
```

**Note:** Run `tool_build.sh` from the project root directory.

## Binary Location

After installation, binaries are located at:

- **Project root:** `<project-root>/bin/<tool-name>`
- **Add to PATH:** `export PATH=$PATH:$(pwd)/bin`

## Available Tools

| Tool | Description | Documentation |
|------|-------------|---------------|
| `substack-reader` | Substack inbox, articles, search | `tools/substack-reader/SKILL.md` |

## Managing Large Outputs

All tools output JSON to stdout. Pipe to files for large results:

```bash
# General pattern
./bin/<tool> <command> [args] > <tool>_<command>_<date>.json

# Process with jq later
jq '.data' <tool>_<command>_<date>.json
```

**Filename conventions:**
- Include unique ID when available: `<tool>_item_123456.json`
- Include query terms for searches: `<tool>_search_<query>_<date>.json`
- Use dated filenames for recurring commands: `<tool>_inbox_YYYYMMDD.json`

## Session Locations

Tools store sessions in platform-specific config directories:

- **Linux:** `~/.config/<tool-name>/session.json`
- **macOS:** `~/Library/Application Support/<tool-name>/session.json`
- **Windows:** `%APPDATA%\<tool-name>\session.json`

## Troubleshooting

**Command not found:**
```bash
# Use full path:
./bin/<tool> <command>

# Or add ./bin to PATH:
export PATH=$PATH:$(pwd)/bin
```

**Permission denied:**
```bash
chmod +x ./bin/<tool>
```

## Requirements

- Go 1.21+
- Platform-specific account (varies by tool)

## Tool-Specific Documentation

See individual tool documentation for detailed usage:

- **substack-reader:** `tools/substack-reader/SKILL.md`
