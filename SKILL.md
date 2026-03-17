---
name: entext-research-tool
description: Research tool for accessing external content platforms via CLI. Use when you need to retrieve inbox posts, article content, or search publications like Substack.
---

# entext-research-tool

Research tool for accessing external content platforms via command line.

## Installation

```bash
git clone https://github.com/dv2811/substack-skill.git
cd substack-skill
./tool_build.sh <tool-name>
```

## Available Tools

### substack-reader

CLI tool for accessing Substack reader data.

**Commands:**

| Command | Description |
|---------|-------------|
| `auth` | Authenticate with Substack via email link |
| `inbox` | Get chronological inbox posts |
| `article` | Get article content by post ID |
| `search` | Search posts (top/all/subscribed modes) |

**Usage:**

```bash
# Install
./tool_build.sh substack-reader

# Authenticate (required first step)
substack auth

# Get inbox posts
substack inbox

# Get article by ID
substack article -post-id 123456

# Search posts
substack search -query "AI" -mode top
```

**Output Format:**

All commands output JSON to stdout for piping to `jq`:

```bash
# Pretty print inbox
substack inbox | jq '.data.posts[] | {title, post_date}'

# Search and extract titles
substack search -query "AI" | jq '.data.results[].title'
```

## Capabilities

- **Authentication** - Email link-based authentication with automatic session management
- **Content Retrieval** - Fetch inbox posts and full article content
- **Search** - Search posts across different modes (top, all, subscribed)
- **Session Management** - Automatic session persistence and refresh

## Session Locations

- **Linux:** `~/.config/substack-reader/session.json`
- **macOS:** `~/Library/Application Support/substack-reader/session.json`
- **Windows:** `%APPDATA%\substack-reader\session.json`

## Requirements

- Go 1.21+
- Substack account (for substack-reader tool)
