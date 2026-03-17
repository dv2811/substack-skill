#!/usr/bin/env python3
"""
Koyfin Snapshot to Excel Export

Converts koyfin snapshot CLI output to Excel with:
- Raw data sheet
- Ratios sheet (calculated metrics)
- Growth sheet (YoY, QoQ growth rates)
- Summary sheet
"""

import json
import sys
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional

import pandas as pd


# Metric categories for organization
METRIC_CATEGORIES = {
    "price": [
        "current_price_adj",
        "prev_1y_price_adj",
        "prev_3y_price_adj",
        "prev_5y_price_adj",
        "start_year_price_adj",
    ],
    "valuation": [
        "market_cap_mil",
        "current_ev",
        "pe_ratio",
        "ev_ebitda",
        "ev_sales",
    ],
    "profitability": [
        "ebitda_last_12m",
        "ebit_last_12m",
        "net_income_ltm",
        "eps_last_12m",
        "revenue_last_12m",
        "cashflow_last_12m",
    ],
    "estimates": [
        "ebitda_next_12m",
        "revenue_next_12m",
        "eps_next_12m",
        "ntm_eps_gaap",
    ],
    "metadata": [
        "name",
        "ticker",
        "sector",
        "industry",
    ],
}


def load_snapshot_data(input_path: Optional[str] = None) -> List[Dict[str, Any]]:
    """
    Load snapshot data from CLI output or file.
    
    Args:
        input_path: Path to JSON file, or None to read from stdin
    
    Returns:
        List of snapshot records
    """
    if input_path:
        with open(input_path, "r") as f:
            data = json.load(f)
    else:
        data = json.load(sys.stdin)
    
    return data.get("results", [])


def flatten_snapshot(records: List[Dict[str, Any]]) -> pd.DataFrame:
    """
    Flatten snapshot records into DataFrame.
    
    Args:
        records: List of snapshot records from CLI
    
    Returns:
        DataFrame with one row per KID
    """
    rows = []
    for record in records:
        kid = record.get("kids", "")
        metrics = record.get("metrics", {})

        row = {"KID": kid}
        for key, metric_data in metrics.items():
            # Extract value from metric object {"date": ..., "value": ..., "currency": ...}
            if isinstance(metric_data, dict):
                if "error" in metric_data:
                    # Skip metrics with errors
                    continue
                row[key] = metric_data.get("value")
            else:
                # Fallback for simple values
                row[key] = metric_data
        rows.append(row)

    return pd.DataFrame(rows)


def calculate_ratios(df: pd.DataFrame) -> pd.DataFrame:
    """
    Calculate financial ratios from raw data.
    
    Args:
        df: Raw snapshot DataFrame
    
    Returns:
        DataFrame with calculated ratios
    """
    ratios = pd.DataFrame({"KID": df["KID"]})
    
    # P/E Ratio (if not already provided)
    if "current_price_adj" in df.columns and "eps_last_12m" in df.columns:
        if "pe_ratio" not in df.columns:
            ratios["pe_ratio_calc"] = df["current_price_adj"] / df["eps_last_12m"].replace(0, pd.NA)
    
    # EV/EBITDA
    if "current_ev" in df.columns and "ebitda_last_12m" in df.columns:
        if "ev_ebitda" not in df.columns:
            ratios["ev_ebitda_calc"] = df["current_ev"] / df["ebitda_last_12m"].replace(0, pd.NA)
    
    # EV/Sales
    if "current_ev" in df.columns and "revenue_last_12m" in df.columns:
        ratios["ev_sales_calc"] = df["current_ev"] / df["revenue_last_12m"].replace(0, pd.NA)
    
    # EBITDA Margin
    if "ebitda_last_12m" in df.columns and "revenue_last_12m" in df.columns:
        ratios["ebitda_margin_calc"] = (df["ebitda_last_12m"] / df["revenue_last_12m"].replace(0, pd.NA)) * 100
    
    # Net Margin
    if "net_income_ltm" in df.columns and "revenue_last_12m" in df.columns:
        ratios["net_margin_calc"] = (df["net_income_ltm"] / df["revenue_last_12m"].replace(0, pd.NA)) * 100
    
    # FCF Margin (if cashflow available)
    if "cashflow_last_12m" in df.columns and "revenue_last_12m" in df.columns:
        ratios["fcf_margin_calc"] = (df["cashflow_last_12m"] / df["revenue_last_12m"].replace(0, pd.NA)) * 100
    
    # P/FCF
    if "current_price_adj" in df.columns and "cashflow_last_12m" in df.columns:
        ratios["p_fcf_calc"] = df["current_price_adj"] / df["cashflow_last_12m"].replace(0, pd.NA)
    
    return ratios


def calculate_growth(df: pd.DataFrame) -> pd.DataFrame:
    """
    Calculate growth metrics from snapshot data.
    
    Args:
        df: Raw snapshot DataFrame
    
    Returns:
        DataFrame with growth rates
    """
    growth = pd.DataFrame({"KID": df["KID"]})
    
    # 1-Year Price Growth
    if "current_price_adj" in df.columns and "prev_1y_price_adj" in df.columns:
        growth["price_1y_growth"] = (
            (df["current_price_adj"] - df["prev_1y_price_adj"]) / 
            df["prev_1y_price_adj"].replace(0, pd.NA) * 100
        )
    
    # 3-Year Price CAGR
    if "current_price_adj" in df.columns and "prev_3y_price_adj" in df.columns:
        growth["price_3y_cagr"] = (
            ((df["current_price_adj"] / df["prev_3y_price_adj"].replace(0, pd.NA)) ** (1/3) - 1) * 100
        )
    
    # 5-Year Price CAGR
    if "current_price_adj" in df.columns and "prev_5y_price_adj" in df.columns:
        growth["price_5y_cagr"] = (
            ((df["current_price_adj"] / df["prev_5y_price_adj"].replace(0, pd.NA)) ** (1/5) - 1) * 100
        )
    
    # YTD Price Growth
    if "current_price_adj" in df.columns and "start_year_price_adj" in df.columns:
        growth["price_ytd_growth"] = (
            (df["current_price_adj"] - df["start_year_price_adj"]) / 
            df["start_year_price_adj"].replace(0, pd.NA) * 100
        )
    
    # Estimate vs LTM growth (for forward-looking metrics)
    if "ebitda_next_12m" in df.columns and "ebitda_last_12m" in df.columns:
        growth["ebitda_growth_est"] = (
            (df["ebitda_next_12m"] - df["ebitda_last_12m"]) / 
            df["ebitda_last_12m"].replace(0, pd.NA) * 100
        )
    
    if "revenue_next_12m" in df.columns and "revenue_last_12m" in df.columns:
        growth["revenue_growth_est"] = (
            (df["revenue_next_12m"] - df["revenue_last_12m"]) / 
            df["revenue_last_12m"].replace(0, pd.NA) * 100
        )
    
    if "eps_next_12m" in df.columns and "eps_last_12m" in df.columns:
        growth["eps_growth_est"] = (
            (df["eps_next_12m"] - df["eps_last_12m"]) / 
            df["eps_last_12m"].replace(0, pd.NA) * 100
        )
    
    return growth


def create_summary(df: pd.DataFrame, ratios: pd.DataFrame, growth: pd.DataFrame) -> pd.DataFrame:
    """
    Create summary sheet with key metrics.
    
    Args:
        df: Raw snapshot DataFrame
        ratios: Calculated ratios DataFrame
        growth: Growth metrics DataFrame
    
    Returns:
        Summary DataFrame
    """
    summary = pd.DataFrame({"KID": df["KID"]})
    
    # Add ticker and name if available
    if "ticker" in df.columns:
        summary["Ticker"] = df["ticker"]
    if "name" in df.columns:
        summary["Name"] = df["name"]
    if "sector" in df.columns:
        summary["Sector"] = df["sector"]
    
    # Key price metrics
    if "current_price_adj" in df.columns:
        summary["Price"] = df["current_price_adj"]
    
    # Market cap
    if "market_cap_mil" in df.columns:
        summary["Market Cap (M)"] = df["market_cap_mil"]
    
    # Key ratios (prefer existing, fallback to calculated)
    if "pe_ratio" in df.columns:
        summary["P/E"] = df["pe_ratio"]
    elif "pe_ratio_calc" in ratios.columns:
        summary["P/E"] = ratios["pe_ratio_calc"]
    
    if "ev_ebitda" in df.columns:
        summary["EV/EBITDA"] = df["ev_ebitda"]
    elif "ev_ebitda_calc" in ratios.columns:
        summary["EV/EBITDA"] = ratios["ev_ebitda_calc"]
    
    # Margins
    if "ebitda_margin_calc" in ratios.columns:
        summary["EBITDA Margin %"] = ratios["ebitda_margin_calc"]
    
    # Growth rates
    if "price_1y_growth" in growth.columns:
        summary["Price 1Y Growth %"] = growth["price_1y_growth"]
    
    if "ebitda_growth_est" in growth.columns:
        summary["EBITDA Growth Est %"] = growth["ebitda_growth_est"]
    
    return summary


def export_to_excel(
    records: List[Dict[str, Any]],
    output_path: str,
    include_raw: bool = True,
    include_ratios: bool = True,
    include_growth: bool = True,
    include_summary: bool = True,
) -> str:
    """
    Export snapshot data to Excel with multiple sheets.
    
    Args:
        records: Snapshot records from CLI
        output_path: Output Excel file path
        include_raw: Include raw data sheet
        include_ratios: Include ratios sheet
        include_growth: Include growth sheet
        include_summary: Include summary sheet
    
    Returns:
        Path to created file
    """
    df = flatten_snapshot(records)
    
    with pd.ExcelWriter(output_path, engine="openpyxl") as writer:
        if include_summary:
            ratios = calculate_ratios(df) if include_ratios else pd.DataFrame({"KID": df["KID"]})
            growth = calculate_growth(df) if include_growth else pd.DataFrame({"KID": df["KID"]})
            summary = create_summary(df, ratios, growth)
            summary.to_excel(writer, sheet_name="Summary", index=False)
        
        if include_raw:
            df.to_excel(writer, sheet_name="Raw Data", index=False)
        
        if include_ratios:
            ratios = calculate_ratios(df)
            ratios.to_excel(writer, sheet_name="Ratios", index=False)
        
        if include_growth:
            growth = calculate_growth(df)
            growth.to_excel(writer, sheet_name="Growth", index=False)
    
    return output_path


def main():
    """Main entry point for CLI usage."""
    import argparse
    
    parser = argparse.ArgumentParser(
        description="Convert koyfin snapshot to Excel"
    )
    parser.add_argument(
        "input",
        nargs="?",
        help="Input JSON file (default: stdin)"
    )
    parser.add_argument(
        "-o", "--output",
        required=True,
        help="Output Excel file path"
    )
    parser.add_argument(
        "--no-raw",
        action="store_true",
        help="Exclude raw data sheet"
    )
    parser.add_argument(
        "--no-ratios",
        action="store_true",
        help="Exclude ratios sheet"
    )
    parser.add_argument(
        "--no-growth",
        action="store_true",
        help="Exclude growth sheet"
    )
    parser.add_argument(
        "--no-summary",
        action="store_true",
        help="Exclude summary sheet"
    )
    
    args = parser.parse_args()
    
    # Load data
    records = load_snapshot_data(args.input)
    
    if not records:
        print("Error: No snapshot data found", file=sys.stderr)
        sys.exit(1)
    
    # Export to Excel
    output_path = export_to_excel(
        records,
        args.output,
        include_raw=not args.no_raw,
        include_ratios=not args.no_ratios,
        include_growth=not args.no_growth,
        include_summary=not args.no_summary,
    )
    
    print(f"Exported {len(records)} records to {output_path}")
    
    # Print summary
    df = flatten_snapshot(records)
    print(f"\nSheets created:")
    if not args.no_summary:
        print("  - Summary: Key metrics and ratios")
    if not args.no_raw:
        print("  - Raw Data: All snapshot metrics")
    if not args.no_ratios:
        print("  - Ratios: Calculated financial ratios")
    if not args.no_growth:
        print("  - Growth: YoY and estimated growth rates")


if __name__ == "__main__":
    main()
