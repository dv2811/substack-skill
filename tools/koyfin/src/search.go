package main

import (
	"fmt"
	"os"
	"strings"

	"entext-applications/internal/koyfin"
)

func runSearch(client *koyfin.Client, session *koyfin.Session, args []string) {
	var (
		query, categories string
		primaryOnly bool
	)
	fs := newFlagSet("search")
	fs.StringVar(&query,"q", "", "Ticker or ETF name to search for (required)")
	fs.StringVar(&categories, "categories", "Equity,ETF", "Search categories (comma-separated)")
	fs.BoolVar(&primaryOnly, "primary-only", false, "Use primary exchange only")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if query == "" {
		fmt.Fprintf(os.Stderr, "Error: -q flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	// Parse categories
	cats := strings.Split(categories, ",")

	req := koyfin.SearchRequest{
		SearchString: query,
		Categories:   cats,
		PrimaryOnly:  primaryOnly,
	}

	results, err := client.LookUpByName(session, req)
	if err != nil {
		exitWithError(err)
	}

	printJSON(map[string]any{"data": results})
}
