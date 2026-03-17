# Suggested Improvements to SKILL.md

## Issue Encountered

When invoking the substack-reader skill, the command `./substack` was attempted from the skill directory (`~/.qwen/skills/substack-reader`), but the binary doesn't exist there after installation.

## Root Cause

The documentation states:
```bash
# Install
./tool_build.sh substack-reader
```

But doesn't clarify:
1. Where the binary is installed
2. How to invoke it after installation
3. That the binary lives in `./bin/substack` in the **project root**, not the skill directory

## Suggested Changes to SKILL.md

### 1. Add Binary Location Section

Add a clear section explaining where the binary is located:

```markdown
**Binary Location:**

After installation, the `substack` binary is located at:
- `<project-root>/bin/substack`
- Add `./bin` to your PATH for global access: `export PATH=$PATH:$(pwd)/bin`
```

### 2. Clarify Installation Context

Update the installation section to clarify the build happens from the project root:

```markdown
## Installation

```bash
# Clone to your projects directory
git clone https://github.com/dv2811/substack-skill.git
cd substack-skill

# Build the tool (creates binary in ./bin/)
./tool_build.sh substack-reader
```

**Note:** Run `tool_build.sh` from the project root directory (`/home/dv/substack-skill`), not from the skill directory.
```

### 3. Add Troubleshooting Section

```markdown
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
```

### 4. Update Usage Examples

Make it clear commands are run from project root:

```markdown
**Usage:**

```bash
# From project root directory
cd /path/to/substack-skill

# Authenticate
./bin/substack auth
# or if ./bin is in PATH:
substack auth

# Get inbox posts
./bin/substack inbox

# Get article by ID
./bin/substack article -post-id 123456

# Search posts
./bin/substack search -query "AI" -mode top
```
```

## Summary

The key improvements needed:
1. **Explicit binary path** - Document that binary is at `./bin/substack`
2. **Working directory clarity** - Specify commands run from project root
3. **PATH setup instructions** - Show how to make `substack` globally accessible
4. **Troubleshooting section** - Help users resolve common "command not found" errors
5. **No guidance on managing large outputs** - Add examples for piping to files

---

## Add: Managing Large Outputs

When working with search results or article content, outputs can be large and consume significant context. Add a section on piping to files:

```markdown
## Managing Large Outputs

Search results and article content can be large. Pipe intermediate results to files to reduce context load:

```bash
# Naming convention: <tool>_<command>_<unique_id|query>_<date>.json

# Article by ID - use post_id for unique filename
./bin/substack article -post-id 191095022 > substack_article_191095022.json

# Search - use query terms + timestamp
./bin/substack search -query "openclaw china" > substack_search_openclaw_china_20260317.json

# Inbox - use command + date
./bin/substack inbox > substack_inbox_20260317.json

# Process file later with jq
jq '.data.post.title' substack_article_191095022.json
jq '.data.results[] | {title, id}' substack_search_openclaw_china_20260317.json
```

**Why pipe to files?**
- Reduces context token usage when working with large result sets
- Allows incremental processing of search results
- Enables re-processing without re-fetching from API

**Filename conventions to avoid collisions:**
- **Always include unique ID when available** (post_id, etc.): `substack_article_191095022.json`
- **Include query terms for searches**: `substack_search_<query>_<date>.json`
- **Use dated filenames for recurring commands**: `substack_inbox_YYYYMMDD.json`
- **Consider a dedicated output directory**: `mkdir -p output && ./bin/substack ... > output/...`
```
