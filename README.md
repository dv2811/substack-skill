# entext-research-tools

Go-based tools for accessing external content platforms.

## Manual Installation

```bash
# Clone the repository
git clone https://github.com/dv2811/entext-research-tools.git

# Change to the project directory
cd entext-research-tools

# Install a tool
./tool_build.sh <tool-name>
```
Then copy tool-specific SKILL.md to target AI tool's skills directory

## Project Structure

```
entext-research-tools/
├── .claude-plugin/
│   └── plugin.json        # Plugin manifest (Claude Code, auto-converts for Qwen)
├── tools/
│   ├── <tool-name>/       # CLI tools (one per directory)
│   │   ├── setup.sh       # Installation script
│   │   ├── README.md      # Tool documentation
│   │   └── src/           # Source code
│   └── ...                # Add more tools here
├── internal/
│   └── <service>/         # Service API clients
├── go.mod
└── README.md
```

## Available Tools

| Tool | Description | Documentation |
|------|-------------|---------------|
| `substack` | Substack inbox, articles, search | `tools/substack/README.md` |
| `koyfin` | Tool for interacting with Koyfin | `tools/koyfin/README.md` |

## Tool-Specific Documentation

See individual tool directories for detailed usage:

- **substack:** `tools/substack/README.md`
- **koyfin:** `tools/koyfin/README.md`
