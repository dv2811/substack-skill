---
name: substack-reader
description: Substack CLI tool for accessing inbox posts, articles, and search. Use when you need to retrieve Substack reader data via command line.
---

# Substack Reader CLI

Part of the **entext-research-tool** skill suite.

## Installation

```bash
# Clone to your projects directory
git clone https://github.com/dv2811/substack-skill.git
cd substack-skill

# Build and install the tool
./tool_build.sh substack-reader
```

Or run the setup script directly:

```bash
cd tools/substack-reader
./setup.sh
```

After installation, `substack` is available in your PATH.

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

## Managing Large Outputs

Article content and search results can be large. Pipe to files for processing:

```bash
# Article by ID
substack article -post-id 191095022 > article_191095022.json

# Search with query terms + timestamp
substack search -query "openclaw china" > search_openclaw_china_$(date +%Y%m%d).json

# Inbox with date
substack inbox > inbox_$(date +%Y%m%d).json

# Process file later
jq '.data.post.title' article_191095022.json
jq '.data.results[] | {title, id}' search_openclaw_china_$(date +%Y%m%d).json
```

**Filename conventions:**
- Include unique ID when available: `article_191095022.json`
- Include query terms for searches: `search_<query>_<date>.json`
- Use dated filenames for recurring commands: `inbox_YYYYMMDD.json`

## Output Format

Commands output JSON to stdout for piping to `jq` or files:

```bash
# Pretty print
substack inbox | jq '.data.posts[] | {title, post_date}'

# Extract search titles
substack search -query "AI" | jq '.data.results[].title'

# Save to file
substack inbox > inbox_$(date +%Y%m%d).json
```

## Session Location

- **Linux:** `~/.config/substack-reader/session.json`
- **macOS:** `~/Library/Application Support/substack-reader/session.json`
- **Windows:** `%APPDATA%\substack-reader\session.json`

## Troubleshooting

**Command not found:**

The `substack` binary should be in your PATH after installation. Check platform-specific locations:

- **macOS:** `$HOME/bin/substack`
- **Linux:** `$HOME/.local/bin/substack`
- **Windows:** `%LOCALAPPDATA%\Programs\substack-reader\substack.exe`

Add to PATH if needed:
```bash
# macOS
export PATH=$PATH:$HOME/bin

# Linux
export PATH=$PATH:$HOME/.local/bin
```

Or use the full path:
```bash
$HOME/bin/substack <command>  # macOS
$HOME/.local/bin/substack <command>  # Linux
```

**Permission denied:**
```bash
chmod +x $(which substack)
```

## Requirements

- Go 1.21+
- Substack account
