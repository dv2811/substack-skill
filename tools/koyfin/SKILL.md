---
name: koyfin
description: Koyfin CLI tools for financial data access including stock search, snapshot data, time series, earnings transcripts, ETF holdings, and stock screener. Use when you need to query Koyfin financial data from the command line.
---

# Koyfin CLI Tools

## Installation

The tool is built to the AI skills directory's `scripts/` folder. Python utilities are copied to the same directory.

## CLI Commands

### auth

**Authenticate with Koyfin account email and password:**

```bash
./scripts/koyfin auth -email "user@example.com" -password "secret"
```

| Flag | Description | Required |
|------|-------------|----------|
| `-email` | Koyfin email address | Yes |
| `-password` | Koyfin password | Yes |

Session persists across CLI invocations and tokens refreshed automatically.

### search

Search for stocks/tickers by name.

```bash
./scripts/koyfin search -q "Apple"
./scripts/koyfin search -q "SPY" -categories "ETF"
```

| Flag | Description | Default |
|------|-------------|---------|
| `-q` | Ticker/ETF name (required) | - |
| `-categories` | Categories (comma-separated) | Equity,ETF |
| `-primary-only` | Primary exchange only | false |

### snapshot

Get snapshot data with financial metrics.

```bash
./scripts/koyfin snapshot -kids <list_of_koyfin_ids> -category Equity
```

| Flag | Description | Default |
|------|-------------|---------|
| `-kids` | Comma-separated Koyfin IDs (required) | - |
| `-category` | Equity or ETF | Equity |

**Limits:** Max 32 KIDs (Equity), 2 KIDs (ETF)

### ticker-data

Get time series data for a ticker.

```bash
./scripts/koyfin ticker-data -id <koyfin_id> -key "p_candle_range" -date-from "2024-01-01"
./scripts/koyfin ticker-data -id <koyfin_id> -key "f_r" -date-from "2020-01-01" -fin-period "quarterly"
```

| Flag | Description | Default |
|------|-------------|---------|
| `-id` | Koyfin ID (required) | - |
| `-key` | Indicator key (required) | - |
| `-date-from` | Start date YYYY-MM-DD (required) | - |
| `-date-to` | End date YYYY-MM-DD | Today |
| `-currency` | Data currency | USD |
| `-agg-period` | day, monthly, quarterly, annually | day |
| `-fin-period` | quarterly, annual, LTM | - |

### transcript

Earnings call transcripts (list, get, summary).

```bash
# List transcripts
./scripts/koyfin transcript -action list -kid <koyfin_id> -limit 5

# Get specific transcript
./scripts/koyfin transcript -action get -transcript-id <transcript_id>

# Get transcript summary
./scripts/koyfin transcript -action summary -transcript-id <transcript_id>
```

| Flag | Description | Default |
|------|-------------|---------|
| `-action` | list, get, summary (required) | list |
| `-kid` | Koyfin ID (required for list) | - |
| `-transcript-id` | Key dev ID (required for get/summary) | - |
| `-limit` | Max results for list (1-64) | 10 |

### schema

Get indicator schema.

```bash
./scripts/koyfin schema -asset-type Equity -indicator-type financials
./scripts/koyfin schema -asset-type Equity -indicator-type ratios
```

| Flag | Description | Default |
|------|-------------|---------|
| `-asset-type` | Asset type | Equity |
| `-indicator-type` | financials, ratios, forward_estimates, market_data (required) | - |

### etf-holdings

Get ETF holdings.

```bash
./scripts/koyfin etf-holdings -kids <list_of_koyfin_ids> -category ETF
```

| Flag | Description | Default |
|------|-------------|---------|
| `-kids` | Koyfin IDs, max 2 (required) | - |
| `-category` | Must be ETF | ETF |

### screener-schema

Get screener filter schema.

```bash
./scripts/koyfin screener-schema -asset-type Equity
```

| Flag | Description | Default |
|------|-------------|---------|
| `-asset-type` | Asset type | Equity |

### screener

Run stock screener.

```bash
# Large cap (>10B)
./scripts/koyfin screener -filters '[{"key":"mkt","min":10000,"max":9007199254740991}]'

# Tech sector, 1B-10B market cap
./scripts/koyfin screener -filters '[{"key":"t_sec","values":["Information Technology"]},{"key":"mkt","min":1000,"max":10000}]'

# EV/EBITDA < 10
./scripts/koyfin screener -filters '[{"key":"evebitdaltm","min":0,"max":10}]'
```

| Flag | Description | Default |
|------|-------------|---------|
| `-filters` | JSON filter array (required) | - |
| `-page-size` | Max 300 | 100 |

**Filter Types & Units:**

| Filter Key | Unit | Example |
|------------|------|---------|
| `mkt` | Millions USD | `{"min":1000}` = >$1B |
| `evebitdaltm` | Ratio (not %) | `{"max":10}` = <10x |
| `pf_fcf_margin-LTM` | Decimal (not %) | `{"min":0.10}` = >10% |
| `f_bet1yr` | Decimal (not %) | `{"min":0.05}` = >5% |
| `chg1mPct_L`, `chg3mPct_L`, `chgYTDPct_L` | Decimal (not %) | `{"min":0.10}` = >10% |
| `pe_ratio` | Ratio | `{"max":20}` = <20x |
| `pb_ratio` | Ratio | `{"max":3}` = <3x |
| `div_yield` | Decimal (not %) | `{"min":0.03}` = >3% |
| `iso2` | Country code | `{"values":["JP","US"]}` |
| `t_sec` | Sector name | `{"values":["Technology"]}` |

## Python Utilities

Location: `scripts/` (same directory as binary)

| File | Description |
|------|-------------|
| `excel_export.py` | Export snapshot data to Excel with multiple sheets |
| `process.py` | Format output as tables |
| `requirements.txt` | Python dependencies for Excel export |

### Excel Export (Snapshot to XLSX)

Export snapshot data to Excel with multiple sheets:

```bash
./scripts/koyfin snapshot -kids <list_of_koyfin_ids> -category Equity | \
    python3 ./scripts/excel_export.py -o <target_xlsx_file>
```

**Excel Sheets Created:**

| Sheet | Description |
|-------|-------------|
| **Summary** | Key metrics: Price, Market Cap, P/E, EV/EBITDA, Margins, Growth |
| **Raw Data** | All snapshot metrics from API |
| **Ratios** | Calculated ratios: P/E, EV/EBITDA, EV/Sales, Margins, P/FCF |
| **Growth** | Growth rates: 1Y/3Y/5Y price CAGR, YTD, estimate vs LTM |

### Install Dependencies (for Excel export)

```bash
pip3 install -r ./scripts/requirements.txt
```

### Format Output

```bash
# Format search results as table
./scripts/koyfin search -q "Apple" | python3 ./scripts/process.py search

# Format snapshot
./scripts/koyfin snapshot -kids <list_of_koyfin_ids> | python3 ./scripts/process.py snapshot

# Format time series
./scripts/koyfin ticker-data -id <koyfin_id> -key "p_candle_range" -date-from "2024-01-01" | \
    python3 ./scripts/process.py ticker
```

## Examples

```bash
# Interactive authentication
./scripts/koyfin auth

# Non-interactive authentication (automation)
./scripts/koyfin auth -email "user@example.com" -password "secret"

# Search for Apple
./scripts/koyfin search -q "Apple"

# Get snapshot
./scripts/koyfin snapshot -kids <list_of_koyfin_ids> -category Equity

# Get 1 year price data
./scripts/koyfin ticker-data -kid <koyfin_id> -key "p_candle_range" -date-from "2024-01-01"

# List recent transcripts
./scripts/koyfin transcript -action list -kid <koyfin_id> -limit 5

# Get transcript
./scripts/koyfin transcript -action get -transcript-id <transcript_id>

# Get transcript summary
./scripts/koyfin transcript -action summary -transcript-id <transcript_id>

# Screen for large cap tech
./scripts/koyfin screener -filters '[{"key":"t_sec","values":["Information Technology"]},{"key":"mkt","min":10000}]' -page-size 50

# Export snapshot to Excel
./scripts/koyfin snapshot -kids <list_of_koyfin_ids> | python3 ./scripts/excel_export.py -o snapshot.xlsx
```

## Session Location

Session file is stored in the scripts directory:

```
scripts/session.json
```

**Note:** This directory is excluded from version control (`.gitignore`) to protect authentication tokens.

## Troubleshooting

**Command not found:**

Always use the full path from the skills directory:
```bash
./scripts/koyfin <command>
```

**Permission denied:**
```bash
chmod +x ./scripts/koyfin
```

**Session expired:**
```bash
./scripts/koyfin auth -email "user@example.com" -password "secret"
```

## Requirements

- Go 1.21+
- Koyfin account
- Python 3 (optional, for Excel export utilities)
