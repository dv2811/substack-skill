package main

import (
	"fmt"
	"os"

	"entext-applications/internal/koyfin"
)

func runSchema(client *koyfin.Client, session *koyfin.Session, args []string) {
	fs := newFlagSet("schema")
	assetType := fs.String("asset-type", "Equity", "Asset type (currently only Equity is supported)")
	indicatorType := fs.String("indicator-type", "", "Indicator type: financials, ratios, forward_estimates, market_data (required for Equity)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *indicatorType == "" {
		fmt.Fprintf(os.Stderr, "Error: -indicator-type flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	req := koyfin.SchemaRequest{
		AssetType:     *assetType,
		IndicatorType: *indicatorType,
	}

	schemas := koyfin.GetAvailableSchema(&req)
	if schemas == nil {
		fmt.Fprintf(os.Stderr, "Error: no schema found for the given asset type and indicator type\n")
		os.Exit(1)
	}

	printJSON(map[string]any{"schemas": schemas})
}
