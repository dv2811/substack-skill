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

# Build and install a tool
./tool_build.sh <tool-name>
```

This installs the tool to your platform-specific binary directory and adds it to PATH.

## Available Tools

| Tool | Description | Documentation |
|------|-------------|---------------|
| `substack-reader` | Substack inbox, articles, search | `tools/substack-reader/SKILL.md` |

## Managing Large Outputs

Tools may produce large outputs. Pipe to files for processing:

```bash
# General pattern
<tool> <command> [args] > <tool>_<command>_<identifier>_<date>.txt

# Append for incremental results
<tool> <command> [args] >> <tool>_<command>_<date>.txt
```

**Filename conventions:**
- Include unique ID when available: `<tool>_item_123456.txt`
- Include query terms for searches: `<tool>_search_<query>_<date>.txt`
- Use dated filenames for recurring commands: `<tool>_inbox_YYYYMMDD.txt`

## Session Locations

Tools store sessions in platform-specific config directories:

- **Linux:** `~/.config/<tool-name>/session.json`
- **macOS:** `~/Library/Application Support/<tool-name>/session.json`
- **Windows:** `%APPDATA%\<tool-name>\session.json`

## Troubleshooting

**Command not found:**

The tool should be in your PATH after installation. Check platform-specific locations:

- **macOS:** `$HOME/bin/<tool>`
- **Linux:** `$HOME/.local/bin/<tool>`
- **Windows:** `%LOCALAPPDATA%\Programs\<tool-name>\<tool>.exe`

Add to PATH if needed:
```bash
# macOS
export PATH=$PATH:$HOME/bin

# Linux  
export PATH=$PATH:$HOME/.local/bin
```

## Requirements

- Go 1.21+
- Platform-specific account (varies by tool)

## Tool-Specific Documentation

See individual tool documentation for detailed usage:

- **substack-reader:** `tools/substack-reader/SKILL.md`
