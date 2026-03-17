# entext-research-tool

Go-based tools for accessing external content platforms.

## Manual Installation

```bash
# Clone the repository
git clone https://github.com/dv2811/substack-skill.git

# Change to the project directory
cd substack-skill

# Install a tool
./tool_build.sh <tool-name>
```

### Install as AI Skill

To use tools with AI assistants that support skills:

**Claude Code:**
```bash
cp tools/<tool-name>/SKILL.md ~/.claude/skills/<tool-name>/SKILL.md
```

**Qwen Code:**
```bash
cp tools/<tool-name>/SKILL.md ~/.qwen/skills/<tool-name>/SKILL.md
qwen --experimental-skills
```

**OpenClaw:**
```bash
openclaw skill install dv2811/substack-skill
```

## Project Structure

```
entext-research-tool/
├── SKILL.md               # Root skill definition
├── manifest.json          # OpenClaw permission manifest
├── internal/
│   └── <service>/         # Service API clients
├── tools/
│   ├── <tool-name>/       # CLI tools (one per directory)
│   │   ├── setup.sh       # Installation script
│   │   ├── README.md      # Tool documentation
│   │   ├── SKILL.md       # Tool skill definition
│   │   └── src/           # Source code
│   └── ...                # Add more tools here
├── go.mod
└── README.md
```

## Available Tools

| Tool | Description | Documentation |
|------|-------------|---------------|
| `substack-reader` | Substack inbox, articles, search | `tools/substack-reader/README.md` |

## Tool-Specific Documentation

See individual tool directories for detailed usage:

- **substack-reader:** `tools/substack-reader/README.md`, `tools/substack-reader/SKILL.md`
