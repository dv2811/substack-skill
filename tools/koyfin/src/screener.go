package main

import (
	"encoding/json"
	"fmt"
	"os"

	"entext-applications/internal/koyfin"
	"entext-applications/internal/utils"
)

func runScreener(client *koyfin.Client, session *koyfin.Session, args []string) {
	fs := newFlagSet("screener")
	filtersJSON := fs.String("filters", "", "JSON array of filter conditions (required)")
	pageSize := fs.Uint("page-size", 100, "Page size (max 300)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *filtersJSON == "" {
		fmt.Fprintf(os.Stderr, "Error: -filters flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	// Parse filters from JSON
	var filters []koyfin.Filter
	if err := json.Unmarshal([]byte(*filtersJSON), &filters); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing filters JSON: %v\n", err)
		os.Exit(1)
	}

	req := koyfin.ScreenCriteria{
		Conditions: filters,
		OrderBy:    koyfin.CreateOrder("mkt", "DESC"),
		PageSize:   uint32(*pageSize),
	}

	// Add default primary filter
	req.Conditions = append(req.Conditions, koyfin.DefaultPrimaryFilter)

	// Handle empty value ranges
	for i := range req.Conditions {
		f := req.Conditions[i]
		if len(f.Values) > 0 {
			continue
		}
		if f.Max != nil && f.Min == nil {
			req.Conditions[i].Min = utils.Ptr(koyfin.MinFacetValue)
		} else if f.Max == nil && f.Min != nil {
			req.Conditions[i].Max = utils.Ptr(koyfin.MaxFacetValue)
		}
	}

	result, err := client.ScreenForStocks(session, req)
	if err != nil {
		exitWithError(err)
	}

	printJSON(map[string]any{"kids": result.Kids})
}
