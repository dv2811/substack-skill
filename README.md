# Substack Skill

Go-based tools and libraries for accessing Substack reader data.

## Project Structure

```
substack-skill/
├── internal/
│   └── substack/          # Core Substack API client
├── tools/
│   └── substack-reader/   # CLI tool
│       ├── setup.sh       # Installation script
│       └── README.md      # Tool documentation
├── go.mod
└── README.md
```

## Components

### Internal Substack Package

The `internal/substack` package provides Go APIs for interacting with Substack:

- HTTP client with authentication
- Session management
- Inbox retrieval
- Search functionality

### CLI Tool

See [`tools/substack-reader/README.md`](tools/substack-reader/README.md) for installation and usage documentation.

## Library Usage

```go
import "entext-applications/internal/substack"

client := substack.NewClient()
session, err := substack.NewSessionFromFile("~/.config/substack-reader/session.json")
```
