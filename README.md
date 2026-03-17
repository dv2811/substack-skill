# Substack Skill

Go-based tools and libraries for accessing Substack reader data, including inbox posts, article content, and search functionality.

## Project Structure

```
substack-skill/
├── internal/
│   ├── substack/          # Core Substack API client
│   │   ├── client.go      # HTTP client and authentication
│   │   ├── session.go     # Session management
│   │   ├── inbox.go       # Inbox retrieval
│   │   ├── feed.go        # Feed operations
│   │   ├── search.go      # Post search
│   │   └── token.go       # Token handling
│   ├── models/            # Data models
│   │   └── session.go     # Session interface
│   └── validator/         # Input validation
│       └── validator.go   # Validation utilities
├── tools/
│   └── substack-reader/   # CLI tool
│       ├── src/
│       │   ├── main.go    # CLI entry point
│       │   ├── inbox.go   # inbox command
│       │   ├── article.go # article command
│       │   └── search.go  # search command
│       ├── setup.sh       # Installation script
│       ├── README.md      # CLI documentation
│       └── SKILL.md       # AI assistant skill doc
├── go.mod
└── README.md
```

## Components

### Internal Substack Package

The `internal/substack` package provides Go APIs for interacting with Substack:

- **Client** - HTTP client with Substack API endpoints
- **Session** - Authentication session management with cookie handling
- **Inbox** - Retrieve chronological inbox posts
- **Search** - Search posts with different modes (top, all, subscribed)

### Substack CLI Tool

Command-line interface for accessing Substack data. See [`tools/substack-reader/README.md`](tools/substack-reader/README.md) for full documentation.

#### Quick Start

```bash
# Install
./tools/substack-reader/setup.sh

# Get inbox posts
substack inbox

# Get article by ID
substack article -post-id 123456

# Search posts
substack search -query "AI" -mode top
```

#### Commands

| Command | Description |
|---------|-------------|
| `auth` | Authenticate with Substack via email link |
| `inbox` | Get chronological inbox posts |
| `article` | Get article content by post ID |
| `search` | Search posts (top/all/subscribed modes) |

## Installation

### Prerequisites

- Go 1.21+
- Substack account

### CLI Tool

**Using Makefile (recommended for development):**

```bash
# Build for current platform
make build

# Full setup (build + install + auth session)
make setup

# Install binary only (no session setup)
make install

# Uninstall binary
make uninstall

# Remove binary + session file
make uninstall-all

# Remove session file only (for security)
make clean-session
```

**Using the build script:**

```bash
# From project root
./tool_build.sh substack-reader
```

Or run the setup script directly:

```bash
cd substack-skill
./tools/substack-reader/setup.sh
```

This will:
1. Build the `substack` binary for your platform
2. Install to platform-specific location:
   - **Linux**: `~/.local/bin/substack`
   - **macOS**: `~/bin/substack`
   - **Windows**: `%LOCALAPPDATA%\Programs\substack-reader\substack.exe`
3. Create session file at platform-specific config directory:
   - **Linux**: `~/.config/substack-reader/session.json`
   - **macOS**: `~/Library/Application Support/substack-reader/session.json`
   - **Windows**: `%APPDATA%\substack-reader\session.json`

### Cross-Platform Builds

Build binaries for all platforms:

```bash
# Linux
GOOS=linux go build -o tools/substack-reader/bin/linux/substack ./tools/substack-reader/src/

# macOS
GOOS=darwin go build -o tools/substack-reader/bin/darwin/substack ./tools/substack-reader/src/

# Windows
GOOS=windows go build -o tools/substack-reader/bin/windows/substack.exe ./tools/substack-reader/src/
```

Place pre-built binaries in `tools/substack-reader/bin/<os>/` and the setup script will use them automatically.

### Library Usage

```go
import "entext-applications/internal/substack"

// Create client
client := substack.NewClient()

// Load session
session, err := substack.NewSessionFromFile("~/.config/substack-reader/session.json")

// Get inbox
inbox, err := client.GetChronologicalInbox(session, "")

// Search posts
results, err := client.SearchPosts(session, substack.SearchRequest{
    Query: "AI",
    Mode:  "top",
})
```

## Authentication

### Using the auth command

Authenticate with Substack using the `auth` command:

```bash
substack auth
```

**Flow:**
1. Enter your email address when prompted
2. Check your email and copy the authentication link
3. Paste the link into the terminal
4. Session is automatically saved

### Session locations

- **Linux**: `~/.config/substack-reader/session.json`
- **macOS**: `~/Library/Application Support/substack-reader/session.json`
- **Windows**: `%APPDATA%\substack-reader\session.json`
- **Custom**: Set `SUBSTACK_SESSION_FILE` environment variable

## Output Format

All CLI commands output JSON to stdout:

```bash
# Pretty print with jq
substack inbox | jq '.data.posts[] | {title, post_date}'

# Count posts
substack inbox | jq '.data.posts | length'

# Search and extract titles
substack search -query "AI" | jq '.data.results[].title'
```

## Session Management

- **Auto-renewal**: Sessions are automatically refreshed at mid-point
- **Signal handling**: Sessions saved on SIGINT/SIGTERM
- **Custom path**: Set `SUBSTACK_SESSION_FILE` environment variable

## Examples

### Get Recent Inbox Posts

```bash
substack inbox | jq '.data.posts[] | {title, author: .publishedBylines[0].name}'
```

### Search and Retrieve Article

```bash
# Find post ID
POST_ID=$(substack search -query "AI" -mode top | jq '.data.results[0].id')

# Get full article
substack article -post-id $POST_ID | jq '.data.post.body_html'
```

### Export Inbox

```bash
substack inbox > inbox_$(date +%Y%m%d).json
```

## License

This project is for personal use and educational purposes.
