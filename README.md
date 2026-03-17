# Substack Skill

Go-based tools for accessing external services.

## Project Structure

```
substack-skill/
├── internal/
│   └── <service>/         # Service API clients
├── tools/
│   ├── <tool-name>/       # CLI tools (one per directory)
│   │   ├── setup.sh       # Installation script
│   │   ├── README.md      # Tool documentation
│   │   └── src/           # Source code
│   └── ...                # Add more tools here
├── go.mod
└── README.md
```

## Tools

- **substack-reader** - Substack CLI tool (see `tools/substack-reader/README.md`)

## Installation

Use the build script to install tools:

```bash
./tool_build.sh <tool-name>
```
