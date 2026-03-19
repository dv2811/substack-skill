package main

import (
	"fmt"
	"os"

	"entext-applications/internal/koyfin"
)

func runScreenerSchema(client *koyfin.Client, session *koyfin.Session, args []string) {
	var assetType string
	fs := newFlagSet("screener-schema")
	fs.StringVar(&assetType, "asset-type", "Equity", "Asset type (currently only Equity is supported)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	req := koyfin.ScreenerSchemaRequest{
		AssetType: assetType,
	}

	filters := koyfin.GetScreenerSchema(req.AssetType)
	if filters == nil {
		fmt.Fprintf(os.Stderr, "Error: no screener schema found for the given asset type\n")
		os.Exit(1)
	}

	printJSON(map[string]any{"filters": filters})
}
