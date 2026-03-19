package main

import (
	"fmt"
	"os"
	"strconv"

	"entext-applications/internal/substack"
)

const articleHelp = `
Get article content by post ID.

Usage: substack article [flags]

Flags:
  -post-id uint
        Post ID to retrieve content for (required)
  -base-url string
        Custom Substack domain (optional, e.g., 'substack.com')
  -compact-mode boolean
        compact content mode. (optional, default: false)

Examples:
  substack article -post-id 123456
  substack article -post-id 123456 -base-url "substack.com"`

func runArticle(client *substack.Client, session *substack.Session, args []string) {
	// check valid session before authenticate
	checkValidSession(session)

	fs := newFlagSet("article")
	postID := fs.Uint("post-id", 0, "Post ID to retrieve content for (required)")
	baseURL := fs.String("base-url", "", "Custom Substack domain (optional, e.g., 'substack.com')")
	compactMode := fs.Bool("compact-mode", false, "compact content mode")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *postID == 0 {
		fmt.Fprintf(os.Stderr, "Error: -post-id flag is required\n")
		fs.Usage()
		os.Exit(1)
	}

	article, err := client.GetArticleContent(session, &substack.PostContentRequest{
		BaseURL: *baseURL,
		PostID:  uint32(*postID),
		CmpMode: *compactMode,
	})
	if err != nil {
		exitWithError(err)
	}

	printJSON(map[string]any{"data": article})
}