package main

import (
	"fmt"
	"os"

	"entext-applications/internal/substack"
)

const SearchCmdHelp = `
Search Substack posts with different modes.
Usage: substack search [flags]
Flags:
	-query string
		Search query (required)
	-mode string
		Search mode: top, all, subscribed (default "all")
	-page int
		Page number for pagination (0-10, not used for top mode)
	-language string
		Language code (2-letter, e.g., 'en')
Examples:
	substack search -query "AI" -mode top
	substack search -query "technology" -mode all -page 1
	substack search -query "newsletter" -mode subscribed -language en`

// Implementation for search CLI command
func runSearch(client *substack.Client, session *substack.Session, args []string) {
	// check valid session before authenticate
	checkValidSession(session)

	var (
		query, mode, language string
		page int
	)
	fs := newFlagSet("search")
	fs.StringVar(&query, "query", "", "Search query (required)")
	fs.StringVar(&mode, "mode", "all", "Search mode: top, all, subscribed")
	fs.IntVar(&page, "page", 0, "Page number for pagination (0-10, not used for top mode)")
	fs.StringVar(&language, "language", "", "Language code (2-letter, e.g., 'en')")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if query == "" {
		fmt.Fprintf(os.Stderr, "Error: -query flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	req := substack.SearchRequest{
		Query: query,
		Mode:  substack.SbkSearchFeedMode(mode),
		Lang:  language,
		Page:  page,
	}

	results, err := client.SearchPosts(session, req)
	if err != nil {
		exitWithError(err)
	}
	printJSON(map[string]any{"data": results})
}
