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

# Build the tool (creates binary in ./bin/)
./tool_build.sh substack-reader
```

Or run the setup script directly:

```bash
cd tools/substack-reader
./setup.sh
```

**Note:** Run `tool_build.sh` from the project root directory, not from the skill directory.

## Binary Location

After installation, the `substack` binary is located at:

- **Project root:** `<project-root>/bin/substack`
- **Add to PATH:** `export PATH=$PATH:$(pwd)/bin`

## Authentication

**First-time setup:**

```bash
./bin/substack auth
# Or if ./bin is in PATH:
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
./bin/substack auth
```

### inbox

Get chronological inbox posts.

```bash
./bin/substack inbox
./bin/substack inbox -after "2024-01-01T00:00:00.000Z"
```

### article

Get article content by post ID.

```bash
./bin/substack article -post-id 123456
```

### search

Search Substack posts.

```bash
./bin/substack search -query "AI" -mode top
./bin/substack search -query "tech" -mode all
./bin/substack search -query "news" -mode subscribed
```

**Search Modes:**
- `top` - Top-ranked results
- `all` - All posts (default)
- `subscribed` - Subscribed publications only

## Examples

```bash
# Get inbox and extract titles
./bin/substack inbox | jq '.data.posts[] | .title'

# Search and get article
POST_ID=$(./bin/substack search -query "AI" -mode top | jq '.data.results[0].id')
./bin/substack article -post-id $POST_ID

# Count posts in inbox
./bin/substack inbox | jq '.data.posts | length'
```

## Managing Large Outputs

Search results and article content can be large. Pipe to files to reduce context load:

```bash
# Article by ID - use post_id for unique filename
./bin/substack article -post-id 191095022 > substack_article_191095022.json

# Search - use query terms + timestamp
./bin/substack search -query "openclaw china" > substack_search_openclaw_china_$(date +%Y%m%d).json

# Inbox - use command + date
./bin/substack inbox > substack_inbox_$(date +%Y%m%d).json

# Process file later with jq
jq '.data.post.title' substack_article_191095022.json
jq '.data.results[] | {title, id}' substack_search_openclaw_china_$(date +%Y%m%d).json
```

**Why pipe to files?**
- Reduces context token usage when working with large result sets
- Allows incremental processing of search results
- Enables re-processing without re-fetching from API

**Filename conventions:**
- Include unique ID when available: `substack_article_191095022.json`
- Include query terms for searches: `substack_search_<query>_<date>.json`
- Use dated filenames for recurring commands: `substack_inbox_YYYYMMDD.json`

## Output Format

All commands output JSON to stdout for piping to `jq` or files:

```bash
# Pretty print
./bin/substack inbox | jq '.data.posts[] | {title, post_date}'

# Extract search titles
./bin/substack search -query "AI" | jq '.data.results[].title'

# Save to file
./bin/substack inbox > inbox_$(date +%Y%m%d).json
```

## Session Location

- **Linux:** `~/.config/substack-reader/session.json`
- **macOS:** `~/Library/Application Support/substack-reader/session.json`
- **Windows:** `%APPDATA%\substack-reader\session.json`

## Troubleshooting

**Command not found:**
```bash
# If `substack` command is not found, use the full path:
./bin/substack <command>

# Or add ./bin to your PATH:
export PATH=$PATH:$(pwd)/bin
```

**Permission denied:**
```bash
chmod +x ./bin/substack
```

## Requirements

- Go 1.21+
- Substack account
