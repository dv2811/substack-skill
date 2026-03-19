---
name: substack
description: Substack CLI tool for accessing inbox posts, articles, and search. Use when you need to retrieve Substack reader data via command line.
---

# Substack Reader CLI

Part of the **entext-research-tool** skill suite.

## Installation

The tool is built to the AI skills directory's `scripts/` folder.

## Authentication

**First-time setup (2-step process):**

```bash
# Step 1: Set email address
./scripts/substack profile -email "user@example.com"

# Step 2: Complete with authentication code or link (obtained from email)
./scripts/substack auth -auth-link "https://substack.com/auth?token=..."
```

Session persists across CLI invocations and auto-refreshes.

## Commands

### profile

Set Substack email address for authentication.

```bash
./scripts/substack profile -email "user@example.com"
```

| Flag | Description | Required |
|------|-------------|----------|
| `-email` | Substack email address | Yes |

### auth

Complete Substack authentication with URL or code from email.

```bash
./scripts/substack auth -auth_string "https://substack.com/auth?token=..."
```

| Flag | Description | Required |
|------|-------------|----------|
| `-auth_string` | Authentication code or URL from email | Yes |

### inbox

Get chronological inbox posts.

```bash
./scripts/substack inbox
./scripts/substack inbox -after "2024-01-01T00:00:00.000Z"
```

### article

Get article content by post ID.

```bash
./scripts/substack article -post-id 123456
```

### search

Search Substack posts.

```bash
./scripts/substack search -query "AI" -mode top
./scripts/substack search -query "tech" -mode all
./scripts/substack search -query "news" -mode subscribed
```

**Search Modes:**
- `top` - Top-ranked results
- `all` - All posts (default)
- `subscribed` - Subscribed publications only

## Examples

```bash
# Non-interactive authentication (automation)
./scripts/substack profile -email "user@example.com"
./scripts/substack auth -auth-link "https://substack.com/auth?token=..."

# Get inbox and extract titles
./scripts/substack inbox | jq '.data.posts[] | .title'

# Search and get article
POST_ID=$(./scripts/substack search -query "AI" -mode top | jq '.data.results[0].id')
./scripts/substack article -post-id $POST_ID

# Count posts in inbox
./scripts/substack inbox | jq '.data.posts | length'
```

## Managing Large Outputs

Article content and search results can be large. Pipe to files for processing:

```bash
# Article by ID
./scripts/substack article -post-id 191095022 > article_191095022.json

# Search with query terms + timestamp
./scripts/substack search -query "openclaw china" > search_openclaw_china_$(date +%Y%m%d).json

# Inbox with date
./scripts/substack inbox > inbox_$(date +%Y%m%d).json

# Process file later
jq '.data.post.title' article_191095022.json
jq '.data.results[] | {title, id}' search_openclaw_china_$(date +%Y%m%d).json
```

**Filename conventions:**
- Include unique ID when available: `article_191095022.json`
- Include query terms for searches: `search_<query>_<date>.json`
- Use dated filenames for recurring commands: `inbox_YYYYMMDD.json`

## Session Location

Session file is stored in the scripts directory:

```
./scripts/session.json
```

**Note:** This directory is excluded from version control (`.gitignore`) to protect authentication tokens.

## Troubleshooting

**Command not found:**

Always use the full path from the skills directory:
```bash
./scripts/substack <command>
```

**Permission denied:**
```bash
chmod +x ./scripts/substack
```

**Session expired:**
```bash
./scripts/substack profile -email "user@example.com"
./scripts/substack auth
```

**Email not set:**
```bash
./scripts/substack profile -email "user@example.com"
```

## Requirements

- Go 1.21+
- Substack account
