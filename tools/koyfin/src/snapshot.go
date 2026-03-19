package main

import (
	"fmt"
	"os"
	"strings"

	"entext-applications/internal/koyfin"
)

func runSnapshot(client *koyfin.Client, session *koyfin.Session, args []string) {
	var kids, category string
	fs := newFlagSet("snapshot")
	fs.StringVar(&kids, "kids", "", "Comma-separated list of Koyfin IDs (required, max 32 for Equity, 2 for ETF)")
	fs.StringVar(&category, "category", "Equity", "Category of the instrument (Equity or ETF)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if kids == "" {
		fmt.Fprintf(os.Stderr, "Error: -kids flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	req := koyfin.SnapshotRequest{
		KIDs:     strings.Split(kids, ","),
		Category: category,
	}

	data, err := client.GetSnapshotData(session, req)
	if err != nil {
		exitWithError(err)
	}

	printJSON(map[string]any{"results": data})
}
