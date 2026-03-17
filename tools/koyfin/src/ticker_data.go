package main

import (
	"fmt"
	"os"

	"entext-applications/internal/koyfin"
)

func runTickerData(client *koyfin.Client, session *koyfin.Session, args []string) {
	fs := newFlagSet("ticker-data")
	id := fs.String("kid", "", "Koyfin ID for the ticker (required)")
	key := fs.String("key", "", "Indicator key to search for (required)")
	dateFrom := fs.String("date-from", "", "Start date in YYYY-MM-DD format (required)")
	dateTo := fs.String("date-to", "", "End date in YYYY-MM-DD format (optional, defaults to today)")
	currency := fs.String("currency", "USD", "Data currency (default: USD)")
	aggPeriod := fs.String("agg-period", "day", "Series granularity: day, monthly, quarterly, annually")
	priceFormat := fs.String("price-format", "", "Price format: both, standard, adj (auto-set based on key)")
	finPeriod := fs.String("fin-period", "", "Financial period: quarterly, annual, LTM")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *id == "" {
		fmt.Fprintf(os.Stderr, "Error: -id flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	if *key == "" {
		fmt.Fprintf(os.Stderr, "Error: -key flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	if *dateFrom == "" {
		fmt.Fprintf(os.Stderr, "Error: -date-from flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	req := koyfin.TickerDataRequest{
		ID:        *id,
		Key:       *key,
		Currency:  *currency,
		DateFrom:  *dateFrom,
		DateTo:    *dateTo,
		AggPeriod: *aggPeriod,
		FinPeriod: *finPeriod,
	}

	// Set price format if explicitly provided
	if *priceFormat != "" {
		req.PriceFormat = *priceFormat
	}

	data, err := client.GetDataSeries(session, req)
	if err != nil {
		exitWithError(err)
	}

	printJSON(map[string]any{"data": data})
}
