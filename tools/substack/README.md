# Substack Reader CLI Tools

Command-line tools for accessing Substack reader data. Get your inbox posts, retrieve article content, and search Substack publications.

## Installation

### Prerequisites

- **Go 1.21+** (for building from source)
- **Substack account** (for API access)

### Quick Install

```bash
# From project root
./tool_build.sh substack-reader
```

Or run the setup script directly:

```bash
./tools/substack-reader/setup.sh
```

This will:
1. Build the `substack` binary from source
2. Install to platform-specific location:
   - **Linux**: `~/.local/bin/substack`
   - **macOS**: `~/bin/substack`
   - **Windows (WSL/Git Bash)**: `%LOCALAPPDATA%\Programs\substack-reader\substack`
3. Create session file at platform-specific config directory:
   - **Linux**: `~/.config/substack-reader/session.json`
   - **macOS**: `~/Library/Application Support/substack-reader/session.json`
   - **Windows**: `%APPDATA%\substack-reader\session.json`
4. Prompt you for Substack credentials

### Manual PATH Setup

If the binary directory is not in your PATH:

**Linux:**
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**macOS:**
```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zprofile
source ~/.zprofile
```

**Windows:**
Add to System Properties > Environment Variables > Path:
```
%LOCALAPPDATA%\Programs\substack-reader
```

### Pre-built Binaries

For distribution without requiring Go installation, place pre-built binaries in:

```
tools/substack-reader/bin/
├── linux/
│   └── substack
├── darwin/
│   └── substack
└── windows/
    └── substack.exe
```

Build all platforms:

```bash
# Linux
GOOS=linux go build -o tools/substack-reader/bin/linux/substack ./tools/substack-reader/src/

# macOS
GOOS=darwin go build -o tools/substack-reader/bin/darwin/substack ./tools/substack-reader/src/

# Windows
GOOS=windows go build -o tools/substack-reader/bin/windows/substack.exe ./tools/substack-reader/src/
```

The setup script will automatically use the pre-built binary for the current platform if available.

## Quick Start

```bash
# Get your inbox posts
substack inbox

# Get inbox posts after a specific timestamp
substack inbox -after "2024-01-01T00:00:00.000Z"

# Get article content by post ID
substack article -post-id 123456

# Search for posts (top results)
substack search -query "AI" -mode top

# Search all posts with pagination
substack search -query "technology" -mode all -page 1

# Search subscribed publications only
substack search -query "newsletter" -mode subscribed
```

## Commands

| Command | Description |
|---------|-------------|
| `auth` | Authenticate with Substack via email link |
| `inbox` | Get chronological inbox posts |
| `article` | Get article content by post ID |
| `search` | Search posts with different modes |

### auth

Authenticate with Substack using passwordless email link flow.

```bash
substack auth
```

**Authentication Flow:**

1. Run `substack auth`
2. Enter your email address when prompted
3. Check your email for a login link from Substack
4. Copy the full authentication URL from the email
5. Paste the URL into the terminal
6. Session is automatically saved

### inbox

Get chronologically sorted inbox posts.

```bash
substack inbox
substack inbox -after "2024-01-01T00:00:00.000Z"
```

| Flag | Description | Default |
|------|-------------|---------|
| `-after` | Timestamp cursor for pagination (RFC3339 format) | - |

### article

Get article content by post ID.

```bash
substack article -post-id 123456
substack article -post-id 123456 -base-url "substack.com"
```

| Flag | Description | Default |
|------|-------------|---------|
| `-post-id` | Post ID to retrieve (required) | - |
| `-base-url` | Custom Substack domain (optional) | - |

### search

Search Substack posts with different modes.

```bash
# Search top results
substack search -query "AI" -mode top

# Search all posts with pagination
substack search -query "technology" -mode all -page 1

# Search subscribed publications
substack search -query "newsletter" -mode subscribed -language en
```

| Flag | Description | Default |
|------|-------------|---------|
| `-query` | Search query (required) | - |
| `-mode` | Search mode: top, all, subscribed | all |
| `-page` | Page number (0-10, not used for top mode) | 0 |
| `-language` | Language code (2-letter, e.g., 'en') | - |

**Search Modes:**

- **top**: Returns top-ranked results for the query (pagination via cursor)
- **all**: Search all Substack posts (pagination via page number)
- **subscribed**: Search only within your subscribed publications

## Authentication

### Initial Setup

Run the setup script to create the configuration directory:

```bash
./tools/substack-reader/setup.sh
```

### Email Link Authentication

Use the `auth` command to authenticate with Substack:

```bash
substack auth
```

**How it works:**

1. Run `substack auth` and enter your email when prompted
2. Check your email for a login link from Substack
3. Copy the full authentication URL from the email
4. Paste the URL into the terminal
5. Session is automatically saved

### Session Location

Session files are stored at platform-specific locations:

- **Linux**: `~/.config/substack-reader/session.json`
- **macOS**: `~/Library/Application Support/substack-reader/session.json`
- **Windows**: `%APPDATA%\substack-reader\session.json`

## Output Format

All commands output JSON to stdout, making it easy to pipe to other tools:

```bash
# Pretty print with jq
substack inbox | jq '.data.posts[] | {title, post_date}'

# Count posts in inbox
substack inbox | jq '.data.posts | length'

# Search and extract article titles
substack search -query "AI" -mode top | jq '.data.results[].title'
```

## Session Management

- Session file location: `~/.config/substack-reader/session.json`
- Custom session file: Set `SUBSTACK_SESSION_FILE` environment variable
- Sessions are automatically saved after each command
- Session expiry is handled automatically with mid-point renewal

## Examples

### Get recent inbox posts

```bash
substack inbox | jq '.data.posts[] | {title, author: .publishedBylines[0].name, date: .post_date}'
```

### Search and get article content

```bash
# Find a post ID
POST_ID=$(substack search -query "AI" -mode top | jq '.data.results[0].id')

# Get the full article
substack article -post-id $POST_ID | jq '.data.post.body_html'
```

### Export inbox to file

```bash
substack inbox > inbox_$(date +%Y%m%d).json
```
