---
name: substack-reader
description: Substack CLI tool for accessing inbox posts, articles, and search. Use when you need to retrieve Substack reader data via command line.
---

# Substack Reader CLI

Part of the **entext-research-tool** skill suite.

## Installation

```bash
# From project root
./tool_build.sh substack-reader
```

Or run the setup script directly:

```bash
cd tools/substack-reader
./setup.sh
```

## Authentication

**First-time setup:**

```bash
substack auth
```

1. Enter your email address when prompted
2. Check your email and copy the authentication link
3. Paste the link into the terminal
4. Session is automatically saved

Session persists across CLI invocations and auto-refreshes.

## Commands

### auth

Authenticate with Substack via email link.

```bash
substack auth
```

### inbox

Get chronological inbox posts.

```bash
substack inbox
substack inbox -after "2024-01-01T00:00:00.000Z"
```

### article

Get article content by post ID.

```bash
substack article -post-id 123456
```

### search

Search Substack posts.

```bash
substack search -query "AI" -mode top
substack search -query "tech" -mode all
substack search -query "news" -mode subscribed
```

**Search Modes:**
- `top` - Top-ranked results
- `all` - All posts (default)
- `subscribed` - Subscribed publications only

## Examples

```bash
# Get inbox and extract titles
substack inbox | jq '.data.posts[] | .title'

# Search and get article
POST_ID=$(substack search -query "AI" -mode top | jq '.data.results[0].id')
substack article -post-id $POST_ID

# Count posts in inbox
substack inbox | jq '.data.posts | length'
```

## Output Format

All commands output JSON to stdout for piping to `jq`:

```bash
# Pretty print
substack inbox | jq '.data.posts[] | {title, post_date}'

# Extract search titles
substack search -query "AI" | jq '.data.results[].title'
```

## Session Location

- **Linux:** `~/.config/substack-reader/session.json`
- **macOS:** `~/Library/Application Support/substack-reader/session.json`
- **Windows:** `%APPDATA%\substack-reader\session.json`

## Requirements

- Go 1.21+
- Substack account
