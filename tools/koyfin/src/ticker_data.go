package main

import (
	"fmt"
	"os"

	"entext-applications/internal/koyfin"
	"entext-applications/internal/validator"
)

func runTickerData(client *koyfin.Client, session *koyfin.Session, args []string) {
	var (
		id, key, dateFrom, dateTo, aggPeriod, priceFormat, finPeriod, currency string
	)
	fs := newFlagSet("ticker-data")
	fs.StringVar(&id, "kid", "", "Koyfin ID for the ticker (required)")
	fs.StringVar(&key, "key", "", "Indicator key to search for (required)")
	fs.StringVar(&dateFrom, "date-from", "", "Start date in YYYY-MM-DD format (required)")
	fs.StringVar(&dateTo, "date-to", "", "End date in YYYY-MM-DD format (optional, defaults to today)")
	fs.StringVar(&currency, "currency", "USD", "Data currency (default: USD)")
	fs.StringVar(&aggPeriod, "agg-period", "day", "Series granularity: day, monthly, quarterly, annually")
	fs.StringVar(&priceFormat, "price-format", "", "Price format: both, standard, adj (auto-set based on key)")
	fs.StringVar(&finPeriod, "fin-period", "", "Financial period: quarterly, annual, LTM")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	v := validator.New()
	v.Check(id != "", "id", "id must not be empty")
	v.Check(key != "", "key", "key must not be empty")
	v.Check(dateFrom != "", "dateFrom", "dateFrom must not be empty")

	if !v.Valid() {
		for k, v := range v.Errors {
			fmt.Fprintf(os.Stderr, "%s: %s", k, v)
		}
	}

	req := koyfin.TickerDataRequest{
		ID:        id,
		Key:       key,
		Currency:  currency,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		AggPeriod: aggPeriod,
		FinPeriod: finPeriod,
	}

	// Set price format if explicitly provided
	if priceFormat != "" {
		req.PriceFormat = priceFormat
	}

	data, err := client.GetDataSeries(session, req)
	if err != nil {
		exitWithError(err)
	}

	printJSON(map[string]any{"data": data})
}
