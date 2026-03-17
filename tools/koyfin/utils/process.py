#!/usr/bin/env python3
"""
Koyfin CLI Result Processing Utilities

Process and format output from koyfin CLI commands.
"""

import json
import sys
from typing import Any, Dict, List, Optional


def parse_cli_output(output: str) -> Any:
    """Parse JSON output from koyfin CLI."""
    return json.loads(output)


def format_snapshot(data: Dict[str, Any], metrics: List[str] = None) -> str:
    """
    Format snapshot data as readable table.
    
    Args:
        data: Raw snapshot JSON from CLI
        metrics: Specific metrics to display (default: all)
    
    Returns:
        Formatted table string
    """
    results = data.get("results", [])
    lines = []
    
    for item in results:
        kid = item.get("kids", "N/A")
        metrics_data = item.get("metrics", {})
        
        lines.append(f"\n{kid}")
        lines.append("-" * 40)
        
        display_metrics = metrics or metrics_data.keys()
        for key in display_metrics:
            if key in metrics_data:
                value = metrics_data[key]
                if isinstance(value, float):
                    value = f"{value:,.2f}"
                lines.append(f"  {key}: {value}")
    
    return "\n".join(lines)


def format_search_results(data: Dict[str, Any]) -> str:
    """
    Format search results as table.
    
    Args:
        data: Raw search JSON from CLI
    
    Returns:
        Formatted table string
    """
    results = data.get("data", [])
    lines = []
    
    lines.append(f"{'KID':<15} {'Ticker':<12} {'Name':<30} {'Category':<10}")
    lines.append("-" * 70)
    
    for item in results:
        kid = item.get("KID", "N/A")[:14]
        ticker = item.get("identifier", "N/A")[:11]
        name = item.get("name", "N/A")[:29]
        category = item.get("category", "N/A")[:9]
        lines.append(f"{kid:<15} {ticker:<12} {name:<30} {category:<10}")
    
    return "\n".join(lines)


def format_ticker_data(data: Dict[str, Any], rows: int = 10) -> str:
    """
    Format time series data.
    
    Args:
        data: Raw ticker-data JSON from CLI
        rows: Number of rows to display
    
    Returns:
        Formatted table string
    """
    ts = data.get("data", {})
    labels = ts.get("labels", [])
    values = ts.get("values", [])
    
    lines = []
    lines.append(f"{'Date':<15} {'Value':>15}")
    lines.append("-" * 32)
    
    for i in range(min(rows, len(labels))):
        lines.append(f"{labels[i]:<15} {values[i]:>15,.2f}")
    
    if len(labels) > rows:
        lines.append(f"... ({len(labels) - rows} more rows)")
    
    return "\n".join(lines)


def format_transcripts(data: Dict[str, Any]) -> str:
    """
    Format transcript list.
    
    Args:
        data: Raw transcript JSON from CLI
    
    Returns:
        Formatted list string
    """
    transcripts = data.get("data", [])
    lines = []
    
    for t in transcripts:
        date = t.get("announcedDate", "N/A")
        title = t.get("transcriptTitle", "N/A")
        key_id = t.get("keyDevId", "N/A")
        lines.append(f"[{date}] {title}")
        lines.append(f"  ID: {key_id}")
        lines.append("")
    
    return "\n".join(lines)


def format_screener_results(data: Dict[str, Any], max_display: int = 20) -> str:
    """
    Format screener results.
    
    Args:
        data: Raw screener JSON from CLI
        max_display: Maximum KIDs to display
    
    Returns:
        Formatted list string
    """
    kids = data.get("kids", [])
    count = len(kids)
    
    lines = [f"Found {count} results:"]
    lines.append("")
    
    display_kids = kids[:max_display]
    for kid in display_kids:
        lines.append(f"  {kid}")
    
    if count > max_display:
        lines.append(f"  ... and {count - max_display} more")
    
    return "\n".join(lines)


def extract_metric(data: Dict[str, Any], kid: str, metric: str) -> Optional[float]:
    """
    Extract specific metric for a KID from snapshot data.
    
    Args:
        data: Raw snapshot JSON
        kid: Koyfin ID to find
        metric: Metric name to extract
    
    Returns:
        Metric value or None
    """
    results = data.get("results", [])
    for item in results:
        if item.get("kids") == kid:
            metrics = item.get("metrics", {})
            return metrics.get(metric)
    return None


def compare_snapshots(before: Dict[str, Any], after: Dict[str, Any], 
                      metric: str) -> str:
    """
    Compare a metric between two snapshots.
    
    Args:
        before: Earlier snapshot data
        after: Later snapshot data
        metric: Metric to compare
    
    Returns:
        Comparison report
    """
    lines = [f"Comparison: {metric}"]
    lines.append("=" * 40)
    
    before_results = {r.get("kids"): r.get("metrics", {}).get(metric) 
                      for r in before.get("results", [])}
    after_results = {r.get("kids"): r.get("metrics", {}).get(metric) 
                     for r in after.get("results", [])}
    
    all_kids = set(before_results.keys()) | set(after_results.keys())
    
    for kid in sorted(all_kids):
        before_val = before_results.get(kid)
        after_val = after_results.get(kid)
        
        if before_val is not None and after_val is not None:
            change = after_val - before_val
            pct = (change / before_val * 100) if before_val != 0 else 0
            sign = "+" if change > 0 else ""
            lines.append(f"{kid}: {before_val:,.2f} → {after_val:,.2f} "
                        f"({sign}{change:,.2f}, {sign}{pct:.1f}%)")
        elif before_val is not None:
            lines.append(f"{kid}: {before_val:,.2f} → N/A")
        else:
            lines.append(f"{kid}: N/A → {after_val:,.2f}")
    
    return "\n".join(lines)


if __name__ == "__main__":
    """Process stdin from koyfin CLI pipe."""
    if len(sys.argv) < 2:
        print("Usage: koyfin <command> | python -m utils.koyfin.process <format>")
        print("Formats: snapshot, search, ticker, transcripts, screener")
        sys.exit(1)
    
    fmt = sys.argv[1]
    data = parse_cli_output(sys.stdin.read())
    
    if fmt == "snapshot":
        print(format_snapshot(data))
    elif fmt == "search":
        print(format_search_results(data))
    elif fmt == "ticker":
        print(format_ticker_data(data))
    elif fmt == "transcripts":
        print(format_transcripts(data))
    elif fmt == "screener":
        print(format_screener_results(data))
    else:
        print(f"Unknown format: {fmt}")
        sys.exit(1)
