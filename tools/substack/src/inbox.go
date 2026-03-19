package main

import (
	"entext-applications/internal/substack"
	"fmt"
	"os"
)

const InboxCmdHelp = `
Get chronological inbox posts.
Usage: substack inbox [flags]
Flags:
	-after: string. Timestamp cursor for pagination (RFC3339 format)
Examples:
	substack inbox
	substack inbox -after "2024-01-01T00:00:00.000Z"`

// command line implementation for list Substack inbox feed
func runInbox(client *substack.Client, session *substack.Session, args []string) {
	// check valid session before authenticate
	checkValidSession(session)

	fs := newFlagSet("inbox")
	after := fs.String("after", "", "Timestamp cursor for pagination (RFC3339 format)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	inbox, err := client.GetChronologicalInbox(session, *after)
	if err != nil {
		exitWithError(err)
	}
	printJSON(map[string]any{"data": inbox})
}
