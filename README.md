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

Example:

```bash
./tool_build.sh substack-reader
```

## Install as AI Skill

### Claude Code

Automatically discovers `SKILL.md` from repository.

### Qwen Code

```bash
qwen --experimental-skills
```

### OpenClaw

```bash
openclaw skill install dv2811/substack-skill
```

## Project Structure

```
entext-research-tool/
├── SKILL.md               # AI skill definition
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

- **substack-reader** - Substack CLI tool (see `tools/substack-reader/README.md`)
