---
name: koyfin-cli
description: Koyfin CLI tools for financial data access including stock search, snapshot data, time series, earnings transcripts, ETF holdings, and stock screener. Use when you need to query Koyfin financial data from the command line.
---

# Koyfin CLI Tools

## Authentication

Run the setup script to authenticate:

```bash
./tool_build.sh koyfin
```

You will be prompted for your Koyfin email and password. Credentials are stored securely and tokens are auto-generated on first API call.

## CLI Commands

### search

Search for stocks/tickers by name.

```bash
koyfin search -q "Apple"
koyfin search -q "SPY" -categories "ETF"
```

| Flag | Description | Default |
|------|-------------|---------|
| `-q` | Ticker/ETF name (required) | - |
| `-categories` | Categories (comma-separated) | Equity,ETF |
| `-primary-only` | Primary exchange only | false |

### snapshot

Get snapshot data with financial metrics.

```bash
koyfin snapshot -kids <list_of_koyfin_ids> -category Equity
```

| Flag | Description | Default |
|------|-------------|---------|
| `-kids` | Comma-separated Koyfin IDs (required) | - |
| `-category` | Equity or ETF | Equity |

**Limits:** Max 32 KIDs (Equity), 2 KIDs (ETF)

### ticker-data

Get time series data for a ticker.

```bash
koyfin ticker-data -id <koyfin_id> -key "p_candle_range" -date-from "2024-01-01"
koyfin ticker-data -id <koyfin_id> -key "f_r" -date-from "2020-01-01" -fin-period "quarterly"
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
koyfin transcript -action list -kid <koyfin_id> -limit 5

# Get specific transcript
koyfin transcript -action get -transcript-id <transcript_id>

# Get transcript summary
koyfin transcript -action summary -transcript-id <transcript_id>
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
koyfin schema -asset-type Equity -indicator-type financials
koyfin schema -asset-type Equity -indicator-type ratios
```

| Flag | Description | Default |
|------|-------------|---------|
| `-asset-type` | Asset type | Equity |
| `-indicator-type` | financials, ratios, forward_estimates, market_data (required) | - |

### etf-holdings

Get ETF holdings.

```bash
koyfin etf-holdings -kids <list_of_koyfin_ids> -category ETF
```

| Flag | Description | Default |
|------|-------------|---------|
| `-kids` | Koyfin IDs, max 2 (required) | - |
| `-category` | Must be ETF | ETF |

### screener-schema

Get screener filter schema.

```bash
koyfin screener-schema -asset-type Equity
```

| Flag | Description | Default |
|------|-------------|---------|
| `-asset-type` | Asset type | Equity |

### screener

Run stock screener.

```bash
# Large cap (>10B)
koyfin screener -filters '[{"key":"mkt","min":10000,"max":9007199254740991}]'

# Tech sector, 1B-10B market cap
koyfin screener -filters '[{"key":"t_sec","values":["Information Technology"]},{"key":"mkt","min":1000,"max":10000}]'

# EV/EBITDA < 10
koyfin screener -filters '[{"key":"evebitdaltm","min":0,"max":10}]'
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

Location: `<binary_dir>/koyfin-utils/` (platform-specific)

| Platform | Python Utilities Path |
|----------|----------------------|
| **Linux** | `~/.local/bin/koyfin-utils/` |
| **macOS** | `~/bin/koyfin-utils/` |
| **Windows** | `%LOCALAPPDATA%\Programs\koyfin\koyfin-utils\` |

### Excel Export (Snapshot to XLSX)

Export snapshot data to Excel with multiple sheets:

```bash
koyfin snapshot -kids <list_of_koyfin_ids> -category Equity | \
    python3 $UTILS_DIR/excel_export.py -o <target_xlsx_file>
```

**Note:** Replace `$UTILS_DIR` with your platform-specific path from the table above.

**Excel Sheets Created:**

| Sheet | Description |
|-------|-------------|
| **Summary** | Key metrics: Price, Market Cap, P/E, EV/EBITDA, Margins, Growth |
| **Raw Data** | All snapshot metrics from API |
| **Ratios** | Calculated ratios: P/E, EV/EBITDA, EV/Sales, Margins, P/FCF |
| **Growth** | Growth rates: 1Y/3Y/5Y price CAGR, YTD, estimate vs LTM |

### Format Output

```bash
# Format search results as table
koyfin search -q "Apple" | python3 $UTILS_DIR/process.py search

# Format snapshot
koyfin snapshot -kids <list_of_koyfin_ids> | python3 $UTILS_DIR/process.py snapshot

# Format time series
koyfin ticker-data -id <koyfin_id> -key "p_candle_range" -date-from "2024-01-01" | \
    python3 $UTILS_DIR/process.py ticker
```

### Install Dependencies (for Excel export)

```bash
pip3 install -r $UTILS_DIR/requirements.txt
```

## Examples

```bash
# Search for Apple
koyfin search -q "Apple"

# Get snapshot
koyfin snapshot -kids <list_of_koyfin_ids> -category Equity

# Get 1 year price data
koyfin ticker-data -kid <koyfin_id> -key "p_candle_range" -date-from "2024-01-01"

# List recent transcripts
koyfin transcript -action list -kid <koyfin_id> -limit 5

# Get transcript
koyfin transcript -action get -transcript-id <transcript_id>

# Get transcript summary
koyfin transcript -action summary -transcript-id <transcript_id>

# Screen for large cap tech
koyfin screener -filters '[{"key":"t_sec","values":["Information Technology"]},{"key":"mkt","min":10000}]' -page-size 50
```
